# Исполнение TASK-006-T2-FT-005-W1

Добавлены native scoped keyword source, поддерживаемый каталог и пассивный payload. Parser отделяет primary DOM от рекомендаций; H1 0 подтверждает пустоту, неопределённая выдача возвращает ошибку. Арендные parser/client и fixtures сохранены.

Compiling baseline RED: small/zero выдавали foreign ID 21000002. Эквивалентный GREEN возвращает [21000001]/[]; 16 native scope requests, ошибки и passive fields проверены. Required focused Docker tests и vet PASS; gofmt выполнен. Подробности: [TASK-006-T2-FT-005-W1-acceptance-evidence.md](TASK-006-T2-FT-005-W1-acceptance-evidence.md).

Изменённые task files: internal/somon/keyword_search.go, internal/somon/keyword_search_test.go, internal/model/ad.go, internal/model/keyword_search.go, testdata/keyword-search-primary-small.html, testdata/keyword-search-primary-empty.html, .memory-bank/contracts/current-integrations.md, .memory-bank/index.md, .memory-bank/changelog.md, .memory-bank/tasks/TASK-006-T2-FT-005-W1.task.json; пять protocol files и task-local probe/log/report artifacts. Hard write_boundary не задан, forbidden scope не затронут. Новые passive ad.go поля — допустимое дополнение advisory scope для того же source outcome.

Fixtures репрезентативные, не live captures; source network не использовался. RSC-only page не считается primary success. Independent проверка ещё требуется. Статус in_progress сохранён; следующий шаг `/verify TASK-006-T2-FT-005-W1`, окончательное решение — /root.
