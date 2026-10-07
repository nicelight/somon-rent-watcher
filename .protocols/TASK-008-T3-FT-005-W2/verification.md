---
description: Independent functional verification of durable keyword search management and rental preservation.
status: active
---
# Verification — TASK-008-T3-FT-005-W2

VERDICT: PASS

## Verification basis and scope

Fresh delegated ROLE: Reviewer, `/verify` only; task remains `in_progress`.
Indexed ID/path/tier/wave/feature, string-array fields, gate shapes, protocol and
both dependency statuses checked directly. TASK006/007 are `done`; their source/price
outcomes remain prerequisites. Exact owned proof: FT-005-AC-001 and FT-005-AC-007,
REQ-010/014/002. Polling/delivery/history/deletion of searches belong to later tasks.

Read AGENTS, Reviewer role, Constitution/MBB/navigation/backbone/index, verify skill
and finding-adjudication, tier-policy obligations/ownership/claim-linked proof,
indexed card and full context/plan/progress/handoff/prior verification. Direct
canonical basis: architecture keyword design; boundary-map Modules, Dependency
Graph, keyword proposal/shapes, Telegram application, Persistence and Shared-data
contracts; runtime keyword state/persistence rules; testing strategy risk-based
checks. Feature, accepted PRD/requirements/epic and manual ownership protocol read.
Planning Revision remains 1; no task-scope interpretation or product decision added.

## Executor claim path

Supporting evidence: [acceptance evidence](../../.tasks/TASK-008-T3-FT-005-W2/TASK-008-T3-FT-005-W2-acceptance-evidence.md),
[progress](progress.md), archived baseline probe and `baseline.log`.
AC001 baseline is a compiling behavioral RED: authorized `ks:new` returns unknown
button and no durable search exists. Earlier pointer compilation error is explicitly
setup-only. Current GREEN uses the same actual `Bot.processUpdate` entrypoint with
real App and temporary SQLite, now exercising complete creation/configuration/enable.
AC007 has the accepted RED_NOT_APPLICABLE reason: preservation already held;
artificial corruption is unnecessary. Initial exact SQLite/rental/auth GREEN is
retained; post-change proof compares protected rows and rental payload after writes/reopen.
No RED was fabricated/backfilled. Executor logs remain supporting observations;
none supplies this review's independent functional proof.

## Reused execute evidence

None. No receipt reuse was requested or accepted. Executor all-package tests,
vet/build/noCGO/gofmt and focused logs were read as supplementary execution evidence.

## Repeated checks

Required focused gate repeated with the repository mounted read-only; same image,
packages, CGO and network isolation as the indexed gate:

```sh
docker run --rm --network none -v "$PWD:/src:ro" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/telegram ./internal/app ./internal/store'
```

Exit 0, [verifier-focused-gate.log](../../.tasks/TASK-008-T3-FT-005-W2/verifier-focused-gate.log).
Direct rerun is cheap and avoids self-attested receipt reuse. No production build or
external request needed for the task outcome.

## New targeted probes

Verifier-authored executable tests remain under `.tasks` and are copied into a
container-only `/tmp` checkout by [verifier-probes.sh](../../.tasks/TASK-008-T3-FT-005-W2/verifier-probes.sh).
The repository is read-only. Safe reproduction:

```sh
docker run --rm --network none -v "$PWD:/src:ro" somon-price-hotfix-builder:latest sh /src/.tasks/TASK-008-T3-FT-005-W2/verifier-probes.sh
```

Exit 0, [verifier-probes.log](../../.tasks/TASK-008-T3-FT-005-W2/verifier-probes.log).
All three new tests passed. Each run allocates new `t.TempDir()/search.db`, local
httptest servers, fake token/IDs; DB/server close and test/container cleanup remove
runtime state. No default/production DB, secrets, real Telegram or external network.

