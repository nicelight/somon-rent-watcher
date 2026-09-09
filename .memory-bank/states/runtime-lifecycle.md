---
description: Current-state lifecycle plus accepted exact-before-fallback and append-only target transitions.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/invariants.md
  - internal/app/app.go
  - internal/store/sqlite_cgo.go
  - internal/filter/settings.go
---

# Runtime lifecycle — current state

## Accepted target delta

The existing baseline, pause, active, recovery, backoff, first-seen, and delivery-retry lifecycles remain. The current implementation changes only these transitions:

### Active poll candidate phases

1. Query `seen_ads` once for the current feed and retain only IDs fresh at poll start.
2. Evaluate exact candidates first under the existing filters and shared detail cap.
3. A fully matching exact ad sets the poll's exact-match outcome before delivery; its delivery error leaves it unseen but still suppresses fallback.
4. A cap-deferred or transient-detail exact candidate makes exact evaluation incomplete, leaves the ID unseen, and suppresses fallback for the cycle.
5. Only a complete exact-empty result permits fallback evaluation. Eligible candidates have known price in `(PriceMax, floor(1.5 * PriceMax)]`, pass all other filters, and share the remaining detail cap.
6. Send at most three eligible fallback ads in ascending price with stable feed-order ties. Confirmed delivery becomes seen; failed delivery remains unseen.
7. Cards conclusively rejected by their applicable path become seen under the existing first-seen policy. No durable candidate phase or shown-history state is added.

### Telegram admin transitions

1. An authorized callback is acknowledged promptly.
2. The action reads/mutates current shared settings or queues the existing manual poll.
3. Navigation/mutation result is sent as a new message; old messages are not edited or deleted.
4. Valid price/negative/interval input sends a fresh main menu after persistence.
5. Manual scan uses separate accepted/completion/error messages and retains single-flight/backoff behavior without busy-button edits.

These target transitions are normative through the PRD/invariants even while the remainder of this document preserves mapped current-state evidence for implementation comparison.

## Scope warning

Это descriptive state map фактической реализации. Русские строки `RuntimeStatus.Mode` являются UI/runtime text, а не принятым stable state enum.

## Service lifecycle

| Current condition | Entry evidence | Exit / transition |
|---|---|---|
| Starting | `runService` loads full config, opens SQLite and constructs `App`. | `App.Run` starts Telegram and poll loops or returns construction error. |
| Uninitialized baseline | SQLite `state.initialized` missing/not `1`. | First sane category poll marks every visible ID seen, writes initial ordinary snapshot and `initialized=1`; no group ad is sent. |
| Paused | Default `settings.enabled=false` or admin disables monitoring. | Every sane poll marks unseen visible IDs seen and updates snapshot/time without detail requests or group delivery; admin toggles Enabled to enter active mode. |
| Active / normal | `settings.enabled=true` and scheduler is not blocked. | Each scheduled cycle processes unseen cards, then schedules the next random delay from persisted `poll_min_minutes..poll_max_minutes`. |
| Polling | Scheduler clears `NextPoll` and calls `pollOnce`. | Success returns to paused/normal; ordinary error schedules a normal random delay; block enters backoff. |
| Backoff | Somon block page or HTTP 403/429; Retry-After can extend configured block delay. | Scheduler waits until `NextPoll`, then polls again. |
| Manual poll queued | Authorized callback binds its chat/message to a capacity-one in-process signal and currently changes that menu button to `Сканирую...`. | The single scheduler wakes once; another manual request and active backoff are rejected. Current completion restores the button and reports zero-result, paused or error state; accepted target replaces both edits with separate messages. |
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

### Production diagnostic: website results versus bot, 2026-09-09

Read-only checks found installed commit `7f9c5f50d659`, an active service with zero
restarts, enabled monitoring and 3,765 seen IDs. The last three observed polls
reported `(new_ids, details, sent)` of `(29, 3, 0)`, `(19, 2, 0)`, `(8, 0, 0)`.
The configured source remained the default unfiltered Dushanbe category.
Of 60 unique ad links extracted from the operator's filtered website page,
38 already existed in production `seen_ads`; these cannot be resent by normal
or manual polling. This count covers extracted links, not every search result.

