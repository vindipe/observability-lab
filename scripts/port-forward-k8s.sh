#!/usr/bin/env bash
set -euo pipefail

NS="${K8S_NAMESPACE:-observability-lab}"

echo "Starting port-forwards in namespace ${NS}. Press Ctrl+C to stop all."

kubectl -n "$NS" port-forward svc/orders-api 8080:8080 &
kubectl -n "$NS" port-forward svc/grafana 3000:3000 &
kubectl -n "$NS" port-forward svc/prometheus 9090:9090 &
kubectl -n "$NS" port-forward svc/alertmanager 9093:9093 &
kubectl -n "$NS" port-forward svc/loki 3100:3100 &
kubectl -n "$NS" port-forward svc/tempo 3200:3200 &

wait
