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

type verifierRequest struct { Method string; Form url.Values }
func TestVerifierSearchManagementAndHarm(t *testing.T) {
 ctx:=context.Background()
 var requests []verifierRequest
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if err:=r.ParseForm();err!=nil { t.Error(err) }
  method:=strings.TrimPrefix(r.URL.Path,"/botLOCAL/")
  requests=append(requests,verifierRequest{method,r.Form})
  w.Header().Set("Content-Type","application/json")
  if method=="answerCallbackQuery" { fmt.Fprint(w,`{"ok":true,"result":true}`) } else { fmt.Fprint(w,`{"ok":true,"result":{"message_id":17,"chat":{"id":-909,"type":"supergroup"}}}`) }
 }))
 defer server.Close()
 path:=filepath.Join(t.TempDir(),"search.db")
 db,err:=store.Open(path);if err!=nil {t.Fatal(err)}
 defer func(){ db.Close() }()
 settings:=filter.DefaultSettings();settings.Enabled=true
 rentalMin,rentalMax:=3700,6500;settings.PriceMin=&rentalMin;settings.PriceMax=&rentalMax
 settings.NegativeWords=[]string{"посуточно","подселение"}
 raw,err:=settings.Encode();if err!=nil {t.Fatal(err)}
 if err=db.SaveSettingsJSON(raw);err!=nil {t.Fatal(err)}
 if err=db.MarkSeen([]int64{77001,77002},time.Date(2026,10,6,10,0,0,0,time.UTC));err!=nil{t.Fatal(err)}
 state:=map[string]string{"initialized":"1","previous_ordinary_ids":"[77001,77002]","last_successful_poll_at":"2026-10-06T10:00:00Z","telegram_offset":"719","private_fixture":"unrelated-preserve"}
 if err=db.SetStates(state);err!=nil{t.Fatal(err)}
 cfg:=config.Config{TelegramAPIBase:server.URL,TelegramBotToken:"LOCAL",TelegramTargetChatID:-909,TelegramAdminUserIDs:[]int64{71,72},PollMin:time.Minute,PollMax:3*time.Minute}
 a,err:=app.New(cfg,db,slog.New(slog.NewTextHandler(io.Discard,nil)));if err!=nil{t.Fatal(err)}
 bot:=telegram.NewBot(telegram.NewClient(server.URL,"LOCAL"),a,[]int64{71,72},-909,nil)
 private:=telegram.Chat{ID:71,Type:"private"};group:=telegram.Chat{ID:-909,Type:"supergroup"};wrong:=telegram.Chat{ID:-910,Type:"supergroup"}
 rows:=func() []model.KeywordSearch {ss,e:=a.ListKeywordSearches();if e!=nil{t.Fatal(e)};return ss}
 row:=func(id int64) model.KeywordSearch {s,found,e:=a.KeywordSearch(id);if e!=nil||!found{t.Fatalf("id=%d found=%v err=%v",id,found,e)};return s}
 unchanged:=func(want []model.KeywordSearch){if got:=rows();!reflect.DeepEqual(want,got){t.Fatalf("protected searches changed want=%+v got=%+v",want,got)}}
 protectRental:=func(){
  got,found,e:=db.LoadSettingsJSON();if e!=nil||!found||got!=raw{t.Fatalf("rental settings changed %q %v",got,e)}
  count,e:=db.CountSeen();if e!=nil||count!=2{t.Fatalf("rental seen count %d %v",count,e)}
  seen,e:=db.SeenIDs([]int64{77001,77002,77003});if e!=nil||!seen[77001]||!seen[77002]||seen[77003]{t.Fatalf("rental seen IDs %v %v",seen,e)}
  for k,v:=range state{got,ok,e:=db.GetState(k);if e!=nil||!ok||got!=v{t.Fatalf("state %s=%q want=%q %v",k,got,v,e)}}
 }
 callback:=func(user int64,chat telegram.Chat,data string) error {
  start:=len(requests)
  e:=telegram.ProcessKeywordUpdateForTest(bot,ctx,telegram.Update{CallbackQuery:&telegram.CallbackQuery{ID:fmt.Sprintf("v-%d",start),From:telegram.User{ID:user},Message:&telegram.Message{MessageID:53,Chat:chat},Data:data}})
  got:=requests[start:];if strings.HasPrefix(data,"ks:")&&(len(got)==0||got[0].Method!="answerCallbackQuery"){t.Fatalf("ack order %s %+v",data,got)}
  if strings.HasPrefix(data,"ks:"){for _,r:=range got{if r.Method!="answerCallbackQuery"&&r.Method!="sendMessage"{t.Fatalf("nonappend-only %s %s",data,r.Method)}}}
  return e
 }
 click:=func(user int64,chat telegram.Chat,data string){if e:=callback(user,chat,data);e!=nil{t.Fatal(data,e)}}
 text:=func(user int64,chat telegram.Chat,value string){if e:=telegram.ProcessKeywordUpdateForTest(bot,ctx,telegram.Update{Message:&telegram.Message{From:&telegram.User{ID:user},Chat:chat,Text:value}});e!=nil{t.Fatal(e)}}
 summary:=func(contains ...string){got:=requests[len(requests)-1];if got.Method!="sendMessage"{t.Fatal(got)};for _,s:=range contains{if !strings.Contains(got.Form.Get("text"),s){t.Fatalf("summary missing %q: %s",s,got.Form.Get("text"))}}}
 trace:=func() []verifierRequest{start:=len(requests);if e:=bot.SendAd(ctx,model.Ad{Card:model.Card{ID:77003,URL:"https://somon.tj/adv/77003_x/",Title:"Квартира"}});e!=nil{t.Fatal(e)};return append([]verifierRequest(nil),requests[start:]...)}
 rentalTrace:=trace()
 // Rejected callbacks must acknowledge only; rejected text cannot consume another user's action.
 for _,c:=range []struct{u int64;chat telegram.Chat}{{73,private},{71,wrong},{73,group}}{
  start:=len(requests);click(c.u,c.chat,"ks:new");text(c.u,c.chat,"intruder")
  if len(requests)!=start+1||requests[start].Method!="answerCallbackQuery"{t.Fatal("authorization sent fresh output")}
  if len(rows())!=0{t.Fatal("unauthorized created search")}
 }
 text(71,private,"/searches");summary("Мои поиски")
 click(71,private,"ks:new");text(71,private,"\t ");if len(rows())!=0{t.Fatal("empty phrase saved")}
 text(71,private,"  полка <A&B>  ");s1:=rows()[0]
 if s1.ID<=0||s1.Phrase!="полка <A&B>"||s1.CategoryKey!="all"||s1.CityKey!="country"||s1.Enabled||s1.Revision!=1{t.Fatal(s1)}
 summary("полка &lt;A&amp;B&gt;","Все категории","Вся страна","на паузе")
 click(72,group,"ks:new");text(72,group,"сборка мебели");s2:=rows()[1]
 if s2.ID<=s1.ID||s2.Enabled{t.Fatal(s2)}
 click(71,private,fmt.Sprintf("ks:city:%d",s1.ID))
 for _,label:=range []string{"Вся страна","Душанбе","Восе","Дангара"}{if !strings.Contains(requests[len(requests)-1].Form.Get("reply_markup"),label){t.Fatal("city catalog label",label)}}
 click(71,private,fmt.Sprintf("ks:setcity:%d:dushanbe",s1.ID))
 click(71,private,fmt.Sprintf("ks:category:%d",s1.ID))
 for _,label:=range []string{"Все категории","Мебель","Столы и стулья","Услуги"}{if !strings.Contains(requests[len(requests)-1].Form.Get("reply_markup"),label){t.Fatal("category catalog label",label)}}
 click(71,private,fmt.Sprintf("ks:setcategory:%d:furniture",s1.ID))
 if row(s1.ID).CityKey!="dushanbe"{t.Fatal("category erased city")}
 click(71,private,fmt.Sprintf("ks:price:%d",s1.ID));text(71,private,"0-0")
 s1=row(s1.ID);if s1.PriceMin==nil||s1.PriceMax==nil||*s1.PriceMin!=0||*s1.PriceMax!=0{t.Fatal("zero bounds lost",s1)}
 click(71,private,fmt.Sprintf("ks:enable:%d",s1.ID));if !row(s1.ID).Enabled{t.Fatal("not enabled")}
 click(72,group,fmt.Sprintf("ks:setcity:%d:dangara",s2.ID));click(72,group,fmt.Sprintf("ks:setcategory:%d:services",s2.ID))
 click(72,group,fmt.Sprintf("ks:price:%d",s2.ID));text(72,group,"-900");click(72,group,fmt.Sprintf("ks:enable:%d",s2.ID))
 s2=row(s2.ID);if s2.PriceMin!=nil||s2.PriceMax==nil||*s2.PriceMax!=900||!s2.Enabled{t.Fatal(s2)}
 // Invalid addressed values cannot write even revision, enabled, or the other row.
 for _,data:=range []string{fmt.Sprintf("ks:setcity:%d:no-city",s1.ID),fmt.Sprintf("ks:setcategory:%d:/x/",s1.ID)}{before:=rows();if callback(71,private,data)==nil{t.Fatal("invalid catalog accepted")};unchanged(before)}
 for _,bad:=range []string{"-1-7","25-24","1.5-2","NaN","999999999999999999999999999999"}{before:=rows();click(71,private,fmt.Sprintf("ks:price:%d",s1.ID));text(71,private,bad);unchanged(before)}
 before:=rows();click(71,private,fmt.Sprintf("ks:phrase:%d",s1.ID));text(71,private," ");unchanged(before)
 // Same group, two different admins: each input addresses its own selected ID.
 click(71,group,fmt.Sprintf("ks:phrase:%d",s1.ID));click(72,group,fmt.Sprintf("ks:price:%d",s2.ID))
 text(73,group,"foreign input");unchanged(before)
 text(72,group,"300-");text(71,group,"настенная полка")
 if row(s1.ID).Phrase!="настенная полка"||row(s2.ID).PriceMin==nil||*row(s2.ID).PriceMin!=300||row(s2.ID).PriceMax!=nil{t.Fatal("cross-admin input lost or crossed")}
 // Same admin across chats: group input cannot consume the private edit.
 click(71,private,fmt.Sprintf("ks:phrase:%d",s1.ID));click(71,group,fmt.Sprintf("ks:phrase:%d",s2.ID))
 text(71,group,"монтаж");text(71,private,"книжная полка")
 if row(s1.ID).Phrase!="книжная полка"||row(s2.ID).Phrase!="монтаж"{t.Fatal("cross-chat input crossed")}
 // Address change replaces the selected search, never applies previous pending input to it.
 before1:=row(s1.ID);click(71,private,fmt.Sprintf("ks:phrase:%d",s1.ID));click(71,private,fmt.Sprintf("ks:price:%d",s2.ID));text(71,private,"80-160")
 if !reflect.DeepEqual(before1,row(s1.ID))||*row(s2.ID).PriceMin!=80||*row(s2.ID).PriceMax!=160{t.Fatal("selection crossed IDs")}
 // An intervening other-admin update invalidates earlier pending revision.
 click(71,private,fmt.Sprintf("ks:phrase:%d",s1.ID));click(72,group,fmt.Sprintf("ks:setcity:%d:vose",s1.ID));before=rows();text(71,private,"stale-write");unchanged(before)
 // Leaving keyword input, and replacing rental input, preserve rental values.
 click(71,private,fmt.Sprintf("ks:price:%d",s1.ID));click(71,private,"ks:rental");before=rows();text(71,private,"100-200");unchanged(before);protectRental()
 click(71,private,"i:price");click(71,private,"ks:new");text(71,private,"/cancel");before=rows();text(71,private,"100-200");unchanged(before);protectRental()
 for _,action:=range []string{"view","enable","phrase","setcity:country"}{parts:=strings.Split(action,":");data:=fmt.Sprintf("ks:%s:987654",parts[0]);if len(parts)==2{data+=":"+parts[1]};before=rows();click(71,private,data);unchanged(before);summary("Мои поиски")}
 // Enable remains independent of settings writes and repeat-enable is idempotent.
 old:=row(s1.ID);click(71,private,fmt.Sprintf("ks:enable:%d",s1.ID));if !reflect.DeepEqual(old,row(s1.ID)){t.Fatal("repeated enable mutated revision")}
 click(71,private,fmt.Sprintf("ks:disable:%d",s1.ID));if row(s1.ID).Enabled||row(s1.ID).Revision!=old.Revision+1||!row(s2.ID).Enabled{t.Fatal("disable crossed IDs")};protectRental()
 want:=rows();if !reflect.DeepEqual(rentalTrace,trace()){t.Fatal("rental delivery changed")}
 if err=db.Close();err!=nil{t.Fatal(err)};db,err=store.Open(path);if err!=nil{t.Fatal(err)}
 a,err=app.New(cfg,db,nil);if err!=nil{t.Fatal(err)};bot=telegram.NewBot(telegram.NewClient(server.URL,"LOCAL"),a,[]int64{71,72},-909,nil)
 unchanged(want);protectRental();if !reflect.DeepEqual(rentalTrace,trace()){t.Fatal("reopened rental trace changed")}
 text(71,private,"/searches");summary("Мои поиски")
 data,_:=json.Marshal(want);t.Logf("Independent AC001/AC007: durable searches=%s; private/targetgroup auth, all pending dimensions, invalid unchanged rows, ack-first zero edits, exact public rental/state+payload after reopen PASS",data)
}
