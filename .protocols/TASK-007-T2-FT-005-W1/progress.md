# Progress — TASK-007-T2-FT-005-W1

## Current status
- state: verifying
- last update: 2026-10-06
- Executor handoff complete; indexed lifecycle remains in_progress. No final verdict asserted.

## What was done
- Added pure bounds validator and strict price predicate; no bounds accepts any monetary availability, inclusive bounds require confirmed TJS. Invalid negative/reversed bounds reject.
- Added explicit expected table: 54 price vectors, 11 bounds-validation cases, two commodity vectors. Production API has no phrase/housing/seller/promotion inputs.
- Minimal Memory Bank routing/changelog appended. Rental settings.go/settings_test.go SHA256 unchanged: a0d8649517d97c2c5d6e57e745eedf5f5f983da17cc36f8792278c3a39e1c312 / 8f354d286b8688015ebe911d9e334bc15609b0bcc1db321b32bdd78feeec95ae.
- Actual files: the two keyword filter files; selected card status; five task protocol files; task evidence files; boundary-map/changelog; PAPERCUTS/gpt-6.1-sol __ 10-06-2026 22.57.md. Temporary baseline test removed.
- Advisory source hints kept. No hard allowed-write list; no forbidden scope touched. Existing dirty work preserved, tier unchanged, no new contract/state/ownership decisions.

## Commands run (with results)
- Baseline command below: exit 1, compiling behavioral RED.
- GREEN command below: exit 0.
- Required gate: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/filter'` → exit 0, focused-gate.log reports `ok .../internal/filter 0.004s`.
- Static: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go vet ./internal/filter && test -z "$(gofmt -l internal/filter/keyword_search.go internal/filter/keyword_search_test.go)"'` → exit 0, static-gate.log empty (success); first equivalent static pass also exit 0.
- Formatting before GREEN: existing builder `gofmt -w internal/filter/keyword_search.go internal/filter/keyword_search_test.go` → exit 0.
- `git diff --check` → exit 0.

## Claim-linked RED / GREEN (T2/T3)
- attempt: 1
- applicability: applicable
- accepted claim locator(s): FT-005-AC-004
- accepted not-applicable reason and alternative proof: none; initial GREEN subcases preserved.
- RED command/probe: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 -v ./internal/filter -run ^TestKeywordPriceClaimBaseline$'` using unchanged compiling baseline_probe_test.go copied temporarily to internal/filter/keyword_price_baseline_test.go.
- RED observation and evidence: pre-production 19 failing and 36 passing subcases; bounded missing/negotiable/unconfirmed/foreign prices admitted by CardMatches, eligible commodity rejected by AdMatches (`комнаты не указаны`). baseline-red.log, baseline_probe_test.go; exact AC locator above. Setup-only Docker nested mount exit 125 excluded, workaround identical probe.
- GREEN command/probe: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'gofmt -w internal/filter/keyword_search.go internal/filter/keyword_search_test.go; CGO_ENABLED=1 go test -count=1 -v ./internal/filter -run "^Test(Keyword|ValidateKeyword)"'`.
- GREEN observation and evidence: all 54 equivalent price vectors accepted/rejected exactly as required; 11 validation cases and two commodity vectors pass. claim-green.log and keyword_search_test.go. Required final package gate passes existing rental tests too; supporting evidence only, fresh verifier proof due.
- claim-equivalent probe changes and rationale: baseline rental entrypoints switched to keyword-only predicate per plan; same monetary matrix; computed baseline oracle replaced with explicit expected rows without weakening vectors. Commodity fields expanded within accepted no-housing rule; price validation added. No baseline API modified or intentionally broken.
- T3 isolation/cleanup/permission evidence: not T3; in-memory values, no DB/network except network-disabled Docker. Temporary baseline package test removed before final gates.

## Reuse Candidates (optional)
None offered; all execute results are supporting evidence. Independent /verify runs its own claim probes.

## Evidence links
- `.tasks/TASK-007-T2-FT-005-W1/TASK-007-T2-FT-005-W1-acceptance-evidence.md`
- `.tasks/TASK-007-T2-FT-005-W1/baseline_probe_test.go`, baseline-red.log, claim-green.log, focused-gate.log, static-gate.log.

## Open issues / risks
No material blocker. Price gate proves filtering only; source/creation/delivery/deletion are other task owners.

## Next step (single concrete action)
`/verify TASK-007-T2-FT-005-W1` in fresh Reviewer context; /root then decides closure and wave sync.
