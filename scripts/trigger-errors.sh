#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
COUNT="${COUNT:-80}"

echo "Triggering ${COUNT} /api/error requests against ${BASE_URL}"
for i in $(seq 1 "$COUNT"); do
  curl -sS -o /dev/null -w "request=%{http_code}\n" "${BASE_URL}/api/error" || true
  sleep 0.15
done

echo "Check Prometheus alerts at http://localhost:9090/alerts and Alertmanager at http://localhost:9093"
