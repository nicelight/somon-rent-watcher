---
description: Лог изменений Memory Bank.
status: active
---
# Changelog

## [2026-10-07] Wave 5 — два finding FT-005 исправлены
- TASK011: detail body подтверждается независимо от fallback; visible blocked body сообщает backoff, foreign body оставляет retry; sparse/hidden-modal cases сохранены.
- TASK012: production polling читает состояния только текущих IDs; full diagnostic API совместим, вся история сохраняется.
- Отдельные independent functional PASS, свежий feature semantic-pass, final native build и non-CGO compile PASS. Feature/epic verified, historical tasks/evidence сохранены.

## [2026-10-07] Wave 4 — FT-005 verified
- TASK010 done after independent functional PASS and separate T3 semantic-pass: both delete routes, atomic rollback, protected rows and stale work.
- Separate feature semantic-pass covers all AC001…007; owner closes FT005/EP002 and REQ010…014 as verified.
- Final independent native build passed on unchanged source: all tests, formatting, vet, CGO and SQLite linkage. Source/integration probes remain local; no live deployment or agent commit.
- Reconciled feature/epic/RTM/implementation plan and current coverage; old FT001…004 unchanged.

## [2026-10-06] Wave 3 — FT-005 delivery and history
- TASK009 done after independent functional PASS for AC005/006; fresh probes cover per-search history, retry, shared limits/backoff and revision guards.
- REQ012 implemented; feature/epic remain planned pending TASK010 deletion and feature semantic verification.
- Feature evidence and implementation plan reconciled; runtime source/architecture contracts unchanged.

## [2026-10-06] Wave 2 / Saved search management

- TASK008 закрыта explicit owner после независимых functional PASS и semantic-pass.
- Создание, адресная настройка и включение поисков доступны в Telegram; добавочное
  хранение сохраняет ID/revision и арендные данные. Конкурентные правки проверены.
- REQ010/014 и FT005/EP002 остаются planned: мониторинг и удаление ещё предстоят.
- Evidence: [functional](../.protocols/TASK-008-T3-FT-005-W2/verification.md),
  [semantic](../.protocols/TASK-008-T3-FT-005-W2/red-verification.md).

## [2026-10-06] Wave 1 / Keyword source and strict price

- TASK-006 и TASK-007 закрыты explicit owner после отдельных independent functional PASS.
- REQ-011 implemented; FT-005/EP-002 остаются planned до завершения остальных результатов.
  Feature-level semantic verification ещё предстоит.
- Источник и строгая цена готовы для UI/storage интеграции TASK-008; production не менялся.
- Evidence: [source](../.protocols/TASK-006-T2-FT-005-W1/verification.md),
  [price](../.protocols/TASK-007-T2-FT-005-W1/verification.md).

## [2026-10-06] Strict keyword price eligibility

- TASK-007: pure inclusive keyword price predicate и independent nonnegative bounds validation;
  при любой границе неизвестная/договорная/неподтверждённая или другая валюта исключается.
- Ноль отличается от отсутствующей границы; без bounds цена не ограничивает поиск.
  Квартирные/авторские criteria, phrase matching и rental fallback не участвуют.
- Compiling baseline RED и claim-equivalent GREEN сохранены; focused Docker filter tests,
  vet и formatting прошли. Existing rental implementation/tests неизменны.
