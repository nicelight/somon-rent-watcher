# Context — TASK-013-T3-FT-005-W6

## Purpose
Deploy ALL current source authorized by operator, preserving one existing watcher and state.

## Execution Attempt
- attempt: 1
- started: 2026-10-07T16:23:59.510687+00:00

## Inputs
Indexed TASK013; feature AC010/REQ008/014; runbook accepted-ft-005-release-procedure plus FT003 sequence/proof; runtime keyword persistence rules. Global Planning Revision1; fresh review pending before start. Dependencies TASK006…012 done. Read governing AGENTS/constitution/MB policies and tier rules; no architecture/code changes.

## Decisions / assumptions
Operator resolved rollout 2026-10-07: «все что есть свежего в коде - выкладывай». Includes rental +50% and all keyword code/fixes. Root GENERAL explicitly owns this manual TASK013 workflow, verification routing and final closure/sync. No old queue adoption.

## Environment
ssh igorprod; /root/somon-rent-watcher clean main old93cbe9; one healthy somonwatch.service /opt/somonwatch/somonwatch; /var/lib/somonwatch/somonwatch.db; private env path unchanged. Historical probes reused as methodology only. Fresh preflight due before write. Local Docker native builder.

## Open questions / blockers
none; awaiting fresh task-plan review.

## Next session
Read context/plan/progress/task. Next: fresh review gate then start and local native build.
