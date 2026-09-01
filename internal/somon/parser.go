package somon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nicelight/somon-rent-watcher/internal/htmlx"
	"github.com/nicelight/somon-rent-watcher/internal/model"
)

var (
	adPathRE    = regexp.MustCompile(`(?i)/adv/(\d+)(?:_|/|\?|#|$)`)
	priceRE     = regexp.MustCompile(`(?i)([0-9][0-9 \x{00A0}\x{202F}]*)\s*(?:[cс]\.?|сомони|tjs)`)
	roomsRE     = regexp.MustCompile(`(?i)(\d+)\s*(?:[-–—]\s*)?комн(?:\.|ат)?`)
	floorRE     = regexp.MustCompile(`(?i)(\d+)\s*(?:-?й\s*)?этаж`)
	ageRE       = regexp.MustCompile(`(?i)(только что|\d+\s+(?:секунда|секунды|секунд|минута|минуты|минут|час|часа|часов|день|дня|дней|неделя|недели|недель)\s+назад|сегодня|вчера)`)
	activeAdsRE = regexp.MustCompile(`(?i)(\d{1,7})\s+активн(?:ое|ых|ые)\s+объявлен(?:ие|ия|ий)?`)
	idTextRE    = regexp.MustCompile(`(?i)ID\s*:\s*(\d+)`)
	priceOnlyRE = regexp.MustCompile(`(?i)^(?:[0-9][0-9 \x{00A0}\x{202F}]*\s*(?:[cс]\.?|сомони|tjs))\s*(?:[0-9][0-9 \x{00A0}\x{202F}]*\s*(?:[cс]\.?|сомони|tjs))?$`)
	numberRE    = regexp.MustCompile(`\d+`)

	// Next.js App Router embeds server data in React Server Component chunks.
	// Somon has used this representation; it is parsed only as a deterministic
	// structured fallback when/alongside ordinary server-rendered card markup.
	rscPushRE = regexp.MustCompile(`(?s)self\.__next_f\.push\(\[\s*\d+\s*,\s*("(?:\\.|[^"\\])*")\s*\]\s*\)`)
)

// ParseCategory reads both ordinary server-rendered cards and, when present,
// Somon's deterministic Next.js/RSC adverts payload. Visible DOM values win;
// RSC fills missing cards and fields. No JavaScript is executed.
func ParseCategory(pageURL string, body []byte) ([]model.Card, error) {
	root, err := htmlx.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parse category HTML: %w", err)
	}
	domCards := parseCategoryDOM(pageURL, root)
	rscCards := parseCategoryRSC(pageURL, body)
	cards := mergeCategorySources(domCards, rscCards)
	if len(cards) == 0 {
		// Somon templates can contain hidden anti-blocking modal text even on a
		// normal page. Treat it as a block only when no real cards were found.
		if reason := blockedReason(root); reason != "" {
			return nil, &BlockedPageError{URL: pageURL, Reason: reason}
		}
		// Fail closed instead of guessing from arbitrary /adv/ strings in scripts
		// or neighbouring cards. The application saves the raw HTML for inspection
		// and leaves seen state untouched.
		return nil, fmt.Errorf("parse category HTML: ad cards not found (DOM=0, RSC=0)")
	}
	return cards, nil
}

func parseCategoryDOM(pageURL string, root *htmlx.Node) []model.Card {
	anchors := categoryAdAnchors(root)

	byID := make(map[int64]model.Card)
	order := make([]int64, 0, len(anchors))
	for _, anchor := range anchors {
		id := adIDFromURL(anchor.GetAttr("href"))
		if id == 0 {
			continue
		}
		container := cardContainer(anchor)
		text := htmlx.Text(container)
		title := titleFromNode(anchor)
		if title == "" {
			title = extractTitle(text)
		}
		if title == "" {
			continue
		}

		card := model.Card{
			ID:       id,
			URL:      normalizeAdURL(pageURL, anchor.GetAttr("href")),
			Title:    title,
			Price:    parsePrice(text),
			Rooms:    parseRooms(title),
			Floor:    parseFloor(title),
			ImageURL: firstImageURL(pageURL, container, title),
			AgeText:  parseAgeText(text),
		}
		card.Promoted, card.PromotionLabel = parsePromotionDOM(container, text)

		if existing, ok := byID[id]; ok {
			byID[id] = mergeCardMissing(existing, card)
			continue
		}
		card.Position = len(order)
		byID[id] = card
		order = append(order, id)
	}

	cards := make([]model.Card, 0, len(order))
	for _, id := range order {
		card := byID[id]
		if card.URL == "" || card.Title == "" {
			continue
		}
		cards = append(cards, card)
	}
	return cards
}

