---
description: Memory Bank Bible — правила, инварианты и стандарты документации.
status: active
---
# Memory Bank Bible (MBB)

## Constitution precedence
- [.memory-bank/constitution.md](../constitution.md) is the top governing policy for agent decisions.
- MBB, spec-index, spec-backbone, invariants, contracts, states, testing, and workflow docs refine the Constitution and MUST NOT contradict it.

## SSOT pyramid
- **Code**: WHAT/HOW — implementation truth.
- **Docstrings**: contracts + @docs pointers.
- **Memory Bank**: WHY/WHERE — boundaries, invariants, navigation.

## Hard rules
1. Every `.memory-bank/**/*.md` MUST have frontmatter with `description:`.
2. Use annotated links: `[.memory-bank/path](rel-path): короткое описание`.
3. Atomic docs: one cohesive concern per doc; split by boundary, change cadence,
   consumers, or reuse, not by file length alone.
4. Duo docs remain valid: `architecture/` (WHAT/WHY) + `guides/` (HOW), cross-link both ways for concepts that use the classic pair model.
5. C4 layering: L1 product → L2 epics → L3 features → L4 plans/tasks.
6. Docs First: update MB immediately after finishing a task.
7. Refactor MB every 5–10 updates (split, merge, archive).
8. Separate facts from interpretations: mark hypotheses explicitly ("предположительно", "требует проверки").
9. After merge/rebase conflicts: re-check MB consistency.
10. MB-SYNC after each wave/significant change (see `workflows/mb-sync.md`).
11. When present, `constitution.md`, `spec-backbone.md`, `spec-index.md`, `glossary.md`, `invariants.md`, `contracts/*`, `states/*`, `runbooks/*`, and `testing/*` act as an explicit normative layer and should be linked from relevant docs.
12. Features compose product behavior and exact canonical spec links. New design
    specs use subject-based paths without `FT-<NNN>` names or feature ownership.

## Forbidden
- Copy-paste implementation details / pseudocode
- Duplicating configs instead of linking to source
- Speculative claims without evidence from code/metrics/tests

## Allowed / encouraged
- Invariants (MUST/NEVER)
- Contracts at boundaries
- Decision rationale + pointers
- Runbooks and verification procedures
