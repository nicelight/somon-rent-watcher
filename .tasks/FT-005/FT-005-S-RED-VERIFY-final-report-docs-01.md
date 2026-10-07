# FT-005 — Независимая семантическая проверка исправлений и production release

## Принятый результат

Свежая проверка 2026-10-07 охватывает FT-005-AC-001…010 / REQ-008, REQ-010…014. Bounded refresh добавляет разрешённый production release AC010 к прежнему независимо проверенному результату AC001…009.
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
Доказательства исходных AC001…009 ограничены representative fixtures, local httptest и temporary SQLite; широкая live DOM compatibility этим verdict не подтверждается. Production acceptance AC010 подтверждается отдельным свежим release evidence ниже.

## Bounded production refresh — AC010

После TASK013 independent functional PASS и отдельного per-task semantic-pass проверены exact8b48c4cef11237716e1dbc471cf363018e48c613 release и его actual operational surface. Оператор явно разрешил все текущие исходники, включая rental fallback; старые TASK001…005 claims/statuses не присваиваются. Все indexed TASK006…013 теперь done: `/root` закрыл TASK013 после отдельных functional PASS и semantic-pass; feature lifecycle/sync остаётся у `/root`.

Повторное независимое вычисление SHA256 показало: все51 source files точно совпадают с прежним финальным TASK012 snapshot; `git diff HEAD^1 HEAD` пуст, ancestor93cbe9 сохранён. Поэтому прежние два фокуса AC008/detail→retry/backoff и AC009/current-feed→history/restart/revision и governing code evidence сохранены без повторных probes. AC001…009 доказательства и границы остаются прежними.

Для новой AC010 проверены два фокуса: exact current-source publication/build/runtime/host isolation и persisted-state/clone-vs-live/additive-init/backup. [TASK013 semantic protocol](../../.protocols/TASK-013-T3-FT-005-W6/red-verification.md) связывает прямые canonical release inputs, реальный release-procedure.py/ordered log, native local/target gates, staged/live doctor и свежие verifier-live-probe.py/json с outcome. Exact release checkout/build/install/process совпадают; single active unit NRestarts0 связана с DB. Все14675 прежних seen rows, settings/stable state и env сохранены, offset monotonic, original DB identity и unrelated host fingerprints прежние. Поиски/history пусты, automatic activation отсутствует; rental paused, baseline не сброшен. Backup root-only на host, stage/scratch удалены. Первое FF refusal было до runtime mutation; fresh retry привёл к exact FF без source drift/reset.

Required Codex Luna/xhigh launches каждого нового фокуса и один retry отказали по thread limit; pack разрешает продолжение без co-reviewer, без новой модели или дополнительного gate. Reviewer выполнил bounded adversarial coverage самостоятельно. Новых production writes/restarts/messages/tests, code или lifecycle изменений не было.

## Verdict и передача владельцу

SEMANTIC_VERDICT: semantic-pass

`/root` может принять TASK013 closure и feature completion после production acceptance и выполнить
финальный `/mb-sync` с обязательными boundary gates. Reviewer обновил только
этот отчёт и matching feature `Semantic Verification`; lifecycle/task state,
code и normative semantics не менялись.
