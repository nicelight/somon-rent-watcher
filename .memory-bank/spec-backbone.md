---
description: Pre-PRD spec framing and global SDD backbone state.
status: active
last_updated: 2026-09-04
---
# SDD Spec Backbone

## Brownfield Current-State Baseline

- Mapping status: complete on 2026-09-01.
- Baseline kind: evidence-backed `as-is`; it is not a target architecture or Global Backbone decision.
- PRD status: `.memory-bank/prd.md` is clarified and Constitution-checked; EP-001 and FT-001..FT-003 are decomposed, while task records remain pending design/tasking.
- Primary routes:
  - [.memory-bank/product.md](product.md): accepted product scope with current-state evidence routes.
  - [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md): accepted C4 L1-L3 target with implementation alignment.
  - [.memory-bank/contracts/current-integrations.md](contracts/current-integrations.md): observed external/internal boundaries.
  - [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md): lifecycle and persistence current state.
  - [.memory-bank/runbooks/almalinux-9-operations.md](runbooks/almalinux-9-operations.md): operations routing.
  - [.memory-bank/runbooks/docker-local-operations.md](runbooks/docker-local-operations.md): local Docker operations routing.
  - [.memory-bank/testing/current-coverage.md](testing/current-coverage.md): proof paths and verification gaps.
- Verification update: local Docker formatting/tests/vet/CGO/linkage and live SQLite/Telegram/Somon doctor passed on 2026-09-01. Target AlmaLinux production build, doctor, systemd startup, fresh paused baseline and host-isolation checks passed on 2026-09-02 for commit `7f9c5f50d659`; Git history before the 2026-09-01 initialization is unavailable.
- Downstream rule: use this baseline as evidence after the operator supplies product intent/PRD/delta; do not infer accepted target decisions from it.

## Pre-PRD Spec Status
- Status: ready_for_prd
- Last updated: 2026-09-04
- Notes: Clarified PRD plus the evidence-backed brownfield baseline supported completed L1-L3 decomposition. Global architecture and Foundation decisions remain pending for `/spec-design` after feature-plan approval.

