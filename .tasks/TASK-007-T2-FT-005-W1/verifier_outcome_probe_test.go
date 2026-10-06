package filter_test

import (
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func verifierInt(n int) *int { return &n }

// Expected outcomes are taken from AC004 and PRD Monitoring rules, independently
// of the implementation and executor's tests. Each case creates fresh values.
func TestVerifierKeywordPriceOutcome(t *testing.T) {
	profiles := []struct {
		name string
		min, max *int
		allowed [5]bool // 0, 99, 100, 200, 201 TJS
	}{
		{"unbounded", nil, nil, [5]bool{true,true,true,true,true}},
		{"minimum", verifierInt(100), nil, [5]bool{false,false,true,true,true}},
		{"maximum", nil, verifierInt(200), [5]bool{true,true,true,true,false}},
		{"range", verifierInt(100), verifierInt(200), [5]bool{false,false,true,true,false}},
		{"zero_minimum", verifierInt(0), nil, [5]bool{true,true,true,true,true}},
		{"zero_maximum", nil, verifierInt(0), [5]bool{true,false,false,false,false}},
	}
	for _, p := range profiles {
		t.Run(p.name, func(t *testing.T) {
			for i, amount := range []int{0,99,100,200,201} {
				if got := filter.KeywordPriceMatches(p.min,p.max,verifierInt(amount),"TJS"); got != p.allowed[i] {
					t.Errorf("price=%d TJS: got %v, expected %v",amount,got,p.allowed[i])
				}
			}
			for _, v := range []struct{name string; price *int; currency string}{
				{"missing",nil,""}, {"negotiable",nil,"TJS"},
				{"foreign",verifierInt(150),"USD"}, {"unconfirmed_currency",verifierInt(150),""},
			} {
				want := p.name == "unbounded"
				if got := filter.KeywordPriceMatches(p.min,p.max,v.price,v.currency); got != want {
					t.Errorf("%s: got %v, expected %v",v.name,got,want)
				}
			}
		})
	}
	for _, amount := range []int{0,100} {
		if !filter.KeywordPriceMatches(verifierInt(amount),verifierInt(amount),verifierInt(amount),"TJS") {
			t.Errorf("equal bounds reject their exact price %d",amount)
		}
	}
}

func TestVerifierKeywordBoundValidation(t *testing.T) {
	for _, c := range []struct{name string; min,max *int; valid bool}{
		{"none",nil,nil,true},{"minimum",verifierInt(100),nil,true},
		{"maximum",nil,verifierInt(200),true},{"range",verifierInt(100),verifierInt(200),true},
		{"zero_minimum",verifierInt(0),nil,true},{"zero_maximum",nil,verifierInt(0),true},
		{"equal_zero",verifierInt(0),verifierInt(0),true},{"equal",verifierInt(100),verifierInt(100),true},
		{"negative_minimum",verifierInt(-1),nil,false},{"negative_maximum",nil,verifierInt(-1),false},
		{"reversed",verifierInt(201),verifierInt(200),false},
	} {
		t.Run(c.name,func(t *testing.T){
			if got := filter.ValidateKeywordPriceBounds(c.min,c.max); (got == nil) != c.valid {
				t.Errorf("valid=%v expected %v; error=%v",got==nil,c.valid,got)
			}
			if !c.valid && filter.KeywordPriceMatches(c.min,c.max,verifierInt(150),"TJS") {
				t.Error("invalid range admits commodity")
			}
		})
	}
}

func TestVerifierCommodityOutcome(t *testing.T) {
	for _, amount := range []int{100,150,200} {
		commodity := model.Ad{Card:model.Card{Title:"Стол",Price:verifierInt(amount),Currency:"TJS"}}
		if !filter.KeywordPriceMatches(verifierInt(100),verifierInt(200),commodity.Price,commodity.Currency) {
			t.Errorf("commodity with no housing/seller fields rejected at %d",amount)
		}
		commodity.Rooms=verifierInt(9); commodity.Floor=verifierInt(99)
		commodity.Promoted=true; commodity.Description="Посуточно без детей"
		commodity.SellerAds=verifierInt(100)
		if !filter.KeywordPriceMatches(verifierInt(100),verifierInt(200),commodity.Price,commodity.Currency) {
			t.Errorf("apartment/text/seller/promotion fields affect commodity at %d",amount)
		}
	}
}
