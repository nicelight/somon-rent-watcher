//go:build cgo
package app
import (
 "context"
 "fmt"
 "io"
 "log/slog"
 "net/http"
 "net/http/httptest"
 "path/filepath"
 "testing"
 "time"
 "github.com/nicelight/somon-rent-watcher/internal/config"
 "github.com/nicelight/somon-rent-watcher/internal/model"
 "github.com/nicelight/somon-rent-watcher/internal/store"
)
func TestKeywordPollingBaseline(t *testing.T) {
 group:=0; requests:=0
 src:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){requests++;fmt.Fprint(w,testCategoryHTML([]int64{9001}))}));defer src.Close()
 tg:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){r.ParseForm();if r.Form.Get("chat_id")=="-200"{group++};fmt.Fprint(w,`{"ok":true,"result":{"message_id":1}}`)}));defer tg.Close()
 db,e:=store.Open(filepath.Join(t.TempDir(),"search.db"));if e!=nil{t.Fatal(e)};defer db.Close()
 a,e:=New(config.Config{CategoryURL:src.URL+"/category/",TelegramBotToken:"test",TelegramAPIBase:tg.URL,TelegramTargetChatID:-200,PollMin:time.Minute,PollMax:time.Minute,MinCards:1,MaxDetailsPerPoll:1,MaxBodyBytes:1<<20,HTTPTimeout:time.Second},db,slog.New(slog.NewTextHandler(io.Discard,nil)));if e!=nil{t.Fatal(e)}
 for _,phrase:=range []string{"стол","стул"}{s,e:=a.CreateKeywordSearch(model.KeywordSearch{Phrase:phrase,CategoryKey:"all",CityKey:"country"});if e!=nil{t.Fatal(e)};if _,ok,e:=a.SetKeywordSearchEnabled(s.ID,true);e!=nil||!ok{t.Fatal(ok,e)}}
 if e:=a.pollOnce(context.Background());e!=nil{t.Fatal(e)};if e:=a.pollOnce(context.Background());e!=nil{t.Fatal(e)}
 t.Logf("AC005/AC006 baseline: enabled searches=2, source requests=%d (rental only), keyword/group deliveries=%d; keyword history schema/API absent in source snapshot",requests,group)
 if group!=0||requests!=2{t.Fatal("baseline unexpectedly integrates searches")}
}
