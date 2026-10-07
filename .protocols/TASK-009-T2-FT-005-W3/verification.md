---
description: Independent functional proof of TASK-009 keyword polling delivery and history.
status: active
---
# Verification — TASK-009-T2-FT-005-W3

## What was verified
Fresh Reviewer GPT-6.1 Sol independently checked actual source, indexed task and
canonical rules; ran required package gate/build and six new probe groups (15 leaf
scenarios) with race detector. All passed. Scope is FT-005-AC-005 / FT-005-AC-006 polling integration;
dependencies006/007/008 are done prerequisites, deletion010 is excluded.

## Verification basis
Valid unique task index/file/ID; tierT2/waveW3/featureFT005; string-array reqs,
dependencies/verify and gate shape checked. Required context/plan/progress/handoff/
verification and substantive executor artifacts present; status in_progress.

Governing inputs: constitution, MBB, spec-backbone/index, Reviewer role, task-linked
feature FT-005-AC-005 / FT-005-AC-006, PRD accepted keyword decisions, requirements012/013 and constraint014,
EP002, testing strategy risk-based checks, architecture keyword-monitoring design,
boundary-map modules/dependency graph, exact polling/persistence/Telegram and keyword
contract/shape headings, runtime-lifecycle keyword state/mutation rules; relevant
tier-policy acceptance/claim/dependency/hard-boundary/closure rules.
Accepted graph owners and current task semantics are consistent. Keyword additive
history is the explicitly accepted keyword delta; earlier rental schema clauses keep
their original rental scope. No material unresolved interpretation or tier escalation.

## Executor claim path
Original supporting RED is `.tasks/TASK-009-T2-FT-005-W3/baseline.log` plus
`keyword_polling_baseline_test.go` and initial byte snapshot: two real persisted/enabled
searches, first+next scheduler polls issue2 rental-only requests and0 keyword sends;
keyword delivery/history/retry path absent. Test compiles and checks behavioral absence.
Preserved initial GREEN covers manual coalescing/backoff and rental cap. Final
claim-equivalent GREEN uses real creation/enabling/poll lifecycle with independent
history and failure/cap/payload/restart assertions, captured in green.log and
progress.md. The documented broken-detail-fixture correction did not falsify RED:
actual keyword parse-error remains, detail failure uses HTTP error. All executor log
checksums independently match commands-results.json; current source matches executor
final byte hashes. Evidence is honest supporting attempt history, not current PASS proof.

## Reused execute evidence
None. Executor offered no reusable receipt; direct safe reruns were cheaper than
receipt eligibility analysis. Executor logs support the claim path only.

## Repeated checks
`sh .tasks/TASK-009-T2-FT-005-W3/verifier-probes.sh gates` exit0: required
`CGO_ENABLED=1 go test -count=1 ./internal/app ./internal/store ./internal/telegram`
on disposable copy of read-only repo, plus CGO executable build exit0. Includes
app rental baseline/paused/no-backfill/exact/fallback/cap regressions, store additive
schema/legacy preservation and Telegram compatibility tests. No implementation changes.

## New targeted probes
`sh .tasks/TASK-009-T2-FT-005-W3/verifier-probes.sh probes` exit0:
`CGO_ENABLED=1 go test -race -count=1 -timeout=60s -v ./internal/app -run TestVerifierSearch`.
Verifier-authored assertions reuse only local fixture wiring. Complete claims-to-probes
mapping and decisive row/request observations are in
[verifier acceptance evidence](../../.tasks/TASK-009-T2-FT-005-W3/verifier-acceptance-evidence.md).

- FT-005-AC-005: existing rental-seen7311 independently delivered in two searches, including
  second enabled later; full captured phrase/title/price/city/photo/link; new7313 delivered;
  rejected7312 reevaluated only after each search's condition edit; exactly6 sends;
  3 delivered rows per search across reopen/edit/price drop; original7311 rev2 preserved.
- FT-005-AC-006: five source/detail/Telegram failure modes preserve nonempty prior history and
  leave failed7422 absent; subsequent reset/reopen poll delivers7422 only after success.
- FT-005-AC-006: one shared detail per cycle with cap1 across A/B/rental, startsA→B→rental,
  deferredB7555 retries, rental9002 receives budget; maximum HTTP inFlight1 and minimum
  observed server spacing5.024301ms for configured5ms delay (1ms arrival tolerance).
- FT-005-AC-006: listing/detail403/429 each halt remaining requests, activate shared backoff
  and refuse manual trigger; queued/running manual duplicates refuse, one send after release.
- FT-005-AC-006: controlled budget-edit barriers prevent stale rejection and stale send; current
  revision subsequently delivers. Send-success/store-write-failure leaves absent history,
  permits one repeat next poll, then durable dedup after reopen.

Isolation: existing builder Docker network-none; original `/src`/evidence read-only,
disposable source copy, cfg.DBPath points to t.TempDir DB; httptest endpoints only;
cancel/join scheduler, close DB/server and automatic temp/container cleanup. Source
hashes unchanged before/after; no real network, production, secrets or Git mutation.
Exact commands/image ID/checksums: `.tasks/TASK-009-T2-FT-005-W3/verifier-commands-results.json`;
source basis: verifier-source-state.json; logs: verifier-gates.log/verifier-probes.log.

## Scope and architecture
Read actual changed App/store/model/Telegram sources and change-surface byte snapshot.
App owns ordering/current-revision checks/history decisions; shared sequential poll and
remaining cap retained. Store alone writes SQLite with conditional revision SQL and
monotonic delivered flag. Telegram renders/sends using existing transport success/error
policy; SharedData passive. Registered App→Somon/Filtering/Persistence/Telegram/SharedData
edges and exact contracts used. No direct DB bypass, second source of truth, reverse
edge, new dependency/worker/infrastructure, auth or rental-policy/source/filter/deletion
scope change. Existing tolerant ParseDetail is an executor-recorded pre-existing papercut;
no out-of-scope redesign required for proved integration outcomes.

Applied installed finding-adjudication pack: main independent functional/architecture
focus plus fresh implementation concurrency/retry/persistence co-review focus. Per
explicit user/model and parent instruction, attempted one GPT-6.1 Sol xhigh co-review
instead of Luna; agent thread limit prevented launch, no retry. No candidates received;
final judgment rests on independently inspected source and fresh observations.

## Verdict
VERDICT: PASS

All owned current outcomes and required gates have credible independent evidence;
no evidenced task-relevant functional, scope or boundary violation found.

## Handoff
T2 task is closure-eligible for root explicit standalone owner; reviewer leaves lifecycle
in_progress. Owner may close TASK009 and perform W3 MB-SYNC. Feature completion remains
pending TASK010 and the separate feature-level semantic review. No red-verify/mb-sync/
scheduler transition executed inside this verification.
