---
description: Evidence-backed current-state integration and observed dependency contracts; non-authoritative for target design.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-01
source_of_truth:
  - internal/somon/client.go
  - internal/telegram/client.go
  - internal/store/sqlite_cgo.go
  - deploy/somonwatch.service
---

# Current integrations and observed boundaries

## Scope warning

Здесь зафиксировано то, что текущий код делает на границах. Эти rows не являются allowed target edges и не заменяют accepted [.memory-bank/contracts/boundary-map.md](boundary-map.md).

## External boundaries

### Somon HTML over HTTPS

- Consumer: current `internal/app` through `internal/somon.Client`.
- Provider: configurable category URL and `/adv/...` detail pages on Somon.tj; recovery uses known room-specific category paths.
- Requests: sequential HTTP GET, static User-Agent, optional Referer, shared client/keep-alive, configured minimum delay, timeout and maximum response body.
- Inputs: server-rendered HTML; deterministic Next.js/RSC data is a fallback/enrichment source. Visible DOM controls city-feed membership/order when visible cards exist.
- Price extraction uses structured price metadata or dedicated visible price nodes (including `SidebarPrice`), then standalone price-only lines (including an explicit `Цена:` label). Whole-card/page text must not join image counts or address numbers to the price. Regression coverage: `TestParseCategoryPriceExcludesPhotoCount` and `TestParseDetailPriceExcludesAddressNumber` in `internal/somon/parser_test.go`.
- Failures: non-2xx becomes typed `HTTPError`; 403/429 or detected block page enter blocked handling; parse/sanity failures return raw body to the caller for private diagnostics and do not authorize state advancement.
- No current interaction with Somon private `/api`, `/author/`, pagination, browser automation, proxy rotation or CAPTCHA bypass.
- Evidence: [internal/somon/client.go](../../internal/somon/client.go), [internal/somon/parser.go](../../internal/somon/parser.go), [internal/app/app.go](../../internal/app/app.go).

### Telegram Bot API

- Consumer: current `internal/telegram` client/bot.
- Provider: configurable Telegram-compatible API base, default `https://api.telegram.org`.
- Requests: form-encoded POST to `getMe`, `deleteWebhook`, `getWebhookInfo`, `getChat`, `getUpdates`, `sendMessage`, `sendPhoto`, `editMessageText`, `answerCallbackQuery`.
- Input topology: `getUpdates` long polling only; `run` removes a stale webhook without dropping queued updates.
- Authorization: only users in the configured admin allowlist can process commands/callbacks, either privately or in the configured target chat; other users and other groups are ignored. All administrators share one filter and polling range, while pending text input is isolated per administrator and chat.
- Scheduler controls: interval input persists an inclusive minute range used for each next normal random delay; an accepted immediate-scan callback binds its chat/message to the single scheduler, exposes a temporary busy button and is rejected while another manual scan or HTTP backoff is active. Completion restores the button and sends an explicit zero-result/paused/error message when applicable.
- Recovery notifications: loss of ordinary-ID overlap still starts the room-specific recovery sweep but is log-only; administrators receive a private gap warning only when the last successful poll exceeds the adaptive time threshold.
- Delivery: remote photo is preferred. Caption descriptions start each emoji-led field on a separate line for readability. A deterministic Telegram 400 photo validation error falls back to text; ambiguous/network/server failure does not immediately fall back to avoid an immediate duplicate.
- Durability: the next Telegram update offset is written to SQLite only after processing an update attempt.
- Token handling: token is embedded in the Bot API URL internally; client-side transport errors replace the full base URL with `<telegram-bot-api>` before returning the message.
- Evidence: [internal/telegram/client.go](../../internal/telegram/client.go), [internal/telegram/bot.go](../../internal/telegram/bot.go), [internal/telegram/render.go](../../internal/telegram/render.go).

### System SQLite through CGO

