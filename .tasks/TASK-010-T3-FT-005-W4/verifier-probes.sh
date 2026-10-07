#!/bin/sh
set -eu
docker run --rm --network none -v "$PWD:/src:ro" -i somon-price-hotfix-builder:latest sh -es <<'SH'
 work=$(mktemp -d)
 trap 'rm -rf "$work"' EXIT
 cp -a /src/internal /src/cmd /src/testdata /src/go.mod "$work/"
 cp /src/.tasks/TASK-010-T3-FT-005-W4/verifier_ui_outcome_test.go "$work/internal/telegram/verifier_delete_test.go"
 cp /src/.tasks/TASK-010-T3-FT-005-W4/verifier_store_outcome_test.go "$work/internal/store/verifier_delete_test.go"
 cp /src/.tasks/TASK-010-T3-FT-005-W4/verifier_app_outcome_test.go "$work/internal/app/verifier_delete_test.go"
 cd "$work"
 CGO_ENABLED=1 go test -race -count=1 -timeout=45s -v ./internal/app ./internal/store ./internal/telegram -run TestVerifier
SH
