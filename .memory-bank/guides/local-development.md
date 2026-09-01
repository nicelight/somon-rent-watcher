---
description: Current local build and verification guide for the Go/CGO repository.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
---

# Local development and build — current state

## Prerequisites

- Go 1.21 or newer.
- GCC and CGO enabled.
- SQLite development headers and system `libsqlite3`.
- Bash plus `gofmt`, `go test`, `go vet`, `file`, `ldd` and `sha256sum` as used by the build script.

The repository has no third-party Go modules and no `go.sum`; [go.mod](../../go.mod) contains only the module path and Go version.

## Canonical current gate

Run [scripts/build.sh](../../scripts/build.sh). It performs:

1. Go version and SQLite header preflight.
2. `gofmt -l cmd internal` and fails on any output.
3. `CGO_ENABLED=1 go test -count=1 ./...`.
4. `CGO_ENABLED=1 go vet ./...`.
5. `CGO_ENABLED=1 go build` with trimmed paths/linker metadata into `dist/somonwatch`.
6. Binary/linkage/version checks and `dist/somonwatch.sha256`.

The delivery archive keeps only `dist/.gitkeep`; the final binary is intentionally built on the target AlmaLinux host for glibc compatibility.

## Additional recorded check

[docs/BUILD_VERIFICATION.md](../../docs/BUILD_VERIFICATION.md) records a successful race run (`CGO_ENABLED=1 go test -race -count=1 ./...`) in the packaging environment. Race is recorded evidence, but it is not part of the current `scripts/build.sh` gate.

## Current environment limitation

During the 2026-09-01 brownfield mapping, `go version` returned `go: command not found`, so the Go gate could not be reproduced in this workspace. `sha256sum -c MANIFEST.sha256` passed for every file listed in the delivery manifest. Install the prerequisites above before relying on local build/test status.

## Related docs

- [.memory-bank/architecture/system-architecture.md](../architecture/system-architecture.md): current component and runtime map.
- [.memory-bank/testing/current-coverage.md](../testing/current-coverage.md): exact existing tests and unverified surfaces.
- [.memory-bank/runbooks/almalinux-9-operations.md](../runbooks/almalinux-9-operations.md): production sequence.
