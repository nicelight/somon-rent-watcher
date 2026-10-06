---
description: Execution handoff to fresh functional and T3 semantic reviewers.
status: active
---
# Handoff — TASK-008-T3-FT-005-W2

## Summary
Authorized administrators can reach “Мои поиски” from existing main menu or `/searches`, create a trimmed phrase with default All/country disabled, inspect summary, choose readable category/city, set optional inclusive price bounds, explicitly enable/disable and address edits to one stable ID. App validation/revision and additive SQLite preserve protected rental state. New callbacks ack before append-only output; pending admin/chat/search/revision rejects stale input and cancels safely on navigation.

## Where to look
- `internal/app/app.go`
- `internal/app/keyword_search.go`
- `internal/app/keyword_search_test.go`
- `internal/store/sqlite_cgo.go`
- `internal/store/keyword_search_cgo.go`
- `internal/store/keyword_search_nocgo.go`
- `internal/store/keyword_search_test.go`
- `internal/telegram/bot.go`
- `internal/telegram/render.go`
- `internal/telegram/keyword_search.go`
- `internal/telegram/keyword_search_harness_test.go`
- `internal/telegram/keyword_search_integration_test.go`

Additional app.go owner mutex and test-only exported Telegram processUpdate bridge are advisory deviations for the same current outcome. No production exported test route or new dependency. Actual source delta is task-start hashes, not HEAD: `.tasks/TASK-008-T3-FT-005-W2/change-surface.json`; prior dirty TASK007 files/deletion preserved. Hard path allow-list absent; forbidden scope untouched. Only search_monitors added; polling/delivery/history/deletion are subsequent cards.

## How to run / verify
Required gate: `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/telegram ./internal/app ./internal/store'`.
All executor commands/result paths are in progress.md; focused/local gates exit0. Current code test suite includes realApp processUpdate management+harm, App validation/revision, exact legacy SQL row comparison. Initial baseline fixture is archived in .tasks only and intentionally fails authorized unsupported-create expectation against baseline; do not copy it into current source tests.

## Claim-linked RED / GREEN evidence
- AC001: progress.md “Attempt 1 baseline — before production” + “Completed execution / current GREEN”; baseline.log honest compiling RED, claim-green.log current corresponding outcome/harm GREEN.
- AC007 accepted RED_NOT_APPLICABLE: same progress sections, baseline.log initial exact SQLite bytes/rental/auth GREEN; legacy exact SQL EXCEPT and captured rental trace remain GREEN after additive initialization/two writes/reopen.
- Full acceptance artifact: `.tasks/TASK-008-T3-FT-005-W2/TASK-008-T3-FT-005-W2-acceptance-evidence.md`.
- Current-attempt reuse candidates: none; logs supporting-only. No independent provenance claim.

## Known issues
No unresolved blocker/material branch. No live compatibility/production acceptance claim. First render gate index regression corrected without weakening original tests. Baseline preparation pointer typo was setup-only, never RED.

## Follow-ups
Fresh `/verify TASK-008-T3-FT-005-W2`, then T3 `/red-verify TASK-008-T3-FT-005-W2`. Explicit manual owner /root closes after required verdicts and performs W2 sync. Status remains in_progress; executor did not run verify/red-verify/mb-sync or promote next task.
