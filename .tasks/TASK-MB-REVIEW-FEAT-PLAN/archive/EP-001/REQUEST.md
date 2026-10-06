# Review request — product feature plan

- TASK_ID: `TASK-MB-REVIEW-FEAT-PLAN`
- STAGE_ID: `S-FEAT`
- Role: fresh-context `ROLE: Reviewer`
- Skill: `.agents/skills/review-feat-plan/SKILL.md`
- Semantic pack: `.agents/skills/review-feat-plan/references/finding-adjudication.md`
- Inputs: `AGENTS.md`, `.memory-bank/constitution.md`, `.memory-bank/analysis/product-brief.md`, `.memory-bank/analysis/index.md`, `.memory-bank/prd.md`, `.memory-bank/product.md`, `.memory-bank/requirements.md`, `.memory-bank/epics/EP-001-reliable-watcher-upgrade.md`, `.memory-bank/features/FT-001-price-fallback.md`, `.memory-bank/features/FT-002-append-only-telegram-ui.md`, `.memory-bank/features/FT-003-isolated-production-release.md`, `.memory-bank/spec-index.md`, `.memory-bank/spec-backbone.md`.
- Review target: PRD → REQ → EP → FT decomposition readiness for `/spec-design`.
- Required verdict vocabulary: `VERDICT: APPROVE|REJECT`.
- Required output: `.tasks/TASK-MB-REVIEW-FEAT-PLAN/TASK-MB-REVIEW-FEAT-PLAN-S-FEAT-final-report-docs-01.md`.
- Prohibition: do not review JSON task design or implementation detail; task records do not yet exist.
- Reviewer is read-only and returns its report to the parent; the parent persists the required report verbatim in substance.