- [WHY/WHERE](contracts/boundary-map.md#keyword-price-implementation-routing): границы и handoff.
  TASK-007 остаётся in_progress до independent `/verify` и решения /root.

## [2026-10-06] Keyword scoped source implementation

- TASK-006: отдельный native category/city + q/relevance source и начальный каталог;
  keyword parser отделяет primary DOM от region recommendations и подтверждённого нуля.
- Пассивные search/price-currency/city поля добавлены в Shared Data; rental parser,
  price fixtures и остальные owners сохранены. Новых dependencies/runtime нет.
- Compiling baseline RED показал RSC recommendation leakage; эквивалентный GREEN,
  focused Docker tests и vet прошли. Fixtures репрезентативные, не live captures.
- [WHY/WHERE и границы доказательства](contracts/current-integrations.md#keyword-source-implementation-routing).
  TASK-006 остаётся in_progress до independent `/verify` и решения владельца.

## [2026-10-06] Keyword search task plan

- [IMPL-FT-005](tasks/plans/IMPL-FT-005.md): пять задач TASK-006…010; source/price,
  создание/настройка, мониторинг и удаление. AC-001…007 закреплены за владельцами proof.
- Уточнены существующие canonical boundary/state owners: payloads, revision,
  runtime DB path, условная запись истории и адресное атомарное удаление.
- Planning Revision 1, Foundation `not_required`, прежние task identities/status/approval
  сохранены. Fresh `/review-tasks-plan FT-005` следующий; реализация и deployment не выполнялись.

## [2026-10-06] Keyword monitoring design proposal

- [Проект поиска](prd.md#keyword-monitoring-proposal--2026-10-06): правила, архитектура, история и открытые вопросы. По просьбе оператора убраны повторы и избыточные детали.
- Приняты несколько поисков помимо аренды, сначала существующие объявления, география каждого поиска и удаление из списка/настроек. Соответствующие вопросы закрыты.
- Минимальный дизайн сохраняет owners/scheduler и изолирует историю поисков от аренды; необоснованные baseline/schema-version/framework/rollback детали убраны. Прежний запрет schema/profiles ограничен арендным upgrade.
- Оператор подтвердил строгий бюджет, исключение договорной цены при лимите, штатный matching Somon и отсутствие повторов при снижении цены.
- URL штатной формы подтвердил `/search/?q=…`; semantic design готов к новой decomposition. Native category/city scope mapping подтверждён следующими URL оператора; «Все категории» поддерживается, смена категории сохраняет город. HTML fixtures относятся к execution proof.
- EP-002/FT-005 product review APPROVE сохранён; FT-005 SDD design complete, следующий свежий `/feature-to-tasks FT-005`. Impact bounded только FT-005; revision 1 → 1, Foundation/старые задачи/approvals и production не менялись.

## [2026-09-09] Price extraction hotfix

- Corrected category/detail DOM price extraction and added regressions for adjacent photo counts and house numbers.
- Isolated release checkout based on installed `7f9c5f50d659` excludes unfinished local feature changes. Native local tests/vet/build passed; offline replay matched all 60 category prices and 38 cached detail prices.
- Completed TASK-005-T3-FT-004-W1 after functional and independent semantic verification. Operational evidence: [hotfix protocol](../.protocols/TASK-005-T3-FT-004-W1/progress.md).
- Final deployed commit `93cbe9ebcfe6` also preserves the existing `Цена:` fixture behavior. Target build/doctor and all normalized final host comparisons passed; verified running binary, rollback backup, intact SQLite/settings/environment. First cycle sent one valid new ad. History retained without replay; first-attempt measurement limitation stays documented in historical receipts.

## [2026-09-09] Production empty-result diagnosis

- Read-only production and website comparison established working polls, differing filters, and 38 already-seen IDs among 60 extracted website ad links. Evidence and limits: [runtime lifecycle](states/runtime-lifecycle.md#production-diagnostic-website-results-versus-bot-2026-09-09).
- Production and implementation were unchanged; individual rejection causes for remaining links are unconfirmed.
- Follow-up checked all 60 detail pages: 14 match real prices and production filters, 46 fail seller count; 11 matching IDs are already seen. Reproduced inflated-price parsing for 11 matches, correcting the initial incomplete explanation. [Full comparison](../.protocols/diagnostics/production-search-audit-2026-09-09/report.md). No implementation or production changes.

## [2026-09-01] Initial setup
- Created Memory Bank skeleton
- Seeded core docs (product, requirements, testing, task registry)

## [2026-09-01] Brownfield current-state mapping

- Mapped implemented product, C4 architecture, integrations, runtime/state lifecycle, operations and verification surfaces from repository evidence.
- Added explicit as-is/target separation; accepted Boundary Map and Global Backbone remain undecided.
- Registered missing PRD and prevented roadmap/task generation under the PRD-less rule.
- Recorded unresolved verification: no local Go toolchain, no live external/host probe, no Git metadata and draft-TZ SQLite-driver drift.

## [2026-09-01] Public repository preparation

- Reframed the public README around rental discovery on popular Central Asian housing platforms without naming the current source platform.
- Kept the public claim accurate by documenting one current platform-specific adapter rather than claiming implemented multi-platform support.
- Replaced the concrete production host address and neighbouring-workload identifiers in the public runbook/preflight script with generic operator-safe checks.
- Recomputed delivery-manifest checksums for the changed README, runbook and preflight script.
- Initialized the local Git repository on branch `main` with GitHub-compatible local author identity; pre-initialization history is not recoverable.
- Published branch `main` to the public origin `https://github.com/nicelight/somon-rent-watcher`.

## [2026-09-01] Wave 1 / Local Docker, multi-admin and scheduler controls

- Added: local hardened Docker Compose route with one permanent container and persisted SQLite volume.
- Updated: Telegram control to an administrator allowlist shared across private chats and the configured target group.
- Added: persisted configurable random polling range plus single-flight `Сканировать сейчас` with busy and zero-result feedback.
- Fixed: live Somon promoted-card classification for Tailwind positional classes.
- Verified: Docker formatting/tests/vet/CGO/linkage gate and local SQLite/Telegram/Somon doctor passed.
- Synchronized: draft technical specification, product/architecture/integration/lifecycle/testing/glossary/runbook routes and verification state.

## [2026-09-01] Adaptive continuity threshold

- Fixed: a persisted polling range above the static 45-minute gap threshold no longer makes a normal scheduled delay look like downtime.
- Preserved: a complete loss of overlap between consecutive ordinary-ID snapshots still independently triggers one recovery sweep.
- Verified by an app-level regression for the 10–70 minute polling range.

## [2026-09-01] Telegram description readability

- Changed: each emoji-led field in an advertisement description now starts on its own Telegram caption line.
- Preserved: description normalization, escaping, truncation and compound emoji sequences.
- Verified by a caption-rendering regression based on an observed live advertisement.

## [2026-09-02] Silent snapshot-gap recovery

- Refactored: continuity detection now keeps snapshot discontinuity and overdue-poll age as typed causes instead of coupling behavior to a rendered reason string.
- Changed: ordinary-ID turnover continues to trigger recovery and structured warning logs but no longer sends repetitive private Telegram messages.
- Preserved: an adaptive overdue successful-poll condition still triggers recovery and an hourly rate-limited private warning.
- Verified: an app integration test proves silent snapshot recovery and overdue-only administrator notification.

## [2026-09-02] AlmaLinux production deployment

- Deployed GitHub commit `7f9c5f50d659` from a retained clean checkout as a native, unprivileged systemd service.
- Verified the production build/test/vet/CGO gate and live SQLite, Telegram and Somon doctor before service startup.
- Created a fresh paused 60-card baseline with zero service restarts and no current-ad delivery.
- Confirmed existing containers, listening sockets, firewall configuration and SELinux behavior were unchanged.
- Kept the local Compose container stopped to prevent competing Telegram long polling with the same bot token.


## 2026-10-06 — TASK-008 keyword search management execution

- Создание выключенных поисков, адресные настройки/city-preserving category,
  summary/enable и pending admin/chat/search/revision через существующий Telegram UI.
- Добавочная транзакционная `search_monitors`, стабильные ID и атомарная revision;
  App владеет validation/mutations, store — SQLite, Telegram — authorized append-only routes.
- Compiling baseline create RED; auth/rental initial GREEN; real App/temp SQLite
  management/harm tests и точное сравнение legacy rental rows после writes/reopen.
- [Execution evidence](../.protocols/TASK-008-T3-FT-005-W2/handoff.md):
  independent `/verify` и T3 `/red-verify` впереди; task in_progress, closure `/root`.

## 2026-10-06 — TASK009 monitoring execution

- Сохранённые поиски подключены к существующему последовательному scheduler с общим
  detail cap, delay/backoff и меняющимся началом обхода; rental behavior сохранено.
- Независимая SQLite история хранит delivered и evaluated revision; текущая revision
  проверяется перед send, устаревший reject не подавляет новые условия. Уведомление
  именует поиск и использует прежнюю Telegram transport/retry semantics.
- [Execution evidence](../.protocols/TASK-009-T2-FT-005-W3/handoff.md): честный исходный
  RED, текущий GREEN/race и required package gate. Task in_progress, independent
  verification впереди; closure/production/deletion не выполнялись.

## 2026-10-07 Wave6 — Whole-source production release

- Published and deployed exact8b48c4cef11237716e1dbc471cf363018e48c613, including keyword search, both repaired findings and explicitly accepted current rental fallback.
- Native local/target/staged clone/live doctor PASS; stopped and fresh independent live proof preserve old data/settings/offset/env/DB identity and unrelated host state. Root-only rollback backup retained; scratch removed. Single service zero restarts, prior pause preserved, no automatic monitor.
- TASK013 closed by explicit root manual owner after independent functional/T3 semantic PASS and fresh feature semantic-pass. FT005 AC001…010 verified; older unfinished queue untouched. Memory Bank release/RTM/testing/navigation reconciled; [proof](../.protocols/TASK-013-T3-FT-005-W6/verification.md).
