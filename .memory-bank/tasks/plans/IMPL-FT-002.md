---
description: Implementation plan for FT-002 append-only Telegram administration.
status: active
feature: FT-002
planning_revision: 1
---

# Implementation Plan — FT-002 append-only Telegram administration

## Goal

Make Telegram administration visibly append-only through two independently
completable changes: reliable callback navigation/mutation with preserved text
input, and separate manual-scan feedback with unchanged scheduler ownership.

## Scope and non-goals

In scope:

- Acknowledge authorized navigation/mutation callbacks before their action and
  send the resulting current menu/state as a new message.
- Keep valid price, negative-word, and interval completion on fresh main-menu
  messages and prove every mode plus administrator/chat isolation.
- Replace manual-scan busy/restore edits with separate start and
  completion/error messages while preserving the scheduler contract.
- Prove authorized private/target-chat behavior, unauthorized/wrong-chat
  silence, stale-menu current-settings semantics, and zero runtime
  `editMessageText` requests in each owned flow.

Out of scope:

- Client detection, edit retry/cache, old-message deletion, compatibility
  layers, persistence, schema, configuration, dependency, worker, or another
  profile/group.
- Authorization-policy, Telegram offset/storage, Polling Application
  production behavior, polling candidates, Somon behavior, or release work.
- Removing the generic Telegram client edit method after accepted runtime
  control flows stop using it.

## Cohesive implementation strategy

### Callback navigation/mutation and preserved text input

`TASK-002-T3-FT-002-W1` is owned by Telegram Adapter at
`internal/telegram`. Reuse its authorization check, pending key
`(userID, chatID)`, Backend settings methods, renderers, and local Bot API test
server. Navigation and mutation callbacks acknowledge promptly and send the
applicable current menu/state without an edit. Old callback message IDs do not
select stale settings. Existing text-input production behavior remains
unchanged; missing mode/isolation coverage is claim-equivalent preservation
proof, not another implementation outcome.

This task owns FT-002-AC-001, FT-002-AC-002, and FT-002-AC-004. It is T3 because
AC-001 changes output at the administrator authorization boundary, including
silent unauthorized and wrong-chat callbacks. Its probes remain disposable and
in-process.

### Manual-scan feedback

`TASK-003-T2-FT-002-W1` is separately owned by Telegram Adapter. It changes
only the accepted `e:scan` start feedback and `Bot.CompleteManualPoll`
completion/error output to fresh messages. Polling Application remains the sole
owner of request capacity, busy-state release, pause, and backoff; app tests
provide regression evidence without an app production change.

This task owns only FT-002-AC-003. It is T2 because it changes a user-visible
manual-scan lifecycle interaction across the existing Telegram/Polling
Application contract but does not change an authorization decision,
permission, secret, production runtime, or irreversible state.

These results share an adapter and some files but neither implementation is
needed for the other to compile or satisfy its exact AC. They therefore remain
two tasks under the execution-cohesion rule.

## Dependencies and waves

- Foundation: `not_required`; no FT-000 dependency exists.
- Product-task dependencies: none. FT-001 does not supply an FT-002
  prerequisite, and neither FT-002 task supplies the other's outcome.
- Wave: both tasks are W1 and initially `ready`.
- Canonical execution is sequential. Overlapping Telegram paths mean these
  tasks are not candidates for experimental parallel execution.

## Expected advisory change surface

`TASK-002-T3-FT-002-W1`:

- `internal/telegram/bot.go` — callback ordering and fresh
  navigation/mutation output.
- `internal/telegram/bot_test.go` — captured method order/counts,
  authorization, stale-menu/current-settings, and input modes/isolation.
- `dist/somonwatch` and `dist/somonwatch.sha256` — deterministic ignored outputs
  of the required native build gate, authorized only in the hard boundary and
  not product-authored source.

`TASK-003-T2-FT-002-W1`:

- `internal/telegram/bot.go` — manual-scan start and completion/error output.
- `internal/telegram/bot_test.go` — accepted start and completion outcome
  traces.
- `internal/telegram/render.go` and `internal/telegram/render_test.go` — only
  minimum removal/reconciliation of the no-longer-used busy scan rendering.
- `internal/app/app_test.go` — only scheduler single-flight/backoff regression
  proof; `internal/app/app.go` remains unchanged.

Each task card declares its own hard write boundary. Evidence and workflow
bookkeeping remain owned by `/exe`, `/verify`, `/red-verify`, and lifecycle
synchronization contracts.

