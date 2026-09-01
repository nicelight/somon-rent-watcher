package telegram

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

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

	bot := NewBot(NewClient(server.URL, "TOKEN"), nil, 1, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
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

	bot := NewBot(NewClient(server.URL, "TOKEN"), nil, 1, -100, slog.New(slog.NewTextHandler(io.Discard, nil)))
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
