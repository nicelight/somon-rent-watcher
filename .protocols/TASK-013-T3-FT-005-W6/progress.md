# Progress — TASK-013-T3-FT-005-W6

## Current status
- state: implementing
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
