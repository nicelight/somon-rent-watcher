package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
	"github.com/nicelight/somon-rent-watcher/internal/store"
	"github.com/nicelight/somon-rent-watcher/internal/telegram"
)

const (
	stateInitialized       = "initialized"
	stateLastSuccessful    = "last_successful_poll_at"
	statePreviousOrdinary  = "previous_ordinary_ids"
	stateTelegramOffset    = "telegram_offset"
	maxDebugFiles          = 20
	adminNotificationPause = time.Hour
)

type processStats struct {
	NewIDs         int
	DetailRequests int
	Sent           int
}

type App struct {
	cfg    config.Config
	store  *store.DB
	somon  *somon.Client
	bot    *telegram.Bot
	logger *slog.Logger

	statusMu sync.RWMutex
	status   model.RuntimeStatus

	randMu sync.Mutex
	rand   *rand.Rand

	notifyMu   sync.Mutex
	lastNotify map[string]time.Time
}

func New(cfg config.Config, db *store.DB, logger *slog.Logger) (*App, error) {
	if db == nil {
		return nil, errors.New("nil store")
	}
	if logger == nil {
		logger = slog.Default()
	}
	a := &App{
		cfg:        cfg,
		store:      db,
		somon:      somon.NewClient(cfg.UserAgent, cfg.RequestDelay, cfg.HTTPTimeout, cfg.MaxBodyBytes),
		logger:     logger,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		lastNotify: make(map[string]time.Time),
	}
	if _, err := a.LoadSettings(); err != nil {
		return nil, fmt.Errorf("initialize settings: %w", err)
	}
	if count, err := db.CountSeen(); err == nil {
		a.status.SeenCount = count
	} else {
		return nil, fmt.Errorf("count seen ads: %w", err)
	}
	if last, ok, err := db.GetState(stateLastSuccessful); err != nil {
		return nil, err
	} else if ok {
		if parsed, err := time.Parse(time.RFC3339Nano, last); err == nil {
			a.status.LastSuccessfulPoll = parsed
		}
	}
	a.status.Mode = "запуск"

	tgClient := telegram.NewClient(cfg.TelegramAPIBase, cfg.TelegramBotToken)
	a.bot = telegram.NewBot(tgClient, a, cfg.TelegramAdminUserIDs, cfg.TelegramTargetChatID, logger)
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)
	go func() { errCh <- a.bot.Run(runCtx) }()
	go func() { errCh <- a.pollLoop(runCtx) }()

	select {
	case <-ctx.Done():
		cancel()
		return nil
	case err := <-errCh:
		cancel()
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
}

// Backend implementation used by the Telegram package.
func (a *App) LoadSettings() (filter.Settings, error) {
	raw, ok, err := a.store.LoadSettingsJSON()
	if err != nil {
		return filter.Settings{}, err
	}
	if !ok {
		settings := filter.DefaultSettings()
		encoded, err := settings.Encode()
		if err != nil {
			return filter.Settings{}, err
		}
		if err := a.store.SaveSettingsJSON(encoded); err != nil {
			return filter.Settings{}, err
		}
		return settings, nil
	}
	return filter.DecodeSettings(raw)
}

func (a *App) SaveSettings(settings filter.Settings) error {
	encoded, err := settings.Encode()
	if err != nil {
		return err
	}
	if err := a.store.SaveSettingsJSON(encoded); err != nil {
		return err
	}
	a.logger.Info("filter settings updated")
	return nil
}

func (a *App) RuntimeStatus() model.RuntimeStatus {
	a.statusMu.RLock()
	status := a.status
	a.statusMu.RUnlock()
	if count, err := a.store.CountSeen(); err == nil {
		status.SeenCount = count
	}
	return status
}

func (a *App) TelegramOffset() (int64, error) {
	offset, ok, err := a.store.GetStateInt64(stateTelegramOffset)
	if err != nil || !ok {
		return 0, err
	}
	return offset, nil
}

func (a *App) SetTelegramOffset(offset int64) error {
	return a.store.SetState(stateTelegramOffset, strconv.FormatInt(offset, 10))
}

