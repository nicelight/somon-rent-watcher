---
description: Независимая функциональная проверка TASK-011 и FT-005-AC-008.
status: active
---
# Verification — TASK-011-T2-FT-005-W5

## What was verified
- Только FT-005-AC-008: подтверждение detail из body, valid sparse fallback, foreign/blocked failure без send/history, retry и shared backoff. REQ-013/REQ-014 — ограничения этого результата; остальные AC/доказанные dependencies не перепроверялись как новые claims.
- Task index/file/ID/tier/wave/feature согласованы; reqs/depends_on/gates/verify имеют допустимую форму. TASK009 — done, prerequisite.
- Прочитаны context, plan, progress, handoff и существующая verification; .tasks acceptance-evidence, baseline.log, green.log, source-state.json.

## Verification basis
- Прямые canonical task links: architecture/system-architecture.md#accepted-shape; contracts/boundary-map.md#modules, #dependency-graph, #detail-response-validation, #polling-orchestration-contract, #persistence-contract; states/runtime-lifecycle.md#keyword-persistence-and-mutation-rules. Конституция, tier obligations/closure/claim ownership/acceptance/RED-GREEN также прочитаны.
- Критерий: fallback не подтверждает foreign HTML; visible block остаётся typed error; DOM имеет собственный непустой h1 + detail-признак, structured advert — совпадающий ID. Sparse body/hidden modal допустимы; ошибка не продвигает историю и не доставляет; HTTP200 block использует существующий backoff.
- Реальная change surface — parser.go/parser_test.go/keyword_polling_test.go. Somon по-прежнему выдаёт model/error; App сохраняет orchestration/retry/backoff и вызывает Store; Store остаётся SQL writer. Новых dependency/worker/schema/bypass нет, graph edges сохранены.
- Не задан hard write allowlist. Forbidden production/real Telegram/Git/deploy scope соблюдён; /src монтировался read-only, runtime/testfiles только в disposable /work; server/DB — httptest/t.TempDir, container --network none и --rm.

## Executor claim path
- FT-005-AC-008 prospective applicable, attempt 1: baseline.log реально компилирует и показывает принятие foreign/blocked detail из полной fallback; polling detail-parse записывал map[1001:{2 true}]. Это отсутствие принятого поведения, не setup/искусственный RED.
- Sparse/hidden modal — исходный GREEN; он сохранён. Claim-equivalent tests retained после validation. green.log — somon/app/store PASS; новая ветка HTTP200 blocked body использует общий backoff.
- Locators: .protocols/TASK-011-T2-FT-005-W5/progress.md#attempt-1-completed-execution и .tasks/TASK-011-T2-FT-005-W5/{baseline.log,green.log,source-state.json}. Hashes всех трёх изменённых файлов совпали до и после независимых проверок.

## Task-scoped checklist
- [x] FT-005-AC-008 foreign/heading-only/no-own-heading rejected с полной fallback; visible blocked marker даже рядом с detail-признаками — typed blocked.
- [x] FT-005-AC-008 sparse description/visible ID/hidden aria/style modal accepted; body title используется, отсутствующая price/currency/image сохраняется из fallback. Matching RSC ID481 accepted, foreign ID482 rejected.
- [x] FT-005-AC-008 HTTP200 foreign/blocked initial detail: captured target sends=0, SQLite history rows=0; reopen сохранил пустую историю. Следующий valid sparse detail: calls=2 total, sends=1, ID1001 delivered at revision2. Ещё reopen не повторяет доставку.
- [x] FT-005-AC-008 HTTP200 visible block: shared backoff HTTP403; следующий monitor не запрашивался и его history пустая; manual request отклонён.

## Regression / non-goals
- [x] Требуемый somon/app/store gate прошёл, включая существующие rental/price fixtures. Никакая другая feature acceptance не присвоена этой задаче.
- [x] Изменения сохраняют owners/registered edges; production/history cleanup/schema/new worker/dependency scope не появился. Дополнительной tier escalation не требуется.

## Reused execute evidence
Нет. Executor GREEN и RED/GREEN path — supporting evidence; receipts не переиспользовались.

## Repeated checks
Required focused gate заново выполнен на read-only source copy с testdata:
```sh
docker run --rm --network none -v "$PWD:/src:ro" -w /work somon-price-hotfix-builder:latest sh -c 'cp /src/go.mod /work/; cp -a /src/internal /src/testdata /work/; CGO_ENABLED=1 go test -count=1 ./internal/somon ./internal/app ./internal/store'
```
Exit0, все 3 пакета PASS; .tasks/TASK-011-T2-FT-005-W5/verifier-focused-gate-complete.log. Первая copy setup не включала root testdata; её fixture-not-found output в verifier-focused-gate.log не используется как product result. Corrected copy выше устраняет setup omission.

## New targeted probes
Verifier-owned fixtures и assertions записаны только как artifacts; production tests/code не менялись:
- .tasks/TASK-011-T2-FT-005-W5/verifier_detail_probe_test.go
- .tasks/TASK-011-T2-FT-005-W5/verifier_polling_probe_test.go
```sh
docker run --rm --network none -v "$PWD:/src:ro" -w /work somon-price-hotfix-builder:latest sh -c 'cp /src/go.mod /work/; cp -a /src/internal /work/; cp /src/.tasks/TASK-011-T2-FT-005-W5/verifier_detail_probe_test.go /work/internal/somon/verifier_detail_probe_test.go; cp /src/.tasks/TASK-011-T2-FT-005-W5/verifier_polling_probe_test.go /work/internal/app/verifier_polling_probe_test.go; CGO_ENABLED=1 go test -count=1 -v -run "^TestVerifier" ./internal/somon ./internal/app'
```
Exit0; наблюдения и полное AC008 mapping перечислены в checklist и .tasks/TASK-011-T2-FT-005-W5/verifier-probe.log. Runtime: Docker image sha256:802934c231fbe4a359d09258bceccb57fda28128d72b78f390006e73250dbafa; input hashes .tasks/TASK-011-T2-FT-005-W5/verifier-source-state.json. Критерии fail/pass — реальные parser errors, typed error, captured send count, persisted selected-monitor history и App runtime status, а не детали реализации.

## Co-review
Semantic pack finding-adjudication и agents/review-code.md прочитаны. Одна best-effort попытка fresh co-review GPT-6.1 Sol/xhigh (user override Luna): не запустилась, agent thread limit reached. Без retries и без блокировки, как требует /verify. Candidate findings отсутствуют, вывод основан на собственном code/spec review и новых outcome probes.

## Verdict
VERDICT: PASS

## Handoff
T2 closure-eligible для root explicit standalone owner; verifier lifecycle не менял, status остаётся in_progress. Рекомендуется owner closure и последующий отдельный TASK012 workflow. Planning repair/tier escalation/debug не требуются. Feature semantic completion — отдельная граница, эта проверка её не подменяет.
