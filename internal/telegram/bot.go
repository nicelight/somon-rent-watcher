package telegram

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
)

// Backend is intentionally narrow: Telegram only needs settings, runtime status,
// and a persisted update offset. The application implements this interface.
type Backend interface {
	LoadSettings() (filter.Settings, error)
	SaveSettings(filter.Settings) error
	RuntimeStatus() model.RuntimeStatus
	TelegramOffset() (int64, error)
	SetTelegramOffset(int64) error
}

type Bot struct {
	client       *Client
	backend      Backend
	adminUserIDs []int64
	adminUserSet map[int64]struct{}
	targetChatID int64
	logger       *slog.Logger

	pendingMu sync.Mutex
	pending   map[pendingKey]string
}

type pendingKey struct {
	userID int64
	chatID int64
}

func NewBot(client *Client, backend Backend, adminUserIDs []int64, targetChatID int64, logger *slog.Logger) *Bot {
	if logger == nil {
		logger = slog.Default()
	}
	adminSet := make(map[int64]struct{}, len(adminUserIDs))
	cleanAdminIDs := make([]int64, 0, len(adminUserIDs))
	for _, id := range adminUserIDs {
		if id <= 0 {
			continue
		}
		if _, exists := adminSet[id]; exists {
			continue
		}
		adminSet[id] = struct{}{}
		cleanAdminIDs = append(cleanAdminIDs, id)
	}
	return &Bot{
		client:       client,
		backend:      backend,
		adminUserIDs: cleanAdminIDs,
		adminUserSet: adminSet,
		targetChatID: targetChatID,
		logger:       logger,
		pending:      make(map[pendingKey]string),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	// getUpdates and webhooks are mutually exclusive. Removing a stale webhook is
	// safe for this bot and does not drop queued updates.
	if err := b.client.DeleteWebhook(ctx, false); err != nil {
		return fmt.Errorf("delete Telegram webhook: %w", err)
	}

	offset, err := b.backend.TelegramOffset()
	if err != nil {
		return fmt.Errorf("load Telegram offset: %w", err)
	}
	b.logger.Info("Telegram long polling started", "offset", offset)

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		updates, err := b.client.GetUpdates(ctx, offset, 50)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			delay := telegramRetryDelay(err)
			b.logger.Error("Telegram getUpdates failed", "error", err, "retry_in", delay)
			if !sleepContext(ctx, delay) {
				return ctx.Err()
			}
			continue
		}
		for _, update := range updates {
			if err := b.processUpdate(ctx, update); err != nil {
				b.logger.Error("Telegram update processing failed", "update_id", update.UpdateID, "error", err)
			}
			next := update.UpdateID + 1
			if err := b.backend.SetTelegramOffset(next); err != nil {
				return fmt.Errorf("persist Telegram offset %d: %w", next, err)
			}
			offset = next
		}
	}
}

func (b *Bot) processUpdate(ctx context.Context, update Update) error {
	if update.CallbackQuery != nil {
		return b.processCallback(ctx, update.CallbackQuery)
	}
	if update.Message != nil {
		return b.processMessage(ctx, update.Message)
	}
	return nil
}

func (b *Bot) processMessage(ctx context.Context, msg *Message) error {
	if msg == nil || msg.From == nil || !b.canManage(msg.From.ID, msg.Chat) {
		return nil
	}
	text := strings.TrimSpace(msg.Text)
	if text == "" {
		return nil
	}

	if strings.HasPrefix(text, "/") {
		b.setPending(msg.From.ID, msg.Chat.ID, "")
		command := strings.Fields(text)[0]
		if i := strings.IndexByte(command, '@'); i >= 0 {
			command = command[:i]
		}
		switch strings.ToLower(command) {
		case "/start", "/filter":
			return b.sendMainMenu(ctx, msg.Chat.ID)
		case "/status":
			return b.sendStatus(ctx, msg.Chat.ID)
		case "/cancel":
			_, err := b.client.SendMessage(ctx, msg.Chat.ID, "Ввод отменён.", nil)
			return err
		default:
			_, err := b.client.SendMessage(ctx, msg.Chat.ID, "Команды: /filter, /status, /cancel", nil)
			return err
		}
	}

	switch b.getPending(msg.From.ID, msg.Chat.ID) {
	case "price":
		return b.applyPriceInput(ctx, msg.From.ID, msg.Chat.ID, text)
	case "negative":
		return b.applyNegativeInput(ctx, msg.From.ID, msg.Chat.ID, text)
	default:
		if msg.Chat.Type != "private" {
			return nil
		}
		return b.sendMainMenu(ctx, msg.Chat.ID)
	}
}

