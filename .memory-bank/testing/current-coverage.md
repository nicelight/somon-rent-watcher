---
description: Evidence-backed inventory of current automated, fixture, build and live verification paths.
status: active
baseline_kind: as-is
last_verified: 2026-09-01
---

# Current test and verification coverage

## Native gate

[scripts/build.sh](../../scripts/build.sh) is the project-native build gate: gofmt check, `go test ./...`, `go vet ./...`, CGO build, linkage/version inspection and checksum generation. Shell syntax and race tests are recorded separately in [docs/BUILD_VERIFICATION.md](../../docs/BUILD_VERIFICATION.md).

No repository CI workflow/configuration was found during file inventory.

## Automated coverage by package

| Surface | Existing proof path |
|---|---|
| Configuration | Defaults, admin-list parsing/legacy fallback, invalid admin IDs and invalid poll-range validation in `internal/config/config_test.go`. |
| HTML utility | Tolerant malformed HTML parsing, normalized text and URL resolution in `internal/htmlx/html_test.go`. |
| Somon parser | DOM cards, discount/current price, visible floor over stale slug, seller fields/unknown, other-city boundary, block detection, RSC fallback, age and recovery URLs in `internal/somon/parser_test.go`. |
| Filters/settings | Inclusive price, room/floor buckets, seller limits/unknown, normalized substrings, polling-range parsing/legacy defaults, input/JSON normalization and non-empty choices in `internal/filter/settings_test.go`. |
| SQLite | Settings, deduplicated seen IDs and state roundtrip in `internal/store/sqlite_test.go` (CGO build tag). |
| App pipeline | Continuity helper including the adaptive stale-time threshold, merge order, single-flight manual-poll rejection/backoff guard, first baseline plus one-time delivery, and paused consumption without later backfill in `internal/app/app_test.go` using `httptest` servers and temp SQLite. |
| Telegram | Menus/caption escaping, size and emoji-field line formatting; polling interval/manual scan busy/completion controls; token URL redaction, Retry-After, photo fallback semantics, multi-admin private/target-group authorization, isolated group text input and notification fan-out in `internal/telegram/*_test.go`. |

## Fixtures

- [testdata/category.html](../../testdata/category.html): category with promoted/ordinary cards and other-city boundary.
- [testdata/detail.html](../../testdata/detail.html): detail values including stale URL-floor conflict, discount and seller fields.
- [testdata/detail_unknown_seller.html](../../testdata/detail_unknown_seller.html): missing seller information.

These are repository fixtures, not proof of the current live Somon page.

## Live and operational verification

- `somonwatch doctor` checks SQLite open/schema, Telegram bot/chat/webhook state and a live category parse with minimum card/ordinary-card sanity. It does not create baseline or mark seen IDs.
- [scripts/preflight-almalinux.sh](../../scripts/preflight-almalinux.sh) captures read-only host/service/container/socket/firewall/SELinux/capacity evidence.
- The post-start procedure in [docs/RUNBOOK_ALMALINUX_9.md](../../docs/RUNBOOK_ALMALINUX_9.md) compares host state and checks service resource use/logs.

## Current local verification status

| Check | Result |
|---|---|
| `sha256sum -c MANIFEST.sha256` | Historical PASS before the current working-tree changes; manifest was not regenerated because no release/archive was requested. |
| Docker native package/test/vet/CGO/linkage gate | PASS on 2026-09-01 through `docker compose build somonwatch`; host Go is not required for this route. |
| Local live Somon/Telegram doctor | PASS on 2026-09-01 inside the permanent Compose container: SQLite, bot, target group and 60-card live parse were OK. |
| Target host preflight/post-start | NOT RUN: production host was not accessed. |

Historical native evidence remains in [docs/BUILD_VERIFICATION.md](../../docs/BUILD_VERIFICATION.md); the Docker rows above are fresh local reruns from the current wave.

## Needs verification

- Target-host build/linkage and systemd behavior.
- Current compliance posture based on time-sensitive external rules.
