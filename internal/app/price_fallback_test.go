package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
	"github.com/nicelight/somon-rent-watcher/internal/store"
)

type fallbackDetail struct {
	price  *int
	status int
}

type fallbackHarness struct {
	t           *testing.T
	application *App
	db          *store.DB
	somon       *httptest.Server
	telegram    *httptest.Server

	mu             sync.Mutex
	details        map[int64]fallbackDetail
	detailRequests []int64
	deliveries     []int64
	failDelivery   map[int64]bool
}

func newFallbackHarness(t *testing.T, maxDetails int, details map[int64]fallbackDetail) *fallbackHarness {
	t.Helper()
	h := &fallbackHarness{t: t, details: details, failDelivery: make(map[int64]bool)}
	h.somon = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := fallbackIDFromPath(r.URL.Path)
		h.mu.Lock()
		h.detailRequests = append(h.detailRequests, id)
		detail, ok := h.details[id]
		h.mu.Unlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		if detail.status != 0 {
			w.WriteHeader(detail.status)
			return
		}
		priceMeta := ""
		if detail.price != nil {
			priceMeta = fmt.Sprintf(`<meta property="product:price:amount" content="%d">`, *detail.price)
		}
		fmt.Fprintf(w, `<html><head>%s</head><body><h1>2-комн. квартира, 3 этаж</h1><div>Этаж: 3</div><h2>Описание</h2><p>Долгосрочная аренда</p><div>1 активное объявление</div><div>ID: %d</div></body></html>`, priceMeta, id)
	}))
	h.telegram = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		id := fallbackIDFromCaption(r.Form.Get("text"))
		h.mu.Lock()
		h.deliveries = append(h.deliveries, id)
		fail := h.failDelivery[id]
		h.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if fail {
			fmt.Fprint(w, `{"ok":false,"error_code":500,"description":"test delivery failure"}`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":-200,"type":"group"}}}`)
	}))

	db, err := store.Open(filepath.Join(t.TempDir(), "fallback.db"))
	if err != nil {
		t.Fatal(err)
	}
	h.db = db
	cfg := config.Config{
		TelegramBotToken:     "token",
		TelegramTargetChatID: -200,
		TelegramAPIBase:      h.telegram.URL,
		DBPath:               db.Path(),
		DebugDir:             filepath.Join(t.TempDir(), "debug"),
		CategoryURL:          h.somon.URL + "/category/",
		UserAgent:            "test",
		PollMin:              time.Minute,
		PollMax:              2 * time.Minute,
		RequestDelay:         0,
		HTTPTimeout:          5 * time.Second,
		MinCards:             1,
		MaxDetailsPerPoll:    maxDetails,
		MaxBodyBytes:         1 << 20,
	}
	h.application, err = New(cfg, db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		h.close()
		t.Fatal(err)
	}
	t.Cleanup(h.close)
	return h
}

func (h *fallbackHarness) close() {
	if h.db != nil {
		h.db.Close()
		h.db = nil
	}
	if h.somon != nil {
		h.somon.Close()
		h.somon = nil
	}
	if h.telegram != nil {
		h.telegram.Close()
		h.telegram = nil
	}
}

func (h *fallbackHarness) card(id int64, price *int, position int) model.Card {
	return model.Card{
		ID:       id,
		URL:      h.somon.URL + "/adv/" + strconv.FormatInt(id, 10) + "_x/",
		Title:    "2-комн. квартира, 3 этаж",
		Price:    price,
		Rooms:    fallbackInt(2),
		Floor:    fallbackInt(3),
		Position: position,
	}
}

func (h *fallbackHarness) run(settings filter.Settings, cards []model.Card) (processStats, error) {
	h.t.Helper()
	return h.application.processNewCards(context.Background(), settings, cards, time.Now().UTC())
}

func (h *fallbackHarness) seen(ids ...int64) map[int64]bool {
	h.t.Helper()
	seen, err := h.db.SeenIDs(ids)
	if err != nil {
		h.t.Fatal(err)
	}
	return seen
}

func (h *fallbackHarness) captured() (details, deliveries []int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]int64(nil), h.detailRequests...), append([]int64(nil), h.deliveries...)
}

