---
description: Status of accepted global invariants and routing to descriptive current-state guardrails.
status: draft
last_verified: 2026-09-01
---

# Invariants

## Authority status

Authoritative PRD и accepted Global Backbone отсутствуют. Поэтому brownfield-маппинг не превращает наблюдаемое поведение кода в новые нормативные `MUST/NEVER` rules.

## Accepted MUST

- Не определены владельцем product/spec workflow.

## Accepted NEVER

- Не определены владельцем product/spec workflow.

## Descriptive current-state guardrails

Текущая реализация фактически:

- не отправляет initial baseline;
- на паузе продвигает seen-baseline без backfill;
- не меняет seen/snapshot при category parse/sanity failure;
- помечает подходящий ID seen только после подтверждённой Telegram delivery;
- не делает обход 403/429, proxy rotation или CAPTCHA bypass;
- ограничивает detail requests на poll и сохраняет debug HTML с private file mode.

Это as-is observations, не target authority. Подробные переходы и исключения: [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md): current lifecycle evidence.

## Evidence

- [internal/app/app.go](../internal/app/app.go): writers and guard conditions.
- [internal/telegram/bot.go](../internal/telegram/bot.go): delivery/fallback behavior and admin-only control.
- [internal/store/sqlite_cgo.go](../internal/store/sqlite_cgo.go): private SQLite permissions and transactional writes.
