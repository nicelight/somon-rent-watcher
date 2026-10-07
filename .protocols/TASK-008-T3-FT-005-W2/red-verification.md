---
description: Independent adversarial verification of keyword-search management, authority isolation and durable mutations.
status: active
---
# Red Verification — TASK-008-T3-FT-005-W2

## Semantic target

- TASK008 владеет AC001/AC007, REQ010/014/002: несколько сохраняемых поисков,
  адресная настройка, сводка и включение; изоляция admin/chat/search и сохранность аренды.
- Прямая normative basis: indexed card; Constitution; accepted PRD/requirements;
  architecture keyword design; boundary-map Modules/Dependency Graph, Telegram,
  Persistence, Shared Data, keyword proposal/shapes; runtime keyword state/mutation
  rules; testing strategy. Backbone/review — Planning Revision 1; TASK006/007 done.
  Polling/history/deletion остаются outcomes последующих TASK009/010.

## Evidence and adversarial coverage

- Независимый Reviewer применил `/red-verify`, finding-adjudication и tier-policy
  obligations/closure authority. Протокол создан по red-verification-template;
  lifecycle принадлежит `/root` согласно FT-005/plan.md.
- [Functional verification](verification.md): PASS; прочитаны реальные verifier
  probes, результаты required focused Docker gate, protected-row/reopen evidence и
  rental poll с отсутствием повторов. Functional matrix повторно не запускался.
- Независимо просмотрены все 12 source/test files фактического change-surface,
  adjacent passive model/catalog/filter и task-start diff существующих файлов.
  [Сверка source hashes](../../.tasks/TASK-008-T3-FT-005-W2/red-verifier-source-state.json)
  подтверждает тот же исходный код, что в successful functional review.
- Два фокуса: полномочия/pending/адресная revision в Telegram; атомарные durable
  mutations/совместимость аренды/registered ownership. Для каждого отдельный
  co-review launch и retry на operator-selected GPT-6.1-Sol/xhigh отклонены лимитом
  agent threads; самостоятельная проверка завершена согласно semantic pack.
- Авторизация выполняется до keyword route; callback ack предшествует новым
  сообщениям. Pending содержит admin/chat key, search ID/action/revision; смена
  выбранного меню отменяет прежний ввод, устаревший input не перезаписывает условия.
  Category/city независимы; missing ID не upsert-ится. Проверены реальные обращения
  transport к App и App к store, включая error/missing/stale branches.
- App владеет validation, disabled create, enable и revision; store — SQL,
  conditional UPDATE и транзакционная повторяемая additive initialization.
  SQLite AUTOINCREMENT сохраняет identity; rental schema/data не переписываются.
  Passive model, process/client/config/scheduler/allowlist и разрешённые edges
  сохранены; производственного прямого SQL в Telegram или composition нет.
- Свежий [targeted probe](../../.tasks/TASK-008-T3-FT-005-W2/red_verifier_atomic_management_test.go)
  проверил 16 раундов: из восьми одновременных updates с одной expected revision
  проходит ровно один полный commit; concurrent enable/edit сохраняет enabled,
  category/city и целостные условия; stale patch после чтения текущей revision
  успешно повторяется; отдельный поиск не меняется. [Лог — PASS, exit 0](../../.tasks/TASK-008-T3-FT-005-W2/red-verifier-probes.log).

Воспроизведение: `docker run --rm --network none -v "$PWD:/src:ro" somon-price-hotfix-builder:latest sh /src/.tasks/TASK-008-T3-FT-005-W2/red-verifier-probes.sh`.
Исходники монтируются read-only; probe копируется в container `/tmp`; SQLite —
`t.TempDir()`, cleanup закрывает DB и удаляет временные ресурсы. Рабочие/production
данные, real Telegram, secrets, source/spec semantics и Git не менялись.

## Admitted findings

none

## Operator questions

none

## Verdict

SEMANTIC_VERDICT: semantic-pass

## Owner handoff

- [Итоговый отчёт](../../.tasks/TASK-008-T3-FT-005-W2/TASK-008-T3-FT-005-W2-S-RED-VERIFY-final-report-docs-01.md).
- Functional PASS и semantic-pass позволяют explicit manual owner `/root`
  принять closure TASK008, записать evidence/lifecycle и выполнить W2 `/mb-sync`.
  Reviewer оставил задачу `in_progress`; dependents не продвигал.
- Resume route: n/a; следующую selection TASK009 выполняет owner после W2 boundary.
