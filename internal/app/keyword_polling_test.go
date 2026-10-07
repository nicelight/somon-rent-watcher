//go:build cgo

package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/store"
)

type keywordTrace struct {
	path, phrase string
	at           time.Time
}

func (r keywordTrace) String() string {
	return fmt.Sprintf("%s?q=%s@%s", r.path, r.phrase, r.at.Format("15:04:05.000000"))
}

type keywordPollingHarness struct {
	t                                             *testing.T
	mu                                            sync.Mutex
	cfg                                           config.Config
	db                                            *store.DB
	a                                             *App
	source, tg                                    *httptest.Server
	ids                                           []int64
	price                                         int
	sourceStatus, detailStatus, telegramStatus    int
	sourceBroken, detailBroken, telegramAmbiguous bool
	rentalNew                                     bool
	trace                                         []keywordTrace
	captions, photos, buttons                     []string
	detailBlocked                                 bool
	detailHook                                    func()
	sendHook                                      func()
	detailPrice                                   int
	inFlight, maxInFlight                         int
}

// Test-only access keeps the accepted adapter API unchanged while routing all
// native HTTPS URLs to a disposable httptest endpoint without external traffic.
type keywordTestTransport struct{ base *url.URL }

func (tr keywordTestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = tr.base.Scheme
	u.Host = tr.base.Host
	c.URL = &u
	c.Host = tr.base.Host
	return http.DefaultTransport.RoundTrip(c)
}
func (h *keywordPollingHarness) newApp() {
	var err error
	h.a, err = New(h.cfg, h.db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		h.t.Fatal(err)
	}
	field := reflect.ValueOf(h.a.somon).Elem().FieldByName("httpClient")
	client := *(**http.Client)(unsafe.Pointer(field.UnsafeAddr()))
	endpoint, _ := url.Parse(h.source.URL)
	client.Transport = keywordTestTransport{endpoint}
}
func newKeywordPollingHarness(t *testing.T, cap int) *keywordPollingHarness {
	h := &keywordPollingHarness{t: t, ids: []int64{1001}, price: 150}
	h.source = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.mu.Lock()
		h.inFlight++
		if h.inFlight > h.maxInFlight {
			h.maxInFlight = h.inFlight
		}
		h.trace = append(h.trace, keywordTrace{r.URL.Path, r.URL.Query().Get("q"), time.Now()})
		ids := append([]int64(nil), h.ids...)
		price, status, broken := h.price, h.sourceStatus, h.sourceBroken
		detailStatus, detailBroken := h.detailStatus, h.detailBroken
		detailBlocked := h.detailBlocked
		rentalNew := h.rentalNew
		hook := h.detailHook
		detailPrice := h.detailPrice
		h.mu.Unlock()
		defer func() { h.mu.Lock(); h.inFlight--; h.mu.Unlock() }()
		switch {
		case r.URL.Path == "/category/":
			if rentalNew {
				fmt.Fprint(w, testCategoryHTML([]int64{9001, 9002}))
			} else {
				fmt.Fprint(w, testCategoryHTML([]int64{9001}))
			}
		case strings.HasPrefix(r.URL.Path, "/adv/"):
			if hook != nil {
				hook()
			}
			if detailPrice != 0 {
				price = detailPrice
			}
			if detailStatus != 0 {
				w.WriteHeader(detailStatus)
				return
			}
			if detailBlocked {
				fmt.Fprint(w, "<html><body>Access denied</body></html>")
				return
			}
			if detailBroken {
				fmt.Fprint(w, "broken detail")
				return
			}
			if strings.Contains(r.URL.Path, "9002_") {
				fmt.Fprint(w, `<html><body><h1>2-комн. квартира, 3 этаж</h1><div class="price">4 500 c.</div><div>Этаж: 3</div><h2>Описание</h2><p>Долгосрочная аренда</p></body></html>`)
				return
			}
			fmt.Fprintf(w, `<html><body><h1>Стол деревянный</h1><div class="price">%d c.</div><h2>Описание</h2><p>Мебель</p></body></html>`, price)
		default:
			if status != 0 {
				w.Header().Set("Retry-After", "2")
				w.WriteHeader(status)
				return
			}
			if broken {
				fmt.Fprint(w, "broken search")
				return
			}
			fmt.Fprintf(w, "<html><body><h1>Результаты %d</h1>", len(ids))
			for _, id := range ids {
				fmt.Fprintf(w, `<article><a href="/adv/%d_table/">Стол деревянный</a><div class="price">%d c.</div><div>Душанбе</div><img src="https://images.example/%d.jpg"></article>`, id, price, id)
			}
			fmt.Fprint(w, "</body></html>")
		}
	}))
	h.tg = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		h.mu.Lock()
		status, ambiguous, hook := h.telegramStatus, h.telegramAmbiguous, h.sendHook
		if r.Form.Get("chat_id") == "-200" {
			h.captions = append(h.captions, r.Form.Get("caption")+r.Form.Get("text"))
			h.photos = append(h.photos, r.Form.Get("photo"))
			h.buttons = append(h.buttons, r.Form.Get("reply_markup"))
		}
		h.mu.Unlock()
		if hook != nil {
			hook()
		}
		if status != 0 {
			w.WriteHeader(status)
			fmt.Fprint(w, `{"ok":false,"error_code":500,"description":"failed"}`)
			return
		}
		if ambiguous {
			fmt.Fprint(w, `{"ok":`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":1}}`)
	}))
	var err error
	path := filepath.Join(t.TempDir(), "search.db")
	h.db, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	h.cfg = config.Config{DBPath: path, CategoryURL: h.source.URL + "/category/", TelegramBotToken: "local", TelegramAPIBase: h.tg.URL, TelegramTargetChatID: -200, PollMin: time.Hour, PollMax: time.Hour, BlockBackoff: time.Hour, RequestDelay: 5 * time.Millisecond, HTTPTimeout: time.Second, MinCards: 1, MaxDetailsPerPoll: cap, MaxBodyBytes: 1 << 20}
	h.newApp()
	t.Cleanup(func() { h.db.Close(); h.source.Close(); h.tg.Close() })
	return h
}
func (h *keywordPollingHarness) search(phrase string, max *int) model.KeywordSearch {
	s, err := h.a.CreateKeywordSearch(model.KeywordSearch{Phrase: phrase, CategoryKey: "all", CityKey: "dushanbe", PriceMax: max})
	if err != nil {
		h.t.Fatal(err)
	}
	s, ok, err := h.a.SetKeywordSearchEnabled(s.ID, true)
	if err != nil || !ok {
		h.t.Fatal(ok, err)
	}
	return s
}
func (h *keywordPollingHarness) poll() {
	h.t.Helper()
	if err := h.a.pollOnce(context.Background()); err != nil {
		h.t.Fatal(err)
	}
}
func (h *keywordPollingHarness) history(s model.KeywordSearch) map[int64]model.KeywordAdState {
	h.t.Helper()
	v, err := h.db.KeywordAdStates(s.ID)
	if err != nil {
		h.t.Fatal(err)
	}
	return v
}
func (h *keywordPollingHarness) reopen() {
	h.t.Helper()
	if err := h.db.Close(); err != nil {
		h.t.Fatal(err)
	}
	var err error
	h.db, err = store.Open(h.cfg.DBPath)
	if err != nil {
		h.t.Fatal(err)
	}
	h.newApp()
}

