//go:build cgo
package app
import (
 "context"
 "fmt"
 "net/http"
 "net/http/httptest"
 "path/filepath"
 "testing"
 "time"
 "github.com/nicelight/somon-rent-watcher/internal/config"
 "github.com/nicelight/somon-rent-watcher/internal/filter"
 "github.com/nicelight/somon-rent-watcher/internal/model"
 "github.com/nicelight/somon-rent-watcher/internal/store"
)
func TestVerifierRentalSeenNoRepeatAcrossManagementAndRestart(t *testing.T){
 sends,details:=0,0
 src:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){if r.URL.Path=="/category/"{fmt.Fprint(w,testCategoryHTML([]int64{77001,77002}));return};details++;http.NotFound(w,r)}));defer src.Close()
 dst:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){sends++;fmt.Fprint(w,`{"ok":true,"result":{"message_id":1,"chat":{"id":-909,"type":"supergroup"}}}`)}));defer dst.Close()
 path:=filepath.Join(t.TempDir(),"search.db");db,e:=store.Open(path);if e!=nil{t.Fatal(e)};defer func(){db.Close()}()
 settings:=filter.DefaultSettings();settings.Enabled=true;raw,e:=settings.Encode();if e!=nil{t.Fatal(e)};if e=db.SaveSettingsJSON(raw);e!=nil{t.Fatal(e)}
 if e=db.MarkSeen([]int64{77001,77002},time.Now().Add(-time.Hour));e!=nil{t.Fatal(e)}
 if e=db.SetStates(map[string]string{"initialized":"1","telegram_offset":"719","previous_ordinary_ids":"[77001,77002]","last_successful_poll_at":time.Now().UTC().Format(time.RFC3339Nano)});e!=nil{t.Fatal(e)}
 cfg:=config.Config{CategoryURL:src.URL+"/category/",TelegramAPIBase:dst.URL,TelegramBotToken:"LOCAL",TelegramTargetChatID:-909,TelegramAdminUserIDs:[]int64{71},PollMin:time.Minute,PollMax:3*time.Minute,GapAfter:time.Hour,HTTPTimeout:time.Second}
 a,e:=New(cfg,db,nil);if e!=nil{t.Fatal(e)}
 if e=a.pollOnce(context.Background());e!=nil{t.Fatal(e)}
 s1,e:=a.CreateKeywordSearch(model.KeywordSearch{Phrase:"полка",CategoryKey:"all",CityKey:"country"});if e!=nil{t.Fatal(e)}
 s2,e:=a.CreateKeywordSearch(model.KeywordSearch{Phrase:"монтаж",CategoryKey:"services",CityKey:"vose"});if e!=nil{t.Fatal(e)}
 for _,s:=range []model.KeywordSearch{s1,s2}{if _,ok,e:=a.SetKeywordSearchEnabled(s.ID,true);e!=nil||!ok{t.Fatal(ok,e)}}
 if e=a.pollOnce(context.Background());e!=nil{t.Fatal(e)}
 if e=db.Close();e!=nil{t.Fatal(e)};db,e=store.Open(path);if e!=nil{t.Fatal(e)};a,e=New(cfg,db,nil);if e!=nil{t.Fatal(e)}
 if e=a.pollOnce(context.Background());e!=nil{t.Fatal(e)}
 if sends!=0||details!=0{t.Fatalf("rental repeated existing ads sends=%d details=%d",sends,details)}
 if n,e:=db.CountSeen();e!=nil||n!=2{t.Fatal(n,e)}
 if raw2,ok,e:=db.LoadSettingsJSON();e!=nil||!ok||raw2!=raw{t.Fatal("rental settings changed",raw2,e)}
 if offset,ok,e:=db.GetState("telegram_offset");e!=nil||!ok||offset!="719"{t.Fatal(offset,ok,e)}
 t.Log("Independent AC007: existing rental polling before/after two search writes and DB restart never re-requests details or repeats seen ads; settings/seen/offset preserved PASS")
}
