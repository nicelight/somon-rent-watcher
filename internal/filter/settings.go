package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

const (
	UnknownChoice         = "unknown"
	DefaultPollMinMinutes = 10
	DefaultPollMaxMinutes = 30
	MinimumPollMinutes    = 1
	MaximumPollMinutes    = 24 * 60
)

var RoomChoices = []string{"1", "2", "3", "4", "5", "6+", UnknownChoice}

var FloorChoices = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14+", UnknownChoice,
}

type Settings struct {
	Enabled         bool     `json:"enabled"`
	PriceMin        *int     `json:"price_min,omitempty"`
	PriceMax        *int     `json:"price_max,omitempty"`
	Rooms           []string `json:"rooms"`
	Floors          []string `json:"floors"`
	SellerAdsLimit  int      `json:"seller_ads_limit"`
	NegativeWords   []string `json:"negative_words"`
	IncludeOrdinary bool     `json:"include_ordinary"`
	IncludePromoted bool     `json:"include_promoted"`
	PollMinMinutes  int      `json:"poll_min_minutes"`
	PollMaxMinutes  int      `json:"poll_max_minutes"`
}

func DefaultSettings() Settings {
	return Settings{
		Enabled:         false,
		Rooms:           append([]string(nil), RoomChoices...),
		Floors:          append([]string(nil), FloorChoices...),
		SellerAdsLimit:  0,
		IncludeOrdinary: true,
		IncludePromoted: true,
		PollMinMinutes:  DefaultPollMinMinutes,
		PollMaxMinutes:  DefaultPollMaxMinutes,
	}
}

func DecodeSettings(data string) (Settings, error) {
	return DecodeSettingsWithPollRange(data, DefaultPollMinMinutes, DefaultPollMaxMinutes)
}

func DecodeSettingsWithPollRange(data string, pollMinMinutes, pollMaxMinutes int) (Settings, error) {
	if strings.TrimSpace(data) == "" {
		s := DefaultSettings()
		s.PollMinMinutes = pollMinMinutes
		s.PollMaxMinutes = pollMaxMinutes
		if err := s.NormalizeAndValidate(); err != nil {
			return Settings{}, err
		}
		return s, nil
	}
	var s Settings
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return Settings{}, fmt.Errorf("decode settings: %w", err)
	}
	// Older persisted settings predate the runtime poll range. Seed only when
	// both fields are absent so a partially invalid range is not hidden.
	if s.PollMinMinutes == 0 && s.PollMaxMinutes == 0 {
		s.PollMinMinutes = pollMinMinutes
		s.PollMaxMinutes = pollMaxMinutes
	}
	if err := s.NormalizeAndValidate(); err != nil {
		return Settings{}, err
	}
	return s, nil
}

func (s Settings) Encode() (string, error) {
	copy := s
	if err := copy.NormalizeAndValidate(); err != nil {
		return "", err
	}
	data, err := json.Marshal(copy)
	if err != nil {
		return "", fmt.Errorf("encode settings: %w", err)
	}
	return string(data), nil
}

