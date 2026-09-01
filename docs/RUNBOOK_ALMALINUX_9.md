# Runbook: Somon Rent Watcher на AlmaLinux 9

**Целевой сервер:** `<server-ip-or-hostname>`  
**Сервис:** `somonwatch.service`  
**Режим:** отдельный непривилегированный host-service, не контейнер  
**Входящие порты:** отсутствуют

Этот порядок специально составлен для production-хоста, где уже работают важные контейнеры и host-сервисы.

---

## 1. Что установка не должна менять

Somon Watcher:

- не добавляется в Docker;
- не перезапускает Docker или containerd;
- не меняет Docker networks;
- не трогает Traefik и его конфигурацию;
- не открывает порты firewalld;
- не меняет маршрутизацию;
- не меняет DNS сервера;
- не отключает SELinux;
- не устанавливает Nginx/Caddy;
- не использует PostgreSQL существующих проектов;
- не занимает `80`, `443`, `8000`, `9000`, `5432`, `19999`, `37525` и другие listening ports.

Ему нужны только исходящие HTTPS-соединения к Somon и Telegram.

Перед установкой нужно заново снять состояние всех существующих containers, host-services, listening sockets и security controls: production мог измениться. Не публикуйте inventory конкретного хоста в открытом репозитории.

---

## 2. Подготовить Telegram

### 2.1. Создать бота

Через `@BotFather`:

1. `/newbot`;
2. сохранить token;
3. добавить бота в целевую группу;
4. дать ему право отправлять сообщения и фотографии;
5. каждому будущему администратору открыть личный чат с ботом и нажать Start;
6. если администраторы будут вводить цену и минус-слова прямо в группе, выполнить `@BotFather` → `/setprivacy` → выбрать бота → `Disable`.

Боту не нужны права администратора, если обычным участникам группы разрешено отправлять сообщения и фотографии.

### 2.2. Получить ID администраторов и группы

Это удобнее сделать после сборки бинарника, но до запуска systemd-service:

```bash
cd /root/somon-rent-watcher
TELEGRAM_BOT_TOKEN='REAL_TOKEN' ./dist/somonwatch ids
```

Затем:

1. каждому администратору отправить `/start` боту в личном чате — строки покажут положительные `user_id`;
2. отправить в группе `/start@BotUsername` — строка покажет отрицательный `chat_id` группы;
3. остановить `ids` через `Ctrl+C`.

Полученные значения:

```text
TELEGRAM_ADMIN_USER_IDS=<user_id_1>,<user_id_2>
TELEGRAM_TARGET_CHAT_ID=<chat_id группы, обычно начинается с -100>
```

`ids` и работающий `somonwatch.service` нельзя запускать одновременно: оба используют Telegram `getUpdates`.

---

## 3. Передача архива

На рабочей машине:

```bash
scp somon-rent-watcher-v1.0.0.zip <admin-user>@<server-ip-or-hostname>:/tmp/
scp somon-rent-watcher-v1.0.0.zip.sha256 <admin-user>@<server-ip-or-hostname>:/tmp/
```

Имя SSH-пользователя намеренно не предполагается. Использовать обычную административную учётную запись сервера.

На сервере:

```bash
ssh <admin-user>@<server-ip-or-hostname>
sudo -i
cd /tmp
sha256sum -c somon-rent-watcher-v1.0.0.zip.sha256
mkdir -p /root/somon-rent-watcher-build
unzip -q somon-rent-watcher-v1.0.0.zip -d /root/somon-rent-watcher-build
cd /root/somon-rent-watcher-build/somon-rent-watcher
```

---

## 4. Снять preflight snapshot

Скрипт только читает состояние:

```bash
./scripts/preflight-almalinux.sh | tee /root/somonwatch-preflight-$(date +%F-%H%M%S).log
```

До продолжения проверить:

```bash
systemctl --failed
docker ps
ss -lntup
firewall-cmd --list-all
free -h
df -h /
getenforce
```

Особенно важно зафиксировать:

- какие контейнеры запущены;
- их health/status;
- существующие listening sockets;
- отсутствие неожиданных failed units;
- запас диска и RAM;
- текущее состояние SELinux/firewalld.

Не продолжать установку, если уже есть проблемы с важными workload. Сначала стабилизировать хост.

---

## 5. Установить только build-зависимости

Проект не скачивает Go-модули из интернета. Для сборки нужны компилятор Go, GCC и системная SQLite.

Сначала посмотреть план транзакции без применения:

