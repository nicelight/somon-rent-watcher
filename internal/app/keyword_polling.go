package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
)

// One scheduler rotates the starting monitor, including rental, and passes one
// remaining detail budget through all owners. No search has a silent baseline.
func (a *App) pollOnce(ctx context.Context) error {
	searches, err := a.store.ListKeywordSearches()
	if err != nil {
		return err
	}
	enabled := make([]model.KeywordSearch, 0, len(searches))
	for _, s := range searches {
		if s.Enabled {
			enabled = append(enabled, s)
		}
	}
	count := len(enabled) + 1
	start := a.pollStart % count
	a.pollStart = (start + 1) % count
	remaining := a.cfg.MaxDetailsPerPoll
	total := processStats{}
	cardCount := 0
	var cycleErrors []error
	a.setCycleStats(0, 0, 0)
	for offset := 0; offset < count; offset++ {
		index := (start + offset) % count
		if index == 0 {
			a.setCycleStats(0, 0, 0)
			err = a.pollRentalOnce(ctx, &remaining)
			a.statusMu.RLock()
			cardCount += a.status.LastCardCount
			total.NewIDs += a.status.LastNewCount
			total.Sent += a.status.LastSentCount
			a.statusMu.RUnlock()
		} else {
			var stats processStats
			var cards int
			stats, cards, err = a.pollKeywordSearch(ctx, enabled[index-1], &remaining)
			total.NewIDs += stats.NewIDs
			total.Sent += stats.Sent
			cardCount += cards
		}
		if err != nil {
			if _, _, blocked := somon.IsBlocked(err); blocked {
				a.setCycleStats(cardCount, total.NewIDs, total.Sent)
				return err
			}
			cycleErrors = append(cycleErrors, err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	a.setCycleStats(cardCount, total.NewIDs, total.Sent)
	return errors.Join(cycleErrors...)
}

func (a *App) pollKeywordSearch(ctx context.Context, s model.KeywordSearch, remaining *int) (processStats, int, error) {
	stats := processStats{}
	cards, _, err := a.somon.FetchKeywordSearch(ctx, s.Phrase, s.CategoryKey, s.CityKey)
	if err != nil {
		return stats, 0, err
	}
	ids := make([]int64, len(cards))
	for i, card := range cards {
		ids[i] = card.ID
	}
	history, err := a.store.KeywordAdStatesForIDs(s.ID, ids)
	if err != nil {
		return stats, len(cards), err
	}
	referer, _ := somon.KeywordSearchURL(s.Phrase, s.CategoryKey, s.CityKey)
	reject := func(id int64) error { return a.store.RecordKeywordAdState(s.ID, id, s.Revision, false) }
	for _, card := range cards {
		old := history[card.ID]
		if old.Delivered || old.EvaluatedRevision == s.Revision {
			continue
		}
		stats.NewIDs++
		if !filter.KeywordPriceMatches(s.PriceMin, s.PriceMax, card.Price, card.Currency) {
			if err = reject(card.ID); err != nil {
				return stats, len(cards), err
			}
			continue
		}
		if *remaining <= 0 {
			continue
		}
		*remaining--
		stats.DetailRequests++
		ad, _, detailErr := a.somon.FetchDetail(ctx, card, referer)
		if detailErr != nil {
			if _, _, blocked := somon.IsBlocked(detailErr); blocked {
				return stats, len(cards), detailErr
			}
			a.logger.Error("keyword detail failed", "search_id", s.ID, "ad_id", card.ID, "error", detailErr)
			continue
		}
		if !filter.KeywordPriceMatches(s.PriceMin, s.PriceMax, ad.Price, ad.Currency) {
			if err = reject(card.ID); err != nil {
				return stats, len(cards), err
			}
			continue
		}
		sent, sendErr := a.deliverKeywordAd(ctx, s, ad)
		if sendErr != nil {
			return stats, len(cards), sendErr
		}
		if sent {
			stats.Sent++
		}
	}
	return stats, len(cards), nil
}

// Mutations and the start of delivery share the application lock. Once send has
// started it cannot be recalled; a failed send leaves no terminal history.
func (a *App) deliverKeywordAd(ctx context.Context, s model.KeywordSearch, ad model.Ad) (bool, error) {
	a.keywordMu.Lock()
	defer a.keywordMu.Unlock()
	current, found, err := a.store.KeywordSearch(s.ID)
	if err != nil {
		return false, err
	}
	if !found || !current.Enabled || current.Revision != s.Revision {
		return false, nil
	}
	if err = a.bot.SendKeywordAd(ctx, s.Phrase, ad); err != nil {
		a.logger.Error("keyword Telegram delivery failed; will retry", "search_id", s.ID, "ad_id", ad.ID, "error", err)
		return false, nil
	}
	if err = a.store.RecordKeywordAdState(s.ID, ad.ID, s.Revision, true); err != nil {
		return false, fmt.Errorf("record keyword delivery: %w", err)
	}
	a.logger.Info("keyword ad sent", "search_id", s.ID, "ad_id", ad.ID, "sent_at", time.Now().UTC())
	return true, nil
}
