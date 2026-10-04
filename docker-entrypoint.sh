#!/bin/sh
# ---------------------------------------------------------------------------
# Container entrypoint for xyhelper-nav.
#
# The application itself runs as the unprivileged user "nav" (uid 10001).
# Docker creates a bind mount from an empty host directory as root:root, which
# that user cannot write to, so SQLite would fail with the rather cryptic
# "unable to open database file". When this container is started as root the
# data directory is therefore handed over to "nav" first and the process is
# re-executed without privileges. When it is started as any other user the
# entrypoint gets out of the way and simply runs the command.
# ---------------------------------------------------------------------------
set -e

DATA_DIR="${DATA_DIR:-/app/resource/data}"
APP_USER="${APP_USER:-nav}"

if [ "$(id -u)" = "0" ]; then
    mkdir -p "$DATA_DIR"
    # Only touch the ownership when it is actually wrong, so a large icon cache
    # is not rewritten on every start.
    if [ "$(stat -c '%u:%g' "$DATA_DIR")" != "$(id -u "$APP_USER"):$(id -g "$APP_USER")" ]; then
        echo "xyhelper-nav: taking ownership of $DATA_DIR for user $APP_USER"
        chown -R "$APP_USER:$APP_USER" "$DATA_DIR"
    fi
    exec su-exec "$APP_USER" "$@"
fi

exec "$@"
