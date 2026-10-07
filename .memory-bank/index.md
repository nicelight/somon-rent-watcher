---
description: Main project knowledge map for accepted product intent, planning, and brownfield evidence.
status: active
last_verified: 2026-09-02
---

# Memory Bank Index

## Product and requirements

- [Keyword monitoring proposal](prd.md#keyword-monitoring-proposal--2026-10-06): принятые независимые поиски, подтверждённые scoped URLs и минимальный дизайн; FT-005 реализована и проверена, включая detail validation и bounded history lookup.

- [.memory-bank/analysis/product-brief.md](analysis/product-brief.md): accepted concise product input.
- [.memory-bank/prd.md](prd.md): clarified, Constitution-checked product requirements.
- [.memory-bank/product.md](product.md): product identity, value, primary flow, constraints, and non-goals.
- [.memory-bank/requirements.md](requirements.md): REQ registry and traceability matrix.
- [.memory-bank/epics/EP-001-reliable-watcher-upgrade.md](epics/EP-001-reliable-watcher-upgrade.md): current upgrade epic.
- [.memory-bank/features/FT-001-price-fallback.md](features/FT-001-price-fallback.md): bounded closest-price alternatives.
- [.memory-bank/features/FT-002-append-only-telegram-ui.md](features/FT-002-append-only-telegram-ui.md): client-visible append-only bot administration.
- [.memory-bank/features/FT-003-isolated-production-release.md](features/FT-003-isolated-production-release.md): ordered state-preserving production delivery.
- [.memory-bank/features/FT-004-price-extraction-hotfix.md](features/FT-004-price-extraction-hotfix.md): completed price parser correction and isolated production hotfix.
- [.memory-bank/epics/EP-002-keyword-monitoring.md](epics/EP-002-keyword-monitoring.md): принятая новая возможность независимых поисков, REQ-010…014.
- [.memory-bank/features/FT-005-keyword-monitoring.md](features/FT-005-keyword-monitoring.md): создание, настройка, удаление и доставка keyword-поисков; native scope подтверждён; AC001…009 verified.
- [.memory-bank/tasks/plans/IMPL-FT-005.md](tasks/plans/IMPL-FT-005.md): семь выполненных задач FT-005 с AC proof, owners и Docker gates.

## Brownfield current-state baseline

- [.memory-bank/product.md](product.md): accepted product scope with links to current-state evidence.
- [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md): C4 context/runtime/component map, entrypoints, data flow and writers.
- [.memory-bank/contracts/current-integrations.md](contracts/current-integrations.md): observed external contracts and internal dependency evidence; non-authoritative for target design. Keyword source API/catalog/fixture limits: [implementation routing](contracts/current-integrations.md#keyword-source-implementation-routing).
- [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md): current polling, delivery, recovery and persisted-state lifecycle.
- [.memory-bank/runbooks/almalinux-9-operations.md](runbooks/almalinux-9-operations.md): routing to the production operations procedure.
- [.memory-bank/runbooks/docker-local-operations.md](runbooks/docker-local-operations.md): local Kubuntu Docker Compose build, runtime and state routing.
- [.memory-bank/guides/local-development.md](guides/local-development.md): current Go/CGO build and verification HOW.
- [.memory-bank/testing/current-coverage.md](testing/current-coverage.md): automated tests, fixtures, live checks and unresolved verification.
- [.memory-bank/glossary.md](glossary.md): shared current and accepted target vocabulary.
- [.memory-bank/invariants.md](invariants.md): accepted-invariant status and routing to descriptive guardrails.
- [.memory-bank/changelog.md](changelog.md): durable log of synchronized implementation/documentation waves.

Baseline scope and remaining gaps are summarized in [.memory-bank/spec-backbone.md#brownfield-current-state-baseline](spec-backbone.md#brownfield-current-state-baseline). The clarified PRD and product decomposition now define target intent; accepted architecture and task records remain pending downstream design.

## Governing and workflow navigation

- [.memory-bank/constitution.md](constitution.md): Project Constitution — top governing policy for agents.
- [.memory-bank/mbb/index.md](mbb/index.md): Memory Bank rules and SSOT conventions.
- [.memory-bank/spec-index.md](spec-index.md): SDD registry plus explicit non-normative baseline registry.
- [.memory-bank/spec-backbone.md](spec-backbone.md): pre-PRD/Global Backbone status and brownfield handoff.
- [.memory-bank/foundation.md](foundation.md): executable baseline sufficiency and Foundation Gate anchors.
- [.memory-bank/requirements.md](requirements.md): accepted REQ registry and RTM.
- [.memory-bank/contracts/boundary-map.md](contracts/boundary-map.md): canonical accepted target graph; currently empty pending SDD design.
- [.memory-bank/workflows/index.md](workflows/index.md): workflow router and shared SDD/execution policies.
- [.memory-bank/testing/index.md](testing/index.md): testing documentation router.
- [.memory-bank/skills/index.md](skills/index.md): installed project skill registry.

## Roles

- [.memory-bank/roles/general.md](roles/general.md): General role contract for one-agent execution.
- [.memory-bank/roles/orchestrator.md](roles/orchestrator.md): Orchestrator role contract.
- [.memory-bank/roles/architect.md](roles/architect.md): Architect role contract.
- [.memory-bank/roles/explorer.md](roles/explorer.md): Explorer role contract.
- [.memory-bank/roles/implementer.md](roles/implementer.md): Implementer role contract.
- [.memory-bank/roles/reviewer.md](roles/reviewer.md): Reviewer role contract.
- [.memory-bank/roles/judge.md](roles/judge.md): Judge supervisory role contract.

## Framework planning stores

- [.memory-bank/prd.md](prd.md): clarified source for decomposition and design.
- [.memory-bank/epics/](epics/): product epic records.
- [.memory-bank/features/](features/): product feature records with stable AC IDs.
- [.memory-bank/tasks/index.json](tasks/index.json): task registry; task records remain pending feature design/tasking.
- [.memory-bank/schemas/task.schema.json](schemas/task.schema.json): JSON schema for future task records.
- [.memory-bank/behavior-specs/](behavior-specs/): optional behavior examples when later linked by product workflow.
