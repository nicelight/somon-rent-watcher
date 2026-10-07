//go:build cgo

package telegram_test

import (
 "fmt"
 "reflect"
 "strings"
 "testing"

 "github.com/nicelight/somon-rent-watcher/internal/app"
 "github.com/nicelight/somon-rent-watcher/internal/model"
 "github.com/nicelight/somon-rent-watcher/internal/store"
 "github.com/nicelight/somon-rent-watcher/internal/telegram"
)

func seedDeleteSearch(t *testing.T,h *keywordHarness,phrase string) model.KeywordSearch {
 t.Helper()
 s,e:=h.a.CreateKeywordSearch(model.KeywordSearch{Phrase:phrase,CategoryKey:"all",CityKey:"country"});if e!=nil{t.Fatal(e)}
 if e=h.db.RecordKeywordAdState(s.ID,s.ID*10,s.Revision,true);e!=nil{t.Fatal(e)}
 if e=h.db.RecordKeywordAdState(s.ID,s.ID*10+1,s.Revision,false);e!=nil{t.Fatal(e)}
 return s
}
func TestKeywordDeleteBothMenus(t *testing.T) {
 for _,route:=range []string{"list","settings"}{t.Run(route,func(t *testing.T){
  h:=newKeywordHarness(t);chat:=telegram.Chat{ID:1,Type:"private"}
  other:=seedDeleteSearch(t,h,"сохранить");selected:=seedDeleteSearch(t,h,"удалить")
  rental:=h.snapshotRental();otherHistory,e:=h.db.KeywordAdStates(other.ID);if e!=nil{t.Fatal(e)}
  if route=="list"{h.click(1,chat,"ks:list")}else{h.click(1,chat,fmt.Sprintf("ks:view:%d",selected.ID))}
  callback:=fmt.Sprintf("ks:delete:%d",selected.ID)
  buttons:=h.requests[len(h.requests)-1].form.Get("reply_markup")
  // Always exercise the actual action before asserting menu availability.
  routeErr:=h.callback(1,chat,callback)
  _,found,e:=h.a.KeywordSearch(selected.ID);if e!=nil{t.Fatal(e)}
  history,e:=h.db.KeywordAdStates(selected.ID);if e!=nil{t.Fatal(e)}
  t.Logf("FT-005-AC-002 %s delete error=%v selected_exists=%v history_rows=%d",route,routeErr,found,len(history))
  if found||len(history)!=0{t.Fatalf("delete left selected monitor/history: exists=%v rows=%d",found,len(history))}
  if routeErr!=nil||!strings.Contains(buttons,callback){t.Fatal("delete route unavailable",buttons,routeErr)}
  if !reflect.DeepEqual(other,h.search(other.ID)){t.Fatal("other monitor changed")}
  got,e:=h.db.KeywordAdStates(other.ID);if e!=nil||!reflect.DeepEqual(otherHistory,got){t.Fatal("other history changed",got,e)}
  if !reflect.DeepEqual(rental,h.snapshotRental()){t.Fatal("rental changed")}
  // Other admin has pending input when the selected search disappears.
  for _,action:=range []string{"view","enable","phrase","delete","setcity"}{data:=fmt.Sprintf("ks:%s:%d",action,selected.ID);if action=="setcity"{data+=":dushanbe"};h.click(2,chat,data)}
  if len(h.list())!=1{t.Fatal("old callback resurrected search")}
  if e=h.db.Close();e!=nil{t.Fatal(e)};h.db,e=store.Open(h.path);if e!=nil{t.Fatal(e)}
  h.a,e=app.New(h.cfg,h.db,nil);if e!=nil{t.Fatal(e)}
  if _,found,e=h.a.KeywordSearch(selected.ID);e!=nil||found{t.Fatal("deletion not durable",e)}
  history,e=h.db.KeywordAdStates(selected.ID);if e!=nil||len(history)!=0{t.Fatal("history restored",history,e)}
  if !reflect.DeepEqual(rental,h.snapshotRental()){t.Fatal("reopen changed rental")}
  next:=seedDeleteSearch(t,h,"новый");if next.ID<=selected.ID{t.Fatal("deleted ID reused",next.ID,selected.ID)}
  t.Logf("FT-005-AC-002 GREEN %s: other=%d unchanged, deleted=%d absent after reopen, next=%d; fresh list and ack-first append-only trace",route,other.ID,selected.ID,next.ID)
 })}
}
func TestKeywordDeleteAuthorizationPreserved(t *testing.T){
 h:=newKeywordHarness(t);s:=seedDeleteSearch(t,h,"сохранить");before:=h.list();history,e:=h.db.KeywordAdStates(s.ID);if e!=nil{t.Fatal(e)};rental:=h.snapshotRental()
 for _,ctx:=range []struct{user int64;chat telegram.Chat}{{3,telegram.Chat{ID:1,Type:"private"}},{1,telegram.Chat{ID:-101,Type:"supergroup"}},{3,telegram.Chat{ID:-100,Type:"supergroup"}}}{
  start:=len(h.requests);h.click(ctx.user,ctx.chat,fmt.Sprintf("ks:delete:%d",s.ID));if len(h.requests)!=start+1{t.Fatal("unauthorized output")}
  got,e:=h.db.KeywordAdStates(s.ID);if e!=nil||!reflect.DeepEqual(history,got)||!reflect.DeepEqual(before,h.list())||!reflect.DeepEqual(rental,h.snapshotRental()){t.Fatal("unauthorized delete altered protected rows")}
 }
 t.Log("FT-005-AC-002 initial/final GREEN: unauthorized/wrong-chat delete preserves full observed snapshots")
}
