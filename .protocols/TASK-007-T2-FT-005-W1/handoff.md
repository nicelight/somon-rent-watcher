# Handoff — TASK-007-T2-FT-005-W1

## Summary
Strict pure keyword price eligibility and independent bounds validation implemented for FT-005-AC-004 / REQ-011. Any bound requires known TJS; bounds inclusive, zero valid, missing bounds unrestricted. Rental logic unchanged.

## Where to look
- `internal/filter/keyword_search.go`: ValidateKeywordPriceBounds and KeywordPriceMatches accept only bounds/monetary values.
- `internal/filter/keyword_search_test.go`: explicit matrix/commodity/validation proof.
- `.memory-bank/contracts/boundary-map.md#keyword-price-implementation-routing`: WHY/WHERE.
- Advisory deviations: task bookkeeping, minimal routing/changelog, environment papercut only; no hard path list, forbidden scope untouched.

## How to run / verify
- Required gate: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/filter'`.
- Execute results: required gate, static and formatting exit 0. Baseline RED exit 1 precedes production writes; claim-equivalent GREEN exit 0.
- Claim-linked RED/GREEN evidence: progress.md attempt 1 / Claim-linked RED / GREEN and `.tasks/TASK-007-T2-FT-005-W1/TASK-007-T2-FT-005-W1-acceptance-evidence.md`.
- Current-attempt reuse candidate locators: none; execute evidence supporting-only.
- Superseded receipts: none. Initial 36 baseline passing subcases preserved; setup bind failure explicitly excluded from RED.

## Known issues
External HEAD changed from 7782bc1 to d99e3a7 during execution; prepared protocol/temporary probe became tracked. Temporary probe removal is now a tracked deletion, with its exact source retained under .tasks. This agent did not commit or run any git mutation.

No material blocker. Source dependency proof not adopted. No commits/push/deployment/live Telegram/working DB/secrets.

## Follow-ups
Fresh `/verify TASK-007-T2-FT-005-W1`. Leave in_progress; /root owns closure then W1 sync, following tasks own integration.

## Explicit owner closure

/root records done after independent /verify PASS; task.verify contains the evidence.
W1 mb-sync and caller-owned lint/strict doctor follow before TASK008.