func (a *App) pollLoop(ctx context.Context) error {
	first := true
	for {
		if !first {
			a.statusMu.RLock()
			next := a.status.NextPoll
			a.statusMu.RUnlock()
			wait := time.Until(next)
			if wait > 0 && !sleepContext(ctx, wait) {
				return ctx.Err()
			}
		}
		first = false

		a.setStatusMode("опрос", "")
		err := a.pollOnce(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}

		delay := a.randomPollDelay()
		mode := "норма"
		lastError := ""
		if err == nil {
			if settings, loadErr := a.LoadSettings(); loadErr == nil && !settings.Enabled {
				mode = "пауза; baseline обновляется"
			}
		}
		if err != nil {
			lastError = err.Error()
			if status, retryAfter, blocked := somon.IsBlocked(err); blocked {
				delay = a.cfg.BlockBackoff
				if retryAfter > delay {
					delay = retryAfter
				}
				mode = fmt.Sprintf("backoff HTTP %d", status)
				a.notifyAdmin(ctx, "blocked", fmt.Sprintf("Somon вернул HTTP %d. Опрос приостановлен на %s; обход блокировки не выполняется.", status, delay), adminNotificationPause)
			} else {
				mode = "ошибка; следующий обычный опрос"
				a.logger.Error("Somon poll failed", "error", err)
			}
		}
		nextPoll := time.Now().Add(delay)
		a.statusMu.Lock()
		a.status.Mode = mode
		a.status.LastError = lastError
		a.status.NextPoll = nextPoll
		if strings.HasPrefix(mode, "backoff") {
			a.status.BackoffUntil = nextPoll
		} else {
			a.status.BackoffUntil = time.Time{}
		}
		a.statusMu.Unlock()
	}
}

func (a *App) pollOnce(ctx context.Context) error {
	now := time.Now().UTC()
	cards, body, err := a.somon.FetchCategory(ctx, a.cfg.CategoryURL)
	if err != nil {
		if len(body) > 0 {
			a.saveDebug("category-parse-error", body)
		}
		return err
	}
	if err := a.validateCategory(cards); err != nil {
		a.saveDebug("category-sanity-error", body)
		a.notifyAdmin(ctx, "category-parser", "Парсер основной страницы Somon не прошёл sanity-check: "+err.Error()+". База seen не изменена.", adminNotificationPause)
		return err
	}

	initialized, err := a.isInitialized()
	if err != nil {
		return err
	}
	ordinary := ordinaryIDs(cards)
	if !initialized {
		ids := cardIDs(cards)
		if err := a.store.MarkSeen(ids, now); err != nil {
			return fmt.Errorf("create baseline: %w", err)
		}
		if err := a.saveSuccessfulPoll(now, ordinary, true); err != nil {
			return err
		}
		a.refreshSeenCount()
		a.setCycleStats(len(cards), 0, 0)
		a.logger.Info("baseline created", "cards", len(cards), "ordinary", len(ordinary))
		a.notifyAdmin(ctx, "baseline", fmt.Sprintf("Somon Watcher запущен. Baseline создан из %d текущих объявлений; они не отправлены. Мониторинг по умолчанию на паузе: настройте /filter и нажмите «Включить мониторинг». Пока он на паузе, новые ID запоминаются без отправки.", len(cards)), 0)
		return nil
	}

	settings, err := a.LoadSettings()
	if err != nil {
		return fmt.Errorf("load filter: %w", err)
	}
	if !settings.Enabled {
		seen, err := a.store.SeenIDs(cardIDs(cards))
		if err != nil {
			return fmt.Errorf("query seen IDs while paused: %w", err)
		}
		newIDs := make([]int64, 0)
		for _, card := range cards {
			if card.ID > 0 && !seen[card.ID] {
				newIDs = append(newIDs, card.ID)
			}
		}
		if err := a.store.MarkSeen(newIDs, now); err != nil {
			return fmt.Errorf("advance fresh baseline while paused: %w", err)
		}
		finished := time.Now().UTC()
		if err := a.saveSuccessfulPoll(finished, ordinary, false); err != nil {
			return err
		}
		a.refreshSeenCount()
		a.setCycleStats(len(cards), len(newIDs), 0)
		a.logger.Info("monitoring paused; fresh baseline advanced", "cards", len(cards), "new_ids", len(newIDs))
		return nil
	}

	previous, err := a.loadPreviousOrdinary()
	if err != nil {
		return err
	}
	lastSuccessful := a.lastSuccessfulPoll()
	gapReason := continuityGapReason(previous, ordinary, lastSuccessful, now, a.cfg.GapAfter)
	if gapReason != "" {
		a.logger.Warn("possible category gap", "reason", gapReason)
		a.notifyAdmin(ctx, "gap", "Возможен разрыв свежей ленты Somon ("+gapReason+"). Выполняется один recovery-опрос выбранных комнат.", adminNotificationPause)
		recovered, recErr := a.recoverCards(ctx, settings, now, lastSuccessful)
		if recErr != nil {
			if _, _, blocked := somon.IsBlocked(recErr); blocked {
				return recErr
			}
			a.logger.Error("recovery sweep partially failed", "error", recErr)
		}
		cards = mergeCards(cards, recovered)
		a.logger.Info("recovery sweep completed", "recovered_candidates", len(recovered))
	}

	stats, err := a.processNewCards(ctx, settings, cards, now)
	if err != nil {
		return err
	}
	finished := time.Now().UTC()
	if err := a.saveSuccessfulPoll(finished, ordinary, false); err != nil {
		return err
	}
	a.refreshSeenCount()
	a.setCycleStats(len(cards), stats.NewIDs, stats.Sent)
	a.logger.Info("Somon poll completed", "cards", len(cards), "ordinary", len(ordinary), "new_ids", stats.NewIDs, "details", stats.DetailRequests, "sent", stats.Sent)
	return nil
}