func TestKeywordPollingFirstNewIndependentRestartRevision(t *testing.T) {
	h := newKeywordPollingHarness(t, 10)
	max := 200
	s1 := h.search("стол", &max)
	s2 := h.search("деревянный стол", &max)
	if err := h.db.MarkSeen([]int64{1001}, time.Now()); err != nil {
		t.Fatal(err)
	}
	h.poll()
	for _, s := range []model.KeywordSearch{s1, s2} {
		if !h.history(s)[1001].Delivered {
			t.Fatal("first existing/rental-seen match missing", s)
		}
	}
	if len(h.captions) != 2 {
		t.Fatal(h.captions)
	}
	for i, c := range h.captions {
		for _, part := range []string{[]string{s1.Phrase, s2.Phrase}[i], "Стол деревянный", "150 c.", "Душанбе"} {
			if !strings.Contains(c, part) {
				t.Fatal("payload missing", part, c)
			}
		}
		if !strings.Contains(h.photos[i], "1001.jpg") || !strings.Contains(h.buttons[i], "1001_table") {
			t.Fatal("photo/link missing")
		}
	}
	h.mu.Lock()
	h.ids = []int64{1001, 1002}
	h.price = 250
	h.mu.Unlock()
	h.poll()
	for _, s := range []model.KeywordSearch{s1, s2} {
		v := h.history(s)
		if !v[1001].Delivered || v[1002].Delivered || v[1002].EvaluatedRevision != s.Revision {
			t.Fatal(v)
		}
	}
	h.mu.Lock()
	h.ids = []int64{1001, 1002, 1003}
	h.price = 150
	h.mu.Unlock()
	h.poll()
	if len(h.captions) != 4 {
		t.Fatal("new delivery or revision rejection broken", h.captions)
	}
	h.reopen()
	h.mu.Lock()
	h.price = 100
	h.mu.Unlock()
	h.poll()
	if len(h.captions) != 4 {
		t.Fatal("restart/price drop repeated")
	}
	raised := 300
	for _, s := range []model.KeywordSearch{s1, s2} {
		current, _, _ := h.a.KeywordSearch(s.ID)
		current.PriceMax = &raised
		updated, ok, err := h.a.UpdateKeywordSearch(current)
		if err != nil || !ok {
			t.Fatal(ok, err)
		}
		t.Logf("search=%d revision %d->%d", s.ID, s.Revision, updated.Revision)
	}
	h.poll()
	if len(h.captions) != 6 {
		t.Fatal("rejected ID not reevaluated", h.captions)
	}
	for _, s := range []model.KeywordSearch{s1, s2} {
		v := h.history(s)
		if len(v) != 3 || !v[1001].Delivered || !v[1002].Delivered || !v[1003].Delivered {
			t.Fatal(v)
		}
		t.Logf("AC005 delivered history search=%d: %+v", s.ID, v)
	}
	h.reopen()
	h.poll()
	if len(h.captions) != 6 {
		t.Fatal("delivered repeated after edit/restart")
	}
}