## Accepted boundaries and invariants

- Owner/topology: [Telegram application boundary](../../contracts/boundary-map.md#telegram-application-boundary) and [Polling orchestration contract](../../contracts/boundary-map.md#polling-orchestration-contract).
- External API: [External boundary rules](../../contracts/boundary-map.md#external-boundary-rules).
- Architecture: [AD-001](../../architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable) and [AD-003](../../architecture/system-architecture.md#ad-003--telegram-administration-is-append-only).
- Lifecycle: [Telegram admin transitions](../../states/runtime-lifecycle.md#telegram-admin-transitions) and [non-durable process state](../../states/runtime-lifecycle.md#non-durable-process-state).
- Rules: [Accepted MUST](../../invariants.md#accepted-must), [Accepted NEVER](../../invariants.md#accepted-never), and [compatibility guardrails](../../invariants.md#existing-compatibility-guardrails).

Telegram must not write SQLite or own scheduler state. Polling Application must
retain manual single-flight/backoff. Authorization and pending-input ownership
remain unchanged, and each changed administration flow makes no message-edit
request.

## Sources and canonical SDD coverage

- Product authority: [PRD REQ-002](../../prd.md#req-002--append-only-telegram-admin-controls), [PRD NFR-001](../../prd.md#nfr-001--simplicity-and-maintainability), and [Requirements](../../requirements.md).
- Exact claims: `.memory-bank/features/FT-002-append-only-telegram-ui.md#FT-002-AC-001`
  through `.memory-bank/features/FT-002-append-only-telegram-ui.md#FT-002-AC-004`.
- Verification: [Risk-based checks](../../testing/strategy.md#risk-based-checks) and the [existing Telegram/application harness](../../testing/current-coverage.md#automated-coverage-by-package).

All applicable concerns reuse the canonical links above. No architecture,
contract, state, data, testing, guide, runbook, ADR, persistence, or behavior
spec extension is required.

## Verification targets and evidence

`TASK-002-T3-FT-002-W1` uses the isolated Telegram HTTP server and in-memory
Backend:

- FT-002-AC-001 starts RED because authorized navigation/mutation currently
  edits before its late acknowledgement and unauthorized/wrong-chat callbacks
  emit rejection output. GREEN is acknowledgement first, one fresh applicable
  current-state result, no edits, and no output/mutation outside authorization.
- FT-002-AC-002 uses accepted alternative proof because current valid inputs
  already persist, clear only their pending key, and send a fresh menu. Capture
  equivalent pre/post GREEN for all modes and cross-admin/chat isolation.
- FT-002-AC-004 uses accepted alternative proof because the baseline already
  lacks the forbidden machinery. Compare production paths, owner/dependency/
  config/storage surfaces, runtime method call sites, and native-build result
  before and after.
- Required gates: `CGO_ENABLED=1 go test -count=1 ./internal/telegram` and
  `./scripts/build.sh`. The hard boundary authorizes only
  `dist/somonwatch` and `dist/somonwatch.sha256` in addition to task source/test
  paths.

`TASK-003-T2-FT-002-W1` uses captured Telegram methods plus scheduler state:

- FT-002-AC-003 starts RED because accepted manual scan currently makes a busy
  edit and completion restore edit, and a successful scan with sent ads can
  omit explicit completion. GREEN sends separate accepted start and each
  success/zero-result/paused/error completion message with no edits while
  preserving duplicate/backoff rejection and busy-state release.
- Required gate:
  `CGO_ENABLED=1 go test -count=1 ./internal/telegram ./internal/app`.

No live Telegram, Somon, SQLite mutation, or production probe is required or
authorized. Each task writes its claim evidence to its own `.tasks/<TASK-ID>/`
acceptance artifact.

## Constitution constraints

Apply schema-backed planning, minimal verifiable change, evidence before done,
and the project KISS gate. Stop if execution requires authorization-policy,
public boundary, persistence/config/dependency, Polling Application production,
polling-candidate, production, tier, dependency, or broader scope changes.

## Completion route

- `TASK-002-T3-FT-002-W1`: `/verify` plus per-task `/red-verify` must pass.
- `TASK-003-T2-FT-002-W1`: `/verify` must pass; per-task `/red-verify` is
  optional under T2 policy.
- After both tasks are implemented, `/red-verify --feature FT-002` must record
  `semantic-pass`, followed by lifecycle-owner closure and wave/feature
  boundary `/mb-sync`.