func TestProcessNewCardsExactSuppressesFallback(t *testing.T) {
	for _, deliveryFails := range []bool{false, true} {
		t.Run(fmt.Sprintf("delivery_fails_%t", deliveryFails), func(t *testing.T) {
			h := newFallbackHarness(t, 10, map[int64]fallbackDetail{
				1: {price: fallbackInt(1000)},
				2: {price: fallbackInt(1200)},
			})
			h.failDelivery[1] = deliveryFails
			settings := filter.DefaultSettings()
			settings.PriceMax = fallbackInt(1000)
			stats, err := h.run(settings, []model.Card{
				h.card(1, fallbackInt(1000), 0),
				h.card(2, fallbackInt(1200), 1),
			})
			if err != nil {
				t.Fatal(err)
			}
			_, deliveries := h.captured()
			if !reflect.DeepEqual(deliveries, []int64{1}) {
				t.Fatalf("deliveries=%v, want exact ID only", deliveries)
			}
			if stats.Sent != 1 && !deliveryFails || stats.Sent != 0 && deliveryFails {
				t.Fatalf("sent=%d deliveryFails=%t", stats.Sent, deliveryFails)
			}
			seen := h.seen(1)
			if seen[1] == deliveryFails {
				t.Fatalf("exact seen=%t deliveryFails=%t", seen[1], deliveryFails)
			}
		})
	}
}

