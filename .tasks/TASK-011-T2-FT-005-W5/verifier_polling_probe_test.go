//go:build cgo

package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestVerifierDetailFailureLifecycle(t *testing.T) {
	for _, blocked := range []bool{false, true} {
		t.Run(fmt.Sprintf("blocked-%v", blocked), func(t *testing.T) {
			h := newKeywordPollingHarness(t, 5)
			s := h.search("стол", nil)
			var bodyMu sync.Mutex
			body := `<html><body><h1>404 missing page</h1></body></html>`
			if blocked { body = `<html><body><h1>Стол</h1><div class="price">150 c.</div><p>Слишком много запросов</p></body></html>` }
			details := 0
			original := h.source.Config.Handler
			h.source.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.HasPrefix(r.URL.Path, "/adv/1001_") { original.ServeHTTP(w, r); return }
				bodyMu.Lock(); current := body; details++; bodyMu.Unlock()
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, current)
			})
			if blocked {
				second := h.search("стул", nil)
				ctx, cancel := context.WithCancel(context.Background())
				done := make(chan error, 1)
				go func() { done <- h.a.pollLoop(ctx) }()
				deadline := time.Now().Add(2 * time.Second)
				for !h.a.RuntimeStatus().BackoffUntil.After(time.Now()) && time.Now().Before(deadline) { time.Sleep(time.Millisecond) }
				status := h.a.RuntimeStatus()
				cancel(); <-done
				if !status.BackoffUntil.After(time.Now()) || !strings.Contains(status.Mode, "backoff") { t.Fatalf("HTTP200 block failed shared backoff: %+v", status) }
				if err := h.a.RequestPollNow(-200, 1); err == nil { t.Fatal("manual request bypassed block") }
				if len(h.history(second)) != 0 { t.Fatal("later monitor advanced history after block") }
				h.mu.Lock()
				for _, request := range h.trace { if request.phrase == "стул" { t.Error("remaining monitor fetched despite shared block") } }
				h.mu.Unlock()
				t.Logf("FT-005-AC-008 HTTP200 body entered shared %s; remaining monitor/manual stopped", status.Mode)
				// Keep recovery proof about the failed ID in the first monitor.
				if _, _, err := h.a.SetKeywordSearchEnabled(second.ID, false); err != nil { t.Fatal(err) }
			} else { h.poll() }
			if len(h.history(s)) != 0 { t.Fatal("detail failure advanced SQLite history", h.history(s)) }
			h.mu.Lock(); sends := len(h.captions); h.mu.Unlock()
			if sends != 0 { t.Fatal("detail failure delivered to Telegram", sends) }
			bodyMu.Lock()
			if details != 1 { t.Fatalf("initial detail requests=%d", details) }
			body = `<html><body><div hidden>Access denied</div><h1>Стол деревянный</h1><h2>Описание</h2><p>Мебель</p></body></html>`
			bodyMu.Unlock()
			h.reopen()
			if len(h.history(s)) != 0 { t.Fatal("failure persisted a row after reopen") }
			h.poll()
			history := h.history(s)
			h.mu.Lock(); sends = len(h.captions); caption := strings.Join(h.captions, " "); h.mu.Unlock()
			bodyMu.Lock(); calls := details; bodyMu.Unlock()
			if !history[1001].Delivered || sends != 1 || calls != 2 || !strings.Contains(caption, "150") { t.Fatalf("retry outcome history=%+v sends=%d calls=%d caption=%s", history, sends, calls, caption) }
			h.reopen(); h.poll()
			h.mu.Lock(); sends = len(h.captions); h.mu.Unlock()
			if sends != 1 { t.Fatal("retry success repeated after another reopen", sends) }
			t.Logf("FT-005-AC-008 failed sends=0 rows=0; reopen+valid sparse retry sends=%d detail_calls=%d delivered=%+v; next reopen skips delivered", sends, calls, history)
		})
	}
}
