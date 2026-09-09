# Verification — TASK-005-T3-FT-004-W1

## Basis and scope
Operator-authorized maintenance correction of reproduced price corruption and isolated
shared-host deployment. Somon adapter contract, current integration contract, existing
runtime lifecycle and AlmaLinux runbook apply. No unfinished feature changes or seen replay.

## Functional evidence
- Regression RED: regression-red.log, obtained before parser mutation.
- Same regression GREEN and full native gate: local-build.log; target independently
  ran the same native tests/vet/build before stopping service (production-release.log).
- Offline prices: offline-prices-green.log, 60 category and 38 saved detail prices
  match independently inspected visible prices.
- Fresh verifier-owned postflight: final-runtime-proof.json compares actual process
  executable, installed binary and target-built SHA; exact clean source commit;
  backup-to-current SQLite row subset and settings/environment equality; valid old
  rollback version; SQLite integrity; no watcher restarts.
- production-poll.log proves a completed post-start cycle without selected errors.
- host-comparison.json proves container/service/listener/route/SELinux/failed-unit
  stability. Initial firewall digest is not reusable because it included traffic
  counters. Two corrected postflight rule hashes match; command review confirms no
  firewall rule writes. Pre/post firewall rule equality by hash remains unproven.

## Independent co-review focuses
- Parser correctness/compatibility: Codex Luna xhigh, no candidate findings.
- Concrete deployment/state-preservation/rollback: Codex Luna xhigh; error handling,
  identity and retry concerns addressed before execution. Adjudication in progress.md.

## Result
All price-correction and safe installation outcomes are verified for this actual run.
The firewall measurement limitation is explicit and is not treated as a passing
pre/post hash comparison. No evidence of unrelated service or configuration mutation.

VERDICT: PASS

Next: independent per-task red-verify, then explicit-owner task closure.
