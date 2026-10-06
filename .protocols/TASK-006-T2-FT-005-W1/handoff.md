# Handoff — TASK-006-T2-FT-005-W1

## Summary
Native scoped keyword source/catalog, primary-only DOM parsing and confirmed zero/error. New Shared Data KeywordSearch and Card.Currency/City remain passive. Rental parser/client untouched.

## Where to look
- internal/somon/keyword_search.go and keyword_search_test.go.
- internal/model/keyword_search.go and additive ad.go fields.
- testdata/keyword-search-primary-small.html / primary-empty.html (provenance in comments).
- .memory-bank/contracts/current-integrations.md#keyword-source-implementation-routing; index/changelog navigation.
Advisory deviation: ad.go passive fields needed for accepted cards; client/parser untouched because separate entrypoints preserve rental behavior. No hard write_boundary; forbidden scope untouched.

## How to verify
Required gate:
```
docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/somon ./internal/htmlx'
```
Focused vet also passed. Claim FT-005-AC-003, attempt 1, honest compiling RED→equivalent GREEN: progress.md and .tasks/TASK-006-T2-FT-005-W1/TASK-006-T2-FT-005-W1-acceptance-evidence.md. All local httptest isolated/cleaned; network-disabled Docker. No reuse receipts; all logs supporting-only.

## Known limits
Representative fixtures, not live captures; no live source compatibility claim. RSC-only cannot establish primary scope and returns error. Unknown currency/missing price/city/image stay absent; source recognizes visible initial catalog city labels. Existing FetchDetail keeps its rental extraction; future keyword monitoring must use confirmed monetary evidence when applying strict price rules, never assume bare detail Price implies TJS.

## Follow-up
/verify TASK-006-T2-FT-005-W1. TASK status in_progress, no closure verdict. /root standalone owner recorded .protocols/FT-005/plan.md owns lifecycle and later sync. No commit/push/deployment/Telegram/DB changes. Other tasks untouched.