// categoryAdAnchors returns advertisement links only from the requested city
// section. Somon can append recommendation cards after a heading such as
// "Объявления ... в других городах"; those cards must not enter the Dushanbe
// fresh feed.
func categoryAdAnchors(root *htmlx.Node) []*htmlx.Node {
	anchors := make([]*htmlx.Node, 0, 128)
	stopped := false
	var visit func(*htmlx.Node)
	visit = func(n *htmlx.Node) {
		if n == nil || stopped {
			return
		}
		if n.Tag == "#text" {
			text := strings.ToLower(htmlx.NormalizeSpace(n.TextData))
			if strings.Contains(text, "объявлен") && strings.Contains(text, "в других городах") {
				stopped = true
			}
			return
		}
		if n.Tag == "a" && adIDFromURL(n.GetAttr("href")) > 0 {
			anchors = append(anchors, n)
		}
		for _, child := range n.Children {
			visit(child)
			if stopped {
				return
			}
		}
	}
	visit(root)
	return anchors
}

func parseCategoryRSC(pageURL string, body []byte) []model.Card {
	chunks := decodeRSCChunks(body)
	adverts := findLargestAdvertsArray(chunks)
	if len(adverts) == 0 {
		return nil
	}

	seen := make(map[int64]struct{}, len(adverts))
	cards := make([]model.Card, 0, len(adverts))
	for _, raw := range adverts {
		advert, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		card, ok := cardFromRSC(pageURL, advert)
		if !ok {
			continue
		}
		if _, exists := seen[card.ID]; exists {
			continue
		}
		seen[card.ID] = struct{}{}
		card.Position = len(cards)
		cards = append(cards, card)
	}
	return cards
}

func cardFromRSC(pageURL string, advert map[string]any) (model.Card, bool) {
	id, ok := firstInt64(advert, "id", "pk", "advert_id", "advertId")
	if !ok || id <= 0 {
		return model.Card{}, false
	}
	rawURL := firstString(advert, "url", "absolute_url", "absoluteUrl", "link")
	if rawURL == "" {
		return model.Card{}, false
	}
	title := firstString(advert, "title", "name")
	if title == "" {
		return model.Card{}, false
	}
	card := model.Card{
		ID:       id,
		URL:      normalizeAdURL(pageURL, rawURL),
		Title:    htmlx.NormalizeSpace(title),
		Price:    firstIntPointer(advert, "price", "current_price", "currentPrice"),
		Rooms:    parseRooms(title),
		Floor:    parseFloor(title),
		ImageURL: normalizeResourceURL(pageURL, firstString(advert, "first_thumb", "firstThumb", "image", "image_url", "imageUrl")),
		AgeText:  firstString(advert, "published", "published_at_display", "publishedAtDisplay", "date_display", "dateDisplay"),
	}
	if card.Rooms == nil {
		card.Rooms = firstIntPointer(advert, "rooms", "room_count", "roomCount")
	}
	if card.Floor == nil {
		card.Floor = firstIntPointer(advert, "floor", "floor_number", "floorNumber")
	}
	card.Promoted, card.PromotionLabel = parsePromotionRSC(advert)
	return card, true
}

func mergeCategorySources(domCards, rscCards []model.Card) []model.Card {
	if len(domCards) == 0 {
		return reindexCards(rscCards)
	}
	if len(rscCards) == 0 {
		return reindexCards(domCards)
	}

	// The visible DOM defines the actual first-page city feed and its order.
	// RSC may contain unrelated recommendation arrays, so when DOM cards exist
	// it is used only to fill missing fields of matching IDs, never to introduce
	// an additional ID that was not visibly present before the other-city marker.
	structured := make(map[int64]model.Card, len(rscCards))
	for _, card := range rscCards {
		if card.ID > 0 {
			structured[card.ID] = card
		}
	}
	out := make([]model.Card, 0, len(domCards))
	seen := make(map[int64]struct{}, len(domCards))
	for _, visible := range domCards {
		if visible.ID <= 0 {
			continue
		}
		if _, ok := seen[visible.ID]; ok {
			continue
		}
		seen[visible.ID] = struct{}{}
		card := visible
		if rsc, ok := structured[visible.ID]; ok {
			card = overlayVisibleCard(mergeCardMissing(rsc, visible), visible)
		}
		if card.URL == "" || card.Title == "" {
			continue
		}
		card.Position = len(out)
		out = append(out, card)
	}
	return out
}

func reindexCards(cards []model.Card) []model.Card {
	out := make([]model.Card, 0, len(cards))
	seen := make(map[int64]struct{}, len(cards))
	for _, card := range cards {
		if card.ID <= 0 || card.URL == "" || card.Title == "" {
			continue
		}
		if _, ok := seen[card.ID]; ok {
			continue
		}
		seen[card.ID] = struct{}{}
		card.Position = len(out)
		out = append(out, card)
	}
	return out
}

