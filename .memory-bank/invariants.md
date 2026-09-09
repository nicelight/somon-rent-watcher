---
description: Accepted product and architecture invariants for the current watcher delta.
status: active
last_verified: 2026-09-02
last_updated: 2026-09-04
source_of_truth:
  - .memory-bank/prd.md
  - .memory-bank/requirements.md
  - .memory-bank/architecture/system-architecture.md
---

# Invariants

## Accepted MUST

- Exact candidates MUST be evaluated before fallback, and any exact match MUST suppress fallback for that poll.
- Exact and fallback detail work MUST share the configured per-poll cap; incomplete exact evaluation MUST suppress fallback.
- Only IDs fresh at poll start MAY be delivered; final rejection/confirmed delivery becomes seen, while deferred/transient/failed-delivery IDs remain unseen.
- A fallback ad MUST pass every non-maximum-price filter and have known price in `(PriceMax, floor(1.5 * PriceMax)]`; delivery order is ascending price with stable feed-order ties and count is at most three.
- Authorized Telegram callbacks MUST be acknowledged promptly and their resulting current UI MUST be sent as a new message.
- Deployment MUST follow the accepted ordered route, preserve SQLite/settings, stop on ambiguous/unhealthy preflight, and change only the watcher runtime.

## Accepted NEVER

- NEVER run fallback without `PriceMax`, after an exact match, or after incomplete exact evaluation.
- NEVER mark ambiguous/failed Telegram delivery seen.
- NEVER use runtime `editMessageText` for admin callback, text-input completion, or manual-scan feedback.
- NEVER add schema/config/dependency/worker/client-detection machinery for this delta.
- NEVER modify unrelated production services, containers, networks, firewall, routing, reverse proxy, SELinux, or databases.

## Existing compatibility guardrails

- Initial baseline is not delivered; paused polls advance the baseline without backfill.
- Category parse/sanity failure advances neither seen IDs nor successful snapshot.
- Manual poll remains single-flight and does not bypass backoff.
- Somon 403/429 is not bypassed through restart pressure, proxies, or CAPTCHA circumvention.

## Verification routes

- FT-001/FT-002 deterministic tests and repository-native build gate.
- FT-003 ordered release receipts and production preflight/postflight.
- Detailed current transitions: [Runtime lifecycle](states/runtime-lifecycle.md).
