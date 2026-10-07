# FT-005 — Независимая семантическая проверка после двух исправлений

## Принятый результат

Свежая проверка 2026-10-07 охватывает FT-005-AC-001…009 / REQ-010…014.
Независимые поиски доставляют существующие и новые совпадения, сохраняют свою
историю и аренду. По поручению оператора исправлены оба пункта
[advisory-отчёта](<../../PAPERCUTS/TECHDEBTS/FT-005 __ 10-07-2026 08.15.md>):
неподтверждённый/blocked detail сохраняет retry, а polling читает только историю
текущих IDs выбранного поиска, сохраняя накопленные записи.

Основание: indexed TASK-006…012 (все `done`), feature/PRD/REQ, прямые task-linked
architecture, boundary-map и runtime-lifecycle, Constitution/MBB, spec-backbone
и registry, testing strategy, tier-policy obligations/closure и установленный
`/red-verify` с finding-adjudication pack. Closure принадлежит `/root` по
[feature protocol](../../.protocols/FT-005/plan.md).

## Доказательства и покрытие

Два самостоятельных фокуса проверены по текущим исходникам и независимым probes:

- **Detail body → parser → Client → App retry/backoff (AC-008).** Проверка body
  предшествует fallback merge. Visible blocked marker возвращает typed error;
  чужой body без принятого detail evidence отвергается. Собственный непустой
  DOM heading с detail-признаком либо matching structured advert подтверждает
  detail; sparse fields и hidden modal сохраняют допустимый fallback. Client
  передаёт ошибку в App; generic detail error не отправляет/не записывает ID,
  typed block прекращает общий обход и включает прежний backoff. Прочитанные
  [parser/polling probes и лог](../TASK-011-T2-FT-005-W5/verifier-probe.log)
  проверяют no-send/no-history, reopen → valid retry → один delivered ID,
  subsequent reopen dedup и остановку следующего monitor/manual request.
  [TASK011 independent functional PASS](../../.protocols/TASK-011-T2-FT-005-W5/verification.md).
- **Current feed → Store query → durable history (AC-009).** Production caller
  передаёт IDs фактических cards и selected monitor в `KeywordAdStatesForIDs`.
  Store связывает `monitor_id` и `ad_id IN` параметрами; empty IDs возвращаются
  до lock/SQL. Full diagnostic reader совместим. Query не удаляет и не меняет
  историю; delivered/revision, stale-reject и pre-send проверки сохранены.
  Прочитанные [Store/App probes и exact comparisons](../TASK-012-T2-FT-005-W5/verifier-acceptance-evidence.md)
  подтверждают независимость двух monitor, отсутствие DB access на empty path,
  сохранность семи historical rows через edit/reopen и production reevaluation
  только текущих отказов с последующим restart dedup.
  [TASK012 independent functional PASS](../../.protocols/TASK-012-T2-FT-005-W5/verification.md).

Исходные AC-001…007 сохраняют ранее доказанный результат: source/strict price
([TASK006](../../.protocols/TASK-006-T2-FT-005-W1/verification.md),
[TASK007](../../.protocols/TASK-007-T2-FT-005-W1/verification.md)),
management/rental preservation
([TASK008 functional](../../.protocols/TASK-008-T3-FT-005-W2/verification.md)
и [semantic](../../.protocols/TASK-008-T3-FT-005-W2/red-verification.md)),
delivery/shared limits/history/revision
([TASK009](../../.protocols/TASK-009-T2-FT-005-W3/verification.md)),
обе кнопки удаления и stale work
([TASK010 functional](../../.protocols/TASK-010-T3-FT-005-W4/verification.md)
и [semantic](../../.protocols/TASK-010-T3-FT-005-W4/red-verification.md)).
Взаимодействия двух исправлений с этими owners проверены; новая schema,
dependency, worker или смена SQL/business ownership отсутствуют.

Все 51 текущих source hashes совпадают с
[финальным independent TASK012 snapshot](../TASK-012-T2-FT-005-W5/verifier-source-state.json).
Относительно прежнего semantic basis TASK010 изменены ровно семь ожидаемых
repair files; остальные 44 файла и все девять ранее инспектированных artifact
hashes сохранены. Executor snapshots TASK011/012 также совпадают. Свежие
[focused gate](../TASK-012-T2-FT-005-W5/verifier-focused-gate.log) и
[native gate](../TASK-012-T2-FT-005-W5/verifier-native-gate.log) имеют exit 0
на этих combined sources: formatting, все package tests, vet, CGO build и
SQLite linkage. Проверены сами probes, logs и команды. Повторный runtime run
не потребовался: источник не изменился, достаточное outcome evidence имеется.

Для каждого фокуса fresh GPT-6.1 Sol/xhigh co-review launch и один retry
завершились `agent thread limit reached`; применён предусмотренный pack fallback.
Модель указана оператором. Итоговое суждение принадлежит этому Reviewer.
Доказательства ограничены representative fixtures, local httptest и temporary
SQLite; live DOM compatibility и production этим verdict не подтверждаются.

## Verdict и передача владельцу

SEMANTIC_VERDICT: semantic-pass

`/root` может принять feature completion после обоих исправлений и выполнить
финальный `/mb-sync` с обязательными boundary gates. Reviewer обновил только
этот отчёт и matching feature `Semantic Verification`; lifecycle/task state,
code и normative semantics не менялись.
