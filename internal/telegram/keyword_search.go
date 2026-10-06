package telegram

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
)

// KeywordBackend extends the existing backend only for the new capability;
// existing rental consumers and their interface remain compatible.
type KeywordBackend interface {
	ListKeywordSearches() ([]model.KeywordSearch, error)
	KeywordSearch(int64) (model.KeywordSearch, bool, error)
	CreateKeywordSearch(model.KeywordSearch) (model.KeywordSearch, error)
	UpdateKeywordSearch(model.KeywordSearch) (model.KeywordSearch, bool, error)
	SetKeywordSearchEnabled(int64, bool) (model.KeywordSearch, bool, error)
}

func (b *Bot) keywordBackend() (KeywordBackend, error) {
	backend, ok := b.backend.(KeywordBackend)
	if !ok {
		return nil, errors.New("поиски недоступны")
	}
	return backend, nil
}

func (b *Bot) sendKeywordList(ctx context.Context, chatID int64) error {
	backend, err := b.keywordBackend()
	if err != nil {
		return err
	}
	searches, err := backend.ListKeywordSearches()
	if err != nil {
		return err
	}
	text := "<b>Мои поиски</b>\nВыберите поиск для настройки или создайте новый."
	if len(searches) == 0 {
		text += "\nСохранённых поисков пока нет."
	}
	rows := make([][]InlineKeyboardButton, 0, len(searches)+2)
	for _, s := range searches {
		status := "⏸"
		if s.Enabled {
			status = "▶"
		}
		rows = append(rows, []InlineKeyboardButton{{Text: status + " " + truncate(s.Phrase, 60), CallbackData: fmt.Sprintf("ks:view:%d", s.ID)}})
	}
	rows = append(rows, []InlineKeyboardButton{{Text: "Новый поиск", CallbackData: "ks:new"}}, []InlineKeyboardButton{{Text: "Фильтр квартир", CallbackData: "ks:rental"}})
	_, err = b.client.SendMessage(ctx, chatID, text, &InlineKeyboardMarkup{InlineKeyboard: rows})
	return err
}

func (b *Bot) sendKeywordSummary(ctx context.Context, chatID int64, s model.KeywordSearch) error {
	category, city := "", ""
	for _, c := range somon.KeywordCategories() {
		if c.Key == s.CategoryKey {
			category = c.Label
		}
	}
	for _, c := range somon.KeywordCities() {
		if c.Key == s.CityKey {
			city = c.Label
		}
	}
	price := "без ограничений"
	if s.PriceMin != nil || s.PriceMax != nil {
		min, max := "—", "—"
		if s.PriceMin != nil {
			min = strconv.Itoa(*s.PriceMin)
		}
		if s.PriceMax != nil {
			max = strconv.Itoa(*s.PriceMax)
		}
		price = min + " – " + max + " сомони"
	}
	state, toggle, action := "на паузе", "Включить поиск", "enable"
	if s.Enabled {
		state, toggle, action = "включён", "Поставить поиск на паузу", "disable"
	}
	text := "<b>Поиск: " + html.EscapeString(truncate(s.Phrase, 180)) + "</b>\n" +
		"Категория: " + html.EscapeString(category) + "\nГеография: " + html.EscapeString(city) + "\nЦена: " + price + "\nСостояние: " + state + "\n\nПроверьте условия перед включением."
	button := func(label, action string) InlineKeyboardButton {
		return InlineKeyboardButton{Text: label, CallbackData: fmt.Sprintf("ks:%s:%d", action, s.ID)}
	}
	rows := [][]InlineKeyboardButton{{button("Изменить фразу", "phrase")}, {button("Категория", "category"), button("География", "city")}, {button("Цена от / до", "price")}, {button(toggle, action)}, {{Text: "Мои поиски", CallbackData: "ks:list"}}}
	_, err := b.client.SendMessage(ctx, chatID, text, &InlineKeyboardMarkup{InlineKeyboard: rows})
	return err
}

func (b *Bot) keywordMissing(ctx context.Context, chatID int64) error {
	if _, err := b.client.SendMessage(ctx, chatID, "Поиск больше не существует. Выберите поиск из текущего списка.", nil); err != nil {
		return err
	}
	return b.sendKeywordList(ctx, chatID)
}

func (b *Bot) processKeywordCallback(ctx context.Context, q *CallbackQuery) error {
	// Authorization is checked by processCallback before this route. Acknowledge
	// before reading settings or sending any fresh output.
	if err := b.client.AnswerCallbackQuery(ctx, q.ID, "", false); err != nil {
		return err
	}
	chatID := q.Message.Chat.ID
	b.setPending(q.From.ID, chatID, "")
	err := b.applyKeywordCallback(ctx, q)
	if err != nil {
		_, sendErr := b.client.SendMessage(ctx, chatID, "Ошибка: "+html.EscapeString(err.Error()), nil)
		if sendErr != nil {
			return sendErr
		}
		return err
	}
	return nil
}