```bash
dnf install --assumeno golang gcc sqlite sqlite-devel
```

Проверить, что пакетный менеджер не предлагает удалить или заменить критичные компоненты. Затем выполнить:

```bash
dnf install golang gcc sqlite sqlite-devel
```

Полный `dnf update` ради этого приложения не требуется.

Проверка:

```bash
go version
gcc --version | head -1
sqlite3 --version
test -f /usr/include/sqlite3.h && echo OK
```

Минимальная версия Go — 1.21. После установки обязательно проверить фактическую версию командой `go version`; при более старой версии сборку не продолжать.

---

## 6. Собрать на самом AlmaLinux-хосте

Сборка на целевом сервере исключает несовместимость glibc между разными дистрибутивами:

```bash
cd /root/somon-rent-watcher-build/somon-rent-watcher
./scripts/build.sh
```

Ожидается:

```text
go test ./...       PASS
go vet ./...        PASS
dist/somonwatch     ELF x86-64
libsqlite3.so       found
SHA256              создан
```

Дополнительная проверка:

```bash
./dist/somonwatch version
ldd ./dist/somonwatch
```

В `ldd` не должно быть `not found`.

---

## 7. Установить файлы, но пока не запускать

```bash
sudo ./scripts/install-almalinux.sh
```

Скрипт:

- создаёт только пользователя/группу `somonwatch`;
- создаёт `/opt/somonwatch`, `/etc/somonwatch`, `/var/lib/somonwatch`;
- устанавливает бинарник и systemd unit;
- сохраняет существующий env-файл при upgrade;
- делает `systemctl daemon-reload`;
- **не запускает и не enable-ит сервис без `--start`**;
- не выполняет Docker/firewall/Traefik/SELinux changes.

Проверить установленные файлы:

```bash
ls -la /opt/somonwatch
ls -ld /etc/somonwatch /var/lib/somonwatch
systemctl cat somonwatch.service
```

---

## 8. Заполнить secrets/config

```bash
vi /etc/somonwatch/somonwatch.env
chmod 0600 /etc/somonwatch/somonwatch.env
chown root:root /etc/somonwatch/somonwatch.env
```

Обязательно заменить:

```text
TELEGRAM_BOT_TOKEN=...
TELEGRAM_ADMIN_USER_IDS=...,...
TELEGRAM_TARGET_CHAT_ID=...
```

Production polling оставить:

```text
SOMON_POLL_MIN=10m
SOMON_POLL_MAX=30m
SOMON_GAP_AFTER=45m
SOMON_BLOCK_BACKOFF=2h
SOMON_MAX_DETAILS_PER_POLL=20
```

Основные пути не менять без причины:

```text
DB_PATH=/var/lib/somonwatch/somonwatch.db
DEBUG_DIR=/var/lib/somonwatch/debug
```

Проверить, что placeholders удалены:

```bash
if grep -q 'replace_me' /etc/somonwatch/somonwatch.env; then echo 'PLACEHOLDER LEFT'; else echo OK; fi
```

Не выводить env-файл в общий terminal log и не отправлять его в чат: там token.

---

## 9. Выполнить doctor до запуска

Env-файл совместим с shell; загрузить его только в текущую root-сессию и передать процессу service-user:

```bash
set -a
. /etc/somonwatch/somonwatch.env
set +a
runuser -u somonwatch --preserve-environment -- /opt/somonwatch/somonwatch doctor
unset TELEGRAM_BOT_TOKEN TELEGRAM_ADMIN_USER_IDS TELEGRAM_ADMIN_USER_ID TELEGRAM_TARGET_CHAT_ID
```

Doctor проверяет:

- создание/доступ к SQLite;
- Telegram token через `getMe`;
- доступ к целевой группе через `getChat`;
- отсутствие конфликтующего webhook;
- получение live HTML Somon;
- количество распознанных карточек;
- наличие обычной части выдачи.

Doctor **не создаёт baseline и не помечает объявления seen**.

Если parser нашёл меньше `SOMON_MIN_CARDS=20`, сервис не запускать. Сначала исследовать HTML и debug output.

---

## 10. Первый запуск

```bash
systemctl enable --now somonwatch.service
systemctl --no-pager --full status somonwatch.service
journalctl -u somonwatch.service -n 100 --no-pager
```

Ожидаемое первое поведение:

