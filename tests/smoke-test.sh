#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
PROM_URL="${PROM_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"
LOKI_URL="${LOKI_URL:-http://localhost:3100}"
TEMPO_URL="${TEMPO_URL:-http://localhost:3200}"

check() {
  local name="$1"
  local url="$2"
  echo "Checking ${name}: ${url}"
  curl -fsS "$url" >/dev/null
}

check "orders-api health" "${BASE_URL}/health"
check "orders-api metrics" "${BASE_URL}/metrics"
check "Prometheus" "${PROM_URL}/-/ready"
check "Grafana" "${GRAFANA_URL}/api/health"
check "Alertmanager" "${ALERTMANAGER_URL}/-/ready"
check "Loki" "${LOKI_URL}/ready"
check "Tempo" "${TEMPO_URL}/ready"

echo "Generating a small amount of traffic..."
DURATION_SECONDS=10 SLEEP_SECONDS=0.1 BASE_URL="$BASE_URL" ./scripts/generate-traffic.sh >/dev/null || true

echo "Smoke test completed."
