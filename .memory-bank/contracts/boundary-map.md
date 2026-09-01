---
description: Canonical accepted module/change-unit dependency graph and boundary contracts.
status: draft
last_verified: 2026-09-01
---

# Boundary Map

## Purpose

- Keep one accepted inventory of project modules/change units and every allowed significant dependency between them.
- Treat `Consumer -> Provider` as the direction of dependency. Observed imports or calls are evidence, not accepted edges by themselves.

## Authority status

Accepted Global Backbone ещё не создан. Brownfield-маппинг не зарегистрировал observed code packages как accepted target modules/edges.

Текущие imports/calls, external integrations и runtime boundaries описаны отдельно: [.memory-bank/contracts/current-integrations.md](current-integrations.md): non-authoritative as-is evidence for future design work.

## Modules

| Module / Change Unit | Parent Architecture Unit | Code Root | Responsibility |
|---|---|---|---|

## Dependency Graph

`Consumer -> Provider` means Consumer depends on Provider through the linked contract.

| Consumer | Provider | Contract |
|---|---|---|

## Inline Contracts

Добавлять contract blocks сюда можно только после принятия corresponding module/change-unit graph владельцем SDD workflow. Current-state evidence само по себе не авторизует edge.

## Update Rules

- `Module / Change Unit` is the unique graph key. Use stable functional responsibility names, not feature/task IDs, current paths, or generic technical layers.
- Every graph row names registered modules and links to one exact contract heading. The graph row alone owns consumer, provider, and direction.
- Include every accepted significant inter-module dependency. An absent edge is not authorized.
- Keep the detailed accepted module inventory here. `system-architecture.md` owns larger architecture units and current-state topology.
- Plans and tasks link relevant graph/contract blocks through existing fields; they do not copy the subgraph or introduce graph-specific task fields.
