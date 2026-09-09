---
description: Implementation plan for FT-001 exact-before-fallback polling.
status: active
feature: FT-001
planning_revision: 1
---

# Implementation Plan — FT-001 price fallback

## Goal

Make one active poll evaluate fresh exact matches first and, only after a
complete exact-empty result, send up to three closest eligible higher-price ads
without weakening the shared detail cap, seen-state, backoff, or delivery-retry
semantics.

## Scope and non-goals

In scope:

- Rework the existing polling candidate flow under the Polling Application
  owner so exact, fallback, rejected, deferred, and delivery-failed outcomes
  follow the accepted lifecycle.
- Add the minimum pure filtering support needed to evaluate all non-maximum-price
  conditions without changing `Settings` persistence.
- Add deterministic integration/unit coverage for FT-001-AC-001 through
  FT-001-AC-005 and run the repository-native build gate.

Out of scope:

- SQLite schema or stored-row interpretation changes.
- A new setting, dependency, worker, compatibility layer, history model,
  duplicate candidate/state-transition owner, profile, or group.
- Somon crawling/parser expansion, Telegram UI work, or production release work.
- Candidate selection or state writes in transport, persistence, composition,
  or shared-data owners.

## Cohesive implementation strategy

Primary owner is Polling Application at `internal/app`. It retains the fresh-ID
snapshot, orders exact evaluation before fallback, accounts for the single
detail-request cap, and alone requests every seen transition. Filtering at
`internal/filter` may expose a pure value-level way to evaluate the existing
non-price-max rules. Somon, Persistence, and Telegram remain unchanged
providers of detail/error, seen storage, and delivery outcomes.

The implementation result and all tests form one T2 task. Exact priority,
fallback ordering/count, cap completeness, and seen/delivery behavior are
coupled state transitions in `processNewCards`; none is a safe independently
releasable intermediate outcome. Tests, RED/GREEN proof, and quality review do
not create separate tasks.

## Dependencies and waves

- Foundation: `not_required`; no FT-000 dependency exists.
- Product-task dependencies: none.
- Wave: W1.
- Task: `TASK-001-T2-FT-001-W1`.

## Expected advisory change surface

- `internal/app/app.go` — phase and order fresh candidate processing while
  preserving cap, error, delivery, and seen transitions.
- `internal/app/price_fallback_test.go` — focused httptest/temp-SQLite polling
  scenarios and captured Telegram output.
- `internal/filter/settings.go` — only if a small pure match variant/helper is
  required for the non-maximum-price predicate.
- `internal/filter/settings_test.go` — boundary proof for that helper if added.

The executor may choose an immaterial test filename or a smaller same-outcome
surface after preflight. No hard write allow-list is imposed.

## Accepted boundaries and invariants

- Owner/topology: [Polling orchestration contract](../../contracts/boundary-map.md#polling-orchestration-contract), [Filtering contract](../../contracts/boundary-map.md#filtering-contract), [Somon adapter contract](../../contracts/boundary-map.md#somon-adapter-contract), [Persistence contract](../../contracts/boundary-map.md#persistence-contract), and [Telegram application boundary](../../contracts/boundary-map.md#telegram-application-boundary).
- Architecture: [AD-001](../../architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable) and [AD-002](../../architecture/system-architecture.md#ad-002--polling-application-owns-candidate-and-seen-transitions).
- Lifecycle: [Active poll candidate phases](../../states/runtime-lifecycle.md#active-poll-candidate-phases).
- MUST/NEVER rules: [Accepted MUST](../../invariants.md#accepted-must) and [Accepted NEVER](../../invariants.md#accepted-never).

The Polling Application must not bypass `internal/store`, and Filtering must
remain deterministic and free of I/O, scheduling, delivery, or durable state.
Exact and fallback detail work share one cap. An exact match suppresses
fallback even after failed delivery; incomplete exact evaluation also
suppresses fallback. No failed delivery is marked seen.

## Sources and canonical SDD coverage

- Product authority: [PRD REQ-003](../../prd.md#req-003--bounded-higher-price-fallback), [PRD REQ-004](../../prd.md#req-004--freshness-deduplication-and-failures), and [Requirements](../../requirements.md).
- Exact claims: `.memory-bank/features/FT-001-price-fallback.md#FT-001-AC-001`
  through `.memory-bank/features/FT-001-price-fallback.md#FT-001-AC-005`.
- Verification policy: [Risk-based checks](../../testing/strategy.md#risk-based-checks) and [existing package harness](../../testing/current-coverage.md#automated-coverage-by-package).

All concrete concerns use the canonical links above (`reuse`). No contract,
state, data, testing, guide, runbook, or ADR extension/creation is required.
No behavior-spec example adds useful authority beyond the exact ACs and
deterministic scenarios, so none is created.

## Verification targets and evidence

Use focused `internal/app` integration tests with httptest Somon/Telegram
servers and temporary SQLite to distinguish:

- exact success and exact delivery failure from fallback output;
- absent `PriceMax` with no fallback pass or fallback-only delivery; when
  `PriceMax` is present, open lower and inclusive upper price boundaries,
  stable tie ordering, and the three-ad limit;
- existing seen IDs, unknown price, over-ceiling and other rejection, plus
  exact-empty no-output behavior;
- shared-cap exhaustion, transient detail failure, HTTP 403/429, and failed
  delivery with the required unseen/seen post-state.

Record the honest pre-change result for every claim with explicit `RED:` and
`GREEN:` observations and a decisive comparison. Preserve pre-implementation
GREEN where the baseline already satisfies a claim; use `RED_NOT_APPLICABLE`
instead of manufacturing a regression. Where fallback behavior is absent,
capture the real baseline RED and then claim-equivalent GREEN. For
FT-001-AC-005 / REQ-006, the pre/post repository diff, owner review, and native
build must prove no new schema, setting, dependency, worker, compatibility
layer, second history model, or duplicate candidate/state-transition owner.

Required gates:

- `CGO_ENABLED=1 go test -count=1 ./internal/app ./internal/filter`
- `./scripts/build.sh`

UAT is the deterministic captured poll output and temporary-SQLite post-state;
no live Somon, Telegram, or production probe is required for this feature.

## Constitution constraints

Apply schema-backed planning, minimal verifiable change, evidence before done,
and the project KISS gate. The task must reuse current owners and stop if
execution reveals a required public contract, schema/config/dependency, T3, or
unresolved product-behavior change.

## Completion route

After task implementation, `/verify TASK-001-T2-FT-001-W1` must pass. When all
FT-001 tasks are implemented, feature completion also requires
`/red-verify --feature FT-001`, followed by lifecycle-owner closure and the
wave/feature-boundary `/mb-sync`.