func (b *Bot) applyKeywordCallback(ctx context.Context, q *CallbackQuery) error {
	chatID := q.Message.Chat.ID
	switch q.Data {
	case "ks:list":
		return b.sendKeywordList(ctx, chatID)
	case "ks:rental":
		return b.sendMainMenu(ctx, chatID)
	case "ks:new":
		if _, err := b.keywordBackend(); err != nil {
			return err
		}
		b.setPending(q.From.ID, chatID, "ks:0:0:phrase")
		_, err := b.client.SendMessage(ctx, chatID, "Введите слово или фразу для нового поиска.\n/cancel — отмена.", nil)
		return err
	}
	parts := strings.Split(q.Data, ":")
	if len(parts) < 3 || len(parts) > 4 {
		return errors.New("неизвестная кнопка поиска")
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || id <= 0 {
		return errors.New("неверный поиск")
	}
	backend, err := b.keywordBackend()
	if err != nil {
		return err
	}
	s, found, err := backend.KeywordSearch(id)
	if err != nil {
		return err
	}
	if !found {
		return b.keywordMissing(ctx, chatID)
	}
	action := parts[1]
	if len(parts) == 4 {
		switch action {
		case "setcategory":
			s.CategoryKey = parts[3]
		case "setcity":
			s.CityKey = parts[3]
		default:
			return errors.New("неизвестная кнопка поиска")
		}
		s, found, err = backend.UpdateKeywordSearch(s)
	} else {
		switch action {
		case "view":
			return b.sendKeywordSummary(ctx, chatID, s)
		case "category", "city":
			var rows [][]InlineKeyboardButton
			if action == "category" {
				for _, c := range somon.KeywordCategories() {
					rows = append(rows, []InlineKeyboardButton{{Text: c.Label, CallbackData: fmt.Sprintf("ks:setcategory:%d:%s", id, c.Key)}})
				}
			} else {
				for _, c := range somon.KeywordCities() {
					rows = append(rows, []InlineKeyboardButton{{Text: c.Label, CallbackData: fmt.Sprintf("ks:setcity:%d:%s", id, c.Key)}})
				}
			}
			rows = append(rows, []InlineKeyboardButton{{Text: "К настройкам поиска", CallbackData: fmt.Sprintf("ks:view:%d", id)}})
			text := "Выберите категорию:"
			if action == "city" {
				text = "Выберите географию:"
			}
			_, err = b.client.SendMessage(ctx, chatID, text, &InlineKeyboardMarkup{InlineKeyboard: rows})
			return err
		case "phrase", "price":
			b.setPending(q.From.ID, chatID, fmt.Sprintf("ks:%d:%d:%s", id, s.Revision, action))
			text := "Введите новое слово или фразу.\n/cancel — отмена."
			if action == "price" {
				text = "Введите цену от / до в сомони: <code>100-200</code>, <code>-200</code>, <code>100-</code> или <code>-</code> без ограничений.\n/cancel — отмена."
			}
			_, err = b.client.SendMessage(ctx, chatID, text, nil)
			return err
		case "enable", "disable":
			s, found, err = backend.SetKeywordSearchEnabled(id, action == "enable")
		default:
			return errors.New("неизвестная кнопка поиска")
		}
	}
	if err != nil {
		return err
	}
	if !found {
		return b.keywordMissing(ctx, chatID)
	}
	return b.sendKeywordSummary(ctx, chatID, s)
}

func (b *Bot) applyKeywordInput(ctx context.Context, userID, chatID int64, text, pending string) error {
	parts := strings.Split(pending, ":")
	if len(parts) != 4 {
		return errors.New("неверное ожидаемое действие поиска")
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return err
	}
	revision, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return err
	}
	backend, err := b.keywordBackend()
	if err != nil {
		return err
	}
	var s model.KeywordSearch
	found := true
	if id == 0 {
		s = model.KeywordSearch{CategoryKey: "all", CityKey: "country"}
	} else {
		s, found, err = backend.KeywordSearch(id)
		if err != nil {
			return err
		}
		if !found {
			b.setPending(userID, chatID, "")
			return b.keywordMissing(ctx, chatID)
		}
		if s.Revision != revision {
			b.setPending(userID, chatID, "")
			if _, err = b.client.SendMessage(ctx, chatID, "Поиск изменился после начала ввода. Проверьте текущие условия и повторите действие.", nil); err != nil {
				return err
			}
			return b.sendKeywordSummary(ctx, chatID, s)
		}
	}
	switch parts[3] {
	case "phrase":
		s.Phrase = text
	case "price":
		s.PriceMin, s.PriceMax, err = filter.ParsePriceInput(text)
	default:
		err = errors.New("неизвестное действие поиска")
	}
	if err == nil {
		if id == 0 {
			s, err = backend.CreateKeywordSearch(s)
		} else {
			s, found, err = backend.UpdateKeywordSearch(s)
		}
	}
	if err != nil {
		_, sendErr := b.client.SendMessage(ctx, chatID, "Ошибка: "+html.EscapeString(err.Error())+"\nПовторите ввод или /cancel.", nil)
		return sendErr
	}
	b.setPending(userID, chatID, "")
	if !found {
		return b.keywordMissing(ctx, chatID)
	}
	return b.sendKeywordSummary(ctx, chatID, s)
}
