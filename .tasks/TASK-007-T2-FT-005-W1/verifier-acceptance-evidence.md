# Независимая проверка TASK-007-T2-FT-005-W1

Проверен только FT-005-AC-004 / ценовая часть REQ-011. Карточка индексирована один раз, ID/T2/FT-005/W1 согласованы, status остаётся in_progress; gates/verify/reqs/depends_on структурно корректны. Все пять обязательных протоколов присутствуют. Scope не требует повышения tier.

Нормативная основа: FT-005-AC-004; PRD Monitoring rules; REQ-011; direct-linked boundary-map Modules, Dependency Graph, Filtering contract, Shared-data contract, Keyword monitoring contract proposal, Keyword search boundary shapes; architecture Keyword monitoring design proposal; testing Risk-based checks и tier-policy.

| Проверяемое правило | Свежая независимая проверка | Наблюдение |
|---|---|---|
| Отсутствующие/одна/две границы, включённые 100/200, исключённые 99/201 | TestVerifierKeywordPriceOutcome: 54 сочетания шести профилей и девяти денежных значений, плюс две равные границы | Все точные ожидания совпали; без границ цена не ограничивает результат |
| Любая граница требует известной цены TJS; unknown/negotiable/USD/unknown currency; отсутствие цены не равно нулю | Та же матрица, min/max/range/zero_min/zero_max | Все четыре неподтверждённых варианта отклонены при каждой границе; известный ноль проходит соответствующие нулевые границы |
| Целые границы ≥ 0 и min ≤ max | TestVerifierKeywordBoundValidation: 11 случаев; дополнительно invalid bounds не дают matching | Отрицательные и обратные границы отклонены; nil/zero/equal допустимы |
| Товар без квартирных полей; отсутствие квартирных/seller/text/promotion ограничений | TestVerifierCommodityOutcome: 100/150/200 с пустыми и с несовместимыми rental полями | Все шесть товаров проходят |
| Нет rental fallback/локального phrase matching/I/O; корректный owner | Независимое чтение keyword_search.go и direct-linked owner/graph | Только scalar monetary inputs; Filtering owns чистые функции; новых module edges, state writes или adapters нет |
| Rental API сохранены; регрессия | SHA256 settings.go/settings_test.go совпадают с началом execution; git diff этих файлов пуст; свежий пакетный gate | Rental исходники неизменны; существующие тесты проходят |

## Executor claim path

Attempt 1 в progress.md и acceptance-evidence.md: сохранён компилирующий baseline_probe_test.go и baseline-red.log с 19 behavioral failures / 36 initial GREEN; claim-green.log с эквивалентными 54 денежными векторами, 11 validation cases и commodity проверками. RED показывает bounded unknown/unconfirmed price и housing coupling через старые entrypoints; GREEN использует новую keyword-only функцию. Initial GREEN для корректных границ сохранён. Setup-only exit 125 исключён из RED. Проверены source, строки логов и SHA256; path не искусственный. Это supporting execution evidence, не независимое доказательство текущего outcome.

## Reused execute evidence

Ни один execute receipt не переиспользован как gate. Executor results остаются supporting-only.

## Repeated checks

Свежий focused Docker package test, go vet и gofmt-check: exit 0. Дешёвые локальные gates выполнены повторно для текущих исходников. Команда/время/exit code: verifier-gates-command.json; вывод: verifier-gates.log. git diff --check также exit 0.

## New targeted probes

Новый verifier-owned external-package probe в verifier_outcome_probe_test.go построен по пользовательским результатам AC004/PRD, с независимыми explicit expectations. Command/время/exit code: verifier-outcome-command.json; вывод: verifier-outcome-probe.log (exit 0). Probe не подменяет source/UI/poll/history/deletion proof других задач.

Репозиторий смонтирован read-only; probe исполнялся в одноразовой копии go.mod/internal/filter/internal/model внутри контейнера, без production/network/Telegram/DB. Все fixtures — in-memory. Контейнер удалён --rm; рабочие implementation файлы не изменены. Инструмент и hashes relevant source записаны в verifier-input-state.json; новые production hashes до/после probe совпали.

## Co-review и ограничения

Installed finding-adjudication pack применён; user override модели GPT-6.1-sol/xhigh передан в две попытки fresh code co-review. Обе отклонены по agent thread limit; продолжение допускается fallback pack. Собственный review не обнаружил материальных task-relevant дефектов. Независимого второго агентского мнения нет.

Внешний user commit d99e3a7 не принят за pre-task baseline. Tracked deletion keyword_price_baseline_test.go соответствует удалению executor-owned временного probe; его exact source сохранён в baseline_probe_test.go. Git mutation/restore/revert/commit verifier не выполнял. TASK-006 source outcome остаётся prerequisite другой задачи, не claim этой проверки.

Результат: AC004 доказан; lifecycle решение принадлежит /root. Feature-wide completion и integration в следующие tasks этой проверкой не подтверждаются.
