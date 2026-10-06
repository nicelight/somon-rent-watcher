---
description: Bounded task execution plan and lifecycle handoff.
status: active
---
# Plan — TASK-008-T3-FT-005-W2

## Goal
Authorized multiple saved keyword searches, targeted settings/summary/enable, independent persistence; protect all rental rows and transport behavior.

## Non-goals
Keyword polling/delivery/history/deletion are later tasks; rental callback/manual scan redesign, dependency claim ownership, production and external side effects are excluded.

## Inputs / source specs
Task and inputs resolved in context.md. Governing accepted owners: App business mutations/revision, store single SQLite writer, Telegram auth/render/pending, passive model, existing Somon catalog/filter bounds APIs. Backbone/review revision1 valid; dependencies done.

## Constraints / invariants (MUST / NEVER)
Keep admin allowlist and target-chat policy; acknowledge new callbacks before fresh output. Stable search IDs, disabled create, city preserved on category edit; invalid/stale/missing inputs cannot mutate other data or upsert missing IDs. Protect rental settings/seen/state/offset. No new source contract or graph edge.

## Preflight-confirmed change surface
12 production/test files listed in `.tasks/TASK-008-T3-FT-005-W2/change-surface.json`. Advisory additions: app.go mutex keeps mutations with accepted owner; two Telegram test files provide exact processUpdate+real App harness without production test API/import cycle. Existing model needs no change. Search history table/operations remain TASK009. Hard write_boundary not set; forbidden_scope untouched. Baseline probe created after durable in_progress and archived after honest RED.

## Applicable quality gates
- Required exact Docker network-none CGO tests telegram/app/store — PASS focused-gate.log.
- All ./... regression; focused vet; CGO build; noCGO composition compile; gofmt task files — PASS local-gates.log.
- git diff --check tracked task source — PASS.

## Claim-linked RED / GREEN (T2/T3)
- applicability: applicable AC001; accepted RED_NOT_APPLICABLE AC007.
- accepted claim locator(s): FT-005-AC-001 / FT-005-AC-007 in feature file.
- planned test/probe and environment: exact Bot.processUpdate, real App, disposable SQLite, local Telegram httptest capture.
- observable RED: authorized ks:new unsupported and no durable search schema/rows.
- corresponding GREEN: two persisted independently configurable/enabled IDs across reopen; invalid/cross-context/stale input leaves protected values unchanged; callback ack first/zero edits.
- accepted not-applicable reason and alternative proof: AC007 preservation already holds baseline, RED would require artificial corruption; retain initial exact closed SQLite/rental trace GREEN, compare same protected rows/trace after additive initialization and writes/reopen (bidirectional SQL EXCEPT for all stored rental columns).
- T3 isolation/cleanup/permission: new t.TempDir databases, fake IDs/TOKEN, httptest Close and DB Close; no production/local working DB.

## MB-SYNC handoff / owner
Explicit top-level manual owner /root handles closure and wave sync after independent gates. Executor updated WHY/WHERE in boundary-map/features/task plan/changelog. Index/spec registry need no new route because canonical docs reused. RTM requirement closure remains with owner after gates. Existing unrelated changes preserved.

## Definition of done
Executor outcome/local gates/evidence complete; independent `/verify` + per-task `/red-verify` and owner decision remain due. No lifecycle closure by executor.
