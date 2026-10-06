package filter

import (
 "testing"
 "github.com/nicelight/somon-rent-watcher/internal/model"
)

func TestKeywordPriceClaimBaseline(t *testing.T) {
 bounds := []struct {name string; min,max *int}{
  {"absent",nil,nil},{"min",ptr(100),nil},{"max",nil,ptr(200)},
  {"both",ptr(100),ptr(200)},{"zero_min",ptr(0),nil},{"zero_max",nil,ptr(0)},
 }
 values := []struct{name string; price *int; currency string}{
  {"99",ptr(99),"TJS"},{"100",ptr(100),"TJS"},{"200",ptr(200),"TJS"},{"201",ptr(201),"TJS"},
  {"zero",ptr(0),"TJS"},{"missing",nil,""},{"negotiable",nil,"TJS"},
  {"foreign",ptr(150),"USD"},{"unknown_currency",ptr(150),""},
 }
 for _,b := range bounds {for _,v := range values {
  t.Run(b.name+"/"+v.name,func(t *testing.T) {
   want := true
   if b.min!=nil || b.max!=nil {
    want = v.price!=nil && v.currency=="TJS"
    if want && b.min!=nil {want = *v.price>=*b.min}
    if want && b.max!=nil {want = *v.price<=*b.max}
   }
   settings:=DefaultSettings(); settings.PriceMin=b.min; settings.PriceMax=b.max
   got,_:=CardMatches(settings,model.Card{Price:v.price,Currency:v.currency})
   if got!=want {t.Errorf("baseline card: got %v, keyword criterion wants %v",got,want)}
  })
 }}
 t.Run("commodity_without_housing_or_seller_fields",func(t *testing.T){
  settings:=DefaultSettings(); settings.PriceMin=ptr(100);settings.PriceMax=ptr(200)
  settings.Rooms=[]string{"2"};settings.Floors=[]string{"3"};settings.SellerAdsLimit=5
  commodity:=model.Ad{Card:model.Card{Title:"Стол",Price:ptr(150),Currency:"TJS"}}
  got,reason:=AdMatches(settings,commodity)
  if !got {t.Errorf("eligible commodity rejected by rental rules: %s",reason)}
 })
}
