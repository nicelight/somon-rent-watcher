---
description: AlmaLinux 9 operations route with the accepted isolated production-release contract.
status: active
baseline_kind: as-is
last_verified: 2026-10-07
last_updated: 2026-10-07
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
  - docs/RUNBOOK_ALMALINUX_9.md
  - deploy/somonwatch.service
  - scripts/install-almalinux.sh
  - scripts/backup-installed.sh
---

# AlmaLinux 9 operations

## Source of procedural detail

[docs/RUNBOOK_ALMALINUX_9.md](../../docs/RUNBOOK_ALMALINUX_9.md) is the existing
step-by-step operator document. Its current procedures remain as-is evidence;
the accepted FT-003 delta below is the target contract that an implementation
task must reflect in executable procedures before deployment.

## Accepted FT-003 release procedure

FT-003 reuses the existing systemd upgrade route and adds only the accepted
publication, git-synchronization, runtime-identity, and isolation boundaries.
The executable task must preserve this order:

1. On the intended local release tree, run the repository-native formatting,
   tests, vet, CGO build, linkage, and version gate.
2. Without changing the gated source tree, create or confirm the release
   commit, require a clean worktree, record the full commit ID, push it to the
   configured GitHub remote, and verify that the intended remote ref resolves
   to that exact commit.
3. Immediately before the first production write, take a fresh read-only host
   snapshot and re-establish one healthy watcher identity from current evidence:
   production checkout, systemd unit, active process/executable, environment
   and SQLite paths, and absence of a competing watcher process/container.
   Also confirm that unrelated workloads are healthy. Missing, conflicting, or
   ambiguous identity stops the release without a production write.
4. Synchronize only the identified production checkout to the already-pushed
   commit. The update must be fast-forward/exact-commit and must not discard,
   overwrite, stash, or reset production-local work; a dirty or diverged
   checkout stops the release, and a post-sync HEAD mismatch also stops.
5. Build in that checkout with `scripts/build.sh`, verify the produced binary
   reports the intended commit and has valid target linkage, then run that
   produced binary through the existing redacted `doctor` environment and
   service-user route. A failure leaves the installed runtime unchanged and
   stops the release.
6. Stop only the identified `somonwatch.service`, record the now-stable
   read-only SQLite/settings evidence, use the existing scoped install route,
   and start only that unit. Preserve `/etc/somonwatch/somonwatch.env`,
   `/var/lib/somonwatch/somonwatch.db`, and the debug directory; no database
   deletion, replacement, migration, or reset is permitted.
7. Verify the installed version/commit, unit health and restart count, redacted
   doctor result, and scoped recent logs. Repeat the stable host/state probes
   and compare them with preflight, allowing changes only to the watcher
   checkout, build/install artifacts, process, and unit lifecycle.

The accepted production runtime remains the documented host
`somonwatch.service`. If fresh evidence identifies another runtime shape, the
release stops for an explicit operator decision; FT-003 does not introduce a
generic multi-runtime deployment mechanism or a runtime migration.

## FT-003 release proof

- Known initial state: local gate and GitHub ref identify one commit; production
  is healthy; exactly one documented watcher runtime and its checkout/state
  paths resolve; the checkout is safe for exact-commit synchronization; SQLite
  passes a read-only integrity probe. Historical runtime evidence cannot
  substitute for this fresh preflight.
- State preservation: compare a non-secret settings digest and `seen_ads` count
  from the stopped pre-install state with the post-start state, require settings
  equality and no loss of previously counted seen rows, and keep the database
  path in place. Normal watcher activity after restart may increase the seen
  count.
- Host isolation: compare stable identities/statuses for unrelated containers
  and services, failed units, listening sockets, firewall rules, routes, and
  SELinux mode. Do not publish host inventory or secrets.
- Safe rerun: any retry starts again from the fresh read-only identity/health
  preflight and exact commit checks. Target build/doctor failure requires no
  runtime mutation or unrelated cleanup; later failure handling stays confined
  to the existing watcher upgrade/troubleshooting route.
- Evidence: keep ordered, redacted receipts in the task-selected operational
  evidence path; include local/remote/production commit IDs, preflight identity
  and health, target build/version/linkage/doctor, scoped unit/log status,
  state comparison, and unrelated-host comparison.

## Current deployment shape

- Target described by the source runbook: AlmaLinux 9 host, dedicated `somonwatch` user, one host `systemd` service and no inbound port.
- Installed executable: `/opt/somonwatch/somonwatch`.
- Root-readable environment: `/etc/somonwatch/somonwatch.env` (`0600`).
- Writable data: `/var/lib/somonwatch/somonwatch.db` and `/var/lib/somonwatch/debug/`.
- Unit: `/etc/systemd/system/somonwatch.service`.
- Backups: `/var/backups/somonwatch/`; backup output contains the Telegram token and remains root-only.

