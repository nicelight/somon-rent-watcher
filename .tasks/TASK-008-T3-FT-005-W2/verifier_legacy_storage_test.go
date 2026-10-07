//go:build cgo
package store
import (
 "path/filepath"
 "testing"
 "time"
 "github.com/nicelight/somon-rent-watcher/internal/model"
)
func TestVerifierLegacyStorageExactAndStableIdentity(t *testing.T){
 path:=filepath.Join(t.TempDir(),"search.db")
 db,e:=Open(path);if e!=nil{t.Fatal(e)};defer func(){db.Close()}()
 if e=db.SaveSettingsJSON(`{"enabled":true,"price_max":6500,"fixture_unknown":"preserve"}`);e!=nil{t.Fatal(e)}
 if e=db.MarkSeen([]int64{77001,77002},time.Date(2023,1,2,3,4,5,6,time.UTC));e!=nil{t.Fatal(e)}
 if e=db.SetStates(map[string]string{"initialized":"1","telegram_offset":"719","unrelated_key":"verifier-keep"});e!=nil{t.Fatal(e)}
 db.mu.Lock();e=db.execLocked(`DROP TABLE search_monitors; CREATE TABLE verifier_settings AS SELECT * FROM settings; CREATE TABLE verifier_seen AS SELECT * FROM seen_ads; CREATE TABLE verifier_state AS SELECT * FROM state;`);db.mu.Unlock();if e!=nil{t.Fatal(e)}
 if e=db.Close();e!=nil{t.Fatal(e)};db,e=Open(path);if e!=nil{t.Fatal(e)}
 if ss,e:=db.ListKeywordSearches();e!=nil||len(ss)!=0{t.Fatal("legacy init",ss,e)}
 min,max:=0,250
 one,e:=db.CreateKeywordSearch(model.KeywordSearch{Phrase:"A",CategoryKey:"all",CityKey:"country",PriceMin:&min,Revision:1});if e!=nil{t.Fatal(e)}
 two,e:=db.CreateKeywordSearch(model.KeywordSearch{Phrase:"B",CategoryKey:"services",CityKey:"vose",PriceMax:&max,Revision:1});if e!=nil{t.Fatal(e)}
 old:=one;one.Phrase="A-edited";one.Enabled=true
 one,ok,e:=db.UpdateKeywordSearch(one);if e!=nil||!ok||one.Revision!=2{t.Fatal(one,ok,e)}
 if _,ok,e=db.UpdateKeywordSearch(old);e!=nil||ok{t.Fatal("stale overwrote",ok,e)}
 missing:=old;missing.ID=777777;if _,ok,e=db.UpdateKeywordSearch(missing);e!=nil||ok{t.Fatal("missing resurrected",ok,e)}
 // Fixture-only removal proves AUTOINCREMENT survives largest-ID removal/reopen;
 // deletion UX/history is a later task and is not claimed here.
 db.mu.Lock();e=db.execLocked(`DELETE FROM search_monitors WHERE id=2`);db.mu.Unlock();if e!=nil||two.ID!=2{t.Fatal(two,e)}
 if e=db.Close();e!=nil{t.Fatal(e)};db,e=Open(path);if e!=nil{t.Fatal(e)}
 three,e:=db.CreateKeywordSearch(model.KeywordSearch{Phrase:"C",CategoryKey:"all",CityKey:"dangara",Revision:1});if e!=nil||three.ID<=two.ID{t.Fatal("identity reused",three,e)}
 db.mu.Lock();e=db.execLocked(`CREATE TEMP TABLE verifier_check(x INTEGER CHECK(x=1)); INSERT INTO verifier_check SELECT CASE WHEN NOT EXISTS(SELECT * FROM settings EXCEPT SELECT * FROM verifier_settings) AND NOT EXISTS(SELECT * FROM verifier_settings EXCEPT SELECT * FROM settings) AND NOT EXISTS(SELECT * FROM seen_ads EXCEPT SELECT * FROM verifier_seen) AND NOT EXISTS(SELECT * FROM verifier_seen EXCEPT SELECT * FROM seen_ads) AND NOT EXISTS(SELECT * FROM state EXCEPT SELECT * FROM verifier_state) AND NOT EXISTS(SELECT * FROM verifier_state EXCEPT SELECT * FROM state) THEN 1 ELSE 0 END;`);db.mu.Unlock();if e!=nil{t.Fatal("rental exact columns changed",e)}
 got,found,e:=db.KeywordSearch(one.ID);if e!=nil||!found||got.Phrase!=one.Phrase||!got.Enabled||got.Revision!=2||got.PriceMin==nil||*got.PriceMin!=0{t.Fatal("reopen mutation lost",got,found,e)}
 t.Log("Independent AC007: legacy initialization+reopen preserves every rental SQL column incl first_seen_at/unknown state; stale/missing writes unchanged; never-reused ID after high-water removal PASS")
}