func overlayVisibleCard(base, visible model.Card) model.Card {
	if visible.URL != "" {
		base.URL = visible.URL
	}
	if visible.Title != "" {
		base.Title = visible.Title
	}
	if visible.Price != nil {
		base.Price = visible.Price
	}
	if visible.Rooms != nil {
		base.Rooms = visible.Rooms
	}
	if visible.Floor != nil {
		base.Floor = visible.Floor
	}
	if visible.ImageURL != "" {
		base.ImageURL = visible.ImageURL
	}
	if visible.AgeText != "" {
		base.AgeText = visible.AgeText
	}
	if visible.Promoted {
		base.Promoted = true
		base.PromotionLabel = visible.PromotionLabel
	}
	return base
}

func mergeCardMissing(a, b model.Card) model.Card {
	if a.Title == "" || len([]rune(b.Title)) > len([]rune(a.Title)) {
		a.Title = b.Title
	}
	if a.Price == nil {
		a.Price = b.Price
	}
	if a.Rooms == nil {
		a.Rooms = b.Rooms
	}
	if a.Floor == nil {
		a.Floor = b.Floor
	}
	if a.ImageURL == "" {
		a.ImageURL = b.ImageURL
	}
	if a.AgeText == "" {
		a.AgeText = b.AgeText
	}
	if b.Promoted {
		a.Promoted = true
		a.PromotionLabel = b.PromotionLabel
	}
	if a.URL == "" {
		a.URL = b.URL
	}
	return a
}

func ParseDetail(pageURL string, body []byte, fallback model.Card) (model.Ad, error) {
	root, err := htmlx.Parse(body)
	if err != nil {
		return model.Ad{}, fmt.Errorf("parse detail HTML: %w", err)
	}
	ad := model.Ad{Card: fallback}
	ad.URL = pageURL
	if id := adIDFromURL(pageURL); id > 0 {
		ad.ID = id
	}

	if advert := findBestAdvertMap(decodeRSCChunks(body), ad.ID); advert != nil {
		mergeDetailRSC(&ad, pageURL, advert)
	}
	mergeVisibleDetail(&ad, pageURL, root)

	if ad.ID == 0 || ad.Title == "" {
		if reason := blockedReason(root); reason != "" {
			return model.Ad{}, &BlockedPageError{URL: pageURL, Reason: reason}
		}
		return model.Ad{}, fmt.Errorf("detail page lacks mandatory ID/title")
	}
	return ad, nil
}

func mergeDetailRSC(ad *model.Ad, pageURL string, advert map[string]any) {
	if id, ok := firstInt64(advert, "id", "pk", "advert_id", "advertId"); ok && id > 0 {
		ad.ID = id
	}
	if rawURL := firstString(advert, "url", "absolute_url", "absoluteUrl", "link"); rawURL != "" {
		ad.URL = normalizeAdURL(pageURL, rawURL)
	}
	if title := firstString(advert, "title", "name"); title != "" {
		ad.Title = htmlx.NormalizeSpace(title)
		ad.Rooms = parseRooms(ad.Title)
		ad.Floor = parseFloor(ad.Title)
	}
	if price := firstIntPointer(advert, "price", "current_price", "currentPrice"); price != nil {
		ad.Price = price
	}
	if description, exists := firstPresentString(advert, "description", "body", "text"); exists {
		ad.Description = strings.TrimSpace(description)
	}
	if image := firstString(advert, "first_thumb", "firstThumb", "image", "image_url", "imageUrl"); image != "" {
		ad.ImageURL = normalizeResourceURL(pageURL, image)
	}
	if age := firstString(advert, "published", "published_at_display", "publishedAtDisplay", "date_display", "dateDisplay"); age != "" {
		ad.AgeText = age
	}
	if rooms := firstIntPointer(advert, "rooms", "room_count", "roomCount"); ad.Rooms == nil && rooms != nil {
		ad.Rooms = rooms
	}
	if floor := firstIntPointer(advert, "floor", "floor_number", "floorNumber"); ad.Floor == nil && floor != nil {
		ad.Floor = floor
	}
	if promoted, label := parsePromotionRSC(advert); promoted {
		ad.Promoted = true
		ad.PromotionLabel = label
	}
	if user, ok := advert["user"].(map[string]any); ok {
		if name := firstString(user, "name", "username", "display_name", "displayName"); name != "" {
			ad.SellerName = name
		}
		if count, ok := findActiveAdsCount(user); ok {
			ad.SellerAds = &count
		}
	}
	if ad.SellerAds == nil {
		if count, ok := findActiveAdsCount(advert); ok {
			ad.SellerAds = &count
		}
	}
	if ad.Floor == nil {
		if value, ok := findLabeledInt(advert, []string{"этаж", "floor"}); ok {
			ad.Floor = &value
		}
	}
	if ad.Rooms == nil {
		if value, ok := findLabeledInt(advert, []string{"комнат", "rooms"}); ok {
			ad.Rooms = &value
		}
	}
}

