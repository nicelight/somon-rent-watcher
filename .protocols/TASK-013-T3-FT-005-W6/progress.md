# Progress — TASK-013-T3-FT-005-W6

## Current status
- state: done
- last update: 2026-10-07

## What was done
Fresh task-plan APPROVE Revision1; mb-lint58/strict doctor PASS0errors0warnings. Root preflight confirmed all done dependencies, direct canonical context/scope, protocol initialized from required templates; planned→ready→in_progress attempt1 before prospective probe/mutation. No production writes yet. Prepared adapted watcher-only release and read-only host probe, syntax compiled.

## Claim-linked RED / GREEN
- attempt: 1
- accepted claim: FT-005-AC-010
- applicability: applicable exact release; preservation/isolation RED_NOT_APPLICABLE: inducing production data loss or host harm is unsafe; alternative proof stable backup/stopped/post-start data and unrelated fingerprints.
- RED: fresh read-only old installed version, to be bound to exact release commit after local gate/publication.
- GREEN: pending exact deployed commit and state/host proof.
- T3: staged doctor disposable clone; backup/root-only and scratch cleanup; no source/DB/manual migration or unrelated runtime edits.

## Evidence
.tasks/TASK-013-T3-FT-005-W6/; fresh review report.

## Next step
Local native gate, release commit/push, fresh production RED/preflight and scoped release.

## Local gate receipt
- attempt: 1; receipt_status: current
- claim: FT-005-AC-010 release source gate
- command: docker run --rm --network none -v "$PWD:/src" -w /src somon-price-hotfix-builder:latest ./scripts/build.sh
- cwd: repository root
- exit_code: 0
- input_state_basis: current entire accepted code, verified repairs, no source changes since independent TASK011/012 verification; c04c287 HEAD plus described repairs.
- completed_at: 2026-10-07T16:26Z
- evidence: .tasks/TASK-013-T3-FT-005-W6/local-native-gate.log; format/all tests/vet/CGO/SQLite linkage PASS. Precommit binary is supporting-only version evidence; exact new release commit will be embedded on target.

## Publication / pre-write RED
Clean gated source committed b7a8c53c1fdf7ca926d35bac0c9a320e1dc8578d; pushed origin main and ls-remote exact ref confirmed. No Go/source changed after gate. Fresh read-only preflight.json captured immediately before first production write: old93cbe9 installed/checkout, healthy one service, integrityOK, unrelated baseline. FT-005-AC-010 RED: installed version differs from expected b7a8c53. This is claim absence, no artificial failure. Next exact FF/build/clone doctor/backup/install/start and comparison.

SSH launch transport timed out during banner exchange before remote Python started (exit255). Fresh preflight-retry.json confirmed exact unchanged old checkout/process/version/state/host. Same attempt safe retry, no uncertain side effect. Initial transport receipt supporting-only; fresh preflight is current.

## Attempt 2 — Git ancestry correction
Attempt1 remote fetch succeeded, ff-only rejected before build/stage/backup/service/data changes because installed93c is separate published historical hotfix ancestry. Fresh remote inspection confirms non-shallow clean checkout still93c and old runtime. Initial ancestry assumption was false; retained evidence honest, attempt1 receipt supporting-only, original AC010 RED retained. Native current parser includes all historical hotfix changes (diff old93c→current parser is only new validated detail-body checks); old price tests and local native gates passed. Root will create a non-rewriting merge commit retaining EXACT current tree with historical93c as second parent, making target FF possible without resetting production or changing source. This is Git lineage correction within exact-release tactic, no task scope/spec/product/architecture change. New attempt2 exact gate/ref/preflight required before retry.

