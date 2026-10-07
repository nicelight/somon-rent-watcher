---
description: TASK-009 claim-linked execution evidence and reproducibility.
status: active
---
# Progress — TASK-009-T2-FT-005-W3

## Current status
- state: verifying
- last update: 2026-10-06
- task lifecycle: in_progress; execution complete, independent verification next.

## What was done
App now rotates rental plus enabled searches in one poll; passes a shared detail budget; retains retry on source/detail/Telegram failures and cap; stops all remaining requests on typed block. Rental behavior preserved. Search-local history uses revision-conditional rejection and durable delivered flag; send rechecks current existence/enabled/revision under mutation lock. Telegram names the search and uses existing photo/text success/ambiguity rules.
Actual changed source files: `.tasks/TASK-009-T2-FT-005-W3/change-surface.json` compares initial/current bytes, including untracked creations. All source was clean initially; unrelated user dirty docs preserved. Additional advisory paths justified in plan; no hard allow-list, no forbidden scope touched.

## Commands run (with results)
Exact commands/exit0/log checksums: `.tasks/TASK-009-T2-FT-005-W3/commands-results.json`.
- baseline (pre-production): local true creation/enable + current pollOnce; PASS test observes behavior RED. Existing manual/backoff/cap tests PASS.
- required package gate: app/store/Telegram PASS on final source.
- claim-equivalent GREEN with race detector: all keyword polling integration cases PASS on final source.
- supporting package race gate PASS and non-CGO cmd build exit0 before final simplification (all keyword detail errors now remain retryable; stubs unchanged). No independent or exactly-once assertion.

## Claim-linked RED / GREEN (T2/T3)
- attempt: 1
- applicability: applicable
- accepted claim locator(s): .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-005; .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-006
- RED command/probe: baseline command in commands-results.json; `.tasks/TASK-009-T2-FT-005-W3/keyword_polling_baseline_test.go` copied temporarily into internal/app after in_progress and removed before production.
- RED observation and evidence: 2 enabled persisted searches via real App operations; first/next poll cycles produce exactly 2 rental-only source requests and 0 keyword/group deliveries. Search history schema/API absent in initial source snapshot. Keyword delivery/retry lifecycle absent for AC006 because scheduler ignores searches. Compiling behavioral absence, not syntax/setup failure. `.tasks/TASK-009-T2-FT-005-W3/baseline.log`.
- initial GREEN: existing manual coalescing/backoff and rental shared cap PASS in baseline.log; shared Somon transport is unchanged, one existing client retained.
- GREEN command/probe: exact green command in commands-results.json, final `internal/app/keyword_polling_test.go`.
- GREEN observation and evidence: `.tasks/TASK-009-T2-FT-005-W3/green.log`; AC005 both searches durably deliver1001/revision2 and1003/revision2, reevaluate rejected1002 into delivered/revision3 after budget change; 6 total sends across restart/edit/price drop with no repeat. Payload assertions cover phrase/title/price/city/photo/link. AC006 source500/true keyword parse error/detail500/Telegram500/malformed response leave empty rows then deliver on allowed next poll/restart. cap1 trace gives details1/1/0, changing start search1→search2→rental; maximum source HTTP in-flight1; configured5ms delay, minimum observed server-arrival spacing4.746043ms (client gate unchanged;1ms measurement tolerance). HTTP403/429 stop remaining searches and set shared backoff; manual rejected. Five running manual duplicates rejected. Controlled stale-rejection/stale-send barriers preserve new revision. Send-success/store-close records no history, then accepted possible repeat (2 sends) on reopen.
- claim-equivalent probe changes and rationale: final same creation/enabling/poll path expands assertions with new persisted history API, payload/source/Telegram trace and restart/interleavings. Removed invalid broken-detail fixture because existing ParseDetail accepts it with fallback card; real keyword source parse failure remains covered, detail retry uses HTTP error. Original initial RED retained; no artificial break.
- isolation/cleanup: all DB paths from t.TempDir()/cfg.DBPath; close/reopen; httptest close; scheduler context cancel/join; no external network/Telegram/production.

## Reuse Candidates (optional)
None offered. Logs and GREEN are supporting execution evidence; fresh `/verify` owns independent proof and may rerun deterministic gates.

## Boundary/spec compliance
Accepted App→Somon/Filter/Store/Telegram paths only. App decides revision/delivery/history; store performs conditional SQL in its adapter; Telegram renders/delivers; Shared Data passive. No graph change, new dependencies/workers/deletion routes or rental exclusion reuse.

## Open issues / risks
No blocker. Existing tolerant detail parser papercut recorded; allowed send-success/write-failure possible repeat demonstrated. No live source compatibility claim.

## Next step (single concrete action)
Fresh `/verify TASK-009-T2-FT-005-W3` by root; execution does not close/promote/run mb-sync.
