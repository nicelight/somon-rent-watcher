---
description: Evidence-backed inventory of current automated, fixture, build and live verification paths.
status: active
baseline_kind: as-is
last_verified: 2026-09-02
---

# Current test and verification coverage

## Native gate

[scripts/build.sh](../../scripts/build.sh) is the project-native build gate: gofmt check, `go test ./...`, `go vet ./...`, CGO build, linkage/version inspection and checksum generation. Shell syntax and race tests are recorded separately in [docs/BUILD_VERIFICATION.md](../../docs/BUILD_VERIFICATION.md).

No repository CI workflow/configuration was found during file inventory.

## Automated coverage by package

| Surface | Existing proof path |
|---|---|
| Configuration | Defaults, admin-list parsing/legacy fallback, invalid admin IDs and invalid poll-range validation in `internal/config/config_test.go`. |
| HTML utility | Tolerant malformed HTML parsing, normalized text and URL resolution in `internal/htmlx/html_test.go`. |
| Somon parser | DOM cards, discount/current price, visible floor over stale slug, seller fields/unknown, other-city boundary, block detection, RSC fallback, age and recovery URLs in `internal/somon/parser_test.go`. |
| Filters/settings | Inclusive price, room/floor buckets, seller limits/unknown, normalized substrings, polling-range parsing/legacy defaults, input/JSON normalization and non-empty choices in `internal/filter/settings_test.go`. |
| SQLite | Settings, deduplicated seen IDs and state roundtrip in `internal/store/sqlite_test.go` (CGO build tag). |
| App pipeline | Typed continuity detection including adaptive stale-time threshold, silent snapshot-gap recovery and overdue-only admin warning; merge order, single-flight manual-poll rejection/backoff guard, first baseline plus one-time delivery, and paused consumption without later backfill in `internal/app/app_test.go` using `httptest` servers and temp SQLite. |
| Telegram | Menus/caption escaping, size and emoji-field line formatting; polling interval/manual scan busy/completion controls; token URL redaction, Retry-After, photo fallback semantics, multi-admin private/target-group authorization, isolated group text input and notification fan-out in `internal/telegram/*_test.go`. |

## Fixtures

- [testdata/category.html](../../testdata/category.html): category with promoted/ordinary cards and other-city boundary.
- [testdata/detail.html](../../testdata/detail.html): detail values including stale URL-floor conflict, discount and seller fields.
- [testdata/detail_unknown_seller.html](../../testdata/detail_unknown_seller.html): missing seller information.

These are repository fixtures, not proof of the current live Somon page.

## Live and operational verification

- `somonwatch doctor` checks SQLite open/schema, Telegram bot/chat/webhook state and a live category parse with minimum card/ordinary-card sanity. It does not create baseline or mark seen IDs.
- [scripts/preflight-almalinux.sh](../../scripts/preflight-almalinux.sh) captures read-only host/service/container/socket/firewall/SELinux/capacity evidence.
- The post-start procedure in [docs/RUNBOOK_ALMALINUX_9.md](../../docs/RUNBOOK_ALMALINUX_9.md) compares host state and checks service resource use/logs.

## Current local verification status

| Check | Result |
|---|---|
| `sha256sum -c MANIFEST.sha256` | Historical PASS before the current working-tree changes; manifest was not regenerated because no release/archive was requested. |
| Docker native package/test/vet/CGO/linkage gate | PASS on 2026-09-01 through `docker compose build somonwatch`; host Go is not required for this route. |
| Local live Somon/Telegram doctor | PASS on 2026-09-01 inside the permanent Compose container: SQLite, bot, target group and 60-card live parse were OK. |
| Target host build and live doctor | PASS on 2026-09-02: formatting/tests/vet, CGO SQLite linkage, Telegram bot/group and 60-card Somon parse were OK for commit `7f9c5f50d659`. |
| Target host systemd/post-start isolation | PASS on 2026-09-02: fresh paused baseline stored 60 IDs; service active with zero restarts; containers, sockets, firewall and SELinux behavior remained unchanged. |

Historical native evidence remains in [docs/BUILD_VERIFICATION.md](../../docs/BUILD_VERIFICATION.md); the Docker rows above are fresh local reruns from the current wave.

## Needs verification

- Current compliance posture based on time-sensitive external rules.


## Keyword deletion execution coverage

FT-005-AC-002 переиспользует existing application mutation lock и SQLite
`ON DELETE CASCADE`: начатая отправка может завершиться перед удалением;
удаление до начала отправки запрещает новый запрос. Схема/арендная история не меняются.

- [Telegram deletion checks](../../internal/telegram/keyword_delete_integration_test.go): обе кнопки, auth/chat, stale callbacks/input, независимость данных, reopen и непереиспользуемый ID.
- [Application deletion checks](../../internal/app/keyword_delete_test.go): local callback/detail barrier и уже начатый HTTP запрос; после удаления нет новой отправки или восстановления истории.
- [SQLite deletion checks](../../internal/store/keyword_delete_test.go): сбой cascade с откатом всей операции и точное сравнение rental settings/seen timestamps/state/offset.
- [TASK010 execution evidence](../../.protocols/TASK-010-T3-FT-005-W4/progress.md): compiling RED/GREEN, isolated harm/race checks и required gates.
- [Independent functional PASS](../../.protocols/TASK-010-T3-FT-005-W4/verification.md) и
  [T3 semantic-pass](../../.protocols/TASK-010-T3-FT-005-W4/red-verification.md) подтверждены;
  [FT005 semantic-pass](../../.tasks/FT-005/FT-005-S-RED-VERIFY-final-report-docs-01.md) закрывает весь цикл.
- Final local native scripts/build.sh PASS (2026-10-07): formatting/all tests/vet/CGO/SQLite linkage;
  keyword integration/harm probes также прошли race detector. Live keyword DOM/production этим не проверены.

## Keyword debt repair verification

- [TASK011 independent PASS](../../.protocols/TASK-011-T2-FT-005-W5/verification.md):
  HTTP200 foreign/visible blocked detail rejection, sparse fallback/hidden modal
  preservation, no false delivery/history, valid retry and shared backoff.
- [TASK012 independent PASS](../../.protocols/TASK-012-T2-FT-005-W5/verification.md):
  current-ID/monitor bounded reads, empty no-query, exact persisted history across
  edit/reopen and production polling dedup/revision reevaluation.
- [Fresh feature semantic-pass](../../.tasks/FT-005/FT-005-S-RED-VERIFY-final-report-docs-01.md)
  covers AC001…009. Final native scripts/build.sh PASS on unchanged combined source;
  non-CGO compile PASS. One existing server-arrival delay test failed then passed
  unchanged rerun and fresh independent gate; logs retained in TASK012.

## Whole-source production release verification

[TASK013 independent functional PASS](../../.protocols/TASK-013-T3-FT-005-W6/verification.md) and [T3 semantic-pass](../../.protocols/TASK-013-T3-FT-005-W6/red-verification.md) cover AC010: release8b48c4c current source, local/target native gates, staged SQLite-clone service-user doctor and installed live doctor, healthy exact binary/unit/checksum, all prior DB rows/state/offset/env and unrelated-host fingerprints.14675→14681 seen with no loss; serviceNRestarts0. Ordinary first poll retained paused state; no baseline recreation/errors. Keyword live delivery wasn't initiated: no monitor configured; product behavior is independently proved by unchanged local task harness/fixtures. [Operational receipt](../../.tasks/TASK-013-T3-FT-005-W6/release-receipt.json).
