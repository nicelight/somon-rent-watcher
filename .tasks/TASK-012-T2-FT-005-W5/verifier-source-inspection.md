# FT-005-AC-009 source and architecture observation

Verifier read current source and task diff independently.

- `internal/app/keyword_polling.go:74–78`: IDs allocated at current `cards` length, populated from each card.ID, passed with selected s.ID to bounded Store lookup. No full diagnostic lookup in production polling.
- `internal/store/keyword_search_cgo.go:158–163`: empty path returns before helper, lock, prepare, SQL access.
- `internal/store/keyword_search_cgo.go:170–185`: SELECT has monitor_id=? and ad_id IN bound placeholders; monitor parameter first, requested IDs after. No full-read/filter-afterward path.
- Full `KeywordAdStates(monitorID)` continues to invoke the same reader without ID predicate, preserving diagnostic compatibility. Existing production matches are compared in App using unchanged delivered/evaluated revision rules.
- Actual TASK012 diff is only these reader/caller changes, no table initialization, data-write/delete, scheduler/Telegram/price-filter changes; nocgo gets the corresponding existing errCGODisabled stub.
- Allowed graph row Polling Application -> Persistence Adapter / Persistence contract used. App selects current feed IDs and controls dedup/orchestration; Store alone owns SQL. No reverse dependency, separate state source, or new architecture unit.

Reproduce inspection with `git diff -- internal/app/keyword_polling.go internal/store/keyword_search_cgo.go internal/store/keyword_search_nocgo.go internal/store/keyword_search_test.go` and `rg -n 'KeywordAdStates\(|KeywordAdStatesForIDs\(' internal cmd`.
