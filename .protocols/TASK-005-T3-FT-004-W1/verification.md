# Verification — TASK-005-T3-FT-004-W1

## Current attempt
- attempt: 2
- Final commit: 93cbe9ebcfe6421040461f4ee4e049430bf5028d.
- Attempt1 receipts/verdicts under attempt-1/ are supporting-only.

## Scope and normative basis
Operator-authorized price extraction correction and isolated shared-host deployment.
Somon adapter owns HTML conversion. Existing filtering/seen/settings/delivery behavior
is preserved; only dedicated DOM price extraction plus explicit Цена: compatibility.

## Functional proof
- Original image-count/address RED is retained in regression-red.log.
- Existing Цена: fixture price assertion failed (nil) on attempt1 before correction:
  labeled-price-red.log; same fixture now passes: labeled-price-green.log.
- Full native local and target tests/vet/build passed after final source changes:
  local-build.log and production-release.log.
- Offline replay: all60 category and38 cached detail prices match independent
  visible prices (offline-prices-green.log).
- Staged target doctor passed before stop: SQLite, Telegram,60 category cards.
- Fresh verifier-owned final-runtime-proof.json: running/installed/built SHA
  c99616820cc979fb6257cfa959dfcc323b4b9ea487ecf92719ec548c8503e9d8;
  exact clean source93cbe9e; SQLite integrity; every3786 pre-final-release row remains
  among3792 rows; settings/env unchanged; backup old executable verified; zero restarts.
- Fresh normalized host comparison BEFORE and AFTER final release passes ALL9 fields:
  containers, services, listeners, routes, firewall rules, SELinux, failed units,
  settings and environment. See production-before.json, production-after.json,
  host-comparison.json. The firewall probe now removes timestamps AND traffic counters.
  Attempt1's invalid firewall hash remains historical evidence only, not a PASS claim.
- First post-final-release cycle completed with6 new IDs,1 detail,1 successful delivery
  (ID17094676, price5000, seller count1). No selected poll-error logs.

## Review
Prior independent parser/deployment co-reviews remain supporting evidence. The final
bounded Цена: compatibility delta and fresh final deployment receipts route to independent
semantic re-verification. No broader feature/architecture change is claimed.

VERDICT: PASS

Next: final per-task semantic verification then explicit-owner closure.
