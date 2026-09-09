---
description: Shared product vocabulary for Somon Rent Watcher.
status: active
last_verified: 2026-09-01
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/analysis/product-brief.md
  - internal/model/ad.go
  - internal/app/app.go
  - internal/filter/settings.go
---

# Glossary

Определения объединяют подтверждённое текущее поведение и принятые target-термины из Product Brief. Target-термины не утверждают, что изменение уже реализовано.

## Terms

- **Ad ID / ID объявления** — положительный числовой идентификатор Somon из URL `/adv/<ID>...`; ключ дедупликации в `seen_ads`.
- **Baseline** — первая успешно распознанная выдача, чьи ID записываются seen без отправки в целевую группу.
- **Fresh / новое объявление** — ID, которого ещё нет в локальном `seen_ads`; это first-seen semantics, а не доказанная дата публикации.
- **Seen** — ID, который больше не должен штатно входить в fresh pipeline; запись хранится в SQLite.
- **Exact match** — свежее объявление, соответствующее текущему диапазону цены и всем остальным условиям общего фильтра.
- **Price fallback** *(target)* — до трёх свежих объявлений, нарушающих только верхнюю границу цены, которые отправляются лишь при отсутствии exact match в текущем цикле.
- **Fallback price ceiling** *(target)* — включительная верхняя граница `150%` от текущего настроенного `PriceMax`; без `PriceMax` price fallback не применяется.
- **Card** — данные из category listing: ID, URL, title, price, rooms, floor, promotion, image, age and position.
- **Ad** — Card, обогащённая данными detail page: description и seller fields.
- **Ordinary** — карточка без признака продвижения.
- **Promoted / VIP** — единый пользовательский тип для Somon VIP и TOP.
- **Card-prefilter** — дешёвая проверка price/rooms/floor/type до detail request; unknown card fields пропускаются к detail stage.
- **Detail-filter** — окончательная проверка detail data, включая unknown fields, seller-ad limit и negative substrings.
- **Monitoring pause** — рассылка отключена, но polls и продвижение fresh baseline продолжаются.
- **Ordinary snapshot** — сохранённый список ID обычной части последнего успешного poll; promoted cards не участвуют в continuity check.
- **Gap suspicion** — отсутствие пересечения ordinary snapshots и/или слишком большой интервал после последнего успешного poll.
- **Recovery sweep** — однократный опрос room-specific pages при gap suspicion.
- **Backoff** — увеличенная задержка до следующего poll после block/HTTP 403/429; обход блокировки не выполняется.
- **Polling range** — общий persisted inclusive-диапазон минут для случайной задержки следующего normal poll; defaults новой/legacy записи берутся из `SOMON_POLL_MIN/MAX`.
- **Append-only admin UI** *(target)* — callback подтверждается, а актуальное меню или результат отправляется новым Telegram-сообщением без редактирования и удаления старых меню.
- **Manual poll / «Сканировать сейчас»** — single-flight пробуждение того же scheduler с итоговой Telegram feedback; не создаёт параллельный poll и не обходит backoff.
- **Admin allowlist** — положительные Telegram user IDs из `TELEGRAM_ADMIN_USER_IDS`; только они могут менять общий фильтр в личном чате или настроенной target group.
- **Doctor** — read-mostly operational check SQLite/Telegram/live Somon parser; может создать/инициализировать пустой SQLite-файл, но не создаёт baseline и не меняет seen IDs.

## Evidence

- [internal/model/ad.go](../internal/model/ad.go): `Card`, `Ad`, `RuntimeStatus`.
- [internal/app/app.go](../internal/app/app.go): baseline, seen, pause, gap, recovery and backoff semantics.
- [internal/filter/settings.go](../internal/filter/settings.go): filter vocabulary and unknown choices.
- [.memory-bank/analysis/product-brief.md](analysis/product-brief.md): accepted target vocabulary for price fallback and append-only admin UI.
- [README.md](../README.md): user-facing terminology.
