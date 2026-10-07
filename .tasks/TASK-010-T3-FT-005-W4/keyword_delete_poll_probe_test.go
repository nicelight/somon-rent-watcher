//go:build cgo

package app

import (
 "context"
 "fmt"
 "io"
 "log/slog"
 "net/http"
 "net/http/httptest"
 "strings"
 "sync"
 "testing"
 "time"

 "github.com/nicelight/somon-rent-watcher/internal/model"
 "github.com/nicelight/somon-rent-watcher/internal/telegram"
)

// Exercise the real authorized Telegram callback through its local long-poll
// entrypoint. Returning to getUpdates establishes that the mutation completed.
func deleteViaLocalBot(t *testing.T,h *keywordPollingHarness,s model.KeywordSearch) {
 t.Helper();ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second);defer cancel()
 completed:=make(chan struct{});var once sync.Once;sent:=false
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  switch {
  case strings.HasSuffix(r.URL.Path,"/getUpdates"):
   if !sent{sent=true;fmt.Fprintf(w,`{"ok":true,"result":[{"update_id":1,"callback_query":{"id":"delete","from":{"id":1},"message":{"message_id":7,"chat":{"id":-200,"type":"supergroup"}},"data":"ks:delete:%d"}}]}`,s.ID);return}
   once.Do(func(){close(completed)});<-ctx.Done()
  case strings.HasSuffix(r.URL.Path,"/deleteWebhook"),strings.HasSuffix(r.URL.Path,"/answerCallbackQuery"):fmt.Fprint(w,`{"ok":true,"result":true}`)
  default:fmt.Fprint(w,`{"ok":true,"result":{"message_id":8}}`)
  }
 }));defer server.Close()
 bot:=telegram.NewBot(telegram.NewClient(server.URL,"local"),h.a,[]int64{1},-200,slog.New(slog.NewTextHandler(io.Discard,nil)))
 done:=make(chan error,1);go func(){done<-bot.Run(ctx)}()
 select{case <-completed:case <-ctx.Done():cancel();<-done;t.Fatal("delete callback did not finish")};cancel();<-done
}
func TestKeywordDeleteBeforeSendStart(t *testing.T){
 h:=newKeywordPollingHarness(t,2);s:=h.search("стол",nil)
 entered,release:=make(chan struct{}),make(chan struct{});var once sync.Once
 h.mu.Lock();h.detailHook=func(){once.Do(func(){close(entered);<-release})};h.mu.Unlock()
 ctx,cancel:=context.WithTimeout(context.Background(),4*time.Second);defer cancel();done:=make(chan error,1)
 go func(){done<-h.a.pollOnce(ctx)}()
 select{case <-entered:case <-ctx.Done():close(release);<-done;t.Fatal("detail barrier not reached")}
 deleteViaLocalBot(t,h,s)
 _,found,e:=h.a.KeywordSearch(s.ID)
 close(release);if err:=<-done;err!=nil{t.Fatal(err)}
 h.mu.Lock();sends:=len(h.captions);trace:=append([]keywordTrace(nil),h.trace...);h.mu.Unlock()
 history:=h.history(s)
 t.Logf("FT-005-AC-002 ordered detail-start → delete-callback-complete → detail-release → send starts=%d; exists=%v history=%v source_trace=%v",sends,found,history,trace)
 if e!=nil||found||sends!=0||len(history)!=0{t.Fatal("deleted stale work delivered or resurrected",e)}
}
