# Plan — TASK-006-T2-FT-005-W1

## Goal / non-goals
Implement only FT-005-AC-003 native scoped primary feed and supported catalog. No filter, UI, SQLite, scheduler, delivery, deletion, rental behavior changes or production work.

## Inputs
Indexed task card, FT-005 AC-003, REQ-011/013 and direct canonical inputs recorded context.md; latest reviewed revision equals Global Planning Revision 1.

## Preflight-confirmed change surface
- internal/somon/keyword_search.go + keyword_search_test.go: catalog, URL, fetch/parser.
- internal/model/keyword_search.go + additive internal/model/ad.go fields: passive payload.
- testdata/keyword-search-primary-small.html / primary-empty.html: representative source-observation fixtures.
- Task-local historical compiling baseline probe/evidence, full protocol; existing MB source navigation/changelog update.
- Existing client/parser helpers reused without changing rental behavior. Additional passive Card fields stay same owned outcome; hard write_boundary not set, forbidden_scope clear.

## Claim-linked RED / GREEN
Applicability applicable; owned FT-005-AC-003. First run existing compiling ParseCategory/FetchCategory with checked-in small/empty fixtures; expected IDs only primary (small 21000001, empty none) versus region/RSC recommendation leakage = RED. Preserve initial GREEN for encoded requests/transport errors/malformed bodies. After production change run equivalent ParseKeywordSearch/FetchKeywordSearch fixtures and exact native paths/q/ordering for four category keys and country/city, empty and recommendations, block/malformed. No artificial missing-symbol RED. All probes run fresh httptest with cleanup in network-disabled Docker. No production side effects.

## Applicable gates
- Required focused Docker: CGO_ENABLED=1 go test -count=1 ./internal/somon ./internal/htmlx.
- gofmt modified Go files; go vet focused packages when implementation complete.

## MB-SYNC handoff / owner
/root explicit standalone owner in .protocols/FT-005/plan.md closes only after independent /verify. Minimal WHY/WHERE/source routing + changelog during execution; final RTM/index status sync remains owner. No feature closure/commit.
