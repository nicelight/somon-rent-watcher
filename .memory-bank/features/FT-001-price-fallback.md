---
description: Show up to three closest fresh higher-price ads only when exact matching is empty.
status: draft
lifecycle: planned
spec_design_status: complete
spec_design_links:
  - ".memory-bank/architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable"
  - ".memory-bank/architecture/system-architecture.md#ad-002--polling-application-owns-candidate-and-seen-transitions"
  - ".memory-bank/contracts/boundary-map.md#polling-orchestration-contract"
  - ".memory-bank/contracts/boundary-map.md#filtering-contract"
  - ".memory-bank/contracts/boundary-map.md#somon-adapter-contract"
  - ".memory-bank/contracts/boundary-map.md#persistence-contract"
  - ".memory-bank/contracts/boundary-map.md#telegram-application-boundary"
  - ".memory-bank/states/runtime-lifecycle.md#active-poll-candidate-phases"
  - ".memory-bank/testing/strategy.md#risk-based-checks"
  - ".memory-bank/testing/current-coverage.md#automated-coverage-by-package"
---

# FT-001 — Closest higher-price alternatives

## Value and use cases

When a completed active poll has no fresh exact match, group members can still see a very small set of relevant alternatives without manually widening the filter. Exact results remain primary and existing first-seen behavior prevents repeat spam.

## Source and constraints

- [.memory-bank/prd.md#req-003--bounded-higher-price-fallback](../prd.md#req-003--bounded-higher-price-fallback)
- [.memory-bank/prd.md#req-004--freshness-deduplication-and-failures](../prd.md#req-004--freshness-deduplication-and-failures)
- Existing request cap, backoff, sequential fetch, SQLite and confirmed-delivery semantics remain mandatory.

## Edge and failure behavior

- Missing `PriceMax`, unknown price, price above the ceiling, or any other failed filter excludes fallback.
- Exact match or incomplete exact evaluation suppresses fallback.
- Delivery failure leaves the affected ID unseen and does not authorize substitution.
- Existing seen IDs never re-enter the candidate pool.

## Acceptance Criteria

### FT-001-AC-001 — Exact results suppress fallback

- REQ: REQ-001, REQ-003, REQ-004, REQ-007
- Observable criterion: if at least one fresh ad passes the complete exact filter, no higher-price ad is sent in that poll, including when exact delivery fails.
- Verification method: app integration test with mixed exact/fallback cards and Telegram success/failure responses.

### FT-001-AC-002 — Select the three closest eligible prices

- REQ: REQ-003, REQ-007
- Observable criterion: without `PriceMax` fallback sends nothing; otherwise an exact-empty completed poll sends at most three otherwise matching fresh ads with `PriceMax < price <= floor(1.5 * PriceMax)`, ascending by price with stable feed-order ties.
- Verification method: deterministic app integration test covering absent `PriceMax`, boundary prices, ordering, tie order, and fourth-candidate exclusion.

### FT-001-AC-003 — Seen and ineligible ads stay silent

- REQ: REQ-003, REQ-004, REQ-007
- Observable criterion: existing seen IDs are ignored; unknown-price, over-ceiling, and other fully rejected fresh candidates are not sent and become seen; no eligible fresh candidate means no ad output.
- Verification method: temporary-SQLite integration test asserting post-poll seen state plus captured Telegram requests.

### FT-001-AC-004 — Preserve cap, backoff, and retry semantics

- REQ: REQ-001, REQ-003, REQ-004, REQ-007
- Observable criterion: exact and fallback details share the existing request cap; cap exhaustion, unresolved detail error, or Telegram delivery failure suppresses fallback as applicable and leaves deferred/retryable/failed-delivery IDs unseen; 403/429 still enters existing backoff.
- Verification method: app tests with a low detail cap and typed detail/Telegram failures.

### FT-001-AC-005 — No new durable mechanism

- REQ: REQ-006
- Observable criterion: the implementation adds no schema, setting, dependency, worker, compatibility layer, second history model, or duplicate candidate/state-transition owner.
- Verification method: code/spec review and repository-native build gate.

## Acceptance closure

- Exact priority, price boundaries/order/count, freshness exclusions, cap/backoff, delivery ambiguity, and no-result silence are covered by FT-001-AC-001 through FT-001-AC-004.
- KISS/no-migration quality is covered by FT-001-AC-005.

## SDD Design

Feature design is complete by reusing the accepted owners; no feature-local
technical spec, ADR, schema, setting, dependency, worker, or new status model is
required.

- Architecture Spine: [AD-001](../architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable) and [AD-002](../architecture/system-architecture.md#ad-002--polling-application-owns-candidate-and-seen-transitions).
- Boundaries: [polling orchestration](../contracts/boundary-map.md#polling-orchestration-contract), [filtering](../contracts/boundary-map.md#filtering-contract), [Somon](../contracts/boundary-map.md#somon-adapter-contract), [persistence](../contracts/boundary-map.md#persistence-contract), and [Telegram delivery](../contracts/boundary-map.md#telegram-application-boundary).
- State and data: [active poll candidate phases](../states/runtime-lifecycle.md#active-poll-candidate-phases) and [accepted freshness/write invariants](../invariants.md#accepted-must); existing SQLite state remains unchanged.
- Verification: the AC methods above use the [risk-based testing policy](../testing/strategy.md#risk-based-checks) and extend the existing [app/filter/store/Telegram harness](../testing/current-coverage.md#automated-coverage-by-package).

Global Backbone Planning Revision 1 is complete and Foundation is
`not_required`.