func (a *App) validateCategory(cards []model.Card) error {
	if len(cards) < a.cfg.MinCards {
		return fmt.Errorf("found %d cards, expected at least %d", len(cards), a.cfg.MinCards)
	}
	valid := 0
	ordinary := 0
	seen := make(map[int64]struct{}, len(cards))
	for _, card := range cards {
		if card.ID > 0 && card.URL != "" && card.Title != "" {
			valid++
		}
		if !card.Promoted {
			ordinary++
		}
		if _, exists := seen[card.ID]; exists {
			return fmt.Errorf("duplicate card ID %d after parser dedupe", card.ID)
		}
		seen[card.ID] = struct{}{}
	}
	if valid < a.cfg.MinCards {
		return fmt.Errorf("only %d structurally valid cards", valid)
	}
	if ordinary == 0 {
		return errors.New("ordinary section is empty")
	}
	return nil
}

func (a *App) processNewCards(ctx context.Context, settings filter.Settings, cards []model.Card, now time.Time) (processStats, error) {
	cards = prioritizeOrdinary(cards)
	ids := cardIDs(cards)
	seen, err := a.store.SeenIDs(ids)
	if err != nil {
		return processStats{}, fmt.Errorf("query seen IDs: %w", err)
	}

	unseenCards := make([]model.Card, 0, len(cards))
	for _, card := range cards {
		if card.ID > 0 && !seen[card.ID] {
			unseenCards = append(unseenCards, card)
		}
	}
	stats := processStats{NewIDs: len(unseenCards)}
	for _, card := range unseenCards {
		if ok, reason := filter.CardMatches(settings, card); !ok {
			a.logger.Debug("new card rejected by prefilter", "ad_id", card.ID, "reason", reason)
			if err := a.store.MarkSeen([]int64{card.ID}, now); err != nil {
				return stats, err
			}
			continue
		}

		if stats.DetailRequests >= a.cfg.MaxDetailsPerPoll {
			a.logger.Warn("detail request cap reached; remaining IDs stay unseen", "cap", a.cfg.MaxDetailsPerPoll)
			break
		}
		stats.DetailRequests++
		ad, detailBody, err := a.somon.FetchDetail(ctx, card, a.cfg.CategoryURL)
		if err != nil {
			var httpErr *somon.HTTPError
			if errors.As(err, &httpErr) {
				switch httpErr.StatusCode {
				case http.StatusNotFound, http.StatusGone:
					if markErr := a.store.MarkSeen([]int64{card.ID}, now); markErr != nil {
						return stats, markErr
					}
					continue
				case http.StatusForbidden, http.StatusTooManyRequests:
					return stats, err
				}
			}
			if len(detailBody) > 0 {
				a.saveDebug(fmt.Sprintf("detail-%d-parse-error", card.ID), detailBody)
				a.notifyAdmin(ctx, "detail-parser", fmt.Sprintf("Не удалось разобрать detail-page Somon для ID %d. ID не помечен просмотренным и будет повторён, пока остаётся в ленте.", card.ID), adminNotificationPause)
			}
			a.logger.Error("detail fetch/parse failed; ad remains unseen", "ad_id", card.ID, "error", err)
			continue
		}

		if ok, reason := filter.AdMatches(settings, ad); !ok {
			a.logger.Debug("new ad rejected by detail filter", "ad_id", ad.ID, "reason", reason)
			if err := a.store.MarkSeen([]int64{card.ID}, now); err != nil {
				return stats, err
			}
			continue
		}
		if err := a.bot.SendAd(ctx, ad); err != nil {
			a.logger.Error("Telegram delivery failed; ad remains unseen", "ad_id", ad.ID, "error", err)
			continue
		}
		if err := a.store.MarkSeen([]int64{card.ID}, now); err != nil {
			return stats, fmt.Errorf("mark sent ad %d seen: %w", card.ID, err)
		}
		stats.Sent++
		a.logger.Info("new ad sent", "ad_id", ad.ID, "price", pointerValue(ad.Price), "seller_ads", pointerValue(ad.SellerAds))
	}
	if stats.NewIDs > 0 {
		a.logger.Info("new IDs processed", "count", stats.NewIDs)
	}
	return stats, nil
}