func (s *Settings) NormalizeAndValidate() error {
	if s.PriceMin != nil && *s.PriceMin < 0 {
		return errors.New("минимальная цена не может быть отрицательной")
	}
	if s.PriceMax != nil && *s.PriceMax < 0 {
		return errors.New("максимальная цена не может быть отрицательной")
	}
	if s.PriceMin != nil && s.PriceMax != nil && *s.PriceMin > *s.PriceMax {
		return errors.New("минимальная цена больше максимальной")
	}

	var err error
	s.Rooms, err = normalizeChoices(s.Rooms, RoomChoices)
	if err != nil {
		return fmt.Errorf("rooms: %w", err)
	}
	if len(s.Rooms) == 0 {
		return errors.New("должен быть выбран хотя бы один вариант комнат")
	}
	s.Floors, err = normalizeChoices(s.Floors, FloorChoices)
	if err != nil {
		return fmt.Errorf("floors: %w", err)
	}
	if len(s.Floors) == 0 {
		return errors.New("должен быть выбран хотя бы один этаж")
	}

	switch s.SellerAdsLimit {
	case 0, 2, 5, 10:
	default:
		return errors.New("лимит объявлений автора должен быть 0, 2, 5 или 10")
	}
	if !s.IncludeOrdinary && !s.IncludePromoted {
		return errors.New("должен быть выбран хотя бы один тип объявления")
	}
	if s.PollMinMinutes < MinimumPollMinutes || s.PollMaxMinutes < MinimumPollMinutes {
		return fmt.Errorf("интервал опроса должен быть не меньше %d минуты", MinimumPollMinutes)
	}
	if s.PollMinMinutes > s.PollMaxMinutes {
		return errors.New("минимальный интервал опроса больше максимального")
	}
	if s.PollMaxMinutes > MaximumPollMinutes {
		return fmt.Errorf("интервал опроса должен быть не больше %d минут", MaximumPollMinutes)
	}

	seen := make(map[string]struct{})
	clean := make([]string, 0, len(s.NegativeWords))
	for _, word := range s.NegativeWords {
		word = normalizeText(word)
		if word == "" {
			continue
		}
		if len([]rune(word)) > 80 {
			return fmt.Errorf("слишком длинное негативное выражение: %q", word)
		}
		if _, ok := seen[word]; ok {
			continue
		}
		seen[word] = struct{}{}
		clean = append(clean, word)
		if len(clean) > 30 {
			return errors.New("не более 30 негативных слов/фраз")
		}
	}
	s.NegativeWords = clean
	return nil
}

func normalizeChoices(values, allowed []string) ([]string, error) {
	allowedSet := make(map[string]int, len(allowed))
	for i, v := range allowed {
		allowedSet[v] = i
	}
	seen := make(map[string]struct{})
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if _, ok := allowedSet[v]; !ok {
			return nil, fmt.Errorf("неизвестное значение %q", v)
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return allowedSet[out[i]] < allowedSet[out[j]] })
	return out, nil
}

func (s Settings) HasRoom(choice string) bool  { return contains(s.Rooms, choice) }
func (s Settings) HasFloor(choice string) bool { return contains(s.Floors, choice) }

func (s *Settings) ToggleRoom(choice string) error {
	if !contains(RoomChoices, choice) {
		return fmt.Errorf("неизвестный вариант комнат %q", choice)
	}
	values, err := toggle(s.Rooms, choice)
	if err != nil {
		return err
	}
	s.Rooms = values
	return s.NormalizeAndValidate()
}

func (s *Settings) ToggleFloor(choice string) error {
	if !contains(FloorChoices, choice) {
		return fmt.Errorf("неизвестный этаж %q", choice)
	}
	values, err := toggle(s.Floors, choice)
	if err != nil {
		return err
	}
	s.Floors = values
	return s.NormalizeAndValidate()
}

func (s *Settings) ToggleOrdinary() error {
	s.IncludeOrdinary = !s.IncludeOrdinary
	if !s.IncludeOrdinary && !s.IncludePromoted {
		s.IncludeOrdinary = true
		return errors.New("нельзя отключить оба типа")
	}
	return nil
}

func (s *Settings) TogglePromoted() error {
	s.IncludePromoted = !s.IncludePromoted
	if !s.IncludeOrdinary && !s.IncludePromoted {
		s.IncludePromoted = true
		return errors.New("нельзя отключить оба типа")
	}
	return nil
}

