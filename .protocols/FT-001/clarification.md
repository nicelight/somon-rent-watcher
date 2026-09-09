# FT-001 finding clarification

## Findings

### Missing no-PriceMax proof ownership

- Semantic basis: REQ-003 and the feature edge behavior already require no fallback without `PriceMax`.
- Resolution: FT-001-AC-002 owns the complete fallback eligibility/selection outcome and now names the absent-`PriceMax` case explicitly.
- Operator decision: not required; accepted product authority permits only this outcome.
- Design impact: none
- Behavior spec impact: none

### Incomplete REQ-006 feature criterion

- Semantic basis: REQ-006 already forbids compatibility layers and duplicated owners in addition to schema/settings/dependencies/workers/history.
- Resolution: FT-001-AC-005 now preserves the full feature-applicable REQ-006 pass/fail surface.
- Operator decision: not required; this is missing decomposition wording, not a new target.
- Design impact: none
- Behavior spec impact: none

### AC-001/003/005 alternative proof labels are nested under RED

- Semantic basis: FT-001-AC-001, FT-001-AC-003, and FT-001-AC-005 preserve
  behavior or absence already present in the baseline. Manufacturing claim
  absence would respectively require adding fallback despite an exact match,
  breaking established seen/rejection/silence behavior, or introducing a
  mechanism forbidden by REQ-006. Under
  `.memory-bank/workflows/tier-policy.md#claim-linked-red--green-for-t2t3`,
  those claims therefore require a grounded alternative-proof contract rather
  than an artificial RED.
- Finding disposition: confirmed. The acceptance parser recognizes a concrete
  label only as `LABEL: value`. Each current `RED: RED_NOT_APPLICABLE because
  ...` string is consequently classified as ordinary RED plus GREEN, not as
  the intended not-applicable branch. The latest persistent Judge assessment
  redirects this exact repeated label pattern before scheduler handoff.
- Option analysis: accepted authority permits only the minimal label repair.
  In each of the three existing `evidence_required` strings, replace the prefix
  `RED: RED_NOT_APPLICABLE because` with standalone
  `RED_NOT_APPLICABLE: because`. Preserve the grounded reason, alternative
  proof, GREEN, decisive comparison, artifact, and every other task field
  verbatim. An actual RED is not viable because it would falsify preserved
  product behavior or add a forbidden mechanism.
- Likely affected artifact: only the AC-001, AC-003, and AC-005
  `evidence_required` strings in
  `.memory-bank/tasks/TASK-001-T2-FT-001-W1.task.json`. Planning Revision 1,
  task ID, tier, wave, status, scope, order, product semantics, implementation
  plan, and code remain unchanged.
- Remaining decision: none; product authority, baseline evidence, tier policy,
  parser behavior, and the Judge condition determine the repair uniquely.
- Design impact: none
- Behavior spec impact: none
- Immediate route: `/feature-to-tasks FT-001` applies only those three prefix
  substitutions, then fresh `/review-tasks-plan FT-001`, strict doctor, and the
  same Judge checkpoint revalidate the queue.

## Status

- Findings: complete
- Remaining decisions: none
- Feature design remains `complete`; canonical owners and links are unchanged.
- Durable repair owner: `/feature-to-tasks FT-001`.
