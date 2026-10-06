package somon

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func keywordFixture(t *testing.T, name string) []byte {
	t.Helper()
	body, err := os.ReadFile("../../testdata/keyword-search-primary-" + name + ".html")
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func keywordIDs(cards []model.Card) []int64 {
	ids := []int64{}
	for _, card := range cards {
		ids = append(ids, card.ID)
	}
	return ids
}

func TestParseKeywordPrimaryFixtures(t *testing.T) {
	for _, tc := range []struct {
		name string
		want []int64
	}{
		{"small", []int64{21000001}},
		{"empty", []int64{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := ParseKeywordSearch("https://somon.tj/search/vose/?q=стол", keywordFixture(t, tc.name))
			if err != nil {
				t.Fatal(err)
			}
			if ids := keywordIDs(cards); !reflect.DeepEqual(ids, tc.want) {
				t.Fatalf("primary IDs=%v want %v", ids, tc.want)
			}
			if len(cards) > 0 {
				card := cards[0]
				if card.Title != "Стол" || card.Price == nil || *card.Price != 200 || card.Currency != "TJS" || card.City != "Восе" || card.ImageURL != "https://somon.tj/primary.jpg" || card.Position != 0 || card.Rooms != nil {
					t.Fatalf("passive primary fields=%+v", card)
				}
			}
		})
	}
}

func TestParseKeywordMultiplePrimaryAndNoLocalMatching(t *testing.T) {
	body := `<h1>Мебель Душанбе 2</h1><article><a href="/adv/22000001_table/">Стол</a><a href="/adv/22000001_table/">Стол</a><span class="price">0 c.</span></article><article><a href="/adv/22000002_chair/">Стул</a></article><h2>Объявления <span>из других регионов</span></h2><article><a href="/adv/22000003_other/">Стол</a></article>`
	cards, err := ParseKeywordSearch("https://somon.tj/search/dushanbe/?q=стол", []byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{22000001, 22000002}; !reflect.DeepEqual(keywordIDs(cards), want) {
		t.Fatalf("primary=%v want %v", keywordIDs(cards), want)
	}
	if cards[0].Price == nil || *cards[0].Price != 0 || cards[0].Currency != "TJS" || cards[0].City != "" || cards[1].Price != nil || cards[1].Currency != "" || cards[1].ImageURL != "" {
		t.Fatalf("missing fields/zero amount must be preserved: %+v", cards)
	}
}

func TestParseKeywordFailClosed(t *testing.T) {
	for name, body := range map[string]string{
		"malformed":                    "broken body",
		"unconfirmed empty":            "<h1>Поиск стол</h1>",
		"recommendations only":         `<h1>Поиск стол</h1><h2>Объявления из других регионов</h2><a href="/adv/23000001_x/">Стол</a>`,
		"RSC cannot establish primary": `<script>self.__next_f.push([1,"0:{\"adverts\":[{\"id\":23000001,\"url\":\"/adv/23000001_x/\",\"title\":\"Стол\"}]}"])</script>`,
		"contradictory zero":           `<h1>Поиск стол 0</h1><a href="/adv/23000001_x/">Стол</a>`,
		"untitled card":                `<h1>Поиск стол 1</h1><a href="/adv/23000001_x/"></a>`,
	} {
		t.Run(name, func(t *testing.T) {
			if cards, err := ParseKeywordSearch(DefaultCategoryURL, []byte(body)); err == nil {
				t.Fatalf("expected parse error, got %+v", cards)
			}
		})
	}
	_, err := ParseKeywordSearch(DefaultCategoryURL, []byte(`<h1>Access denied</h1><h1>Поиск 0</h1>`))
	var blocked *BlockedPageError
	if !errors.As(err, &blocked) {
		t.Fatalf("block must remain typed: %T %v", err, err)
	}
}

func TestKeywordCatalogAndURL(t *testing.T) {
	wantCategories := []KeywordCategory{
		{"all", "Все категории", "/search/"},
		{"furniture", "Мебель", "/vse-dlya-doma/mebel/"},
		{"tables_chairs", "Столы и стулья", "/vse-dlya-doma/mebel/mebel-dlya-kuhni/stolyi-stulya/"},
		{"services", "Услуги", "/biznes-i-uslugi/"},
	}
	wantCities := []KeywordCity{{"country", "Вся страна", ""}, {"dushanbe", "Душанбе", "dushanbe"}, {"vose", "Восе", "vose"}, {"dangara", "Дангара", "dangara"}}
	if !reflect.DeepEqual(KeywordCategories(), wantCategories) || !reflect.DeepEqual(KeywordCities(), wantCities) {
		t.Fatal("catalog differs from accepted source contract")
	}
	copyCategories := KeywordCategories()
	copyCategories[0].Path = "/evil/"
	if KeywordCategories()[0].Path != "/search/" {
		t.Fatal("caller mutated source catalog")
	}
	for _, category := range wantCategories {
		for _, city := range wantCities {
			raw, err := KeywordSearchURL("  стол & стул + /?  ", category.Key, city.Key)
			if err != nil {
				t.Fatal(err)
			}
			u, _ := url.Parse(raw)
			path := category.Path
			if city.Slug != "" {
				path += city.Slug + "/"
			}
			if u.Host != "somon.tj" || u.Scheme != "https" || u.Path != path || u.Query().Get("q") != "стол & стул + /?" || u.Query().Get("ordering") != "relevance" || len(u.Query()) != 2 {
				t.Fatalf("wrong native request: %s", raw)
			}
		}
	}
	for _, tc := range [][3]string{{" ", "all", "country"}, {"стол", "https://evil.test/", "country"}, {"стол", "all", "../dushanbe"}, {"стол", "", "country"}, {"стол", "all", ""}} {
		if _, err := KeywordSearchURL(tc[0], tc[1], tc[2]); err == nil {
			t.Fatalf("accepted invalid request %v", tc)
		}
	}
}

type keywordLocalTransport struct {
	target *url.URL
	seen   func(*http.Request)
}

func (transport keywordLocalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport.seen(req)
	local := req.Clone(req.Context())
	local.URL.Scheme, local.URL.Host = transport.target.Scheme, transport.target.Host
	return http.DefaultTransport.RoundTrip(local)
}

func TestFetchKeywordSearchNativeScope(t *testing.T) {
	body := keywordFixture(t, "small")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(body) }))
	defer server.Close()
	endpoint, _ := url.Parse(server.URL)
	client := NewClient("keyword-test", 0, time.Second, 1<<20)
	var captured *url.URL
	client.httpClient.Transport = keywordLocalTransport{endpoint, func(r *http.Request) {
		value := *r.URL
		captured = &value
		if r.Header.Get("User-Agent") != "keyword-test" {
			t.Error("existing request policy not reused")
		}
	}}
	for _, category := range KeywordCategories() {
		for _, city := range KeywordCities() {
			cards, raw, err := client.FetchKeywordSearch(context.Background(), "стол & стул", category.Key, city.Key)
			if err != nil || !reflect.DeepEqual(keywordIDs(cards), []int64{21000001}) || string(raw) != string(body) {
				t.Fatalf("fetch primary=%v err=%v", keywordIDs(cards), err)
			}
			path := category.Path
			if city.Slug != "" {
				path += city.Slug + "/"
			}
			if captured.Host != "somon.tj" || captured.Path != path || captured.Query().Get("q") != "стол & стул" || captured.Query().Get("ordering") != "relevance" {
				t.Fatalf("captured native scope=%s", captured)
			}
		}
	}
}

