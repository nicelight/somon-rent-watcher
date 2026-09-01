---
description: C4 L1 current-state product baseline для Somon Rent Watcher.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-01
source_of_truth:
  - README.md
  - internal/app/app.go
---

# Product — current state

## Scope warning

Этот документ описывает только реализованное состояние репозитория на 2026-09-01. Он не является PRD, roadmap или утверждённым target-state решением.

## What this is

Somon Rent Watcher — небольшой Go-сервис для настроенного списка администраторов и одной Telegram-группы. Он периодически читает свежую серверную HTML-выдачу аренды квартир в Душанбе на Somon.tj, применяет общий настраиваемый фильтр и отправляет подходящие впервые замеченные объявления в Telegram.

## Core value

- Сократить ручной просмотр выдачи Somon до уведомлений о новых подходящих объявлениях.
- Не спамить текущей выдачей при первом запуске и не отправлять один известный ID штатно повторно.
- Дать настроенным администраторам управление общим фильтром и состоянием мониторинга через личный чат Telegram или целевую группу.

## Actors and external systems

| Actor / system | Current interaction |
|---|---|
| Настроенный администратор | Настраивает общий фильтр, включает/ставит на паузу мониторинг и смотрит `/status` в личном чате с ботом или целевой группе. |
| Участники целевой группы | Получают подходящие объявления; пользователи вне admin allowlist не могут управлять ботом. |
| Somon.tj | Отдаёт category/detail HTML по исходящему HTTPS; private/internal API и browser automation не используются. |
| Telegram Bot API | Отдаёт updates через long polling и принимает сообщения/фото; webhook и входящий порт не нужны. |

## Primary current flows

1. `somonwatch run` открывает SQLite и параллельно запускает Telegram long polling и Somon polling loop.
2. Первый успешный poll создаёт seen-baseline без рассылки и оставляет мониторинг на паузе.
3. На паузе новые ID продолжают становиться seen без отправки и без последующего backfill.
4. В активном режиме новые карточки проходят дешёвый card-prefilter; detail HTML запрашивается только для кандидатов и в пределах per-poll cap.
5. Прошедшее detail-filter объявление отправляется в одну группу и только после успешной отправки помечается seen.
6. Подозрение на разрыв обычной ленты запускает один recovery sweep по выбранным room pages; 403/429 переводят scheduler в длинный backoff.

## Current implementation constraints

- Один Linux process; production route использует `systemd`, локальный Kubuntu route — один Docker Compose container. Входящие TCP/UDP-порты отсутствуют.
- Go 1.21+, стандартная библиотека Go и system `libsqlite3` через CGO; внешних Go modules нет.
- Один SQLite-файл хранит seen IDs, JSON settings и служебное state.
- Список Telegram admin IDs, один общий filter profile и одна target group; управление из других групп запрещено.
- Никаких browser automation, Redis/PostgreSQL, queue, mandatory Docker, proxy rotation или CAPTCHA bypass.
- Production-эксплуатация scraper имеет явно задокументированный compliance-риск относительно опубликованных правил Somon; актуальность правил требует внешней проверки перед запуском.

## Current non-goals

Текущая версия не является полным crawler/archive и не реализует pagination/lazy-load reverse engineering, историю цен, дедупликацию квартир, несколько профилей/групп, web UI, NLP/LLM, телефон автора или обход блокировок.

## Evidence

- [README.md](../README.md): публичное описание реализованного сервиса, stack и flows.
- [docs/IMPLEMENTATION_NOTES.md](../docs/IMPLEMENTATION_NOTES.md): намеренные safety additions и фактический dependency choice.
- [cmd/somonwatch/main.go](../cmd/somonwatch/main.go): CLI entrypoints `run`, `doctor`, `ids`, `version`, `help`.
- [internal/app/app.go](../internal/app/app.go): scheduler, baseline, pause, filtering, delivery, recovery и backoff.
- [docs/TZ_Somon_Rent_Watcher.md](../docs/TZ_Somon_Rent_Watcher.md): draft pre-implementation intent; использовать как historical product evidence, а не как принятый Memory Bank PRD.

## Routing

- Current architecture: [.memory-bank/architecture/system-architecture.md](architecture/system-architecture.md): C4 as-is topology and component map.
- Current lifecycle: [.memory-bank/states/runtime-lifecycle.md](states/runtime-lifecycle.md): polling, delivery and persisted-state transitions.
- Operations: [.memory-bank/runbooks/almalinux-9-operations.md](runbooks/almalinux-9-operations.md): source runbook routing and safety boundaries.
- Local Docker operations: [.memory-bank/runbooks/docker-local-operations.md](runbooks/docker-local-operations.md): Kubuntu Compose runtime and persistent-state route.
