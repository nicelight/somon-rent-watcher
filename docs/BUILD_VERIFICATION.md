# Build verification

Package version: `1.0.0`

Verification performed in the packaging environment:

```text
CGO_ENABLED=1 go test -count=1 ./...        PASS
CGO_ENABLED=1 go vet ./...                  PASS
CGO_ENABLED=1 go test -race -count=1 ./...  PASS
bash -n scripts/*.sh                        PASS
./scripts/build.sh                          PASS
```

The build produced an x86-64 dynamically linked Linux binary using `libsqlite3`. That local binary is deliberately excluded from the delivery archive because it was built against a different Linux distribution/glibc. The runbook builds the final executable on the target AlmaLinux 9 host and checks it with `ldd` before installation.

Coverage includes:

- first-run baseline without group notifications;
- paused monitoring consuming IDs without later backfill;
- one-time delivery of a new matching ID;
- price/room/floor/author/type/negative-word filters;
- discounted current price;
- visible floor overriding a stale URL slug;
- VIP/TOP classification;
- exclusion of advertisements after the “other cities” boundary;
- seller active-ad count and unknown seller behavior;
- RSC fallback;
- blocked-page detection;
- SQLite persistence;
- Telegram photo fallback without immediate duplicate after an ambiguous failure.
