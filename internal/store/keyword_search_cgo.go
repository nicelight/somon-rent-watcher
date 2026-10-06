//go:build cgo

package store

/*
#include <sqlite3.h>
*/
import "C"

import (
	"errors"
	"github.com/nicelight/somon-rent-watcher/internal/model"
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
