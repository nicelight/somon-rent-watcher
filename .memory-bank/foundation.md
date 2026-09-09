---
description: Foundation Dev Path evidence and feature pressure map.
status: active
last_updated: 2026-09-04
---

# Foundation Dev Path

## Gate Anchors

- Foundation Required: false
- Foundation Requirement: REQ-000
- Foundation Pseudo-Feature: FT-000
- Foundation Gate Task: not_required

## Minimal Work Path

- Build command: `./scripts/build.sh` on a Go/CGO/SQLite-capable host; documented local equivalent `docker compose build somonwatch`.
- Start command: `somonwatch run` through the detected local/production runtime.
- Primary entrypoint: `cmd/somonwatch/main.go`.
- Smoke path: `somonwatch version` plus `somonwatch doctor` with the intended environment.
- Test command: `CGO_ENABLED=1 go test -count=1 ./...` as part of `scripts/build.sh`.
- Evidence: `.memory-bank/testing/current-coverage.md` records passing local Docker and AlmaLinux build/test/vet/CGO/linkage, live doctor and production systemd/isolation checks.

## Feature Pressure Map

| Feature | Pressure | Foundation Response | Probe | Status |
|---|---|---|---|---|
| FT-001 | Poll ordering, shared request cap and SQLite seen transitions | Existing application/store/test harness is sufficient; feature task owns behavior delta. | App integration tests plus native build. | baseline_sufficient |
| FT-002 | Telegram callback/message behavior and scheduler feedback | Existing Bot API test server and bot/application boundary are sufficient. | Telegram/app tests plus native build. | baseline_sufficient |
| FT-003 | Production-sensitive publication and runtime isolation | Existing runbook, preflight, target build/doctor and verified systemd route are sufficient; current runtime is re-detected before write. | Ordered release receipts and pre/post host comparison. | baseline_sufficient |

## Deferred Decisions

| Decision | Why deferred | Trigger to revisit |
|---|---|---|
| Alternative production runtime | No target migration is accepted; current runtime is discovered at deployment. | Operator accepts a runtime migration or preflight proves the documented runtime no longer exists. |

## Foundation Exit Criteria

- Existing minimal build/start/smoke/test path is evidenced and needs no separate enabling work.
- Compatibility probes exist for SQLite, Telegram, Somon and target CGO linkage.
- No P0/P1 design pressure blocks feature tasking.
- Feature development path is allowed after feature-local design/task review.
