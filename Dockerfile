# syntax=docker/dockerfile:1

FROM golang:1.21-bookworm AS builder

RUN apt-get update \
    && apt-get install -y --no-install-recommends file libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod VERSION ./
COPY cmd ./cmd
COPY internal ./internal
COPY testdata ./testdata
COPY scripts/build.sh ./scripts/build.sh

RUN ./scripts/build.sh

FROM debian:bookworm-slim AS runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 somonwatch \
    && useradd --system --uid 10001 --gid somonwatch \
        --home-dir /var/lib/somonwatch --shell /usr/sbin/nologin somonwatch \
    && install -d -o somonwatch -g somonwatch -m 0700 /var/lib/somonwatch

COPY --from=builder --chown=root:root /src/dist/somonwatch /usr/local/bin/somonwatch

USER somonwatch:somonwatch
WORKDIR /var/lib/somonwatch

STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/local/bin/somonwatch"]
CMD ["run"]
