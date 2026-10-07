---
description: Accepted requirements and traceability for the current Somon Rent Watcher delta.
status: active
last_updated: 2026-10-06
source_of_truth:
  - .memory-bank/prd.md
---

# Requirements

## Functional requirements

### REQ-001 — Preserve monitoring semantics

Keep baseline, first-seen, pause, recovery, backoff, shared detail cap, and confirmed-delivery-before-seen behavior except for the explicit in-poll fallback classification.

### REQ-002 — Append-only Telegram controls

Acknowledge authorized callbacks promptly and send the resulting menu/state as a new message without runtime `editMessageText`; text input and manual scan also complete through fresh messages.

### REQ-003 — Bounded price fallback

Only after a completed poll finds no fresh exact match, send at most three fresh ads satisfying all other filters with known price in `(PriceMax, floor(1.5 * PriceMax)]`, ordered by price with stable feed-order ties; no `PriceMax` means no fallback.

### REQ-004 — Freshness and failure behavior

Use existing `seen_ads` as the sole freshness authority, ignore all existing seen IDs, mark fully rejected fresh cards seen under the current first-seen policy, preserve unseen retry after ambiguous/failed delivery, and suppress fallback when an exact match exists or exact evaluation is incomplete.

### REQ-005 — Ordered release

Run local gates before GitHub publication, then current production preflight, production git synchronization, target build, scoped runtime update, and postflight while preserving settings and SQLite state.

### REQ-009 — Correct Somon price extraction

