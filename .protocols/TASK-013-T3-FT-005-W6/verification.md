---
description: Independent functional production-release verification for TASK-013.
status: active
---
# Verification — TASK-013-T3-FT-005-W6

## What was verified
Independent ROLE Reviewer /verify of FT-005-AC-010 (REQ-008/REQ-014 release subset), not dependency product claims. Exact indexed T3/W6/FT005 card, dependencies, gates/verify array and full context/plan/progress/handoff/protocol were checked. Task remains in_progress.

## Verification basis
Direct canonical inputs: runbooks/almalinux-9-operations.md#accepted-ft-005-release-procedure, #accepted-ft-003-release-procedure, #ft-003-release-proof; states/runtime-lifecycle.md#keyword-persistence-and-mutation-rules; contracts/boundary-map.md#persistence-contract (Composition/Polling Application → Persistence Adapter accepted rows). Constitution, spec registry/backbone Revision1 and applicable tier-policy claim ownership, execution evidence, T3 obligations and closure authority applied. Purpose: all authorized current source in one existing watcher; no product implementation, DB reset/manual migration, env edits, unrelated host operations or delivery probes.

## Executor claim path
- Exact-version claim: original healthy installed93cbe9 RED in preflight.json and fresh attempt2-preflight.json → exact8b48c4cef11237716e1dbc471cf363018e48c613 GREEN in attempt2-release.log/release-receipt.json. Honest attempt1 ff-only rejection occurred before runtime modification; attempt2 merge retained first-parent tree (a061c0c43a10bdfeb25ad54996081de0b21a4f37) and historical93c ancestor, enabling exact FF without source changes/reset.
- State-loss/isolation RED_NOT_APPLICABLE accepted: deliberately inducing production loss/host harm is unsafe. Alternative backup/stopped/post-start state comparison and unrelated fingerprints are recorded in release-procedure.py, attempt2-release.log and preflight/postflight JSON. This does not stand alone for PASS: fresh verifier proof below covers current harm-driving outcomes.

## Task-scoped checklist and new targeted probes
Verifier-owned read-only `ssh igorprod python3 - < .tasks/TASK-013-T3-FT-005-W6/verifier-live-probe.py > .tasks/TASK-013-T3-FT-005-W6/verifier-live-probe.json` exited0 at 2026-10-07T16:37:14Z. Only SELECT/read-only SQLite URI; no doctor, restart, install, env edit, activation, message or backup copying.

- [x] Exact publication/checkout/installed/running release: fresh local HEAD and remote ls-remote main both exact8b48; target HEAD/origin main clean exact8b48. `/proc/1063284/exe`, installed and target-built binary SHA256 all66a8cdcb83419c0868c4015170495c47f1da3da41672ce578a0b3910ada9d0d2; version8b48, target linkage valid. Process bound to somonwatch.service and actual DB file; one watcher process/no competing watcher container; active/running User somonwatch NRestarts0.
- [x] All old data preserved: fresh transactional read snapshots live and root-only backup `/var/backups/somonwatch/keyword-release-20261007T163309Z/somonwatch.db` have integrityOK. Every prior row in every old non-state table preserved (seen_ads14675→14681, settings1→1); settings exactly equal. All prior state keys and stable values retained; only accepted ordinary poll snapshot/time may advance, Telegram offset monotonic. No automatic monitor/history: search_monitors0/search_ad_state0, both additive tables present.
- [x] Env/unit/DB identity: live env equals backup and original preflight hash; unit equals backup and accepted deploy unit. Live DB inode263941/device64770 and birth time2026-09-02 independently show original file predates release, grounding executor's stopped inode-equality check; all old data retained. Backup directory root-owned0700; DB/env/checksum files root-owned0600. Backup remains on target.
- [x] Host isolation: fresh normalized container identities/status/start/restarts (8 containers), unrelated running services, listeners, routes, firewall rules without counters/timestamps, SELinux Enforcing and failed units all match attempt2-preflight.json. Comparison and stage markers in verifier-comparison.json.
- [x] Safe staging/cleanup and runtime health: inspected release script routes pre-install Store.Open doctor exclusively to disposable .backup clone, overrides DB_PATH after env sourcing, service user and scratch directory, removes both scratch and staged binary before install. Ordered release log shows staged clone=True/live clone=False doctor success and stopped-state exact old table equality. Fresh probes independently confirm stage/scratch absent, no scoped ERROR logs or new baseline, healthy expected process. Inspected doctor implementation does identity/chat/source GET reads and no polling/delivery. Current SQLite mutation remains normal Store initializer only; no architecture/source code changed by this release.

## Repeated checks
Fresh local `docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest ./scripts/build.sh` exited0; verifier-local-native-gate.log records formatting, all unit/integration tests, vet, CGO build, SQLite linkage and exact8b48 version. Chosen as cheap fresh gate instead of reusing earlier local result. Current cmd/internal/scripts/deploy/go.mod/go.sum/VERSION diff against8b48 is empty; only workflow documents/evidence changed.

## Reused execute evidence
Attempt2 target native gate (`nice -n10 ./scripts/build.sh`, cwd `/root/somon-rent-watcher`, COMMIT8b48/GOMAXPROCS2/GOFLAGS-p=2) exited0 in ordered attempt2-release.log, completed16:32Z before install; source exact8b48 and clean checkout, built checksum66a8, release procedure's explicit bounded gate command and current redacted target toolchain/binary build information independently grounded in verifier-target-toolchain.log. Source/config/modules/accepted unit and binary are unchanged. Target build was not repeated because verification production scope is read-only. Current build/process/linkage plus fresh equivalent local native gate independently ground outcome; executor target result remains supporting evidence.
Staged/live doctor and stopped comparison receipts support historical sequencing; fresh live state/identity/cleanup/host observations above independently cover all harm-driving claims. No volatile remote doctor result is claimed current solely from receipt.

## Adjudication and scope
Installed finding-adjudication pack applied. Required fresh Codex Luna xhigh code-quality co-review launch attempted once but could not launch (agent thread limit); verify skill says continue without retry/block. Main Reviewer independently inspected task scripts and actual no-code release change surface; no evidenced task-relevant violation. No hypothetical improvements or historical dependency re-verification. Additional verifier task artifacts are necessary same-outcome advisory scope additions; no forbidden production writes or lifecycle changes.

## Verdict
VERDICT: PASS

## Handoff
Separate fresh `/red-verify TASK-013-T3-FT-005-W6` required before T3 closure. Explicit root owner handles lifecycle and wave Memory Bank sync after required semantic PASS. Verifier did not change status or invoke semantic review/sync. No blocking finding or unresolved branch.
