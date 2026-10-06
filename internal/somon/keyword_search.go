package somon

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/nicelight/somon-rent-watcher/internal/htmlx"
	"github.com/nicelight/somon-rent-watcher/internal/model"
)

// KeywordCategory and KeywordCity describe the deliberately bounded native catalog.
type KeywordCategory struct {
	Key, Label, Path string
}

type KeywordCity struct {
	Key, Label, Slug string
}

// Return fresh slices so callers cannot mutate the source-owned catalog.
func KeywordCategories() []KeywordCategory {
	return []KeywordCategory{
		{"all", "Все категории", "/search/"},
		{"furniture", "Мебель", "/vse-dlya-doma/mebel/"},
		{"tables_chairs", "Столы и стулья", "/vse-dlya-doma/mebel/mebel-dlya-kuhni/stolyi-stulya/"},
		{"services", "Услуги", "/biznes-i-uslugi/"},
	}
}

func KeywordCities() []KeywordCity {
	return []KeywordCity{
		{"country", "Вся страна", ""},
		{"dushanbe", "Душанбе", "dushanbe"},
		{"vose", "Восе", "vose"},
		{"dangara", "Дангара", "dangara"},
	}
}

// KeywordSearchURL uses only native allowlisted scope and URL-encoded matching.
func KeywordSearchURL(phrase, categoryKey, cityKey string) (string, error) {
	phrase = strings.TrimSpace(phrase)
	if phrase == "" {
		return "", fmt.Errorf("keyword search phrase is empty")
	}
	var path string
	for _, category := range KeywordCategories() {
		if category.Key == categoryKey {
			path = category.Path
			break
		}
	}
	if path == "" {
		return "", fmt.Errorf("unsupported keyword category %q", categoryKey)
	}
	cityFound := false
	for _, city := range KeywordCities() {
		if city.Key == cityKey {
			cityFound = true
			if city.Slug != "" {
				path += city.Slug + "/"
			}
			break
		}
	}
	if !cityFound {
		return "", fmt.Errorf("unsupported keyword city %q", cityKey)
	}
	u := url.URL{Scheme: "https", Host: "somon.tj", Path: path}
	u.RawQuery = url.Values{"q": {phrase}, "ordering": {"relevance"}}.Encode()
	return u.String(), nil
}

func (c *Client) FetchKeywordSearch(ctx context.Context, phrase, categoryKey, cityKey string) ([]model.Card, []byte, error) {
	pageURL, err := KeywordSearchURL(phrase, categoryKey, cityKey)
	if err != nil {
		return nil, nil, err
	}
	body, err := c.get(ctx, pageURL, "")
	if err != nil {
		return nil, nil, err
	}
	cards, err := ParseKeywordSearch(pageURL, body)
	return cards, body, err
}

var keywordZeroHeadingRE = regexp.MustCompile(`\S.*\s0$`)

// ParseKeywordSearch keeps membership in the visible primary section. RSC data
// cannot establish scope by itself: its largest array can be recommendations.
// Rental parsing and its separate sanity policy are intentionally unchanged.
func ParseKeywordSearch(pageURL string, body []byte) ([]model.Card, error) {
	root, err := htmlx.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse keyword HTML: %w", err)
	}
	primary := keywordPrimaryDOM(root)
	zero := htmlx.FindFirst(primary, func(n *htmlx.Node) bool {
		return n.Tag == "h1" && keywordZeroHeadingRE.MatchString(htmlx.Text(n))
	}) != nil
	anchors := htmlx.FindAll(primary, func(n *htmlx.Node) bool {
		return n.Tag == "a" && adIDFromURL(n.GetAttr("href")) > 0
	})
	if len(anchors) == 0 {
		if reason := blockedReason(root); reason != "" {
			return nil, &BlockedPageError{URL: pageURL, Reason: reason}
		}
		if zero {
			return []model.Card{}, nil
		}
		return nil, fmt.Errorf("parse keyword HTML: primary cards or confirmed empty marker not found")
	}
	if zero {
		return nil, fmt.Errorf("parse keyword HTML: zero count contradicts primary cards")
	}
	cards := make([]model.Card, 0, len(anchors))
	seen := make(map[int64]bool)
	for _, anchor := range anchors {
		id := adIDFromURL(anchor.GetAttr("href"))
		if seen[id] {
			continue
		}
		title := firstNonEmpty(anchor.GetAttr("title"), htmlx.Text(anchor))
		if title == "" {
			if img := htmlx.FindFirst(anchor, func(n *htmlx.Node) bool { return n.Tag == "img" }); img != nil {
				title = img.GetAttr("alt")
			}
		}
		if strings.TrimSpace(title) == "" {
			continue
		}
		container := keywordCardContainer(anchor)
		price, currency := keywordDOMMoney(container)
		card := model.Card{
			ID: id, URL: normalizeAdURL(pageURL, anchor.GetAttr("href")),
			Title: htmlx.NormalizeSpace(title), Price: price, Currency: currency,
			City: keywordDOMCity(container), ImageURL: firstImageURL(pageURL, container, title),
			AgeText: parseAgeText(htmlx.Text(container)), Position: len(cards),
		}
		seen[id] = true
		cards = append(cards, card)
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("parse keyword HTML: primary card title missing")
	}
	return cards, nil
}