1. сервис сразу получает основную выдачу;
2. текущие ID записываются в baseline;
3. текущие объявления **не отправляются в группу**;
4. каждому настроенному администратору приходит личное сообщение о создании baseline;
5. мониторинг остаётся на паузе по умолчанию;
6. следующий poll назначается случайно через 10–30 минут.

На паузе сервис продолжает опрашивать ленту и запоминать новые ID без отправки. После включения он не делает backfill объявлений, появившихся во время паузы.

Если текущие объявления пошли в группу при первом запуске, немедленно остановить только этот unit:

```bash
systemctl stop somonwatch.service
journalctl -u somonwatch.service --since -15min --no-pager
```

Не перезапускать Docker/Traefik/firewalld.

---

## 11. Настроить фильтр

В личном чате с ботом или в настроенной целевой группе:

```text
/filter
```

Сначала, не снимая паузу, настроить:

- цену от/до;
- комнаты;
- этажи;
- лимит активных объявлений автора;
- негативные слова;
- обычные/VIP.

Все администраторы изменяют один общий фильтр. Затем в главном меню нажать:

```text
▶ Включить мониторинг
```

Проверка состояния:

```text
/status
```

Изменение фильтра действует только на новые ID. Старые объявления не переотправляются. Кнопка `⏸ Поставить на паузу` прекращает рассылку, но fresh-baseline продолжает двигаться вперёд.

---

## 12. Контроль после запуска

Через 30–60 минут:

```bash
systemctl is-active somonwatch.service
systemctl show somonwatch.service \
  -p MainPID -p MemoryCurrent -p MemoryPeak -p CPUUsageNSec -p TasksCurrent -p NRestarts
journalctl -u somonwatch.service --since -1h --no-pager
systemctl --failed
```

Повторить ключевые host-проверки и сравнить с preflight:

```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
ss -lntup
firewall-cmd --list-all
```

Ожидается:

- Docker containers в прежнем состоянии;
- новых listening ports нет;
- firewall rules не изменились;
- Traefik не менялся;
- память somonwatch значительно ниже заданного `MemoryMax=256M`;
- CPU практически нулевой между poll.

---

## 13. Повседневные команды

```bash
systemctl status somonwatch.service
journalctl -u somonwatch.service -f
journalctl -u somonwatch.service --since today
systemctl restart somonwatch.service
systemctl stop somonwatch.service
systemctl start somonwatch.service
```

Рестарт не переотправляет известные ID: SQLite сохраняется.

---

## 14. Backup

Онлайн-backup не требует остановки сервиса:

```bash
sudo /root/somon-rent-watcher-build/somon-rent-watcher/scripts/backup-installed.sh
```

Либо указать каталог:

```bash
sudo ./scripts/backup-installed.sh /var/backups/somonwatch/before-upgrade
```

Backup содержит:

- SQLite snapshot через `.backup`;
- env-файл;
- бинарник;
- systemd unit;
- документацию;
- SHA256SUMS.

Каталог root-only, поскольку содержит Telegram token.

---

## 15. Безопасный upgrade

```bash
# 1. Снять backup без остановки
sudo ./scripts/backup-installed.sh /var/backups/somonwatch/before-vNEXT

# 2. Собрать новый архив в отдельном build-каталоге
./scripts/build.sh

# 3. Остановить только новый сервис
systemctl stop somonwatch.service

# 4. Установить новый бинарник/unit; env сохраняется
sudo ./scripts/install-almalinux.sh

# 5. Doctor
set -a
. /etc/somonwatch/somonwatch.env
set +a
runuser -u somonwatch --preserve-environment -- /opt/somonwatch/somonwatch doctor
unset TELEGRAM_BOT_TOKEN TELEGRAM_ADMIN_USER_IDS TELEGRAM_ADMIN_USER_ID TELEGRAM_TARGET_CHAT_ID

# 6. Запустить только somonwatch
systemctl start somonwatch.service
journalctl -u somonwatch.service -n 100 --no-pager
```

Install script откажется заменять файлы, если `somonwatch.service` всё ещё активен. Это защищает от частично заменённого running install.

---

## 16. Rollback

Пример, где `/var/backups/somonwatch/before-vNEXT` — проверенный backup:

```bash
systemctl stop somonwatch.service
cp -a /var/backups/somonwatch/before-vNEXT/somonwatch /opt/somonwatch/somonwatch
cp -a /var/backups/somonwatch/before-vNEXT/somonwatch.service /etc/systemd/system/somonwatch.service
cp -a /var/backups/somonwatch/before-vNEXT/somonwatch.env /etc/somonwatch/somonwatch.env
chmod 0755 /opt/somonwatch/somonwatch
chmod 0600 /etc/somonwatch/somonwatch.env
chown root:root /opt/somonwatch/somonwatch /etc/somonwatch/somonwatch.env
```