## Safe current sequence

1. Prepare the Telegram bot and determine all administrator/target chat IDs with `somonwatch ids`; do not run it concurrently with the service. Disable BotFather Group Privacy when group-based text input is required.
2. Verify archive checksums and take a read-only host preflight snapshot with `scripts/preflight-almalinux.sh`.
3. Install only Go/GCC/SQLite build prerequisites and build on the target host through `scripts/build.sh`.
4. Install without starting through `scripts/install-almalinux.sh`; populate and protect the env file.
5. Run `somonwatch doctor` as the service user before first start. It verifies SQLite/Telegram/live Somon parsing without creating the seen baseline.
6. Enable/start only `somonwatch.service`; confirm the first cycle creates baseline with no group notifications and remains paused.
7. Configure the shared filter privately or in the target group, then explicitly enable monitoring.
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

## Production verification

- On 2026-09-09 final parser hotfix `93cbe9ebcfe6` was built and tested locally
  and on the target, passed staged service-user doctor, and replaced only the
  existing watcher via the backup/install route. Running, installed and target-built
  binaries matched; SQLite/settings/env and all pre-final-release seen rows survived.
  The first cycle successfully delivered one new matching ad. Other service/container/
  listener/route/firewall-rule/SELinux/failed-unit fingerprints matched across the
  final release. Firewall fingerprints omit generated timestamps and traffic counters.
  The first attempt's invalid counter-sensitive fingerprint remains explicitly limited
  in historical receipts and is not substituted for the final normalized comparison.
  [Hotfix verification](../../.protocols/TASK-005-T3-FT-004-W1/verification.md) records
  the exact receipts and rollback backup.

- On 2026-09-02 the target AlmaLinux 9 host passed the read-only preflight, native build/test/vet/linkage gate and live `somonwatch doctor` for commit `7f9c5f50d659`.
- The installed systemd service created a fresh 60-card baseline while paused, remained active with zero restarts and passed unit verification.
- Existing containers, listening sockets, firewall configuration, failed-unit state and SELinux behavior were unchanged after deployment.
- Administrators whom the bot has never contacted must open the bot once before Telegram permits private operational notifications; target-group control is independent of that limitation.

## Accepted FT-005 release procedure

Operator-authorized 2026-10-07 release publishes ALL current source, including verified keyword functionality/repairs and current rental fallback. Reuse Accepted FT-003 release procedure order and isolation, exact commit, backup and native gates. FT-005 extends only its no-migration restriction: existing Store initializer may add search_monitors/search_ad_state tables/indexes at normal startup, preserving every existing table/row/settings/offset and database file identity. No manual migration, deletion, reset or baseline replay. This exception applies to this authorized FT-005 release and does not rewrite historical FT-003 acceptance.

Before runtime change, staged doctor MUST use a disposable SQLite .backup clone owned by somonwatch and a scratch DEBUG_DIR under watcher data directory; override DB_PATH after sourcing existing env. New Store.Open MUST NOT initialize live SQLite while the old binary is running. Doctor may perform existing read-only Telegram identity/chat and Somon fetch checks, MUST NOT run polling or delivery. Remove scratch/staged artifacts afterward. Keep root-only rollback backup on host. Stop/install/start only somonwatch.service; compare stopped pre-install old rows in all existing tables plus env/settings/DB identity to post-start rows (allow ordinary new seen rows and offset advancement, no reset), healthy single process and unrelated host fingerprints. No automatically enabled keyword monitor. On later failure restore only backed-up binary/unit lifecycle; never restore/delete database or alter unrelated workloads. Evidence belongs to TASK013 operational receipts.

## FT-005 installed release

2026-10-07 exact8b48c4cef11237716e1dbc471cf363018e48c613 installed/running/build binary checksum66a8cdcb83419c0868c4015170495c47f1da3da41672ce578a0b3910ada9d0d2. Local+target native and staged clone/live doctor PASS, single active service zero restarts, existing state/settings/env/DB identity and every prior seen row preserved; all unrelated normalized fingerprints equal. Old rental remains paused, no keyword monitor auto-created/enabled. Backup remains host/root-only /var/backups/somonwatch/keyword-release-20261007T163309Z, stage/scratch removed. [Independent verification](../../.protocols/TASK-013-T3-FT-005-W6/verification.md) and [semantic review](../../.protocols/TASK-013-T3-FT-005-W6/red-verification.md).

Historical deployed93c was a separate branch from main; first ff-only refused before runtime change. Current main already contained its price fixes; local non-rewriting merge retained exact current tree and joined ancestry, then fresh native/ref/preflight yielded safe target FF. For future releases check merge-base exit explicitly before publication; never use reset/stash on production to conceal divergence.
