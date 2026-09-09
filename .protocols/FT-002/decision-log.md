# FT-002 planning decisions

## Authoritative decisions applied

- Authorized callbacks are acknowledged before their action and represent the
  resulting menu/state with a fresh `sendMessage`; runtime control flows make
  zero `editMessageText` requests.
- Existing administrator allowlist and private/target-chat authorization are
  preserved. Unauthorized and wrong-group callbacks do not mutate settings or
  emit bot output; invalid authorized callbacks retain bounded rejection
  feedback without message edits.
- Valid price, negative-word, and interval input continue to use pending state
  isolated by administrator and chat, persist through the existing Backend,
  clear only that pending key, and send a fresh main menu.
- An accepted manual scan sends separate start and completion/error messages;
  it retains the existing capacity-one single-flight, pause, and backoff rules
  without a busy-button edit or restore edit.
- Old menus remain actionable against current settings. No old-message
  deletion, client detection, edit retry/cache, persistence, configuration,
  dependency, worker, or compatibility layer is authorized.

## Rebuild decision

The rejected one-task queue is rebuilt into two W1 tasks with no dependency:

- `TASK-002-T3-FT-002-W1` owns FT-002-AC-001, FT-002-AC-002, and
  FT-002-AC-004. T3 remains required because AC-001 changes observable behavior
  at the administrator authorization boundary. AC-002 is preserved behavior,
  not a separate implementation result. AC-004 is allocated here once as the
  KISS proof for the primary Telegram protocol conversion.
- `TASK-003-T2-FT-002-W1` owns FT-002-AC-003. It changes the manual-scan
  interaction/lifecycle output across the accepted Telegram/Polling
  Application boundary but does not change authorization decisions or
  permissions, so T2 is sufficient.

The tasks are independent implementation outcomes even though canonical
execution is sequential and their advisory Telegram paths overlap. Neither
task inherits or duplicates the other's exact claim proof.

The task requiring `./scripts/build.sh` keeps a non-empty hard boundary and
adds only `dist/somonwatch` and `dist/somonwatch.sha256`; the tracked
`dist/.gitkeep` makes broader `dist/` authorization unnecessary.

No unresolved operator decision or blocker remains.