func TestKeywordPollingSharedCapRotationAndDelay(t *testing.T) {
	h := newKeywordPollingHarness(t, 1)
	s1 := h.search("стол", nil)
	s2 := h.search("стул", nil)
	h.poll()
	if !h.history(s1)[1001].Delivered || len(h.history(s2)) != 0 {
		t.Fatal("cap deferral wrong")
	}
	settings, _ := h.a.LoadSettings()
	settings.Enabled = true
	if err := h.a.SaveSettings(settings); err != nil {
		t.Fatal(err)
	}
	h.mu.Lock()
	h.rentalNew = true
	h.mu.Unlock()
	for cycle := 0; cycle < 3; cycle++ {
		h.mu.Lock()
		before := len(h.trace)
		h.mu.Unlock()
		h.poll()
		h.mu.Lock()
		trace := append([]keywordTrace(nil), h.trace[before:]...)
		h.mu.Unlock()
		details := 0
		for _, r := range trace {
			if strings.HasPrefix(r.path, "/adv/") {
				details++
			}
		}
		if details > 1 {
			t.Fatal("shared cap exceeded", trace)
		}
		t.Logf("AC006 cycle=%d details=%d trace=%+v", cycle, details, trace)
	}
	if !h.history(s2)[1001].Delivered {
		t.Fatal("round robin starved second search")
	}
	seen, err := h.db.SeenIDs([]int64{9002})
	if err != nil || !seen[9002] {
		t.Fatal("rental starved", seen, err)
	}
	h.mu.Lock()
	trace := append([]keywordTrace(nil), h.trace...)
	h.mu.Unlock()
	if h.maxInFlight != 1 {
		t.Fatal("HTTP requests overlapped", h.maxInFlight)
	}
	t.Logf("AC006 maximum in-flight=%d, configured shared delay=%s", h.maxInFlight, h.cfg.RequestDelay)
	minDelay := time.Hour
	for i := 1; i < len(trace); i++ {
		if d := trace[i].at.Sub(trace[i-1].at); d < minDelay {
			minDelay = d
		}
		if d := trace[i].at.Sub(trace[i-1].at); d < 4*time.Millisecond {
			t.Fatal("shared delay missing", d)
		}
	}
	t.Logf("AC006 observed minimum HTTP start spacing=%s (target configured 5ms, tolerance 1ms)", minDelay)
}

