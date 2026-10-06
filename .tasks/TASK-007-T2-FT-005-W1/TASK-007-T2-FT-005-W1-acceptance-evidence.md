# Acceptance evidence — TASK-007-T2-FT-005-W1

Claim: FT-005-AC-004 (REQ-011), attempt 1. Supporting execution evidence; independent /verify pending.

Baseline probe compiles against existing CardMatches/AdMatches. Ran only after task in_progress and before production changes. Matrix has absent/min/max/both/zero bounds, 99/100/200/201/zero/missing/negotiable/USD/unknown currency; commodity without housing/seller fields.

- RED: 19 failures, 36 initial passes. Unknown/negotiable admitted with bounds; foreign/unconfirmed 150 admitted for eligible numeric range; commodity rejected for absent rooms. Inclusive 100 and 200 already pass and remain preserved.
- GREEN: equivalent 54 price vectors pass with keyword predicate; 11 bounds validation cases and two commodity vectors pass. 100/200 included, 99/201 excluded for 100–200; bounded unknown/negotiable/foreign/unconfirmed fail; absent bounds never restrict price; zero bounds distinguish known zero from missing.
- Gate: exact required focused Docker tests exit 0, includes existing rental tests. Vet/formatting exit 0. Source/rental implementation untouched.

Commands: progress.md / Commands run and Claim-linked RED / GREEN. Logs:

| Artifact | SHA256 |
|---|---|
| baseline-red.log | e663dfe0c067ce09f7ab64003ab983d43208aa27218dca36c4e334eceef020ad |
| claim-green.log | 73b6b1a8aaa6e2f1f0134df2c1a8602f58f1abd9de6a668a915e3757b2aca5c5 |
| focused-gate.log | 015ff30226b538d532809f19f0bfc944a9a051479be8faddff496ee8b9b6ae6c |
| static-gate.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |

Baseline source retained as baseline_probe_test.go. To reproduce baseline safely, temporarily copy it into internal/filter/keyword_price_baseline_test.go, run baseline command, then remove it before current gate; this compares historical rental behavior to the keyword criterion and is expected to fail.

Nested file mount setup failure (exit 125) was not RED; same copied probe was used successfully. No fake break or new-symbol compile failure used.

No final PASS/closure claimed. /root owns lifecycle after fresh /verify.
