package somon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestParseCategoryDOM(t *testing.T) {
	cards, err := ParseCategory(DefaultCategoryURL, fixture(t, "category.html"))
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 3 {
		t.Fatalf("cards=%d IDs=%v", len(cards), sortedCardIDs(cards))
	}
	if cards[0].ID != 17000001 || cards[1].ID != 17000002 || cards[2].ID != 17000003 {
		t.Fatalf("unexpected order/IDs: %v", sortedCardIDs(cards))
	}
	vip := cards[0]
	if !vip.Promoted || vip.PromotionLabel != "VIP" {
		t.Fatalf("VIP not detected: %+v", vip)
	}
	if vip.Price == nil || *vip.Price != 5999 {
		t.Fatalf("current discounted price=%v", vip.Price)
	}
	if vip.Rooms == nil || *vip.Rooms != 2 || vip.Floor == nil || *vip.Floor != 4 {
		t.Fatalf("visible rooms/floor not used: %+v", vip)
	}
	if !strings.Contains(vip.ImageURL, "17000001.jpg") || vip.AgeText != "3 минуты назад" {
		t.Fatalf("image/age: %+v", vip)
	}
	if !cards[1].Promoted || cards[1].PromotionLabel != "ТОП" {
		t.Fatalf("TOP not detected: %+v", cards[1])
	}
	if cards[2].Promoted || cards[2].Price == nil || *cards[2].Price != 4200 {
		t.Fatalf("ordinary card: %+v", cards[2])
	}
	for _, card := range cards {
		if card.ID == 17999999 {
			t.Fatal("card after 'other cities' boundary leaked into Dushanbe feed")
		}
	}
}

func TestParseDetailUsesVisibleDataNotStaleURLSlug(t *testing.T) {
	fallback := model.Card{
		ID:       17000001,
		URL:      "https://somon.tj/adv/17000001_2-komn-kvartira-12-etazh-79m2-somoni/",
		Title:    "2-комн. квартира, 12 этаж, 79м², Сомони",
		Floor:    intPtr(12),
		Promoted: true,
	}
	ad, err := ParseDetail(fallback.URL, fixture(t, "detail.html"), fallback)
	if err != nil {
		t.Fatal(err)
	}
	if ad.ID != 17000001 || ad.Floor == nil || *ad.Floor != 4 {
		t.Fatalf("visible floor must override stale URL/fallback: %+v", ad)
	}
	if ad.Price == nil || *ad.Price != 5999 {
		t.Fatalf("price=%v", ad.Price)
	}
	if !strings.Contains(ad.Description, "уютная квартира") || !strings.Contains(ad.Description, "посуточной") {
		t.Fatalf("description=%q", ad.Description)
	}
	if ad.SellerName != "Rofi Estate R.M" || ad.SellerSince != "Июля 2025" || ad.SellerAds == nil || *ad.SellerAds != 17 {
		t.Fatalf("seller parse: name=%q since=%q ads=%v", ad.SellerName, ad.SellerSince, ad.SellerAds)
	}
	if !strings.Contains(ad.ImageURL, "detail-17000001.jpg") {
		t.Fatalf("image=%q", ad.ImageURL)
	}
}

func TestParseDetailUnknownSellerIsAllowed(t *testing.T) {
	url := "https://somon.tj/adv/17000003_1-komn-kvartira-2-etazh-45m2/"
	ad, err := ParseDetail(url, fixture(t, "detail_unknown_seller.html"), model.Card{ID: 17000003, URL: url})
	if err != nil {
		t.Fatal(err)
	}
	if ad.SellerAds != nil || ad.SellerName != "" {
		t.Fatalf("seller should be unknown: %+v", ad)
	}
}

func TestHiddenBlockedModalDoesNotRejectValidCards(t *testing.T) {
	body := []byte(`<html><body><div hidden>Ваш аккаунт был заблокирован на somon.tj</div><article><span>4000 c.</span><a href="/adv/12345678_x/">1-комн. квартира, 2 этаж, 40м²</a></article></body></html>`)
	cards, err := ParseCategory(DefaultCategoryURL, body)
	if err != nil || len(cards) != 1 {
		t.Fatalf("cards=%v err=%v", cards, err)
	}
}

func TestPromotionClassDoesNotTreatTailwindTopPositionAsPromoted(t *testing.T) {
	body := []byte(`<html><body>
		<article class="advert-card absolute top-4"><span>4000 c.</span><a href="/adv/12345678_x/">1-комн. квартира, 2 этаж, 40м²</a></article>
		<article class="advert-card bg-top"><span>5000 c.</span><a href="/adv/12345679_x/">2-комн. квартира, 3 этаж, 60м²</a></article>
	</body></html>`)
	cards, err := ParseCategory(DefaultCategoryURL, body)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 2 {
		t.Fatalf("cards=%v", cards)
	}
	if cards[0].Promoted {
		t.Fatalf("Tailwind top positioning must stay ordinary: %+v", cards[0])
	}
	if !cards[1].Promoted || cards[1].PromotionLabel != "ТОП" {
		t.Fatalf("bg-top promotion not detected: %+v", cards[1])
	}
}

func TestBlockedPageWithoutCardsFailsAsBlocked(t *testing.T) {
	_, err := ParseCategory(DefaultCategoryURL, []byte(`<html><body><h1>Access denied</h1></body></html>`))
	var blocked *BlockedPageError
	if !errors.As(err, &blocked) {
		t.Fatalf("expected BlockedPageError, got %T %v", err, err)
	}
}

func TestParseCategoryRSCFallback(t *testing.T) {
	payload := `0:{"adverts":[{"id":18000001,"url":"/adv/18000001_x/","title":"3-комн. квартира, 15 этаж, 90м²","price":7000,"first_thumb":"//cdntj.somon.tj/a.jpg","is_top":true,"published":"7 минут назад"}]}`
	quoted, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte(`<html><body><script>self.__next_f.push([1,` + string(quoted) + `])</script></body></html>`)
	cards, err := ParseCategory(DefaultCategoryURL, body)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 {
		t.Fatalf("cards=%v", cards)
	}
	card := cards[0]
	if card.ID != 18000001 || card.Price == nil || *card.Price != 7000 || card.Floor == nil || *card.Floor != 15 || !card.Promoted {
		t.Fatalf("RSC card=%+v", card)
	}
}

func TestAgeAndRecoveryURLHelpers(t *testing.T) {
	cases := []struct {
		input string
		want  time.Duration
	}{
		{"Сегодня", 0},
		{"Вчера", 24 * time.Hour},
		{"2 недели назад", 14 * 24 * time.Hour},
		{"15 минут назад", 15 * time.Minute},
	}
	for _, tc := range cases {
		got, ok := ParseAgeDuration(tc.input)
		if !ok || got != tc.want {
			t.Fatalf("ParseAgeDuration(%q)=%s,%v want %s", tc.input, got, ok, tc.want)
		}
	}
	urls := RecoveryURLs(DefaultCategoryURL, []string{"1", "6+", "unknown"})
	if len(urls) != 2 || !strings.Contains(urls[1], "/6-i-bolee-komnat/dushanbe/") {
		t.Fatalf("recovery URLs=%v", urls)
	}
}

func intPtr(v int) *int { return &v }
