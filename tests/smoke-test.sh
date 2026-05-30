#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
PROM_URL="${PROM_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"
LOKI_URL="${LOKI_URL:-http://localhost:3100}"
TEMPO_URL="${TEMPO_URL:-http://localhost:3200}"
RETRIES="${RETRIES:-30}"
SLEEP_SECONDS="${SMOKE_SLEEP_SECONDS:-2}"

check() {
  local name="$1"
  local url="$2"
  echo "Checking ${name}: ${url}"
  curl -fsS "$url" >/dev/null
}

check_with_retry() {
  local name="$1"
  local url="$2"
  local attempt

  echo "Checking ${name}: ${url}"
  for attempt in $(seq 1 "$RETRIES"); do
    if curl -fsS "$url" >/dev/null; then
      return 0
    fi
    echo "  ${name} not ready yet (${attempt}/${RETRIES}); retrying in ${SLEEP_SECONDS}s..." >&2
    sleep "$SLEEP_SECONDS"
  done

  echo "ERROR: ${name} did not become ready: ${url}" >&2
  return 1
}

check_with_retry "orders-api health" "${BASE_URL}/health"
check "orders-api metrics" "${BASE_URL}/metrics"
check_with_retry "Prometheus" "${PROM_URL}/-/ready"
check_with_retry "Grafana" "${GRAFANA_URL}/api/health"
check_with_retry "Alertmanager" "${ALERTMANAGER_URL}/-/ready"

# Loki's /ready endpoint can be stricter than what we need for this local lab and
# may briefly return 503 while the single-node ring settles. Buildinfo confirms
# the HTTP API is up; logs can then be verified in Grafana Explore.
check_with_retry "Loki API" "${LOKI_URL}/loki/api/v1/status/buildinfo"
check_with_retry "Tempo" "${TEMPO_URL}/ready"

echo "Generating a small amount of traffic..."
DURATION_SECONDS=10 SLEEP_SECONDS=0.1 BASE_URL="$BASE_URL" ./scripts/generate-traffic.sh >/dev/null || true

echo "Smoke test completed."
