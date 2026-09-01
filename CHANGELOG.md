# Changelog

## Unreleased

- Allow multiple configured Telegram administrators to share one filter.
- Allow administrator commands, callbacks and scoped text input in the configured target group while rejecting other users and groups.
- Isolate pending text input by administrator and chat, and fan out technical notifications to every configured administrator.
- Avoid classifying Tailwind positioning classes such as `top-4` as promoted Somon cards.

## 1.0.0 — 2026-09-01

- Initial production-oriented MVP.
- Fresh-only Somon category monitoring with random 10–30 minute polling.
- Price, room, floor, seller-ad-count, negative-word and ordinary/VIP filters.
- Telegram private control menu and group notifications.
- SQLite state, baseline initialization, gap recovery, rate limiting and 403/429 backoff.
- AlmaLinux 9 systemd deployment with a dedicated user and no inbound ports.
- Safe initial paused state and per-poll detail request cap.