func mergeVisibleDetail(ad *model.Ad, pageURL string, root *htmlx.Node) {
	linesText := htmlx.TextLines(root)
	if h1 := htmlx.FindFirst(root, func(n *htmlx.Node) bool { return n.Tag == "h1" }); h1 != nil {
		if title := htmlx.Text(h1); title != "" {
			ad.Title = title
		}
	}
	if ad.Title == "" {
		ad.Title = htmlx.FirstMetaContent(root, map[string]string{"property": "og:title", "name": "twitter:title"})
	}
	if idMatch := idTextRE.FindStringSubmatch(linesText); len(idMatch) == 2 {
		if id, err := strconv.ParseInt(idMatch[1], 10, 64); err == nil && id > 0 {
			ad.ID = id
		}
	}
	if rooms := parseRooms(ad.Title); rooms != nil {
		ad.Rooms = rooms
	}
	if floor := extractLabeledInt(linesText, "этаж"); floor != nil {
		ad.Floor = floor
	} else if floor := parseFloor(ad.Title); floor != nil {
		ad.Floor = floor
	}
	if price := detailPrice(root, linesText); price != nil {
		ad.Price = price
	}
	if description := extractDescriptionFromDOM(root, linesText); description != "" {
		ad.Description = description
	} else if ad.Description == "" {
		ad.Description = htmlx.FirstMetaContent(root, map[string]string{
			"property": "og:description",
			"name":     "description",
		})
	}
	if age := parseAgeText(linesText); age != "" {
		ad.AgeText = age
	}
	if image := detailImageURL(pageURL, root, ad.Title); image != "" {
		ad.ImageURL = image
	}

	sellerText := sellerBlockText(root, linesText)
	name, since, count := parseSeller(sellerText)
	if name != "" {
		ad.SellerName = name
	}
	if since != "" {
		ad.SellerSince = since
	}
	if count != nil {
		ad.SellerAds = count
	}
	if ad.SellerAds == nil {
		if m := activeAdsRE.FindStringSubmatch(linesText); len(m) == 2 {
			if n, err := strconv.Atoi(m[1]); err == nil {
				ad.SellerAds = &n
			}
		}
	}
}

func blockedReason(root *htmlx.Node) string {
	visible := strings.ToLower(htmlx.Text(root))
	for _, marker := range []string{
		"аккаунт был заблокирован",
		"доступ временно ограничен",
		"слишком много запросов",
		"access denied",
		"temporarily blocked",
		"just a moment",
	} {
		if strings.Contains(visible, marker) {
			return marker
		}
	}
	return ""
}

func cardContainer(anchor *htmlx.Node) *htmlx.Node {
	best := anchor
	anchorHasTitle := titleFromNode(anchor) != ""
	for n, depth := anchor, 0; n != nil && depth < 14; n, depth = n.Parent, depth+1 {
		ids := uniqueAdIDs(n)
		if len(ids) > 1 {
			break
		}
		if len(ids) == 1 {
			text := htmlx.Text(n)
			if parsePrice(text) != nil && (anchorHasTitle || strings.Contains(strings.ToLower(text), "квартир") || roomsRE.MatchString(text)) {
				best = n
				if n.Tag == "li" || n.Tag == "article" {
					return n
				}
				class := strings.ToLower(n.GetAttr("class"))
				if strings.Contains(class, "card") || strings.Contains(class, "advert") || strings.Contains(class, "listing") || strings.Contains(class, "item") {
					return n
				}
			}
		}
	}
	return best
}

func uniqueAdIDs(root *htmlx.Node) map[int64]struct{} {
	ids := make(map[int64]struct{})
	htmlx.Walk(root, func(n *htmlx.Node) bool {
		if len(ids) > 1 {
			return false
		}
		if n.Tag == "a" {
			if id := adIDFromURL(n.GetAttr("href")); id > 0 {
				ids[id] = struct{}{}
			}
		}
		return true
	})
	return ids
}

func titleFromNode(n *htmlx.Node) string {
	for _, candidate := range []string{htmlx.Text(n), n.GetAttr("title")} {
		if title := extractTitle(candidate); title != "" {
			return title
		}
	}
	if img := htmlx.FindFirst(n, func(cur *htmlx.Node) bool { return cur.Tag == "img" }); img != nil {
		if title := extractTitle(img.GetAttr("alt")); title != "" {
			return title
		}
	}
	return ""
}

func extractTitle(text string) string {
	text = htmlx.NormalizeSpace(text)
	loc := roomsRE.FindStringIndex(text)
	if loc == nil {
		return ""
	}
	start := loc[0]
	end := len(text)
	for _, marker := range []string{" назад", " Сегодня", " Вчера", " Душанбе", " VIP", " ТОП", " TOP"} {
		if i := strings.Index(text[start:], marker); i >= 0 && start+i < end {
			end = start + i
		}
	}
	title := strings.TrimSpace(text[start:end])
	if len([]rune(title)) > 240 {
		r := []rune(title)
		title = string(r[:240])
	}
	return title
}