Operator-authorized production correction on 2026-09-09: extract the actual price
from the ad's price data without adjacent photo counts or address numbers. Preserve
supported explicit Цена: price-line parsing. Deploy under REQ-005/REQ-008, preserving
settings and seen semantics. Source: [accepted follow-up](prd.md#accepted-production-hotfix-2026-09-09).

## Quality requirements

### REQ-006 — KISS maintainability (PRD NFR-001)

For FT-001/FT-002, reuse current settings, storage, scheduler, clients, and limits with no new dependency, database schema, configuration option, worker, or compatibility layer. Pass/fail changes when that rental delta introduces one of those concepts or duplicates an existing owner. Verification: code/spec review and project-native build gate.

### REQ-007 — Filtering and delivery reliability (PRD NFR-002)

Fallback must never hide an exact match, exceed three ads or the shared detail cap, exceed the inclusive 150% ceiling, or mark a failed delivery seen. Pass/fail changes with mixed candidates, cap exhaustion, boundary prices, or Telegram failure. Verification: deterministic unit/integration tests.

### REQ-008 — Production isolation (PRD NFR-003)

Only Somon Rent Watcher may change during deployment; unrelated workload/service/network/security/database state must remain unchanged. Pass/fail changes with runtime detection, scoped restart/install commands, or before/after host-state differences. Verification: operator-authorized read-only preflight/postflight and scoped health checks.

## Out of scope

- Separate shown-history or backfill of existing `seen` IDs.
- Configurable fallback count/percentage and client-specific Telegram handling.
- Broader scraper, storage, product-profile, or production-infrastructure changes beyond the accepted independent keyword searches.

## Accepted keyword-monitoring delta

Источник всех REQ-010…014: [принятая дельта PRD](prd.md#keyword-monitoring-proposal--2026-10-06).
REQ-004/REQ-006 сохраняют арендный scope; новая история самостоятельна. REQ-002
сохраняет append-only UI, REQ-008 — production isolation; первоначальная разработка была локальной; production acceptance AC010 разрешена оператором 2026-10-07.

### REQ-010 — Управление независимыми поисками

Администратор создаёт несколько сохраняемых поисков: обязательное слово/фраза,
«Все категории» либо выбранная реальная категория Somon, собственная география и необязательные целые цены от/до
в сомони ≥ 0, от ≤ до. Создание проходит через сводку и включение; настройки адресуют
конкретный поиск. Удаление доступно в списке и настройках, удаляет выбранный поиск
с его историей, сохраняет остальные и аренду; старое меню не восстанавливает удалённый ID.

### REQ-011 — Штатный scoped search и строгая цена

Слово/фраза сопоставляется поиском Somon в выбранных категории и географии;
категория необязательна («Все категории»), её смена сохраняет город.
рекомендации других регионов не подходят. Обрабатывается доступная публичная выдача,
без полного crawler/archive и без замены scoped search глобальной первой страницей
с local city filter. Обе ценовые границы включены; rental fallback и квартирные фильтры
не применяются. При любой заданной границе нужна известная цена в сомони;
неизвестная/договорная цена исключается, конвертация валют отсутствует.

### REQ-012 — Первые совпадения и независимая история доставки

При включении доставляются доступные существующие совпадения, затем новые.
Уведомление содержит название поиска, товар/предложение, цену, город, фото и ссылку.
История независима для каждого поиска: доставленный ID не повторяется после restart,
правки условий или снижения цены, но может подходить другому поиску. Изменение условий
позволяет переоценить ранее отклонённые ID, сохраняя историю доставки.

### REQ-013 — Надёжность нового мониторинга

Ошибочная выдача не продвигает историю; пустая выдача успешна только при подтверждении
source empty marker. История доставки фиксируется после Telegram success; ошибка/неоднозначность
доставки и исчерпание лимита сохраняют retry. Общие HTTP delay, detail cap, backoff и single-flight
не обходятся; stale evaluation после правки не закрепляет отказ для новых условий, а удалённый
поиск не начинает новую доставку. Уже отправленный Telegram запрос отозвать нельзя; сбой после
отправки до записи истории сохраняет существующую неоднозначность повтора.
Pass/fail меняется при source/Telegram error, cap exhaustion, 403/429, restart и конкурентной
правке/удалении. Verification: воспроизводимые проверки lifecycle и захваченных запросов.

### REQ-014 — Минимальная интеграция и сохранность аренды

Новая возможность переиспользует существующие Go-процесс, SQLite, scheduler, clients,
Telegram-группу и admin allowlist; без новой dependency, worker или инфраструктуры.
Допускается только принятое добавочное хранение поисков/истории. Rental settings,
`seen_ads`, state, Telegram offset и существующее rental поведение сохраняются без повторной
рассылки. Авторизация, append-only UI и pending input по admin/chat/search остаются изолированными.
Pass/fail меняется при повторном открытии SQLite, работе двух поисков с одним ID, чужом callback,
пересечении input или новой runtime ownership. Verification: review и проверки на временной БД/Telegram harness.

## Requirements Traceability Matrix

| Requirement | Epic | Feature | Acceptance / proof route | Lifecycle |
|---|---|---|---|---|
| REQ-001 | EP-001 | FT-001, FT-002 | FT-001-AC-001, FT-001-AC-004, FT-002-AC-003 | planned |
| REQ-002 | EP-001, EP-002 | FT-002, FT-005 | FT-002-AC-001, FT-002-AC-002, FT-002-AC-003; FT-005-AC-001 | planned |
| REQ-003 | EP-001 | FT-001 | FT-001-AC-001, FT-001-AC-002, FT-001-AC-003, FT-001-AC-004 | planned |
| REQ-004 | EP-001 | FT-001 | FT-001-AC-001, FT-001-AC-003, FT-001-AC-004 | planned |
| REQ-005 | EP-001 | FT-003, FT-004 | FT-003-AC-001, FT-003-AC-002; FT-004-AC-002 (hotfix done) | planned |
| REQ-006 | EP-001 | FT-001, FT-002 | FT-001-AC-005, FT-002-AC-004 | planned |
| REQ-007 | EP-001 | FT-001 | FT-001-AC-001, FT-001-AC-002, FT-001-AC-003, FT-001-AC-004 | planned |
| REQ-008 | EP-001, EP-002 | FT-003, FT-004, FT-005 | FT-003-AC-001, FT-003-AC-002; FT-004-AC-002 (hotfix done); FT-005-AC-010 (production verified, TASK013 done) | planned |
| REQ-009 | EP-001 | FT-004 | FT-004-AC-001; TASK-005-T3-FT-004-W1 | done |
| REQ-010 | EP-002 | FT-005 | FT-005-AC-001, FT-005-AC-002 — Telegram/store lifecycle checks | verified |
| REQ-011 | EP-002 | FT-005 | FT-005-AC-003, FT-005-AC-004 — TASK006/007 functional PASS; FT005 semantic-pass | verified |
| REQ-012 | EP-002 | FT-005 | FT-005-AC-005 — TASK009 functional PASS; FT005 semantic-pass | verified |
| REQ-013 | EP-002 | FT-005 | FT-005-AC-002, FT-005-AC-003, FT-005-AC-006, FT-005-AC-008 — failure/concurrency/detail-validation checks | verified |
| REQ-014 | EP-002 | FT-005 | FT-005-AC-001, FT-005-AC-002, FT-005-AC-007, FT-005-AC-008, FT-005-AC-009 — authorization, preservation, validation and bounded lookup; FT-005-AC-010 production preservation — TASK013 functional/semantic PASS | verified |
