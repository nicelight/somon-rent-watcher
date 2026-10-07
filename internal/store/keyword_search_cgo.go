//go:build cgo

package store

/*
#include <sqlite3.h>
*/
import "C"

import (
	"errors"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"strings"
)

const keywordSearchColumns = `id, phrase, category_key, city_key, price_min, price_max, enabled, revision`

func keywordSearchRow(stmt *C.sqlite3_stmt) model.KeywordSearch {
	s := model.KeywordSearch{ID: int64(C.sqlite3_column_int64(stmt, 0)), Phrase: columnText(stmt, 1), CategoryKey: columnText(stmt, 2), CityKey: columnText(stmt, 3), Enabled: C.sqlite3_column_int(stmt, 6) != 0, Revision: int64(C.sqlite3_column_int64(stmt, 7))}
	if C.sqlite3_column_type(stmt, 4) != C.SQLITE_NULL {
		v := int(C.sqlite3_column_int64(stmt, 4))
		s.PriceMin = &v
	}
	if C.sqlite3_column_type(stmt, 5) != C.SQLITE_NULL {
		v := int(C.sqlite3_column_int64(stmt, 5))
		s.PriceMax = &v
	}
	return s
}

func (db *DB) ListKeywordSearches() ([]model.KeywordSearch, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`SELECT ` + keywordSearchColumns + ` FROM search_monitors ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	searches := []model.KeywordSearch{}
	for {
		rc := C.sqlite3_step(stmt)
		switch rc {
		case C.SQLITE_ROW:
			searches = append(searches, keywordSearchRow(stmt))
		case C.SQLITE_DONE:
			return searches, nil
		default:
			return nil, db.sqliteErr(rc)
		}
	}
}

func (db *DB) KeywordSearch(id int64) (model.KeywordSearch, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`SELECT ` + keywordSearchColumns + ` FROM search_monitors WHERE id=?`)
	if err != nil {
		return model.KeywordSearch{}, false, err
	}
	defer C.sqlite3_finalize(stmt)
	if rc := C.sqlite3_bind_int64(stmt, 1, C.sqlite3_int64(id)); rc != C.SQLITE_OK {
		return model.KeywordSearch{}, false, db.sqliteErr(rc)
	}
	rc := C.sqlite3_step(stmt)
	if rc == C.SQLITE_DONE {
		return model.KeywordSearch{}, false, nil
	}
	if rc != C.SQLITE_ROW {
		return model.KeywordSearch{}, false, db.sqliteErr(rc)
	}
	return keywordSearchRow(stmt), true, nil
}

func bindKeywordSearch(stmt *C.sqlite3_stmt, s model.KeywordSearch) error {
	for i, value := range []string{s.Phrase, s.CategoryKey, s.CityKey} {
		if err := bindText(stmt, i+1, value); err != nil {
			return err
		}
	}
	for i, value := range []*int{s.PriceMin, s.PriceMax} {
		var rc C.int
		if value == nil {
			rc = C.sqlite3_bind_null(stmt, C.int(i+4))
		} else {
			rc = C.sqlite3_bind_int64(stmt, C.int(i+4), C.sqlite3_int64(*value))
		}
		if rc != C.SQLITE_OK {
			return errors.New("bind keyword price")
		}
	}
	enabled := 0
	if s.Enabled {
		enabled = 1
	}
	if rc := C.sqlite3_bind_int(stmt, 6, C.int(enabled)); rc != C.SQLITE_OK {
		return errors.New("bind keyword enabled")
	}
	if rc := C.sqlite3_bind_int64(stmt, 7, C.sqlite3_int64(s.Revision)); rc != C.SQLITE_OK {
		return errors.New("bind keyword revision")
	}
	return nil
}

func (db *DB) CreateKeywordSearch(s model.KeywordSearch) (model.KeywordSearch, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`INSERT INTO search_monitors(phrase,category_key,city_key,price_min,price_max,enabled,revision) VALUES(?,?,?,?,?,?,?)`)
	if err != nil {
		return model.KeywordSearch{}, err
	}
	defer C.sqlite3_finalize(stmt)
	if err = bindKeywordSearch(stmt, s); err != nil {
		return model.KeywordSearch{}, err
	}
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return model.KeywordSearch{}, db.sqliteErr(rc)
	}
	s.ID = int64(C.sqlite3_last_insert_rowid(db.handle))
	return s, nil
}

// Update is a single atomic statement. A stale revision or missing ID cannot
// overwrite newer values or resurrect a search. SQLite remains the only writer.
func (db *DB) UpdateKeywordSearch(s model.KeywordSearch) (model.KeywordSearch, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	expected := s.Revision
	s.Revision++
	stmt, err := db.prepareLocked(`UPDATE search_monitors SET phrase=?,category_key=?,city_key=?,price_min=?,price_max=?,enabled=?,revision=? WHERE id=? AND revision=?`)
	if err != nil {
		return model.KeywordSearch{}, false, err
	}
	defer C.sqlite3_finalize(stmt)
	if err = bindKeywordSearch(stmt, s); err != nil {
		return model.KeywordSearch{}, false, err
	}
	if rc := C.sqlite3_bind_int64(stmt, 8, C.sqlite3_int64(s.ID)); rc != C.SQLITE_OK {
		return model.KeywordSearch{}, false, db.sqliteErr(rc)
	}
	if rc := C.sqlite3_bind_int64(stmt, 9, C.sqlite3_int64(expected)); rc != C.SQLITE_OK {
		return model.KeywordSearch{}, false, db.sqliteErr(rc)
	}
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return model.KeywordSearch{}, false, db.sqliteErr(rc)
	}
	if C.sqlite3_changes(db.handle) == 0 {
		return model.KeywordSearch{}, false, nil
	}
	return s, true, nil
}

func (db *DB) KeywordAdStates(monitorID int64) (map[int64]model.KeywordAdState, error) {
	return db.keywordAdStates(monitorID, nil)
}

// KeywordAdStatesForIDs bounds polling reads to the current feed. It does not
// prune delivered or rejected history; the full diagnostic reader stays compatible.
func (db *DB) KeywordAdStatesForIDs(monitorID int64, ids []int64) (map[int64]model.KeywordAdState, error) {
	if len(ids) == 0 {
		return map[int64]model.KeywordAdState{}, nil
	}
	return db.keywordAdStates(monitorID, ids)
}

func (db *DB) keywordAdStates(monitorID int64, ids []int64) (map[int64]model.KeywordAdState, error) {

	db.mu.Lock()
	defer db.mu.Unlock()
	query := `SELECT ad_id,evaluated_revision,delivered FROM search_ad_state WHERE monitor_id=?`
	if len(ids) > 0 {
		query += ` AND ad_id IN (` + strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",") + `)`
	}
	stmt, err := db.prepareLocked(query)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	if rc := C.sqlite3_bind_int64(stmt, 1, C.sqlite3_int64(monitorID)); rc != C.SQLITE_OK {
		return nil, db.sqliteErr(rc)
	}
	for i, id := range ids {
		if rc := C.sqlite3_bind_int64(stmt, C.int(i+2), C.sqlite3_int64(id)); rc != C.SQLITE_OK {
			return nil, db.sqliteErr(rc)
		}
	}
	states := make(map[int64]model.KeywordAdState)
	for {
		rc := C.sqlite3_step(stmt)
		switch rc {
		case C.SQLITE_ROW:
			states[int64(C.sqlite3_column_int64(stmt, 0))] = model.KeywordAdState{EvaluatedRevision: int64(C.sqlite3_column_int64(stmt, 1)), Delivered: C.sqlite3_column_int(stmt, 2) != 0}
		case C.SQLITE_DONE:
			return states, nil
		default:
			return nil, db.sqliteErr(rc)
		}
	}
}

// Rejections are conditional on the evaluated revision. Confirmed deliveries
// survive edits, but neither path creates history for a missing monitor.
func (db *DB) RecordKeywordAdState(monitorID, adID, revision int64, delivered bool) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`INSERT INTO search_ad_state(monitor_id,ad_id,evaluated_revision,delivered)
 SELECT id,?, ?, ? FROM search_monitors WHERE id=? AND (?=1 OR revision=?)
 ON CONFLICT(monitor_id,ad_id) DO UPDATE SET evaluated_revision=excluded.evaluated_revision, delivered=MAX(search_ad_state.delivered,excluded.delivered)
 WHERE search_ad_state.delivered=0`)
	if err != nil {
		return err
	}
	defer C.sqlite3_finalize(stmt)
	flag := int64(0)
	if delivered {
		flag = 1
	}
	for i, v := range []int64{adID, revision, flag, monitorID, flag, revision} {
		if rc := C.sqlite3_bind_int64(stmt, C.int(i+1), C.sqlite3_int64(v)); rc != C.SQLITE_OK {
			return db.sqliteErr(rc)
		}
	}
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return db.sqliteErr(rc)
	}
	return nil
}

// SQLite executes this statement and the existing ON DELETE CASCADE within one
// transaction. A failed history deletion rolls back the monitor deletion too.
func (db *DB) DeleteKeywordSearch(id int64) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`DELETE FROM search_monitors WHERE id=?`)
	if err != nil {
		return false, err
	}
	defer C.sqlite3_finalize(stmt)
	if rc := C.sqlite3_bind_int64(stmt, 1, C.sqlite3_int64(id)); rc != C.SQLITE_OK {
		return false, db.sqliteErr(rc)
	}
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return false, db.sqliteErr(rc)
	}
	return C.sqlite3_changes(db.handle) != 0, nil
}
