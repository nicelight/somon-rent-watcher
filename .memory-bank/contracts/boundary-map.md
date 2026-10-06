---
description: Canonical accepted module dependency graph and inline boundary contracts.
status: active
last_verified: 2026-09-02
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/architecture/system-architecture.md
  - .memory-bank/prd.md
---

# Boundary Map

## Modules

| Module / Change Unit | Parent Architecture Unit | Code Root | Responsibility |
|---|---|---|---|
| Composition | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `cmd/somonwatch` | Wire configuration, adapters, CLI and lifecycle. |
| Configuration | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/config` | Parse and validate environment configuration. |
| Polling Application | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/app` | Own scheduling, candidate ordering and durable state-transition decisions. |
| Filtering | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/filter` | Pure settings and match rules. |
| Somon Adapter | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/somon`, `internal/htmlx` | Fetch and parse public HTML with rate/backoff-compatible errors. |
| Persistence Adapter | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/store` | Own SQLite access and transactions. |
| Telegram Adapter | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/telegram` | Own Telegram transport, admin interaction, rendering and delivery. |
| Shared Data | [Somon Watcher runtime](../architecture/system-architecture.md#c4-l2-runtime-composition) | `internal/model` | Passive cross-module data structures only. |

## Dependency Graph

`Consumer -> Provider` means Consumer may depend on Provider only through the linked contract.

| Consumer | Provider | Contract |
|---|---|---|
| Composition | Configuration | [Composition contract](#composition-contract) |
| Composition | Polling Application | [Composition contract](#composition-contract) |
| Composition | Persistence Adapter | [Composition contract](#composition-contract) |
| Composition | Somon Adapter | [Composition contract](#composition-contract) |
| Composition | Telegram Adapter | [Composition contract](#composition-contract) |
| Configuration | Somon Adapter | [Configuration contract](#configuration-contract) |
| Polling Application | Configuration | [Polling orchestration contract](#polling-orchestration-contract) |
| Polling Application | Filtering | [Polling orchestration contract](#polling-orchestration-contract) |
| Polling Application | Somon Adapter | [Polling orchestration contract](#polling-orchestration-contract) |
| Polling Application | Persistence Adapter | [Persistence contract](#persistence-contract) |
| Polling Application | Telegram Adapter | [Telegram application boundary](#telegram-application-boundary) |
| Polling Application | Shared Data | [Shared-data contract](#shared-data-contract) |
| Filtering | Shared Data | [Filtering contract](#filtering-contract) |
| Somon Adapter | Shared Data | [Somon adapter contract](#somon-adapter-contract) |
| Telegram Adapter | Filtering | [Telegram application boundary](#telegram-application-boundary) |
| Telegram Adapter | Shared Data | [Shared-data contract](#shared-data-contract) |
| Telegram Adapter | Somon Adapter | [Telegram application boundary](#telegram-application-boundary) |

No reverse dependency or direct cross-owner state write is authorized.

## Inline Contracts

### Composition contract

Composition loads configuration, creates adapters/application, dispatches CLI commands, and controls start/stop. It does not own filtering, polling, Telegram interaction, or persistence rules.

### Configuration contract

Configuration supplies validated immutable process settings and may reuse the Somon adapter's canonical default category URL. Runtime filter/settings ownership remains in Polling Application through Persistence Adapter.

### Polling orchestration contract

Polling Application is the sole owner of poll phases, shared detail-cap consumption, exact-before-fallback decisions, manual-scan single-flight, recovery/backoff, runtime status, and requests to persist seen/snapshot state. Filtering supplies pure predicates; Somon supplies cards/ads/errors; Telegram supplies delivery outcomes.

### Filtering contract

Filtering validates/normalizes `Settings` and evaluates supplied `Card`/`Ad` values without I/O, persistence, scheduling or delivery side effects. A feature-local price-bound variant may be expressed through settings/value helpers but does not create new durable configuration.

### Somon adapter contract

Somon Adapter owns sequential/rate-limited HTTP and deterministic HTML-to-model conversion, including typed block/status errors. It does not decide business eligibility or seen state.

### Persistence contract

Persistence Adapter alone opens/writes SQLite. Polling Application owns when a business transition is requested; Persistence owns serialization, transactions, WAL and private paths. This delta does not change schema or reinterpret stored rows.

### Telegram application boundary

Telegram Adapter authorizes supported command/callback contexts, acknowledges callbacks, owns pending input, sends append-only menus/status, and returns delivery success/error. It may render normalized Somon text, but it does not choose poll candidates or write SQLite directly.

### Shared-data contract

Shared Data contains passive structures and no I/O, orchestration, mutable singleton or owner-specific business operations.

## External boundary rules

- Somon: outbound public HTML only; sequential requests, configured delay/timeout/body limit, sanity failure without state advancement, and existing 403/429 backoff.
- Telegram: long polling plus callback acknowledgement/message/photo output; admin allowlist and target-chat constraints remain mandatory; callback UI is append-only.
- Production: GitHub → read-only preflight → production git/build/doctor → update only the detected watcher runtime → postflight. Secrets are never printed.

## Update Rules

- Module identity is functional and stable; code paths are discovery roots, not hard task boundaries.
- Add an edge only with an exact contract heading and accepted need.
- Preserve one writer for each mutable invariant; a shared SQLite file does not grant shared write authority.
- Keep feature/task-specific details in feature/task artifacts and do not duplicate a feature subgraph here.

## Keyword monitoring contract proposal

Дельта к [PRD](../prd.md#keyword-monitoring-proposal--2026-10-06); ранее принятые
контракты аренды сохраняются.

- Приложение владеет созданием, настройкой и удалением поиска по стабильному ID.
  Обе кнопки удаления вызывают одну операцию; остальные поиски и аренда сохраняются.
- Somon adapter принимает фразу, категорию и географию; цена проверяется локально.
  Результат — primary cards в выбранном scope, подтверждённая пустая выдача либо ошибка.
  URL: allowlisted category base («Все категории» = `/search/`) + выбранный city
  path segment либо вся страна; `q` кодируется через `net/url`,
  `ordering=relevance` — подтверждённый порядок. [Source evidence](current-integrations.md#keyword-search-source-observations)
  подтверждает совместный scope; global page + local city filter его не заменяет.
  Начальный catalog: `/search/`, `/vse-dlya-doma/mebel/`,
  `/vse-dlya-doma/mebel/mebel-dlya-kuhni/stolyi-stulya/`, `/biznes-i-uslugi/`;
  city: вся страна либо `dushanbe`, `vose`, `dangara`. Произвольный URL не принимается.
  Primary results отделяются от рекомендаций других регионов: на реальных страницах
  заголовок «Объявления из других регионов» есть даже при нуле локальных результатов.
  Ноль карточек допустим только с подтверждённым empty marker; иначе это parse error.
- Штатный Somon matching определяет слово/фразу; локальный фильтр строго проверяет
  цену от/до в сомони без fallback и исключает неизвестную/договорную цену при границе.
  Квартирные фильтры и sanity-пороги нельзя применять к товарной выдаче.
  Неизвестная цена не превращается в ноль; источник определяет доступные поля.
- Telegram сохраняет авторизацию и append-only UI; callbacks и pending input
  адресуют конкретный поиск. Категория/город выбираются читаемыми кнопками из
  [начального catalog](../prd.md#monitoring-rules); смена категории сохраняет город.
  Удалённый ID не создаётся заново из старого меню.
- До доставки приложение проверяет существование/актуальность поиска; уже
  отправленный Telegram запрос отозвать нельзя. SQLite пишет только store.

Проверки: независимость двух поисков с одним ad ID, первый запуск, оба места удаления,
restart, неизменность аренды, авторизация, общий cap/backoff и ошибки доставки.
Фикстуры проверяют пустую/малую локальную выдачу с чужими рекомендациями;
существующие rental/price fixtures остаются regression proof.

### Keyword search boundary shapes

Новый search payload содержит стабильный ID, phrase, category key, city key,
необязательные integer price min/max, enabled и revision. Phrase после trim MUST
быть непустой; значения catalog MUST быть allowlisted. Название уведомления — phrase.
Search создаётся выключенным; включение следует после сводки. Search values не
используют apartment `Settings`; отсутствие границы отличается от нулевой границы.

Somon принимает phrase/category/city, возвращает primary cards с ID, URL, title,
известной ценой/валютой, city и доступным image URL; details добавляют описание.
Пассивные поля MUST сохранять отсутствие цены/города/фото без выдумывания данных.
Price для строгого фильтра считается известной только при подтверждённых сомони;
неизвестная валюта MUST NOT превращаться в сомони. Без price bounds цена не
ограничивает eligibility. Matching phrase MUST NOT подменяться локальным substring.
Filter принимает только search bounds и доступные monetary values, без I/O и
квартирных правил. Catalog/source owner — Somon Adapter; passive payloads — Shared Data.

Telegram передаёт операции list/create/read/update/enable/delete приложению.
Приложение MUST сохранять ID/revision и проверять актуальность; transport MUST NOT
писать SQLite. Pending state задаёт admin/chat/search/action; адресованный input
MUST NOT изменять другой search или арендные settings. Категория и city независимы.
Отсутствующий ID возвращает missing result и свежий список, без upsert/resurrection.
Эта feature добавляет собственные callback routes; переделка rental callbacks,
manual-scan UX и выполнение незавершённых FT-002 claims в неё не входят.

Verification targets: FT-005-AC-001…006; temporary SQLite и существующие
httptest Telegram/Somon harness. Payload symbols/методы выбираются исполнителем
в указанных owners; изменение существующих consumers должно оставаться совместимым.