// Copy only the tree preceding the documented foreign-region heading. Cutting
// membership before parsing also stops card-container text/image leakage.
func keywordPrimaryDOM(root *htmlx.Node) *htmlx.Node {
	stopped := false
	var copyNode func(*htmlx.Node, *htmlx.Node) *htmlx.Node
	copyNode = func(n, parent *htmlx.Node) *htmlx.Node {
		if stopped {
			return nil
		}
		if n.Tag == "#text" || len(n.Tag) == 2 && n.Tag[0] == 'h' && n.Tag[1] >= '1' && n.Tag[1] <= '6' {
			text := strings.ToLower(htmlx.Text(n))
			if strings.Contains(text, "объявлен") && (strings.Contains(text, "из других регионов") || strings.Contains(text, "в других городах")) {
				stopped = true
				return nil
			}
		}
		out := &htmlx.Node{Tag: n.Tag, TextData: n.TextData, Attr: n.Attr, Parent: parent}
		for _, child := range n.Children {
			if copied := copyNode(child, out); copied != nil {
				out.Children = append(out.Children, copied)
			}
		}
		return out
	}
	return copyNode(root, nil)
}

func keywordCardContainer(anchor *htmlx.Node) *htmlx.Node {
	best := anchor
	for n, depth := anchor.Parent, 0; n != nil && depth < 14; n, depth = n.Parent, depth+1 {
		if len(uniqueAdIDs(n)) != 1 {
			break
		}
		best = n
		if n.Tag == "article" || n.Tag == "li" || n.HasClassPart("card") || n.HasClassPart("advert") {
			return n
		}
	}
	return best
}

func keywordDOMCity(root *htmlx.Node) string {
	// Preserve only a city actually present as a visible field, never infer it
	// from the requested scope. Missing source values stay missing.
	for _, line := range strings.Split(htmlx.TextLines(root), "\n") {
		for _, city := range KeywordCities()[1:] {
			if strings.EqualFold(strings.TrimSpace(line), city.Label) {
				return city.Label
			}
		}
	}
	return ""
}

func keywordDOMMoney(root *htmlx.Node) (*int, string) {
	// A plain structured amount does not imply TJS. Require currency evidence.
	currency := htmlx.FirstMetaContent(root, map[string]string{"itemprop": "priceCurrency"})
	if currency == "" {
		currency = htmlx.FirstMetaContent(root, map[string]string{"property": "product:price:currency"})
	}
	if strings.EqualFold(currency, "TJS") {
		raw := firstNonEmpty(
			htmlx.FirstMetaContent(root, map[string]string{"itemprop": "price"}),
			htmlx.FirstMetaContent(root, map[string]string{"property": "product:price:amount"}),
		)
		clean := strings.NewReplacer(" ", "", "\u00a0", "", "\u202f", "").Replace(raw)
		if value, err := strconv.Atoi(clean); err == nil && value >= 0 {
			return &value, "TJS"
		}
	}
	if currency != "" && !strings.EqualFold(currency, "TJS") {
		return nil, currency
	}
	var dedicated *int
	htmlx.Walk(root, func(n *htmlx.Node) bool {
		if dedicated != nil {
			return false
		}
		if n.GetAttr("itemprop") == "price" || classHasMarker(n.GetAttr("class"), "price") || n.GetAttr("data-component") == "SidebarPrice" {
			raw := firstNonEmpty(n.GetAttr("content"), htmlx.Text(n))
			if match := priceRE.FindStringSubmatch(raw); priceOnlyRE.MatchString(raw) && len(match) == 2 {
				clean := strings.NewReplacer(" ", "", "\u00a0", "", "\u202f", "").Replace(match[1])
				if value, err := strconv.Atoi(clean); err == nil && value >= 0 {
					dedicated = &value
				}
			}
		}
		return true
	})
	if dedicated != nil {
		return dedicated, "TJS"
	}
	for _, line := range strings.Split(htmlx.TextLines(root), "\n") {
		line = htmlx.NormalizeSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "цена:") {
			line = strings.TrimSpace(line[len("цена:"):])
		}
		if priceOnlyRE.MatchString(line) {
			match := priceRE.FindStringSubmatch(line)
			clean := strings.NewReplacer(" ", "", "\u00a0", "", "\u202f", "").Replace(match[1])
			if value, err := strconv.Atoi(clean); err == nil && value >= 0 {
				return &value, "TJS"
			}
		}
	}
	return nil, currency
}
