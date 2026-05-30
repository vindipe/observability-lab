#!/usr/bin/env bash
set -euo pipefail

docker compose down -v --remove-orphans || true
kind delete cluster --name "${KIND_CLUSTER:-observability-lab}" || true
