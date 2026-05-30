#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION_SECONDS="${DURATION_SECONDS:-90}"
SLEEP_SECONDS="${SLEEP_SECONDS:-0.25}"

end=$((SECONDS + DURATION_SECONDS))
echo "Generating traffic for ${DURATION_SECONDS}s against ${BASE_URL}"

while [ "$SECONDS" -lt "$end" ]; do
  curl -fsS "${BASE_URL}/health" >/dev/null || true
  curl -fsS "${BASE_URL}/api/orders" >/dev/null || true
  curl -fsS "${BASE_URL}/api/slow" >/dev/null || true
  curl -fsS "${BASE_URL}/api/error" >/dev/null || true
  sleep "${SLEEP_SECONDS}"
done

echo "Done. Open Grafana at http://localhost:3000 and inspect the provisioned dashboard."
