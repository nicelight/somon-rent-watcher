---
description: TASK-009 completed execution handoff to fresh verification.
status: active
---
# Handoff — TASK-009-T2-FT-005-W3

## Summary
Searches now deliver existing/new primary matches through the shared rotating sequential poll scheduler. Search-local delivered/revision history persists across restart and edits; cap/errors defer IDs. Named payload uses existing Telegram success/ambiguity transport. Task remains in_progress pending independent verification.

## Where to look
- `internal/app/keyword_polling.go`/test: orchestration, shared budget, revision/send barrier, claim integration tests.
- `internal/app/app.go`: rental poll split and explicit remaining budget; aggregate runtime status and existing backoff/single-flight retained.
- Store keyword adapters + sqlite_cgo schema: independent history, conditional expected revision, delivered survives edits.
- model keyword passive state; Telegram keyword renderer and reused sendAdCaption transport.
- Full actual byte change set: `.tasks/TASK-009-T2-FT-005-W3/change-surface.json`. Additional advisory paths justified in plan; no hard path allow-list, forbidden scopes untouched.

## How to run / verify
- Required command exactly in task card/commands-results.json; final package PASS.
- Focused local probes: Docker network-none `CGO_ENABLED=1 go test -race -count=1 -v ./internal/app -run TestKeywordPolling`; final GREEN PASS.
- AC005/AC006 RED/GREEN: progress.md#claim-linked-red--green-t2t3; `.tasks/TASK-009-T2-FT-005-W3/TASK-009-T2-FT-005-W3-acceptance-evidence.md`, baseline.log/green.log.
- Probe isolation: t.TempDir DB, httptest URLs, source native URL rewrite only inside test reflection transport, cancel/join/close cleanup.
- Current-attempt reuse candidates: none; all execute evidence supporting-only. Fresh verifier-owned proof required.
- Final source basis: initial-source-snapshot.json vs change-surface.json; preserve all unrelated user work.

## Known issues
Existing ParseDetail tolerates arbitrary detail body when card fallback supplies ID/title (papercut). No source/parser changes. Accepted possible repeat after send success/store failure demonstrated. No live source/Telegram/prod claim.

## Follow-ups
Root GENERAL manual owner: fresh `/verify TASK-009-T2-FT-005-W3`, then closure/MB-SYNC when obligations pass. Do not replay completed execution. TASK010 deletion/stale-delete remains separate and unimplemented. No required unresolved branch.

## Owner closure
Explicit standalone owner /root closes TASK009 done after independent functional PASS
in verification.md and task.verify. W3 complete; feature semantic gate remains due.
