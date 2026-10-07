#!/bin/sh
set -eu
review_work=/tmp/task008-red-verifier
mkdir -p "$review_work"
cp -a /src/internal /src/cmd /src/go.mod "$review_work/"
cp /src/.tasks/TASK-008-T3-FT-005-W2/red_verifier_atomic_management_test.go "$review_work/internal/app/red_verifier_atomic_management_test.go"
cd "$review_work"
gofmt -w internal/app/red_verifier_atomic_management_test.go
CGO_ENABLED=1 go test -count=1 -v ./internal/app -run '^TestRedVerifierManagementAtomicMutations$'
