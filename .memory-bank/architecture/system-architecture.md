---
description: C4 current-state architecture baseline and observed runtime topology for Somon Rent Watcher.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-01
source_of_truth:
  - cmd/somonwatch/main.go
  - internal/app/app.go
  - internal/store/sqlite_cgo.go
  - deploy/somonwatch.service
  - Dockerfile
  - compose.yaml
---

# System Architecture — current state

## Scope warning

Этот C4 baseline описывает фактический repository shape. Наблюдаемые imports/calls и runtime topology не являются accepted target graph или Architecture Decisions.

## C4 L1 — System Context

Somon Rent Watcher связывает настроенный список Telegram administrators, одну target group, публичные HTML pages Somon.tj и Telegram Bot API. Все network connections исходящие; production process не слушает inbound ports.

| External actor/system | Relationship to Somon Rent Watcher |
|---|---|
| Configured administrators | Совместно управляют одним filter/status через private updates или updates из настроенной target group. |
| Target Telegram group | Получает подходящие объявления. |
| Somon.tj | Источник category/detail HTML для polling pipeline. |
| Telegram Bot API | Long polling input и message/photo output. |
| AlmaLinux 9 + systemd | Запускает один hardened host process и предоставляет local filesystem. |
| Kubuntu + Docker Compose | Локально запускает тот же process в одном non-root container и предоставляет named volume. |

## C4 L2 — Containers / runtime units

| Runtime unit | Current responsibility | State |
|---|---|---|
| `somonwatch` executable | CLI (`run`, `doctor`, `ids`, `version`, `help`), application process и diagnostics. | In-memory runtime state plus SQLite-backed durable state. |
| Telegram loop goroutine | `deleteWebhook`, persisted-offset `getUpdates`, admin UI and outbound notifications. | Offset persisted in SQLite `state`. |
| Somon scheduler goroutine | Persisted random polling range, single-flight manual trigger with Telegram completion feedback, baseline/pause/active pipeline, recovery and block backoff. | Poll range lives in singleton JSON `settings`; poll timestamps/snapshot live in SQLite `state`; current manual request is in-memory only. |
| SQLite database | `seen_ads`, singleton JSON `settings`, key/value `state`; WAL enabled. | Default local `./somonwatch.db`, production `/var/lib/somonwatch/somonwatch.db`. |
| Debug HTML directory | Raw category/detail HTML saved only after relevant parse/sanity failures; oldest files are pruned above 20. | Default next to DB, production `/var/lib/somonwatch/debug/`. |

`App.Run` starts exactly two long-lived goroutines and cancels both when the parent context ends or either loop returns an error.

## C4 L3 — Observed components

| Component / code root | Current responsibility |
|---|---|
| `cmd/somonwatch` | Composition root, OS signals and CLI commands. |
| `internal/config` | Environment defaults, parsing and validation. |
| `internal/app` | Scheduler, pipeline orchestration, runtime status, recovery and writers. |
| `internal/somon` | Rate-limited HTTP client, category/detail parser and recovery URLs. |
| `internal/htmlx` | Project-local tolerant HTML tree/parser utilities. |
| `internal/filter` | Filter settings, normalization, validation and two-stage match logic. |
| `internal/store` | Serialized CGO access to system SQLite and schema/state operations. |
| `internal/telegram` | Bot API client, admin-only long-poll UI, rendering and delivery fallback. |
| `internal/model` | Shared `Card`, `Ad` and `RuntimeStatus` data structures. |

Observed internal call/import relationships are catalogued separately in [.memory-bank/contracts/current-integrations.md](../contracts/current-integrations.md): non-authoritative current boundary evidence. The accepted target graph in [.memory-bank/contracts/boundary-map.md](../contracts/boundary-map.md) remains empty until the SDD owner decides it.

## Runtime entrypoints