Восстанавливать DB только если новая версия реально изменила/повредила данные:

```bash
cp -a /var/lib/somonwatch/somonwatch.db /var/lib/somonwatch/somonwatch.db.failed
cp -a /var/backups/somonwatch/before-vNEXT/somonwatch.db /var/lib/somonwatch/somonwatch.db
chown somonwatch:somonwatch /var/lib/somonwatch/somonwatch.db
chmod 0600 /var/lib/somonwatch/somonwatch.db
```

Затем:

```bash
restorecon -RF /opt/somonwatch /etc/somonwatch /var/lib/somonwatch 2>/dev/null || true
systemctl daemon-reload
systemctl start somonwatch.service
```

---

## 17. Диагностика

### 17.1. Telegram `getUpdates` conflict

Симптом: Telegram сообщает, что другой `getUpdates` уже выполняется.

Причина: одновременно запущены `somonwatch run`, `somonwatch ids` или другой клиент того же token.

```bash
ps aux | grep '[s]omonwatch'
systemctl status somonwatch.service
```

Оставить только systemd-service.

### 17.2. Бот не пишет в группу

Проверить:

- корректность отрицательного `TELEGRAM_TARGET_CHAT_ID`;
- бот состоит в группе;
- у него есть право отправлять сообщения и фотографии;
- `doctor` проходит `getChat`;
- в журнале нет Telegram API error.

### 17.3. HTTP 403/429 от Somon

Сервис сам переходит в backoff минимум на 2 часа и уведомляет каждого настроенного администратора лично.

Не делать:

- частые ручные restarts;
- уменьшение poll interval;
- proxy rotation;
- CAPTCHA bypass;
- смену fingerprint для обхода.

Проверить позже:

```bash
journalctl -u somonwatch.service --since -6h --no-pager
```

### 17.4. Изменился HTML Somon

При sanity-check error база seen и snapshot не изменяются. HTML сохраняется в:

```text
/var/lib/somonwatch/debug/
```

Просмотр:

```bash
ls -lt /var/lib/somonwatch/debug | head
sudo -u somonwatch wc -c /var/lib/somonwatch/debug/*.html
```

Файлы могут содержать публичный текст объявлений; не публиковать их без необходимости.

### 17.5. SQLite

Проверка:

```bash
sqlite3 /var/lib/somonwatch/somonwatch.db 'PRAGMA quick_check;'
sqlite3 /var/lib/somonwatch/somonwatch.db 'SELECT COUNT(*) FROM seen_ads;'
```

Не удалять DB для «лечения»: это сбросит фильтр и историю seen. При полной потере DB сервис создаст новый baseline без рассылки текущих карточек, но пользовательские настройки придётся задать заново.

### 17.6. SELinux

Не отключать SELinux и не выполнять `setenforce 0`.

Стандартные пути и unit обычно работают без custom policy. При AVC:

```bash
ausearch -m AVC -ts recent | tail -100
ls -lZ /opt/somonwatch /etc/somonwatch /var/lib/somonwatch
restorecon -RF /opt/somonwatch /etc/somonwatch /var/lib/somonwatch
```

Не включать произвольные SELinux booleans без понимания причины.

---

## 18. Мягкое отключение/удаление

Остановить и отключить только этот unit:

```bash
systemctl disable --now somonwatch.service
```

Сначала сохранить backup. Затем при необходимости удалить executable/unit:

```bash
rm -f /etc/systemd/system/somonwatch.service
rm -rf /opt/somonwatch
systemctl daemon-reload
```

По умолчанию сохранить:

```text
/etc/somonwatch/
/var/lib/somonwatch/
/var/backups/somonwatch/
```

Они содержат конфигурацию, token, seen IDs и backup.

---

## 19. Compliance

На дату подготовки опубликованные условия Somon содержали ограничения на доступ сторонними программами и автоматизированный сбор данных. Технический runbook не заменяет разрешение владельца сайта.

Перед постоянной эксплуатацией проверить актуальные версии:

```text
https://somon.tj/about/license/
https://somon.tj/about/rules/
https://somon.tj/robots.txt
```

Приложение намеренно не использует внутренний `/api`, `/author/`, пагинацию, browser automation, прокси-ротацию или обход блокировок.
