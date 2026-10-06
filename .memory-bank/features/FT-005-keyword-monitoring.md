---
description: Создание, настройка и мониторинг независимых поисков Somon в Telegram без повторов доставленных объявлений.
status: draft
lifecycle: planned
spec_design_status: complete
spec_design_links:
  - ".memory-bank/architecture/system-architecture.md#keyword-monitoring-design-proposal"
  - ".memory-bank/contracts/boundary-map.md#keyword-monitoring-contract-proposal"
  - ".memory-bank/states/runtime-lifecycle.md#keyword-monitoring-state-proposal"
  - ".memory-bank/invariants.md#accepted-must"
  - ".memory-bank/invariants.md#accepted-never"
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
---

# FT-005 — Независимые keyword-поиски

## Value and use cases

Администратор создаёт поиск, например «стол» до 200 сомони, с собственной категорией
и географией, проверяет сводку и включает. Группа сначала получает доступные существующие
совпадения, затем новые. Настройка и удаление относятся к выбранному поиску; аренда работает рядом.

## Sources and constraints

- [Accepted PRD delta](../prd.md#keyword-monitoring-proposal--2026-10-06): все продуктовые правила и acceptance examples.
- [REQ-010…014](../requirements.md#accepted-keyword-monitoring-delta): требования этой feature; REQ-002 сохраняет append-only UI.
- [EP-002](../epics/EP-002-keyword-monitoring.md): продуктовая ценность.
- [Source evidence](../contracts/current-integrations.md#keyword-search-source-observations): native category/city paths, `q` и relevance ordering подтверждены URL оператора.

## Acceptance Criteria

### FT-005-AC-001 — Создание и адресная настройка

- REQ: REQ-010, REQ-014, REQ-002
- Observable criterion: авторизованный администратор создаёт несколько сохраняемых поисков с обязательным словом/фразой, «Все категории» либо выбранной реальной категорией, собственной географией и необязательными целыми границами цены ≥ 0, от ≤ до; проходит сводку и включение, меняет условия выбранного поиска; смена категории сохраняет город. Категория/география выбираются читаемыми кнопками из поддерживаемого начального каталога. Невалидный ввод не сохраняется. Callbacks/input изолированы по admin/chat/search, чужие действия не меняют данные; результат отправляется новым сообщением.
- Verification method: Telegram harness и временная SQLite: создание двух поисков, невалидный ввод, адресная правка, авторизация, пересекающийся pending input и повторное открытие БД.

### FT-005-AC-002 — Два места удаления и текущие условия

- REQ: REQ-010, REQ-013, REQ-014
- Observable criterion: удаление из списка и настроек убирает выбранный поиск и его историю, сохраняя аренду и остальные поиски; старое меню не восстанавливает удалённый ID. Перед новой доставкой проверяется существование/актуальность поиска; устаревшая оценка после правки не становится отказом для новых условий. Уже отправленный Telegram запрос не отзывается.
- Verification method: lifecycle-проверки обоих UI routes, stale callbacks и конкурентной правки/удаления на временной БД с захватом Telegram запросов.

### FT-005-AC-003 — Штатный поиск в выбранном scope

- REQ: REQ-011, REQ-013
- Observable criterion: Somon выполняет штатный matching фразы с выбранными категорией и географией; учитываются только primary results доступной публичной выдачи. Рекомендации других регионов исключаются. Global first-page с local city filter не заменяет scoped search. Подтверждённая пустая выдача успешна, неподтверждённая пустота считается ошибкой и не продвигает историю.
- Verification method: подтверждённый native scoped URL и HTML fixtures с непустой, малой и нулевой primary выдачей, чужими рекомендациями и повреждённой страницей.

### FT-005-AC-004 — Строгие ценовые границы

- REQ: REQ-011
- Observable criterion: при диапазоне 100–200 цены 100 и 200 проходят, 99/201 и неизвестная/договорная цена не проходят. Любая заданная граница требует известной цены в сомони; валюта не конвертируется. Границы необязательны; rental fallback и квартирные фильтры к поиску не применяются.
- Verification method: детерминированные проверки отсутствующих/одной/двух границ, граничных сумм, unknown/negotiable и другой валюты; сценарий товара, которому не подходят квартирные поля.

### FT-005-AC-005 — Первые совпадения и история каждого поиска

- REQ: REQ-012
- Observable criterion: включение доставляет доступные существующие совпадения без тихого baseline, затем новые; уведомление содержит название поиска, товар/предложение, цену, город, фото и ссылку. Доставленный ID не повторяется в этом поиске после restart, изменения условий или снижения цены; тот же ID независимо доставляется другому подходящему поиску. Правка условий позволяет переоценить ранее отклонённый ID, сохраняя доставленные.
- Verification method: captured notification и temporary-SQLite integration: первый/следующий опрос, два поиска с одним ID, restart, price drop, правка бюджета и ранее отклонённый ID.

### FT-005-AC-006 — Ошибки, лимиты и retry

- REQ: REQ-013
- Observable criterion: source error/parse failure не продвигает историю; Telegram failure/неоднозначность и общий detail cap оставляют повторную попытку. Delivery history фиксируется только после success; сбой между отправкой и записью сохраняет принятую неоднозначность возможного повтора. Последовательный общий scheduler соблюдает HTTP delay, общий cap/backoff и single-flight; 403/429 не обходятся.
- Verification method: integration с controlled source/Telegram failures, низким cap, 403/429, interleaving scheduler/UI и повторным открытием SQLite.

### FT-005-AC-007 — Сохранность аренды и минимальная интеграция

- REQ: REQ-014
- Observable criterion: повторное открытие с добавочным хранением сохраняет rental settings, `seen_ads`, state, Telegram offset и существующее rental поведение без повторной рассылки. Поиски используют прежние Go-процесс, SQLite, scheduler, clients, группу и allowlist; новая dependency, worker или инфраструктура не добавляется.
- Verification method: temporary-SQLite проверка сохранности и независимости данных, существующие rental/price regression checks и review затронутых owners.

## Acceptance closure

Edge/failure outcomes включены прямо в FT-005-AC-001…007: невалидный ввод/авторизация,
stale menu/правка/удаление, чужие рекомендации/empty error, неизвестная цена,
restart/revision/dedup, retry/backoff/cap и сохранность аренды.
Полный crawler/archive, currency conversion и новые группы: `Disposition: out_of_scope`
по [PRD](../prd.md#keyword-monitoring-proposal--2026-10-06); `Change route: /write-prd`.

## SDD Design Gate

Первоначальный `/spec-design` завершён: Planning Revision 1, Foundation `not_required`.
Bounded redesign сохраняет revision; feature design complete. Product decomposition
reviewed `APPROVE`; task plan создан, следующий этап — свежий `/review-tasks-plan FT-005`.

- [Architecture](../architecture/system-architecture.md#keyword-monitoring-design-proposal): существующие owners и общий scheduler.
- [Contracts](../contracts/boundary-map.md#keyword-monitoring-contract-proposal): source/filter/UI semantics.
- [State](../states/runtime-lifecycle.md#keyword-monitoring-state-proposal): независимая история, revision, сохранность аренды.
- [Invariants](../invariants.md): общие ограничения и scope арендной дельты.

Native category/city + `q`/`ordering=relevance` подтверждён
[source evidence](../contracts/current-integrations.md#keyword-search-source-observations).
Source blocker закрыт; raw HTML fixtures — execution proof. FT-001…004 и их
очереди/approval сохранены. Этот этап не включает публикацию или deployment.

## Task plan

[IMPL-FT-005](../tasks/plans/IMPL-FT-005.md): пять независимо проверяемых результатов,
полное покрытие AC-001…007 и локальные Docker gates.
Source/price W1 → creation W2 → monitoring W3 → deletion W4;
authoritative records — [task index](../tasks/index.json). Все новые cards `planned`;
старые identities/status/approvals и Planning Revision 1 сохранены.
