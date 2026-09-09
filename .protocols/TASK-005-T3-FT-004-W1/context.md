---
description: Template for .protocols/TASK-NNN-TN-FT-NNN-WN/context.md (clean-session context set).
status: active
---
# Context — <TASK_ID>

## Purpose
This file captures the **minimal reproducible context** so a fresh session can resume work safely.

## Execution Attempt
- attempt:
- started:

## Inputs (what drives this task)
- Task record: `.memory-bank/tasks/<TASK_ID>.task.json`
- Task index: `.memory-bank/tasks/index.json`
- Specs: (FT/EP/REQ docs you opened)
- Acceptance criteria source: (FT or verification section)

## Richer inputs (optional)
- Source Artifacts: ...
- Normative Inputs: ...
- Constraints / Invariants: ...
- Verification Targets: ...

## Fallback basis (if richer inputs were absent)
- Classic feature doc: ...
- Requirements / RTM: ...
- Duo docs: ...

## Loaded context set (what was read)
Keep this list short (2–8 items). Prefer SSOT pointers.
- `AGENTS.md`
- `.memory-bank/index.md`
- ...

## Decisions / assumptions
- Decision: ...
- Assumption (needs verification): ...

## Commands run / environment notes
- `...` → OK/FAIL (logs/evidence → `.tasks/<TASK_ID>/...`)

## Open questions / blockers
- ...

## Next session
- Start by reading: `context.md`, `plan.md`, `progress.md`
- Next action (one concrete step): ...

## Execution Attempt
- attempt: 1
- started: 2026-09-09T08:37:24.194789+00:00

## Manual hotfix authority
Operator requested correction and careful deployment. Standalone maintenance of installed baseline, separate from unfinished FT-001..003 queue. No history replay or new product behavior. Production preflight identifies one healthy documented systemd unit and a clean checkout at 7f9c5f50d659.

## Execution Attempt
- attempt: 2
- started: 2026-09-09T09:01:58.960624+00:00
- correction basis: existing detail_unknown_seller fixture loses its previously parsed labeled price after narrowing fallback; restore only explicit Цена: price line and test price. Initial attempt evidence is supporting-only under attempt-1/.
