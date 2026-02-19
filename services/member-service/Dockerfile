# ============================================================
# Build context: repo root (tellawl/)
#
# Build:
#   docker build \
#     -f services/member-service/Dockerfile \
#     -t <registry>/member-service:latest \
#     .
#
# Push:
#   docker push <registry>/member-service:latest
# ============================================================

# ── Builder ─────────────────────────────────────────────────
FROM golang:1.25-bookworm AS builder

# CGO dependencies required by confluent-kafka-go (builds librdkafka from source)
RUN apt-get update && apt-get install -y --no-install-recommends \
  gcc \
  g++ \
  make \
  pkg-config \
  libssl-dev \
  libsasl2-dev \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy shared packages first to leverage layer caching
COPY packages/ packages/

# Copy service source
COPY services/member-service/ services/member-service/

WORKDIR /build/services/member-service

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux \
  go build -ldflags="-s -w" -o /bin/member-service ./cmd/api/

# ── Runtime ──────────────────────────────────────────────────
FROM debian:bookworm-slim AS runtime

# Runtime shared libraries needed by confluent-kafka-go and TLS
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  libssl3 \
  libsasl2-2 \
  && rm -rf /var/lib/apt/lists/*

RUN useradd -r -u 1001 -g root appuser

WORKDIR /app

COPY --from=builder /bin/member-service .
COPY --from=builder /build/services/member-service/db ./db

USER appuser

EXPOSE 8080

ENTRYPOINT ["./member-service"]
