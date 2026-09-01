# Local Docker runbook (Kubuntu)

Этот вариант запускает один постоянный контейнер `somon-rent-watcher`. Входящие порты, reverse proxy и отдельный database container не нужны: сервис использует исходящий HTTPS и SQLite в named volume.

## 1. Telegram и конфигурация

Скопируйте template, если локального `.env` ещё нет:

```bash
cp deploy/somonwatch.env.example .env
chmod 600 .env
```

Обязательные значения:

```text
TELEGRAM_BOT_TOKEN=...
TELEGRAM_ADMIN_USER_IDS=123456789,987654321
TELEGRAM_TARGET_CHAT_ID=-1001234567890
```

Все ID администраторов должны быть положительными и разделяться запятыми. Они управляют одним общим фильтром. Команды принимаются в личных чатах и только в `TELEGRAM_TARGET_CHAT_ID`; остальные пользователи и группы игнорируются.

Каждый администратор должен сначала открыть личный чат с ботом и нажать Start, иначе Telegram не позволит боту присылать ему технические уведомления. Для ввода цены и минус-слов в группе отключите Group Privacy через `@BotFather` → `/setprivacy` → выбрать бота → `Disable`.

Не публикуйте `.env`, token или вывод команд, раскрывающих environment. `.env` исключён из Git и Docker build context.

## 2. Сборка и запуск

```bash
docker compose config --quiet
docker compose up --build -d
docker compose ps
docker compose logs --tail=100 somonwatch
```

Docker build выполняет `gofmt` check, `go test`, `go vet`, CGO build и linkage check. Конечный container работает от непривилегированного пользователя, имеет read-only root filesystem, не публикует ports и пишет только в `/var/lib/somonwatch`.

## 3. Проверка

После запуска выполните doctor внутри уже работающего постоянного контейнера:

```bash
docker compose exec somonwatch somonwatch doctor
```

Doctor проверяет SQLite, Telegram bot/target chat и live Somon parser. Команда не создаёт baseline и не меняет seen IDs, но основной процесс к этому моменту уже мог выполнить первый poll.

Первый успешный poll создаёт baseline без отправки текущих объявлений. Общий мониторинг остаётся на паузе, пока один из администраторов не настроит `/filter` и не включит его.

## 4. Управление

```bash
docker compose logs -f somonwatch
docker compose restart somonwatch
docker compose stop somonwatch
docker compose start somonwatch
docker compose down
```

`docker compose down` сохраняет named volume и SQLite. Не используйте `docker compose down -v`, если не хотите безвозвратно удалить локальную базу, настройки и seen history.
