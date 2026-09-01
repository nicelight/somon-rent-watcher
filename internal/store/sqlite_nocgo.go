//go:build !cgo

package store

import (
	"errors"
	"time"
)

var errCGODisabled = errors.New("SQLite support requires CGO_ENABLED=1 and sqlite3 development headers")

type DB struct{}

func Open(string) (*DB, error)                        { return nil, errCGODisabled }
func (*DB) Close() error                              { return nil }
func (*DB) Path() string                              { return "" }
func (*DB) SeenIDs([]int64) (map[int64]bool, error)   { return nil, errCGODisabled }
func (*DB) MarkSeen([]int64, time.Time) error         { return errCGODisabled }
func (*DB) CountSeen() (int64, error)                 { return 0, errCGODisabled }
func (*DB) LoadSettingsJSON() (string, bool, error)   { return "", false, errCGODisabled }
func (*DB) SaveSettingsJSON(string) error             { return errCGODisabled }
func (*DB) GetState(string) (string, bool, error)     { return "", false, errCGODisabled }
func (*DB) SetState(string, string) error             { return errCGODisabled }
func (*DB) SetStates(map[string]string) error         { return errCGODisabled }
func (*DB) GetStateInt64(string) (int64, bool, error) { return 0, false, errCGODisabled }
