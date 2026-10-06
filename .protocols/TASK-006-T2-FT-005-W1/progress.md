# Progress — TASK-006-T2-FT-005-W1

## Current status
- state: done
- attempt: 1
- Task lifecycle: done by /root after independent PASS; attempt 1 evidence retained.

## Claim-linked RED / GREEN
- applicability: applicable
- owned claim: FT-005-AC-003 (REQ-011/013)
- Planned baseline existing ParseCategory/FetchCategory on identical small/empty fixtures; GREEN new source entrypoint same primary IDs/empty marker and request fields.

## Next step
Create fixtures and run compiling baseline probe before production edits.

## Initial claim observation before production writes
Docker historical compiling probe exit 1, baseline-red.log. small IDs [21000002] instead of [21000001]; empty IDs [21000002] instead of []. Existing rental DOM parser only accepts apartment titles; RSC recommendation becomes feed. This is actual claim RED, no syntax/setup failure. Initial GREEN preserved: captured encoded native request, blocked parser, malformed parser, HTTP 403. Probe uses fixtures unchanged for final GREEN.

## Execution complete, verification next
Final source behavior GREEN; required focused Docker gate exit 0; focused vet exit 0; gofmt exit 0. Logs and precise comparison: .tasks/TASK-006-T2-FT-005-W1/TASK-006-T2-FT-005-W1-acceptance-evidence.md.
Actual production/data changes: internal/somon/keyword_search.go, internal/somon/keyword_search_test.go, internal/model/ad.go, internal/model/keyword_search.go, testdata/keyword-search-primary-small.html, testdata/keyword-search-primary-empty.html. Additional Card fields required by accepted passive boundary shapes; source kept in new file, client.go/parser.go unmodified. Runtime protocol/evidence bookkeeping scoped to this task. Minimal MB WHY/WHERE/index/changelog updated preserving outer-workflow edits.
Hard write_boundary not set; forbidden_scope untouched. Somon→Shared Data accepted edge unchanged, rental public API compatible. No new dependency/state writer/runtime or material branch, tier remains T2.
GREEN probe switches from existing category entrypoint to new keyword entrypoint while preserving fixture ID sets/empty-marker expectation and native request fields; tests add direct negative/error coverage without weakening original claim. Baseline/verbose/final logs supporting-only; no reuse candidate offered.

## Next step
/verify TASK-006-T2-FT-005-W1 in fresh independent Reviewer context. Task remains in_progress; closure/sync owner /root. No unresolved blockers. Live source compatibility remains unproven by representative fixture evidence.

## Owner closure evidence mapping

- accepted claim locator(s): .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-003
- RED observation and evidence: attempt 1 compiling baseline returned foreign ID 21000002 instead of primary 21000001 / empty set; exit 1. .tasks/TASK-006-T2-FT-005-W1/baseline-red.log and TASK-006-T2-FT-005-W1-acceptance-evidence.md retain the original pre-change observation.
- GREEN observation and evidence: same fixtures yield only primary 21000001 / empty set; final focused Docker tests exit 0. .tasks/TASK-006-T2-FT-005-W1/final-gate.log and TASK-006-T2-FT-005-W1-acceptance-evidence.md retain the result; independent verifier PASS is in verification.md.

This maps existing evidence into the framework field names; no new probe or backfilled
observation. /root closure: done. Next FT-005 task: TASK-007-T2-FT-005-W1.