## Attempt2 GREEN and handoff
FT-005-AC-010 exact source unchanged merge8b48c4cef11237716e1dbc471cf363018e48c613 (HEAD tree equals first-parent tree, old93c ancestor) passed fresh local native gate with embedded8b48 version. Clean release pushed/ref exact, attempt2-preflight old93c healthy14675seen. Remote FF exact, target native format/alltests/vet/CGO/linkage/version PASS. Staged service-user doctor clone PASS; removed stage/scratch; backup root-only /var/backups/somonwatch/keyword-release-20261007T163309Z. Stopped install exact preexisting tables/settings/env/DB identity unchanged. Started only watcher, live doctor PASS; one healthy service NRestarts0, installed/running/build checksum66a8cdcb83419c0868c4015170495c47f1da3da41672ce578a0b3910ada9d0d2. All old rows/state keys/Telegram offset preserved,14675→14681 seen from ordinary runtime; added search tables empty/no automatic monitor. GREEN same AC010 exact version probe; preservation/isolation RED_NOT_APPLICABLE alternative proof met. attempt2-postflight and release-receipt JSON assert all unrelated normalized fingerprints equal. No source code or forbidden host/data changes. Additional task-owned probes/receipts are advisory scope deviations supporting same release.

## Execution handoff status
- state: verifying
- next step: independent /verify TASK013, then T3 semantic review and feature refresh/owner closure/sync.

## Current reuse candidates
- attempt: 2; receipt_status: current; claim FT-005-AC-010 local exact native gate; command: Docker network-none ./scripts/build.sh; exit0; completed2026-10-07T16:31Z; evidence attempt2-local-native-gate.log; unchanged release source/merge tree.
- attempt: 2; receipt_status: current; claim FT-005-AC-010 scoped release; command ssh igorprod python3 - 8b48c4cef11237716e1dbc471cf363018e48c613 < release-procedure.py; exit0; evidence attempt2-release.log/release-receipt.json; completed2026-10-07T16:33Z; runtime/target exact release, backup/clone/install/live doctor and state comparison.
- attempt: 2; receipt_status: current; claim FT-005-AC-010 unrelated/state identity; command ssh igorprod python3 - < host-state-probe.py, local compare; exit0; attempt2-preflight/postflight JSON + release-receipt; expected exact8b48 release. Earlier transport/FF receipts supporting-only, original old-version RED retained, no unsafe replay.

## Owner lifecycle decision
Root explicit standalone manual owner writes done after fresh independent functional PASS and per-task semantic-pass. All TASK013-owned AC010 outcomes and tier obligations complete, no production remainder. Source exact8b unchanged, no forbidden scope touched. Feature semantic refresh and wave sync/post-sync gates next; these do not reopen completed task.

## Final claim-linked evidence mapping
- attempt: 2
- applicability: applicable (exact release); state-loss/isolation not applicable
- accepted claim locator(s): .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-010
- accepted not-applicable reason and alternative proof: FT-005-AC-010 RED_NOT_APPLICABLE for deliberately inducing production data loss/unrelated host harm; alternative proof stopped exact table/settings/DB identity/env comparison plus independent backup/live old-row/state/offset and unrelated fingerprint comparison.
- RED command/probe: ssh igorprod python3 - < host-state-probe.py, attempt2-preflight.json; installed93c compared to accepted exact8b48 release.
- RED observation and evidence: FT-005-AC-010 installed old93c absent intended current code before change; attempt2-preflight.json and start of attempt2-release.log. Original attempt1 honest RED retained, failure receipts supporting-only.
- GREEN command/probe: scoped ssh release-procedure exact8b48, postflight read-only host probe and local release-receipt comparisons.
- GREEN observation and evidence: FT-005-AC-010 exact8b48 built/installed/running,66a8 checksum, healthNRestarts0, preserved old data/settings/state/offset/env/DB identity and unrelated fingerprints; attempt2-release.log, attempt2-postflight.json, release-receipt.json; independent verifier-live-probe.json and verifier-comparison.json ground current result.
- claim-equivalent probe changes and rationale: same version/state/host method; intended commit8b48 changed only by non-rewriting ancestry merge retaining exact current source tree. No weakening or inherited claims.
- T3 isolation/cleanup/permission evidence: explicit all-source deploy authorization, clone staged doctor, root-only backup, stage/scratch absent, only watcher FF/install/service touched, no manual DB migration/reset/delivery probes or unrelated workloads changed.

## Final sync
Root applied /mb-sync to already-written owner closure; feature/epic lifecycle verified, RTM/testing/operations/navigation/changelog agree. Read changed links/status/evidence; mb-lint58 PASS, post-sync strict doctor0errors0warnings in .tasks/TASK013/post-sync-doctor.json. Code source unchanged from8b48, all feature tasks done; older queues retained. No deployment remainder.
