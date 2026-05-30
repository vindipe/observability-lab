#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION_SECONDS="${DURATION_SECONDS:-90}"
SLEEP_SECONDS="${SLEEP_SECONDS:-0.25}"

request() {
  local path="$1"
  local expected_fail="${2:-false}"
  local status

  status="$(curl -sS -o /dev/null -w '%{http_code}' "${BASE_URL}${path}" || echo '000')"

  if [[ "$expected_fail" == "true" ]]; then
    # /api/error intentionally returns HTTP 500 sometimes. That is useful demo
    # traffic, not a script failure.
    return 0
  fi

  if [[ "$status" =~ ^2|3 ]]; then
    return 0
  fi

  echo "WARN: ${path} returned HTTP ${status}" >&2
  return 0
}

end=$((SECONDS + DURATION_SECONDS))
echo "Generating traffic for ${DURATION_SECONDS}s against ${BASE_URL}"

while [ "$SECONDS" -lt "$end" ]; do
  request "/health"
  request "/api/orders"
  request "/api/slow"
  request "/api/error" true
  sleep "${SLEEP_SECONDS}"
done

echo "Done. Open Grafana at http://localhost:3000 and inspect the provisioned dashboard."
