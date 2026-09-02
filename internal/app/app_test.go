package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/store"
)

func TestContinuityGapSeparatesRecoveryFromAdminNotification(t *testing.T) {
	now := time.Now()
	gap := detectContinuityGap([]int64{1, 2}, []int64{2, 3}, now.Add(-20*time.Minute), now, 45*time.Minute)
	if gap.detected() || gap.shouldNotifyAdmin() {
		t.Fatalf("unexpected gap: %+v", gap)
	}

	gap = detectContinuityGap([]int64{1}, []int64{2}, now.Add(-20*time.Minute), now, 45*time.Minute)
	if !gap.detected() || gap.shouldNotifyAdmin() || !strings.Contains(gap.reason(), "нет пересечения") {
		t.Fatalf("snapshot-only gap=%+v reason=%q", gap, gap.reason())
	}

	gap = detectContinuityGap([]int64{1}, []int64{2}, now.Add(-60*time.Minute), now, 45*time.Minute)
	if !gap.detected() || !gap.shouldNotifyAdmin() || !strings.Contains(gap.reason(), "последний успешный") {
		t.Fatalf("overdue gap=%+v reason=%q", gap, gap.reason())
	}
}

func TestEffectiveGapAfterAccountsForConfiguredPollRange(t *testing.T) {
	if got := effectiveGapAfter(45*time.Minute, 30); got != 45*time.Minute {
		t.Fatalf("default range gap=%s, want 45m", got)
	}
	if got := effectiveGapAfter(45*time.Minute, 70); got != 75*time.Minute {
		t.Fatalf("extended range gap=%s, want 75m", got)
	}

	now := time.Now()
	gapAfter := effectiveGapAfter(45*time.Minute, 70)
	if gap := detectContinuityGap([]int64{1, 2}, []int64{2, 3}, now.Add(-61*time.Minute), now, gapAfter); gap.detected() {
		t.Fatalf("scheduled 61m delay must not look stale: %+v", gap)
	}
	gap := detectContinuityGap([]int64{1}, []int64{2}, now.Add(-61*time.Minute), now, gapAfter)
	if !gap.detected() || gap.shouldNotifyAdmin() || !strings.Contains(gap.reason(), "нет пересечения") || strings.Contains(gap.reason(), "последний успешный") {
		t.Fatalf("snapshot gap must remain independent from stale-time gap: %+v reason=%q", gap, gap.reason())
	}
}

func TestMergeCardsPreservesPrimaryOrder(t *testing.T) {
	primary := []model.Card{{ID: 2}, {ID: 1}}
	secondary := []model.Card{{ID: 1}, {ID: 3}}
	got := mergeCards(primary, secondary)
	if len(got) != 3 || got[0].ID != 2 || got[1].ID != 1 || got[2].ID != 3 {
		t.Fatalf("merged=%v", got)
	}
}

