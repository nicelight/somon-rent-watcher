#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${VERSION:-$(tr -d '[:space:]' < VERSION)}"
COMMIT="${COMMIT:-$(git rev-parse --short=12 HEAD 2>/dev/null || printf 'archive')}"
BUILD_TIME="${BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

command -v go >/dev/null || { echo "ERROR: Go is not installed" >&2; exit 1; }
command -v gcc >/dev/null || { echo "ERROR: gcc is not installed (CGO is required)" >&2; exit 1; }
[[ -f /usr/include/sqlite3.h || -f /usr/local/include/sqlite3.h ]] || {
  echo "ERROR: sqlite3.h not found; install sqlite-devel" >&2
  exit 1
}

GO_VERSION="$(go env GOVERSION | sed 's/^go//')"
GO_MAJOR="${GO_VERSION%%.*}"
GO_MINOR="${GO_VERSION#*.}"; GO_MINOR="${GO_MINOR%%.*}"
if (( GO_MAJOR < 1 || (GO_MAJOR == 1 && GO_MINOR < 21) )); then
  echo "ERROR: Go >= 1.21 is required; found $(go version)" >&2
  exit 1
fi

echo "==> Formatting check"
UNFORMATTED="$(gofmt -l cmd internal)"
if [[ -n "$UNFORMATTED" ]]; then
  echo "ERROR: unformatted Go files:" >&2
  echo "$UNFORMATTED" >&2
  exit 1
fi

echo "==> Unit and integration tests"
CGO_ENABLED=1 go test -count=1 ./...

echo "==> go vet"
CGO_ENABLED=1 go vet ./...

echo "==> Building dist/somonwatch"
mkdir -p dist
CGO_ENABLED=1 go build \
  -buildvcs=false \
  -trimpath \
  -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
  -o dist/somonwatch ./cmd/somonwatch
chmod 0755 dist/somonwatch

file dist/somonwatch
ldd dist/somonwatch | grep -E 'sqlite|libc|not found' || true
./dist/somonwatch version
sha256sum dist/somonwatch | tee dist/somonwatch.sha256

echo "Build complete: $ROOT_DIR/dist/somonwatch"
