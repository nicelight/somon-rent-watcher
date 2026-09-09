# FT-002 planning resume state

## Current outcome

- Feature design: complete by reuse; no canonical spec extension or creation is needed.
- Global Backbone: complete at Planning Revision 1.
- Foundation: `not_required`.
- Queue action: `rebuild_required`; replace the rejected one-task scope with
  exactly two independently completable implementation outcomes in W1.

## Execution-cohesive boundary

Unmerged implementation outcomes:

1. Prompt authorized callback acknowledgement, silent unauthorized/wrong-chat
   rejection, and fresh navigation/mutation output, with claim-equivalent proof
   that valid text input remains isolated and returns a fresh main menu.
2. Accepted manual-scan start plus completion/error feedback through separate
   messages while Polling Application retains single-flight, pause, and backoff.

The callback navigation/mutation changes and preserved text-input proof stay
together because text input needs no production change and its regression proof
supports the same Telegram administration result; proof alone does not create a
third task. FT-002-AC-004 is assigned once to that primary protocol-conversion
task. Manual-scan feedback is a separate implementation result in the
`e:scan`/`Bot.CompleteManualPoll` paths and can compile and pass its own
accepted trace without the callback/navigation conversion.

Final candidates:

- `TASK-002-T3-FT-002-W1` — implement append-only callback navigation/mutation
  and preserve text-input behavior; owns FT-002-AC-001, FT-002-AC-002, and
  FT-002-AC-004.
- `TASK-003-T2-FT-002-W1` — implement append-only manual-scan feedback while
  preserving Polling Application behavior; owns FT-002-AC-003.

No production-only configuration or acceptance task belongs to FT-002.

## Ownership and routes

- Primary owner/root for both outcomes: Telegram Adapter, `internal/telegram`.
- The manual-scan outcome crosses the existing Polling Application interaction
  defined by `.memory-bank/contracts/boundary-map.md#polling-orchestration-contract`
  and `.memory-bank/contracts/boundary-map.md#telegram-application-boundary`;
  no Polling Application production change is planned.
- Forbidden bypasses: settings persistence outside the Backend/Persistence
  path, scheduler state ownership inside Telegram, authorization decisions
  outside the existing allowlist/target-chat rule, or orchestration in
  `cmd/somonwatch`.
- Both tasks retain narrow non-empty hard source/test boundaries. The callback
  task also authorizes exactly `dist/somonwatch` and
  `dist/somonwatch.sha256`, the deterministic outputs of its required
  `./scripts/build.sh` gate; no `dist/` subtree permission is granted.

## Queue and next action

- Dependencies: none. The two results are independently completable; FT-000 is
  not required and FT-001 is unrelated.
- Wave: both W1. Canonical execution remains sequential; overlapping advisory
  source paths do not authorize experimental parallel execution.
- Initial status: both `ready`, because there are no dependencies or blockers;
  execution still requires current-revision task-plan approval and the
  applicable queue doctor gate because the queue contains T3 work.
- Next action: `/review-tasks-plan FT-002` in a fresh Reviewer context.
