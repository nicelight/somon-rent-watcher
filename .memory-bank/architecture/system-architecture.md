---
description: Accepted system architecture with current implementation alignment for Somon Rent Watcher.
status: active
last_verified: 2026-09-02
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
  - .memory-bank/contracts/boundary-map.md
  - cmd/somonwatch/main.go
  - internal/app/app.go
  - deploy/somonwatch.service
  - compose.yaml
---

# System Architecture

## Accepted shape

Somon Rent Watcher remains one deployable Go modular monolith with one composition root, two cooperating runtime loops, one SQLite database, and outbound-only Somon/Telegram integrations. Existing package boundaries remain the implementation change units; the accepted delta does not add a service, worker, queue, schema, dependency, inbound API, or compatibility layer.

The architecture is intentionally conservative: product behavior changes inside the current polling and Telegram owners, while persistence and external adapter contracts stay stable.

## C4 L1 — System context

| Actor/system | Accepted relationship |
|---|---|
| Configured administrators | Manage one shared filter/status from private chat or the configured target group. |
| Target Telegram group | Receives exact and bounded fallback ad notifications. |
| Somon.tj | Supplies category/detail HTML through sequential outbound HTTPS. |
| Telegram Bot API | Supplies long-polled updates and accepts callback acknowledgements/messages/photos. |
| Local/production host | Runs exactly one watcher runtime and provides its private SQLite/debug storage. |

## C4 L2 — Runtime composition

| Runtime unit | Responsibility | State owner |
|---|---|---|
| `somonwatch` executable | CLI, wiring, lifecycle and diagnostics. | In-memory process plus SQLite-backed application state. |
| Telegram loop | Persisted-offset updates, authorization, append-only admin UI and outbound delivery. | Offset in SQLite; pending text input in memory. |
| Polling scheduler | Baseline/pause/active polls, exact-before-fallback orchestration, recovery, backoff and single-flight manual trigger. | Poll settings/snapshots in SQLite; active request in memory. |
| SQLite | `seen_ads`, singleton settings JSON and string state. | Sole durable application store. |
| Debug directory | Private bounded parser/sanity evidence. | Local filesystem, oldest files pruned. |

Production runtime identity is determined by read-only preflight immediately before deployment. Current verified evidence is a host `systemd` service; local Docker Compose remains a separate development route. No deployment step may operate on an unrelated runtime.

## C4 L3 — Accepted components

| Component | Code root | Responsibility |
|---|---|---|
| Composition | `cmd/somonwatch` | Config/adapters wiring, commands, signals, start/stop. |
| Configuration | `internal/config` | Environment defaults, parsing and validation. |
| Polling application | `internal/app` | Scheduler, poll orchestration, candidate phases, status and all seen/state transition decisions. |
| Filtering | `internal/filter` | Pure settings validation and deterministic card/ad matching. |
| Somon adapter | `internal/somon`, `internal/htmlx` | Rate-limited HTTP, HTML parsing and recovery URL/data acquisition. |
| Persistence adapter | `internal/store` | Serialized SQLite access, schema and transactional operations requested by the application. |
| Telegram adapter | `internal/telegram` | Bot API transport, authorization/UI handling, rendering and ad delivery. |
| Shared data | `internal/model` | Passive `Card`, `Ad`, and `RuntimeStatus` structures. |

Exact module identities and allowed dependencies are canonical in [Boundary Map](../contracts/boundary-map.md).

## Main target data flow

1. The scheduler fetches and sanity-checks category cards before durable advancement.
2. Baseline/pause behavior keeps existing first-seen semantics without group delivery.
3. Active polling queries current seen IDs and partitions fresh cards under the existing filters and request cap.
4. Exact candidates are fully evaluated first. A passing exact ad suppresses fallback even if its Telegram delivery fails.
5. Only after a complete exact-empty evaluation, eligible higher-price candidates are detail-filtered and up to three are delivered in price/feed order.
6. The polling application alone decides `MarkSeen`: confirmed deliveries and final rejections become seen; cap-deferred, transient-detail, and failed-delivery IDs remain unseen.
7. Telegram callbacks are acknowledged before their potentially slower action; the resulting current menu/status is sent as a new message.

## Durable ownership

