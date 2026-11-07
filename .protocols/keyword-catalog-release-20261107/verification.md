# Keyword category release — 2026-11-07

Operator authorized production deployment via ssh igorprod.

- Exact release `4ca196aaa73820fdf8bff792af6ce544a4b556c3` published and installed; built/installed/running checksums matched. Native local and target format/test/vet/build/linkage gates passed. Telegram integration checks cover the new menu, removal of tables/chairs, category selection and city preservation; URL tests cover every category/city combination.
- Before and after installation doctor passed using disposable SQLite clones. Only somonwatch.service was stopped/installed/started. Single active process, zero restarts; first paused poll succeeded with 60 cards and two fresh IDs.
- Database identity, exact stopped-state data, settings/environment and previous seen/history rows preserved. Seen 23159 → 23161. No saved searches existed in the procedure's initial/stopped/post-start snapshots; none was created, modified or activated by deployment. Preflight rejects a tables_chairs search rather than silently changing its scope.
- All unrelated host fingerprints matched: [comparison](comparison.json), [preflight](immediate-preflight.json), [postflight](postflight.json). [Receipt](release-receipt.json), [procedure](release-procedure.py).
- Rollback backup `/var/backups/somonwatch/keyword-catalog-release-20261107T125127Z`; stage and doctor scratch cleaned. No database restoration or replacement.
