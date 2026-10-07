//go:build cgo
package store
import (
 "fmt"
 "path/filepath"
 "reflect"
 "strings"
 "testing"
 "time"
 "github.com/nicelight/somon-rent-watcher/internal/model"
)
func TestVerifierDeleteRollbackExactRows(t *testing.T){
 path:=filepath.Join(t.TempDir(),"search.db");db,e:=Open(path);if e!=nil{t.Fatal(e)};defer func(){db.Close()}()
 if e=db.SaveSettingsJSON(`{"enabled":false,"price_max":4321}`);e!=nil{t.Fatal(e)};if e=db.MarkSeen([]int64{711,722},time.Date(2026,10,7,4,3,2,1,time.UTC));e!=nil{t.Fatal(e)};if e=db.SetStates(map[string]string{"initialized":"1","telegram_offset":"671","custom":"untouched"});e!=nil{t.Fatal(e)}
 selected,e:=db.CreateKeywordSearch(model.KeywordSearch{Phrase:"selected",CategoryKey:"all",CityKey:"country",Revision:7});if e!=nil{t.Fatal(e)};other,e:=db.CreateKeywordSearch(model.KeywordSearch{Phrase:"other",CategoryKey:"all",CityKey:"country",Revision:4});if e!=nil{t.Fatal(e)}
 for _,s:=range []model.KeywordSearch{selected,other}{for _,id:=range []int64{30,31,32}{if e=db.RecordKeywordAdState(s.ID,id,s.Revision,id==30);e!=nil{t.Fatal(e)}}}
 before,e:=db.ListKeywordSearches();if e!=nil{t.Fatal(e)};history,e:=db.KeywordAdStates(selected.ID);if e!=nil{t.Fatal(e)};otherHistory,e:=db.KeywordAdStates(other.ID);if e!=nil{t.Fatal(e)}
 db.mu.Lock();e=db.execLocked(fmt.Sprintf(`CREATE TABLE verifier_settings AS SELECT * FROM settings;CREATE TABLE verifier_seen AS SELECT * FROM seen_ads;CREATE TABLE verifier_state AS SELECT * FROM state;
 CREATE TRIGGER verifier_abort AFTER DELETE ON search_ad_state WHEN OLD.monitor_id=%d AND NOT EXISTS(SELECT 1 FROM search_ad_state WHERE monitor_id=%d) BEGIN SELECT RAISE(ABORT,'verifier-last-child-abort'); END;`,selected.ID,selected.ID));db.mu.Unlock();if e!=nil{t.Fatal(e)}
 if found,e:=db.DeleteKeywordSearch(selected.ID);e==nil||found||!strings.Contains(e.Error(),"verifier-last-child-abort"){t.Fatal("last child failure not observed",found,e)}
 rows,e:=db.ListKeywordSearches();if e!=nil||!reflect.DeepEqual(before,rows){t.Fatal("monitor rollback incomplete")};got,e:=db.KeywordAdStates(selected.ID);if e!=nil||!reflect.DeepEqual(history,got){t.Fatal("3-child rollback incomplete",got,e)}
 db.mu.Lock();e=db.execLocked(`DROP TRIGGER verifier_abort;`);db.mu.Unlock();if e!=nil{t.Fatal(e)}
 if found,e:=db.DeleteKeywordSearch(selected.ID);e!=nil||!found{t.Fatal("successful delete",found,e)};if e=db.Close();e!=nil{t.Fatal(e)};db,e=Open(path);if e!=nil{t.Fatal(e)}
 if got,found,e:=db.KeywordSearch(other.ID);e!=nil||!found||!reflect.DeepEqual(other,got){t.Fatal("other monitor changed")};got,e=db.KeywordAdStates(other.ID);if e!=nil||!reflect.DeepEqual(otherHistory,got){t.Fatal("other history changed")};got,e=db.KeywordAdStates(selected.ID);if e!=nil||len(got)!=0{t.Fatal("deleted history returned")};if _,found,e:=db.KeywordSearch(selected.ID);e!=nil||found{t.Fatal("deleted monitor returned")}
 db.mu.Lock();e=db.execLocked(`CREATE TEMP TABLE verifier_exact(ok INTEGER CHECK(ok=1)); INSERT INTO verifier_exact SELECT CASE WHEN
 NOT EXISTS(SELECT * FROM settings EXCEPT SELECT * FROM verifier_settings) AND NOT EXISTS(SELECT * FROM verifier_settings EXCEPT SELECT * FROM settings) AND
 NOT EXISTS(SELECT * FROM seen_ads EXCEPT SELECT * FROM verifier_seen) AND NOT EXISTS(SELECT * FROM verifier_seen EXCEPT SELECT * FROM seen_ads) AND
 NOT EXISTS(SELECT * FROM state EXCEPT SELECT * FROM verifier_state) AND NOT EXISTS(SELECT * FROM verifier_state EXCEPT SELECT * FROM state) THEN 1 ELSE 0 END;`);db.mu.Unlock();if e!=nil{t.Fatal("exact persistent rental snapshot mismatch",e)}
 t.Log("AC002 final-child cascade ABORT restored parent+all3 children; successful deletion reopened absent; exact SQL bidirectional settings/seen/timestamps/state/offset and other monitor/history preserved")
}
