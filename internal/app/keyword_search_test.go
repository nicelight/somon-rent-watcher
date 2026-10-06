//go:build cgo

package app

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/store"
)

func TestKeywordApplicationOwnsValidationAndRevision(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "search.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	a, err := New(config.Config{PollMin: time.Minute, PollMax: 2 * time.Minute}, db, nil)
	if err != nil {
		t.Fatal(err)
	}
	base := model.KeywordSearch{Phrase: "  стол  ", CategoryKey: "all", CityKey: "country", ID: 500, Enabled: true, Revision: 900}
	s, err := a.CreateKeywordSearch(base)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID == 500 || s.ID <= 0 || s.Enabled || s.Revision != 1 || s.Phrase != "стол" {
		t.Fatal(s)
	}
	original, _ := a.ListKeywordSearches()
	neg, invertedMin, invertedMax := -1, 200, 100
	invalid := []model.KeywordSearch{base, base, base, base, base}
	invalid[0].Phrase = " "
	invalid[1].CategoryKey = "/arbitrary/"
	invalid[2].CityKey = "invalid"
	invalid[3].PriceMin = &neg
	invalid[4].PriceMin = &invertedMin
	invalid[4].PriceMax = &invertedMax
	for _, bad := range invalid {
		if _, err = a.CreateKeywordSearch(bad); err == nil {
			t.Fatal("invalid create accepted", bad)
		}
		bad.ID = s.ID
		bad.Revision = s.Revision
		if _, _, err = a.UpdateKeywordSearch(bad); err == nil {
			t.Fatal("invalid update accepted", bad)
		}
		current, _ := a.ListKeywordSearches()
		if !reflect.DeepEqual(original, current) {
			t.Fatal("invalid values changed DB")
		}
	}
	old := s
	s.Phrase = "стулья"
	s.Enabled = true
	s, found, err := a.UpdateKeywordSearch(s)
	if err != nil || !found || s.Revision != 2 || s.Enabled {
		t.Fatal(s, found, err)
	}
	if _, _, err = a.UpdateKeywordSearch(old); err == nil {
		t.Fatal("stale revision accepted")
	}
	s, found, err = a.SetKeywordSearchEnabled(s.ID, true)
	if err != nil || !found || !s.Enabled || s.Revision != 3 {
		t.Fatal(s, found, err)
	}
	again, found, err := a.SetKeywordSearchEnabled(s.ID, true)
	if err != nil || !found || !reflect.DeepEqual(again, s) {
		t.Fatal("repeated enable changed revision")
	}
	if _, found, err = a.SetKeywordSearchEnabled(999999, true); err != nil || found {
		t.Fatal("missing enable created row")
	}
	if _, found, err = a.UpdateKeywordSearch(model.KeywordSearch{ID: 999999}); err != nil || found {
		t.Fatal("missing update created row")
	}
}
