# Review request — keyword-monitoring decomposition

- TASK_ID: `TASK-MB-REVIEW-FEAT-PLAN`; STAGE_ID: `S-FEAT`.
- Role: fresh-context `ROLE: Reviewer`; skill: `.agents/skills/review-feat-plan/SKILL.md`, including `references/finding-adjudication.md`.
- Scope: accepted keyword delta only — REQ-010…014, EP-002, FT-005. Inspect FT-001…004 only for actual contradiction; preserve their approvals and task queues.
- Inputs: `AGENTS.md`, `.memory-bank/constitution.md`, `.memory-bank/analysis/product-brief.md`, `.memory-bank/analysis/index.md`, `.memory-bank/prd.md`, `.memory-bank/product.md`, `.memory-bank/requirements.md`, `.memory-bank/epics/EP-002-keyword-monitoring.md`, `.memory-bank/features/FT-005-keyword-monitoring.md`, `.memory-bank/spec-index.md`, `.memory-bank/spec-backbone.md`; applicable keyword sections in registered architecture, invariants, boundary, lifecycle and source-evidence docs.
- Target: PRD → REQ → EP → FT coherence, traceability, acceptance closure and bounded feature value. No JSON task design or implementation-detail review.
- Governing context: accepted bounded redesign keeps Planning Revision 1 and Foundation `not_required`. Native category/city with `q` is source evidence `needed_before_tasks`, not a pending product choice. Do not invent source parameters or replace scoped search with global first-page/local city filtering.
- Semantic focuses: traceability/Constitution; acceptance/failures/NFRs and one bounded feature-boundary probe. Co-reviewers use GPT-6.1 Sol under the operator's explicit model instruction.
- Required verdict: `VERDICT: APPROVE|REJECT`.
- Output: `.tasks/TASK-MB-REVIEW-FEAT-PLAN/TASK-MB-REVIEW-FEAT-PLAN-S-FEAT-final-report-docs-01.md`. Reviewer edits only these skill-owned operational outputs; no source/plan/feature repair, commits, deployment or network probes.
- Previous EP-001 request/report are preserved unchanged in `archive/EP-001/`.
