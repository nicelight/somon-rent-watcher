---
description: Claim-linked execution evidence for TASK-001-T2-FT-001-W1.
status: active
---
# Acceptance evidence — TASK-001-T2-FT-001-W1

## Attempt 1 pre-change basis

- Production source basis: Git `5b6a76e2cc4e59134f25d23897974f75b2747f5e`; expected production/test paths were clean before execution.
- Probe addition before production behavior change: `internal/app/price_fallback_test.go` only.
- Environment: disposable `golang:1.21-bookworm` container with `libsqlite3-dev`; repository-mounted tests use httptest Somon/Telegram and temporary SQLite only.
- Host toolchain attempt: `gofmt -w ... && CGO_ENABLED=1 go test ...` could not start because host `gofmt` is absent, matching the recorded local-development limitation. The probe was formatted and executed in the repository-native Docker toolchain.

## FT-001-AC-001 — alternative pre-change GREEN

Command:

`CGO_ENABLED=1 go test -count=1 ./internal/app -run 'TestProcessNewCards(ExactSuppressesFallback|WithoutPriceMaxDoesNotRunFallback|SeenAndFinalRejectionsStaySilent)$'`

Result: exit 0. Both successful and failed exact delivery captured only exact ID `1`, zero fallback delivery, and the failed exact ID remained unseen. The no-`PriceMax` scenario stayed on the ordinary exact path.

## FT-001-AC-002 — honest pre-change RED

Command (shared with AC-004):

`CGO_ENABLED=1 go test -count=1 ./internal/app -run 'TestProcessNewCards(SelectsClosestFallbacks|IncompleteExactSuppressesFallback|FallbackDeliveryFailureDoesNotSubstitute)$'`

Result: exit 1. `TestProcessNewCardsSelectsClosestFallbacks` captured `deliveries=[]`; required order is `[14 12 13]`, representing prices `[1001, 1200(first), 1200(second)]`. Baseline had no fallback delivery path.

## FT-001-AC-003 — alternative pre-change GREEN

The exit-0 command above captured no detail request for pre-seen ID `51`, one detail request for unknown-price ID `52`, no delivery, and final seen state for pre-seen/unknown/over-ceiling IDs `51/52/53`. Baseline already preserved this claim.

## FT-001-AC-004 — honest RED plus preserved behavior

The exit-1 command above showed:

- low-cap/transient exact scenario used one detail request and sent nothing, but prematurely marked fallback ID `32` seen instead of deferring it unseen;
- fallback delivery-failure scenario captured `deliveries=[]` and `DetailRequests=0`, proving the required fallback retry lifecycle was absent.

Existing exact failed-delivery unseen behavior passed under the AC-001 command. Existing typed 403/429 propagation remains covered by the repository app/Somon tests and will receive claim-equivalent post-change coverage in the focused suite.

## FT-001-AC-005 / REQ-006 — alternative pre-change GREEN

- Pre-change inventory: no task change in `go.mod`, `go.sum`, `internal/store`, `internal/config`, deploy/runtime, schema, persisted settings, worker, compatibility, or history ownership.
- Disposable pre-change source copy omitted only the new RED probe and ran `./scripts/build.sh`: exit 0; formatting, all tests, vet, CGO build, SQLite/libc linkage and version completed.
- Baseline checksum in the disposable container: `8fdb00e24432188442a4562ad68178303ad2c8b50a34ba40666125bc6c5a9190`.

## Post-change GREEN and final gates

### FT-001-AC-001

Focused suite exit 0. Exact ID `1` was the only attempted delivery with either successful or failed Telegram response; fallback was suppressed in both cases. Successful exact delivery became seen; failed exact delivery remained unseen.

### FT-001-AC-002

Focused suite exit 0. With `PriceMax=1000`, the captured fallback delivery IDs were `[14, 12, 13]`, corresponding to prices `[1001, 1200(first), 1200(second)]`. The price-1500 fourth candidate, price-1501 candidate, unknown-price candidate, and later equal-price ID `17` were not delivered and became seen after conclusive evaluation. Equal-price order remained stable across candidates originating from both exact-detail and known fallback partitions. Without `PriceMax`, the same ordinary ad followed the one exact path and was neither duplicated nor evaluated by a fallback phase.

### FT-001-AC-003

Focused suite exit 0. Poll-start seen ID `51` caused no detail/delivery request. Unknown-price ID `52` and over-ceiling ID `53` were silent and seen after final rejection. Complete no-result scenarios emitted no ad output.

### FT-001-AC-004

Focused suite exit 0. Evidence covers:

- one shared detail counter across phases;
- exact transient failure plus cap exhaustion suppressing fallback, with IDs `31/32` unseen;
- fallback-phase cap exhaustion after ID `35` suppressing all fallback delivery, with IDs `35/36` unseen;
- HTTP 403 and 429 returned as typed Somon block errors, with ID `37` unseen;
- failed first fallback delivery attempted only selected IDs `[41, 42, 43]`, did not substitute fourth ID `44`, left failed ID `41` unseen, and marked confirmed IDs `42/43` plus final fourth candidate `44` seen.

### FT-001-AC-005 / REQ-006

- Actual production change is confined to `internal/app/app.go`; Filtering remains unchanged because copying `Settings` in memory with `PriceMax=nil` reuses its existing pure predicates.
- `go.mod`, absent `go.sum`, `internal/store`, `internal/config`, deploy/runtime, schema, stored settings, scheduler/worker topology, compatibility and history surfaces are unchanged.
- Polling Application still exclusively owns phase order and every `MarkSeen` request; Somon, Telegram, Persistence, Filtering and Shared Data contracts are unchanged.

### Required gates

Host Go remains unavailable, so the exact task commands ran inside the documented repository-native `golang:1.21-bookworm` Docker toolchain with CGO and SQLite headers:

- `CGO_ENABLED=1 go test -count=1 ./internal/app ./internal/filter` — exit 0; `internal/app` and `internal/filter` passed.
- `./scripts/build.sh` — exit 0; formatting, all package tests, vet, CGO build, SQLite/libc linkage and version passed.
- Final binary checksum after the stable cross-partition tie assertion: `5f5b8752bd5958c82122e729e70d46a89bc812fc0f21de4dd62116fbbbbb070c`.
- `git diff --check` — exit 0.

The final Docker build is not offered as a reuse candidate: its package-install environment is external and the script writes generated `dist/` outputs. Fresh `/verify` should repeat the mapped checks.