func (a *App) recoverCards(ctx context.Context, settings filter.Settings, now, lastSuccessful time.Time) ([]model.Card, error) {
	urls := somon.RecoveryURLs(a.cfg.CategoryURL, settings.Rooms)
	if len(urls) == 0 {
		return nil, nil
	}
	window := a.cfg.GapAfter + a.cfg.PollMax
	if !lastSuccessful.IsZero() {
		window = now.Sub(lastSuccessful) + a.cfg.PollMax
	}
	if window < a.cfg.PollMax {
		window = a.cfg.PollMax
	}

	byID := make(map[int64]model.Card)
	var errs []string
	for _, pageURL := range urls {
		cards, body, err := a.somon.FetchCategory(ctx, pageURL)
		if err != nil {
			if _, _, blocked := somon.IsBlocked(err); blocked {
				return nil, err
			}
			if len(body) > 0 {
				a.saveDebug("recovery-parse-error", body)
			}
			errs = append(errs, pageURL+": "+err.Error())
			continue
		}
		for _, card := range cards {
			age, ok := somon.ParseAgeDuration(card.AgeText)
			if !ok || age > window {
				continue
			}
			if existing, ok := byID[card.ID]; ok {
				byID[card.ID] = preferCard(existing, card)
			} else {
				byID[card.ID] = card
			}
		}
	}
	ids := make([]int64, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	cards := make([]model.Card, 0, len(ids))
	for _, id := range ids {
		cards = append(cards, byID[id])
	}
	if len(errs) > 0 {
		return cards, errors.New(strings.Join(errs, "; "))
	}
	return cards, nil
}

func (a *App) isInitialized() (bool, error) {
	value, ok, err := a.store.GetState(stateInitialized)
	return ok && value == "1", err
}

func (a *App) saveSuccessfulPoll(now time.Time, ordinary []int64, initialized bool) error {
	data, err := json.Marshal(ordinary)
	if err != nil {
		return err
	}
	values := map[string]string{
		stateLastSuccessful:   now.UTC().Format(time.RFC3339Nano),
		statePreviousOrdinary: string(data),
	}
	if initialized {
		values[stateInitialized] = "1"
	}
	if err := a.store.SetStates(values); err != nil {
		return fmt.Errorf("save successful poll state: %w", err)
	}
	a.statusMu.Lock()
	a.status.LastSuccessfulPoll = now
	a.status.LastError = ""
	a.statusMu.Unlock()
	return nil
}

func (a *App) loadPreviousOrdinary() ([]int64, error) {
	raw, ok, err := a.store.GetState(statePreviousOrdinary)
	if err != nil || !ok || strings.TrimSpace(raw) == "" {
		return nil, err
	}
	var ids []int64
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, fmt.Errorf("decode previous ordinary IDs: %w", err)
	}
	return ids, nil
}

func (a *App) lastSuccessfulPoll() time.Time {
	a.statusMu.RLock()
	defer a.statusMu.RUnlock()
	return a.status.LastSuccessfulPoll
}

func (a *App) setStatusMode(mode, lastError string) {
	a.statusMu.Lock()
	a.status.Mode = mode
	a.status.LastError = lastError
	a.status.NextPoll = time.Time{}
	a.statusMu.Unlock()
}

func (a *App) refreshSeenCount() {
	if count, err := a.store.CountSeen(); err == nil {
		a.statusMu.Lock()
		a.status.SeenCount = count
		a.statusMu.Unlock()
	}
}

func (a *App) setCycleStats(cards, newIDs, sent int) {
	a.statusMu.Lock()
	a.status.LastCardCount = cards
	a.status.LastNewCount = newIDs
	a.status.LastSentCount = sent
	a.statusMu.Unlock()
}

func prioritizeOrdinary(cards []model.Card) []model.Card {
	out := append([]model.Card(nil), cards...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Promoted == out[j].Promoted {
			return out[i].Position < out[j].Position
		}
		return !out[i].Promoted
	})
	return out
}

