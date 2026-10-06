package app

import (
	"errors"
	"strings"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
)

func validateKeywordSearch(s *model.KeywordSearch) error {
	s.Phrase = strings.TrimSpace(s.Phrase)
	if s.Phrase == "" {
		return errors.New("введите слово или фразу для поиска")
	}
	if _, err := somon.KeywordSearchURL(s.Phrase, s.CategoryKey, s.CityKey); err != nil {
		return errors.New("выберите категорию и город из списка")
	}
	return filter.ValidateKeywordPriceBounds(s.PriceMin, s.PriceMax)
}

func (a *App) ListKeywordSearches() ([]model.KeywordSearch, error) {
	return a.store.ListKeywordSearches()
}
func (a *App) KeywordSearch(id int64) (model.KeywordSearch, bool, error) {
	return a.store.KeywordSearch(id)
}

func (a *App) CreateKeywordSearch(s model.KeywordSearch) (model.KeywordSearch, error) {
	if err := validateKeywordSearch(&s); err != nil {
		return model.KeywordSearch{}, err
	}
	// The application owns identity/revision and the summary-before-enable flow.
	s.ID = 0
	s.Revision = 1
	s.Enabled = false
	return a.store.CreateKeywordSearch(s)
}

func (a *App) UpdateKeywordSearch(s model.KeywordSearch) (model.KeywordSearch, bool, error) {
	a.keywordMu.Lock()
	defer a.keywordMu.Unlock()
	current, found, err := a.store.KeywordSearch(s.ID)
	if err != nil || !found {
		return model.KeywordSearch{}, found, err
	}
	if s.Revision != current.Revision {
		return model.KeywordSearch{}, true, errors.New("поиск изменён; откройте его настройки снова")
	}
	if err = validateKeywordSearch(&s); err != nil {
		return model.KeywordSearch{}, true, err
	}
	s.Enabled = current.Enabled
	return a.store.UpdateKeywordSearch(s)
}

func (a *App) SetKeywordSearchEnabled(id int64, enabled bool) (model.KeywordSearch, bool, error) {
	a.keywordMu.Lock()
	defer a.keywordMu.Unlock()
	s, found, err := a.store.KeywordSearch(id)
	if err != nil || !found {
		return model.KeywordSearch{}, found, err
	}
	if s.Enabled == enabled {
		return s, true, nil
	}
	s.Enabled = enabled
	return a.store.UpdateKeywordSearch(s)
}
