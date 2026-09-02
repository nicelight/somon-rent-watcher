package telegram

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
)

func SettingsText(s filter.Settings) string {
	monitoring := "⏸ на паузе"
	if s.Enabled {
		monitoring = "▶ включён"
	}
	text := "<b>Фильтр квартир</b>\n\n" +
		"Мониторинг: <b>" + monitoring + "</b>\n" +
		"Цена: <b>" + html.EscapeString(priceSummary(s)) + "</b>\n" +
		"Комнаты: <b>" + html.EscapeString(choicesSummary(s.Rooms, filter.RoomChoices, roomLabel)) + "</b>\n" +
		"Этаж: <b>" + html.EscapeString(choicesSummary(s.Floors, filter.FloorChoices, floorLabel)) + "</b>\n" +
		"Автор: <b>" + html.EscapeString(authorSummary(s.SellerAdsLimit)) + "</b>\n" +
		"Минус-слова: <b>" + html.EscapeString(negativeSummary(s.NegativeWords)) + "</b>\n" +
		"Тип: <b>" + html.EscapeString(typeSummary(s)) + "</b>\n" +
		"Интервал: <b>" + html.EscapeString(pollIntervalSummary(s)) + "</b>\n\n" +
		"Изменения применяются только к новым объявлениям."
	if !s.Enabled {
		text += "\nПока мониторинг на паузе, новые ID запоминаются без отправки."
	}
	return text
}

func StatusText(s model.RuntimeStatus) string {
	last := "ещё не было"
	if !s.LastSuccessfulPoll.IsZero() {
		last = formatTime(s.LastSuccessfulPoll)
	}
	next := "не назначен"
	if !s.NextPoll.IsZero() {
		next = formatTime(s.NextPoll)
	}
	text := "<b>Статус Somon Watcher</b>\n\n" +
		"Режим: <b>" + html.EscapeString(nonEmpty(s.Mode, "запуск")) + "</b>\n" +
		"Последний успешный опрос: <b>" + html.EscapeString(last) + "</b>\n" +
		"Следующий опрос: <b>" + html.EscapeString(next) + "</b>\n" +
		"Последний цикл: <b>" + strconv.Itoa(s.LastCardCount) + " карточек · " + strconv.Itoa(s.LastNewCount) + " новых · " + strconv.Itoa(s.LastSentCount) + " отправлено</b>\n" +
		"Известных ID: <b>" + strconv.FormatInt(s.SeenCount, 10) + "</b>"
	if !s.BackoffUntil.IsZero() {
		text += "\nBackoff до: <b>" + html.EscapeString(formatTime(s.BackoffUntil)) + "</b>"
	}
	if s.LastError != "" {
		text += "\nПоследняя ошибка: <code>" + html.EscapeString(truncate(s.LastError, 500)) + "</code>"
	}
	return text
}

func AdCaption(ad model.Ad) string {
	title := truncate(somon.NormalizeText(ad.Title), 180)
	if title == "" {
		title = "Квартира"
	}
	price := "цена не указана"
	if ad.Price != nil {
		price = formatNumber(*ad.Price) + " c."
	}

	var lines []string
	lines = append(lines, "<b>"+html.EscapeString(title)+" — "+html.EscapeString(price)+"</b>")

	var features []string
	if ad.Floor != nil {
		features = append(features, fmt.Sprintf("%d этаж", *ad.Floor))
	} else {
		features = append(features, "этаж не указан")
	}
	if ad.Promoted {
		features = append(features, "VIP")
	} else {
		features = append(features, "обычное")
	}
	lines = append(lines, html.EscapeString(strings.Join(features, " · ")))

	if ad.SellerName != "" {
		lines = append(lines, "", "Автор: <b>"+html.EscapeString(truncate(ad.SellerName, 100))+"</b>")
	}
	if ad.SellerAds != nil {
		lines = append(lines, "Активных объявлений: <b>"+strconv.Itoa(*ad.SellerAds)+"</b>")
	} else {
		lines = append(lines, "Активных объявлений: <b>не указано</b>")
	}

	if description := formatAdDescription(ad.Description); description != "" {
		lines = append(lines, "", html.EscapeString(description))
	}
	if ad.AgeText != "" {
		lines = append(lines, "", "Опубликовано на Somon: "+html.EscapeString(ad.AgeText))
	}
	lines = append(lines, "ID: <code>"+strconv.FormatInt(ad.ID, 10)+"</code>")

	// Every user-controlled field is truncated before HTML escaping. The visible
	// caption remains comfortably below Telegram's 1024-character photo limit,
	// and markup is never cut in the middle of a tag or entity.
	return strings.Join(lines, "\n")
}