func (a *App) randomPollDelay() time.Duration {
	if a.cfg.PollMax <= a.cfg.PollMin {
		return a.cfg.PollMin
	}
	delta := a.cfg.PollMax - a.cfg.PollMin
	a.randMu.Lock()
	n := a.rand.Int63n(int64(delta) + 1)
	a.randMu.Unlock()
	return a.cfg.PollMin + time.Duration(n)
}

func (a *App) notifyAdmin(ctx context.Context, key, message string, minInterval time.Duration) {
	now := time.Now()
	a.notifyMu.Lock()
	last := a.lastNotify[key]
	if minInterval > 0 && !last.IsZero() && now.Sub(last) < minInterval {
		a.notifyMu.Unlock()
		return
	}
	a.lastNotify[key] = now
	a.notifyMu.Unlock()
	if err := a.bot.SendAdmin(ctx, message); err != nil {
		a.logger.Error("admin notification failed", "key", key, "error", err)
	}
}

func (a *App) saveDebug(prefix string, body []byte) {
	if len(body) == 0 || strings.TrimSpace(a.cfg.DebugDir) == "" {
		return
	}
	if err := os.MkdirAll(a.cfg.DebugDir, 0o700); err != nil {
		a.logger.Error("create debug directory failed", "error", err)
		return
	}
	prefix = sanitizeFilePart(prefix)
	name := fmt.Sprintf("%s-%s.html", time.Now().UTC().Format("20060102T150405.000000000Z"), prefix)
	path := filepath.Join(a.cfg.DebugDir, name)
	if err := os.WriteFile(path, body, 0o600); err != nil {
		a.logger.Error("write debug HTML failed", "path", path, "error", err)
		return
	}
	entries, err := os.ReadDir(a.cfg.DebugDir)
	if err != nil {
		return
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	for len(files) > maxDebugFiles {
		_ = os.Remove(filepath.Join(a.cfg.DebugDir, files[0]))
		files = files[1:]
	}
}

func continuityGapReason(previous, current []int64, lastSuccessful, now time.Time, gapAfter time.Duration) string {
	var reasons []string
	if len(previous) > 0 && len(current) > 0 && !hasIntersection(previous, current) {
		reasons = append(reasons, "нет пересечения обычных ID")
	}
	if !lastSuccessful.IsZero() && now.Sub(lastSuccessful) > gapAfter {
		reasons = append(reasons, "последний успешный опрос был "+now.Sub(lastSuccessful).Round(time.Minute).String()+" назад")
	}
	return strings.Join(reasons, "; ")
}

func hasIntersection(a, b []int64) bool {
	set := make(map[int64]struct{}, len(a))
	for _, id := range a {
		set[id] = struct{}{}
	}
	for _, id := range b {
		if _, ok := set[id]; ok {
			return true
		}
	}
	return false
}

func ordinaryIDs(cards []model.Card) []int64 {
	ids := make([]int64, 0, len(cards))
	for _, card := range cards {
		if !card.Promoted && card.ID > 0 {
			ids = append(ids, card.ID)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func cardIDs(cards []model.Card) []int64 {
	ids := make([]int64, 0, len(cards))
	for _, card := range cards {
		if card.ID > 0 {
			ids = append(ids, card.ID)
		}
	}
	return ids
}

func mergeCards(primary, secondary []model.Card) []model.Card {
	seen := make(map[int64]struct{}, len(primary)+len(secondary))
	out := make([]model.Card, 0, len(primary)+len(secondary))
	for _, card := range append(append([]model.Card(nil), primary...), secondary...) {
		if card.ID <= 0 {
			continue
		}
		if _, ok := seen[card.ID]; ok {
			continue
		}
		seen[card.ID] = struct{}{}
		out = append(out, card)
	}
	return out
}

func preferCard(a, b model.Card) model.Card {
	if a.Title == "" || len(b.Title) > len(a.Title) {
		a.Title = b.Title
	}
	if a.Price == nil {
		a.Price = b.Price
	}
	if a.Rooms == nil {
		a.Rooms = b.Rooms
	}
	if a.Floor == nil {
		a.Floor = b.Floor
	}
	if a.ImageURL == "" {
		a.ImageURL = b.ImageURL
	}
	if a.AgeText == "" {
		a.AgeText = b.AgeText
	}
	if a.URL == "" {
		a.URL = b.URL
	}
	if b.Promoted {
		a.Promoted = true
		a.PromotionLabel = b.PromotionLabel
	}
	return a
}

func sanitizeFilePart(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func pointerValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
