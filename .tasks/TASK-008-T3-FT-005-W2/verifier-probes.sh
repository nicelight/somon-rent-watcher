#!/bin/sh
set -eu
work=/tmp/task008-verifier
mkdir -p "$work"
cp -a /src/internal /src/cmd /src/go.mod "$work/"
cp /src/.tasks/TASK-008-T3-FT-005-W2/verifier_search_management_test.go "$work/internal/telegram/verifier_search_management_test.go"
cp /src/.tasks/TASK-008-T3-FT-005-W2/verifier_legacy_storage_test.go "$work/internal/store/verifier_legacy_storage_test.go"
cp /src/.tasks/TASK-008-T3-FT-005-W2/verifier_rental_no_repeat_test.go "$work/internal/app/verifier_rental_no_repeat_test.go"
cd "$work"
gofmt -w internal/telegram/verifier_search_management_test.go internal/store/verifier_legacy_storage_test.go internal/app/verifier_rental_no_repeat_test.go
CGO_ENABLED=1 go test -count=1 -v ./internal/telegram ./internal/store ./internal/app -run '^TestVerifier'
