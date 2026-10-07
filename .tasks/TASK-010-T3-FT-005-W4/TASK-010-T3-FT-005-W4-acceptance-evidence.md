# TASK-010-T3-FT-005-W4 — execution acceptance evidence

- attempt: 1
- applicability: applicable
- accepted claim locator(s): .memory-bank/features/FT-005-keyword-monitoring.md#FT-005-AC-002
- Governing REQ: REQ-010/013/014; only FT-005-AC-002 owned.

RED observation and evidence: baseline-probes.log compiles and runs real callbacks: list/settings ks:delete report unknown action, selected_exists=true/history_rows=2. Real authorized delete callback completes at paused detail barrier; releasing old detail starts one Telegram request and durable delivered history. Exit1 for behavioral assertions. initial-ui-and-harness-timeout.log has an invalid cleanup failure and is not app RED.

GREEN observation and evidence: unchanged probes in green-probes.log exit0: both actual actions remove search/history, unaffected other search/history/rental snapshot, reopen absence, next ID3 > deleted ID2; detail-start → delete-complete → detail-release → zero send starts/history. Existing auth and stale-edit/reject tests initially/finally GREEN.

Expanded necessary harm proof (harm-probes.log, targeted race exit0): pending edit from another admin returns missing after deletion; unauthorized/wrong chat acknowledges without changing snapshots; cascade abort at second child row restores monitor and both histories, then successful deletion durable with exact bidirectional SQL EXCEPT rental settings/seen timestamp/all state/offset equality and other search/history equality. Already-started send is received while lock held, deletion serializes after completion, one request allowed, then deleted history absent; both supported conditional writes and stale candidate cause no recreation/send after deletion and reopen. Next scheduler cycle sends none for deleted monitor.

T3 disposable surface: every scenario t.TempDir search.db, fake admin1/2 and target groups, local httptest source/Telegram; safe rerun from fresh seed, cancel/join goroutines, close DB/server, --network none container removed. No real DB/Telegram/production/Git writes.

Commands: probes.sh (baseline exit1, GREEN exit0); Docker targeted `CGO_ENABLED=1 go test -race -count=1 -timeout=45s -v ./internal/app ./internal/store ./internal/telegram -run "TestKeywordDelete|TestKeywordPollingStaleRejectionAndStaleDelivery"` exit0. Required Docker focused packages and scripts/build.sh exit0: commands-results.json, logs and gate-input-state snapshots. Native gate covers formatting/all tests/vet/CGO build/linkage/version/checksum.

Source attribution: initial/final-source-snapshot.json and change-surface.json. HEAD unchanged (7db6f00f64bfb698b6a9bdec138e69d870892811); existing dirty prerequisite source preserved. Source delta only internal/app/keyword_delete_test.go, internal/app/keyword_search.go, internal/store/keyword_delete_test.go, internal/store/keyword_search_cgo.go, internal/store/keyword_search_nocgo.go, internal/telegram/keyword_delete_integration_test.go, internal/telegram/keyword_search.go. Durable coverage navigation updated. Dependency claims not adopted. No reuse candidate/independent verdict. Task in_progress; fresh /verify and T3 /red-verify required before parent /root closure/sync.
