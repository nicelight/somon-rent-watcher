---
description: Compact Memory Bank routing for current AlmaLinux 9 deployment, operation, backup and recovery procedures.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
last_updated: 2026-09-01
source_of_truth:
  - docs/RUNBOOK_ALMALINUX_9.md
  - deploy/somonwatch.service
  - scripts/install-almalinux.sh
  - scripts/backup-installed.sh
---

# AlmaLinux 9 operations — current route

## Source of procedural detail

[docs/RUNBOOK_ALMALINUX_9.md](../../docs/RUNBOOK_ALMALINUX_9.md) is the existing step-by-step operator document. This Memory Bank page records where the current operational contracts live; it does not duplicate every shell command.

## Current deployment shape

- Target described by the source runbook: AlmaLinux 9 host, dedicated `somonwatch` user, one host `systemd` service and no inbound port.
- Installed executable: `/opt/somonwatch/somonwatch`.
- Root-readable environment: `/etc/somonwatch/somonwatch.env` (`0600`).
- Writable data: `/var/lib/somonwatch/somonwatch.db` and `/var/lib/somonwatch/debug/`.
- Unit: `/etc/systemd/system/somonwatch.service`.
- Backups: `/var/backups/somonwatch/`; backup output contains the Telegram token and remains root-only.

## Safe current sequence

1. Prepare the Telegram bot and determine administrator/target chat IDs with `somonwatch ids`; do not run it concurrently with the service.
2. Verify archive checksums and take a read-only host preflight snapshot with `scripts/preflight-almalinux.sh`.
3. Install only Go/GCC/SQLite build prerequisites and build on the target host through `scripts/build.sh`.
4. Install without starting through `scripts/install-almalinux.sh`; populate and protect the env file.
5. Run `somonwatch doctor` as the service user before first start. It verifies SQLite/Telegram/live Somon parsing without creating the seen baseline.
6. Enable/start only `somonwatch.service`; confirm the first cycle creates baseline with no group notifications and remains paused.
7. Configure the filter privately, then explicitly enable monitoring.
8. Compare service/container/socket/firewall state with preflight after deployment.

## Current maintenance routes

| Need | Existing procedure/evidence |
|---|---|
| Build/package | [scripts/build.sh](../../scripts/build.sh) and [.memory-bank/guides/local-development.md](../guides/local-development.md). |
| Host survey | [scripts/preflight-almalinux.sh](../../scripts/preflight-almalinux.sh). |
| Install/upgrade | [scripts/install-almalinux.sh](../../scripts/install-almalinux.sh); refuses replacement while the service is active and preserves existing env. |
| Online backup | [scripts/backup-installed.sh](../../scripts/backup-installed.sh); uses SQLite `.backup`. |
| Rollback/troubleshooting/removal | Sections 16–18 of [docs/RUNBOOK_ALMALINUX_9.md](../../docs/RUNBOOK_ALMALINUX_9.md). |

## Safety boundaries

- Do not expose/log the env file or Telegram token.
- Do not delete the DB as routine troubleshooting: it removes settings and seen history.
- Do not change Docker, Traefik, firewall, routes or SELinux for this service; they are outside the current install contract.
- Do not react to 403/429 with rapid restarts, lower poll intervals or blocking circumvention.
- Re-check current Somon legal/robots terms before sustained production operation.

## Needs verification

The mapped workspace contains only repository evidence. The current state of the documented server, Telegram bot/group and live Somon HTML has not been observed; run preflight and `somonwatch doctor` on the intended host.
