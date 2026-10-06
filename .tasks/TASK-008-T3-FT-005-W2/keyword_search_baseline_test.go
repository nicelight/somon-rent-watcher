//go:build cgo

package telegram_test

import (
 "bytes"
 "context"
 "fmt"
 "io"
 "log/slog"
 "net/http"
 "net/http/httptest"
 "os"
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

func TestKeywordBaselineRentalAndAuthPreserved(t *testing.T) {
 path:=filepath.Join(t.TempDir(),"search.db")
 db,err:=store.Open(path);if err!=nil{t.Fatal(err)}
 settings:=filter.DefaultSettings();settings.Enabled=true;max:=6000;settings.PriceMax=&max
 raw,err:=settings.Encode();if err!=nil{t.Fatal(err)}
 if err=db.SaveSettingsJSON(raw);err!=nil{t.Fatal(err)}
 if err=db.MarkSeen([]int64{101,202},time.Date(2026,10,6,12,0,0,0,time.UTC));err!=nil{t.Fatal(err)}
 states:=map[string]string{"initialized":"1","previous_ordinary_ids":"[101,202]","last_successful_poll_at":"2026-10-06T12:00:00Z","telegram_offset":"88"}
 if err=db.SetStates(states);err!=nil{t.Fatal(err)}
 if err=db.Close();err!=nil{t.Fatal(err)}
 before,err:=os.ReadFile(path);if err!=nil{t.Fatal(err)}
 db,err=store.Open(path);if err!=nil{t.Fatal(err)}
 var methods []string
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){methods=append(methods,strings.TrimPrefix(r.URL.Path,"/botTOKEN/"));w.Header().Set("Content-Type","application/json");if strings.HasSuffix(r.URL.Path,"answerCallbackQuery"){fmt.Fprint(w,`{"ok":true,"result":true}`)}else{fmt.Fprint(w,`{"ok":true,"result":{"message_id":1,"chat":{"id":-100,"type":"supergroup"}}}`)}}));defer server.Close()
 a,err:=app.New(config.Config{PollMin:time.Minute,PollMax:2*time.Minute,TelegramAPIBase:server.URL,TelegramBotToken:"TOKEN",TelegramAdminUserIDs:[]int64{1,2},TelegramTargetChatID:-100},db,slog.New(slog.NewTextHandler(io.Discard,nil)));if err!=nil{t.Fatal(err)}
 bot:=telegram.NewBot(telegram.NewClient(server.URL,"TOKEN"),a,[]int64{1,2},-100,nil)
 for _,q:=range []telegram.CallbackQuery{{ID:"nonadmin",From:telegram.User{ID:3},Message:&telegram.Message{Chat:telegram.Chat{ID:-100,Type:"supergroup"}},Data:"ks:new"},{ID:"wrongchat",From:telegram.User{ID:1},Message:&telegram.Message{Chat:telegram.Chat{ID:-200,Type:"supergroup"}},Data:"ks:new"}}{if err:=telegram.ProcessKeywordUpdateForTest(bot,context.Background(),telegram.Update{CallbackQuery:&q});err!=nil{t.Fatal(err)}}
 if !reflect.DeepEqual(methods,[]string{"answerCallbackQuery","answerCallbackQuery"}){t.Fatal(methods)}
 methods=nil
 ad:=model.Ad{Card:model.Card{ID:303,Title:"Квартира",URL:"https://somon.tj/adv/303_x/"}}
 if err=bot.SendAd(context.Background(),ad);err!=nil{t.Fatal(err)}
 if !reflect.DeepEqual(methods,[]string{"sendMessage"}){t.Fatal(methods)}
 got,ok,err:=db.LoadSettingsJSON();if err!=nil||!ok||got!=raw{t.Fatalf("settings changed %v",err)}
 seen,err:=db.SeenIDs([]int64{101,202,303});if err!=nil||!reflect.DeepEqual(seen,map[int64]bool{101:true,202:true}){t.Fatal(seen,err)}
 for k,want:=range states{got,ok,err:=db.GetState(k);if err!=nil||!ok||got!=want{t.Fatal(k,got,err)}}
 if err=db.Close();err!=nil{t.Fatal(err)}
 after,err:=os.ReadFile(path);if err!=nil{t.Fatal(err)}
 if !bytes.Equal(before,after){t.Fatal("baseline rental SQLite bytes changed across reopen/auth/delivery")}
 t.Log("initial GREEN AC007 exact SQLite bytes + settings/seen/state/offset preserved across reopen; rental sendMessage trace unchanged; AC001 auth rejected with acknowledgements only")
}

func TestKeywordCreateBaseline(t *testing.T){
 path:=filepath.Join(t.TempDir(),"search.db");db,err:=store.Open(path);if err!=nil{t.Fatal(err)}
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","application/json");fmt.Fprint(w,`{"ok":true,"result":true}`)}));defer server.Close()
 a,err:=app.New(config.Config{PollMin:time.Minute,PollMax:2*time.Minute},db,nil);if err!=nil{t.Fatal(err)}
 bot:=telegram.NewBot(telegram.NewClient(server.URL,"TOKEN"),a,[]int64{1},-100,nil)
 err=telegram.ProcessKeywordUpdateForTest(bot,context.Background(),telegram.Update{CallbackQuery:&telegram.CallbackQuery{ID:"create",From:telegram.User{ID:1},Message:&telegram.Message{Chat:telegram.Chat{ID:1,Type:"private"}},Data:"ks:new"}})
 if err==nil{t.Fatal("baseline unexpectedly supports route")}
 if !strings.Contains(err.Error(),"неизвестная кнопка"){t.Fatalf("unrelated failure: %v",err)}
 if e:=db.Close();e!=nil{t.Fatal(e)}
 raw,e:=os.ReadFile(path);if e!=nil{t.Fatal(e)}
 if bytes.Contains(raw,[]byte("search_monitors")){t.Fatal("baseline unexpectedly has durable search table")}
 t.Fatalf("behavioral RED AC001: authorized create route rejected: %v; no durable search table or search created",err)
}
