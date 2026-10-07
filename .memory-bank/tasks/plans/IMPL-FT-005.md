---
description: Минимальный implementation plan независимых keyword-поисков Somon.
status: active
feature: FT-005
planning_revision: 1
---

# FT-005 — Implementation plan

## Goal and scope

Несколько независимо сохраняемых поисков рядом с арендой: phrase, optional category,
своя география и строгие optional price bounds. Создать → сводка → включить;
первые существующие и последующие совпадения; два места удаления. Первоначальная разработка была локальной. Принятое расширение W6 ниже разрешает deployment всей свежей версии.

## Sources and coverage

- [PRD](../../prd.md#keyword-monitoring-proposal--2026-10-06), [REQ-010…014](../../requirements.md#accepted-keyword-monitoring-delta), [FT-005 AC](../../features/FT-005-keyword-monitoring.md#acceptance-criteria): принятые результаты.
- [Architecture](../../architecture/system-architecture.md#keyword-monitoring-design-proposal), [graph](../../contracts/boundary-map.md#dependency-graph): reuse; существующие owners и edges.
- [Boundary shapes](../../contracts/boundary-map.md#keyword-search-boundary-shapes): extend существующего canonical contract, payloads, catalog и отсутствие ownership bypass.
- [Persistence](../../states/runtime-lifecycle.md#keyword-persistence-and-mutation-rules): extend existing canonical state/data owner, runtime path и mutation semantics.
- [Source evidence](../../contracts/current-integrations.md#keyword-search-source-observations): reuse; подтверждённые category/city + `q`/`ordering=relevance`, без выдуманных query fields.
- [Invariants](../../invariants.md), [testing](../../testing/strategy.md#risk-based-checks), [Constitution](../../constitution.md): reuse; KISS, один writer и evidence before done.

Новых specs не требуется. Deployment/event/agent-tool concerns не изменяются.
Global Backbone complete, Planning Revision 1; Foundation `not_required`.

## Outcomes, owners and dependencies

| Task | Owned AC | Primary owner / root | Depends on |
|---|---|---|---|
| TASK-006-T2-FT-005-W1 | AC-003 | Somon Adapter / `internal/somon` | none |
| TASK-007-T2-FT-005-W1 | AC-004 | Filtering / `internal/filter` | none |
| TASK-008-T3-FT-005-W2 | AC-001, AC-007 | Polling Application / `internal/app`; Telegram owns UI | 006, 007 |
| TASK-009-T2-FT-005-W3 | AC-005, AC-006 | Polling Application / `internal/app` | 006, 007, 008 |
| TASK-010-T3-FT-005-W4 | AC-002 | Polling Application / `internal/app`; Telegram owns delete routes | 008, 009 |

Source owns native catalog/scoped URL and primary-only parsing, independent of search
storage. Filter owns pure bounds eligibility independent of apartment settings.
Creation consumes those catalog/validation contracts and delivers one usable saved
search; UI/persistence are not separate task outcomes. Monitoring consumes saved
searches/source/filter; scheduler/history/notification establish one delivery outcome.
Delete is independently observable after creation and monitoring, including cancellation
of stale work. Proof stays with its owner; dependencies do not transfer claims.

Все callbacks идут через [Telegram application boundary](../../contracts/boundary-map.md#telegram-application-boundary).
Source/filter interactions use [polling orchestration](../../contracts/boundary-map.md#polling-orchestration-contract),
[Somon](../../contracts/boundary-map.md#somon-adapter-contract) и [filter](../../contracts/boundary-map.md#filtering-contract).
SQLite operations use [persistence](../../contracts/boundary-map.md#persistence-contract);
passive payloads use [Shared Data](../../contracts/boundary-map.md#shared-data-contract).
Leaf keyword details supplement those existing edges; providers retain compatible rental consumers.
Telegram MUST NOT write DB or choose candidates; store MUST NOT own scheduling;
composition MUST NOT own business orchestration. No graph/revision change.

## Advisory change surface

Use descriptive subject files: `keyword_search.go`, `keyword_search_test.go` within
the owning packages; source catalog/URL/parser in `internal/somon/keyword_search.go`,
passive payloads in `internal/model/keyword_search.go`. Additive SQLite operations
may use `internal/store/keyword_search_cgo.go` and `_nocgo.go`; preserve build identities
and existing `sqlite_cgo.go` initializer. Extend existing bot/render/app/client integration
only as needed. Exact immaterial filenames remain executor discretion. No hard
write allow-list; task semantic scope and forbidden scope remain binding.

## Verification and UAT

Use existing `httptest`, app/bot/store harness and `t.TempDir()` SQLite. Source fixtures
derive from documented public scope pages and preserve primary/recommendation/empty
markers; if current public fetch blocks, report missing parser evidence without bypass.
No live request/production acceptance is a completion requirement. RED must demonstrate
a missing behavior through a compiling existing entrypoint, then equivalent GREEN;
preserve initial GREEN for unchanged authorization, prices, rental rows and limits.

Focused Docker package gates are in each card. Final native `scripts/build.sh` gate
uses existing `somon-price-hotfix-builder:latest`, Go 1.21/CGO/SQLite, with network disabled;
it covers formatting, all tests, vet, build/linkage/version. No extra dependency or
architecture framework/check is introduced. Full native suite remains regression
evidence for the current delta, not adoption of old FT-001/002 claims.

UAT is deterministic local harness: two searches, first results, budget edit/rejected-ID
reevaluation, restart/price drop, both delete routes and a stale menu. Exact evidence,
initial state, rerun and cleanup are in each task. `/verify` is required for all cards;
T3 also requires per-task `/red-verify`; feature completion requires
`/red-verify --feature FT-005`. No new operational approval checkpoint for disposable probes.

## Routing

Tasks are initially `planned`. Fresh `/review-tasks-plan FT-005` precedes the
applicable `/mb-doctor --strict` and sequential execution. Existing FT-001…004 cards,
identities, lifecycle and approvals are unchanged; unfinished TASK-001 is not a
dependency. New feature owns its append-only menus, preserving existing rental UI.

## Execution completion

TASK006…010 done по решениям explicit standalone owner /root. Все пять результатов
имеют independent functional PASS; T3 TASK008 и TASK010 имеют отдельный semantic-pass.
[Feature semantic-pass](../../../.tasks/FT-005/FT-005-S-RED-VERIFY-final-report-docs-01.md)
подтверждает сквозной цикл AC001…007. Authoritative task status/evidence — в indexed cards.

- W1 source/price: [006 verification](../../../.protocols/TASK-006-T2-FT-005-W1/verification.md),
  [007 verification](../../../.protocols/TASK-007-T2-FT-005-W1/verification.md).
- W2 management/preservation: [008 verification](../../../.protocols/TASK-008-T3-FT-005-W2/verification.md),
  [008 semantic review](../../../.protocols/TASK-008-T3-FT-005-W2/red-verification.md).
- W3 polling/history: [009 verification](../../../.protocols/TASK-009-T2-FT-005-W3/verification.md).
- W4 deletion: [010 verification](../../../.protocols/TASK-010-T3-FT-005-W4/verification.md),
  [010 semantic review](../../../.protocols/TASK-010-T3-FT-005-W4/red-verification.md).

Final native scripts/build.sh independently passed: all tests, formatting, vet,
CGO build and SQLite linkage. Current source hashes match verified state.
Only local isolated fixtures/httptest/temporary SQLite were used. No deployment,
production data changes or agent commits; FT001…004 identities/status/approvals preserved.

## Debt repair extension

Operator-authorized outcomes remain independent: TASK011 AC008 validates shared
Somon detail body preserving valid sparse/fallback and App retry/backoff; TASK012
AC009 limits App→Persistence reads to current IDs preserving all history. Both T2
W5 depend on done TASK009 and run sequentially. Existing architecture/source/store
contracts plus boundary-map detail-response-validation/current-feed-history-lookup
apply; no new owner, schema, worker or production action. Local parser/store/poll
RED→GREEN plus focused Docker checks, final native build and fresh feature review.

Repair completion: TASK011/012 done after separate independent functional PASS;
fresh feature semantic-pass covers AC001…009. Final combined-source native build
(all tests/format/vet/CGO) and non-CGO compile PASS; no history cleanup or schema change.
Evidence: .protocols/TASK-011-T2-FT-005-W5/verification.md and
.protocols/TASK-012-T2-FT-005-W5/verification.md. Root owns W5 final sync/gates.

## Production acceptance extension

Единственный новый результат: точная текущая версия работает в существующем production service с сохранением данных и окружения (AC010). Source fixes уже завершены; build/probes/backup/install/postflight принадлежат одному release outcome. Оператор явно разрешил всю свежую версию 2026-10-07, включая rental fallback. TASK013 T3 W6 зависит от done TASK006…012 и не имеет dependents. Owner — существующий operations route scripts/build.sh/install-almalinux.sh/backup-installed.sh; без изменения архитектуры, graph или Planning Revision 1. Canonical reuse/extend: runbooks/almalinux-9-operations.md#accepted-ft-005-release-procedure; states/runtime-lifecycle.md#keyword-persistence-and-mutation-rules. Claim AC010: read-only baseline old installed commit RED, exact deployed commit GREEN; data-loss/isolation RED_NOT_APPLICABLE (вред production), alternative stable before/after preserved-state proof. Disposable clone для staged doctor; no real Telegram delivery test, only existing doctor read calls. Release creates no monitor. Full native local/target gates, fresh independent /verify and T3 /red-verify; owner closure then /mb-sync. Older queue status/claims retained.
