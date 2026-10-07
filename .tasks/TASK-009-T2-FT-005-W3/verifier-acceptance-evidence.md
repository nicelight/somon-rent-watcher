# TASK-009 independent acceptance evidence

Reviewer: GPT-6.1 Sol. Scope: FT-005-AC-005/006, REQ-012/013 integration delta;
REQ-014/architecture are constraints. Dependencies006/007/008 are done prerequisites;
delete/stale-delete AC002 remains TASK010. No task lifecycle change.

Fresh execution commands from project root:

```sh
sh .tasks/TASK-009-T2-FT-005-W3/verifier-probes.sh gates
sh .tasks/TASK-009-T2-FT-005-W3/verifier-probes.sh probes
```

Both exit0. Script mounts repo `/src:ro`, evidence read-only and copies required Go
source/fixtures into disposable `mktemp`; injects only verifier artifact, formats
that copy and uses network-disabled existing builder. DB always `t.TempDir()/search.db`
via `cfg.DBPath`; HTTP only httptest, native URL rewrite confined to test transport.
Source snapshot before/after unchanged. Receipt metadata/checksums:
[verifier-commands-results.json](verifier-commands-results.json).

New independent probe assertions are
[verifier_polling_outcome_test.go](verifier_polling_outcome_test.go);
they reuse executor **fixture wiring only**, not executor behavioral assertions.
They observe real App, source parser, Telegram requests and durable store rows.
Six probe groups contain fifteen leaf scenarios; race detector clean.

| Owned result | Fresh independent observation | Probe / log |
|---|---|---|
| AC005 first match despite rental seen | Existing7311 seeded rental-seen, first search sends it immediately and stores delivered rev2. Payload contains phrase/title/190 c./Душанбе; captured photo7311/link7311. Second search enabled later independently delivers same7311. | DurableOutcome; verifier-probes.log |
| AC005 new/restart/edit/drop/dedup | New7312 initially rejected at rev2 by both; only A budget change promotes it at rev3. Price drop with new7313 sends7313 to both and keeps B7312 rejected until B conditions change. Close/reopen and subsequent polls retain exactly6 sends and3 delivered rows in each search. Delivered7311 original rev2 row unchanged. | DurableOutcome; verifier-probes.log |
| AC006 source error/parse and detail retry | With existing delivered1001 row, source503/real keyword parse failure/detail502 leaves existing row byte-value-equivalent and no7422 row. Reset failure+close/reopen+next poll delivers7422 once without resending1001. | FailuresPreserveExistingHistory/source-http,source-parse,detail-http |
| AC006 Telegram failure/ambiguity | Telegram500/malformed success reply leaves7422 absent and1001 unchanged. After reset/reopen next poll records7422 only after confirmed success; request count increases by exactly1 retry. | FailuresPreserveExistingHistory/telegram-error,telegram-ambiguity |
| AC006 shared cap/rotation/deferral | cap1, rental plus two enabled searches. Consecutive starts A→B→rental; exactly one detail per cycle (A7555/B1001/rental9002). B7555 remains absent while deferred then succeeds next cycle. Rental9002 seen only after delivery. HTTP max inFlight1, configured5ms minimum server start spacing5.024301ms. | CapFairness; verifier-probes.log |
| AC006 shared block/backoff/manual | HTTP403 and429 from keyword listing **and details** stop all following searches. Runtime backoff≥configured1h, correct HTTP mode; RequestPollNow refuses. History remains empty, exact traces contain2 requests for listing-block or3 for detail-block. | BackoffAndRunningManual/detailfalse-403/429,detailtrue-403/429 |
| AC006 single-flight | One queued manual trigger accepted, queued duplicate refused. While actual scheduler is blocked inside detail,3 more manual triggers refused. Released cycle sends once, max inFlight1. | BackoffAndRunningManual/running-manual |
| AC006 concurrent budget/revision | detail barrier while revision2 is read; update budget to revision3 before response. Stale reject and stale eligible-send separately leave no row/no send. Next current-revision poll succeeds at revision3. | RevisionBarrier/true,false |
| AC006 success/write-failure ambiguity | Telegram success closes temp DB before history write; poll reports persistence error and reopened history empty. Next poll sends one possible repeat, stores success; further reopen/poll does not repeat. No exactly-once claim. | SendWriteAmbiguity |

Package gate re-executed with all app/store/Telegram tests, including rental baseline,
paused/no-backfill, exact/fallback and cap regression checks, exit0. CGO executable build
also exit0. [verifier-gates.log](verifier-gates.log),
[verifier-probes.log](verifier-probes.log).

Independent architecture review: changed App polling orchestrates only registered
App→Somon/Filter/Store/Telegram/SharedData boundaries. Sole `store` SQLite writer;
conditional SQL rejects stale revisions and retains delivered flag; model passive;
Telegram renders/delivers through existing success/error transport. Existing clients,
allowlist/target policy and scheduler retained; no dependency/worker/infrastructure or
source/filter/authorization/deletion changes. T2 scope remains valid.

Executor prospective evidence checked separately: actual creation/enabling baseline
observed only2 rental requests and0 keyword notifications, compiling absence (not broken
setup). Original RED retained, equivalent enabled-search/poll GREEN integrated durable
rows/failure/cap observations; baseline preserved manual/rental-cap GREEN. Recorded
hashes of baseline/green/package/race logs match contents. Executor evidence supports
history of the attempt; every current owned result above has fresh independent proof.

No reuse of execute receipts. One allowed fresh GPT-6.1 Sol xhigh co-review attempt
could not launch (agent thread limit), no retry. Its absence is not a verification gate.
Known pre-existing permissive ParseDetail papercut is outside this integration claim;
actual keyword parse-error and detail-HTTP failure behavior were independently tested.