## Decomposition Inputs
- User scenarios: [.memory-bank/prd.md#ux--interaction-flow](prd.md#ux--interaction-flow) and [#users--actors](prd.md#users--actors) define the admin, group, polling and delivery scenarios; no separate scenario artifact is needed for decomposition.
- Domain model: [.memory-bank/prd.md#data--domain-model](prd.md#data--domain-model) preserves `Card`, `Ad`, `Settings`, `seen_ads`, and per-poll candidate classes.
- Constraints: [.memory-bank/prd.md#functional-requirements](prd.md#functional-requirements) and [#non-functional-requirements](prd.md#non-functional-requirements) fix KISS, request-cap, retry and production-isolation constraints.
- Non-goals: [.memory-bank/prd.md#non-goals](prd.md#non-goals) excludes schema/history, configurable fallback machinery, client detection and unrelated production changes.
- Risks: [.memory-bank/prd.md#edge-cases--failure-handling](prd.md#edge-cases--failure-handling) covers incomplete exact evaluation, delivery ambiguity, backoff and stale menu behavior.
- Boundary hints: extend the existing `internal/app` polling owner for exact/fallback orchestration and the existing `internal/telegram` owner for append-only callback output; preserve `internal/store`, Somon and deployment boundaries pending `/spec-design` acceptance.
- Lifecycle hints: preserve first-seen and delivery retry lifecycle while adding only an in-memory per-poll exact-before-fallback phase; manual scan keeps single-flight scheduler ownership.

## Open Design Questions

- None that block product decomposition. `/spec-design` must determine the minimum accepted architecture/contracts and truthful Foundation path without inventing a new runtime or storage layer.

## Backbone Area Matrix
| Area | Status | Authoritative source | Notes |
|---|---|---|---|
| architecture_style | authoritative | [.memory-bank/architecture/system-architecture.md#accepted-shape](architecture/system-architecture.md#accepted-shape) | Preserve one Go modular monolith and existing package owners under KISS. |
| source_of_truth | authoritative | [.memory-bank/architecture/system-architecture.md#durable-ownership](architecture/system-architecture.md#durable-ownership) | PRD/Memory Bank own target; SQLite remains sole durable runtime store. |
| module_boundaries | authoritative | [.memory-bank/contracts/boundary-map.md#modules](contracts/boundary-map.md#modules) | Existing functional packages and exact allowed edges/contracts are accepted. |
| user_scenarios | authoritative | [.memory-bank/prd.md#ux--interaction-flow](prd.md#ux--interaction-flow) | Polling, group delivery and administrator flows are sufficient; no separate scenario artifact required. |
| constraints | authoritative | [.memory-bank/invariants.md](invariants.md) | Accepted MUST/NEVER and compatibility guardrails. |
| non_goals | authoritative | [.memory-bank/prd.md#non-goals](prd.md#non-goals) | No storage/config/dependency/client-detection or unrelated infrastructure expansion. |
| domain_model | authoritative | [.memory-bank/prd.md#data--domain-model](prd.md#data--domain-model) | Existing Card/Ad/Settings/seen model plus ephemeral poll candidate classes is sufficient. |
| data_flow | authoritative | [.memory-bank/architecture/system-architecture.md#main-target-data-flow](architecture/system-architecture.md#main-target-data-flow) | Exact-before-fallback and append-only callback flows have one orchestration owner each. |
| storage | authoritative | [.memory-bank/states/runtime-lifecycle.md#persisted-current-state](states/runtime-lifecycle.md#persisted-current-state) | Existing schema/write ownership is preserved; no migration. |
| api_contracts | authoritative | [.memory-bank/contracts/boundary-map.md#external-boundary-rules](contracts/boundary-map.md#external-boundary-rules) | Existing Somon/Telegram outbound contracts are retained and callback UI becomes append-only. |
| event_message_contracts | not_applicable | [.memory-bank/architecture/system-architecture.md#accepted-shape](architecture/system-architecture.md#accepted-shape) | No event bus, queue or asynchronous message envelope exists or is accepted. |
| agent_io_contracts | not_applicable | [.memory-bank/product.md#non-goals](product.md#non-goals) | Runtime contains no AI agent/tool protocol boundary. |
| security_safety | authoritative | [.memory-bank/invariants.md#accepted-never](invariants.md#accepted-never) | Secrets, block handling and shared-host isolation remain bounded. |
| deployment | authoritative | [.memory-bank/architecture/system-architecture.md#deployment-boundary](architecture/system-architecture.md#deployment-boundary) | Ordered, runtime-detected, scoped update with state preservation and postflight. |
| risks | authoritative | [.memory-bank/prd.md#edge-cases--failure-handling](prd.md#edge-cases--failure-handling) | Cap, delivery ambiguity, backoff, stale menu and runtime ambiguity have accepted outcomes. |
| open_questions | authoritative | [.memory-bank/spec-backbone.md#open-design-questions](spec-backbone.md#open-design-questions) | No unresolved design question blocks feature design/tasking. |

## Handoff To /prd-to-features
- Ready: yes
- Required reads: `.memory-bank/prd.md`, `.memory-bank/spec-index.md`, this file, and linked brownfield evidence.
- Stop conditions: Pre-PRD Spec Status becomes stale/blocked or PRD clarification/Constitution checks no longer pass.

## Handoff To /spec-design
- Global Backbone Status: complete under `/spec-design`
- Downstream readiness: feature-local `/spec-auto FT-<NNN>` and `/feature-to-tasks FT-<NNN>` are allowed; product scheduler waits for reviewed tasks.
- Backbone areas to revisit: none before current feature design.
- Candidate specs: extend existing registered owners only when feature design proves concrete missing detail.

## Global Backbone Status
- Status: complete
- Planning Revision: 1
- Mode: standard_architecture_scaffold
- Architecture artifact strategy: split-core-docs
- Not applicable areas:
  - event_message_contracts: not_applicable - no event bus, queue or message-envelope boundary exists in the accepted one-process design.
  - agent_io_contracts: not_applicable - the product runtime has no AI agent/tool protocol surface.
- Notes: Existing architecture, boundary, lifecycle, invariant, testing and runbook owners form the minimum production-sensitive scaffold. Foundation is not required because the executable/test/storage/runtime baseline is already proven.
