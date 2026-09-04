---
description: Лог изменений Memory Bank.
status: active
---
# Changelog

## [2026-09-01] Initial setup
- Created Memory Bank skeleton
- Seeded core docs (product, requirements, testing, task registry)

## [2026-09-01] Brownfield current-state mapping

- Mapped implemented product, C4 architecture, integrations, runtime/state lifecycle, operations and verification surfaces from repository evidence.
- Added explicit as-is/target separation; accepted Boundary Map and Global Backbone remain undecided.
- Registered missing PRD and prevented roadmap/task generation under the PRD-less rule.
- Recorded unresolved verification: no local Go toolchain, no live external/host probe, no Git metadata and draft-TZ SQLite-driver drift.

## [2026-09-01] Public repository preparation

- Reframed the public README around rental discovery on popular Central Asian housing platforms without naming the current source platform.
- Kept the public claim accurate by documenting one current platform-specific adapter rather than claiming implemented multi-platform support.
- Replaced the concrete production host address and neighbouring-workload identifiers in the public runbook/preflight script with generic operator-safe checks.
- Recomputed delivery-manifest checksums for the changed README, runbook and preflight script.
- Initialized the local Git repository on branch `main` with GitHub-compatible local author identity; pre-initialization history is not recoverable.
- Published branch `main` to the public origin `https://github.com/nicelight/somon-rent-watcher`.

## [2026-09-01] Wave 1 / Local Docker, multi-admin and scheduler controls

- Added: local hardened Docker Compose route with one permanent container and persisted SQLite volume.
- Updated: Telegram control to an administrator allowlist shared across private chats and the configured target group.
- Added: persisted configurable random polling range plus single-flight `Сканировать сейчас` with busy and zero-result feedback.
- Fixed: live Somon promoted-card classification for Tailwind positional classes.
- Verified: Docker formatting/tests/vet/CGO/linkage gate and local SQLite/Telegram/Somon doctor passed.
- Synchronized: draft technical specification, product/architecture/integration/lifecycle/testing/glossary/runbook routes and verification state.

## [2026-09-01] Adaptive continuity threshold

- Fixed: a persisted polling range above the static 45-minute gap threshold no longer makes a normal scheduled delay look like downtime.
- Preserved: a complete loss of overlap between consecutive ordinary-ID snapshots still independently triggers one recovery sweep.
- Verified by an app-level regression for the 10–70 minute polling range.

## [2026-09-01] Telegram description readability

- Changed: each emoji-led field in an advertisement description now starts on its own Telegram caption line.
- Preserved: description normalization, escaping, truncation and compound emoji sequences.
- Verified by a caption-rendering regression based on an observed live advertisement.

## [2026-09-02] Silent snapshot-gap recovery

- Refactored: continuity detection now keeps snapshot discontinuity and overdue-poll age as typed causes instead of coupling behavior to a rendered reason string.
- Changed: ordinary-ID turnover continues to trigger recovery and structured warning logs but no longer sends repetitive private Telegram messages.
- Preserved: an adaptive overdue successful-poll condition still triggers recovery and an hourly rate-limited private warning.
- Verified: an app integration test proves silent snapshot recovery and overdue-only administrator notification.

## [2026-09-02] AlmaLinux production deployment

- Deployed GitHub commit `7f9c5f50d659` from a retained clean checkout as a native, unprivileged systemd service.
- Verified the production build/test/vet/CGO gate and live SQLite, Telegram and Somon doctor before service startup.
- Created a fresh paused 60-card baseline with zero service restarts and no current-ad delivery.
- Confirmed existing containers, listening sockets, firewall configuration and SELinux behavior were unchanged.
- Kept the local Compose container stopped to prevent competing Telegram long polling with the same bot token.