func TestManualPollRequestIsCoalescedAndRejectedDuringBackoff(t *testing.T) {
	application := &App{
		pollNowCh: make(chan manualPollRequest, 1),
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if err := application.RequestPollNow(-100, 10); err != nil {
		t.Fatal(err)
	}
	if err := application.RequestPollNow(-100, 11); err == nil {
		t.Fatal("expected duplicate request rejection")
	}
	if len(application.pollNowCh) != 1 {
		t.Fatalf("queued requests=%d", len(application.pollNowCh))
	}
	request := <-application.pollNowCh
	if request.chatID != -100 || request.messageID != 10 {
		t.Fatalf("request=%+v", request)
	}
	application.finishManualPollState()
	application.status.BackoffUntil = time.Now().Add(time.Hour)
	if err := application.RequestPollNow(-100, 12); err == nil {
		t.Fatal("expected request rejection during backoff")
	}
}

func TestPollBaselineThenSendsOnlyNewAd(t *testing.T) {
	var mu sync.Mutex
	categoryVersion := 1
	groupMessages := 0
	adminMessages := 0
	recoveryRequests := 0

	somonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/category/":
			mu.Lock()
			version := categoryVersion
			mu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if version == 1 {
				fmt.Fprint(w, testCategoryHTML([]int64{1001}))
			} else {
				fmt.Fprint(w, testCategoryHTML([]int64{1002}))
			}
		case "/adv/1002_x/":
			fmt.Fprint(w, `<html><head><meta property="og:image" content="https://images.example/1002.jpg"><meta property="product:price:amount" content="4500"></head><body><h1>2-комн. квартира, 3 этаж, 60м²</h1><div>Этаж: 3</div><h2>Описание</h2><p>Долгосрочная аренда</p><a href="/author/x">Иван На сайте с Июня 2026 1 активное объявление</a><div>2 минуты назад</div><div>ID: 1002</div></body></html>`)
		default:
			if strings.Contains(r.URL.Path, "/nedvizhimost/") {
				mu.Lock()
				recoveryRequests++
				mu.Unlock()
			}
			http.NotFound(w, r)
		}
	}))
	defer somonServer.Close()

	telegramServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		chatID := r.Form.Get("chat_id")
		mu.Lock()
		if chatID == "-200" {
			groupMessages++
		} else if chatID == "100" {
			adminMessages++
		}
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":1,"type":"private"}}}`)
	}))
	defer telegramServer.Close()

	db, err := store.Open(filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	cfg := config.Config{
		TelegramBotToken:     "token",
		TelegramAdminUserIDs: []int64{100},
		TelegramTargetChatID: -200,
		TelegramAPIBase:      telegramServer.URL,
		DBPath:               db.Path(),
		DebugDir:             filepath.Join(t.TempDir(), "debug"),
		CategoryURL:          somonServer.URL + "/category/",
		UserAgent:            "test",
		PollMin:              time.Minute,
		PollMax:              2 * time.Minute,
		GapAfter:             45 * time.Minute,
		BlockBackoff:         2 * time.Hour,
		RequestDelay:         0,
		HTTPTimeout:          5 * time.Second,
		MinCards:             1,
		MaxDetailsPerPoll:    20,
		MaxBodyBytes:         1 << 20,
	}
	application, err := New(cfg, db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := application.pollOnce(ctx); err != nil {
		t.Fatal(err)
	}
	settings, err := application.LoadSettings()
	if err != nil {
		t.Fatal(err)
	}
	settings.Enabled = true
	if err := application.SaveSettings(settings); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	if groupMessages != 0 || adminMessages != 1 {
		t.Fatalf("after baseline group=%d admin=%d", groupMessages, adminMessages)
	}
	categoryVersion = 2
	mu.Unlock()

	if err := application.pollOnce(ctx); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	if adminMessages != 1 || recoveryRequests == 0 {
		t.Fatalf("snapshot gap admin=%d recovery_requests=%d, want silent recovery", adminMessages, recoveryRequests)
	}
	mu.Unlock()
	if err := application.pollOnce(ctx); err != nil {
		t.Fatal(err)
	}
	application.statusMu.Lock()
	application.status.LastSuccessfulPoll = time.Now().Add(-2 * time.Hour)
	application.statusMu.Unlock()
	if err := application.pollOnce(ctx); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	if groupMessages != 1 {
		t.Fatalf("new ad messages=%d, want 1", groupMessages)
	}
	if adminMessages != 2 {
		t.Fatalf("admin messages=%d, want baseline plus overdue warning", adminMessages)
	}
	if count, err := db.CountSeen(); err != nil || count != 2 {
		t.Fatalf("seen count=%d err=%v", count, err)
	}
}

func testCategoryHTML(ids []int64) string {
	var b strings.Builder
	b.WriteString("<html><body><main>")
	for i, id := range ids {
		fmt.Fprintf(&b, `<article class="item"><span>%d c.</span><a href="/adv/%d_x/">2-комн. квартира, 3 этаж, 60м²</a><span>%d минуты назад Душанбе</span></article>`, 4500+i, id, i+1)
	}
	b.WriteString("</main></body></html>")
	return b.String()
}

func TestPausedMonitoringConsumesFreshIDsWithoutSending(t *testing.T) {
	var mu sync.Mutex
	categoryVersion := 1
	groupMessages := 0
	detailRequests := 0

	somonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		version := categoryVersion
		mu.Unlock()
		switch r.URL.Path {
		case "/category/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			switch version {
			case 1:
				fmt.Fprint(w, testCategoryHTML([]int64{2001}))
			case 2:
				fmt.Fprint(w, testCategoryHTML([]int64{2002, 2001}))
			default:
				fmt.Fprint(w, testCategoryHTML([]int64{2003, 2002, 2001}))
			}
		case "/adv/2003_x/":
			mu.Lock()
			detailRequests++
			mu.Unlock()
			fmt.Fprint(w, `<html><head><meta property="product:price:amount" content="4500"></head><body><h1>2-комн. квартира, 3 этаж, 60м²</h1><div>Этаж: 3</div><h2>Описание</h2><p>Долгосрочная аренда</p><div>1 активное объявление</div><div>ID: 2003</div></body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer somonServer.Close()

	telegramServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("chat_id") == "-200" {
			mu.Lock()
			groupMessages++
			mu.Unlock()
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":1,"type":"private"}}}`)
	}))
	defer telegramServer.Close()

	db, err := store.Open(filepath.Join(t.TempDir(), "paused.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	cfg := config.Config{
		TelegramBotToken:     "token",
		TelegramAdminUserIDs: []int64{100},
		TelegramTargetChatID: -200,
		TelegramAPIBase:      telegramServer.URL,
		DBPath:               db.Path(),
		DebugDir:             filepath.Join(t.TempDir(), "debug"),
		CategoryURL:          somonServer.URL + "/category/",
		UserAgent:            "test",
		PollMin:              time.Minute,
		PollMax:              2 * time.Minute,
		GapAfter:             45 * time.Minute,
		BlockBackoff:         2 * time.Hour,
		HTTPTimeout:          5 * time.Second,
		MinCards:             1,
		MaxDetailsPerPoll:    20,
		MaxBodyBytes:         1 << 20,
	}
	application, err := New(cfg, db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := application.pollOnce(ctx); err != nil { // baseline, paused by default
		t.Fatal(err)
	}
	mu.Lock()
	categoryVersion = 2
	mu.Unlock()
	if err := application.pollOnce(ctx); err != nil { // 2002 is consumed while paused
		t.Fatal(err)
	}
	settings, err := application.LoadSettings()
	if err != nil {
		t.Fatal(err)
	}
	settings.Enabled = true
	if err := application.SaveSettings(settings); err != nil {
		t.Fatal(err)
	}
	if err := application.pollOnce(ctx); err != nil { // 2002 must not backfill
		t.Fatal(err)
	}
	mu.Lock()
	categoryVersion = 3
	mu.Unlock()
	if err := application.pollOnce(ctx); err != nil { // only 2003 is delivered
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	if groupMessages != 1 {
		t.Fatalf("group messages=%d, want 1", groupMessages)
	}
	if detailRequests != 1 {
		t.Fatalf("detail requests=%d, want 1", detailRequests)
	}
	if count, err := db.CountSeen(); err != nil || count != 3 {
		t.Fatalf("seen count=%d err=%v", count, err)
	}
}