func parsePrice(text string) *int {
	text = htmlx.NormalizeSpace(text)
	m := priceRE.FindStringSubmatch(text)
	if len(m) != 2 {
		return nil
	}
	clean := strings.NewReplacer(" ", "", "\u00a0", "", "\u202f", "").Replace(m[1])
	n, err := strconv.Atoi(clean)
	if err != nil || n <= 0 {
		return nil
	}
	return &n
}

func parseRooms(text string) *int {
	m := roomsRE.FindStringSubmatch(htmlx.NormalizeSpace(text))
	if len(m) != 2 {
		return nil
	}
	n, err := strconv.Atoi(m[1])
	if err != nil || n <= 0 || n > 99 {
		return nil
	}
	return &n
}

func parseFloor(text string) *int {
	m := floorRE.FindStringSubmatch(htmlx.NormalizeSpace(text))
	if len(m) != 2 {
		return nil
	}
	n, err := strconv.Atoi(m[1])
	if err != nil || n <= 0 || n > 200 {
		return nil
	}
	return &n
}

func parseAgeText(text string) string {
	m := ageRE.FindStringSubmatch(htmlx.NormalizeSpace(text))
	if len(m) != 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

func ParseAgeDuration(text string) (time.Duration, bool) {
	text = strings.ToLower(htmlx.NormalizeSpace(text))
	if strings.Contains(text, "только что") || text == "сегодня" {
		return 0, true
	}
	if text == "вчера" {
		return 24 * time.Hour, true
	}
	m := numberRE.FindString(text)
	if m == "" {
		return 0, false
	}
	n, err := strconv.Atoi(m)
	if err != nil {
		return 0, false
	}
	switch {
	case strings.Contains(text, "секунд"):
		return time.Duration(n) * time.Second, true
	case strings.Contains(text, "минут"):
		return time.Duration(n) * time.Minute, true
	case strings.Contains(text, "час"):
		return time.Duration(n) * time.Hour, true
	case strings.Contains(text, "день") || strings.Contains(text, "дня") || strings.Contains(text, "дней"):
		return time.Duration(n) * 24 * time.Hour, true
	case strings.Contains(text, "недел"):
		return time.Duration(n) * 7 * 24 * time.Hour, true
	default:
		return 0, false
	}
}

func parsePromotionDOM(n *htmlx.Node, text string) (bool, string) {
	upper := strings.ToUpper(text)
	if containsMarker(upper, "VIP") || treeHasClassMarker(n, "vip") {
		return true, "VIP"
	}
	if containsMarker(upper, "ТОП") || containsMarker(upper, "TOP") ||
		treeHasClassMarker(n, "top") || treeHasClassMarker(n, "premium") || treeHasClassMarker(n, "promoted") {
		return true, "ТОП"
	}
	return false, ""
}

func parsePromotionRSC(advert map[string]any) (bool, string) {
	for _, key := range []string{"is_vip", "isVip", "vip"} {
		if asBool(advert[key]) {
			return true, "VIP"
		}
	}
	for _, key := range []string{"is_top", "isTop", "top", "is_promoted", "isPromoted", "promoted"} {
		if asBool(advert[key]) {
			return true, "ТОП"
		}
	}
	raw := advert["ad_type"]
	if raw == nil {
		raw = advert["adType"]
	}
	value := ""
	switch typed := raw.(type) {
	case map[string]any:
		value = firstString(typed, "type", "name", "code", "slug")
	case string:
		value = typed
	}
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || value == "ordinary" || value == "normal" || value == "regular" || value == "standard" || value == "default" || value == "free" {
		return false, ""
	}
	if strings.Contains(value, "vip") {
		return true, "VIP"
	}
	return true, "ТОП"
}

func containsMarker(text, marker string) bool {
	for _, token := range strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if token == marker {
			return true
		}
	}
	return false
}

func classHasMarker(className, marker string) bool {
	className = strings.ToLower(className)
	marker = strings.ToLower(marker)
	for _, token := range strings.FieldsFunc(className, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if token == marker {
			return true
		}
	}
	return false
}

func treeHasClassMarker(root *htmlx.Node, marker string) bool {
	found := false
	htmlx.Walk(root, func(n *htmlx.Node) bool {
		if classHasMarker(n.GetAttr("class"), marker) {
			found = true
			return false
		}
		return true
	})
	return found
}

