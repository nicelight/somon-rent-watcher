package filter

import "errors"

// ValidateKeywordPriceBounds preserves absent bounds and accepts zero inclusively.
func ValidateKeywordPriceBounds(min, max *int) error {
	if min != nil && *min < 0 {
		return errors.New("минимальная цена не может быть отрицательной")
	}
	if max != nil && *max < 0 {
		return errors.New("максимальная цена не может быть отрицательной")
	}
	if min != nil && max != nil && *min > *max {
		return errors.New("минимальная цена больше максимальной")
	}
	return nil
}

// KeywordPriceMatches applies only strict monetary bounds. Native phrase matching
// belongs to Somon; rental criteria and fallback do not participate.
func KeywordPriceMatches(min, max, price *int, currency string) bool {
	if ValidateKeywordPriceBounds(min, max) != nil {
		return false
	}
	if min == nil && max == nil {
		return true
	}
	if price == nil || currency != "TJS" {
		return false
	}
	return (min == nil || *price >= *min) && (max == nil || *price <= *max)
}
