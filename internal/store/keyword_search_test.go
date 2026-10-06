//go:build cgo

package store

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func TestKeywordAdditiveInitializationPreservesExactRentalRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { db.Close() }()
	if err = db.SaveSettingsJSON(`{"enabled":true,"price_max":6000}`); err != nil {
		t.Fatal(err)
	}
	if err = db.MarkSeen([]int64{101, 202}, time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if err = db.SetStates(map[string]string{"initialized": "1", "telegram_offset": "88", "previous_ordinary_ids": "[101,202]", "last_successful_poll_at": "2026-10-06T12:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	// Build the legacy SQLite fixture from existing adapter methods. These
	// snapshot tables are test-fixture data only, never production schema.
	db.mu.Lock()
	err = db.execLocked(`DROP TABLE search_monitors;
 CREATE TABLE fixture_settings AS SELECT * FROM settings;
 CREATE TABLE fixture_seen AS SELECT * FROM seen_ads;
 CREATE TABLE fixture_state AS SELECT * FROM state;`)
	db.mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	searches, err := db.ListKeywordSearches()
	if err != nil || len(searches) != 0 {
		t.Fatal(searches, err)
	}
	min, max := 0, 200
	s1, err := db.CreateKeywordSearch(model.KeywordSearch{Phrase: "стол", CategoryKey: "all", CityKey: "vose", PriceMin: &min, PriceMax: &max, Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	s2, err := db.CreateKeywordSearch(model.KeywordSearch{Phrase: "ремонт", CategoryKey: "services", CityKey: "country", Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	if s1.ID <= 0 || s2.ID <= s1.ID {
		t.Fatal("stable IDs invalid")
	}
	s1.Phrase = "стулья"
	s1.Enabled = true
	updated, found, err := db.UpdateKeywordSearch(s1)
	if err != nil || !found || updated.Revision != 2 {
		t.Fatal(updated, found, err)
	}
	if _, found, err = db.UpdateKeywordSearch(s1); err != nil || found {
		t.Fatal("stale write changed row", found, err)
	}
	missing := s1
	missing.ID = 999999
	if _, found, err = db.UpdateKeywordSearch(missing); err != nil || found {
		t.Fatal("missing row recreated")
	}
	want, err := db.ListKeywordSearches()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := db.ListKeywordSearches()
	if err != nil || !reflect.DeepEqual(want, got) {
		t.Fatal(got, err)
	}
	// SQLite EXCEPT compares every stored column, including first_seen_at and
	// all state rows. A CHECK failure makes a mismatch observable as test failure.
	db.mu.Lock()
	err = db.execLocked(`CREATE TEMP TABLE preservation_check(ok INTEGER CHECK(ok=1));
 INSERT INTO preservation_check SELECT CASE WHEN
 NOT EXISTS(SELECT * FROM settings EXCEPT SELECT * FROM fixture_settings) AND
 NOT EXISTS(SELECT * FROM fixture_settings EXCEPT SELECT * FROM settings) AND
 NOT EXISTS(SELECT * FROM seen_ads EXCEPT SELECT * FROM fixture_seen) AND
 NOT EXISTS(SELECT * FROM fixture_seen EXCEPT SELECT * FROM seen_ads) AND
 NOT EXISTS(SELECT * FROM state EXCEPT SELECT * FROM fixture_state) AND
 NOT EXISTS(SELECT * FROM fixture_state EXCEPT SELECT * FROM state)
 THEN 1 ELSE 0 END;`)
	db.mu.Unlock()
	if err != nil {
		t.Fatalf("exact rental rows changed: %v", err)
	}
	t.Log("GREEN AC007: legacy seeded settings/seen_ads first_seen_at/all state + telegram_offset exact bidirectional SQL EXCEPT equality after additive initialization/two search writes/update/reopen; stable IDs and atomic revision confirmed")
}
