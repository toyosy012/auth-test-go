#!/bin/ash
# wait-for-it.sh

set -e

until nc -z auth-test-db 3306; do
  >&2 echo "auth-test-db is unavailable - sleeping"
  sleep 3
done
>&2 echo "redis and auth-test-db is up - executing command"

exec $@