func TestKeywordPollingFailuresRetryWithoutHistory(t *testing.T) {
	for _, failure := range []string{"source-status", "source-parse", "detail-status", "detail-parse", "telegram-error", "telegram-ambiguous"} {
		t.Run(failure, func(t *testing.T) {
			h := newKeywordPollingHarness(t, 2)
			s := h.search("стол", nil)
			h.mu.Lock()
			switch failure {
			case "source-status":
				h.sourceStatus = 500
			case "source-parse":
				h.sourceBroken = true
			case "detail-status":
				h.detailStatus = 500
			case "detail-parse":
				h.detailBroken = true
			case "telegram-error":
				h.telegramStatus = 500
			case "telegram-ambiguous":
				h.telegramAmbiguous = true
			}
			h.mu.Unlock()
			err := h.a.pollOnce(context.Background())
			if strings.HasPrefix(failure, "source") && err == nil {
				t.Fatal("source failure ignored")
			}
			if len(h.history(s)) != 0 {
				t.Fatal("failure advanced history", h.history(s))
			}
			h.mu.Lock()
			failedTrace := append([]keywordTrace(nil), h.trace...)
			attempts := len(h.captions)
			h.sourceStatus = 0
			h.sourceBroken = false
			h.detailStatus = 0
			h.detailBroken = false
			h.telegramStatus = 0
			h.telegramAmbiguous = false
			h.mu.Unlock()
			h.reopen()
			h.poll()
			if !h.history(s)[1001].Delivered {
				t.Fatal("failed candidate did not retry")
			}
			t.Logf("AC006 %s failure trace=%+v failed sends=%d; after restart delivered=%+v", failure, failedTrace, attempts, h.history(s))
		})
	}
}

func TestKeywordPollingSharedBlockedBackoffAndSingleFlight(t *testing.T) {
	for _, code := range []int{403, 429, 0} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			h := newKeywordPollingHarness(t, 2)
			s := h.search("стол", nil)
			h.search("стул", nil)
			h.mu.Lock()
			h.sourceStatus = code
			if code == 0 {
				h.detailBlocked = true
			}
			h.mu.Unlock()
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() { done <- h.a.pollLoop(ctx) }()
			defer func() { cancel(); <-done }()
			deadline := time.After(2 * time.Second)
			for {
				status := h.a.RuntimeStatus()
				if status.BackoffUntil.After(time.Now()) {
					if code != 0 && !strings.Contains(status.Mode, fmt.Sprint(code)) {
						t.Fatal(status)
					}
					break
				}
				select {
				case <-deadline:
					t.Fatal("backoff not entered")
				default:
					time.Sleep(time.Millisecond)
				}
			}
			if err := h.a.RequestPollNow(-200, 1); err == nil {
				t.Fatal("manual bypassed shared block")
			}
			if len(h.history(s)) != 0 {
				t.Fatal("blocked advanced history")
			}
			h.mu.Lock()
			trace := append([]keywordTrace(nil), h.trace...)
			h.mu.Unlock()
			want := 2
			if code == 0 {
				want = 3
			}
			if len(trace) != want {
				t.Fatal("block did not stop remaining searches", trace)
			}
			t.Logf("AC006 HTTP%d shared backoff, manual rejected, trace=%+v", code, trace)
		})
	}
	h := newKeywordPollingHarness(t, 1)
	h.search("стол", nil)
	h.poll()
	if err := h.a.RequestPollNow(-200, 1); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if err := h.a.RequestPollNow(-200, int64(i+2)); err == nil {
			t.Fatal("duplicate accepted")
		}
	}
	if len(h.a.pollNowCh) != 1 {
		t.Fatal("single flight broken")
	}
}