- Consumer: `internal/store` and, through it, `internal/app`/CLI diagnostics.
- Provider: system `libsqlite3`, linked with `-lsqlite3`; CGO-disabled builds expose the same Go methods but return a clear unsupported error.
- Concurrency: one in-process mutex serializes access to a `SQLITE_OPEN_FULLMUTEX` handle; busy timeout is 5 seconds.
- Durability: WAL, `synchronous=NORMAL`, foreign keys and transactional multi-row writes.
- Privacy: parent directory is created with `0700`; DB file is forced to `0600`.
- Schema and writers: [.memory-bank/states/runtime-lifecycle.md#persisted-current-state](../states/runtime-lifecycle.md#persisted-current-state).
- Evidence: [internal/store/sqlite_cgo.go](../../internal/store/sqlite_cgo.go), [internal/store/sqlite_nocgo.go](../../internal/store/sqlite_nocgo.go).

### Host process, environment and filesystem

- `deploy/somonwatch.service` starts `/opt/somonwatch/somonwatch run` as user/group `somonwatch`, reads `/etc/somonwatch/somonwatch.env` and grants writes only to `/var/lib/somonwatch` under `ProtectSystem=strict`.
- The env template owns required Telegram identifiers/secrets and initial Somon/polling defaults; `internal/config` supplies defaults and validates ranges. The persisted Telegram setting owns the runtime normal-poll range after initialization.
- Install/backup scripts write only the documented `/opt`, `/etc`, `/var/lib`, `/var/backups` paths and systemd unit; install refuses to replace an active service.
- Evidence: [deploy/somonwatch.env.example](../../deploy/somonwatch.env.example), [deploy/somonwatch.service](../../deploy/somonwatch.service), [internal/config/config.go](../../internal/config/config.go), [scripts/install-almalinux.sh](../../scripts/install-almalinux.sh), [scripts/backup-installed.sh](../../scripts/backup-installed.sh).

## Observed internal dependency map

This map comes from imports/calls in the current source tree. It is not an accepted architecture graph.

| Current consumer | Current providers used |
|---|---|
| `cmd/somonwatch` | `app`, `config`, `somon`, `store`, `telegram` |
| `internal/app` | `config`, `filter`, `model`, `somon`, `store`, `telegram` |
| `internal/config` | `somon` for the default category URL |
| `internal/filter` | `model` |
| `internal/somon` | `htmlx`, `model` |
| `internal/store` | system SQLite C API |
| `internal/telegram` | `filter`, `model`, `somon` (`NormalizeText` in rendering) |
| `internal/htmlx`, `internal/model` | standard library only |

## Needs verification

- Live compatibility with the current Somon DOM/RSC payload and Telegram account/chat configuration is not proven by repository reads; run `somonwatch doctor` on the intended host.
- The current published Somon rules/robots state is external and time-sensitive; re-check before sustained production use.
- No accepted versioning/compatibility policy exists yet for internal packages, SQLite schema or stored settings JSON; current code has no explicit migration framework.

## Keyword search source observations

Публичные страницы проверены 2026-10-06:
[Все объявления](https://somon.tj/search/) — реальный путь общего поиска;
[Мебель Душанбе](https://somon.tj/vse-dlya-doma/mebel/dushanbe/) — раздел с городом.
[Мебель](https://somon.tj/vse-dlya-doma/mebel/) показывает город и цену от/до;
[Услуги](https://somon.tj/biznes-i-uslugi/) и [Отдам даром](https://somon.tj/otdam-darom/)
— отдельные разделы. Универсальный sale/service/rent не подтверждён.

[Потери и находки Дангара](https://somon.tj/odezhda-i-obuv/poteri-i-nahodki/dangara/):
H1 «Потери и находки Дангара 0», далее «Объявления из других регионов» и карточки.
[Столы Восе](https://somon.tj/vse-dlya-doma/mebel/ofisnaya-mebel/stolyi/vose/):
одна локальная карточка, затем тот же заголовок и рекомендации других регионов.
Поэтому наличие карточек само по себе не подтверждает локальные совпадения.

Оператор прислал URL штатной формы:
[поиск «стол»](https://somon.tj/search/?q=%D1%81%D1%82%D0%BE%D0%BB).
Подтверждён keyword contract: `/search/`, URL-encoded query `q`.
Цена проверяется локально, поэтому price query parameters не требуются для дизайна.
Прямой GET ранее вернул 403; это не отменяет подтверждённый пользователем URL.

Оператор подтвердил штатные ссылки с совместным scope + keyword:
[Все категории, Душанбе](https://somon.tj/search/dushanbe/?q=%D1%81%D1%82%D0%BE%D0%BB&ordering=relevance),
[Столы и стулья, Душанбе](https://somon.tj/vse-dlya-doma/mebel/mebel-dlya-kuhni/stolyi-stulya/dushanbe/?ordering=relevance&q=%D1%81%D1%82%D0%BE%D0%BB),
[Мебель, вся страна](https://somon.tj/vse-dlya-doma/mebel/?q=%D1%81%D1%82%D0%BE%D0%BB).
Подтверждены category/city path scope, `q` и `ordering=relevance`; source blocker закрыт.
Raw HTML/fixtures непустой и нулевой выдачи — execution parser proof; параметры date sort
или иные query fields не подтверждены и не добавляются.

### Keyword source implementation routing

`internal/somon/keyword_search.go` реализует source-owned `KeywordCategories`,
`KeywordCities`, `KeywordSearchURL`, `FetchKeywordSearch` и `ParseKeywordSearch`.
Стабильные keys: `all`, `furniture`, `tables_chairs`, `services`;
`country`, `dushanbe`, `vose`, `dangara`. Native paths соответствуют контракту;
смена category не заменяет city. `internal/model/keyword_search.go` содержит
пассивный `KeywordSearch`; `Card` дополнен `Currency`/`City` без rental-правил.

Keyword parser использует видимый primary DOM до заголовка других регионов/городов.
H1 со счётчиком `0` подтверждает пустую выдачу; карточки при таком счётчике дают ошибку.
RSC-only выдача считается неподтверждённой: её крупнейший массив не доказывает scope
и может содержать рекомендации. Rental parser/client не изменены.
Цена не получает валюту TJS без явного source подтверждения; отсутствующие city/photo
остаются пустыми. Known initial-catalog city names читаются из видимых строк карточки.

Локальные fixtures `testdata/keyword-search-primary-{small,empty}.html` —
репрезентативный sanitized HTML по документированным наблюдениям H1/region heading
и существующим card/RSC patterns, не live captures. Они не доказывают актуальную
полную DOM-совместимость Somon. Execution evidence TASK-006 хранится в
[протоколе](../../.protocols/TASK-006-T2-FT-005-W1/handoff.md);
independent `/verify` ещё требуется, production не изменялся.
