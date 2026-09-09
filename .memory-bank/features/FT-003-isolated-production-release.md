---
description: Publish, synchronize, and deploy the verified watcher without disturbing other production workloads.
status: draft
lifecycle: planned
spec_design_status: complete
spec_design_links:
  - ".memory-bank/architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable"
  - ".memory-bank/architecture/system-architecture.md#ad-004--deployment-is-runtime-detected-and-isolated"
  - ".memory-bank/architecture/system-architecture.md#deployment-boundary"
  - ".memory-bank/contracts/boundary-map.md#external-boundary-rules"
  - ".memory-bank/invariants.md#accepted-must"
  - ".memory-bank/invariants.md#accepted-never"
  - ".memory-bank/foundation.md#feature-pressure-map"
  - ".memory-bank/runbooks/almalinux-9-operations.md#accepted-ft-003-release-procedure"
  - ".memory-bank/runbooks/almalinux-9-operations.md#ft-003-release-proof"
  - ".memory-bank/testing/strategy.md#risk-based-checks"
---

# FT-003 — Isolated production release

## Value and use cases

The operator receives the verified behavior on production with preserved settings/seen history and evidence that more important workloads on the shared host were not changed.

## Source and constraints

- [.memory-bank/prd.md#req-005--ordered-and-isolated-release](../prd.md#req-005--ordered-and-isolated-release)
- [AlmaLinux operations baseline](../runbooks/almalinux-9-operations.md)
- Existing user authorization explicitly permits GitHub publication and a scoped production update, but runtime type must be re-established read-only immediately before writing.

## Edge and failure behavior

- A failed local gate, unhealthy host preflight, target build/doctor failure, or ambiguous runtime identity stops deployment.
- No database deletion/migration, broad restart, or unrelated container/service/network/security mutation is allowed.

## Acceptance Criteria

### FT-003-AC-001 — Ordered verified deployment

- REQ: REQ-005, REQ-008
- Observable criterion: evidence shows this strict order: local gates pass; commit is pushed to GitHub; read-only production preflight identifies a single healthy Somon Rent Watcher runtime; ambiguous runtime identity or unhealthy host stops the deployment; production git synchronizes that commit; target-compatible build/doctor passes; only that runtime is updated/restarted; post-deploy version, health and logs pass.
- Verification method: ordered command receipts with commit IDs, preflight identity/health, build/version output, redacted doctor result, scoped runtime status, and logs.

### FT-003-AC-002 — Preserve state and host isolation

- REQ: REQ-005, REQ-008
- Observable criterion: watcher settings/seen SQLite survives the update and read-only before/after checks show unrelated workloads, sockets, firewall/routing, and SELinux state unchanged.
- Verification method: scoped database counts/settings checksum-safe observation plus preflight/postflight comparison without secret disclosure.

## Acceptance closure

- Publication order, runtime ambiguity stop, target compatibility, state preservation, runtime health/logs, and unrelated-host isolation are covered by FT-003-AC-001 and FT-003-AC-002.

## Global backbone links

- [AD-001 and AD-004](../architecture/system-architecture.md#architecture-spine)
- [Deployment boundary](../architecture/system-architecture.md#deployment-boundary)
- [External production boundary](../contracts/boundary-map.md#external-boundary-rules)
- [Accepted production invariants](../invariants.md#accepted-never)
- [Foundation baseline](../foundation.md)
- [AlmaLinux operations route](../runbooks/almalinux-9-operations.md)

## SDD Design

Feature design is complete by reusing the accepted deployment and operational
owners. No feature-local technical spec, ADR, deployment abstraction,
alternative runtime, database lifecycle, rollback promise, or new
infrastructure is required.

- Architecture Spine: [AD-001](../architecture/system-architecture.md#ad-001--preserve-one-cohesive-deployable) keeps one deployable and [AD-004](../architecture/system-architecture.md#ad-004--deployment-is-runtime-detected-and-isolated) binds publication, preflight, update, and host isolation.
- Deployment and boundaries: the [deployment boundary](../architecture/system-architecture.md#deployment-boundary), [external production rule](../contracts/boundary-map.md#external-boundary-rules), and accepted [MUST](../invariants.md#accepted-must)/[NEVER](../invariants.md#accepted-never) constraints require exact order, runtime ambiguity stop, state preservation, and no unrelated production mutation.
- Operations: the [accepted FT-003 release procedure](../runbooks/almalinux-9-operations.md#accepted-ft-003-release-procedure) reuses the current systemd runbook while adding exact-commit GitHub/production synchronization and fresh read-only runtime identification. A non-systemd or ambiguous runtime stops instead of selecting a new deployment mechanism.
- Verification: [FT-003 release proof](../runbooks/almalinux-9-operations.md#ft-003-release-proof) defines the known initial state, safe rerun, observable receipts, SQLite/settings preservation, redaction, and stable pre/post host comparison under the [risk-based testing policy](../testing/strategy.md#risk-based-checks).

Global Backbone Planning Revision 1 is complete and Foundation is
`not_required`; its [FT-003 pressure probe](../foundation.md#feature-pressure-map)
is satisfied by the existing runbook, target build/doctor, and fresh runtime
preflight.
