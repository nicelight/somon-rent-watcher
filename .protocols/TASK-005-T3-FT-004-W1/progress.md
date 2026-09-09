# Hotfix execution — TASK-005-T3-FT-004-W1

## Execution Attempt
- attempt: 1

## Scope and authority
Operator explicitly requested parser corrections and careful shared-host deployment.
Standalone maintenance branch hotfix/somon-price-extraction from installed
7f9c5f50d659; no pending product feature changes, no history replay/reset.
Somon Adapter remains sole parser owner; no cross-module contract change.

## Price correctness claim — RED/GREEN
- RED before production source mutation: two focused tests failed on the installed
  parser, category image count contaminating price and detail address contaminating
  price. `regression-red.log` contains the result.
- GREEN: same two tests pass after domPrice scopes extraction to price metadata,
  marked nodes/SidebarPrice and price-only lines.
- Native gate: Docker builder executes scripts/build.sh: formatting, all tests,
  vet, CGO/SQLite build and linkage/version. `local-build.log` passed.
- Offline replay: every one of 60 saved category prices and 38 cached detail prices
  matches the independently inspected dedicated price blocks. `offline-prices-green.log`.
- Actual changed source: internal/somon/parser.go and parser_test.go; two documentation
  files in isolated release. Same parser changes copied to original working tree
  without touching unfinished app.go or other user work.

## Production preservation claim
Fresh read-only preflight: production-before.json. One active systemd watcher,
clean installed-baseline checkout, SQLite integrity OK, zero watcher restarts.
Eight existing containers and unrelated services/listeners/routes/firewall/SELinux
fingerprinted without publishing inventory.
Deliberately inducing loss or unrelated-host harm is not applicable; before/after
state and identity comparisons provide the accepted alternative proof.

## Deployment procedure
/tmp/somon-hotfix-release.py: exact branch commit fetch/fast-forward, resource-limited
native target build, service-user doctor using staged binary, existing root-only
backup route, stop/install/start only watcher. Exact settings/seen rows/environment
and database inode checked while stopped; subset preservation after start. Old binary
backup and watcher-only rollback retained. No database reset, filter mutation, manual
Telegram sending, unrelated workload restart or network/security modification.

## Release source and co-review
- Exact pushed hotfix commit: d5fdc120400acba96d088c465aa71db9ab1402ab; remote branch ref independently resolved to the same ID.
- Parser-focused independent co-review (Codex Luna xhigh): candidate_findings: none.
- Concrete release and host-state probes retained alongside this protocol.

## Production execution and postflight
- Target scripts/build.sh passed all tests, vet and native linkage; doctor passed
  SQLite, Telegram identity/target and 60 Somon cards before watcher stop.
- Exact deployed commit d5fdc120400acba96d088c465aa71db9ab1402ab; installed/running/built
  SHA256 ea25d21beefc71e1ccbdd3dec749805598ef35accaecf318e522ae70239895b1.
- Backup /var/backups/somonwatch/price-hotfix-20260909T085033Z contains verified old
  7f9c5f50d659 executable and intact pre-release database/environment.
- Stopped installation preserved exact settings, 3775 seen rows, environment and
  database inode. Fresh independent postflight compared backup rows to live DB:
  every row preserved; 3786 current rows after normal polling; settings/env unchanged.
- First poll completed: 60 cards, 11 new IDs, 3 details, 0 sends; no selected error logs.
- Container/service/listener/route/SELinux/failed-unit fingerprints matched preflight.
- Measurement limitation: original firewall digest included changing policy packet
  counters. It does not prove pre/post rule equality. Corrected rules-only fingerprints
  match across two postflight reads; executed commands contain no firewall modification.
  This is a probe defect, not evidence of a firewall change; do not report that the
  original firewall comparison passed.

## Co-review adjudication
- Stop/rollback error handling and retry-from-prepared-source were corrected before
  execution. Fresh identity probe bound PID, executable, cgroup and actual SQLite FD.
- Final host comparison is a separate mandatory step, not the script installation label.
- restorecon cleanup was moved inside try/finally before execution; staged file removed.
- A requested extra baseline-version guard is not needed to establish this completed
  run: the fresh preflight identified the exact baseline version and final verification
  independently verified the preserved rollback executable is that version. No drift
  was observed. No unexecuted generic deployment improvements were added.

## Attempt 2 — supported labeled price compatibility
Existing detail_unknown_seller fixture previously returned4200 but attempt1 returnednil. Added price assertion to existing test: labeled-price-red.log fails before correction, labeled-price-green.log passes after stripping only explicit Цена: before price-only matching. Native local gate and 60/38 offline replay passed again. Final hotfix commit93cbe9ebcfe6421040461f4ee4e049430bf5028d is pushed. Attempt1 receipts/verdicts are supporting-only under attempt-1/. Normalized firewall probe captured a fresh preflight before this final deployment.

## Final attempt2 result
Final native target build/doctor passed. Only watcher stop/install/start completed;
3786 stopped-state rows/settings/env/inode preserved. Backup:
/var/backups/somonwatch/price-hotfix-20260909T090439Z.
Fresh final-runtime-proof verifies exact93cbe9e running binary, valid rollback d5fdc12,
all3786 rows preserved among3792, settings/env unchanged and zero restarts.
ALL normalized host fingerprints now match pre/post final release, including firewall.
First cycle:60 cards,6 new IDs,1 detail,1 sent (17094676 at5000, seller_ads1).
