#!/bin/sh
set -eu
docker run --rm --network none -v "$PWD:/src:ro" somon-price-hotfix-builder:latest sh -ec '
 work=$(mktemp -d)
 trap '\''rm -rf "$work"'\'' EXIT
 cp -a /src/internal /src/cmd /src/testdata /src/go.mod "$work/"
 cp /src/.tasks/TASK-010-T3-FT-005-W4/keyword_delete_ui_probe_test.go "$work/internal/telegram/keyword_delete_probe_test.go"
 cp /src/.tasks/TASK-010-T3-FT-005-W4/keyword_delete_poll_probe_test.go "$work/internal/app/keyword_delete_probe_test.go"
 cd "$work"
 CGO_ENABLED=1 go test -count=1 -timeout=30s -v ./internal/telegram ./internal/app -run "TestKeywordDelete|TestKeywordPollingStaleRejectionAndStaleDelivery"
'
