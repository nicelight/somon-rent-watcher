---
description: Главная карта знаний проекта и brownfield current-state baseline для агентов.
status: active
last_verified: 2026-09-01
---

# Memory Bank Index

## Brownfield current-state baseline

- [.memory-bank/product.md](product.md): реализованный C4 L1 product scope, actors, flows and constraints.
- [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md): C4 context/runtime/component map, entrypoints, data flow and writers.
- [.memory-bank/contracts/current-integrations.md](contracts/current-integrations.md): observed external contracts and internal dependency evidence; non-authoritative for target design.
- [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md): current polling, delivery, recovery and persisted-state lifecycle.
- [.memory-bank/runbooks/almalinux-9-operations.md](runbooks/almalinux-9-operations.md): routing to the production operations procedure.
- [.memory-bank/runbooks/docker-local-operations.md](runbooks/docker-local-operations.md): local Kubuntu Docker Compose build, runtime and state routing.
- [.memory-bank/guides/local-development.md](guides/local-development.md): current Go/CGO build and verification HOW.
- [.memory-bank/testing/current-coverage.md](testing/current-coverage.md): automated tests, fixtures, live checks and unresolved verification.
- [.memory-bank/glossary.md](glossary.md): current-state vocabulary.
- [.memory-bank/invariants.md](invariants.md): accepted-invariant status and routing to descriptive guardrails.

Baseline scope and remaining gaps are summarized in [.memory-bank/spec-backbone.md#brownfield-current-state-baseline](spec-backbone.md#brownfield-current-state-baseline). No authoritative PRD exists, so brownfield mapping did not create roadmap/task entities or target architecture decisions.

## Governing and workflow navigation

- [.memory-bank/constitution.md](constitution.md): Project Constitution — top governing policy for agents.
- [.memory-bank/mbb/index.md](mbb/index.md): Memory Bank rules and SSOT conventions.
- [.memory-bank/spec-index.md](spec-index.md): SDD registry plus explicit non-normative baseline registry.
- [.memory-bank/spec-backbone.md](spec-backbone.md): pre-PRD/Global Backbone status and brownfield handoff.
- [.memory-bank/requirements.md](requirements.md): requirements authority status; formal REQ/RTM not initialized without PRD.
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

- `.memory-bank/prd.md`: absent; operator input is required before product decomposition.
- [.memory-bank/epics/](epics/): framework directory; no brownfield roadmap entities were generated.
- [.memory-bank/features/](features/): framework directory; no brownfield feature entities were generated.
- [.memory-bank/tasks/index.json](tasks/index.json): framework task registry; mapping did not add task links/records.
- [.memory-bank/schemas/task.schema.json](schemas/task.schema.json): JSON schema for future task records.
- [.memory-bank/behavior-specs/](behavior-specs/): optional behavior examples when later linked by product workflow.
