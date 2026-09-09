# Семантическая проверка — TASK-005-T3-FT-004-W1

Проверена финальная attempt 2: commit `93cbe9ebcfe6421040461f4ee4e049430bf5028d`, исправляющий завышенные цены и сохраняющий поддержанный формат `Цена: 4 200 c.`, установлен в существующий watcher. Основание: indexed task и непосредственно связанные integration/runbook/lifecycle/invariant документы. Незавершённая продуктовая очередь не входит в maintenance hotfix.

Повторная проверка ограничена delta относительно ранее проверенного `d5fdc12` и новыми runtime receipts. `domPrice` снимает только явный префикс `Цена:`, после чего применяет прежний anchored price-only match. Full-page fallback не восстановлен. Assertion существующей fixture доказала RED (`nil`) до исправления и GREEN (4200) после него; полный local/target native gate и replay 60 category/38 detail prices прошли. Source bytes основной копии совпадают с чистым release checkout.

Просмотрены изменения release procedure: новый baseline-version guard и помещение stage copy/chmod в существующий cleanup try. Финальный журнал подтверждает exact commit, target doctor до stop и scoped install/start. Свежий runtime probe связывает running/installed/built SHA `c99616820cc979fb6257cfa959dfcc323b4b9ea487ecf92719ec548c8503e9d8`, сохраняет все 3786 исходных seen rows среди 3792, settings/env, SQLite integrity и zero restarts; rollback binary проверен. Первый штатный poll завершился с одной успешной отправкой.

Независимое сравнение before/after JSON подтвердило все девять стабильных полей: containers, services, listeners, routes, firewall rules, SELinux, failed units, settings, environment. Firewall fingerprint нормализован до и после финальной установки. Невалидный исходный hash attempt 1 остаётся только историческим evidence и не подменяет это сравнение.

Два fresh co-review Codex Luna xhigh проверили только изменённые focuses: labeled-price compatibility и final deployment/state/isolation. Оба — APPROVE, без существенных замечаний. Итоговый Reviewer подтвердил устранение принятой compatibility regression и достаточность финальных receipts.

SEMANTIC_VERDICT: semantic-pass

Evidence: `.protocols/TASK-005-T3-FT-004-W1/verification.md`, `labeled-price-red.log`, `labeled-price-green.log`, `local-build.log`, `offline-prices-green.log`, `production-release.log`, `production-poll.log`, `production-before.json`, `production-after.json`, `host-comparison.json`, `final-runtime-proof.json`, соответствующие probe/release scripts. `attempt-1/` — только supporting history.

Действие explicit owner: оформить closure вместе с functional PASS и выполнить необходимый mb-sync. Reviewer не менял task status, код или production.