func firstImageURL(pageURL string, root *htmlx.Node, title string) string {
	images := htmlx.FindAll(root, func(n *htmlx.Node) bool { return n.Tag == "img" || n.Tag == "source" })
	var fallback string
	for _, img := range images {
		ref := firstNonEmpty(img.GetAttr("src"), img.GetAttr("data-src"), img.GetAttr("data-original"), firstSrcset(img.GetAttr("srcset")), firstSrcset(img.GetAttr("data-srcset")))
		if ref == "" || strings.HasPrefix(ref, "data:") {
			continue
		}
		resolved := normalizeResourceURL(pageURL, ref)
		if fallback == "" {
			fallback = resolved
		}
		alt := strings.ToLower(img.GetAttr("alt"))
		if title != "" && alt != "" && strings.Contains(strings.ToLower(title), strings.TrimSpace(strings.Split(alt, ", Душанбе")[0])) {
			return resolved
		}
		if strings.Contains(strings.ToLower(resolved), "cdntj.somon.tj") {
			return resolved
		}
	}
	return fallback
}

func detailImageURL(pageURL string, root *htmlx.Node, title string) string {
	if image := htmlx.FirstMetaContent(root, map[string]string{"property": "og:image", "name": "twitter:image"}); image != "" {
		return normalizeResourceURL(pageURL, image)
	}
	return firstImageURL(pageURL, root, title)
}

func firstSrcset(srcset string) string {
	if srcset == "" {
		return ""
	}
	first := strings.Split(srcset, ",")[0]
	fields := strings.Fields(first)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func adIDFromURL(raw string) int64 {
	m := adPathRE.FindStringSubmatch(raw)
	if len(m) != 2 {
		return 0
	}
	id, _ := strconv.ParseInt(m[1], 10, 64)
	return id
}

func normalizeAdURL(baseURL, raw string) string {
	resolved := htmlx.ResolveURL(baseURL, strings.ReplaceAll(raw, `\/`, `/`))
	u, err := url.Parse(resolved)
	if err != nil {
		return resolved
	}
	u.Fragment = ""
	return u.String()
}

func normalizeResourceURL(baseURL, raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, `\/`, `/`))
	if raw == "" {
		return ""
	}
	return htmlx.ResolveURL(baseURL, raw)
}

func extractLabeledInt(linesText, label string) *int {
	lines := strings.Split(linesText, "\n")
	label = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(label), ":"))
	for i, line := range lines {
		low := strings.ToLower(strings.TrimSpace(line))
		if low == label || low == label+":" {
			if i+1 < len(lines) {
				if n := firstInt(lines[i+1]); n != nil {
					return n
				}
			}
		}
		if strings.HasPrefix(low, label+":") || strings.HasPrefix(low, label+" ") {
			if n := firstInt(low[len(label):]); n != nil {
				return n
			}
		}
	}
	return nil
}

func firstInt(text string) *int {
	m := numberRE.FindString(text)
	if m == "" {
		return nil
	}
	n, err := strconv.Atoi(m)
	if err != nil {
		return nil
	}
	return &n
}

func detailPrice(root *htmlx.Node, linesText string) *int {
	for _, keys := range []map[string]string{
		{"property": "product:price:amount"},
		{"itemprop": "price"},
	} {
		if raw := htmlx.FirstMetaContent(root, keys); raw != "" {
			if n := firstInt(strings.NewReplacer(" ", "", "\u00a0", "", "\u202f", "").Replace(raw)); n != nil && *n > 0 {
				return n
			}
		}
	}
	var found *int
	htmlx.Walk(root, func(n *htmlx.Node) bool {
		if found != nil {
			return false
		}
		if n.GetAttr("itemprop") == "price" || classHasMarker(n.GetAttr("class"), "price") {
			if raw := firstNonEmpty(n.GetAttr("content"), htmlx.Text(n)); raw != "" {
				found = parsePrice(raw)
			}
		}
		return true
	})
	if found != nil {
		return found
	}
	return parsePrice(linesText)
}

func extractDescriptionFromDOM(root *htmlx.Node, fallbackLines string) string {
	heading := htmlx.FindFirst(root, func(n *htmlx.Node) bool {
		if n.Tag == "#text" {
			return false
		}
		return strings.EqualFold(htmlx.NormalizeSpace(htmlx.Text(n)), "Описание")
	})
	if heading != nil {
		for parent, depth := heading.Parent, 0; parent != nil && depth < 3; parent, depth = parent.Parent, depth+1 {
			text := htmlx.TextLines(parent)
			if description := extractDescription(text); description != "" && len([]rune(description)) <= 5000 {
				return description
			}
		}
	}
	return extractDescription(fallbackLines)
}