func TestProcessNewCardsSelectsClosestFallbacks(t *testing.T) {
	h := newFallbackHarness(t, 10, map[int64]fallbackDetail{
		11: {price: fallbackInt(1500)},
		12: {price: fallbackInt(1200)},
		13: {price: fallbackInt(1200)},
		14: {price: fallbackInt(1001)},
		15: {price: fallbackInt(1501)},
		16: {price: nil},
		17: {price: fallbackInt(1200)},
	})
	settings := filter.DefaultSettings()
	settings.PriceMax = fallbackInt(1000)
	stats, err := h.run(settings, []model.Card{
		h.card(11, fallbackInt(1500), 0),
		h.card(12, fallbackInt(1200), 1),
		h.card(13, fallbackInt(1200), 2),
		h.card(14, fallbackInt(1001), 3),
		h.card(15, fallbackInt(1501), 4),
		h.card(16, nil, 5),
		h.card(17, nil, 3),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, deliveries := h.captured()
	if !reflect.DeepEqual(deliveries, []int64{14, 12, 13}) {
		t.Fatalf("deliveries=%v, want closest stable order [14 12 13]", deliveries)
	}
	if stats.Sent != 3 {
		t.Fatalf("sent=%d, want 3", stats.Sent)
	}
	seen := h.seen(11, 12, 13, 14, 15, 16, 17)
	for _, id := range []int64{11, 12, 13, 14, 15, 16, 17} {
		if !seen[id] {
			t.Fatalf("final candidate %d remains unseen: %v", id, seen)
		}
	}
}

func TestProcessNewCardsWithoutPriceMaxDoesNotRunFallback(t *testing.T) {
	h := newFallbackHarness(t, 10, map[int64]fallbackDetail{21: {price: fallbackInt(1200)}})
	settings := filter.DefaultSettings()
	stats, err := h.run(settings, []model.Card{h.card(21, fallbackInt(1200), 0)})
	if err != nil {
		t.Fatal(err)
	}
	details, deliveries := h.captured()
	if !reflect.DeepEqual(details, []int64{21}) || !reflect.DeepEqual(deliveries, []int64{21}) || stats.Sent != 1 {
		t.Fatalf("ordinary exact path details=%v deliveries=%v stats=%+v", details, deliveries, stats)
	}
}

func TestProcessNewCardsIncompleteExactSuppressesFallback(t *testing.T) {
	h := newFallbackHarness(t, 1, map[int64]fallbackDetail{
		31: {status: http.StatusInternalServerError},
		32: {price: fallbackInt(1200)},
	})
	settings := filter.DefaultSettings()
	settings.PriceMax = fallbackInt(1000)
	stats, err := h.run(settings, []model.Card{
		h.card(31, nil, 0),
		h.card(32, fallbackInt(1200), 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	details, deliveries := h.captured()
	if !reflect.DeepEqual(details, []int64{31}) || len(deliveries) != 0 || stats.DetailRequests != 1 {
		t.Fatalf("details=%v deliveries=%v stats=%+v", details, deliveries, stats)
	}
	seen := h.seen(31, 32)
	if seen[31] || seen[32] {
		t.Fatalf("incomplete candidates must remain unseen: %v", seen)
	}
}

func TestProcessNewCardsSharedCapSuppressesIncompleteFallback(t *testing.T) {
	h := newFallbackHarness(t, 1, map[int64]fallbackDetail{
		35: {price: fallbackInt(1100)},
		36: {price: fallbackInt(1200)},
	})
	settings := filter.DefaultSettings()
	settings.PriceMax = fallbackInt(1000)
	stats, err := h.run(settings, []model.Card{
		h.card(35, fallbackInt(1100), 0),
		h.card(36, fallbackInt(1200), 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	details, deliveries := h.captured()
	if !reflect.DeepEqual(details, []int64{35}) || len(deliveries) != 0 || stats.DetailRequests != 1 {
		t.Fatalf("details=%v deliveries=%v stats=%+v", details, deliveries, stats)
	}
	seen := h.seen(35, 36)
	if seen[35] || seen[36] {
		t.Fatalf("incomplete fallback candidates must remain unseen: %v", seen)
	}
}

func TestProcessNewCardsPreservesSomonBackoffErrors(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusTooManyRequests} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			h := newFallbackHarness(t, 10, map[int64]fallbackDetail{37: {status: status}})
			settings := filter.DefaultSettings()
			settings.PriceMax = fallbackInt(1000)
			_, err := h.run(settings, []model.Card{h.card(37, nil, 0)})
			gotStatus, _, blocked := somon.IsBlocked(err)
			if !blocked || gotStatus != status {
				t.Fatalf("error=%v status=%d blocked=%t", err, gotStatus, blocked)
			}
			if h.seen(37)[37] {
				t.Fatal("blocked detail must remain unseen")
			}
		})
	}
}

func TestProcessNewCardsFallbackDeliveryFailureDoesNotSubstitute(t *testing.T) {
	h := newFallbackHarness(t, 10, map[int64]fallbackDetail{
		41: {price: fallbackInt(1100)},
		42: {price: fallbackInt(1200)},
		43: {price: fallbackInt(1300)},
		44: {price: fallbackInt(1400)},
	})
	h.failDelivery[41] = true
	settings := filter.DefaultSettings()
	settings.PriceMax = fallbackInt(1000)
	stats, err := h.run(settings, []model.Card{
		h.card(41, fallbackInt(1100), 0),
		h.card(42, fallbackInt(1200), 1),
		h.card(43, fallbackInt(1300), 2),
		h.card(44, fallbackInt(1400), 3),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, deliveries := h.captured()
	if !reflect.DeepEqual(deliveries, []int64{41, 42, 43}) || stats.Sent != 2 {
		t.Fatalf("deliveries=%v stats=%+v", deliveries, stats)
	}
	seen := h.seen(41, 42, 43, 44)
	if seen[41] || !seen[42] || !seen[43] || !seen[44] {
		t.Fatalf("delivery/final-rejection seen state=%v", seen)
	}
}

func TestProcessNewCardsSeenAndFinalRejectionsStaySilent(t *testing.T) {
	h := newFallbackHarness(t, 10, map[int64]fallbackDetail{
		51: {price: fallbackInt(1000)},
		52: {price: nil},
	})
	if err := h.db.MarkSeen([]int64{51}, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	settings := filter.DefaultSettings()
	settings.PriceMax = fallbackInt(1000)
	stats, err := h.run(settings, []model.Card{
		h.card(51, fallbackInt(1000), 0),
		h.card(52, nil, 1),
		h.card(53, fallbackInt(1501), 2),
	})
	if err != nil {
		t.Fatal(err)
	}
	details, deliveries := h.captured()
	if !reflect.DeepEqual(details, []int64{52}) || len(deliveries) != 0 || stats.NewIDs != 2 {
		t.Fatalf("details=%v deliveries=%v stats=%+v", details, deliveries, stats)
	}
	seen := h.seen(51, 52, 53)
	if !seen[51] || !seen[52] || !seen[53] {
		t.Fatalf("final seen state=%v", seen)
	}
}

func fallbackInt(value int) *int { return &value }

func fallbackIDFromPath(path string) int64 {
	path = strings.TrimPrefix(path, "/adv/")
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimSuffix(path, "_x")
	id, _ := strconv.ParseInt(path, 10, 64)
	return id
}

func fallbackIDFromCaption(caption string) int64 {
	const prefix = "ID: <code>"
	start := strings.Index(caption, prefix)
	if start < 0 {
		return 0
	}
	start += len(prefix)
	end := strings.Index(caption[start:], "</code>")
	if end < 0 {
		return 0
	}
	id, _ := strconv.ParseInt(caption[start:start+end], 10, 64)
	return id
}
