---
description: Product Requirements Document.
status: draft
type: prd
clarification_status: complete
constitution_checked: true
---

# Product Requirements Document — Somon Rent Watcher

## Source Inputs

- [.memory-bank/analysis/product-brief.md](analysis/product-brief.md), `Decision: proceed`.
- [.memory-bank/constitution.md](constitution.md), explicit interview skip with existing framework principles and project KISS policy.
- [.memory-bank/product.md](product.md), evidence-backed brownfield baseline.
- Operator decisions of 2026-09-04: append-only Telegram controls, limited higher-price fallback, ordered GitHub/server/runtime deployment, and no unrelated production impact.

## Product Summary

Somon Rent Watcher monitors fresh apartment-rental listings, applies one shared filter, and notifies one Telegram group without normal duplicate delivery. The accepted delta keeps the strict filter primary, adds a bounded price fallback only when a completed poll finds no fresh exact match, and makes admin callbacks reliable across Telegram clients by sending new messages instead of editing old ones.

## Goals

- Reduce empty-result periods without turning the price limit into an unbounded suggestion system.
- Keep exact matches visibly and logically primary.
- Preserve first-seen, pause, backoff, recovery, request-cap, and retry behavior.
- Make every inline action produce client-visible feedback without message edits.
- Release through a verified, isolated path that affects only Somon Rent Watcher.

## Non-goals

- No separate shown-history table, SQLite migration, backfill, or reprocessing of existing `seen` IDs.
- No configurable fallback count/percentage, client detection, edit retry, menu cleanup, or multiple profiles/groups.
- No changes to Somon crawling scope, anti-blocking behavior, unrelated host services, containers, networks, firewall, reverse proxy, or shared databases.

## Users / Actors

- Configured administrators manage one shared filter in private chat or the configured target group.
- Group participants receive exact or fallback ad notifications.
- Somon.tj supplies category/detail HTML; Telegram Bot API supplies updates and message delivery.

## Functional Requirements

### REQ-001 — Preserve the monitoring baseline

The service shall retain its current baseline, first-seen, pause, recovery, backoff, detail-request-cap, and confirmed-delivery-before-seen behavior except where REQ-003 explicitly classifies fresh higher-price candidates inside the same poll.

### REQ-002 — Append-only Telegram admin controls

1. Every authorized inline callback shall be acknowledged promptly with `answerCallbackQuery`.
2. Navigation and filter mutations shall send the current menu/state as a new message and shall not call `editMessageText`.
3. Price, negative-word, and interval input completion shall send a fresh main menu after persisting valid input.
4. Manual scan shall use separate start/completion feedback and shall not edit or restore a busy button.
5. Old menus remain visible and actionable; the result of any action is represented by the newest bot message.

### REQ-003 — Bounded higher-price fallback

1. Exact candidates shall be evaluated before fallback candidates.
2. Fallback is eligible only when the poll finds no fresh ad satisfying the complete current filter and exact-price range.
3. A fallback ad shall satisfy every current filter condition except `PriceMax`, shall have a known price strictly greater than `PriceMax`, and shall have a price no greater than `floor(PriceMax * 1.5)`.
4. When `PriceMax` is absent, fallback shall not run.
5. At most three fallback ads shall be sent, ordered by ascending price; stable feed order breaks equal-price ties.
6. The existing `MaxDetailsPerPoll` limit applies to exact and fallback detail requests together. If the exact-candidate evaluation is incomplete because the cap is reached, the service shall not claim an empty exact result and shall defer fallback to a later poll.

### REQ-004 — Freshness, deduplication, and failures

1. Only IDs absent from `seen_ads` at the beginning of the poll are eligible for exact or fallback delivery.
2. Existing `seen` IDs are ignored regardless of whether they were previously delivered or rejected.
3. Fresh cards rejected by all applicable paths are marked `seen` under the existing first-seen policy; cards deferred by the detail cap remain unseen.
4. An exact ad that passes filtering suppresses fallback for that poll even if Telegram delivery fails or is ambiguous.
5. A delivery failure leaves that ad unseen for the existing retry behavior.
6. If neither an exact ad nor an eligible fresh fallback ad exists, no ad message is sent.

### REQ-005 — Ordered and isolated release

1. Applicable local tests, vet, formatting, CGO build, and linkage checks shall pass before publication.
2. Release order shall be local verification, commit/push to the configured GitHub remote, read-only production preflight, production git synchronization, target-compatible build, then update/restart only the verified Somon Rent Watcher runtime.
3. Production settings and SQLite state shall be preserved.
4. Runtime shape shall be detected from current host evidence immediately before deployment; no unrelated container or host-service operation is permitted.
5. Post-deploy checks shall cover version, doctor/runtime health, logs, and unchanged unrelated-workload state.

## Non-functional Requirements

### NFR-001 — Simplicity and maintainability

The change shall reuse existing settings, `seen_ads`, scheduler, clients, and request cap without new dependencies, schema, configuration options, background workers, or compatibility layers. Verification class: code/spec review plus repository-native build gate.

### NFR-002 — Reliability

No fallback path may hide an exact match, mark a failed delivery as seen, bypass backoff, or exceed the shared detail-request cap. Verification class: deterministic unit/integration tests with mixed exact/fallback candidates and delivery failure.