func extractDescription(linesText string) string {
	lines := strings.Split(linesText, "\n")
	start := -1
	for i, line := range lines {
		if strings.EqualFold(strings.TrimSpace(line), "Описание") {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	var out []string
	for i := start; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		low := strings.ToLower(line)
		nextLow := ""
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) != "" {
				nextLow = strings.ToLower(strings.TrimSpace(lines[j]))
				break
			}
		}
		if strings.Contains(low, "показать телефон") || strings.Contains(low, "начать чат") || strings.Contains(low, "whatsapp") || strings.Contains(low, "на сайте с") || strings.Contains(low, "активн") || strings.Contains(low, "пожаловаться") || strings.Contains(low, "похожие объявления") {
			break
		}
		if strings.Contains(nextLow, "на сайте с") || strings.Contains(nextLow, "активн") {
			break
		}
		if priceOnlyRE.MatchString(htmlx.NormalizeSpace(line)) {
			break
		}
		out = append(out, line)
		if len([]rune(strings.Join(out, " "))) > 4000 {
			break
		}
	}
	return strings.TrimSpace(strings.Join(out, " "))
}

func sellerBlockText(root *htmlx.Node, linesText string) string {
	var best string
	htmlx.Walk(root, func(n *htmlx.Node) bool {
		if n.Tag == "#text" {
			return true
		}
		text := htmlx.Text(n)
		low := strings.ToLower(text)
		if !strings.Contains(low, "на сайте с") || !strings.Contains(low, "активн") {
			return true
		}
		length := len([]rune(text))
		if length == 0 || length > 500 {
			return true
		}
		if best == "" || length < len([]rune(best)) {
			best = text
		}
		return true
	})
	if best != "" {
		return best
	}
	lines := strings.Split(linesText, "\n")
	for i, line := range lines {
		low := strings.ToLower(line)
		if strings.Contains(low, "на сайте с") {
			parts := []string{line}
			if i > 0 {
				parts = append([]string{lines[i-1]}, parts...)
			}
			if i+1 < len(lines) {
				parts = append(parts, lines[i+1])
			}
			joined := strings.Join(parts, " ")
			if strings.Contains(strings.ToLower(joined), "активн") {
				return joined
			}
		}
	}
	return ""
}

func parseSeller(text string) (name, since string, ads *int) {
	text = htmlx.NormalizeSpace(text)
	if text == "" {
		return "", "", nil
	}
	lower := strings.ToLower(text)
	idx := strings.Index(lower, "на сайте с")
	if idx >= 0 {
		name = strings.TrimSpace(text[:idx])
		rest := strings.TrimSpace(text[idx+len("на сайте с"):])
		if m := activeAdsRE.FindStringIndex(rest); m != nil {
			since = strings.TrimSpace(rest[:m[0]])
		}
	}
	if m := activeAdsRE.FindStringSubmatch(text); len(m) == 2 {
		if n, err := strconv.Atoi(m[1]); err == nil {
			ads = &n
		}
	}
	return cleanSellerName(name), since, ads
}

func cleanSellerName(name string) string {
	name = strings.TrimSpace(name)
	if len([]rune(name)) > 120 {
		r := []rune(name)
		name = string(r[len(r)-120:])
	}
	return name
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// RSC helpers.
func decodeRSCChunks(page []byte) []any {
	matches := rscPushRE.FindAllSubmatch(page, -1)
	out := make([]any, 0, len(matches))
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		var chunk string
		if err := json.Unmarshal(match[1], &chunk); err != nil {
			continue
		}
		for _, line := range strings.Split(chunk, "\n") {
			colon := strings.IndexByte(line, ':')
			if colon < 0 || colon+1 >= len(line) {
				continue
			}
			payload := strings.TrimSpace(line[colon+1:])
			if payload == "" {
				continue
			}
			decoder := json.NewDecoder(bytes.NewBufferString(payload))
			decoder.UseNumber()
			var value any
			if err := decoder.Decode(&value); err == nil {
				out = append(out, value)
			}
		}
	}
	return out
}

func walkJSON(value any, visit func(any)) {
	visit(value)
	switch typed := value.(type) {
	case map[string]any:
		for _, child := range typed {
			walkJSON(child, visit)
		}
	case []any:
		for _, child := range typed {
			walkJSON(child, visit)
		}
	}
}

func findLargestAdvertsArray(chunks []any) []any {
	var best []any
	for _, chunk := range chunks {
		walkJSON(chunk, func(value any) {
			object, ok := value.(map[string]any)
			if !ok {
				return
			}
			for key, raw := range object {
				normalized := normalizeKey(key)
				if normalized != "adverts" && normalized != "advertisements" && normalized != "listings" {
					continue
				}
				candidate, ok := raw.([]any)
				if ok && len(candidate) > len(best) {
					best = candidate
				}
			}
		})
	}
	return best
}

func findBestAdvertMap(chunks []any, expectedID int64) map[string]any {
	var best map[string]any
	bestScore := 0
	for _, chunk := range chunks {
		walkJSON(chunk, func(value any) {
			object, ok := value.(map[string]any)
			if !ok {
				return
			}
			score := 0
			if id, ok := firstInt64(object, "id", "pk", "advert_id", "advertId"); ok && id == expectedID {
				score += 100
			}
			if strings.Contains(firstString(object, "url", "link"), "/adv/") {
				score += 15
			}
			if firstString(object, "title", "name") != "" {
				score += 10
			}
			if _, exists := object["description"]; exists {
				score += 25
			}
			if _, exists := object["price"]; exists {
				score += 5
			}
			if _, exists := object["user"]; exists {
				score += 5
			}
			if score > bestScore {
				bestScore = score
				best = object
			}
		})
	}
	if bestScore < 100 {
		return nil
	}
	return best
}