func formatAdDescription(raw string) string {
	normalized := truncate(somon.NormalizeText(raw), 360)
	if normalized == "" {
		return ""
	}

	out := make([]rune, 0, len([]rune(normalized)))
	var previous rune
	for _, r := range normalized {
		if isEmojiFieldStart(r) && len(out) > 0 && !isEmojiContinuation(r, previous) {
			for len(out) > 0 && out[len(out)-1] == ' ' {
				out = out[:len(out)-1]
			}
			if len(out) > 0 && out[len(out)-1] != '\n' {
				out = append(out, '\n')
			}
		}
		out = append(out, r)
		previous = r
	}
	return strings.TrimSpace(string(out))
}

func isEmojiFieldStart(r rune) bool {
	return r >= 0x2600 && r <= 0x26ff ||
		r >= 0x2700 && r <= 0x27bf ||
		r >= 0x2b00 && r <= 0x2bff ||
		r >= 0x1f000 && r <= 0x1faff
}

func isEmojiContinuation(current, previous rune) bool {
	if current >= 0x1f3fb && current <= 0x1f3ff {
		return true // skin-tone modifier
	}
	if previous == '\u200d' {
		return true // joined emoji sequence
	}
	return current >= 0x1f1e6 && current <= 0x1f1ff && previous >= 0x1f1e6 && previous <= 0x1f1ff
}

func mainKeyboard(s filter.Settings) *InlineKeyboardMarkup {
	return mainKeyboardWithScanState(s, false)
}

func mainKeyboardWithScanState(s filter.Settings, scanning bool) *InlineKeyboardMarkup {
	toggleLabel := "▶ Включить мониторинг"
	if s.Enabled {
		toggleLabel = "⏸ Поставить на паузу"
	}
	scanLabel := "Сканировать сейчас"
	scanData := "e:scan"
	if scanning {
		scanLabel = "⏳ Сканирую..."
		scanData = "e:scan_busy"
	}
	return &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{
		{{Text: toggleLabel, CallbackData: "e:toggle"}},
		{{Text: "Интервал", CallbackData: "i:interval"}, {Text: scanLabel, CallbackData: scanData}},
		{{Text: "Цена", CallbackData: "i:price"}, {Text: "Комнаты", CallbackData: "m:rooms"}},
		{{Text: "Этаж", CallbackData: "m:floors"}, {Text: "Автор", CallbackData: "m:author"}},
		{{Text: "Минус-слова", CallbackData: "i:negative"}, {Text: "Тип", CallbackData: "m:type"}},
	}}
}

