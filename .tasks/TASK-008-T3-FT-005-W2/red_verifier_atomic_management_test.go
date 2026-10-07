//go:build cgo

package app

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/store"
)

func TestRedVerifierManagementAtomicMutations(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "search.db"))
	if err != nil { t.Fatal(err) }
	defer db.Close()
	a, err := New(config.Config{PollMin: time.Minute, PollMax: 2*time.Minute}, db, nil)
	if err != nil { t.Fatal(err) }
	other, err := a.CreateKeywordSearch(model.KeywordSearch{Phrase:"другой поиск",CategoryKey:"services",CityKey:"dangara"})
	if err != nil { t.Fatal(err) }
	read := func(id int64) model.KeywordSearch {
		t.Helper()
		s, found, err := a.KeywordSearch(id)
		if err != nil || !found { t.Fatal(found,err) }
		return s
	}
	type result struct { s model.KeywordSearch; found bool; err error }
	for i := 0; i < 16; i++ {
		s, err := a.CreateKeywordSearch(model.KeywordSearch{Phrase:"стол",CategoryKey:"all",CityKey:"vose"})
		if err != nil { t.Fatal(err) }
		start := make(chan struct{})
		out := make(chan result,8)
		for j := 0; j < 8; j++ {
			patch := s
			patch.Phrase = fmt.Sprintf("стол %d",j)
			go func(p model.KeywordSearch) {
				<-start
				v,ok,e := a.UpdateKeywordSearch(p)
				out <- result{v,ok,e}
			}(patch)
		}
		close(start)
		wins := 0
		var winner model.KeywordSearch
		for j := 0; j < 8; j++ {
			r := <-out
			if !r.found { t.Fatal("supported search disappeared") }
			if r.err == nil { wins++; winner = r.s }
		}
		if wins != 1 || winner.Revision != 2 || !reflect.DeepEqual(winner,read(s.ID)) {
			t.Fatalf("expected one complete revision commit: wins=%d winner=%+v",wins,winner)
		}
		patch := winner
		patch.Phrase = "новые условия"
		start = make(chan struct{})
		go func() { <-start; v,ok,e := a.UpdateKeywordSearch(patch); out <- result{v,ok,e} }()
		go func() { <-start; v,ok,e := a.SetKeywordSearchEnabled(s.ID,true); out <- result{v,ok,e} }()
		close(start)
		first,second := <-out,<-out
		current := read(s.ID)
		if !first.found || !second.found || !current.Enabled || current.CityKey != "vose" || current.CategoryKey != "all" {
			t.Fatalf("concurrent enable/settings lost addressed fields: %+v %+v %+v",first,second,current)
		}
		if current.Phrase != patch.Phrase {
			if first.err == nil && second.err == nil { t.Fatal("success concealed missing settings") }
			current.Phrase = patch.Phrase
			_,ok,e := a.UpdateKeywordSearch(current)
			if e != nil || !ok { t.Fatal("stale patch cannot recover with current revision",e,ok) }
		}
		current = read(s.ID)
		if current.Revision != 4 || !current.Enabled || current.Phrase != patch.Phrase { t.Fatal(current) }
		if !reflect.DeepEqual(other,read(other.ID)) { t.Fatal("mutation altered another search") }
	}
	t.Log("16 rounds: eight competing expected-revision patches commit exactly one; concurrent enable/edit retains complete values and retry reaches revision4; separate search untouched")
}
