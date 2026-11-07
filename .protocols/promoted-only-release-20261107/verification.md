# Promoted-only category release verification

Operator requested careful deployment via `ssh igorprod` on 2026-11-07.

- Release: `c2baf63de3be2f4ea9a5ec9cce2f66d0e94488a1`; local/remote ref, production checkout, built/installed/running executable confirmed by the release procedure.
- Local and target native formatting/test/vet/build/linkage gates PASS. Doctor PASS before and after installation using disposable SQLite clones, without baseline mutation or message delivery.
- Only the identified `somonwatch.service` was stopped/installed/started. One active process, zero restarts. First cycle succeeded: paused baseline advanced, 60 cards, one fresh ID.
- SQLite integrity OK; unchanged database inode, settings and environment; every prior seen row preserved, 23158 → 23159. Existing keyword monitor and history retained. No activation or reset.
- All unrelated containers/services/listeners/routes/firewall/SELinux/failed-unit fingerprints match: [comparison](comparison.json), [fresh preflight](immediate-preflight.json), [postflight](postflight.json).
- Rollback backup: `/var/backups/somonwatch/promoted-only-release-20261107T121709Z`. No database restore or replacement occurred. Stage and doctor scratch removed.
- [Receipt](release-receipt.json), [procedure](release-procedure.py). Build/release console logs remain local ignored operational logs.

Live pages at doctor time contained four ordinary cards; all-promoted acceptance is covered by the deterministic app regression test. No claim that the live first cycle reproduced that composition.
