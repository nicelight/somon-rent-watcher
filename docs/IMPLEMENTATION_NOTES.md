# Implementation notes

## Intentional safety additions to the agreed MVP

Two small safeguards were added without changing the product model:

1. **Monitoring starts paused.** The first run creates the baseline, but notifications remain disabled until the administrator configures `/filter` and presses `▶ Включить мониторинг`. While paused, new IDs are still consumed, so enabling does not produce a delayed flood.
2. **Detail requests are capped per poll.** `SOMON_MAX_DETAILS_PER_POLL=20` prevents a burst after a long outage or recovery. Unprocessed IDs remain unseen and are retried while they remain in the visible feed.

## Dependency choice

The application has no third-party Go modules:

- HTML is parsed by a small tolerant tree parser included in `internal/htmlx`;
- SQLite is accessed through the stable system `libsqlite3` C API using CGO;
- Telegram Bot API is called directly with `net/http`.

This keeps the build reproducible from the archive without downloading Go modules. The target host needs Go, GCC and `sqlite-devel` only while building; at runtime it needs glibc and `libsqlite3`.

## Parser strategy

The Somon parser intentionally does not depend on one CSS class name. It combines deterministic signals:

- `/adv/<ID>...` links;
- visible title/price/age text;
- card ancestry with one unique advertisement ID;
- VIP/TOP text and class markers;
- current detail-page fields;
- deterministic Next.js React Server Component data as a fallback/enrichment source.

The visible DOM defines the actual city feed and its order. Structured payloads may enrich a matching ID but do not inject unrelated recommendation IDs when visible cards are present.

## Mandatory live verification

The build environment used for packaging had no direct outbound network access, so the compiled parser was verified with fixtures and unit/integration tests rather than a raw live HTTP fetch from that environment. The current page structure and fields were separately researched, but Somon can change HTML at any time.

For that reason `somonwatch doctor` is mandatory on the AlmaLinux host before starting the service. It performs a real low-rate fetch, checks the card count and confirms that ordinary cards are recognized without modifying the baseline.
