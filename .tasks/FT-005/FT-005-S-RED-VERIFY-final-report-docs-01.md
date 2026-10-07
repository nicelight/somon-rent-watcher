# FT-005 — Независимая семантическая проверка

## Принятый результат

Проверен полный пользовательский цикл FT-005-AC-001…007 / REQ-010…014:
администратор создаёт независимые поиски, задаёт native category/geography и
строгую цену, включает доставку существующих и новых совпадений, меняет условия
и удаляет выбранный поиск из двух меню. История каждого поиска сохраняет
доставленные ID, допускает переоценку отказов после правки и retry после ошибок;
аренда, общие ограничения и прежние владельцы данных сохраняются.

Основание: indexed TASK-006…010, accepted PRD/feature/REQ, direct task-linked
architecture, boundary-map, runtime-lifecycle, source observations и testing
strategy; Constitution/MBB, spec-backbone/index, Reviewer role, tier-policy и
installed `/red-verify` с finding-adjudication pack. Все пять task records имеют
`done`; closure принадлежит `/root` по [.protocols/FT-005/plan.md](../../.protocols/FT-005/plan.md).

## Доказательства и сквозное покрытие

Два самостоятельных фокуса проверены по фактическим исходникам:

- **Source → цена → адресный Telegram UI.** Allowlisted native paths, `q` и
  relevance поступают в Somon adapter; primary DOM ограничен до рекомендаций,
  ноль требует подтверждения. Pure keyword predicate соблюдает optional inclusive
  TJS bounds без apartment/fallback правил. Реальные routes проходят прежнюю
  allowlist/chat проверку, подтверждают callback до нового сообщения и обращаются
  к App; pending input содержит admin/chat/search/action/revision. Category и city
  независимы; создание выключено до сводки/enable. Проверены source/filter/UI
  proof bodies и capture logs, включая цену/город/фото/ссылку уведомления.
- **Scheduler → durable history → правка/удаление.** Один последовательный обход
  делит cap и HTTP delay с арендой и меняет начало по кругу; typed 403/429 прекращают
  цикл и включают общий backoff. История разделена по monitor/ad и отделена от
  rental seen; delivered фиксируется после success и переживает правки/reopen.
  Conditional reject учитывает revision. App согласует mutation с началом send
  и повторно проверяет existence/enabled/revision. Обе кнопки используют один
  delete; SQLite CASCADE атомарно удаляет только выбранные rows. Инспектированные
  fault/barrier probes подтверждают retry, stale-edit/delete, rollback и принятую
  неоднозначность уже начатой отправки без восстановления удалённой истории.

Независимые functional доказательства и достигнутые этапы:

- Source AC-003 и strict price AC-004:
  [TASK006 PASS](../../.protocols/TASK-006-T2-FT-005-W1/verification.md),
  [TASK007 PASS](../../.protocols/TASK-007-T2-FT-005-W1/verification.md).
- Создание/настройка и сохранность аренды AC-001/007:
  [TASK008 PASS](../../.protocols/TASK-008-T3-FT-005-W2/verification.md) и
  [T3 semantic-pass](../../.protocols/TASK-008-T3-FT-005-W2/red-verification.md).
- Доставка/история/retry AC-005/006:
  [TASK009 PASS](../../.protocols/TASK-009-T2-FT-005-W3/verification.md),
  [verifier assertions](../TASK-009-T2-FT-005-W3/verifier_polling_outcome_test.go).
- Обе кнопки удаления и concurrent lifecycle AC-002:
  [TASK010 PASS](../../.protocols/TASK-010-T3-FT-005-W4/verification.md) и
  [T3 semantic-pass](../../.protocols/TASK-010-T3-FT-005-W4/red-verification.md).

Проверены сами verifier probes, logs и native gate, а не только PASS prose.
Все 51 текущих source hashes совпали с
[финальным verified source](../TASK-010-T3-FT-005-W4/red-verifier-source-state.json);
все девять перечисленных там artifact hashes совпали. HEAD —
`7db6f00f64bfb698b6a9bdec138e69d870892811`. Source snapshots учитывают изменения,
не видимые в Git diff после внешнего commit. Финальный independent
[native gate](../TASK-010-T3-FT-005-W4/verifier-native-build-gate.log) имеет exit 0:
formatting, все package tests, vet, CGO build и SQLite linkage.
Дополнительный runtime rerun не потребовался: исходники, реальные barriers и
protected-row snapshots достаточны для двух фокусов этого bounded review.

Для каждого фокуса fresh GPT-6.1 Sol/xhigh co-review launch и один retry завершились
`agent thread limit reached`; применён предусмотренный pack fallback.
Модель выбрана по явному указанию оператора. Итоговое суждение принадлежит этому Reviewer.
Source fixtures остаются representative sanitized HTML; live DOM compatibility,
Somon/Telegram production и deployment этим verdict не подтверждаются.

## Verdict и передача владельцу

SEMANTIC_VERDICT: semantic-pass

`/root` может завершить feature/epic lifecycle и выполнить финальный `/mb-sync` и
обязательные boundary gates. Reviewer записал только этот отчёт и matching feature
`Semantic Verification`; task/feature statuses, исходники и normative semantics
не менял. Рабочие данные, production, Git и реальные Telegram сообщения не затронуты.
