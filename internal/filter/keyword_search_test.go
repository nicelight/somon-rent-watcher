package filter

import (
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func TestKeywordPriceMatches(t *testing.T) {
	// Columns: 99, 100, 200, 201, zero, missing, negotiable, USD, unknown currency.
	bounds := []struct {
		name     string
		min, max *int
		want     [9]bool
	}{
		{"absent", nil, nil, [9]bool{true, true, true, true, true, true, true, true, true}},
		{"min", ptr(100), nil, [9]bool{false, true, true, true, false, false, false, false, false}},
		{"max", nil, ptr(200), [9]bool{true, true, true, false, true, false, false, false, false}},
		{"both", ptr(100), ptr(200), [9]bool{false, true, true, false, false, false, false, false, false}},
		{"zero_min", ptr(0), nil, [9]bool{true, true, true, true, true, false, false, false, false}},
		{"zero_max", nil, ptr(0), [9]bool{false, false, false, false, true, false, false, false, false}},
	}
	for _, b := range bounds {
		values := []struct {
			name     string
			price    *int
			currency string
		}{
			{"99", ptr(99), "TJS"}, {"100", ptr(100), "TJS"},
			{"200", ptr(200), "TJS"}, {"201", ptr(201), "TJS"},
			{"zero", ptr(0), "TJS"}, {"missing", nil, ""}, {"negotiable", nil, "TJS"},
			{"foreign", ptr(150), "USD"}, {"unknown_currency", ptr(150), ""},
		}
		for i, v := range values {
			t.Run(b.name+"/"+v.name, func(t *testing.T) {
				if got := KeywordPriceMatches(b.min, b.max, v.price, v.currency); got != b.want[i] {
					t.Errorf("got %v, want %v", got, b.want[i])
				}
			})
		}
	}
}

func TestKeywordCommodityIgnoresApartmentFields(t *testing.T) {
	for _, card := range []model.Card{
		{Title: "Стол", Price: ptr(150), Currency: "TJS"},
		{Title: "Стол", Price: ptr(150), Currency: "TJS", Rooms: ptr(9), Floor: ptr(99), Promoted: true},
	} {
		ad := model.Ad{Card: card, Description: "Посуточно без детей", SellerAds: ptr(100)}
		if !KeywordPriceMatches(ptr(100), ptr(200), ad.Price, ad.Currency) {
			t.Fatal("eligible commodity rejected by apartment/seller/text/promotion fields")
		}
	}
}

func TestValidateKeywordPriceBounds(t *testing.T) {
	for _, tc := range []struct {
		name     string
		min, max *int
		wantErr  bool
	}{
		{"absent", nil, nil, false}, {"min", ptr(100), nil, false},
		{"max", nil, ptr(200), false}, {"range", ptr(100), ptr(200), false},
		{"zero_min", ptr(0), nil, false}, {"zero_max", nil, ptr(0), false},
		{"equal_zero", ptr(0), ptr(0), false}, {"equal", ptr(100), ptr(100), false},
		{"negative_min", ptr(-1), nil, true}, {"negative_max", nil, ptr(-1), true},
		{"reversed", ptr(201), ptr(200), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidateKeywordPriceBounds(tc.min, tc.max); (got != nil) != tc.wantErr {
				t.Errorf("err=%v, want error=%v", got, tc.wantErr)
			}
			if tc.wantErr && KeywordPriceMatches(tc.min, tc.max, ptr(150), "TJS") {
				t.Fatal("invalid bounds matched")
			}
		})
	}
}