func (b *Bot) processCallback(ctx context.Context, query *CallbackQuery) error {
	if query == nil || query.Message == nil || !b.canManage(query.From.ID, query.Message.Chat) {
		if query != nil {
			_ = b.client.AnswerCallbackQuery(ctx, query.ID, "Недоступно", false)
		}
		return nil
	}

	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	data := query.Data
	var actionErr error

	switch {
	case data == "m:main":
		actionErr = b.editMainMenu(ctx, chatID, messageID)
	case data == "e:toggle":
		actionErr = b.toggleEnabled(ctx, chatID, messageID)
	case data == "m:rooms":
		actionErr = b.editRoomsMenu(ctx, chatID, messageID)
	case data == "m:floors":
		actionErr = b.editFloorsMenu(ctx, chatID, messageID)
	case data == "m:author":
		actionErr = b.editAuthorMenu(ctx, chatID, messageID)
	case data == "m:type":
		actionErr = b.editTypeMenu(ctx, chatID, messageID)
	case data == "i:price":
		b.setPending(query.From.ID, chatID, "price")
		_, actionErr = b.client.SendMessage(ctx, chatID,
			"Введите диапазон цены одной строкой:\n<code>3500-6000</code>, <code>-6000</code>, <code>3500-</code> или <code>-</code> для снятия ограничения.\n\n/cancel — отмена.", nil)
	case data == "i:negative":
		b.setPending(query.From.ID, chatID, "negative")
		_, actionErr = b.client.SendMessage(ctx, chatID,
			"Введите негативные слова/фразы через запятую.\nНапример: <code>посуточно, подселение, без детей</code>\n<code>-</code> — очистить список.\n\n/cancel — отмена.", nil)
	case strings.HasPrefix(data, "r:"):
		actionErr = b.toggleRoom(ctx, chatID, messageID, strings.TrimPrefix(data, "r:"))
	case strings.HasPrefix(data, "f:"):
		actionErr = b.toggleFloor(ctx, chatID, messageID, strings.TrimPrefix(data, "f:"))
	case strings.HasPrefix(data, "a:"):
		actionErr = b.setAuthorLimit(ctx, chatID, messageID, strings.TrimPrefix(data, "a:"))
	case data == "t:o":
		actionErr = b.toggleType(ctx, chatID, messageID, true)
	case data == "t:p":
		actionErr = b.toggleType(ctx, chatID, messageID, false)
	default:
		actionErr = errors.New("неизвестная кнопка")
	}

	if actionErr != nil {
		_ = b.client.AnswerCallbackQuery(ctx, query.ID, actionErr.Error(), true)
		return actionErr
	}
	return b.client.AnswerCallbackQuery(ctx, query.ID, "", false)
}

func (b *Bot) applyPriceInput(ctx context.Context, userID, chatID int64, text string) error {
	min, max, err := filter.ParsePriceInput(text)
	if err != nil {
		_, sendErr := b.client.SendMessage(ctx, chatID, "Ошибка: "+html.EscapeString(err.Error()), nil)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	settings.PriceMin, settings.PriceMax = min, max
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	b.setPending(userID, chatID, "")
	return b.sendMainMenu(ctx, chatID)
}

func (b *Bot) applyNegativeInput(ctx context.Context, userID, chatID int64, text string) error {
	words, err := filter.ParseNegativeInput(text)
	if err != nil {
		_, sendErr := b.client.SendMessage(ctx, chatID, "Ошибка: "+html.EscapeString(err.Error()), nil)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	settings.NegativeWords = words
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	b.setPending(userID, chatID, "")
	return b.sendMainMenu(ctx, chatID)
}

func (b *Bot) toggleEnabled(ctx context.Context, chatID, messageID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	settings.Enabled = !settings.Enabled
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	return b.editMainMenu(ctx, chatID, messageID)
}

func (b *Bot) toggleRoom(ctx context.Context, chatID, messageID int64, choice string) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	if err := settings.ToggleRoom(choice); err != nil {
		return err
	}
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	return b.editRoomsMenu(ctx, chatID, messageID)
}

func (b *Bot) toggleFloor(ctx context.Context, chatID, messageID int64, choice string) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	if err := settings.ToggleFloor(choice); err != nil {
		return err
	}
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	return b.editFloorsMenu(ctx, chatID, messageID)
}

func (b *Bot) setAuthorLimit(ctx context.Context, chatID, messageID int64, raw string) error {
	limit, err := strconv.Atoi(raw)
	if err != nil {
		return errors.New("неверный лимит")
	}
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	settings.SellerAdsLimit = limit
	if err := settings.NormalizeAndValidate(); err != nil {
		return err
	}
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	return b.editAuthorMenu(ctx, chatID, messageID)
}

func (b *Bot) toggleType(ctx context.Context, chatID, messageID int64, ordinary bool) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	if ordinary {
		err = settings.ToggleOrdinary()
	} else {
		err = settings.TogglePromoted()
	}
	if err != nil {
		return err
	}
	if err := b.backend.SaveSettings(settings); err != nil {
		return err
	}
	return b.editTypeMenu(ctx, chatID, messageID)
}

