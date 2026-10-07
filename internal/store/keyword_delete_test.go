//go:build cgo

package store

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func TestKeywordDeleteCascadeRollbackAndPreservation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search.db")
	db, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer func() { db.Close() }()
	if e = db.SaveSettingsJSON(`{"enabled":true,"price_max":6000}`); e != nil {
		t.Fatal(e)
	}
	if e = db.MarkSeen([]int64{7, 9}, time.Date(2026, 10, 6, 1, 2, 3, 0, time.UTC)); e != nil {
		t.Fatal(e)
	}
	if e = db.SetStates(map[string]string{"initialized": "1", "telegram_offset": "88", "arbitrary": "protected"}); e != nil {
		t.Fatal(e)
	}
	other, e := db.CreateKeywordSearch(model.KeywordSearch{Phrase: "other", CategoryKey: "all", CityKey: "country", Revision: 1})
	if e != nil {
		t.Fatal(e)
	}
	selected, e := db.CreateKeywordSearch(model.KeywordSearch{Phrase: "selected", CategoryKey: "all", CityKey: "country", Revision: 1})
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []model.KeywordSearch{other, selected} {
		if e = db.RecordKeywordAdState(s.ID, 10, s.Revision, true); e != nil {
			t.Fatal(e)
		}
		if e = db.RecordKeywordAdState(s.ID, 20, s.Revision, false); e != nil {
			t.Fatal(e)
		}
	}
	before, e := db.ListKeywordSearches()
	if e != nil {
		t.Fatal(e)
	}
	selectedHistory, e := db.KeywordAdStates(selected.ID)
	if e != nil {
		t.Fatal(e)
	}
	otherHistory, e := db.KeywordAdStates(other.ID)
	if e != nil {
		t.Fatal(e)
	}
	db.mu.Lock()
	e = db.execLocked(`CREATE TABLE fixture_delete_settings AS SELECT * FROM settings;
 CREATE TABLE fixture_delete_seen AS SELECT * FROM seen_ads;
 CREATE TABLE fixture_delete_state AS SELECT * FROM state;
 CREATE TRIGGER fail_keyword_delete AFTER DELETE ON search_ad_state WHEN OLD.ad_id=20
 BEGIN SELECT RAISE(ABORT,'test cascade failure'); END;`)
	db.mu.Unlock()
	if e != nil {
		t.Fatal(e)
	}
	// The second child deletion fails after SQLite has started cascade work; the
	// entire statement, including parent and preceding child, must roll back.
	found, e := db.DeleteKeywordSearch(selected.ID)
	if e == nil || found || !strings.Contains(e.Error(), "test cascade failure") {
		t.Fatal("failure not injected", found, e)
	}
	got, e := db.ListKeywordSearches()
	if e != nil || !reflect.DeepEqual(before, got) {
		t.Fatal("failed delete changed monitors", got, e)
	}
	gotHistory, e := db.KeywordAdStates(selected.ID)
	if e != nil || !reflect.DeepEqual(selectedHistory, gotHistory) {
		t.Fatal("failed cascade partially deleted history", gotHistory, e)
	}
	db.mu.Lock()
	e = db.execLocked(`DROP TRIGGER fail_keyword_delete;`)
	db.mu.Unlock()
	if e != nil {
		t.Fatal(e)
	}
	found, e = db.DeleteKeywordSearch(selected.ID)
	if e != nil || !found {
		t.Fatal(found, e)
	}
	if found, e = db.DeleteKeywordSearch(selected.ID); e != nil || found {
		t.Fatal("stale delete", found, e)
	}
	if e = db.Close(); e != nil {
		t.Fatal(e)
	}
	db, e = Open(path)
	if e != nil {
		t.Fatal(e)
	}
	if got, found, e := db.KeywordSearch(other.ID); e != nil || !found || !reflect.DeepEqual(other, got) {
		t.Fatal("other search changed", got, e)
	}
	if got, found, e := db.KeywordSearch(selected.ID); e != nil || found {
		t.Fatal("deleted search reappeared", got, e)
	}
	if got, e := db.KeywordAdStates(selected.ID); e != nil || len(got) != 0 {
		t.Fatal("deleted history reappeared", got, e)
	}
	if got, e := db.KeywordAdStates(other.ID); e != nil || !reflect.DeepEqual(otherHistory, got) {
		t.Fatal("other history changed", got, e)
	}
	db.mu.Lock()
	e = db.execLocked(`CREATE TEMP TABLE delete_preservation_check(ok INTEGER CHECK(ok=1));
 INSERT INTO delete_preservation_check SELECT CASE WHEN
 NOT EXISTS(SELECT * FROM settings EXCEPT SELECT * FROM fixture_delete_settings) AND
 NOT EXISTS(SELECT * FROM fixture_delete_settings EXCEPT SELECT * FROM settings) AND
 NOT EXISTS(SELECT * FROM seen_ads EXCEPT SELECT * FROM fixture_delete_seen) AND
 NOT EXISTS(SELECT * FROM fixture_delete_seen EXCEPT SELECT * FROM seen_ads) AND
 NOT EXISTS(SELECT * FROM state EXCEPT SELECT * FROM fixture_delete_state) AND
 NOT EXISTS(SELECT * FROM fixture_delete_state EXCEPT SELECT * FROM state)
 THEN 1 ELSE 0 END;`)
	db.mu.Unlock()
	if e != nil {
		t.Fatal("exact rental rows changed", e)
	}
	t.Log("FT-005-AC-002 GREEN: cascade interrupted at second history row rolls back monitor + both histories; successful deletion durable; other search/history and exact bidirectional settings/seen timestamp/all state/offset preserved")
}
