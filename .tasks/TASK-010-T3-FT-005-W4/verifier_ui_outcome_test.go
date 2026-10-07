//go:build cgo
package telegram_test
import (
 "encoding/json"
 "fmt"
 "reflect"
 "strings"
 "testing"
 "github.com/nicelight/somon-rent-watcher/internal/app"
 "github.com/nicelight/somon-rent-watcher/internal/model"
 "github.com/nicelight/somon-rent-watcher/internal/store"
 "github.com/nicelight/somon-rent-watcher/internal/telegram"
)
func TestVerifierDeleteMenuOutcome(t *testing.T) {
 for _, route:=range []string{"list-private","settings-target"} { t.Run(route,func(t *testing.T){
 h:=newKeywordHarness(t); chat:=telegram.Chat{ID:1,Type:"private"}; if route=="settings-target" {chat=telegram.Chat{ID:-100,Type:"supergroup"}}
 seed:=func(phrase string) model.KeywordSearch {s,e:=h.a.CreateKeywordSearch(model.KeywordSearch{Phrase:phrase,CategoryKey:"all",CityKey:"vose"});if e!=nil{t.Fatal(e)};for _,id:=range []int64{451,452}{if e=h.db.RecordKeywordAdState(s.ID,id,s.Revision,id==451);e!=nil{t.Fatal(e)}};return s}
 selected:=seed("удаляемый & поиск");other:=seed("остаётся")
 rental:=h.snapshotRental(); before:=h.list(); otherHistory,e:=h.db.KeywordAdStates(other.ID);if e!=nil{t.Fatal(e)}
 for _,bad:=range []struct{user int64;chat telegram.Chat}{{99,chat},{1,telegram.Chat{ID:-901,Type:"supergroup"}}} {
 start:=len(h.requests);h.click(bad.user,bad.chat,fmt.Sprintf("ks:delete:%d",selected.ID));hist,e:=h.db.KeywordAdStates(selected.ID)
 if e!=nil||len(hist)!=2||!reflect.DeepEqual(before,h.list())||!reflect.DeepEqual(rental,h.snapshotRental())||len(h.requests)!=start+1{t.Fatal("denied deletion changed state/output",e)}
 }
 // Another administrator owns pending input when deletion occurs.
 h.click(2,chat,fmt.Sprintf("ks:price:%d",selected.ID))
 menu:="ks:list";if route=="settings-target"{menu=fmt.Sprintf("ks:view:%d",selected.ID)};h.click(1,chat,menu)
 var markup telegram.InlineKeyboardMarkup;if e=json.Unmarshal([]byte(h.requests[len(h.requests)-1].form.Get("reply_markup")),&markup);e!=nil{t.Fatal(e)}
 action:="";for _,row:=range markup.InlineKeyboard{for _,button:=range row{if button.CallbackData==fmt.Sprintf("ks:delete:%d",selected.ID){action=button.CallbackData}}};if action==""{t.Fatal("actual menu lacks selected delete button")}
 start:=len(h.requests);h.click(1,chat,action);trace:=h.requests[start:]
 if len(trace)!=3||trace[0].method!="answerCallbackQuery"||trace[1].form.Get("text")!="Поиск удалён."||!strings.Contains(trace[2].form.Get("text"),"Мои поиски"){t.Fatal("ack/confirmation/fresh-list order",trace)}
 if strings.Contains(trace[2].form.Get("reply_markup"),action){t.Fatal("deleted ID remains in new menu")}
 h.text(2,chat,"100-500")
 for _,suffix:=range []string{"view","enable","disable","phrase","price","delete","setcategory"}{data:=fmt.Sprintf("ks:%s:%d",suffix,selected.ID);if suffix=="setcategory"{data+=":services"};h.click(2,chat,data); if !strings.Contains(h.requests[len(h.requests)-2].form.Get("text"),"больше не существует"){t.Fatal("stale callback missing result",data)}}
 stale:=selected;stale.Phrase="восстановить";if _,found,e:=h.a.UpdateKeywordSearch(stale);e!=nil||found{t.Fatal("stale update",found,e)}
 for _,flag:=range []bool{false,true}{if e=h.db.RecordKeywordAdState(selected.ID,999,selected.Revision,flag);e!=nil{t.Fatal(e)}}
 if e=h.db.Close();e!=nil{t.Fatal(e)};h.db,e=store.Open(h.path);if e!=nil{t.Fatal(e)};h.a,e=app.New(h.cfg,h.db,nil);if e!=nil{t.Fatal(e)}
 if _,found,e:=h.a.KeywordSearch(selected.ID);e!=nil||found{t.Fatal("reopened deleted search",found,e)};hist,e:=h.db.KeywordAdStates(selected.ID);if e!=nil||len(hist)!=0{t.Fatal("history resurrection",hist,e)}
 if !reflect.DeepEqual([]model.KeywordSearch{other},h.list())||!reflect.DeepEqual(rental,h.snapshotRental()){t.Fatal("protected search/rental changed")};got,e:=h.db.KeywordAdStates(other.ID);if e!=nil||!reflect.DeepEqual(got,otherHistory){t.Fatal("other history changed")}
 next:=seed("новый ID");if next.ID<=other.ID||next.ID==selected.ID{t.Fatal("ID reused",next.ID)}
 for _,r:=range h.requests {if r.method!="answerCallbackQuery"&&r.method!="sendMessage"{t.Fatal("non append-only output",r.method)}}
 t.Logf("AC002 route=%s deleted=%d other=%d next=%d; real button; deny snapshots; ack→confirm→fresh-list; pending+7 stale actions missing; reopen/history absent; preserved rental+other",route,selected.ID,other.ID,next.ID)
 }) }
}
