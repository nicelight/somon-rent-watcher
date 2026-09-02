---
description: Requirements registry status and current evidence routing; no authoritative PRD is present.
status: draft
last_verified: 2026-09-01
---

# Requirements

## Authority status

- `.memory-bank/prd.md` отсутствует.
- Формальные `REQ-*`, RTM, epics, features и product task records не создавались во время brownfield-маппинга.
- [docs/TZ_Somon_Rent_Watcher.md](../docs/TZ_Somon_Rent_Watcher.md) имеет собственный статус `draft v1.1` и синхронизирован с реализованным локальным MVP; это current-state input, но не замена принятому PRD.
- Реализованное поведение зарегистрировано как current-state evidence в [.memory-bank/product.md](product.md), [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md) и [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md).

## REQ list

Не инициализирован без authoritative PRD/delta.

## Out of scope

Будет определено владельцем PRD; current implementation non-goals перечислены только описательно в [.memory-bank/product.md#current-non-goals](product.md#current-non-goals).

## Traceability

RTM не инициализирована: превращать существующий код или draft-ТЗ в принятые product requirements без решения оператора запрещено.

## Next route

После передачи operator intent/PRD/delta использовать `/brief` → `/write-prd` либо другой явно выбранный product workflow; brownfield baseline остаётся входным evidence.