func findActiveAdsCount(value any) (int, bool) {
	known := map[string]struct{}{
		"activeadscount": {}, "activeadvertcount": {}, "activeadvertscount": {},
		"activeitemscount": {}, "activelistingscount": {}, "advertscount": {},
		"advertcount": {}, "itemscount": {}, "listingscount": {},
	}
	var result int
	found := false
	walkJSON(value, func(current any) {
		if found {
			return
		}
		object, ok := current.(map[string]any)
		if !ok {
			return
		}
		for key, raw := range object {
			normalized := normalizeKey(key)
			_, exact := known[normalized]
			generic := strings.Contains(normalized, "active") && strings.Contains(normalized, "count") &&
				(strings.Contains(normalized, "ad") || strings.Contains(normalized, "advert") || strings.Contains(normalized, "item") || strings.Contains(normalized, "listing"))
			if !exact && !generic {
				continue
			}
			if value, ok := asInt(raw); ok && value >= 0 {
				result = value
				found = true
				return
			}
		}
	})
	return result, found
}

func findLabeledInt(value any, labels []string) (int, bool) {
	var result int
	found := false
	walkJSON(value, func(current any) {
		if found {
			return
		}
		object, ok := current.(map[string]any)
		if !ok {
			return
		}
		label := strings.ToLower(firstString(object, "label", "name", "title", "key"))
		matched := false
		for _, candidate := range labels {
			if strings.Contains(label, strings.ToLower(candidate)) {
				matched = true
				break
			}
		}
		if !matched {
			return
		}
		for _, key := range []string{"value", "display_value", "displayValue", "text"} {
			if number, ok := asInt(object[key]); ok {
				result = number
				found = true
				return
			}
		}
	})
	return result, found
}

func firstString(object map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := asString(object[key]); value != "" {
			return value
		}
	}
	return ""
}

func firstPresentString(object map[string]any, keys ...string) (string, bool) {
	for _, key := range keys {
		value, exists := object[key]
		if !exists {
			continue
		}
		return asString(value), true
	}
	return "", false
}

func firstInt64(object map[string]any, keys ...string) (int64, bool) {
	for _, key := range keys {
		if value, ok := asInt64(object[key]); ok {
			return value, true
		}
	}
	return 0, false
}

func firstIntPointer(object map[string]any, keys ...string) *int {
	for _, key := range keys {
		if value, ok := asInt(object[key]); ok && value >= 0 {
			return &value
		}
	}
	return nil
}

func asString(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return typed.String()
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	default:
		return ""
	}
}

func asInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case json.Number:
		if result, err := typed.Int64(); err == nil {
			return result, true
		}
		if floating, err := typed.Float64(); err == nil {
			return int64(floating), true
		}
	case float64:
		return int64(typed), true
	case float32:
		return int64(typed), true
	case int:
		return int64(typed), true
	case int64:
		return typed, true
	case string:
		cleaned := digitsAndMinus(typed)
		if cleaned == "" || cleaned == "-" {
			return 0, false
		}
		if result, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
			return result, true
		}
	}
	return 0, false
}

func asInt(value any) (int, bool) {
	result, ok := asInt64(value)
	if !ok || result > int64(^uint(0)>>1) || result < -int64(^uint(0)>>1)-1 {
		return 0, false
	}
	return int(result), true
}

func asBool(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		parsed, _ := strconv.ParseBool(strings.TrimSpace(typed))
		return parsed
	case json.Number:
		return typed.String() != "0"
	case float64:
		return typed != 0
	case int:
		return typed != 0
	default:
		return false
	}
}

func digitsAndMinus(value string) string {
	var b strings.Builder
	for i, r := range strings.TrimSpace(value) {
		if unicode.IsDigit(r) || (r == '-' && i == 0) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeKey(key string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(key) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeText is exported for Telegram snippets and parser tests.
func NormalizeText(text string) string {
	var b strings.Builder
	space := false
	for _, r := range strings.TrimSpace(text) {
		if unicode.IsSpace(r) || r == '\u00a0' || r == '\u202f' {
			if !space {
				b.WriteByte(' ')
				space = true
			}
			continue
		}
		space = false
		b.WriteRune(r)
	}
	return b.String()
}

func sortedCardIDs(cards []model.Card) []int64 {
	ids := make([]int64, 0, len(cards))
	for _, card := range cards {
		ids = append(ids, card.ID)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}