func TestFetchKeywordErrorsAndEmpty(t *testing.T) {
	for _, tc := range []struct {
		name        string
		body        []byte
		status      int
		wantError   bool
		wantBlocked bool
	}{
		{"confirmed empty", keywordFixture(t, "empty"), 200, false, false},
		{"unconfirmed empty", []byte("<h1>Поиск</h1>"), 200, true, false},
		{"blocked body", []byte("<h1>Access denied</h1>"), 200, true, true},
		{"403", nil, 403, true, true},
		{"429", nil, 429, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(tc.status)
				w.Write(tc.body)
			}))
			defer server.Close()
			endpoint, _ := url.Parse(server.URL)
			client := NewClient("test", 0, time.Second, 1<<20)
			client.httpClient.Transport = keywordLocalTransport{endpoint, func(*http.Request) {}}
			cards, raw, err := client.FetchKeywordSearch(context.Background(), "стол", "all", "country")
			if (err != nil) != tc.wantError {
				t.Fatalf("error=%v wantError=%v", err, tc.wantError)
			}
			if len(cards) != 0 {
				t.Fatalf("unexpected cards=%v", cards)
			}
			_, retry, blocked := IsBlocked(err)
			if blocked != tc.wantBlocked || tc.status == 429 && retry != time.Minute {
				t.Fatalf("typed block=%v retry=%s error=%v", blocked, retry, err)
			}
			if tc.status == 200 && string(raw) != string(tc.body) {
				t.Fatal("raw diagnostic body lost")
			}
		})
	}
}

func TestKeywordMoneyEvidence(t *testing.T) {
	for _, tc := range []struct {
		name, markup string
		amount       *int
		currency     string
	}{
		{"structured TJS", `<meta itemprop="price" content="200"><meta itemprop="priceCurrency" content="TJS">`, intPtr(200), "TJS"},
		{"structured unknown", `<meta itemprop="price" content="200">`, nil, ""},
		{"foreign currency", `<meta itemprop="price" content="200"><meta itemprop="priceCurrency" content="USD">`, nil, "USD"},
		{"negotiable", `<span class="price">Договорная</span>`, nil, ""},
		{"invalid negative", `<span class="price">-200 c.</span>`, nil, ""},
		{"photo count", `<div><span>6</span><img src="/image.jpg"></div><span class="price">5 000 c.</span>`, intPtr(5000), "TJS"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := ParseKeywordSearch(DefaultCategoryURL, []byte(`<article><a href="/adv/24000001_x/">Товар</a>`+tc.markup+`</article>`))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cards[0].Price, tc.amount) || cards[0].Currency != tc.currency || strings.Contains(cards[0].Title, "квартир") {
				t.Fatalf("money=%+v", cards[0])
			}
		})
	}
}
