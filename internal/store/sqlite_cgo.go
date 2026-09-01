//go:build cgo

package store

/*
#cgo LDFLAGS: -lsqlite3
#include <sqlite3.h>
#include <stdlib.h>

static int sw_bind_text(sqlite3_stmt *stmt, int idx, const char *value) {
    return sqlite3_bind_text(stmt, idx, value, -1, SQLITE_TRANSIENT);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type DB struct {
	mu     sync.Mutex
	handle *C.sqlite3
	path   string
	closed bool
}

func Open(path string) (*DB, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New("empty SQLite path")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create SQLite directory: %w", err)
	}

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var handle *C.sqlite3
	flags := C.int(C.SQLITE_OPEN_READWRITE | C.SQLITE_OPEN_CREATE | C.SQLITE_OPEN_FULLMUTEX)
	if rc := C.sqlite3_open_v2(cpath, &handle, flags, nil); rc != C.SQLITE_OK {
		msg := "unknown SQLite error"
		if handle != nil {
			msg = C.GoString(C.sqlite3_errmsg(handle))
			C.sqlite3_close_v2(handle)
		}
		return nil, fmt.Errorf("open SQLite %s: %s", path, msg)
	}
	db := &DB{handle: handle, path: path}
	// The database contains bot settings and browsing state. Keep it private even
	// when the caller did not set a restrictive umask.
	if err := os.Chmod(path, 0o600); err != nil {
		C.sqlite3_close_v2(handle)
		return nil, fmt.Errorf("chmod SQLite %s: %w", path, err)
	}
	C.sqlite3_extended_result_codes(handle, 1)
	C.sqlite3_busy_timeout(handle, 5000)
	if err := db.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) initialize() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	const schema = `
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA foreign_keys=ON;
PRAGMA busy_timeout=5000;
PRAGMA wal_autocheckpoint=1000;

CREATE TABLE IF NOT EXISTS seen_ads (
    ad_id INTEGER PRIMARY KEY,
    first_seen_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    json TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS state (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
`
	if err := db.execLocked(schema); err != nil {
		return fmt.Errorf("initialize SQLite: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed || db.handle == nil {
		return nil
	}
	if rc := C.sqlite3_close_v2(db.handle); rc != C.SQLITE_OK {
		return db.sqliteErr(rc)
	}
	db.closed = true
	db.handle = nil
	return nil
}

func (db *DB) Path() string { return db.path }

func (db *DB) SeenIDs(ids []int64) (map[int64]bool, error) {
	result := make(map[int64]bool, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	stmt, err := db.prepareLocked(`SELECT ad_id FROM seen_ads WHERE ad_id IN (` + strings.Join(placeholders, ",") + `)`)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	for i, id := range ids {
		if rc := C.sqlite3_bind_int64(stmt, C.int(i+1), C.sqlite3_int64(id)); rc != C.SQLITE_OK {
			return nil, db.sqliteErr(rc)
		}
	}
	for {
		rc := C.sqlite3_step(stmt)
		switch rc {
		case C.SQLITE_ROW:
			result[int64(C.sqlite3_column_int64(stmt, 0))] = true
		case C.SQLITE_DONE:
			return result, nil
		default:
			return nil, db.sqliteErr(rc)
		}
	}
}

func (db *DB) MarkSeen(ids []int64, when time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.execLocked("BEGIN IMMEDIATE"); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = db.execLocked("ROLLBACK")
		}
	}()

	stmt, err := db.prepareLocked(`INSERT OR IGNORE INTO seen_ads(ad_id, first_seen_at) VALUES(?, ?)`)
	if err != nil {
		return err
	}
	defer C.sqlite3_finalize(stmt)
	stamp := when.UTC().Format(time.RFC3339Nano)
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		C.sqlite3_reset(stmt)
		C.sqlite3_clear_bindings(stmt)
		if rc := C.sqlite3_bind_int64(stmt, 1, C.sqlite3_int64(id)); rc != C.SQLITE_OK {
			return db.sqliteErr(rc)
		}
		if err := bindText(stmt, 2, stamp); err != nil {
			return err
		}
		if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
			return db.sqliteErr(rc)
		}
	}
	if err := db.execLocked("COMMIT"); err != nil {
		return err
	}
	committed = true
	return nil
}

func (db *DB) CountSeen() (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`SELECT COUNT(*) FROM seen_ads`)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(stmt)
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_ROW {
		return 0, db.sqliteErr(rc)
	}
	return int64(C.sqlite3_column_int64(stmt, 0)), nil
}

func (db *DB) LoadSettingsJSON() (string, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`SELECT json FROM settings WHERE id = 1`)
	if err != nil {
		return "", false, err
	}
	defer C.sqlite3_finalize(stmt)
	rc := C.sqlite3_step(stmt)
	if rc == C.SQLITE_DONE {
		return "", false, nil
	}
	if rc != C.SQLITE_ROW {
		return "", false, db.sqliteErr(rc)
	}
	return columnText(stmt, 0), true, nil
}

func (db *DB) SaveSettingsJSON(value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`INSERT INTO settings(id, json) VALUES(1, ?) ON CONFLICT(id) DO UPDATE SET json=excluded.json`)
	if err != nil {
		return err
	}
	defer C.sqlite3_finalize(stmt)
	if err := bindText(stmt, 1, value); err != nil {
		return err
	}
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return db.sqliteErr(rc)
	}
	return nil
}

func (db *DB) GetState(key string) (string, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.prepareLocked(`SELECT value FROM state WHERE key = ?`)
	if err != nil {
		return "", false, err
	}
	defer C.sqlite3_finalize(stmt)
	if err := bindText(stmt, 1, key); err != nil {
		return "", false, err
	}
	rc := C.sqlite3_step(stmt)
	if rc == C.SQLITE_DONE {
		return "", false, nil
	}
	if rc != C.SQLITE_ROW {
		return "", false, db.sqliteErr(rc)
	}
	return columnText(stmt, 0), true, nil
}

func (db *DB) SetState(key, value string) error {
	return db.SetStates(map[string]string{key: value})
}

func (db *DB) SetStates(values map[string]string) error {
	if len(values) == 0 {
		return nil
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.execLocked("BEGIN IMMEDIATE"); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = db.execLocked("ROLLBACK")
		}
	}()
	stmt, err := db.prepareLocked(`INSERT INTO state(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`)
	if err != nil {
		return err
	}
	defer C.sqlite3_finalize(stmt)
	for key, value := range values {
		C.sqlite3_reset(stmt)
		C.sqlite3_clear_bindings(stmt)
		if err := bindText(stmt, 1, key); err != nil {
			return err
		}
		if err := bindText(stmt, 2, value); err != nil {
			return err
		}
		if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
			return db.sqliteErr(rc)
		}
	}
	if err := db.execLocked("COMMIT"); err != nil {
		return err
	}
	committed = true
	return nil
}

func (db *DB) GetStateInt64(key string) (int64, bool, error) {
	value, ok, err := db.GetState(key)
	if err != nil || !ok {
		return 0, ok, err
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("state %s is not int64: %w", key, err)
	}
	return n, true, nil
}

func (db *DB) execLocked(sql string) error {
	if db.closed || db.handle == nil {
		return errors.New("SQLite database is closed")
	}
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	var errmsg *C.char
	rc := C.sqlite3_exec(db.handle, csql, nil, nil, &errmsg)
	if rc != C.SQLITE_OK {
		msg := ""
		if errmsg != nil {
			msg = C.GoString(errmsg)
			C.sqlite3_free(unsafe.Pointer(errmsg))
		}
		if msg == "" {
			return db.sqliteErr(rc)
		}
		return fmt.Errorf("sqlite exec: %s", msg)
	}
	return nil
}

func (db *DB) prepareLocked(sql string) (*C.sqlite3_stmt, error) {
	if db.closed || db.handle == nil {
		return nil, errors.New("SQLite database is closed")
	}
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	var stmt *C.sqlite3_stmt
	if rc := C.sqlite3_prepare_v2(db.handle, csql, -1, &stmt, nil); rc != C.SQLITE_OK {
		return nil, db.sqliteErr(rc)
	}
	return stmt, nil
}

func (db *DB) sqliteErr(rc C.int) error {
	msg := "SQLite error"
	if db.handle != nil {
		msg = C.GoString(C.sqlite3_errmsg(db.handle))
	}
	return fmt.Errorf("sqlite code %d: %s", int(rc), msg)
}

func bindText(stmt *C.sqlite3_stmt, index int, value string) error {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	if rc := C.sw_bind_text(stmt, C.int(index), cvalue); rc != C.SQLITE_OK {
		return fmt.Errorf("sqlite bind text code %d", int(rc))
	}
	return nil
}

func columnText(stmt *C.sqlite3_stmt, index int) string {
	value := C.sqlite3_column_text(stmt, C.int(index))
	if value == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(value)))
}