| Command | Current effect |
|---|---|
| `somonwatch run` / no argument | Loads full config, opens/initializes SQLite and starts both loops. |
| `somonwatch doctor` | Opens/initializes SQLite, verifies Telegram and live Somon parsing, may save diagnostic HTML; does not create baseline or mark IDs seen. |
| `somonwatch ids` | Runs a standalone Telegram `getUpdates` reader for discovering user/chat IDs; must not run with the service using the same token. |
| `somonwatch version` | Prints linker-injected version/commit/build time. |

## Main current data flow

1. Scheduler GETs category HTML through the rate-limited Somon client.
2. `internal/somon` builds Card values from visible DOM and deterministic Next.js/RSC fallback; visible DOM owns actual feed membership/order when available.
3. `internal/app` performs sanity checks before any seen/snapshot mutation.
4. Initial/paused paths write IDs directly to SQLite without group delivery.
5. Active path queries seen IDs, prioritizes ordinary cards, applies card-prefilter, fetches detail HTML, applies detail-filter and calls Telegram delivery.
6. Confirmed delivery is followed by `MarkSeen`; rejected cards and terminal 404/410 are also marked seen, while transient detail/delivery failures remain unseen.
7. Successful cycles atomically update last poll time and ordinary snapshot; Telegram loop independently persists the next update offset.

## Writers and filesystem paths

| Writer | Current path/data |
|---|---|
| `internal/store` | SQLite file, WAL/SHM sidecars, schema and all durable application state. |
| `internal/app.saveDebug` | Private diagnostic HTML, max 20 files. |
| `cmd.saveDoctorHTML` | Private `doctor-*.html` diagnostic page. |
| build/install/backup scripts | `dist/`, `/opt/somonwatch`, `/etc/somonwatch`, `/var/lib/somonwatch`, `/var/backups/somonwatch`. |
| Docker Compose | Named volume `somon-rent-watcher-data` mounted at `/var/lib/somonwatch`; root filesystem remains read-only. |

## Current constraints and unresolved verification

- `go.mod` declares Go 1.21 and no module dependencies; SQLite requires CGO, GCC, headers and runtime `libsqlite3`.
- [docs/TZ_Somon_Rent_Watcher.md](../../docs/TZ_Somon_Rent_Watcher.md) is synchronized to the implemented Go 1.21+/project-local parser/system SQLite CGO stack, but remains draft current-state input rather than an authoritative PRD.
- Git repository was initialized on 2026-09-01; history before that initialization is unavailable. No CI configuration was found.
- Host Go remains absent, but the Docker builder ran the complete native formatting/test/vet/CGO/linkage gate successfully on 2026-09-01.
- Local Docker `somonwatch doctor` verified SQLite, Telegram bot/target group and live Somon parsing on 2026-09-01. The documented production host remains unobserved and requires its own preflight/doctor.

## Architecture Spine

### Accepted Architecture Decisions

None. `/map-codebase` does not derive target `AD-*` decisions from repository evidence.

### Deferred target decisions

All target architecture areas remain owned by [.memory-bank/spec-backbone.md](../spec-backbone.md), whose Global Backbone status is blocked pending PRD decomposition and `/spec-design`.

## Evidence and related docs

- [cmd/somonwatch/main.go](../../cmd/somonwatch/main.go): composition root and CLI behavior.
- [internal/app/app.go](../../internal/app/app.go): concurrency, pipeline, lifecycle and state writers.
- [internal/store/sqlite_cgo.go](../../internal/store/sqlite_cgo.go): concrete SQLite schema and CGO implementation.
- [deploy/somonwatch.service](../../deploy/somonwatch.service): production runtime and hardening.
- [.memory-bank/guides/local-development.md](../guides/local-development.md): build and local verification HOW.
- [.memory-bank/states/runtime-lifecycle.md](../states/runtime-lifecycle.md): detailed current transitions.
- [.memory-bank/runbooks/almalinux-9-operations.md](../runbooks/almalinux-9-operations.md): operator routing.
- [.memory-bank/runbooks/docker-local-operations.md](../runbooks/docker-local-operations.md): local Kubuntu container routing.
