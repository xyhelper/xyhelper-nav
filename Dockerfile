# ---------------------------------------------------------------------------
# xyhelper-nav — multi-stage image build.
#
#   docker build -t xyhelper/nav:latest .
#   docker run -d -p 8000:8000 -v $(pwd)/data:/app/resource/data xyhelper/nav:latest
#
# A freshly created host directory is owned by root, so the container entrypoint
# takes ownership of the mounted data directory and then drops privileges to the
# unprivileged "nav" user before starting the server.
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

# ca-certificates is needed to fetch remote site icons over HTTPS, tzdata keeps
# log timestamps in the expected timezone and su-exec lets the entrypoint drop
# root privileges after fixing the ownership of the data volume.
RUN apk add --no-cache ca-certificates tzdata su-exec \
    && adduser -D -u 10001 -h /app nav

WORKDIR /app

COPY --from=builder /out/nav /app/nav
COPY manifest/config /app/manifest/config
COPY resource/public /app/resource/public
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# The SQLite database, the icon cache and the session files live here; mount a
# volume over it to keep the data across container upgrades. The explicit 0755
# matters: the entrypoint must stay readable for the unprivileged user that
# re-executes it.
RUN chmod 0755 /usr/local/bin/docker-entrypoint.sh \
    && mkdir -p /app/resource/data \
    && chown -R nav:nav /app

EXPOSE 8000

# The entrypoint starts as root only to normalise the ownership of the mounted
# data directory and immediately re-executes the server as "nav".
USER root

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/nav"]
