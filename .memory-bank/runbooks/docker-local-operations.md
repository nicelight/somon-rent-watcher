---
description: Compact Memory Bank route for the local Kubuntu Docker Compose deployment.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-01
source_of_truth:
  - Dockerfile
  - compose.yaml
  - docs/RUNBOOK_DOCKER_LOCAL.md
---

# Local Docker operations — current route

## Runtime shape

- One `somon-rent-watcher:local` image built through a multi-stage Debian/Go CGO build.
- One permanent `somon-rent-watcher` container, no published ports and outbound HTTPS only.
- SQLite and debug HTML persist in named volume `somon-rent-watcher-data` at `/var/lib/somonwatch`.
- Root filesystem is read-only; the process runs as non-root with all Linux capabilities dropped.
- `.env` supplies Telegram secrets and IDs but is excluded from Git and Docker build context.

## Operator route

[docs/RUNBOOK_DOCKER_LOCAL.md](../../docs/RUNBOOK_DOCKER_LOCAL.md) owns the local build, launch, doctor, logs and lifecycle commands. The image build runs the repository formatting/test/vet/build gate. `docker compose down` preserves state; `docker compose down -v` destroys the named volume and must not be used as routine cleanup.

## Telegram behavior

Configured administrators share one filter and persisted polling range and may manage them privately or in the configured target group. Manual scan uses the same single scheduler, shows temporary busy/completion feedback and preserves pause/backoff semantics. Group text input requires BotFather Group Privacy to be disabled. Other users and groups remain unauthorized.
