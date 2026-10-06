# Plan — TASK-007-T2-FT-005-W1

## Goal
Pure strict optional inclusive keyword price eligibility and independent bounds validation.

## Non-goals
Native phrase/parser, UI/storage/polling/delivery/deletion, apartment criteria/fallback, old tasks and production.

## Inputs / source specs
- Indexed task card, FT005 AC004, REQ011, boundary-map filtering/shared-data/keyword shapes, architecture proposal and testing risk-based checks.

## Constraints / invariants (MUST / NEVER)
- Filtering owns pure validation/predicate; native phrase matching stays with Somon.
- Preserve missing price and absent-vs-zero bounds; only confirmed TJS matches bounded prices.
- Never alter rental Settings/CardMatches/AdMatches; no source/UI/store calls or forbidden external actions.

## Scope
- In scope: two advisory keyword filter files, task-owned protocol/evidence and minimal Memory Bank WHY/WHERE/changelog.
- Out of scope: all other feature outcomes and TASK001 historical in_progress state.

## Proposed changes
### Preflight-confirmed change surface
- Expected hints kept: internal/filter/keyword_search.go and internal/filter/keyword_search_test.go.
- Additional same-outcome files/areas and rationale: protocol/evidence/status bookkeeping, canonical routing/changelog, one environment papercut.
- Hard `write_boundary` present and satisfied: not set.
- `forbidden_scope` / stop-condition check: clear; no violation.

## Applicable quality gates
- [x] Focused required package gate: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/filter'` (exit 0).
- [x] Local go vet and new-file formatting check (exit 0); diff whitespace check (exit 0).

## Claim-linked RED / GREEN (T2/T3)
- applicability: applicable
- accepted claim locator(s): FT-005-AC-004
- planned test/probe and environment: same in-memory absent/min/max/both/zero bounds and 99/100/200/201/missing/negotiable/currency vectors; commodity with no apartment/seller fields. Network-disabled Docker.
- observable RED: existing CardMatches admits bounded missing/unconfirmed price; AdMatches rejects eligible commodity.
- corresponding GREEN: keyword-only predicate has exact expected eligibility; strict known TJS and no housing/text/promotion requirements.
- accepted not-applicable reason and alternative proof: none; already-GREEN baseline subcases preserved.
- T3 isolation, safe rerun, cleanup, and permission boundary: not T3; no mutable external state, fresh in-memory fixtures.

## MB-SYNC handoff / owner
- Owner: explicit standalone /root, `.protocols/FT-005/plan.md#task006-closure-and-task007-selection`; delegated executor cannot close.
- WHY/WHERE/changelog updated; index already routes boundary-map, no new router needed. RTM remains feature-partial.
- Final lifecycle decision and W1 sync: /root after independent /verify; no full mb-sync or doctor inside /exe.

## Definition of done
Execution handoff complete with truthful RED/GREEN and required gate. Task remains in_progress pending independent /verify and owner closure.
