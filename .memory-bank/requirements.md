---
description: Accepted requirements and traceability for the current Somon Rent Watcher delta.
status: active
last_updated: 2026-09-09
source_of_truth:
  - .memory-bank/prd.md
---

# Requirements

## Functional requirements

### REQ-001 — Preserve monitoring semantics

Keep baseline, first-seen, pause, recovery, backoff, shared detail cap, and confirmed-delivery-before-seen behavior except for the explicit in-poll fallback classification.

### REQ-002 — Append-only Telegram controls

Acknowledge authorized callbacks promptly and send the resulting menu/state as a new message without runtime `editMessageText`; text input and manual scan also complete through fresh messages.

### REQ-003 — Bounded price fallback

Only after a completed poll finds no fresh exact match, send at most three fresh ads satisfying all other filters with known price in `(PriceMax, floor(1.5 * PriceMax)]`, ordered by price with stable feed-order ties; no `PriceMax` means no fallback.

### REQ-004 — Freshness and failure behavior

Use existing `seen_ads` as the sole freshness authority, ignore all existing seen IDs, mark fully rejected fresh cards seen under the current first-seen policy, preserve unseen retry after ambiguous/failed delivery, and suppress fallback when an exact match exists or exact evaluation is incomplete.

### REQ-005 — Ordered release

Run local gates before GitHub publication, then current production preflight, production git synchronization, target build, scoped runtime update, and postflight while preserving settings and SQLite state.

### REQ-009 — Correct Somon price extraction

Operator-authorized production correction on 2026-09-09: extract the actual price
from the ad's price data without adjacent photo counts or address numbers. Preserve
supported explicit Цена: price-line parsing. Deploy under REQ-005/REQ-008, preserving
settings and seen semantics. Source: [accepted follow-up](prd.md#accepted-production-hotfix-2026-09-09).

## Quality requirements

### REQ-006 — KISS maintainability (PRD NFR-001)

Reuse current settings, storage, scheduler, clients, and limits with no new dependency, database schema, configuration option, worker, or compatibility layer. Pass/fail changes when the delta introduces one of those concepts or duplicates an existing owner. Verification: code/spec review and project-native build gate.

### REQ-007 — Filtering and delivery reliability (PRD NFR-002)

Fallback must never hide an exact match, exceed three ads or the shared detail cap, exceed the inclusive 150% ceiling, or mark a failed delivery seen. Pass/fail changes with mixed candidates, cap exhaustion, boundary prices, or Telegram failure. Verification: deterministic unit/integration tests.

### REQ-008 — Production isolation (PRD NFR-003)

Only Somon Rent Watcher may change during deployment; unrelated workload/service/network/security/database state must remain unchanged. Pass/fail changes with runtime detection, scoped restart/install commands, or before/after host-state differences. Verification: operator-authorized read-only preflight/postflight and scoped health checks.

## Out of scope

- Separate shown-history or backfill of existing `seen` IDs.
- Configurable fallback count/percentage and client-specific Telegram handling.
- Broader scraper, storage, product-profile, or production-infrastructure changes.

## Requirements Traceability Matrix

| Requirement | Epic | Feature | Acceptance / proof route | Lifecycle |
|---|---|---|---|---|
| REQ-001 | EP-001 | FT-001, FT-002 | FT-001-AC-001, FT-001-AC-004, FT-002-AC-003 | planned |
| REQ-002 | EP-001 | FT-002 | FT-002-AC-001, FT-002-AC-002, FT-002-AC-003 | planned |
| REQ-003 | EP-001 | FT-001 | FT-001-AC-001, FT-001-AC-002, FT-001-AC-003, FT-001-AC-004 | planned |
| REQ-004 | EP-001 | FT-001 | FT-001-AC-001, FT-001-AC-003, FT-001-AC-004 | planned |
| REQ-005 | EP-001 | FT-003, FT-004 | FT-003-AC-001, FT-003-AC-002; FT-004-AC-002 (hotfix done) | planned |
| REQ-006 | EP-001 | FT-001, FT-002 | FT-001-AC-005, FT-002-AC-004 | planned |
| REQ-007 | EP-001 | FT-001 | FT-001-AC-001, FT-001-AC-002, FT-001-AC-003, FT-001-AC-004 | planned |
| REQ-008 | EP-001 | FT-003, FT-004 | FT-003-AC-001, FT-003-AC-002; FT-004-AC-002 (hotfix done) | planned |
| REQ-009 | EP-001 | FT-004 | FT-004-AC-001; TASK-005-T3-FT-004-W1 | done |
