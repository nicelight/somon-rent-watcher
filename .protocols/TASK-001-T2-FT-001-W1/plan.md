---
description: Execution plan for TASK-001-T2-FT-001-W1.
status: active
---
# Plan — TASK-001-T2-FT-001-W1

## Goal

Evaluate fresh exact candidates first and, only after conclusively complete exact-empty evaluation, send up to three closest eligible higher-price ads while preserving one detail cap and existing seen/failure/backoff semantics.

## Non-goals

- No SQLite schema or stored-row reinterpretation.
- No setting, dependency, worker, compatibility layer, second history, or duplicate owner.
- No FT-002, FT-003, live, production, commit, or push work.

## Inputs / source specs

- Task: `.memory-bank/tasks/TASK-001-T2-FT-001-W1.task.json`
- Feature claims: `FT-001-AC-001` through `FT-001-AC-005`
- REQs: `REQ-001`, `REQ-003`, `REQ-004`, `REQ-006`, `REQ-007`
- Normative owners: AD-001, AD-002, polling/filtering/Somon/persistence/Telegram contracts, active poll candidate phases, accepted MUST/NEVER.

## Constraints / invariants

- MUST obtain honest pre-change RED for AC-002/AC-004 and preserve alternative pre-change GREEN for AC-001/AC-003/AC-005.
- MUST share `MaxDetailsPerPoll`, suppress fallback on exact or incomplete exact evaluation, preserve stable feed ties, and use only poll-start-fresh IDs.
- NEVER enter task forbidden scope or decide a new product/public/data-authority branch.

## Scope

### In scope

- `internal/app/app.go` polling candidate orchestration.
- Focused deterministic `internal/app` integration coverage.
- `internal/filter/settings.go` and its tests only if a pure non-price-max helper is needed.
- Required protocol/evidence and task lifecycle bookkeeping.

### Out of scope

- `internal/store`, `internal/config`, dependency files, deployment/runtime, Telegram UI behavior, and all unrelated cleanup.

## Proposed changes

### Preflight-confirmed change surface

- Expected hints kept: `internal/app/app.go`, `internal/app/price_fallback_test.go`; filter paths remain conditional.
- Additional same-outcome files: protocol and `.tasks/` evidence required by `/exe`.
- Hard `write_boundary`: not set.
- `forbidden_scope` / stop-condition check: clear.

## Applicable quality gates

- [ ] `CGO_ENABLED=1 go test -count=1 ./internal/app ./internal/filter`
- [ ] `./scripts/build.sh`

## Claim-linked RED / GREEN

- AC-001: accepted RED_NOT_APPLICABLE; alternative pre-change GREEN for exact success/failure suppressing fallback.
- AC-002: honest pre-change RED showing zero fallback where three ordered ads are required.
- AC-003: accepted RED_NOT_APPLICABLE; alternative pre-change GREEN for seen exclusion/final rejection/no-output behavior.
- AC-004: honest pre-change RED for missing fallback lifecycle plus preserved GREEN for current cap/block/retry behavior.
- AC-005 / REQ-006: accepted RED_NOT_APPLICABLE; alternative pre-change inventory/build, then post-change diff/build.
- Probe: existing httptest Somon/Telegram and temporary SQLite harness; no T3 or external side effect.

## MB-SYNC handoff / owner

- Owner: scheduler.
- Broader durable Memory Bank sync is due after verification/feature boundary, not owned by this `/exe` run.
- Task registry/status closure owner: scheduler/lifecycle owner after `/verify`.

## Definition of done

- Smallest implementation and focused proof satisfy every FT-001 claim and required gate.
- Exact actual files, RED/GREEN, gate output, scope compliance, and verifier handoff are durable.
- Task remains `in_progress` for fresh `/verify TASK-001-T2-FT-001-W1`.