- `internal/app` owns the meaning and order of `seen_ads` transitions but uses only `internal/store` operations.
- `internal/store` owns SQLite mechanics and file permissions; no other package opens or writes the DB.
- `internal/telegram` owns callback/UI/delivery protocol behavior but does not decide poll candidate lifecycle.
- `internal/filter` has no I/O or mutable global state.
- The composition root wires owners and does not host business orchestration.

## Deployment boundary

- Local proof uses the repository-native build gate or its documented Docker build equivalent.
- GitHub publication precedes production git synchronization.
- Production preflight must identify the watcher runtime and host health before any write.
- Only the watcher checkout/build/artifact/runtime may be updated; SQLite/settings remain in place.
- Version, doctor/runtime health, logs and unchanged unrelated host state are post-deploy evidence.

## Architecture Spine

### AD-001 — Preserve one cohesive deployable

- Binds: composition, all internal components, deployment.
- Prevents: services/workers/queues/DI registries introduced for this delta.
- Rule: implement accepted behavior inside existing package owners and keep one `somonwatch` runtime.
- Verification: dependency/schema/config diff review plus repository-native build.
- Source: REQ-006 and the accepted KISS instruction.

### AD-002 — Polling application owns candidate and seen transitions

- Binds: `internal/app`, filtering, Somon, persistence and Telegram delivery.
- Prevents: fallback hiding exact matches, premature seen writes, cap bypass and duplicated lifecycle ownership.
- Rule: exact evaluation precedes fallback; `internal/app` alone orders detail work and requests every seen transition.
- Verification: FT-001 integration tests and boundary review.
- Source: REQ-001, REQ-003, REQ-004 and REQ-007.

### AD-003 — Telegram administration is append-only

- Binds: Telegram callback, text-input and manual-scan flows.
- Prevents: unreliable client-visible dependence on message edits or client detection.
- Rule: acknowledge authorized callbacks promptly and represent the current result with a new message; runtime control flows do not edit old messages.
- Verification: FT-002 captured Bot API method/order tests.
- Source: REQ-002 and REQ-006.

### AD-004 — Deployment is runtime-detected and isolated

- Binds: publication, production checkout/build, runtime update and host evidence.
- Prevents: wrong-runtime updates and collateral changes on the shared host.
- Rule: stop on unhealthy/ambiguous preflight and update only the identified Somon Rent Watcher runtime in the accepted order.
- Verification: FT-003 ordered receipts and before/after host comparison.
- Source: REQ-005 and REQ-008.

## Current implementation alignment

The mapped code already implements the accepted one-process/package/storage shape and passed local Docker plus AlmaLinux build/doctor/runtime checks. The planned deltas are confined to polling orchestration, Telegram UI behavior, their tests/docs, and the scoped release path.

## Related evidence

- [Current integrations](../contracts/current-integrations.md)
- [Runtime lifecycle](../states/runtime-lifecycle.md)
- [Current coverage](../testing/current-coverage.md)
- [Local development](../guides/local-development.md)
- [AlmaLinux operations](../runbooks/almalinux-9-operations.md)

## Keyword monitoring design proposal

Принятая [возможность](../prd.md#accepted-decisions) добавляется в существующие
Go-процесс, SQLite и Telegram UI; аренда сохраняет свои фильтры и историю.

- `internal/app` управляет поисками, scheduler, кандидатами и историей.
- `internal/somon` получает выдачу; `internal/filter` проверяет условия поиска.
- `internal/store` хранит поиски и историю; `internal/telegram` управляет ими и доставляет.

Один scheduler последовательно обходит аренду и включённые поиски с общими интервалом,
HTTP delay, detail cap и backoff. Начало обхода меняется по кругу, чтобы один поиск
не забирал бюджет постоянно; ручной запуск сохраняет single-flight.

Поток: выдача → история конкретного поиска → фильтр → details → Telegram → запись.
Общий `seen_ads(ad_id)` остаётся арендным; независимые поиски используют
[добавочные записи](../states/runtime-lifecycle.md#keyword-monitoring-state-proposal).
Строгий бюджет, штатный matching и отсутствие повторов после снижения цены приняты.
Source flow: подтверждённый category/city path + `q`/`ordering=relevance` → primary
results → локальная цена. «Все категории» — полноценный native scope. FT-005 готова
к tasking; HTML fixtures относятся к проверке реализации.
