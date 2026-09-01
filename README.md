# Central Asia Rent Watcher

Небольшой Go-сервис для отслеживания новых объявлений об аренде жилья на популярных платформах Центральной Азии и отправки подходящих вариантов в Telegram.

Проект ориентирован на бережный low-rate monitoring публичной HTML-выдачи. Текущая версия поставляется с одним platform-specific adapter; подключение других площадок требует отдельного parser/adapter и собственных HTML fixtures.

## Что реализовано

- first-seen semantics: один ID штатно отправляется не более одного раза;
- первый запуск создаёт baseline без спама текущими объявлениями;
- monitoring по умолчанию стоит на паузе до настройки фильтра;
- на паузе свежие ID запоминаются без последующего backfill;
- случайный polling раз в 10–30 минут;
- двухступенчатая обработка: дешёвый card-prefilter и detail-page только для кандидатов;
- ограничение detail-запросов за один poll;
- recovery sweep по выбранным категориям при подозрении на разрыв ленты;
- фильтры по цене, комнатам, этажу, активности автора, словам/фразам и типу объявления;
- SQLite для seen IDs, настроек и служебного состояния;
- длинный backoff на HTTP 403/429 без обхода блокировок;
- sanity-check HTML и private diagnostic snapshots при изменении страницы;
- Telegram long polling без webhook, reverse proxy и входящего порта;
- один бинарник и один hardened `systemd` unit.

## Стек

- Go 1.21+;
- стандартная библиотека Go;
- system SQLite через CGO (`libsqlite3`);
- прямой Telegram Bot API;
- встроенный tolerant HTML tree parser;
- systemd.

В проекте нет сторонних Go modules и сетевого скачивания зависимостей при сборке.

## Быстрый старт

На AlmaLinux/RHEL-подобной системе:

```bash
sudo dnf install golang gcc sqlite sqlite-devel
./scripts/build.sh
./dist/<binary> version
```

Build script выполняет formatting check, package tests, `go vet`, CGO build и создаёт checksum бинарника.

Перед production-запуском заполните Telegram credentials и platform URL в environment template, затем обязательно выполните:

```bash
./dist/<binary> doctor
```

`doctor` проверяет SQLite, Telegram и live HTML parser, но не создаёт seen-baseline.

## CLI

```text
<binary> run      # основной сервис; команда по умолчанию
<binary> doctor   # SQLite + Telegram + live parser; baseline не меняется
<binary> ids      # вывод Telegram user_id/chat_id входящих сообщений
<binary> version
<binary> help
```

Подробный порядок безопасной установки на AlmaLinux 9 находится в [production runbook](docs/RUNBOOK_ALMALINUX_9.md).

## Telegram-фильтр

Администратор управляет сервисом только в личном чате с ботом:

- `/filter` — открыть меню фильтра;
- `/status` — показать последний/следующий poll, режим, ошибки и число seen IDs;
- кнопка включения/паузы — управляет рассылкой, но не останавливает продвижение baseline.

Объявления отправляются в одну настроенную Telegram-группу. Изменение фильтра действует только на новые ID и не запускает retroactive backfill.

## Поведение свежести

«Новое» означает: ID впервые замечен локальным watcher после первоначальной baseline-инициализации.

- известный ID после поднятия повторно не отправляется;
- изменение цены известного ID не создаёт новое уведомление;
- изменение фильтра не пересматривает старые ID;
- ранее невидимый ID может быть отправлен один раз при первом попадании в наблюдаемую выдачу.

Это fresh-only monitor, а не полный archive/crawler.

## Безопасность эксплуатации

- Telegram token хранится только в root-readable env-файле;
- SQLite и diagnostic HTML используют private permissions;
- `systemd` unit работает от отдельного непривилегированного пользователя;
- процессу не нужны capabilities, устройства, home directories или inbound ports;
- proxy rotation, CAPTCHA bypass и иные механизмы обхода ограничений не реализованы.

## Ответственное использование

Автоматизированный доступ может ограничиваться robots.txt, Terms of Service и применимым законодательством конкретной площадки. Перед постоянной эксплуатацией получите необходимое разрешение или отдельно оцените и примите этот риск. Наличие технической возможности не означает права на сбор или распространение данных.

Проект намеренно использует низкую частоту запросов, не обращается к private/internal APIs и прекращает активный polling при признаках блокировки.

## Структура

```text
cmd/                 CLI composition root
internal/app/        scheduler и processing pipeline
internal/config/     environment configuration
internal/filter/     фильтры и settings validation
internal/htmlx/      tolerant HTML tree
internal/model/      shared data structures
internal/store/      SQLite CGO wrapper
internal/telegram/   Bot API, menu и notifications
internal/<adapter>/  platform HTTP client и parser
deploy/              systemd и environment template
scripts/             build, preflight, install и backup
docs/                technical notes и production runbook
testdata/            HTML fixtures
```

## Проверка

```bash
./scripts/build.sh
```

Existing tests cover baseline/no-spam behavior, pause without backfill, one-time delivery, filter boundaries, HTML parser regressions, SQLite persistence and Telegram delivery fallback.

## License

[MIT](LICENSE)
