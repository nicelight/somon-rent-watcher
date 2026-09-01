---
description: Current-state polling, delivery and persistence lifecycle derived from code and tests.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
---

# Runtime lifecycle — current state

## Scope warning

Это descriptive state map фактической реализации. Русские строки `RuntimeStatus.Mode` являются UI/runtime text, а не принятым stable state enum.

## Service lifecycle

| Current condition | Entry evidence | Exit / transition |
|---|---|---|
| Starting | `runService` loads full config, opens SQLite and constructs `App`. | `App.Run` starts Telegram and poll loops or returns construction error. |
| Uninitialized baseline | SQLite `state.initialized` missing/not `1`. | First sane category poll marks every visible ID seen, writes initial ordinary snapshot and `initialized=1`; no group ad is sent. |
| Paused | Default `settings.enabled=false` or admin disables monitoring. | Every sane poll marks unseen visible IDs seen and updates snapshot/time without detail requests or group delivery; admin toggles Enabled to enter active mode. |
| Active / normal | `settings.enabled=true` and scheduler is not blocked. | Each scheduled cycle processes unseen cards, then schedules the next random delay. |
| Polling | Scheduler clears `NextPoll` and calls `pollOnce`. | Success returns to paused/normal; ordinary error schedules a normal random delay; block enters backoff. |
| Backoff | Somon block page or HTTP 403/429; Retry-After can extend configured block delay. | Scheduler waits until `NextPoll`, then polls again. |
| Stopping | Parent context cancelled by SIGINT/SIGTERM or one loop returns. | `App.Run` cancels both loops and process exits. |

## Ad ID lifecycle

| Current event | Seen result |
|---|---|
| First baseline or newly observed while paused | Marked seen without notification. |
| Active card rejected by card-prefilter | Marked seen. |
| Detail cap already reached | Remains unseen; retried while still visible later. |
| Detail HTTP 404/410 | Marked seen. |
| Detail transient/network/parse error | Remains unseen; diagnostic HTML may be stored. |
| Ad rejected by detail-filter | Marked seen. |
| Telegram delivery confirmed | Marked seen after the successful call. |
| Telegram delivery fails or is ambiguous | Remains unseen. |
| Category parse/sanity failure | No visible ID or ordinary snapshot is advanced. |

This produces first-seen semantics. Reappearance, price changes or filter changes do not backfill a previously seen ID.

## Gap and recovery lifecycle

1. A successful poll stores sorted ordinary IDs and completion time.
2. A later active poll suspects a gap when ordinary snapshots do not intersect and/or the last successful poll is older than `SOMON_GAP_AFTER`.
3. The service notifies the administrator with hourly per-key suppression.
4. It derives room-specific recovery URLs from the active room filter, fetches them sequentially and keeps only cards with a parseable age inside the inferred outage window.
5. Recovery cards merge into the primary feed by ID and enter the normal unseen pipeline once; future polls return to the main category schedule.

## Persisted current state

SQLite schema is created automatically by `store.Open`:

| Table | Current fields | Current writer/reader |
|---|---|---|
| `seen_ads` | `ad_id INTEGER PRIMARY KEY`, `first_seen_at TEXT NOT NULL` | `App.pollOnce/processNewCards` writes through `MarkSeen`; app/status/doctor query. |
| `settings` | singleton row `id=1`, JSON text | `App.LoadSettings/SaveSettings`; Telegram admin actions mutate the JSON. |
| `state` | string `key`, string `value` | App poll state and Telegram offset. |

Known `state` keys:

- `initialized` — baseline completion marker `1`.
- `last_successful_poll_at` — RFC3339Nano completion time.
- `previous_ordinary_ids` — JSON array of sorted ordinary IDs.
- `telegram_offset` — decimal next update offset.

There is no explicit schema version or migration table in the current repository. This is a current fact, not a recommendation.

## Non-durable process state

- `RuntimeStatus` UI snapshot, next poll/backoff time and last cycle counts.
- Random generator and throttled admin-notification timestamps.
- Telegram pending text-input mode (`price` or `negative`); a restart clears it.

## Evidence and proof paths

- [internal/app/app.go](../../internal/app/app.go): lifecycle transitions, recovery and writers.
- [internal/store/sqlite_cgo.go](../../internal/store/sqlite_cgo.go): schema and transactions.
- [internal/filter/settings.go](../../internal/filter/settings.go): settings serialization and validation.
- [internal/app/app_test.go](../../internal/app/app_test.go): baseline/one-time delivery and paused no-backfill integration tests.
- [internal/store/sqlite_test.go](../../internal/store/sqlite_test.go): state/settings/seen roundtrip.
