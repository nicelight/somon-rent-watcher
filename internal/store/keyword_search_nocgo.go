//go:build !cgo

package store

import "github.com/nicelight/somon-rent-watcher/internal/model"

func (*DB) ListKeywordSearches() ([]model.KeywordSearch, error) { return nil, errCGODisabled }
func (*DB) KeywordSearch(int64) (model.KeywordSearch, bool, error) {
	return model.KeywordSearch{}, false, errCGODisabled
}
func (*DB) CreateKeywordSearch(model.KeywordSearch) (model.KeywordSearch, error) {
	return model.KeywordSearch{}, errCGODisabled
}
func (*DB) UpdateKeywordSearch(model.KeywordSearch) (model.KeywordSearch, bool, error) {
	return model.KeywordSearch{}, false, errCGODisabled
}

func (*DB) KeywordAdStates(int64) (map[int64]model.KeywordAdState, error) { return nil, errCGODisabled }
func (*DB) RecordKeywordAdState(int64, int64, int64, bool) error          { return errCGODisabled }
func (*DB) DeleteKeywordSearch(int64) (bool, error)                       { return false, errCGODisabled }

func (*DB) KeywordAdStatesForIDs(int64, []int64) (map[int64]model.KeywordAdState, error) {
	return nil, errCGODisabled
}
