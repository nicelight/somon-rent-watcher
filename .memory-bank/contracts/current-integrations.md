---
description: Evidence-backed current-state integration and observed dependency contracts; non-authoritative for target design.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-01
source_of_truth:
  - internal/somon/client.go
  - internal/telegram/client.go
  - internal/store/sqlite_cgo.go
  - deploy/somonwatch.service
---

# Current integrations and observed boundaries

## Scope warning

Здесь зафиксировано то, что текущий код делает на границах. Эти rows не являются allowed target edges и не заменяют accepted [.memory-bank/contracts/boundary-map.md](boundary-map.md).

## External boundaries

### Somon HTML over HTTPS

- Consumer: current `internal/app` through `internal/somon.Client`.
- Provider: configurable category URL and `/adv/...` detail pages on Somon.tj; recovery uses known room-specific category paths.
- Requests: sequential HTTP GET, static User-Agent, optional Referer, shared client/keep-alive, configured minimum delay, timeout and maximum response body.
- Inputs: server-rendered HTML; deterministic Next.js/RSC data is a fallback/enrichment source. Visible DOM controls city-feed membership/order when visible cards exist.
- Failures: non-2xx becomes typed `HTTPError`; 403/429 or detected block page enter blocked handling; parse/sanity failures return raw body to the caller for private diagnostics and do not authorize state advancement.
- No current interaction with Somon private `/api`, `/author/`, pagination, browser automation, proxy rotation or CAPTCHA bypass.
- Evidence: [internal/somon/client.go](../../internal/somon/client.go), [internal/somon/parser.go](../../internal/somon/parser.go), [internal/app/app.go](../../internal/app/app.go).

### Telegram Bot API

- Consumer: current `internal/telegram` client/bot.
- Provider: configurable Telegram-compatible API base, default `https://api.telegram.org`.
- Requests: form-encoded POST to `getMe`, `deleteWebhook`, `getWebhookInfo`, `getChat`, `getUpdates`, `sendMessage`, `sendPhoto`, `editMessageText`, `answerCallbackQuery`.
- Input topology: `getUpdates` long polling only; `run` removes a stale webhook without dropping queued updates.
- Authorization: only the configured admin user in a private chat can process commands/callbacks; ads go to one configured target chat.
- Delivery: remote photo is preferred. A deterministic Telegram 400 photo validation error falls back to text; ambiguous/network/server failure does not immediately fall back to avoid an immediate duplicate.
- Durability: the next Telegram update offset is written to SQLite only after processing an update attempt.
- Token handling: token is embedded in the Bot API URL internally; client-side transport errors replace the full base URL with `<telegram-bot-api>` before returning the message.
- Evidence: [internal/telegram/client.go](../../internal/telegram/client.go), [internal/telegram/bot.go](../../internal/telegram/bot.go), [internal/telegram/render.go](../../internal/telegram/render.go).

### System SQLite through CGO

- Consumer: `internal/store` and, through it, `internal/app`/CLI diagnostics.
- Provider: system `libsqlite3`, linked with `-lsqlite3`; CGO-disabled builds expose the same Go methods but return a clear unsupported error.
- Concurrency: one in-process mutex serializes access to a `SQLITE_OPEN_FULLMUTEX` handle; busy timeout is 5 seconds.
- Durability: WAL, `synchronous=NORMAL`, foreign keys and transactional multi-row writes.
- Privacy: parent directory is created with `0700`; DB file is forced to `0600`.
- Schema and writers: [.memory-bank/states/runtime-lifecycle.md#persisted-current-state](../states/runtime-lifecycle.md#persisted-current-state).
- Evidence: [internal/store/sqlite_cgo.go](../../internal/store/sqlite_cgo.go), [internal/store/sqlite_nocgo.go](../../internal/store/sqlite_nocgo.go).

### Host process, environment and filesystem

- `deploy/somonwatch.service` starts `/opt/somonwatch/somonwatch run` as user/group `somonwatch`, reads `/etc/somonwatch/somonwatch.env` and grants writes only to `/var/lib/somonwatch` under `ProtectSystem=strict`.
- The env template owns required Telegram identifiers/secrets and Somon/polling limits; `internal/config` supplies defaults and validates ranges.
- Install/backup scripts write only the documented `/opt`, `/etc`, `/var/lib`, `/var/backups` paths and systemd unit; install refuses to replace an active service.
- Evidence: [deploy/somonwatch.env.example](../../deploy/somonwatch.env.example), [deploy/somonwatch.service](../../deploy/somonwatch.service), [internal/config/config.go](../../internal/config/config.go), [scripts/install-almalinux.sh](../../scripts/install-almalinux.sh), [scripts/backup-installed.sh](../../scripts/backup-installed.sh).

## Observed internal dependency map

This map comes from imports/calls in the current source tree. It is not an accepted architecture graph.

| Current consumer | Current providers used |
|---|---|
| `cmd/somonwatch` | `app`, `config`, `somon`, `store`, `telegram` |
| `internal/app` | `config`, `filter`, `model`, `somon`, `store`, `telegram` |
| `internal/config` | `somon` for the default category URL |
| `internal/filter` | `model` |
| `internal/somon` | `htmlx`, `model` |
| `internal/store` | system SQLite C API |
| `internal/telegram` | `filter`, `model`, `somon` (`NormalizeText` in rendering) |
| `internal/htmlx`, `internal/model` | standard library only |

## Needs verification

- Live compatibility with the current Somon DOM/RSC payload and Telegram account/chat configuration is not proven by repository reads; run `somonwatch doctor` on the intended host.
- The current published Somon rules/robots state is external and time-sensitive; re-check before sustained production use.
- No accepted versioning/compatibility policy exists yet for internal packages, SQLite schema or stored settings JSON; current code has no explicit migration framework.
