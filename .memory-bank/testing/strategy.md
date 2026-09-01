---
description: Minimal framework baseline policy for risk-based testing and verification evidence.
status: active
---
# Testing Strategy

## Risk-based checks
- Choose checks from concrete product and regression risks in the PRD,
  Constitution, requirements, features, subject specs, and actual project shape.
- Use the cheapest check that reliably proves the required behavior.
- Add a broader or more expensive test level only when a narrower check cannot
  prove the outcome.
- Do not create tests merely to fill unit, integration, or e2e categories.

## Integrity
- Do not weaken assertions, disable failing checks, or replace meaningful
  verification with decorative coverage to obtain a green result.
- Treat a failing applicable check as evidence to investigate or resolve within
  explicit scope.

## Alembic migration ownership
- For an accepted linear Alembic graph, one project-level database contract
  owns dynamic checks for a single head, no branches, and intact ancestry.
  A feature migration test owns only its revision, direct `down_revision`,
  upgrade/downgrade, schema transition, and data preservation.
- A new revision reruns the project-level contract; it does not require literal
  current-head consumers or historical feature-test updates. Resolve current
  head at execution preflight when needed and keep its ID out of durable testing
  specs, feature plans, and task cards.

## Evidence and ownership
- T2/T3 tasks require an executable verification path. T0/T1 may use compact
  evidence or a documented no-runnable-check route when no meaningful check exists.
- Store commands, results, logs, screenshots, and verdicts in the task-selected
  `.protocols/<TASK_ID>/` and `.tasks/<TASK_ID>/` paths, not in this policy.
- Keep product quality targets in requirements/features and simple verification
  methods in feature AC/task records. Use a subject spec only for a non-trivial
  reproducible measurement method or expert rubric; it never supplies a missing
  product target. Keep executable gates in task records.

## Current project profile

Brownfield evidence for this repository is routed through [.memory-bank/testing/current-coverage.md](current-coverage.md). It records the existing Go/CGO gate, package tests, fixtures, live `doctor` route and checks that were not reproducible in the mapping environment. This descriptive profile does not add a new target quality policy.
