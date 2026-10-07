package somon

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func TestVerifierDetailBodyOutcome(t *testing.T) {
	price := 725
	fallback := model.Card{ID: 481, URL: "https://somon.tj/adv/481_table/", Title: "Fallback table", Price: &price, Currency: "TJS", ImageURL: "https://images.example/table.jpg"}
	cases := []struct{ name, body string; accept, blocked bool }{
		{"foreign", `<html><body>Service unavailable</body></html>`, false, false},
		{"heading-only", `<html><body><h1>Not found</h1></body></html>`, false, false},
		{"no-own-heading", `<html><body><div class="price">725 c.</div><h2>Описание</h2></body></html>`, false, false},
		{"visible-block-with-detail", `<html><body><h1>Table</h1><div class="price">725 c.</div><p>Доступ временно ограничен</p></body></html>`, false, true},
		{"sparse-description", `<html><body><h1>Body table</h1><h2>Описание</h2><p>Мебель</p></body></html>`, true, false},
		{"sparse-id", `<html><body><h1>Body table</h1><div>ID: 481</div></body></html>`, true, false},
		{"sparse-hidden-aria", `<html><body><div aria-hidden="true">Access denied</div><h1>Body table</h1><h2>Описание</h2></body></html>`, true, false},
		{"sparse-hidden-style", `<html><body><div style="visibility: hidden">Access denied</div><h1>Body table</h1><h2>Описание</h2></body></html>`, true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ad, err := ParseDetail(fallback.URL, []byte(tc.body), fallback)
			var blocked *BlockedPageError
			if (err == nil) != tc.accept || errors.As(err, &blocked) != tc.blocked { t.Fatalf("accept=%v blocked=%v ad=%+v err=%v", tc.accept, tc.blocked, ad, err) }
			if tc.accept && (ad.ID != 481 || ad.Price == nil || *ad.Price != price || ad.Currency != "TJS" || ad.ImageURL != fallback.ImageURL || ad.Title != "Body table") { t.Fatalf("sparse merge lost validated fields/fallback: %+v", ad) }
			t.Logf("FT-005-AC-008 accepted=%v typed-block=%v error=%v", err == nil, errors.As(err, &blocked), err)
		})
	}
	for _, id := range []int{481, 482} {
		payload := `0:{"advert":{"id":` + string(mustProbeJSON(id)) + `,"title":"RSC table","description":"Мебель"}}` + "\n"
		quoted, _ := json.Marshal(payload)
		body := []byte(`<html><body><script>self.__next_f.push([1,` + string(quoted) + `])</script></body></html>`)
		ad, err := ParseDetail(fallback.URL, body, fallback)
		if (err == nil) != (id == 481) { t.Fatalf("RSC id=%d ad=%+v error=%v", id, ad, err) }
		t.Logf("FT-005-AC-008 RSC ID=%d accepted=%v", id, err == nil)
	}
}

func mustProbeJSON(v any) []byte { b, _ := json.Marshal(v); return b }
