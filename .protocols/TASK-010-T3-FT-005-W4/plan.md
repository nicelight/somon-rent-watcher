---
description: Bounded execution plan and owner handoff for TASK-010-T3-FT-005-W4.
status: active
---
# Plan — TASK-010-T3-FT-005-W4

## Goal
FT-005-AC-002: deleting from list/settings removes exactly one monitor/history; stale work cannot send/recreate it.

## Non-goals
No adoption of dependency ACs, production/runtime/real data, commits, new dependencies/schema/workers, rental changes or authorization-policy changes.

## Inputs / source specs
Indexed task and .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-002; REQ-010/013/014; direct inputs resolved in context.md. Planning Revision 1 APPROVE; TASK008/009 done.

## Constraints / invariants (MUST / NEVER)
Keep existing owners: transport authorization/UI → app mutation/send coordination → store atomic SQLite transaction. Preserve rental/other searches, missing-result, revision condition, nonreused IDs and append-only callback output.

## Scope
Both delete buttons use ks:delete:<id> and one application operation. Existing ON DELETE CASCADE handles atomic history removal; existing keywordMu serializes send and delete. Tests reuse local temporary harnesses.

## Preflight-confirmed change surface
Advisory app/Telegram/store keyword sources and tests retained. New descriptive keyword_delete_test.go and keyword_delete_integration_test.go use nearest existing package/harness ownership. No source rename. Durable coverage pointer in testing/current-coverage.md satisfies WHY/WHERE and navigation; existing router already links it. Hard write_boundary not set; forbidden_scope untouched.

## Applicable quality gates
Both required Docker focused package test and native scripts/build.sh gates passed, exact command/result/input hashes in commands-results.json. Targeted race harm checks cover destructive outcome only.

## Claim-linked RED / GREEN (T2/T3)
- applicability: applicable
- accepted claim locator(s): .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-002
- planned test/probe and environment: real callback actions, detail/send barriers, cascade-abort trigger; t.TempDir search.db + local httptest in network-none Docker.
- observable RED: both delete callbacks leave selected search/history, detail-release permits stale send.
- corresponding GREEN: selected absent, other/rental preserved, stale callbacks/input missing, deletion durable/ID never reused; delete before send stops request; started request allowed then history removed; interrupted cascade rolls back both; stale-edit/reject stays protected.
- T3 isolation, safe rerun, cleanup, and permission boundary: identical fresh seed per test; cancel/wait goroutines, close DB/server, container removed; no external production state.

## MB-SYNC handoff / owner
Explicit standalone owner /root recorded TASK010 workflow in .protocols/FT-005/plan.md. Parent owns final lifecycle/feature-plan/RTM/changelog sync after independent functional + per-task semantic gates. Feature-level semantic gate remains due. /exe leaves task in_progress and does not invoke verification/sync.

## Definition of done
Required gates + independent /verify PASS + per-task /red-verify semantic-pass; owner closure. Executor implementation/evidence are complete.
