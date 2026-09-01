---
description: Pre-PRD spec framing and global SDD backbone state.
status: active
last_updated: 2026-09-01
---
# SDD Spec Backbone

## Brownfield Current-State Baseline

- Mapping status: complete on 2026-09-01.
- Baseline kind: evidence-backed `as-is`; it is not a target architecture or Global Backbone decision.
- PRD status: `.memory-bank/prd.md` absent; roadmap entities and task records were not generated.
- Primary routes:
  - [.memory-bank/product.md](product.md): C4 L1 product current state.
  - [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md): C4 L1-L3 runtime/component current state.
  - [.memory-bank/contracts/current-integrations.md](contracts/current-integrations.md): observed external/internal boundaries.
  - [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md): lifecycle and persistence current state.
  - [.memory-bank/runbooks/almalinux-9-operations.md](runbooks/almalinux-9-operations.md): operations routing.
  - [.memory-bank/testing/current-coverage.md](testing/current-coverage.md): proof paths and verification gaps.
- Needs verification: Go checks, live Somon/Telegram and target-host state were not available in this workspace; Git history before the 2026-09-01 initialization is unavailable.
- Downstream rule: use this baseline as evidence after the operator supplies product intent/PRD/delta; do not infer accepted target decisions from it.

## Pre-PRD Spec Status
- Status: blocked
- Last updated: 2026-09-01
- Notes: Brownfield baseline exists, but authoritative PRD does not. Run /spec-init after /write-prd to determine whether PRD decomposition is safe.

## Decomposition Inputs
- User scenarios: not_started
- Domain model: not_started
- Constraints: not_started
- Non-goals: not_started
- Risks: not_started
- Boundary hints: not_started
- Lifecycle hints: not_started

## Open Design Questions
- TBD

## Backbone Area Matrix
| Area | Status | Authoritative source | Notes |
|---|---|---|---|
| architecture_style | blocked | - | Decide in /spec-design after /prd-to-features. |
| source_of_truth | blocked | - | Decide in /spec-design after /prd-to-features. |
| module_boundaries | blocked | .memory-bank/contracts/boundary-map.md | Accept parent architecture units in /spec-design; reconcile concrete modules, edges, and contracts in /feature-to-tasks. |
| user_scenarios | blocked | .memory-bank/user-scenarios.md | Create/review when scenarios affect decomposition or architecture. |
| constraints | blocked | - | Capture in /spec-init and refine in /spec-design. |
| non_goals | blocked | - | Capture in /spec-init and refine in /spec-design. |
| domain_model | blocked | .memory-bank/domains/core-domain.md | Create only when domain model affects decomposition or shared design. |
| data_flow | blocked | - | Decide in /spec-design after /prd-to-features. |
| storage | blocked | - | Decide in /spec-design after /prd-to-features. |
| api_contracts | blocked | - | Decide authoritative/needed/not_applicable/blocked in /spec-design. |
| event_message_contracts | blocked | - | Decide authoritative/needed/not_applicable/blocked in /spec-design. |
| agent_io_contracts | blocked | - | Decide authoritative/needed/not_applicable/blocked in /spec-design. |
| security_safety | blocked | - | Decide in /spec-design after /prd-to-features. |
| deployment | blocked | - | Decide in /spec-design after /prd-to-features. |
| risks | blocked | - | Capture in /spec-init and refine in /spec-design. |
| open_questions | blocked | - | Resolve or keep blocked. |

## Handoff To /prd-to-features
- Ready: no
- Required reads: .memory-bank/prd.md, .memory-bank/spec-index.md, this file, and linked pre-PRD specs.
- Stop conditions: Pre-PRD Spec Status is missing, stale, or blocked.

## Handoff To /spec-design
- Global Backbone Status: intentionally pending until /spec-design
- Downstream readiness: /feature-to-tasks, /autopilot, and autonomous scheduler mode must wait for /spec-design.
- Backbone areas to revisit: all
- Candidate specs: see .memory-bank/spec-index.md Planned Specs.

## Global Backbone Status
- Status: blocked
- Planning Revision: 0
- Mode: pending
- Architecture artifact strategy: pending
- Not applicable areas:
  - TBD
- Notes: /spec-design has not completed the global architecture scaffold yet.
