# FT-002 finding clarification

## Findings

### Over-merged callback and manual-scan outcomes

- Semantic basis: FT-002-AC-001 and FT-002-AC-002 cover callback/navigation,
  mutation, authorization-boundary, and preserved text-input behavior in the
  Telegram administration flow. FT-002-AC-003 separately covers accepted
  manual-scan start/completion/error feedback while Polling Application keeps
  scheduler ownership. The current code locates these results in distinct
  production paths (`Bot.processCallback` and its menu/input helpers versus the
  `e:scan` branch and `Bot.CompleteManualPoll`), and each has a deterministic
  proof path that can pass without implementing the other.
- Finding disposition: confirmed. Shared Telegram ownership, transport method,
  files, or a desire to release the full feature together does not make the two
  independently completable implementation results execution-cohesive under
  `.memory-bank/workflows/execute-loop.md#execution-cohesive-task-boundary`.
- Option analysis: keeping one task would preserve the current artifacts but
  violate the accepted slicing rule. The only contract-valid minimal repair is
  two cohesive implementation tasks: one for prompt authorized callback
  acknowledgement plus fresh navigation/mutation output with claim-equivalent
  preservation proof for the already-working text-input completion, and one
  for separate append-only manual-scan start/completion/error feedback with
  scheduler semantics preserved. This changes planning only; it adds no
  product behavior, owner, dependency, or implementation mechanism.
- Likely affected artifacts: `.memory-bank/tasks/plans/IMPL-FT-002.md`,
  `.memory-bank/tasks/TASK-002-T3-FT-002-W1.task.json`,
  `.memory-bank/tasks/index.json`, and `.protocols/FT-002/plan.md` plus the
  planning decision record maintained by the tasking workflow. The tasking
  owner must preserve all FT-002 AC IDs and map every exact claim once without
  duplicating dependency proof.
- Remaining decision: none; accepted feature criteria and the governing task
  boundary permit only the split result.
- Design impact: none
- Behavior spec impact: none
- Immediate route: `/feature-to-tasks FT-002` rebuilds the implementation plan
  and queue, then a fresh `/review-tasks-plan FT-002` revalidates them.

### Native build gate writes outside the hard boundary

- Semantic basis: the required `./scripts/build.sh` gate runs `go build -o`,
  `chmod`, and checksum generation against `dist/somonwatch` and
  `dist/somonwatch.sha256`. The current non-empty
  `runtime_context.write_boundary` lists neither path, while
  `.memory-bank/workflows/tier-policy.md#hard-write-boundary` makes every such
  omitted project write a hard stop.
- Finding disposition: confirmed. The gate and current task card cannot both be
  executed as written.
- Option analysis: omitting/emptying `write_boundary` is contract-valid and
  would leave semantic scope, forbidden scope, and stop conditions in force,
  but it removes useful path-level enforcement. Adding `dist/` would be broader
  than necessary. The minimum repair is to retain the narrow source/test
  allow-list and add only the two deterministic generated artifact paths above
  to every replacement task that requires this gate. The tracked
  `dist/.gitkeep` means `mkdir -p dist` does not require authorizing creation of
  a broader directory in a clean checkout.
- Likely affected artifacts: each replacement task card that keeps
  `./scripts/build.sh`, its matching implementation-plan boundary description,
  and queue/index references created or reconciled by tasking. The gate itself,
  product source, and canonical testing policy do not change.
- Remaining decision: none; exact artifact authorization makes the accepted
  gate executable without a broad write scope.
- Design impact: none
- Behavior spec impact: none
- Immediate route: `/feature-to-tasks FT-002` applies the exact hard-boundary
  repair while rebuilding the split queue.

## Status

- Findings: complete
- Remaining decisions: none
- Feature semantics and accepted SDD design remain unchanged; the current
  `spec_design_status: complete`, links, REQ mappings, and FT-002 AC IDs remain
  valid.
- Durable repair owner: `/feature-to-tasks FT-002`.
