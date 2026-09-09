---
description: Reliable matching, Telegram control, and isolated delivery of the current watcher upgrade.
status: active
lifecycle: planned
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
---

# EP-001 — Reliable watcher upgrade

## Value

Give the group useful nearby-price options during genuinely empty cycles and make administration dependable across Telegram clients, without weakening the watcher's first-seen safety or production isolation.

## Requirements

- REQ-001 through REQ-008 in [.memory-bank/requirements.md](../requirements.md).

## Features

- [FT-001 — Closest higher-price alternatives](../features/FT-001-price-fallback.md)
- [FT-002 — Append-only Telegram administration](../features/FT-002-append-only-telegram-ui.md)
- [FT-003 — Isolated production release](../features/FT-003-isolated-production-release.md)

## Success metrics

- Exact matches always win; exact-empty polls produce zero to three eligible closest-price fresh alternatives.
- Authorized Telegram callbacks produce prompt acknowledgement and a fresh visible result without message edits.
- Local and target gates pass, persisted watcher state survives, and unrelated production workloads remain unchanged.

## Acceptance criteria

- Every feature AC is verified through its declared method.
- Product behavior matches `.memory-bank/prd.md` without a new storage/configuration/dependency layer.
- The release follows the approved ordered path and isolates Somon Rent Watcher.

## Open questions

None.

## Accepted maintenance follow-up

- [FT-004 — Price extraction hotfix](../features/FT-004-price-extraction-hotfix.md): operator-authorized correction completed on2026-09-09; pending FT-001..003 lifecycle remains unchanged.
