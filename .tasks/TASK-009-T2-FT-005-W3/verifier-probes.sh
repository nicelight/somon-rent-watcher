#!/bin/sh
set -eu
TASK_DIR=.tasks/TASK-009-T2-FT-005-W3
case "${1:-}" in
 gates) CHECK='CGO_ENABLED=1 go test -count=1 ./internal/app ./internal/store ./internal/telegram; CGO_ENABLED=1 go build -o /tmp/verifier-somonwatch ./cmd/somonwatch' ;;
 probes) CHECK='CGO_ENABLED=1 go test -race -count=1 -timeout=60s -v ./internal/app -run TestVerifierSearch' ;;
 *) echo 'usage: sh .tasks/TASK-009-T2-FT-005-W3/verifier-probes.sh gates|probes' >&2; exit 2 ;;
esac
# Original source mounted read-only; all Go writes and injected probes are disposable.
docker run --rm --network none -v "$PWD:/src:ro" -v "$PWD/$TASK_DIR:/evidence:ro" somon-price-hotfix-builder:latest sh -c '
 set -eu
 work=$(mktemp -d)
 trap '\''rm -rf "$work"'\'' EXIT
 cp -a /src/internal /src/cmd /src/testdata /src/go.mod "$work/"
 cd "$work"
 cp /evidence/verifier_polling_outcome_test.go internal/app/verifier_polling_outcome_test.go
 gofmt -w internal/app/verifier_polling_outcome_test.go
 sh -ec "$1"
' sh "$CHECK"
