package filter

import (
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func ptr(v int) *int { return &v }

func TestPriceRoomsFloorsAndSellerFilters(t *testing.T) {
	s := DefaultSettings()
	s.PriceMin, s.PriceMax = ptr(3500), ptr(6000)
	s.Rooms = []string{"2", UnknownChoice}
	s.Floors = []string{"2", "14+", UnknownChoice}
	s.SellerAdsLimit = 5

	ad := model.Ad{Card: model.Card{Price: ptr(6000), Rooms: ptr(2), Floor: ptr(14)}, SellerAds: ptr(4)}
	if ok, reason := AdMatches(s, ad); !ok {
		t.Fatalf("expected match: %s", reason)
	}
	ad.Price = ptr(6001)
	if ok, _ := AdMatches(s, ad); ok {
		t.Fatal("price above inclusive max matched")
	}
	ad.Price = ptr(3500)
	ad.SellerAds = ptr(5)
	if ok, _ := AdMatches(s, ad); ok {
		t.Fatal("seller count equal to <5 limit matched")
	}
	ad.SellerAds = nil
	if ok, _ := AdMatches(s, ad); ok {
		t.Fatal("unknown seller count matched active limit")
	}
	s.SellerAdsLimit = 0
	if ok, reason := AdMatches(s, ad); !ok {
		t.Fatalf("unknown seller should match 'not important': %s", reason)
	}
}

func TestUnknownAndPrefilterSemantics(t *testing.T) {
	s := DefaultSettings()
	s.PriceMax = ptr(5000)
	s.Rooms = []string{"2"}
	s.Floors = []string{"3"}
	card := model.Card{ID: 1}
	if ok, reason := CardMatches(s, card); !ok {
		t.Fatalf("unknown card values must reach detail: %s", reason)
	}
	ad := model.Ad{Card: card}
	if ok, _ := AdMatches(s, ad); ok {
		t.Fatal("unknown detail price must fail when price filter is active")
	}
}

func TestNegativeWordsAreSimpleNormalizedSubstrings(t *testing.T) {
	s := DefaultSettings()
	s.NegativeWords = []string{"без детей", "ПОСУТОЧНО"}
	if err := s.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}
	ad := model.Ad{Card: model.Card{Title: "Квартира"}, Description: "Сдаётся   БЕЗ\nДЕТЕЙ."}
	if ok, reason := AdMatches(s, ad); ok || reason != "негативное слово: без детей" {
		t.Fatalf("got ok=%v reason=%q", ok, reason)
	}
}

func TestParseInputsAndJSON(t *testing.T) {
	min, max, err := ParsePriceInput("3 500—6000")
	if err != nil || min == nil || *min != 3500 || max == nil || *max != 6000 {
		t.Fatalf("price parse: %v %v %v", min, max, err)
	}
	min, max, err = ParsePriceInput("-")
	if err != nil || min != nil || max != nil {
		t.Fatalf("clear price: %v %v %v", min, max, err)
	}
	words, err := ParseNegativeInput(" Посуточно, без детей, посуточно ")
	if err != nil || len(words) != 2 || words[0] != "посуточно" {
		t.Fatalf("words=%v err=%v", words, err)
	}

	s := DefaultSettings()
	s.Rooms = []string{"3", "1"}
	raw, err := s.Encode()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeSettings(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded.Rooms) != 2 || decoded.Rooms[0] != "1" || decoded.Rooms[1] != "3" {
		t.Fatalf("normalized rooms=%v", decoded.Rooms)
	}
}

func TestCannotDisableLastOption(t *testing.T) {
	s := DefaultSettings()
	s.Rooms = []string{"1"}
	if err := s.ToggleRoom("1"); err == nil {
		t.Fatal("expected error")
	}
	s.IncludeOrdinary = true
	s.IncludePromoted = false
	if err := s.ToggleOrdinary(); err == nil || !s.IncludeOrdinary {
		t.Fatalf("last type must remain selected: %+v, %v", s, err)
	}
}
