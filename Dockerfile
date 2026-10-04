# ---------------------------------------------------------------------------
# xyhelper-nav — multi-stage image build.
#
#   docker build -t xyhelper/nav:latest .
#   docker run -d -p 8000:8000 -v $(pwd)/data:/app/resource/data xyhelper/nav:latest
# ---------------------------------------------------------------------------

# ------------------------------- build stage -------------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cache the module download layer separately from the sources.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/nav .

# ------------------------------ runtime stage ------------------------------
FROM alpine:3.20

# ca-certificates is needed to fetch remote site icons over HTTPS,
# tzdata keeps log timestamps in the expected timezone.
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 -h /app nav

WORKDIR /app

COPY --from=builder /out/nav /app/nav
COPY manifest/config /app/manifest/config
COPY resource/public /app/resource/public

# The SQLite database and the session files live here; mount a volume over it
# to keep the data across container upgrades.
RUN mkdir -p /app/resource/data \
    && chown -R nav:nav /app

USER nav

EXPOSE 8000

CMD ["/app/nav"]
