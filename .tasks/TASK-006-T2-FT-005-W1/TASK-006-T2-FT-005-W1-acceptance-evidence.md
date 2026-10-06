# Acceptance evidence — TASK-006-T2-FT-005-W1

Owned claim FT-005-AC-003 (REQ-011/013), attempt 1.
Task entered in_progress after preflight/protocol and before fixtures/probe/production writes.

## Before production changes: compiling RED and preserved GREEN
`baseline_probe_test.go` is preserved historical source. Docker copies repository into /tmp/probe and copies this probe into internal/somon; container removal cleans isolation. Existing production entrypoints only (ParseCategory/FetchCategory), identical checked-in fixtures.

Command:
```
docker run --rm --network none -v "$PWD:/src:ro" -w /src somon-price-hotfix-builder:latest sh -c 'cp -a /src /tmp/probe && cp /src/.tasks/TASK-006-T2-FT-005-W1/baseline_probe_test.go /tmp/probe/internal/somon/keyword_baseline_test.go && cd /tmp/probe && CGO_ENABLED=1 go test -count=1 -v ./internal/somon -run TestKeywordSourceBaseline'
```
Exit 1, [baseline-red.log](baseline-red.log).
- small: expected [21000001], observed [21000002], foreign recommendation.
- empty H1 0: expected [], observed [21000002], foreign recommendation.
- Original DOM parser only recognizes apartment title; largest RSC array supplies recommendation. This is behavioral absence of source claim, no syntax/setup failure.
- Initial GREEN: captured /search/vose/ q=стол & стул and relevance, blocked page typed error, malformed body error, HTTP 403 typed error. Preserved transport implementation.

## After production changes: claim-equivalent GREEN
Same small/empty fixtures: ParseKeywordSearch and FetchKeywordSearch return exact primary IDs [21000001] / []; region recommendation ID 21000002 never enters. Multiple primary IDs retain order/dedup and native matching (no local phrase predicate/apartment fields).
16 category×city scopes captured by httptest: all/furniture/tables_chairs/services × country/dushanbe/vose/dangara; selected city retained, HTTPS somon.tj native path, encoded q, relevance; extra fields absent. Empty phrase and unsupported paths/keys rejected before fetching.
Malformed/unconfirmed empty/RSC-only/untitled/contradictory zero bodies error. Blocked body/403/429 typed failures, Retry-After=60 -> one minute. Raw successful HTTP body preserved even parse failure. Missing price/currency/city/image do not become guessed values; explicit 0 TJS preserved, unknown/foreign currency never becomes TJS; photo count excluded from price.

[focused-green.log](focused-green.log): verbose full source/htmlx proof PASS. A subsequent production extraction correction added negative-price rejection and product:price:amount metadata support; final mandatory gate below covers final state.

```
docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/somon ./internal/htmlx'
```
Exit 0, [final-gate.log](final-gate.log): internal/somon and internal/htmlx PASS, including existing rental/price fixtures/tests unchanged.
Focused vet `docker run --rm --network none -v "$PWD:/src:ro" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go vet ./internal/somon ./internal/htmlx'` exit 0, [focused-vet.log](focused-vet.log) empty output.
gofmt final modified Go files completed exit 0.

## Evidence limits
Representative sanitized HTML reflects documented primary/region/H1 observations and existing source markup patterns; not fetched/captured from blocked pages. No live compatibility claim. RSC alone cannot prove primary membership; no fallback resurrection. Source only exposes visible initial catalog city labels, missing fields remain absent. Rental behavior unchanged; no UI/filter/store/scheduler/deletion/production edits. Evidence is supporting executor evidence, not independent verification and not a reuse receipt.