func toggle(values []string, value string) ([]string, error) {
	out := append([]string(nil), values...)
	for i, v := range out {
		if v == value {
			if len(out) == 1 {
				return values, errors.New("нельзя снять последний вариант")
			}
			return append(out[:i], out[i+1:]...), nil
		}
	}
	return append(out, value), nil
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func CardMatches(s Settings, card model.Card) (bool, string) {
	if card.Promoted && !s.IncludePromoted {
		return false, "продвигаемые отключены"
	}
	if !card.Promoted && !s.IncludeOrdinary {
		return false, "обычные отключены"
	}
	if card.Price != nil && !priceMatches(s, *card.Price) {
		return false, "цена"
	}
	if card.Rooms != nil && !s.HasRoom(roomChoice(*card.Rooms)) {
		return false, "комнаты"
	}
	if card.Floor != nil && !s.HasFloor(floorChoice(*card.Floor)) {
		return false, "этаж"
	}
	return true, ""
}

func AdMatches(s Settings, ad model.Ad) (bool, string) {
	if ok, reason := CardMatches(s, ad.Card); !ok {
		return false, reason
	}
	if (s.PriceMin != nil || s.PriceMax != nil) && ad.Price == nil {
		return false, "цена не указана"
	}
	if ad.Price != nil && !priceMatches(s, *ad.Price) {
		return false, "цена"
	}
	if ad.Rooms == nil {
		if !s.HasRoom(UnknownChoice) {
			return false, "комнаты не указаны"
		}
	} else if !s.HasRoom(roomChoice(*ad.Rooms)) {
		return false, "комнаты"
	}
	if ad.Floor == nil {
		if !s.HasFloor(UnknownChoice) {
			return false, "этаж не указан"
		}
	} else if !s.HasFloor(floorChoice(*ad.Floor)) {
		return false, "этаж"
	}
	if s.SellerAdsLimit > 0 {
		if ad.SellerAds == nil {
			return false, "число объявлений автора не указано"
		}
		if *ad.SellerAds >= s.SellerAdsLimit {
			return false, "слишком много объявлений автора"
		}
	}

	text := normalizeText(ad.Title + "\n" + ad.Description)
	for _, negative := range s.NegativeWords {
		if strings.Contains(text, negative) {
			return false, "негативное слово: " + negative
		}
	}
	return true, ""
}

func priceMatches(s Settings, price int) bool {
	if s.PriceMin != nil && price < *s.PriceMin {
		return false
	}
	if s.PriceMax != nil && price > *s.PriceMax {
		return false
	}
	return true
}

func roomChoice(value int) string {
	if value >= 6 {
		return "6+"
	}
	return strconv.Itoa(value)
}

func floorChoice(value int) string {
	if value >= 14 {
		return "14+"
	}
	return strconv.Itoa(value)
}

func normalizeText(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))
	space := false
	for _, r := range s {
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
	return strings.TrimSpace(b.String())
}

func ParsePriceInput(input string) (*int, *int, error) {
	input = strings.TrimSpace(input)
	input = strings.NewReplacer("–", "-", "—", "-", "−", "-").Replace(input)
	if input == "-" || input == "" {
		return nil, nil, nil
	}
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return nil, nil, errors.New("формат: 3500-6000, -6000, 3500- или -")
	}
	parse := func(raw string) (*int, error) {
		raw = strings.ReplaceAll(strings.TrimSpace(raw), " ", "")
		if raw == "" {
			return nil, nil
		}
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 {
			return nil, errors.New("цена должна быть целым неотрицательным числом")
		}
		return &v, nil
	}
	min, err := parse(parts[0])
	if err != nil {
		return nil, nil, err
	}
	max, err := parse(parts[1])
	if err != nil {
		return nil, nil, err
	}
	if min != nil && max != nil && *min > *max {
		return nil, nil, errors.New("минимальная цена больше максимальной")
	}
	return min, max, nil
}

func ParseNegativeInput(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "-" || input == "" {
		return nil, nil
	}
	items := strings.Split(input, ",")
	s := DefaultSettings()
	s.NegativeWords = items
	if err := s.NormalizeAndValidate(); err != nil {
		return nil, err
	}
	return s.NegativeWords, nil
}

func ParsePollIntervalInput(input string) (int, int, error) {
	input = strings.TrimSpace(input)
	input = strings.NewReplacer("–", "-", "—", "-").Replace(input)
	parts := strings.Split(input, "-")
	if len(parts) == 1 {
		parts = []string{parts[0], parts[0]}
	}
	if len(parts) != 2 {
		return 0, 0, errors.New("формат интервала: 10-30 или 15")
	}
	parse := func(raw string) (int, error) {
		value, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil {
			return 0, errors.New("интервал должен быть целым числом минут")
		}
		return value, nil
	}
	min, err := parse(parts[0])
	if err != nil {
		return 0, 0, err
	}
	max, err := parse(parts[1])
	if err != nil {
		return 0, 0, err
	}
	s := DefaultSettings()
	s.PollMinMinutes = min
	s.PollMaxMinutes = max
	if err := s.NormalizeAndValidate(); err != nil {
		return 0, 0, err
	}
	return min, max, nil
}
