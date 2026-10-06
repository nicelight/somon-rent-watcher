//go:build cgo

package telegram_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/app"
	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/store"
	"github.com/nicelight/somon-rent-watcher/internal/telegram"
)

type keywordAPIRequest struct {
	method string
	form   url.Values
}
type keywordHarness struct {
	t        *testing.T
	db       *store.DB
	a        *app.App
	bot      *telegram.Bot
	path     string
	requests []keywordAPIRequest
	cfg      config.Config
}

func newKeywordHarness(t *testing.T) *keywordHarness {
	t.Helper()
	h := &keywordHarness{t: t, path: filepath.Join(t.TempDir(), "search.db")}
	var err error
	h.db, err = store.Open(h.path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { h.db.Close() })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		method := strings.TrimPrefix(r.URL.Path, "/botTOKEN/")
		h.requests = append(h.requests, keywordAPIRequest{method, r.Form})
		w.Header().Set("Content-Type", "application/json")
		if method == "answerCallbackQuery" {
			fmt.Fprint(w, `{"ok":true,"result":true}`)
		} else {
			fmt.Fprint(w, `{"ok":true,"result":{"message_id":1,"chat":{"id":1,"type":"private"}}}`)
		}
	}))
	t.Cleanup(server.Close)
	settings := filter.DefaultSettings()
	settings.Enabled = true
	max := 6000
	settings.PriceMax = &max
	raw, err := settings.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if err = h.db.SaveSettingsJSON(raw); err != nil {
		t.Fatal(err)
	}
	if err = h.db.MarkSeen([]int64{101, 202}, time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if err = h.db.SetStates(map[string]string{"initialized": "1", "previous_ordinary_ids": "[101,202]", "last_successful_poll_at": "2026-10-06T12:00:00Z", "telegram_offset": "88"}); err != nil {
		t.Fatal(err)
	}
	h.cfg = config.Config{PollMin: time.Minute, PollMax: 2 * time.Minute, TelegramAPIBase: server.URL, TelegramBotToken: "TOKEN", TelegramAdminUserIDs: []int64{1, 2}, TelegramTargetChatID: -100}
	h.a, err = app.New(h.cfg, h.db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	h.bot = telegram.NewBot(telegram.NewClient(server.URL, "TOKEN"), h.a, []int64{1, 2}, -100, nil)
	return h
}
func (h *keywordHarness) callback(user int64, chat telegram.Chat, data string) error {
	h.t.Helper()
	start := len(h.requests)
	err := telegram.ProcessKeywordUpdateForTest(h.bot, context.Background(), telegram.Update{CallbackQuery: &telegram.CallbackQuery{ID: fmt.Sprintf("callback-%d", start), From: telegram.User{ID: user}, Message: &telegram.Message{MessageID: 10, Chat: chat}, Data: data}})
	trace := h.requests[start:]
	if len(trace) == 0 || trace[0].method != "answerCallbackQuery" {
		h.t.Fatalf("ack must precede output %s: %+v", data, trace)
	}
	for _, r := range trace {
		if r.method == "editMessageText" {
			h.t.Fatalf("keyword route edited old message: %s", data)
		}
	}
	return err
}
func (h *keywordHarness) text(user int64, chat telegram.Chat, text string) {
	h.t.Helper()
	if err := telegram.ProcessKeywordUpdateForTest(h.bot, context.Background(), telegram.Update{Message: &telegram.Message{From: &telegram.User{ID: user}, Chat: chat, Text: text}}); err != nil {
		h.t.Fatal(err)
	}
}
func (h *keywordHarness) search(id int64) model.KeywordSearch {
	h.t.Helper()
	s, found, err := h.a.KeywordSearch(id)
	if err != nil || !found {
		h.t.Fatalf("search %d found=%v err=%v", id, found, err)
	}
	return s
}
func (h *keywordHarness) list() []model.KeywordSearch {
	h.t.Helper()
	s, err := h.a.ListKeywordSearches()
	if err != nil {
		h.t.Fatal(err)
	}
	return s
}
func (h *keywordHarness) click(user int64, chat telegram.Chat, data string) {
	h.t.Helper()
	if err := h.callback(user, chat, data); err != nil {
		h.t.Fatal(data, err)
	}
}
func (h *keywordHarness) snapshotRental() map[string]any {
	h.t.Helper()
	raw, ok, err := h.db.LoadSettingsJSON()
	if err != nil || !ok {
		h.t.Fatal(err)
	}
	seen, err := h.db.SeenIDs([]int64{101, 202, 303})
	if err != nil {
		h.t.Fatal(err)
	}
	count, err := h.db.CountSeen()
	if err != nil {
		h.t.Fatal(err)
	}
	out := map[string]any{"settings": raw, "seen": seen, "seen_count": count}
	for _, k := range []string{"initialized", "previous_ordinary_ids", "last_successful_poll_at", "telegram_offset"} {
		v, ok, err := h.db.GetState(k)
		if err != nil || !ok {
			h.t.Fatal(k, err)
		}
		out[k] = v
	}
	return out
}
func (h *keywordHarness) deliveryTrace() []keywordAPIRequest {
	h.t.Helper()
	start := len(h.requests)
	if err := h.bot.SendAd(context.Background(), model.Ad{Card: model.Card{ID: 303, Title: "Квартира", URL: "https://somon.tj/adv/303_x/"}}); err != nil {
		h.t.Fatal(err)
	}
	return append([]keywordAPIRequest(nil), h.requests[start:]...)
}

func TestKeywordSearchManagementThroughRealApp(t *testing.T) {
	h := newKeywordHarness(t)
	private := telegram.Chat{ID: 1, Type: "private"}
	group := telegram.Chat{ID: -100, Type: "supergroup"}
	protected := h.snapshotRental()
	rentalTrace := h.deliveryTrace()
	h.text(1, private, "/filter")
	if !strings.Contains(h.requests[len(h.requests)-1].form.Get("reply_markup"), "ks:list") {
		t.Fatal("existing menu must reach searches")
	}
	h.click(1, private, "ks:list")
	h.click(1, private, "ks:new")
	h.text(1, private, "  стол & стулья  ")
	s1 := h.list()[0]
	if s1.ID <= 0 || s1.Enabled || s1.Revision != 1 || s1.Phrase != "стол & стулья" || s1.CategoryKey != "all" || s1.CityKey != "country" {
		t.Fatalf("invalid initial search %+v", s1)
	}
	summary := h.requests[len(h.requests)-1].form.Get("text")
	if !strings.Contains(summary, "стол &amp; стулья") || !strings.Contains(summary, "Все категории") || !strings.Contains(summary, "Вся страна") || !strings.Contains(summary, "на паузе") {
		t.Fatal(summary)
	}
	id := s1.ID
	h.click(1, private, fmt.Sprintf("ks:city:%d", id))
	buttons := h.requests[len(h.requests)-1].form.Get("reply_markup")
	for _, label := range []string{"Вся страна", "Душанбе", "Восе", "Дангара"} {
		if !strings.Contains(buttons, label) {
			t.Fatal(buttons)
		}
	}
	h.click(1, private, fmt.Sprintf("ks:setcity:%d:vose", id))
	h.click(1, private, fmt.Sprintf("ks:category:%d", id))
	buttons = h.requests[len(h.requests)-1].form.Get("reply_markup")
	for _, label := range []string{"Все категории", "Мебель", "Столы и стулья", "Услуги"} {
		if !strings.Contains(buttons, label) {
			t.Fatal(buttons)
		}
	}
	h.click(1, private, fmt.Sprintf("ks:setcategory:%d:tables_chairs", id))
	if s := h.search(id); s.CityKey != "vose" || s.CategoryKey != "tables_chairs" {
		t.Fatal(s)
	}
	h.click(1, private, fmt.Sprintf("ks:price:%d", id))
	h.text(1, private, "100-200")
	h.click(1, private, fmt.Sprintf("ks:enable:%d", id))
	s1 = h.search(id)
	if !s1.Enabled || *s1.PriceMin != 100 || *s1.PriceMax != 200 || s1.Revision != 5 {
		t.Fatalf("summary/enable settings %+v", s1)
	}
	h.click(2, group, "ks:new")
	h.text(2, group, "ремонт")
	s2 := h.list()[1]
	if s2.ID <= id || s2.Enabled {
		t.Fatal(s2)
	}
	h.click(2, group, fmt.Sprintf("ks:setcategory:%d:services", s2.ID))
	h.click(2, group, fmt.Sprintf("ks:setcity:%d:dangara", s2.ID))
	h.click(2, group, fmt.Sprintf("ks:price:%d", s2.ID))
	h.text(2, group, "0-")
	h.click(2, group, fmt.Sprintf("ks:enable:%d", s2.ID))
	h.click(1, private, fmt.Sprintf("ks:phrase:%d", id))
	h.text(1, private, "обеденный стол")
	s1 = h.search(id)
	s2 = h.search(s2.ID)
	if s1.Phrase != "обеденный стол" || s1.Revision != 6 || s2.Phrase != "ремонт" || s2.CityKey != "dangara" || s2.PriceMin == nil || *s2.PriceMin != 0 || s2.PriceMax != nil {
		t.Fatal(s1, s2)
	}
	h.click(1, private, fmt.Sprintf("ks:disable:%d", id))
	h.click(1, private, fmt.Sprintf("ks:enable:%d", id))
	expected := h.list()
	if !reflect.DeepEqual(protected, h.snapshotRental()) || !reflect.DeepEqual(rentalTrace, h.deliveryTrace()) {
		t.Fatal("search management altered rental rows/delivery trace")
	}
	if err := h.db.Close(); err != nil {
		t.Fatal(err)
	}
	var err error
	h.db, err = store.Open(h.path)
	if err != nil {
		t.Fatal(err)
	}
	h.a, err = app.New(h.cfg, h.db, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, h.list()) || !reflect.DeepEqual(protected, h.snapshotRental()) {
		t.Fatal("reopen changed searches or rental")
	}
	h.bot = telegram.NewBot(telegram.NewClient(h.cfg.TelegramAPIBase, "TOKEN"), h.a, []int64{1, 2}, -100, nil)
	if !reflect.DeepEqual(rentalTrace, h.deliveryTrace()) {
		t.Fatal("reopen changed rental delivery payload")
	}
	for _, r := range h.requests {
		if r.method == "editMessageText" {
			t.Fatal("new routes must be append-only")
		}
	}
	data, _ := json.Marshal(expected)
	t.Logf("GREEN AC001: two persisted searches %s; ack before every fresh callback output; zero edits. GREEN AC007: exact public rental snapshot + delivery payload retained through writes/reopen", data)
}

func TestKeywordSearchInvalidAndCrossContextInput(t *testing.T) {
	h := newKeywordHarness(t)
	p1 := telegram.Chat{ID: 1, Type: "private"}
	p2 := telegram.Chat{ID: 2, Type: "private"}
	g := telegram.Chat{ID: -100, Type: "supergroup"}
	bad := telegram.Chat{ID: -200, Type: "supergroup"}
	before := h.snapshotRental()
	for _, action := range []struct {
		user int64
		chat telegram.Chat
	}{{3, g}, {1, bad}} {
		start := len(h.requests)
		h.click(action.user, action.chat, "ks:new")
		h.text(action.user, action.chat, "unauthorized")
		if len(h.requests) != start+1 || len(h.list()) != 0 {
			t.Fatal("unauthorized input mutated or sent messages")
		}
	}
	h.click(1, p1, "ks:new")
	h.text(1, p1, "   ")
	if len(h.list()) != 0 {
		t.Fatal("blank phrase saved")
	}
	h.text(2, g, "other admin")
	h.text(1, g, "other chat")
	if len(h.list()) != 0 {
		t.Fatal("pending create leaked")
	}
	h.text(1, p1, "стол")
	id1 := h.list()[0].ID
	h.click(2, g, "ks:new")
	h.text(2, g, "ремонт")
	id2 := h.list()[1].ID
	snapshot := h.list()
	for _, data := range []string{fmt.Sprintf("ks:setcategory:%d:invalid", id1), fmt.Sprintf("ks:setcity:%d:invalid", id1)} {
		if err := h.callback(1, p1, data); err == nil {
			t.Fatal("invalid catalog accepted")
		}
		if !reflect.DeepEqual(snapshot, h.list()) {
			t.Fatal("invalid catalog saved")
		}
	}
	h.click(1, p1, fmt.Sprintf("ks:price:%d", id1))
	for _, input := range []string{"-1-200", "200-100", "100.5-200", "abc", "9223372036854775808-"} {
		h.text(1, p1, input)
		if !reflect.DeepEqual(snapshot, h.list()) {
			t.Fatal("invalid prices saved", input)
		}
	}
	h.click(1, p1, fmt.Sprintf("ks:phrase:%d", id1))
	h.text(1, p1, " ")
	if !reflect.DeepEqual(snapshot, h.list()) {
		t.Fatal("blank edit saved")
	}
	// Two admins, two chats and two stable IDs retain their own pending action.
	h.click(1, g, fmt.Sprintf("ks:price:%d", id1))
	h.click(2, g, fmt.Sprintf("ks:phrase:%d", id2))
	h.text(2, p2, "wrong chat")
	h.text(3, g, "500-600")
	h.text(1, bad, "500-600")
	if !reflect.DeepEqual(snapshot, h.list()) {
		t.Fatal("cross-context input saved")
	}
	h.text(2, g, "мастер")
	h.text(1, g, "100-200")
	if h.search(id2).Phrase != "мастер" || *h.search(id1).PriceMin != 100 {
		t.Fatal("actions crossed admins/searches")
	}
	// Switching this admin's addressed search cancels its older pending action.
	h.click(1, g, fmt.Sprintf("ks:price:%d", id1))
	h.click(1, g, fmt.Sprintf("ks:price:%d", id2))
	h.text(1, g, "0-50")
	if *h.search(id1).PriceMax != 200 || *h.search(id2).PriceMax != 50 {
		t.Fatal("pending action crossed IDs")
	}
	// A second admin changing the same ID invalidates an earlier revision-bound input.
	h.click(1, g, fmt.Sprintf("ks:price:%d", id1))
	h.click(2, g, fmt.Sprintf("ks:setcity:%d:dushanbe", id1))
	current := h.search(id1)
	h.text(1, g, "900-1000")
	if !reflect.DeepEqual(current, h.search(id1)) {
		t.Fatal("stale input overwrote current search")
	}
	// Rental pending input cannot consume a search phrase; search pending cannot consume rental input.
	h.click(1, g, fmt.Sprintf("ks:price:%d", id1))
	h.click(1, g, "ks:rental")
	h.text(1, g, "900-1000")
	if !reflect.DeepEqual(current, h.search(id1)) {
		t.Fatal("navigation leaked pending search")
	}
	if err := telegram.ProcessKeywordUpdateForTest(h.bot, context.Background(), telegram.Update{CallbackQuery: &telegram.CallbackQuery{ID: "rental-price", From: telegram.User{ID: 1}, Message: &telegram.Message{Chat: g}, Data: "i:price"}}); err != nil {
		t.Fatal(err)
	}
	h.click(1, g, "ks:new")
	h.text(1, g, "диван")
	if len(h.list()) != 3 {
		t.Fatal("rental pending consumed search phrase")
	}
	h.click(1, g, "ks:new")
	h.text(1, g, "/cancel")
	h.text(1, g, "отменено")
	if len(h.list()) != 3 {
		t.Fatal("cancel retained pending")
	}
	prior := h.list()
	h.click(1, g, "ks:view:999999")
	h.click(1, g, "ks:enable:999999")
	h.click(1, g, "ks:setcity:999999:vose")
	if !reflect.DeepEqual(prior, h.list()) {
		t.Fatal("missing ID resurrected")
	}
	if !reflect.DeepEqual(before, h.snapshotRental()) {
		t.Fatal("rental data changed")
	}
	t.Log("GREEN harm AC001: invalid phrase/catalog/negative/inverted prices no writes; auth denied; admin/chat/ID/revision pending isolated; cancel and rental navigation cancel safely; missing ID no upsert")
}
