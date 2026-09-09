# FT-003 planning resume state

## Current outcome

- Feature design: complete by reuse; no canonical spec extension or creation is needed.
- Global Backbone: complete at Planning Revision 1.
- Foundation: `not_required`.
- Queue action: create one final production-acceptance task after all local implementation tasks.

## Execution-cohesive boundary

Unmerged change outcomes are local release verification, GitHub publication,
fresh production identity/preflight, exact-commit production synchronization,
target build/doctor, scoped systemd update, and state/host postflight. They merge
into one production-acceptance outcome because none is an independently
releasable result: every earlier step is a prerequisite or safety proof for the
single authorized production mutation, and the postflight alone establishes
that the release completed without state loss or collateral host impact.

Final candidate:

- `TASK-004-T3-FT-003-W2` — production acceptance: publish the gated commit,
  synchronize it safely, update only the freshly verified systemd watcher, and
  prove version, health, state preservation, and host isolation.

Tests, receipts, RED/GREEN or alternative proof, retry preflight, and postflight
remain in this task and do not create separate implementation tasks.

## Ownership and routes

- Primary owner: the accepted Somon Watcher deployment boundary and AlmaLinux
  operations route.
- Local boundary: the intended release tree, its ordinary Git metadata, the
  configured GitHub remote/ref, and the two ignored `scripts/build.sh` outputs.
- Production boundary: the freshly identified watcher checkout/build, the
  existing scoped installer targets and backup under watcher-owned paths, and
  the lifecycle of `somonwatch.service` only.
- Crossed boundaries: GitHub publication, SSH access resolved from current
  read-only local evidence, systemd, the service-user doctor environment, and
  read-only SQLite/host probes.
- Forbidden bypasses: assumed access commands or runtime identity, non-systemd
  or competing runtimes, container actions, dirty/diverged checkout overwrite,
  force/reset/stash, database mutation/migration/replacement, secret output,
  and unrelated host/service/network/security changes.

## Queue and next action

- Dependencies: `TASK-001-T2-FT-001-W1`, `TASK-002-T3-FT-002-W1`, and
  `TASK-003-T2-FT-002-W1`.
- Wave: W2, strictly after the three W1 local implementation tasks.
- Initial status: `planned`, because all three dependencies are still open.
- Next action: `/review-tasks-plan FT-003` in a fresh Reviewer context.
