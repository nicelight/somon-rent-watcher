# TASK-008-T3-FT-005-W2 — acceptance evidence, attempt 1

Execution evidence only; independent functional/T3 semantic verdict remains due.

## FT-005-AC-001 / REQ-010, REQ-014, REQ-002

Baseline after durable in_progress and before production: existing Bot.processUpdate + real App + temporary SQLite receives authorized private ks:new. Callback fails specifically “неизвестная кнопка”; closed SQLite has no search_monitors schema/rows. Compiling behavioral assertion fails honestly. Source/log: keyword_search_baseline_test.go / baseline.log. baseline-setup-error.log contains an earlier pointer assignment compile failure: setup-only, excluded from RED.

Corresponding GREEN: keyword_search_integration_test.go invokes the same exact processUpdate entrypoint via test-only export with real App and temp SQLite. Private admin creates ID1, trims phrase, sees escaped summary All/country/disabled; category and city menus expose accepted readable catalog. Category change preserves Vose; bounds100–200 then explicit enable. Target-group second admin creates ID2 services/Dangara/min0/noMax, separately enabled. Addressed phrase and disable/enable preserve ID2 and advance revisions. Both full payloads persist through DB close/reopen/new App. Callback traces start answerCallbackQuery before every fresh output; no editMessageText. Initial main menu offers route. Observed final ID1 revision8, ID2 revision5 in claim-green.log.

Harm proof: nonadmin and wrong chat produce ack-only/no writes; empty phrase, invalid catalog, negative/inverted/fractional/noninteger/overflow input leave rows unchanged. Pending two admins/chats/IDs stay isolated; switching selected search replaces that admin/chat action; newer revision causes old pending input rejection; rental navigation/cancel clear search pending, search creation replaces rental pending without mutating rental values. Missing ID view/enable/edit returns fresh list with no resurrection. Valid administrative sharing follows unchanged allowlist policy; no per-search ownership model invented.

## FT-005-AC-007 / REQ-014

Accepted RED_NOT_APPLICABLE: preservation was already correct in baseline; absence cannot be observed without artificial destructive regression. Initial GREEN baseline exact closed SQLite bytes compare across reopen/auth/delivery, protected settings/seen/state/offset unchanged, same rental sendMessage payload.

Current GREEN: store test constructs a legacy-schema fixture using existing store writes and test-only removal of newly added empty search table, captures all settings/seen/state rows, reopens through new additive initializer, writes two searches/targeted revision edit, reopens again. Bidirectional SQLite EXCEPT compares every rental column including first_seen_at and all state keys, with CHECK failure on any mismatch. Result equal. Integration compares the same protected public snapshot and rental SendAd method/form (ad URL ID303) before search writes, after writes and after actual DB reopen/new App; equality. Existing rental/filter/fallback suites pass in all-package regression. No polling/history/deletion behavior claimed by this task.

REQ014 reviewer rubric: same Go process/SQLite/client/config/group/allowlist/scheduler composition; new business writes owned by App, SQL only Persistence; Telegram interface performs transport/input only; no new dependency/worker/infra or graph edge. Code/import assessment supports this, independent reviewer assessment remains due.

## Commands and results

- Baseline exact command in progress.md: exit1 expected honest create behavior RED; rental/auth test PASS.
- Claim GREEN verbose package suite: Docker network-none CGO `go test -count=1 -v ./internal/telegram ./internal/app ./internal/store`, exit0 claim-green.log.
- Required exact focused Docker gate: exit0 focused-gate.log.
- All-package tests, focused vet, CGO build, noCGO compilation, task-source formatting: exit0 local-gates.log; exact command in progress.md.
- Task tracked-source `git diff --check`: exit0.

Task-start source snapshot + change-surface.json prove bounded local delta; existing unrelated source and go.mod hashes unchanged. Twelve production/test files, four WHY/WHERE docs, task ready→in_progress and execution protocol/evidence only. No commit/push/deploy, production DB/secrets/live Telegram, new dependencies/workers, previous FT task statuses or deletion/history/poll outcomes touched.

## Isolation / safe rerun

All tests use fresh t.TempDir/search.db, local httptest, fake user IDs/TOKEN; DB/server close and test temporary cleanup. Required Docker network none; ephemeral build outputs /tmp inside container. Rerun creates new fixture state and never opens default/production DB_PATH.

## Handoff

Independent `/verify TASK-008-T3-FT-005-W2` followed by T3 `/red-verify TASK-008-T3-FT-005-W2`; /root is closure owner. Task remains in_progress. Executor evidence is supporting-only, no reuse receipt offered.
