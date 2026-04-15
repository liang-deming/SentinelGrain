#!/bin/sh
set -e
HOST="${SENTINEL_REDIS_HOST:-redis:6379}"
sed "s|__REDIS_HOST__|${HOST}|g" /app/etc/sentinel.yaml.template > /tmp/sentinel.yaml
exec /app/sentinel-rpc -f /tmp/sentinel.yaml
