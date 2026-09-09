# FT-003 finding clarification

## Findings

### AC-002 alternative proof is hidden behind a RED label

- Semantic basis: REQ-005 and REQ-008 require SQLite/settings preservation and
  unchanged unrelated host state. Deliberately violating either condition to
  obtain RED would falsify the accepted production-safety claim, so
  `.memory-bank/workflows/tier-policy.md#claim-linked-red--green-for-t2t3`
  permits accepted alternative proof instead.
- Finding disposition: confirmed. The acceptance parser recognizes concrete
  `LABEL: value` fields independently. In the current `RED:
  RED_NOT_APPLICABLE because ...` wording it records concrete RED but cannot
  record `RED_NOT_APPLICABLE`, because that label has no colon and is nested in
  the RED value. AC-002 therefore does not satisfy the intended not-applicable
  proof branch. The fresh Reviewer returned this as its sole blocker; the
  absent cached report file is non-authoritative and does not change the
  finding.
- Option analysis: the only contract-valid minimal repair is a standalone
  `RED_NOT_APPLICABLE: <grounded reason>; alternative proof: <proof>` sequence,
  followed by the existing `GREEN:`, decisive comparison, and artifact. An
  actual RED is not viable because it would require an expressly forbidden
  destructive production action. This is a planning/proof-contract correction
  only; it changes no product target, feature behavior, design, tier, task
  identity, scope, dependency, boundary, or verification outcome.
- Likely affected artifact: only the AC-002 `evidence_required` string in
  `.memory-bank/tasks/TASK-004-T3-FT-003-W2.task.json` needs semantic repair.
  All other task fields and existing FT-003 product/design artifacts remain
  valid.
- Remaining decision: none; accepted product authority and the parser contract
  permit only the non-destructive alternative-proof route.
- Design impact: none
- Behavior spec impact: none
- Immediate route: `/feature-to-tasks FT-003` changes only that AC-002 evidence
  contract, preserving all other task fields, then a fresh
  `/review-tasks-plan FT-003` revalidates the repaired card.

## Status

- Findings: complete
- Remaining decisions: none
- Feature semantics and accepted SDD design remain unchanged; the current
  `spec_design_status: complete`, links, REQ mappings, and FT-003 AC IDs remain
  valid.
- Durable repair owner: `/feature-to-tasks FT-003`.
