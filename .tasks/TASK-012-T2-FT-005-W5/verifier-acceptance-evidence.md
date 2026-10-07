# Independent verification — TASK-012-T2-FT-005-W5

Reviewer independently proved FT-005-AC-009, with REQ-012/REQ-014 as governing constraints. No other feature claims were adopted.

| Observation | Decisive comparison | Evidence |
|---|---|---|
| Current feed history selection | First monitor: requested [202,101,202,999] returns exactly states 101/202 of 4 total, never historical 303/404 or missing 999. Second monitor has different delivered/revision values for shared IDs and returns only its own states. Missing monitor returns empty. | `verifier_history_probe_test.go`, `verifier-outcome.log` |
| Empty selection | Nil and zero-length ID lists return empty; same calls on a nil DB receiver succeed, demonstrating no DB access. | Store probe, `verifier-outcome.log` |
| History preserved | Exact full maps for both monitors (7 total rows, including delivered flags and evaluated revisions) remain equal after bounded reads, a revision-changing edit, and close/reopen. Full diagnostic reader stays compatible. | Store probe, `verifier-outcome.log` |
| Production polling lifecycle | 3 current cards alongside 2 off-feed rows: delivered ID and current-revision rejection skipped; fresh rejection stored. After budget edit/reopen, only the two current rejections receive details/delivery, current delivered ID has no detail request. Off-feed states and other monitor unchanged. Another reopen yields zero new/details/sends and identical map. | `verifier_polling_history_probe_test.go`, `verifier-outcome.log` |
| Production caller/query | App supplies selected monitor and IDs directly from current cards; Store SQL binds monitor_id and requested ad_id IN parameters. Production polling has no full-history reader call; no schema/history write/cleanup delta. Existing accepted App -> Persistence edge retained. | `verifier-source-inspection.md`, source snapshot |

Reproduce all fresh checks: `python3 .tasks/TASK-012-T2-FT-005-W5/run-verifier.py`.
The runner copies cmd/internal/scripts/testdata/go.mod/VERSION into a disposable task-local directory, runs network-none Docker with existing builder image, runs required focused gate and original native build gate before injecting verifier-owned probe files into the copy, then runs those probes. Copies, test DBs/servers and build output are removed; repository source is never edited.

- Fresh required focused gate: somon/app/store all passed (`verifier-focused-gate.log`).
- Final combined source native gate: gofmt/all tests/vet/CGO build/SQLite linkage/version/checksum passed (`verifier-native-gate.log`). This is supporting regression evidence, not feature completion.
- Fresh independent outcome probes: both passed (`verifier-outcome.log`).
- Commands/image/completion/exit codes: `verifier-commands-results.json`. Source hashes: `verifier-source-state.json`; executor handoff's four hashes match this independently tested snapshot, and source remained unchanged afterward.
- Executor claim path: actual pre-change baseline returned three historic rows for a current feed of one, exit1 (`baseline.log`, `baseline-test.go.txt`); claim-equivalent final bounded tests and required gate passed (`green-rerun.log`, execution `progress.md`). These remain supporting evidence, not verifier proof.
- No execute receipt reused. Direct reruns were cheap and no current eligible receipt was offered.
- Earlier setup runs failed only because Docker does not share host `/tmp`. They executed no tests; corrected isolation location is inside shared task artifacts. Original setup logs/commands retained as `verifier-setup-*`.
- Executor's initial timing test failure (3.746916ms vs 4ms) did not recur in either fresh focused/native gate. No product finding is inferred from it.
- Loaded finding-adjudication pack and review-code instructions. One requested GPT-6.1-sol/xhigh co-review launch (explicit parent override of Luna) failed due agent thread limit. No retry under `/verify` instructions; absence is best-effort availability, not missing functional evidence.

Functional outcome passes; T2 closure remains with root. Status, source/specs/Git/production untouched by verifier.