### NFR-003 — Production isolation

Deployment shall not alter the state of unrelated containers, services, sockets, firewall, routing, SELinux, or databases. Verification class: before/after read-only host snapshot and scoped Somon Watcher health checks.

## Data / Domain Model

- `Card` remains category-listing data; `Ad` remains detail-enriched data.
- `Settings.PriceMax` is the sole fallback basis; fallback ceiling is derived in memory and is not persisted.
- `seen_ads` remains the sole freshness/deduplication store. No distinction between rejected, exact-delivered, and fallback-delivered IDs is introduced.
- A poll internally distinguishes exact candidates, fallback candidates, and rejected cards without adding durable lifecycle state.

## UX / Interaction Flow

1. An administrator invokes `/filter` or `/status` and receives a new menu message.
2. An inline press is acknowledged immediately; the bot performs the action and sends the resulting submenu/main menu as a new message.
3. Text-input actions send their prompt, accept input scoped by administrator and chat, persist it, and send a fresh main menu.
4. A manual scan posts start feedback and later posts completion/error feedback; matching ads remain normal separate notifications.
5. During automatic/manual polling, exact ads are sent normally. Only an exact-empty completed evaluation may emit up to three closest-price fallback ads.

## Integrations / Dependencies

- Telegram Bot API: retain `getUpdates`, `answerCallbackQuery`, `sendMessage`, and `sendPhoto`; `editMessageText` is no longer used by runtime bot flows.
- Somon.tj: retain sequential category/detail fetching, delay, body limit, sanity validation, and block handling.
- SQLite: retain existing schema and transactions.
- Git/GitHub and production host: release only through the operator-authorized ordered route and current runtime evidence.

## Edge Cases / Failure Handling

- Missing `PriceMax` or unknown candidate price disables eligibility for fallback.
- A candidate exactly at `PriceMax` is exact; exactly at `150%` is eligible fallback.
- Equal-price fallback candidates preserve stable feed priority.
- Detail parse/network errors remain unseen and do not prove that exact evaluation was empty.
- HTTP 403/429 still aborts into existing backoff; fallback does not bypass it.
- Detail-cap exhaustion leaves unprocessed IDs unseen and suppresses fallback for that cycle.
- Telegram delivery errors preserve unseen retry semantics and are not replaced with alternative fallback output.
- Old Telegram menu presses operate against current persisted settings and produce a new canonical message.

## Acceptance Criteria

- **PRD-AC-001:** With at least one fresh fully matching ad, no over-price ad is sent in that poll.
- **PRD-AC-002:** With no exact match, eligible fresh ads priced above `PriceMax` and at or below `150%` are sent in ascending-price order, capped at three.
- **PRD-AC-003:** Ineligible, previously seen, over-ceiling, unknown-price, and fourth-or-later fallback ads are not sent.
- **PRD-AC-004:** Exact/fallback detail work respects one shared request cap; incomplete exact evaluation suppresses fallback and preserves deferred IDs.
- **PRD-AC-005:** Delivery failure leaves the affected ID unseen, and an exact filter match still suppresses fallback.
- **PRD-AC-006:** Authorized callbacks are acknowledged and result in new messages; runtime callback flows make zero `editMessageText` calls.
- **PRD-AC-007:** Text-input completion and manual scan use fresh messages while preserving authorization, per-admin/chat pending input, single-flight, pause, and backoff behavior.
- **PRD-AC-008:** Repository-native gates pass before push; production preflight and postflight prove only Somon Rent Watcher changed and its persisted state survived.

## Verification Strategy

- Filter/unit probes for inclusive/exclusive price boundaries and the derived ceiling.
- App integration tests using temporary SQLite and HTTP test servers for exact suppression, three-item ordering/cap, seen exclusion, detail cap, and delivery failure.
- Telegram HTTP test-server assertions for callback acknowledgement, fresh menu sends, zero message edits, text completion, and manual scan feedback.
- `scripts/build.sh` or the documented Docker equivalent for formatting, tests, vet, CGO build, linkage, and version evidence.
- Read-only production preflight, target-host build/doctor, scoped service update, and postflight comparison.

## Clarifications

### 2026-09-04

- The operator accepted KISS and explicitly skipped a separate Constitution interview.
- Fallback uses the current configured `PriceMax`, a fixed three-item cap, and a fixed inclusive `150%` ceiling.
- Existing `seen_ads` remains the freshness authority; no shown-history model or backfill is added.
- Telegram UI becomes append-only rather than adding client detection or edit retries.
- Deployment remains sequential and isolated; actual runtime type is verified immediately before the production write.

## Constitution Amendment Candidates

None.

## Unresolved Blockers

None.

## Accepted production hotfix (2026-09-09)

After reporting missing listings and receiving the live per-ad diagnostic, the operator
explicitly requested correcting inflated price extraction and carefully deploying on
the shared production host. REQ-009/FT-004 record this bounded correction. Existing
REQ-005 and REQ-008 govern release/isolation; settings/seen history remain unchanged.
Pending fallback and Telegram UI work is not included in this isolated hotfix.
