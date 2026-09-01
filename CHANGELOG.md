# Changelog

## 1.0.0 — 2026-09-01

- Initial production-oriented MVP.
- Fresh-only Somon category monitoring with random 10–30 minute polling.
- Price, room, floor, seller-ad-count, negative-word and ordinary/VIP filters.
- Telegram private control menu and group notifications.
- SQLite state, baseline initialization, gap recovery, rate limiting and 403/429 backoff.
- AlmaLinux 9 systemd deployment with a dedicated user and no inbound ports.
- Safe initial paused state and per-poll detail request cap.
