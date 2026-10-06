---
description: Template for .protocols/TASK-NNN-TN-FT-NNN-WN/plan.md (execution plan + MB-SYNC handoff).
status: active
---
# Plan — TASK-007-T2-FT-005-W1

## Goal
Strict optional inclusive keyword price bounds, independent of apartment criteria.

## Non-goals
Native phrase/parser, UI, persistence, polling, delivery, deletion, fallback and production.

## Inputs / source specs
- Task record: `.memory-bank/tasks/TASK-007-T2-FT-005-W1.task.json`
- Task index: `.memory-bank/tasks/index.json`
- Feature/Epic: ...
- REQ IDs: ...

## Richer execution inputs (optional)
- Source Artifacts: ...
- Normative Inputs: ...
- Verification Targets: ...

## Fallback basis
- If richer inputs are absent, record the classic basis used for execution:
  - feature doc
  - requirements / RTM
  - duo docs
  - related contracts / states / runbooks / testing docs (if needed)

## Constraints / invariants (MUST / NEVER)
- MUST: ...
- NEVER: ...

## Scope
### In scope

### Out of scope

## Proposed changes
### Touched areas (hypotheses OK)
- `path/to/file` — why

### Preflight-confirmed change surface
- Expected hints kept: internal/filter/keyword_search.go and keyword_search_test.go.
- Additional same-outcome files/areas and rationale: task protocol/evidence and minimal Memory Bank WHY/WHERE/changelog.
- Hard `write_boundary` present and satisfied: not set
- `forbidden_scope` / stop-condition check: clear

## Applicable quality gates
List only evidence-backed project-native checks required by the task record,
linked specs/PRD, or repository configuration.

- [ ] Focused Docker package tests: task card exact command — proves keyword matrix and preserves rental regressions.
- [ ] gofmt only the new production/test files; Go vet ./internal/filter — local static checks.
- No meaningful runnable check: `<not applicable | rationale>`

## Claim-linked RED / GREEN (T2/T3)
- applicability: applicable
- accepted claim locator(s): FT-005-AC-004
- planned test/probe and environment: table-driven in-memory values; baseline existing CardMatches/AdMatches, final pure keyword predicate; network-disabled Docker.
- observable RED: bounded unknown-price CardMatches accepts; commodity AdMatches rejects on housing criteria.
- corresponding GREEN: missing/negotiable/unknown/foreign currency reject under any bound; commodity accepted; 99/100/200/201 inclusive matrix and zero.
- accepted not-applicable reason and alternative proof: none; baseline already-correct inclusive subcases preserved.
- T3 isolation, safe rerun, cleanup, and permission boundary: not T3; in-memory values reconstructed per subtest, no I/O or DB.

## Fan-out plan (if needed)
- Delegated agent A: scope ...
- Delegated agent B: scope ...

## MB-SYNC handoff / owner
Scheduler or explicit standalone owner performs sync after verification/status
decision. `/exe` only records handoff notes.

An `explicit standalone owner` exists only when the user directly asked the
current top-level agent to close the task, or when the top-level
agent/orchestrator explicitly runs a manual workflow for one TASK and records
that it owns closure. Subagent prompts do not silently become closure
owners.

Checklist:
- [x] Owner identified: explicit standalone owner /root; `.protocols/FT-005/plan.md#task006-closure-and-task007-selection`.
- [ ] Explicit standalone owner basis recorded if manual closure is expected: user direct instruction | top-level manual workflow ownership | n/a
- [ ] `.memory-bank/` docs needing update (WHY/WHERE, no pseudocode): ...
- [x] `.memory-bank/index.md` router update needed: no (existing canonical boundary-map route).
- [x] RTM update in `.memory-bank/requirements.md` needed: no (feature partial).
- [x] Task registry/status update owner: /exe starts; /root closes after independent /verify.
- [x] Changelog update owner: executor for this unit; /root wave sync.

## Definition of done
- ...