func TestKeywordPollingStaleRejectionAndStaleDelivery(t *testing.T) {
	for _, reject := range []bool{true, false} {
		t.Run(fmt.Sprint(reject), func(t *testing.T) {
			h := newKeywordPollingHarness(t, 2)
			max := 200
			s := h.search("стол", &max)
			entered, release := make(chan struct{}), make(chan struct{})
			once := sync.Once{}
			h.mu.Lock()
			h.detailHook = func() { once.Do(func() { close(entered); <-release }) }
			if reject {
				h.detailPrice = 250
			}
			h.mu.Unlock()
			done := make(chan error, 1)
			go func() { done <- h.a.pollOnce(context.Background()) }()
			<-entered
			current, _, _ := h.a.KeywordSearch(s.ID)
			newMax := 100
			if reject {
				newMax = 300
			}
			current.PriceMax = &newMax
			updated, ok, err := h.a.UpdateKeywordSearch(current)
			if err != nil || !ok {
				t.Fatal(ok, err)
			}
			// Detail is already based on revision 2. Set the detail response above the old
			// budget for rejection, or keep it eligible for the stale-send case.

			close(release)
			if err := <-done; err != nil {
				t.Fatal(err)
			}
			if len(h.history(s)) != 0 || len(h.captions) != 0 {
				t.Fatal("stale work suppressed new revision or sent", h.history(s), h.captions)
			}
			h.mu.Lock()
			h.detailHook = nil
			if !reject {
				h.price = 90
			}
			h.mu.Unlock()
			h.poll()
			if !h.history(updated)[1001].Delivered {
				t.Fatal("new revision not evaluated")
			}
		})
	}
}

func TestKeywordPollingSendSuccessStoreFailurePossibleRepeat(t *testing.T) {
	h := newKeywordPollingHarness(t, 2)
	s := h.search("стол", nil)
	h.mu.Lock()
	h.sendHook = func() { h.db.Close() }
	h.mu.Unlock()
	if err := h.a.pollOnce(context.Background()); err == nil {
		t.Fatal("store failure missing")
	}
	h.mu.Lock()
	h.sendHook = nil
	h.mu.Unlock()
	h.reopen()
	if len(h.history(s)) != 0 {
		t.Fatal("unwritten success became durable")
	}
	h.poll()
	if len(h.captions) != 2 || !h.history(s)[1001].Delivered {
		t.Fatal("accepted ambiguity/retry wrong", h.captions, h.history(s))
	}
	t.Log("AC006 send success + write failure: possible repeat retained (2 sends), second success durable")
}

func TestKeywordPollingManualSingleFlightWhileRunning(t *testing.T) {
	h := newKeywordPollingHarness(t, 1)
	h.search("стол", nil)
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	h.mu.Lock()
	h.detailHook = func() { once.Do(func() { close(entered); <-release }) }
	h.mu.Unlock()
	if err := h.a.RequestPollNow(-200, 1); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- h.a.pollLoop(ctx) }()
	select {
	case <-entered:
	case <-time.After(time.Second):
		cancel()
		close(release)
		<-done
		t.Fatal("detail request never started")
	}
	for i := 0; i < 5; i++ {
		if err := h.a.RequestPollNow(-200, int64(i+2)); err == nil {
			cancel()
			close(release)
			<-done
			t.Fatal("running duplicate request accepted")
		}
	}
	close(release)
	deadline := time.After(time.Second)
	for {
		status := h.a.RuntimeStatus()
		if !status.NextPoll.IsZero() {
			if status.Mode != "норма" {
				cancel()
				<-done
				t.Fatal("active searches shown paused", status)
			}
			break
		}
		select {
		case <-deadline:
			cancel()
			<-done
			t.Fatal("poll did not finish")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	cancel()
	<-done
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.maxInFlight != 1 {
		t.Fatal("scheduler HTTP overlap", h.maxInFlight)
	}
	t.Logf("AC006 running manual duplicates=5 rejected; max in-flight=%d", h.maxInFlight)
}
