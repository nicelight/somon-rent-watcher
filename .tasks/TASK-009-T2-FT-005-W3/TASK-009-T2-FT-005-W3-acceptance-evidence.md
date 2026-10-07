# TASK-009-T2-FT-005-W3 — execution evidence

Owned claims: FT-005-AC-005 / REQ012; FT-005-AC-006 / REQ013. Dependency source/filter/management proof remains with006/007/008.

| Claim | Decisive observations | Artifact |
|---|---|---|
| AC005 RED | Two real created/enabled searches; first/next polls rental-only, no group/keyword delivery/history path | baseline.log, baseline probe, initial-source-snapshot.json |
| AC006 RED | Keyword polling/retry lifecycle absent; existing manual single-flight/backoff and rental cap initially GREEN | baseline.log |
| AC005 GREEN | Shared rental-seen1001 delivered independently to both; new1003 delivered; rejected1002 revision2 reevaluated after budget edit to delivered revision3; restart/edit/drop create no repeats; correct phrase/title/price/city/photo/link | green.log, TestKeywordPollingFirstNewIndependentRestartRevision |
| AC006 errors/retry | Source500/actual search parse error/detail500/Telegram500/ambiguous malformed success response keep rows empty, after reset+reopen deliver1001; send-success/write-failure permits 2 sends then durable success | green.log failures and ambiguity cases |
| AC006 shared constraints | cap1 across rental+2 searches gives details1/1/0, start rotates search1→search2→rental, no starvation; max source HTTP in-flight1; configured5ms, observed min arrival spacing4.746043ms within1ms measurement tolerance; unchanged client wait; HTTP403/429 shared backoff stops following search; manual cannot bypass;5 running duplicates rejected | green.log rotation/backoff/single-flight cases |
| AC006 revision | Detail barrier yields old-budget rejection after new revision; conditional write leaves no stale reject; old eligible work also cannot send after revision change; next poll delivers current revision | green.log stale barrier cases |

History comparison after reopen: each search has1001/1003 delivered at revision2 and1002 delivered at revision3; final total6 sends, no repeat. Original delivered revisions remain stable after edits. Failure rows remain absent until confirmed success.

Required package gate PASS: app/store/Telegram. Final claim probe PASS with race detector. Supporting package race/non-CGO build PASS; exact commands, exit codes and checksums in commands-results.json. Source basis is full initial/final byte hash snapshots; HEAD alone is insufficient and no commits were created.

One fixture correction: broken detail HTML is accepted by existing fallback parser; removed that false error fixture without changing the adapter. Actual source parse-error proof and detail HTTP-error proof remain. Papercut records evidence. No exactly-once or live compatibility claim.

All probes local httptest + disposable t.TempDir DB via cfg.DBPath; containers network none. Cancellation/close cleanup included. Status stays in_progress; next `/verify TASK-009-T2-FT-005-W3`.
