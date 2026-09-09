# PRD bootstrap decision log

- Existing `seen_ads` semantics remain authoritative; no shown-history feature was created.
- Fallback selection and Telegram interaction are separate features because each has its own user-visible result and verification boundary.
- Production release is a separate feature because it has an independently blocking outcome, operator authorization, runtime boundary, and rollback/isolation acceptance.
- No unresolved decomposition decision remains.
