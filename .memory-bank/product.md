---
description: Product identity and accepted target scope for Somon Rent Watcher.
status: active
last_updated: 2026-10-06
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/analysis/product-brief.md
---

# Product — Somon Rent Watcher

## Identity

Somon Rent Watcher is a small Go service for one configured Telegram group and its administrators. It watches apartment rentals with a shared rental filter and supports independent saved keyword searches on Somon.tj, delivering matching listings without normal duplicates within each search.

## Core value

- Replace repeated manual browsing with timely filtered notifications.
- Keep strict matches primary while surfacing a small, bounded set of close higher-price alternatives when a completed poll is otherwise empty.
- Keep administrator actions visible across Telegram clients through append-only feedback.
- Find goods or services within each saved search's geography and strict budget, starting with available matches and continuing with new listings.

## Audience

- Configured administrators who manage the shared rental filter, saved searches, and scheduler controls.
- Members of the configured group who consume ad notifications.

## Primary flow

1. One process polls Somon category HTML and validates the feed.
2. Baseline, pause, first-seen, recovery, backoff, and request-cap behavior protect delivery and state.
3. Fresh exact candidates are evaluated first.
4. Only when exact evaluation completes with no match, up to three closest fresh higher-price ads within `150%` of `PriceMax` may be sent.
5. Authorized Telegram callbacks are acknowledged and their current result/menu is sent as a new message.
6. Verified changes are published and deployed in an isolated sequence affecting only Somon Rent Watcher.
7. An administrator creates a keyword search, checks its summary, and enables it. The service queries its native category/geography scope, delivers existing matching listings first, and then new listings; settings and deletion address that search.

## Constraints

- One process, one SQLite file, one shared rental filter, independent saved searches, one target group, and an admin allowlist.
- The rental upgrade adds no database schema, dependency, setting, worker, inbound port, browser automation, block bypass, or client-specific behavior. Keyword searches permit the accepted additive storage while reusing existing runtime owners and limits.
- A single detail-request cap covers exact and fallback work.
- Production settings/seen history must survive upgrades; unrelated workloads must remain unchanged.

## Non-goals

- Full crawler/archive, price history, apartment-level deduplication, multiple target groups, web UI, NLP/LLM, or private Somon APIs.
- Separate shown-history, old-ID backfill, configurable fallback policy, Telegram message-edit retries, or old-menu cleanup.
- Changes to unrelated production containers, services, networks, firewall, routing, reverse proxy, SELinux, or databases.

## Accepted keyword-monitoring delta

Помимо аренды приняты несколько независимых сохраняемых поисков по слову/фразе,
с необязательными границами цены, категорией и собственной географией. Бюджет строгий;
штатный Somon matching определяет слово/фразу. При границе неизвестная/договорная цена
и цена не в сомони исключаются, без конвертации и rental fallback.

Сначала доставляются доступные существующие совпадения, затем новые. Удаление есть
в списке и настройках конкретного поиска. Доставленные ID не повторяются в этом поиске
после restart, правки или снижения цены; ранее отклонённые можно переоценить после правки.
Ошибки и лимиты сохраняют retry, а данные аренды и других поисков остаются независимыми.

- [PRD](prd.md#keyword-monitoring-proposal--2026-10-06): принятые правила и примеры.
- [EP-002](epics/EP-002-keyword-monitoring.md): ценность и критерии новой возможности.
- [FT-005](features/FT-005-keyword-monitoring.md): наблюдаемое поведение и acceptance.
- [Source evidence](contracts/current-integrations.md#keyword-search-source-observations): `q` подтверждён; native category/city scope вместе с `q` требует подтверждения перед source-dependent задачами. Global first-page с local city filter не заменяет scoped search.

## Current-state evidence

- [Architecture baseline](architecture/system-architecture.md)
- [Integration baseline](contracts/current-integrations.md)
- [Runtime lifecycle baseline](states/runtime-lifecycle.md)
- [Current test coverage](testing/current-coverage.md)
- [AlmaLinux operations route](runbooks/almalinux-9-operations.md)
