//go:build cgo

package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestDBRoundTrip(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "store.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.SaveSettingsJSON(`{"enabled":false}`); err != nil {
		t.Fatal(err)
	}
	settings, ok, err := db.LoadSettingsJSON()
	if err != nil || !ok || settings != `{"enabled":false}` {
		t.Fatalf("settings=%q ok=%v err=%v", settings, ok, err)
	}

	if err := db.MarkSeen([]int64{1, 2, 2}, time.Now()); err != nil {
		t.Fatal(err)
	}
	seen, err := db.SeenIDs([]int64{2, 3})
	if err != nil || !seen[2] || seen[3] {
		t.Fatalf("seen=%v err=%v", seen, err)
	}
	count, err := db.CountSeen()
	if err != nil || count != 2 {
		t.Fatalf("count=%d err=%v", count, err)
	}

	if err := db.SetStates(map[string]string{"key": "value", "number": "42"}); err != nil {
		t.Fatal(err)
	}
	value, found, err := db.GetState("key")
	if err != nil || !found || value != "value" {
		t.Fatalf("state=%q found=%v err=%v", value, found, err)
	}
	number, found, err := db.GetStateInt64("number")
	if err != nil || !found || number != 42 {
		t.Fatalf("number=%d found=%v err=%v", number, found, err)
	}
}
