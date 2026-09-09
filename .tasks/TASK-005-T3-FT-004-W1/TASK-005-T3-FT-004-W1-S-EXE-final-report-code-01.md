# Somon price extraction hotfix

- Final installed commit: `93cbe9ebcfe6421040461f4ee4e049430bf5028d`.
- Published branch: `hotfix/somon-price-extraction`.
- Release isolated from unrelated unfinished local feature work.
- Parser uses scoped price metadata/nodes and standalone price lines, including
  the existing explicit `Цена:` form. Photo counts and address numbers no longer
  inflate prices. Existing fixture price compatibility is retained.
- Regression RED/GREEN, native local and target formatting/tests/vet/CGO build,
  staged target doctor and cached replay (60 category/38 detail prices) passed.
- Running/installed/built binary SHA256:
  `c99616820cc979fb6257cfa959dfcc323b4b9ea487ecf92719ec548c8503e9d8`.
- Latest rollback backup: `/var/backups/somonwatch/price-hotfix-20260909T090439Z`
  (root-only), previous executable version independently verified. Original
  pre-hotfix backup is also retained at `price-hotfix-20260909T085033Z`.
- SQLite/settings/environment preserved; all 3786 pre-final-release seen rows
  remain among 3792 current rows. Zero watcher restarts.
- First final-version cycle: 60 cards, 6 new IDs, 1 detail, 1 successful delivery
  (ID 17094676 at 5000 somoni, seller count 1), no selected poll error.
- Eight containers, other running services, listeners, routes, firewall rules,
  SELinux and failed-unit state matched normalized pre/post final-release probes.
  Historical attempt1 firewall-counter measurement limitation is retained in
  attempt-1/ receipts; final comparison used the corrected method on both sides.
- Seen history is preserved; previously seen ads are not automatically replayed.

Current execution and functional receipts:
[protocol](../../.protocols/TASK-005-T3-FT-004-W1/progress.md) and
[verification](../../.protocols/TASK-005-T3-FT-004-W1/verification.md).
