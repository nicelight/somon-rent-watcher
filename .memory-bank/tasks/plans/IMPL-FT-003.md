---
description: Implementation plan for the isolated FT-003 production release.
status: active
feature: FT-003
planning_revision: 1
---

# Implementation Plan — FT-003 isolated production release

## Goal

Publish the locally verified watcher and deploy that exact commit through the
single accepted systemd runtime while preserving settings/seen state and
proving that unrelated production workloads and host controls did not change.

## Scope and non-goals

In scope:

- Close all three local implementation dependencies, run the repository-native
  release gate, create or confirm one clean release commit, push it to the
  configured GitHub remote, and prove the remote ref resolves to that commit.
- Resolve the actual production access command from fresh read-only local
  evidence; take a fresh read-only host/runtime/checkout/state snapshot before
  the first production write.
- Fast-forward only the identified production checkout to the pushed commit,
  build for the target, run the produced binary's redacted doctor, and update
  only the verified `somonwatch.service` through the existing installer.
- Prove the installed commit, unit health/restarts/logs, preserved settings and
  `seen_ads`, and unchanged unrelated host state with ordered redacted receipts.

Out of scope:

- A new deployment script, CI/CD layer, alternate runtime, container action,
  runtime migration, infrastructure change, or rollback promise.
- Force/reset/stash, discarded production-local work, ref/history rewrite, or
  synchronization of any checkout other than the freshly identified watcher.
- Database deletion, replacement, migration, content repair, or schema change.
- Changes to unrelated services, containers, sockets, firewall, routes,
  reverse proxy, SELinux policy/mode, packages, databases, or secrets.
- Reimplementation or re-proof of FT-001/FT-002 claims owned by dependencies.

## Cohesive implementation strategy

`TASK-004-T3-FT-003-W2` owns one production-acceptance result. Its sequence is:

1. Require all three dependencies closed, identify the intended release tree,
   run `./scripts/build.sh`, and require the gated source tree to remain
   unchanged. Attribute every releasable worktree change to the reviewed task
   queue; an unrelated or unexplained change stops before staging.
2. Create or confirm the ordinary release commit, require a clean worktree,
   record its full ID, push without force to the configured GitHub remote/ref,
   and verify the remote ref resolves to that exact ID.
3. Before connecting, resolve the real access command read-only. The text
   `sh igorprod` and earlier SSH-alias observation are not executable authority.
   Inspect current command/shell/SSH configuration without executing an
   unverified target; ambiguity stops for the operator.
4. Through the resolved authorized route, run a fresh read-only production
   preflight. Continue only when one healthy host `somonwatch.service` maps
   consistently to one active process/executable, service user, environment,
   SQLite/debug paths, and watcher checkout, with no competing process or
   container and healthy unrelated workloads.
5. Make the first production write only after preflight: fetch the already
   verified remote ref into that checkout, require a clean non-diverged
   fast-forward, advance without reset/stash/force, and require production HEAD
   to equal the release commit.
6. Run `./scripts/build.sh` in the synchronized checkout, require target linkage
   and exact embedded commit, then run that produced binary through the existing
   service-user environment with secret-redacted output. Build/doctor failure
   leaves the installed runtime unchanged and stops.
7. Stop only `somonwatch.service`, take stable read-only settings/SQLite
   evidence, run the existing scoped installer without database/config reset,
   and start only `somonwatch.service`.
8. Verify installed commit, active unit, restart count, redacted doctor, scoped
   recent logs, stable settings digest and nondecreasing seen count, then repeat
   and compare the unrelated host snapshot. Any retry restarts at step 3/4.

These steps are inseparable: publication supplies the immutable deployment
identity; preflight authorizes the only production mutation; target proof
prevents incompatible install; and postflight closes the same state/isolation
claim. KISS therefore yields one final task rather than publication, deploy,
and acceptance subtasks.

## Dependencies and waves

- Foundation: `not_required`; there is no FT-000 dependency.
- Direct dependencies: `TASK-001-T2-FT-001-W1`,
  `TASK-002-T3-FT-002-W1`, and `TASK-003-T2-FT-002-W1`.
- Wave: W2, the highest and final feature wave.
- Initial status: `planned` until all dependencies reach an accepted closed
  state. Dependency proof remains with those tasks.

## Expected advisory change surface and hard boundaries

Local project writes during this task are limited to ordinary Git metadata for
the authorized release commit/ref and the ignored `dist/somonwatch` and
`dist/somonwatch.sha256` outputs of `scripts/build.sh`. Workflow-owned task
evidence/lifecycle writes remain governed by their skills. No production source
edit is allowed after the local gate.

External writes are limited to:

- the configured GitHub release ref through a normal non-force push;
- the freshly identified production watcher checkout and its `dist` outputs;
- watcher-owned install/backup paths used by the existing scoped installer,
  `/opt/somonwatch`, `/etc/systemd/system/somonwatch.service`, systemd reload
  metadata, and stop/start lifecycle of `somonwatch.service`;
- no content change to `/etc/somonwatch/somonwatch.env` or
  `/var/lib/somonwatch/somonwatch.db`; both must remain at their established
  paths, and evidence must not reveal the environment contents.

The hard boundary excludes other refs/checkouts, services, containers,
packages, networks, security controls, databases, and host files. If the
existing installer would act outside the accepted watcher-owned targets under
fresh conditions, execution stops before invoking it.

## Accepted boundaries and invariants

- Architecture: [Deployment boundary](../../architecture/system-architecture.md#deployment-boundary) and [AD-004](../../architecture/system-architecture.md#ad-004--deployment-is-runtime-detected-and-isolated).
- External edge: [Production boundary rule](../../contracts/boundary-map.md#external-boundary-rules).
- Runtime/state: [Current deployment shape](../../runbooks/almalinux-9-operations.md#current-deployment-shape) and [persisted state](../../states/runtime-lifecycle.md#persisted-current-state).
- Operations: [Accepted FT-003 release procedure](../../runbooks/almalinux-9-operations.md#accepted-ft-003-release-procedure) and [FT-003 release proof](../../runbooks/almalinux-9-operations.md#ft-003-release-proof).
- Safety: [Accepted MUST](../../invariants.md#accepted-must), [Accepted NEVER](../../invariants.md#accepted-never), and the runbook [safety boundaries](../../runbooks/almalinux-9-operations.md#safety-boundaries).

The production runtime must be exactly one healthy `somonwatch.service`; a
different runtime is not an alternative branch. The exact release commit must
match local, GitHub, production checkout, built binary, and installed binary.
SQLite/settings stay in place, and stable unrelated host observations must
match before and after.

## Sources and canonical SDD coverage

- Product authority: [PRD REQ-005](../../prd.md#req-005--ordered-and-isolated-release), [PRD NFR-003](../../prd.md#nfr-003--production-isolation), [REQ-005](../../requirements.md#req-005--ordered-release), and [REQ-008](../../requirements.md#req-008--production-isolation-prd-nfr-003).
- Exact claims: `.memory-bank/features/FT-003-isolated-production-release.md#FT-003-AC-001` and `.memory-bank/features/FT-003-isolated-production-release.md#FT-003-AC-002`.
- Verification policy: [Risk-based checks](../../testing/strategy.md#risk-based-checks) and the [current operational proof paths](../../testing/current-coverage.md#live-and-operational-verification).

All applicable concerns use the canonical links above (`reuse`). No
architecture, contract, state, data, testing, runbook, guide, ADR, or behavior
spec extension/creation is required.

## Verification targets and evidence

FT-003-AC-001 receives one ordered redacted receipt chain covering dependency
closure, local gate, full commit ID, GitHub remote/ref equality, resolved access
route class, preflight systemd identity/health, safe production Git transition,
target build/version/linkage/doctor, the sole unit stop/install/start, and
post-version/health/log checks. A natural RED is the fresh production version/
HEAD mismatch before release. On a safe rerun already at the exact healthy
commit, do not manufacture RED by downgrade or interruption; preserve
pre-implementation GREEN and perform the accepted no-broader-mutation proof.

FT-003-AC-002 uses accepted alternative proof because deliberately losing
state or modifying an unrelated workload would violate the claim. Record
read-only pre-stop and stable stopped-state database identity/integrity,
settings digest and `seen_ads` count, then post-start equality/non-loss. Compare
redacted stable identities/statuses for unrelated containers/services, failed
units, sockets, firewall, routes, and SELinux. Store no token, env contents, raw
host inventory, or public infrastructure identity in the receipt.

Required local gate:

- `./scripts/build.sh`

The synchronized target runs the same gate plus exact version/linkage and the
redacted service-user `doctor`. All evidence is stored in
`.tasks/TASK-004-T3-FT-003-W2/TASK-004-T3-FT-003-W2-acceptance-evidence.md`.

## Constitution constraints

Apply schema-backed planning, minimal verifiable change, evidence before done,
and the project KISS gate. T3 probes must be already authorized, read-only
until the explicit release write, safely rerunnable from a fresh preflight, and
strictly confined to the watcher. Any need to infer access/runtime, broaden
permissions, expose secrets, overwrite Git state, mutate SQLite, operate a
container, or touch unrelated host state stops execution.

## Completion route

`/verify TASK-004-T3-FT-003-W2` and per-task `/red-verify` must pass before the
lifecycle owner can close the T3 task. After all FT-003 work is implemented,
`/red-verify --feature FT-003` must record `semantic-pass`, followed by the W2
`/mb-sync` boundary.