| Owned result | Fresh independent observation |
|---|---|
| AC001 creation/configuration/summary/enable | `TestVerifierSearchManagementAndHarm`: private admin and second admin in target group create two distinct IDs, trimmed/escaped phrase, All/country/disabled/revision1 initial summary, all readable catalog buttons, independent city/category, category preserves city, zero/two/one optional bounds, summary before explicit enable, addressed edits and disable. Final rows ID1 revision9 disabled; ID2 revision8 enabled, independently persisted through real DB close/open/new App/Bot. |
| AC001 invalid input does not persist | Same test compares whole search payloads including revision/enabled after empty phrase, unknown city/category, negative/inverted/fractional/noninteger/overflow prices; unchanged rows. |
| AC001 authorization and all pending dimensions | Nonadmin and wrong group callbacks are ack-only and rejected text creates no rows. Two admins in one group edit different IDs; same admin in private/group has independent pending actions; switching selected ID replaces pending action; other-admin revision update rejects stale text. Navigation/cancel and replacing rental pending cannot leak into rental settings. Exact rows compared after these branches. |
| AC001 stale/missing/addressing/append-only | Missing view/enable/phrase/catalog targets return fresh list and no upsert. Repeated enable is idempotent; disable increments only selected revision. Every new-route callback acknowledges before output; all new-route requests are ack/sendMessage, never edits. |
| AC007 exact legacy data/additive reopen | `TestVerifierLegacyStorageExactAndStableIdentity`: seeded legacy-schema fixture reopened by additive initializer, two search writes, conditional update, stale/missing rejection and restart; bidirectional SQL EXCEPT checks every settings/seen_ads/state column, including first_seen_at and unknown state/settings data. Stable AUTOINCREMENT high-water identity also survives fixture-only removal of largest ID/reopen. This does not claim deletion functionality. |
| AC007 observable rental behavior | `TestVerifierRentalSeenNoRepeatAcrossManagementAndRestart`: real rental `pollOnce` before/after two enabled search writes and SQLite/App restart; same already-seen category IDs produce zero details and zero Telegram sends in all three polls. Settings, seen IDs/count and offset retained. Separate management test captures equal rental SendAd method/form before/after writes/reopen. |

An initial verifier helper accidentally demanded ack-before-output from the unchanged
rental `i:price` route, which is explicitly outside TASK008's new-route requirement.
Only this assertion's scope was corrected to `ks:`; the initial observation is kept in
`verifier-probe-scope-correction.log`. The final test still exercises rental-to-keyword
pending replacement and protects its rows. No implementation assertion was weakened.

## Architecture and actual change surface

Two review focuses: authorization/pending/state isolation; permitted module ownership
and preservation. A fresh GPT-6.1-Sol/xhigh co-review attempt (operator model override)
failed with agent-thread limit; per supplied one-attempt instruction no retry. Best-effort
co-review absence is nonblocking; final judgment rests with this Reviewer.

Compared all 12 current source hashes with change-surface.json twice; all match.
All other task-start snapshotted source/config/go.mod inputs remain unchanged,
including prior filter work; tracked TASK007 baseline deletion is preserved.
[Verifier source state](../../.tasks/TASK-008-T3-FT-005-W2/verifier-source-state.json)
records the source/toolchain/isolation basis. HEAD diff is not the ownership basis.

Canonical registered edges preserved: App -> Filtering/Somon/Store/Shared Data,
Telegram -> Filtering/Somon/Shared Data and application via existing Backend plus
optional KeywordBackend, Store -> Shared Data. App owns validation, create-disabled,
expected revision and enable decisions (`internal/app/keyword_search.go`); store owns
all SQL/transactions/conditional writes (`internal/store/keyword_search_cgo.go`,
transactional IF NOT EXISTS initializer). Telegram owns auth/pending/parsing/catalog
rendering and transport, with no production SQL/Store import. Shared Data is passive.
The added App mutex stays with the business owner; the exported update bridge is
only `_test.go`. Existing composition, config, Go process, DB/client/group/allowlist,
scheduler and dependencies remain unchanged; no new worker/infrastructure/second
source of truth or forbidden write/bypass was found. Additive search storage is
within accepted keyword delta; the pre-existing rental-only schema statement does
not override that explicit delta.

## Findings and handoff

No material functional defect or canonical ownership violation observed. Every owned
claim and T3 harm-driving branch has fresh verifier evidence; required gate passed.
This is local/disposable-state proof, not production acceptance. No implementation,
canonical spec, lifecycle, dependency or previous task record changed by Reviewer.

T3 PASS requires separate fresh `/red-verify TASK-008-T3-FT-005-W2` before closure.
Explicit manual owner `/root` in `.protocols/FT-005/plan.md` owns that boundary and
later lifecycle/W2 sync. Status remains `in_progress`.