func (b *Bot) sendMainMenu(ctx context.Context, chatID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	_, err = b.client.SendMessage(ctx, chatID, SettingsText(settings), mainKeyboard(settings))
	return err
}

func (b *Bot) editMainMenu(ctx context.Context, chatID, messageID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	_, err = b.client.EditMessageText(ctx, chatID, messageID, SettingsText(settings), mainKeyboard(settings))
	return ignoreNotModified(err)
}

func (b *Bot) editRoomsMenu(ctx context.Context, chatID, messageID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	_, err = b.client.EditMessageText(ctx, chatID, messageID,
		"<b>Количество комнат</b>\nМожно выбрать несколько вариантов.", roomsKeyboard(settings))
	return ignoreNotModified(err)
}

func (b *Bot) editFloorsMenu(ctx context.Context, chatID, messageID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	_, err = b.client.EditMessageText(ctx, chatID, messageID,
		"<b>Этаж</b>\nМожно выбрать несколько вариантов. 14+ включает 14-й этаж и выше.", floorsKeyboard(settings))
	return ignoreNotModified(err)
}

func (b *Bot) editAuthorMenu(ctx context.Context, chatID, messageID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	_, err = b.client.EditMessageText(ctx, chatID, messageID,
		"<b>Активных объявлений автора</b>\nСчётчик Somon включает текущее объявление.", authorKeyboard(settings))
	return ignoreNotModified(err)
}

func (b *Bot) editTypeMenu(ctx context.Context, chatID, messageID int64) error {
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	_, err = b.client.EditMessageText(ctx, chatID, messageID,
		"<b>Тип объявления</b>\nVIP включает все продвигаемые карточки Somon: VIP и ТОП.", typeKeyboard(settings))
	return ignoreNotModified(err)
}

func (b *Bot) sendStatus(ctx context.Context, chatID int64) error {
	status := b.backend.RuntimeStatus()
	settings, err := b.backend.LoadSettings()
	if err != nil {
		return err
	}
	text := StatusText(status) + "\n\n" + SettingsText(settings)
	_, err = b.client.SendMessage(ctx, chatID, text, mainKeyboard(settings))
	return err
}

func (b *Bot) SendAdmin(ctx context.Context, text string) error {
	return b.sendToAdmins(ctx, html.EscapeString(text))
}

func (b *Bot) SendAdminHTML(ctx context.Context, text string) error {
	return b.sendToAdmins(ctx, text)
}

func (b *Bot) sendToAdmins(ctx context.Context, text string) error {
	var errs []error
	for _, adminUserID := range b.adminUserIDs {
		if _, err := b.client.SendMessage(ctx, adminUserID, text, nil); err != nil {
			errs = append(errs, fmt.Errorf("admin %d: %w", adminUserID, err))
		}
	}
	return errors.Join(errs...)
}

func (b *Bot) SendAd(ctx context.Context, ad model.Ad) error {
	caption := AdCaption(ad)
	keyboard := &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{
		{Text: "Открыть на Somon", URL: ad.URL},
	}}}
	if ad.ImageURL != "" {
		if _, err := b.client.SendPhoto(ctx, b.targetChatID, ad.ImageURL, caption, keyboard); err == nil {
			return nil
		} else {
			// Retry as text only for a deterministic Telegram validation error. On a
			// network/timeout error the photo may already have been accepted; an
			// immediate text fallback could create a duplicate notification.
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.Code != 400 {
				return err
			}
			b.logger.Warn("Telegram rejected remote photo, falling back to text", "ad_id", ad.ID, "error", err)
		}
	}
	_, err := b.client.SendMessage(ctx, b.targetChatID, caption, keyboard)
	return err
}

func (b *Bot) canManage(userID int64, chat Chat) bool {
	if _, allowed := b.adminUserSet[userID]; !allowed {
		return false
	}
	return chat.Type == "private" || chat.ID == b.targetChatID
}

func (b *Bot) setPending(userID, chatID int64, value string) {
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()
	key := pendingKey{userID: userID, chatID: chatID}
	if value == "" {
		delete(b.pending, key)
		return
	}
	b.pending[key] = value
}

func (b *Bot) getPending(userID, chatID int64) string {
	b.pendingMu.Lock()
	defer b.pendingMu.Unlock()
	return b.pending[pendingKey{userID: userID, chatID: chatID}]
}

func ignoreNotModified(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) && strings.Contains(strings.ToLower(apiErr.Description), "message is not modified") {
		return nil
	}
	return err
}

func telegramRetryDelay(err error) time.Duration {
	delay := 5 * time.Second
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.RetryAfter > 0 {
		retry := time.Duration(apiErr.RetryAfter) * time.Second
		if retry > delay {
			delay = retry
		}
	}
	return delay
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
