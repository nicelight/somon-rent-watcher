package model

// KeywordSearch is a passive saved-search payload, independent of rental settings.
type KeywordSearch struct {
	ID          int64
	Phrase      string
	CategoryKey string
	CityKey     string
	PriceMin    *int
	PriceMax    *int
	Enabled     bool
	Revision    int64
}
