# Autonomous run plan

## Objective

Deliver the operator-approved Somon Rent Watcher delta end to end:

- replace Telegram inline-message edits with an append-only admin UI;
- when a poll has no fresh exact match, send up to three nearest-price fresh ads that satisfy every other filter and cost no more than 150% of the configured maximum;
- verify locally, publish to GitHub, synchronize the production checkout, then update only the verified Somon Rent Watcher runtime without disturbing other host workloads.

## Mode and ownership

- Entry: `/multiagentic`
- Base owner before product queue: `/autonomous`
- Product scheduler after strict-ready handoff: delegated `/autopilot`
- Execution: sequential
- Judge: one persistent read-only Judge for the entire run

## Current phase

- Product/Design preflight
- Brownfield baseline: complete
- Product Brief: proceed
- Pre-queue health: passed
- Next action: resolve Constitution governance decision through `/constitution`

## Canonical sequence

1. Resolve governing principles.
2. Create and review PRD and global SDD backbone.
3. Decide and, only if required, close Foundation work.
4. Plan and review each product feature in fresh contexts.
5. Pass lint and strict Memory Bank readiness.
6. Delegate the product queue to `/autopilot` for sequential execute, verify, semantic gates, closure, sync, and wave review.
7. Treat GitHub publication and production deployment as external writes requiring the already explicit operator authorization plus a current read-only production preflight.

## Safety boundaries

- No SQLite schema migration or second shown-history store unless an accepted upstream decision changes the target.
- No secret output, unrelated production changes, broad restart, Docker/network/firewall/Traefik mutation, or deletion of the production database.
- Production runtime shape must be proved immediately before deployment; current durable evidence describes a host `systemd` service, not a container.
