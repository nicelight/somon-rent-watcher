# Independent verification evidence — TASK-006-T2-FT-005-W1

Reviewer-owned verification completed 2026-10-06T22:50:36+05:00. Owned result: FT-005-AC-003 (REQ-011/013 source subset). All observations use Go 1.21.13 in existing network-disabled Docker builder, local httptest, fresh temporary container copy; servers close with defer, container removed by --rm. No runtime DB, Telegram, live source, secrets or production changes.

## Governing basis and preflight

Read indexed TASK-006 card/index; full current-attempt context/plan/progress/handoff/verification and acceptance evidence; AGENTS, Reviewer role, Constitution/MBB/spec backbone/spec index/index; full tier policy; direct canonical architecture keyword design, boundary map Modules/Dependency Graph/Somon adapter/Shared Data/keyword contract/shapes, testing risk-based checks and source observations; mapped feature AC003, accepted REQ delta, PRD decisions, EP002 and implementation plan. Index/file/ID, tier T2, wave W1, FT005, string-array reqs/depends_on, required gate and verify array shapes consistent. No dependencies, hard write allow-list or tier escalation. Forbidden scope preserved.

Somon -> Shared Data is registered and retained: transport/parser owns catalog/HTTP/conversion; new model struct and Currency/City are passive. Client/parser/rental/price extraction production files unchanged; only new keyword source code and passive model fields. No I/O/business eligibility/history/price-bound policy in model, no new writer, graph edge, dependency, worker, scheduler or UI.

## Executor claim path (supporting only)

[Acceptance evidence](TASK-006-T2-FT-005-W1-acceptance-evidence.md), [historical compiling probe](baseline_probe_test.go), [baseline log](baseline-red.log), and protocol progress/handoff establish attempt 1 behavioral RED through existing FetchCategory/ParseCategory. Both unchanged representative fixtures yielded foreign ID 21000002 instead of primary [21000001]/[]. Initial request encoding and blocked/malformed/403 behavior were GREEN and existing transport/parser remained unchanged. Final keyword entrypoint GREEN uses identical ID expectations/fixtures. No missing-symbol/setup RED or receipt reuse. These artifacts are supporting; current verdict does not rely on executor GREEN alone.

## Reused execute evidence

None. Execute logs were treated only as supporting evidence. Cheap required gate repeated directly.

## Repeated checks

Command, exit 0; [output](verifier-focused-gate.log):

```sh
docker run --rm --network none -v "$PWD:/src:ro" -w /src somon-price-hotfix-builder:latest sh -c 'CGO_ENABLED=1 go test -count=1 ./internal/somon ./internal/htmlx'
```

Both packages pass. Existing rental DOM/RSC, hidden blocked modal, price/photo/address and seller fixtures remain regression evidence for additive source compatibility. Catalog paths/labels/keys, invalid scope, field currency evidence and final implementation tests pass. Executor focused vet/gofmt logs supporting; no extra gate categories introduced.

## New targeted probes

Independent source [verifier_outcome_probe_test.go](verifier_outcome_probe_test.go); [complete output](verifier-outcome-probe.log). Copied only into disposable container package, never production test tree. Command, exit 0:

```sh
docker run --rm --network none -v "$PWD:/src:ro" -w /src somon-price-hotfix-builder:latest sh -c 'cp -a /src /tmp/reviewer006 && cp /src/.tasks/TASK-006-T2-FT-005-W1/verifier_outcome_probe_test.go /tmp/reviewer006/internal/somon/reviewer006_outcome_test.go && cd /tmp/reviewer006 && CGO_ENABLED=1 go test -count=1 -v ./internal/somon -run TestReviewer006Outcome'
```

| AC003/spec observation | Independent decisive comparison |
|---|---|
| Native category/city and matching | All 16 requests captured before local transport rewrite. HTTPS somon.tj, hardcoded accepted path/category × country/dushanbe/vose/dangara, exact trimmed phrase `стол + & стул / ? # %`, only q+ordering=relevance, GET and configured user-agent preserved. |
| Primary-only nonempty/multiple/small | Independent nested-section fixture returns [26000001,26000002], ordered and deduped; foreign DOM/RSC 26000003 absent. Different title from q still returns, no local substring/apartment predicate. Unchanged small fixture returns exactly [21000001]. |
| Confirmed empty | Unchanged empty fixture and independent H1 zero + foreign recommendations return [] with nil error. |
| Unconfirmed/malformed/contradictory/RSC-only | Each returns error and no cards; unscoped RSC is safe failure, never fabricated primary success. HTTP 200 diagnostic raw body preserved. |
| Block/error transport | Blocked HTML including zero-looking heading remains typed blocked error; 403/429 remain typed HTTPError. Retry-After 90 preserved as 90s. |
| Source passive fields | Explicit 0 сомони remains Price=0/Currency=TJS; visible Восе/image preserved, missing second-card price/currency/city/image remain absent. Currency USD never becomes TJS; invalid keys/blank phrase rejected before HTTP. |
| Rental regression and owners | Required packages green; production rental client/parser unchanged, passive compatible model additions only; registered edge and source-owned catalog retained. |

## Finding adjudication and limits

Loaded and applied installed verify finding-adjudication pack and review-code agent contract. Mandatory fresh best-effort code co-review attempted with user-selected GPT-6.1-sol xhigh; tool rejected launch: agent thread limit reached. No co-review output relied on; /verify explicitly permits continuation after launch failure.

No material task-relevant implementation defect found. Fixtures are expressly authorized representative sanitized HTML with documented source-observation provenance, not live captures; tests prove local source contract behavior, not live DOM compatibility. No live probe/compatibility obligation belongs to this task. UI/storage/filter bounds/polling/history advancement and notification acceptance belong to later indexed cards and were not adopted. Known source limitations in handoff preserved.

## Next owner

Root manual lifecycle owner recorded .protocols/FT-005/plan.md may close T2 task after reading verification and append closure decision; verifier keeps in_progress. Feature completion still requires separate feature red-verify after all tasks. No commit or deployment.
