package telegram

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
)

type botTestBackend struct {
	mu       sync.Mutex
	settings filter.Settings
	offset   int64
}

func (b *botTestBackend) LoadSettings() (filter.Settings, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.settings, nil
}

func (b *botTestBackend) SaveSettings(settings filter.Settings) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.settings = settings
	return nil
}

func (b *botTestBackend) RuntimeStatus() model.RuntimeStatus { return model.RuntimeStatus{} }

func (b *botTestBackend) TelegramOffset() (int64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.offset, nil
}

func (b *botTestBackend) SetTelegramOffset(offset int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.offset = offset
	return nil
}

func TestSendAdFallsBackToTextOnTelegramPhotoValidationError(t *testing.T) {
	var photoCalls, messageCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botTOKEN/sendPhoto":
			photoCalls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok":false,"error_code":400,"description":"Bad Request: failed to get HTTP URL content"}`)
		case "/botTOKEN/sendMessage":
			messageCalls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":-100,"type":"group"}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	bot := NewBot(NewClient(server.URL, "TOKEN"), nil, []int64{1}, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
	err := bot.SendAd(context.Background(), model.Ad{Card: model.Card{
		ID: 1, URL: "https://somon.tj/adv/1_x/", Title: "1-комн. квартира", ImageURL: "https://example.invalid/image.jpg",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if photoCalls.Load() != 1 || messageCalls.Load() != 1 {
		t.Fatalf("photo=%d message=%d", photoCalls.Load(), messageCalls.Load())
	}
}

func TestSendAdDoesNotImmediatelyDuplicateAfterAmbiguousPhotoFailure(t *testing.T) {
	var photoCalls, messageCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/botTOKEN/sendPhoto":
			photoCalls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok":false,"error_code":500,"description":"Internal Server Error"}`)
		case "/botTOKEN/sendMessage":
			messageCalls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":-100,"type":"group"}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	bot := NewBot(NewClient(server.URL, "TOKEN"), nil, []int64{1}, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
	err := bot.SendAd(context.Background(), model.Ad{Card: model.Card{
		ID: 2, URL: "https://somon.tj/adv/2_x/", Title: "2-комн. квартира", ImageURL: "https://example.invalid/image.jpg",
	}})
	if err == nil {
		t.Fatal("expected Telegram error")
	}
	if photoCalls.Load() != 1 || messageCalls.Load() != 0 {
		t.Fatalf("photo=%d message=%d", photoCalls.Load(), messageCalls.Load())
	}
}

func TestAdminsCanManageOnlyPrivateOrTargetGroup(t *testing.T) {
	var mu sync.Mutex
	var sentChats []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/botTOKEN/sendMessage":
			mu.Lock()
			sentChats = append(sentChats, r.Form.Get("chat_id"))
			mu.Unlock()
			fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":1,"type":"private"}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	backend := &botTestBackend{settings: filter.DefaultSettings()}
	bot := NewBot(NewClient(server.URL, "TOKEN"), backend, []int64{1, 2}, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
	ctx := context.Background()

	if err := bot.processMessage(ctx, &Message{From: &User{ID: 2}, Chat: Chat{ID: -100, Type: "supergroup"}, Text: "/status@watcher_bot"}); err != nil {
		t.Fatal(err)
	}
	if err := bot.processMessage(ctx, &Message{From: &User{ID: 3}, Chat: Chat{ID: -100, Type: "supergroup"}, Text: "/status"}); err != nil {
		t.Fatal(err)
	}
	if err := bot.processMessage(ctx, &Message{From: &User{ID: 1}, Chat: Chat{ID: -200, Type: "supergroup"}, Text: "/status"}); err != nil {
		t.Fatal(err)
	}
	if err := bot.processMessage(ctx, &Message{From: &User{ID: 1}, Chat: Chat{ID: 1, Type: "private"}, Text: "/status"}); err != nil {
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(sentChats) != 2 || sentChats[0] != "-100" || sentChats[1] != "1" {
		t.Fatalf("sent chats=%v", sentChats)
	}
}

func TestGroupTextInputIsIsolatedByAdminAndChat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/botTOKEN/sendMessage":
			fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":-100,"type":"supergroup"}}}`)
		case "/botTOKEN/answerCallbackQuery":
			fmt.Fprint(w, `{"ok":true,"result":true}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	backend := &botTestBackend{settings: filter.DefaultSettings()}
	bot := NewBot(NewClient(server.URL, "TOKEN"), backend, []int64{1, 2}, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
	ctx := context.Background()
	group := Chat{ID: -100, Type: "supergroup"}

	if err := bot.processCallback(ctx, &CallbackQuery{
		ID: "callback-1", From: User{ID: 1}, Message: &Message{MessageID: 10, Chat: group}, Data: "i:price",
	}); err != nil {
		t.Fatal(err)
	}
	if got := bot.getPending(1, -100); got != "price" {
		t.Fatalf("admin 1 pending=%q", got)
	}
	if got := bot.getPending(2, -100); got != "" {
		t.Fatalf("admin 2 pending=%q", got)
	}

	if err := bot.processMessage(ctx, &Message{From: &User{ID: 2}, Chat: group, Text: "3500-6000"}); err != nil {
		t.Fatal(err)
	}
	settings, err := backend.LoadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if settings.PriceMin != nil || settings.PriceMax != nil {
		t.Fatalf("another admin changed pending input: %+v", settings)
	}

	if err := bot.processMessage(ctx, &Message{From: &User{ID: 1}, Chat: group, Text: "3500-6000"}); err != nil {
		t.Fatal(err)
	}
	settings, err = backend.LoadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if settings.PriceMin == nil || *settings.PriceMin != 3500 || settings.PriceMax == nil || *settings.PriceMax != 6000 {
		t.Fatalf("price settings=%+v", settings)
	}
	if got := bot.getPending(1, -100); got != "" {
		t.Fatalf("pending after save=%q", got)
	}
}

func TestAdminNotificationsAreSentToEveryAdmin(t *testing.T) {
	var mu sync.Mutex
	var sentChats []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		mu.Lock()
		sentChats = append(sentChats, r.Form.Get("chat_id"))
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":1,"type":"private"}}}`)
	}))
	defer server.Close()

	bot := NewBot(NewClient(server.URL, "TOKEN"), nil, []int64{1, 2, 2, 0}, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err := bot.SendAdmin(context.Background(), "status"); err != nil {
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(sentChats) != 2 || sentChats[0] != "1" || sentChats[1] != "2" {
		t.Fatalf("sent chats=%v", sentChats)
	}
}
