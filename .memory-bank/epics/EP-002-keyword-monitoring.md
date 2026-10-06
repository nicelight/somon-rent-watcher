---
description: Независимые keyword-поиски Somon с собственным бюджетом, географией и доставкой в существующую Telegram-группу.
status: active
lifecycle: planned
last_updated: 2026-10-06
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

Все FT-005-AC-001…007 подтверждены объявленными методами; прежние rental settings,
история и поведение сохраняются. Product scope и native scope wire mapping закрыты.

## Planning dependency

[Source evidence](../contracts/current-integrations.md#keyword-search-source-observations):
штатные category/city URLs с `q`/`ordering=relevance` подтверждены оператором.
Product decomposition reviewed `APPROVE`; [task plan](../tasks/plans/IMPL-FT-005.md)
создан. Следующий этап — свежий `/review-tasks-plan FT-005`.