The website URL selected price 3,000–5,200, floor minimum 11 and area minimum 35.
Production settings selected rooms 1/2/3/unknown, floors
1/2/9/10/11/12/13/14+/unknown, and seller count strictly below 5. There is no
area setting in the current filter model. Thus the two searches are not
equivalent. Per-ad exclusion reasons for the remaining website links were not
established; zero sends alone does not prove a parser or delivery failure.
No production settings, seen history, service state or implementation changed.

Follow-up inspection of all 60 detail pages found 14 matches using their
dedicated visible price blocks and the same production filters; 46 fail the
seller limit. Of the 14 matches, 11 are already seen. Independently, 11 of the
14 have inflated category prices in the installed parser and would fail its
prefilter; the remaining 3 pass the unchanged parser/filter and are all seen.
This corrects the initial incomplete explanation: a reproduced price parsing
defect also suppresses suitable ads. For example, ID 17075880 is 5,000 on the
site but parsed as 65,000; ID 17087959 is 4,800, parsed as 74,800 in the category
and 314,800 in the detail. Parsing price from broad text admits adjacent
numbers. No implementation correction or history replay has been applied.
Full per-ad comparison and diagnostic method:
[production search audit](../../.protocols/diagnostics/production-search-audit-2026-09-09/report.md).

Resolution on 2026-09-09: final price extraction hotfix `93cbe9ebcfe6` was deployed
from the installed baseline in an isolated branch. Price metadata/dedicated
nodes and price-only lines now exclude neighbouring numeric text. The native
local/target gates and offline replay of 60 category prices and 38 detail prices
passed, including existing explicit `Цена:` price-line compatibility. All 3,786
rows present before the final install and the settings were preserved. The first
final-version poll completed with 6 new IDs, 1 detail and 1 successful delivery
(ID 17094676 at 5,000 somoni). Seen-history
semantics are unchanged; this release does not replay previously rejected ads.
Evidence: [hotfix verification](../../.protocols/TASK-005-T3-FT-004-W1/verification.md).

## Gap and recovery lifecycle

1. A successful poll stores sorted ordinary IDs and completion time.
2. A later active poll suspects a gap when ordinary snapshots do not intersect and/or the last successful poll is older than the greater of `SOMON_GAP_AFTER` and the persisted maximum poll interval plus five minutes. This keeps an intentionally long configured schedule from being treated as downtime.
3. A snapshot-only gap starts recovery silently and remains visible in service logs. The service notifies every configured administrator privately only when the successful-poll age exceeds the adaptive time threshold, with hourly per-key suppression.
4. It derives room-specific recovery URLs from the active room filter, fetches them sequentially and keeps only cards with a parseable age inside the inferred outage window.
5. Recovery cards merge into the primary feed by ID and enter the normal unseen pipeline once; future polls return to the main category schedule.

## Persisted current state

SQLite schema is created automatically by `store.Open`:

| Table | Current fields | Current writer/reader |
|---|---|---|
| `seen_ads` | `ad_id INTEGER PRIMARY KEY`, `first_seen_at TEXT NOT NULL` | `App.pollOnce/processNewCards` writes through `MarkSeen`; app/status/doctor query. |
| `settings` | singleton row `id=1`, JSON text including filter, enabled state and polling minute range | `App.LoadSettings/SaveSettings`; Telegram admin actions mutate the JSON. |
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
- Telegram pending text-input mode (`price`, `negative` or `interval`); a restart clears it.
- Capacity-one manual-poll and schedule-change signals plus the current manual request chat/message binding; a restart clears them.

## Evidence and proof paths

- [internal/app/app.go](../../internal/app/app.go): lifecycle transitions, recovery and writers.
- [internal/store/sqlite_cgo.go](../../internal/store/sqlite_cgo.go): schema and transactions.
- [internal/filter/settings.go](../../internal/filter/settings.go): settings serialization and validation.
- [internal/app/app_test.go](../../internal/app/app_test.go): baseline/one-time delivery and paused no-backfill integration tests.
- [internal/store/sqlite_test.go](../../internal/store/sqlite_test.go): state/settings/seen roundtrip.
