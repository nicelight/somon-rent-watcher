---
description: Product identity and accepted target scope for Somon Rent Watcher.
status: active
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/analysis/product-brief.md
---

# Product — Somon Rent Watcher

## Identity

Somon Rent Watcher is a small Go service for one configured Telegram group and its administrators. It watches fresh apartment-rental listings on Somon.tj, applies one shared filter, and sends useful first-seen results without normal duplicates.

## Core value

- Replace repeated manual browsing with timely filtered notifications.
- Keep strict matches primary while surfacing a small, bounded set of close higher-price alternatives when a completed poll is otherwise empty.
- Keep administrator actions visible across Telegram clients through append-only feedback.

## Audience

- Configured administrators who manage one shared filter and scheduler controls.
- Members of the configured group who consume ad notifications.

## Primary flow

1. One process polls Somon category HTML and validates the feed.
2. Baseline, pause, first-seen, recovery, backoff, and request-cap behavior protect delivery and state.
3. Fresh exact candidates are evaluated first.
4. Only when exact evaluation completes with no match, up to three closest fresh higher-price ads within `150%` of `PriceMax` may be sent.
5. Authorized Telegram callbacks are acknowledged and their current result/menu is sent as a new message.
6. Verified changes are published and deployed in an isolated sequence affecting only Somon Rent Watcher.

## Constraints

- One process, one SQLite file, one shared filter, one target group, and an admin allowlist.
- No new database schema, dependency, setting, worker, inbound port, browser automation, block bypass, or client-specific behavior for this delta.
- A single detail-request cap covers exact and fallback work.
- Production settings/seen history must survive upgrades; unrelated workloads must remain unchanged.

## Non-goals

- Full crawler/archive, price history, apartment-level deduplication, multiple profiles/groups, web UI, NLP/LLM, or private Somon APIs.
- Separate shown-history, old-ID backfill, configurable fallback policy, Telegram message-edit retries, or old-menu cleanup.
- Changes to unrelated production containers, services, networks, firewall, routing, reverse proxy, SELinux, or databases.

## Current-state evidence

- [Architecture baseline](architecture/system-architecture.md)
- [Integration baseline](contracts/current-integrations.md)
- [Runtime lifecycle baseline](states/runtime-lifecycle.md)
- [Current test coverage](testing/current-coverage.md)
- [AlmaLinux operations route](runbooks/almalinux-9-operations.md)