func roomsKeyboard(s filter.Settings) *InlineKeyboardMarkup {
	rows := make([][]InlineKeyboardButton, 0, 4)
	var row []InlineKeyboardButton
	for _, value := range filter.RoomChoices {
		label := roomLabel(value)
		if s.HasRoom(value) {
			label = "✓ " + label
		}
		row = append(row, InlineKeyboardButton{Text: label, CallbackData: "r:" + value})
		if len(row) == 3 {
			rows = append(rows, row)
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}
	rows = append(rows, []InlineKeyboardButton{{Text: "← Назад", CallbackData: "m:main"}})
	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

func floorsKeyboard(s filter.Settings) *InlineKeyboardMarkup {
	rows := make([][]InlineKeyboardButton, 0, 5)
	var row []InlineKeyboardButton
	for _, value := range filter.FloorChoices {
		label := floorLabel(value)
		if s.HasFloor(value) {
			label = "✓ " + label
		}
		row = append(row, InlineKeyboardButton{Text: label, CallbackData: "f:" + value})
		if len(row) == 4 {
			rows = append(rows, row)
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}
	rows = append(rows, []InlineKeyboardButton{{Text: "← Назад", CallbackData: "m:main"}})
	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

func authorKeyboard(s filter.Settings) *InlineKeyboardMarkup {
	choices := []struct {
		value int
		label string
	}{{0, "Не важно"}, {2, "< 2"}, {5, "< 5"}, {10, "< 10"}}
	rows := make([][]InlineKeyboardButton, 0, 3)
	for i := 0; i < len(choices); i += 2 {
		row := make([]InlineKeyboardButton, 0, 2)
		for j := i; j < i+2 && j < len(choices); j++ {
			label := choices[j].label
			if s.SellerAdsLimit == choices[j].value {
				label = "✓ " + label
			}
			row = append(row, InlineKeyboardButton{Text: label, CallbackData: "a:" + strconv.Itoa(choices[j].value)})
		}
		rows = append(rows, row)
	}
	rows = append(rows, []InlineKeyboardButton{{Text: "← Назад", CallbackData: "m:main"}})
	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

func typeKeyboard(s filter.Settings) *InlineKeyboardMarkup {
	ordinary := "Обычные"
	if s.IncludeOrdinary {
		ordinary = "✓ " + ordinary
	}
	promoted := "VIP"
	if s.IncludePromoted {
		promoted = "✓ " + promoted
	}
	return &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{
		{{Text: ordinary, CallbackData: "t:o"}, {Text: promoted, CallbackData: "t:p"}},
		{{Text: "← Назад", CallbackData: "m:main"}},
	}}
}

func priceSummary(s filter.Settings) string {
	switch {
	case s.PriceMin == nil && s.PriceMax == nil:
		return "не ограничена"
	case s.PriceMin == nil:
		return "до " + formatNumber(*s.PriceMax)
	case s.PriceMax == nil:
		return "от " + formatNumber(*s.PriceMin)
	default:
		return formatNumber(*s.PriceMin) + " — " + formatNumber(*s.PriceMax)
	}
}

func choicesSummary(values, all []string, label func(string) string) string {
	if len(values) == len(all) {
		return "все"
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, label(value))
	}
	return strings.Join(out, ", ")
}

func roomLabel(value string) string {
	switch value {
	case filter.UnknownChoice:
		return "не указано"
	case "6+":
		return "6+к"
	default:
		return value + "к"
	}
}

func floorLabel(value string) string {
	switch value {
	case filter.UnknownChoice:
		return "не указано"
	case "14+":
		return "14+эт"
	default:
		return value + "эт"
	}
}

func authorSummary(limit int) string {
	if limit == 0 {
		return "не важно"
	}
	return fmt.Sprintf("меньше %d", limit)
}

func negativeSummary(words []string) string {
	if len(words) == 0 {
		return "нет"
	}
	return truncate(strings.Join(words, ", "), 180)
}

func typeSummary(s filter.Settings) string {
	switch {
	case s.IncludeOrdinary && s.IncludePromoted:
		return "Обычные + VIP"
	case s.IncludeOrdinary:
		return "Обычные"
	default:
		return "VIP"
	}
}

func pollIntervalSummary(s filter.Settings) string {
	if s.PollMinMinutes == s.PollMaxMinutes {
		return strconv.Itoa(s.PollMinMinutes) + " мин"
	}
	return strconv.Itoa(s.PollMinMinutes) + "–" + strconv.Itoa(s.PollMaxMinutes) + " мин"
}

func formatNumber(n int) string {
	negative := n < 0
	if negative {
		n = -n
	}
	digits := strconv.Itoa(n)
	for i := len(digits) - 3; i > 0; i -= 3 {
		digits = digits[:i] + " " + digits[i:]
	}
	if negative {
		return "-" + digits
	}
	return digits
}

func formatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05 MST")
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 || utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return strings.TrimSpace(string(runes[:max]))
}

func nonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
