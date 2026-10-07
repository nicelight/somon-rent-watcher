//go:build cgo

package store

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func TestVerifierCurrentFeedHistory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "independent-history.db")
	db, err := Open(path)
	if err != nil { t.Fatal(err) }
	t.Cleanup(func() { db.Close() })
	a, err := db.CreateKeywordSearch(model.KeywordSearch{Phrase:"first", CategoryKey:"all", CityKey:"country", Revision:7})
	if err != nil { t.Fatal(err) }
	b, err := db.CreateKeywordSearch(model.KeywordSearch{Phrase:"second", CategoryKey:"all", CityKey:"country", Revision:11})
	if err != nil { t.Fatal(err) }
	for _, row := range []struct{ monitor, id, revision int64; delivered bool }{
		{a.ID, 101, 7, true}, {a.ID, 202, 7, false}, {a.ID, 303, 7, true}, {a.ID, 404, 7, false},
		{b.ID, 101, 11, false}, {b.ID, 202, 11, true}, {b.ID, 505, 11, true},
	} {
		if err = db.RecordKeywordAdState(row.monitor, row.id, row.revision, row.delivered); err != nil { t.Fatal(err) }
	}
	beforeA, err := db.KeywordAdStates(a.ID); if err != nil { t.Fatal(err) }
	beforeB, err := db.KeywordAdStates(b.ID); if err != nil { t.Fatal(err) }
	query := func(monitor int64, ids []int64, want map[int64]model.KeywordAdState) {
		t.Helper()
		got, e := db.KeywordAdStatesForIDs(monitor, ids)
		if e != nil || !reflect.DeepEqual(got, want) { t.Fatalf("monitor=%d ids=%v got=%v want=%v error=%v", monitor, ids, got, want, e) }
	}
	query(a.ID, []int64{202,101,202,999}, map[int64]model.KeywordAdState{101:beforeA[101],202:beforeA[202]})
	query(b.ID, []int64{101,202,999}, map[int64]model.KeywordAdState{101:beforeB[101],202:beforeB[202]})
	query(99999, []int64{101}, map[int64]model.KeywordAdState{})
	query(a.ID, nil, map[int64]model.KeywordAdState{})
	query(a.ID, []int64{}, map[int64]model.KeywordAdState{})
	// A nil DB receiver makes any DB access fail: the empty path must be a pure return.
	var noDB *DB
	for _, ids := range [][]int64{nil, {}} {
		got,e := noDB.KeywordAdStatesForIDs(a.ID, ids)
		if e != nil || len(got) != 0 { t.Fatalf("empty requires no DB: %v %v",got,e) }
	}
	a.Phrase="edited"
	updated,found,err := db.UpdateKeywordSearch(a)
	if err != nil || !found || updated.Revision!=8 { t.Fatal(updated,found,err) }
	query(a.ID, []int64{101,202}, map[int64]model.KeywordAdState{101:beforeA[101],202:beforeA[202]})
	if err=db.Close(); err!=nil { t.Fatal(err) }
	db,err=Open(path); if err!=nil { t.Fatal(err) }
	for _, row := range []struct{ id int64; want map[int64]model.KeywordAdState }{{a.ID,beforeA},{b.ID,beforeB}} {
		got,e:=db.KeywordAdStates(row.id)
		if e!=nil || !reflect.DeepEqual(got,row.want) { t.Fatalf("full diagnostic history changed after bounded reads/edit/reopen: got=%v want=%v error=%v",got,row.want,e) }
	}
	query(a.ID, []int64{101,202}, map[int64]model.KeywordAdState{101:beforeA[101],202:beforeA[202]})
	t.Log("FT-005-AC-009 selected 2 of 4 first-monitor rows; second monitor independent; nil/empty use no DB; delivered + evaluated_revision exact across edit/reopen; diagnostic history retains all 7 rows")
}
