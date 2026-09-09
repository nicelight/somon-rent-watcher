---
description: Pure SDD spec registry and planned-spec index.
status: active
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/spec-index.md
---

# SDD Spec Index

## Purpose

- Register canonical current/target specs without duplicating their decisions.
- Keep readiness, matrix, open questions and handoffs in [spec-backbone.md](spec-backbone.md).
- Feature `spec_design_status` and exact applicable links live in feature documents.
- Accepted backbone/shared-contract changes route to `/spec-redesign`; feature-local extensions route through `/feature-to-tasks` or `/spec-auto`.

## Current-State Baseline Registry

These artifacts remain descriptive evidence and do not override accepted target specs.

| Type | Path | Status | Current-state scope |
|---|---|---|---|
| integration baseline | [.memory-bank/contracts/current-integrations.md](contracts/current-integrations.md) | active as-is | Observed external boundaries and package dependencies. |
| operations route | [.memory-bank/runbooks/docker-local-operations.md](runbooks/docker-local-operations.md) | active as-is | Local Docker build/runtime and persistent-state route. |
| testing baseline | [.memory-bank/testing/current-coverage.md](testing/current-coverage.md) | active as-is | Existing tests, fixtures, live checks and gaps. |
| development guide | [.memory-bank/guides/local-development.md](guides/local-development.md) | active as-is | Build prerequisites and native gate. |

## Spec Registry

| Type | Path | Status | Scope | Change route |
|---|---|---|---|---|
| governance | [.memory-bank/constitution.md](constitution.md) | active | Top governing policy. | /constitution |
| architecture | [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md) | active | Accepted system shape, owners, data flow, deployment boundary and Architecture Spine. | initial /spec-design; post-acceptance /spec-redesign |
| invariants | [.memory-bank/invariants.md](invariants.md) | active | Accepted global MUST/NEVER and compatibility guardrails. | initial /spec-design; post-acceptance /spec-redesign |
| glossary | [.memory-bank/glossary.md](glossary.md) | active | Shared current and accepted target vocabulary. | /brief, /spec-init, or initial /spec-design; post-acceptance /spec-redesign |
| contract | [.memory-bank/contracts/boundary-map.md](contracts/boundary-map.md) | active | Canonical module identities, allowed dependency graph and boundary contracts. | initial /spec-design, post-acceptance /spec-redesign, or /feature-to-tasks |
| state | [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md) | active | Current lifecycle plus accepted exact/fallback and append-only target transitions. | initial /spec-design; post-acceptance /spec-redesign; feature-local /feature-to-tasks |
| foundation | [.memory-bank/foundation.md](foundation.md) | active | Existing executable baseline sufficiency and Foundation Gate anchors. | initial /spec-design; /foundation-to-tasks only when required |
| testing | [.memory-bank/testing/strategy.md](testing/strategy.md) | active | Risk-based testing and evidence policy. | explicit project-level user decision |
| runbook | [.memory-bank/runbooks/almalinux-9-operations.md](runbooks/almalinux-9-operations.md) | active | Current AlmaLinux operations evidence plus the accepted FT-003 isolated release sequence and proof boundary. | /spec-auto or /feature-to-tasks for feature-local detail; /spec-redesign for shared deployment changes |

## Planned Specs

| Area | Expected path | Needed by | Notes |
|---|---|---|---|
| interface_contract_specs | existing `.memory-bank/contracts/*`, `.memory-bank/testing/*`, and `.memory-bank/runbooks/*` | /feature-to-tasks | Extend only when one feature proves concrete missing boundary/verification detail. |
| data_specs | existing `.memory-bank/states/*` | /feature-to-tasks | Extend only when feature design proves current lifecycle/storage contract insufficient. |
| subject_feature_concerns | existing registered architecture/contract/state/testing/runbook/guide owners | /spec-auto, /feature-to-tasks | Reuse first; create no feature technical-spec hub. |

## Broken / Missing Links

- None.

## Update Rules

- Keep this file to registry/index content: types, paths, statuses, scopes, change routes and broken links.
- Canonical identity is the path; do not add reverse-usage or feature ownership copies.
- Do not add backbone matrices, decision bodies, feature status maps, state machines, schemas or open-question dumps here.
- Use registered specs for detailed decisions and [spec-backbone.md](spec-backbone.md) for readiness/routing.
