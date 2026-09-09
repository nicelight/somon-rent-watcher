# FT-001 planning decisions

## Authoritative decisions applied

- Reuse `Settings.PriceMax`, `seen_ads`, the existing scheduler, shared detail
  cap, Somon client, Telegram delivery, and current package boundaries.
- A completed exact-empty evaluation is the only entry to fallback; any exact
  match or incomplete exact evaluation suppresses fallback.
- Fallback remains in-memory, contains at most three otherwise eligible fresh
  ads in `(PriceMax, floor(1.5 * PriceMax)]`, and uses ascending price with
  stable feed-order ties.
- Exact and fallback delivery retain the current confirmed-delivery-before-seen
  rule; final rejection becomes seen and deferred/transient/failed delivery
  remains unseen.
- No schema, setting, dependency, worker, compatibility layer, separate
  history, or duplicate candidate/state-transition owner is authorized.

## Planning decision

The unattended execution-cohesion pass selected one T2 task. The polling
behavior and its regression proof are one lifecycle change centered on
`internal/app`; splitting filter support or tests would not create an
independently completable product outcome.

No unresolved operator decision or blocker remains.
