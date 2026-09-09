---
description: Completed operator-authorized correction of inflated Somon prices and isolated production hotfix.
status: active
lifecycle: verified
clarification_status: complete
last_clarified: 2026-09-09
clarification_questions: 0
last_updated: 2026-09-09
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
spec_design_status: complete
spec_design_links:
  - ".memory-bank/contracts/boundary-map.md#somon-adapter-contract"
  - ".memory-bank/contracts/current-integrations.md#somon-html-over-https"
  - ".memory-bank/states/runtime-lifecycle.md#ad-id-lifecycle"
  - ".memory-bank/runbooks/almalinux-9-operations.md#safety-boundaries"
---
# FT-004 — Correct Somon price extraction

## Clarifications

Clarification: no critical ambiguity found. The operator explicitly requested the
proven parser correction and careful deployment; no new product behavior was selected.

## Accepted scope and source

The operator reported missing listings, requested verification against the live site,
and explicitly authorized correcting the proven price parser defect and carefully
deploying only the watcher on the shared host on 2026-09-09. This document records that
completed maintenance change; it does not claim a retrospective task-plan approval or
completion of pending FT-001..003 work.

- [REQ-009](../requirements.md#req-009--correct-somon-price-extraction): correct prices.
- [Production diagnosis](../../.protocols/diagnostics/production-search-audit-2026-09-09/report.md): observed 5000→65000 and4800→314800 corruption.
- Preserve the existing parser owner/API, explicit Цена: fixture compatibility,
  settings and seen history; no history replay or other product changes.

## Acceptance Criteria

### FT-004-AC-001 — Read the actual price without adjacent numbers

- REQ: REQ-009
- Outcome: category photo counts and detail address numbers do not enter price;
  the existing explicit Цена: price line still parses correctly.
- Proof: regression RED/GREEN, full native gate, 60 category and38 cached detail prices.

### FT-004-AC-002 — Install the isolated hotfix without losing state

- REQ: REQ-005, REQ-008
- Outcome: only the watcher is updated after local/target gates and doctor; all old
  seen rows/settings/environment and unrelated host state are preserved; rollback
  binary is retained and the installed exact commit runs successfully.
- Proof: exact-commit release, backup-to-live subset check, PID/binary/source identity,
  normalized pre/post host comparisons and successful first-poll delivery.

## Completed task and evidence

- [TASK-005-T3-FT-004-W1](../tasks/TASK-005-T3-FT-004-W1.task.json): final93cbe9ebcfe6421040461f4ee4e049430bf5028d.
- [Functional verification](../../.protocols/TASK-005-T3-FT-004-W1/verification.md).
- [Independent semantic verification](../../.protocols/TASK-005-T3-FT-004-W1/red-verification.md).

Existing global architecture and source owners are reused without an architecture,
schema, setting, dependency or lifecycle change.
