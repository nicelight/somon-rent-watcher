# FT-003 planning decisions

## Authoritative decisions applied

- One final T3 production-acceptance task owns the indivisible ordered release
  and depends on every local implementation task.
- The exact order is local repository-native gate; ordinary commit/push and
  exact GitHub-ref confirmation; fresh read-only production preflight; safe
  exact-commit fast-forward synchronization; target build and redacted doctor;
  stop/install/start of `somonwatch.service` only; then version, health, logs,
  SQLite/settings preservation, and unrelated-host comparison.
- The operator phrase `sh igorprod` and prior evidence of an SSH alias are only
  discovery inputs. Execution must resolve the actual access command from
  current read-only local evidence before connecting and must stop rather than
  infer or repair an ambiguous route.
- The only accepted runtime shape is one healthy host systemd unit named
  `somonwatch.service`, bound consistently to its process, executable,
  environment, SQLite path, and production checkout. Missing, competing,
  ambiguous, non-systemd, or unhealthy evidence stops before production write.
- Production Git synchronization is non-destructive: no force, reset, stash,
  discarded local work, or ref rewrite. Dirty, diverged, non-fast-forward, or
  post-sync commit mismatch stops.
- No container operation, runtime migration, database mutation/migration/
  replacement, broad host change, unrelated service action, or secret-bearing
  receipt is authorized.

## Task boundary decision

Publication, production preflight/synchronization, runtime update, and
postflight remain one task because they jointly produce and prove one release;
splitting would leave an unaccepted partially published or partially deployed
state without independent product value.

No unresolved operator decision or blocker remains for planning. Runtime or
access ambiguity discovered during execution is an explicit stop condition,
not planner discretion.
