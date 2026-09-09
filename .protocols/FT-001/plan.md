# FT-001 planning resume state

## Current outcome

- Feature design: complete by reuse; no canonical spec extension or creation is needed.
- Global Backbone: complete at Planning Revision 1.
- Foundation: `not_required`.
- Queue action: reconcile the existing T2 implementation task in W1.

## Execution-cohesive boundary

Unmerged implementation outcomes are exact-candidate evaluation, bounded
higher-price selection, shared-cap/failure handling, and seen/delivery
transitions. They merge into one outcome because the existing
`internal/app.processNewCards` flow owns all of these state decisions and no
subset can be enabled or proved independently without exposing a partial poll
lifecycle. Filter helpers and all claim-linked tests remain inside that task.

Final candidate:

- `TASK-001-T2-FT-001-W1` — implement and prove the complete exact-before-fallback polling outcome.

No production-only configuration or acceptance task belongs to FT-001.

## Ownership and routes

- Primary owner/root: Polling Application, `internal/app`.
- Supporting provider/root: Filtering, `internal/filter`, through
  `.memory-bank/contracts/boundary-map.md#filtering-contract`.
- Crossed public boundaries: Somon detail/error outcomes, Persistence
  `SeenIDs`/`MarkSeen`, Telegram delivery success/error.
- Forbidden bypasses: direct SQLite access outside `internal/store`, candidate
  choice inside Somon/Telegram, I/O or state ownership inside Filtering, or
  orchestration in `cmd/somonwatch`.

## Queue and next action

- Dependencies: none; FT-000 is not required.
- Initial status: `ready`, because there are no task dependencies or blockers;
  execution still requires current-revision task-plan approval.
- Next action: `/review-tasks-plan FT-001` in a fresh Reviewer context.
