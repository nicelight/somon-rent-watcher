---
description: Make Telegram administration visible through acknowledged callbacks and fresh messages.
status: draft
lifecycle: planned
spec_design_status: complete
spec_design_links:
  - ".memory-bank/architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable"
  - ".memory-bank/architecture/system-architecture.md#ad-003--telegram-administration-is-append-only"
  - ".memory-bank/contracts/boundary-map.md#polling-orchestration-contract"
  - ".memory-bank/contracts/boundary-map.md#telegram-application-boundary"
  - ".memory-bank/contracts/boundary-map.md#external-boundary-rules"
  - ".memory-bank/states/runtime-lifecycle.md#telegram-admin-transitions"
  - ".memory-bank/states/runtime-lifecycle.md#non-durable-process-state"
  - ".memory-bank/invariants.md#accepted-must"
  - ".memory-bank/invariants.md#accepted-never"
  - ".memory-bank/invariants.md#existing-compatibility-guardrails"
  - ".memory-bank/testing/strategy.md#risk-based-checks"
  - ".memory-bank/testing/current-coverage.md#automated-coverage-by-package"
---

# FT-002 — Append-only Telegram administration

## Value and use cases

Administrators can navigate and change the shared filter from Telegram clients that display edited bot messages unreliably. The newest bot message represents the current menu/state; no client detection or cleanup machinery is needed.

## Source and constraints

- [.memory-bank/prd.md#req-002--append-only-telegram-admin-controls](../prd.md#req-002--append-only-telegram-admin-controls)
- Preserve admin/target-chat authorization, per-admin/chat pending text input, current shared settings, and manual-scan single-flight/backoff behavior.

## Edge and failure behavior

- Old menu callbacks act on current persisted settings and produce a new result message.
- Invalid authorized callbacks retain bounded rejection feedback; unauthorized
  users and callbacks from other group chats remain silent and non-mutating.
- Manual scan reports start/completion/error without editing a busy button.

## Acceptance Criteria

### FT-002-AC-001 — Authorized callback acknowledgement and boundary rejection

- REQ: REQ-002
- Observable criterion: every authorized navigation or mutation callback in private/target chat is acknowledged promptly and produces a new menu/state message; unauthorized users and callbacks from other group chats produce no settings mutation or bot output; runtime callback handling makes zero `editMessageText` requests.
- Verification method: Telegram HTTP test server captures method order/request counts while backend assertions cover authorized, unauthorized, private, target-chat, and wrong-chat cases.

### FT-002-AC-002 — Text input returns a fresh main menu

- REQ: REQ-002
- Observable criterion: valid price, negative-word, and interval input persists settings and sends a new main menu while pending input remains isolated by administrator and by chat.
- Verification method: bot tests for all input modes plus cross-admin and cross-chat isolation.

### FT-002-AC-003 — Manual scan uses separate feedback

- REQ: REQ-001, REQ-002
- Observable criterion: accepted manual scan sends start feedback and completion/error feedback without message edits while preserving single-flight, pause, and backoff rules.
- Verification method: bot/app tests with captured Telegram methods and scheduler request state.

### FT-002-AC-004 — KISS interaction model

- REQ: REQ-006
- Observable criterion: no client detection, edit retry/cache, old-message deletion, new persistence, or configuration is introduced.
- Verification method: code/spec review and repository-native build gate.

## Acceptance closure

- Navigation/mutation visibility, stale menus, text input, authorization isolation, manual scan completion/error, and zero edits are covered by FT-002-AC-001 through FT-002-AC-003.
- Simplicity constraints are covered by FT-002-AC-004.

## SDD Design

Feature design is complete by reusing the accepted owners; no feature-local
technical spec, ADR, storage/configuration change, compatibility layer, client
detection, edit retry/cache, or old-message cleanup mechanism is required.

- Architecture Spine: [AD-001](../architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable) and [AD-003](../architecture/system-architecture.md#ad-003--telegram-administration-is-append-only).
- Boundaries: [Telegram administration and delivery](../contracts/boundary-map.md#telegram-application-boundary), the [Telegram Bot API rule](../contracts/boundary-map.md#external-boundary-rules), and [polling orchestration](../contracts/boundary-map.md#polling-orchestration-contract) for manual-scan single-flight/backoff.
- State and data: [Telegram admin transitions](../states/runtime-lifecycle.md#telegram-admin-transitions), [non-durable pending/manual state](../states/runtime-lifecycle.md#non-durable-process-state), and the accepted [MUST](../invariants.md#accepted-must), [NEVER](../invariants.md#accepted-never), and [compatibility](../invariants.md#existing-compatibility-guardrails) rules; existing settings and Telegram offset persistence remain unchanged.
- Verification: the AC methods above use the [risk-based testing policy](../testing/strategy.md#risk-based-checks) and extend the existing [Telegram/application HTTP test-server harness](../testing/current-coverage.md#automated-coverage-by-package) with method-order, fresh-message, authorization, pending-input isolation, and manual-scan feedback assertions.

Global Backbone Planning Revision 1 is complete and Foundation is
`not_required`.
