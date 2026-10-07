---
description: TASK-009 scoped execution plan and owner handoff.
status: active
---
# Plan — TASK-009-T2-FT-005-W3

## Goal
AC005 first/current/new independent delivery; AC006 retry and shared scheduler limits.

## Non-goals
Dependency source/filter/management claims; deletion/stale-delete TASK010; auth/rental policy, deployment and new infrastructure.

## Inputs / source specs
Indexed card, feature AC005/006 and direct canonical links listed in context. Planning Revision1 equals latest FT005 task-plan APPROVE.

## Constraints / invariants (MUST / NEVER)
App owns candidate/state decisions and shared scheduler; store alone writes SQLite; Telegram renders/delivers. Rental seen/history/filter semantics remain separate. No production or external effects.

## Scope
Independent history/schema/payload, rotating sequential scheduler, shared remaining detail budget, revision-conditional rejection/current send, local regression/integration proof.

## Preflight-confirmed change surface
Expected App/store/Telegram areas retained. Additional same-outcome paths: `internal/app/keyword_polling.go`/test (polling separated from management), `internal/model/keyword_search.go` (passive state), `internal/store/sqlite_cgo.go` (canonical additive schema), `internal/telegram/bot.go` (reuse existing transport policy). No rename.
Hard write_boundary not set; forbidden scopes untouched. Exact byte changes in `.tasks/TASK-009-T2-FT-005-W3/change-surface.json`; pre-existing user changes retained.

## Applicable quality gates
- Required Docker network-none package tests app/store/Telegram.
- Focused claim-equivalent integration with race detector for scheduler/UI interleavings.
- Existing non-CGO build compiles newly added adapter stubs; supporting full focused package race run.

## Claim-linked RED / GREEN (T2/T3)
- applicability: applicable
- accepted claim locator(s): `.memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-005`; `.memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-006`.
- RED: real persisted/enabled searches are absent from existing scheduler; first/next cycles have no keyword delivery/history/retry lifecycle.
- GREEN: same creation/enabling/poll path expanded with durable row/restart/revision/payload/error/cap/backoff/barrier trace assertions.
- Initial GREEN existing manual coalescing/backoff and rental shared cap preserved.

## MB-SYNC handoff / owner
Root GENERAL explicit manual owner from `.protocols/FT-005/plan.md` owns independent verification, closure and wave sync. Executor leaves in_progress and does not run verify/red-verify/mb-sync/promote. Feature/changelog updated with implementation evidence, without acceptance closure. Navigation uses existing feature/protocol links.

## Definition of done
Execution delta and gates complete; durable evidence sufficient for fresh `/verify`, task remains open pending that owner decision.
