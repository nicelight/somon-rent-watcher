---
description: Независимые keyword-поиски Somon с собственным бюджетом, географией и доставкой в существующую Telegram-группу.
status: active
lifecycle: verified
last_updated: 2026-10-07
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
---

# EP-002 — Мониторинг по ключевым словам

## Value

Участники группы находят нужные товары или услуги без повторного ручного поиска;
каждый поиск соблюдает собственные географию и строгий бюджет, не меняя аренду.

## Sources and requirements

- [PRD](../prd.md#keyword-monitoring-proposal--2026-10-06): правила приняты оператором 2026-10-06.
- [REQ-010…014](../requirements.md#accepted-keyword-monitoring-delta): lifecycle, matching, доставка, надёжность и сохранность.
- [FT-005](../features/FT-005-keyword-monitoring.md): полный пользовательский цикл поиска.

## Success metrics

- Включённый поиск доставляет существующие и новые совпадения в своём scope и бюджете.
- Один доставленный ID не повторяется в том же поиске; история другого поиска независима.
- Оба места удаления убирают выбранный поиск, сохраняя остальные данные и аренду.
- Ошибки/лимиты сохраняют retry и общие ограничения сервиса.

## Acceptance criteria

Все FT-005-AC-001…010 подтверждены объявленными методами; прежние rental settings,
история и поведение сохраняются. Product scope и native scope wire mapping закрыты.

## Planning dependency

[Source evidence](../contracts/current-integrations.md#keyword-search-source-observations):
штатные category/city URLs с `q`/`ordering=relevance` подтверждены оператором.
Product decomposition reviewed `APPROVE`; [task plan](../tasks/plans/IMPL-FT-005.md)
reviewed APPROVE для Planning Revision 1. TASK006…010 done, independent functional
проверки и T3 gates пройдены; [feature semantic-pass](../features/FT-005-keyword-monitoring.md#semantic-verification)
подтверждает локальную реализацию. Первоначальная разработка была локальной; разрешённый оператором production acceptance TASK013 W6 завершён для всей версии8b48c4c с сохранением данных и окружения (REQ-008/014). Все TASK006…013 done, функциональные и семантические gates пройдены.
