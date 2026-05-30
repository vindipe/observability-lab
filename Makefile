SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

COMPOSE ?= docker compose
KIND_CLUSTER ?= observability-lab
K8S_NAMESPACE ?= observability-lab
APP_IMAGE ?= observability-lab-app:local

.PHONY: help
help: ## Show available commands
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-18s %s\n", $$1, $$2}'

.PHONY: up
up: ## Start the Docker Compose lab
	$(COMPOSE) up -d --build

.PHONY: down
down: ## Stop the Docker Compose lab
	$(COMPOSE) down

.PHONY: logs
logs: ## Follow Docker Compose logs
	$(COMPOSE) logs -f --tail=120

.PHONY: traffic
traffic: ## Generate demo traffic against the Go app
	./scripts/generate-traffic.sh

.PHONY: errors
errors: ## Generate error traffic to trigger alerts
	./scripts/trigger-errors.sh

.PHONY: test
test: ## Run smoke tests against the local Docker Compose lab
	./tests/smoke-test.sh

.PHONY: clean
clean: ## Stop and remove demo volumes
	$(COMPOSE) down -v --remove-orphans

.PHONY: mimir-up
mimir-up: ## Start optional Mimir profile in addition to the base stack
	$(COMPOSE) --profile advanced up -d --build

.PHONY: promtail-up
promtail-up: ## Start optional Promtail side pipeline for log shipping comparison
	$(COMPOSE) --profile promtail up -d --build

.PHONY: k8s-up
k8s-up: ## Create a local kind cluster
	kind create cluster --name $(KIND_CLUSTER) --config k8s/kind-config.yaml

.PHONY: k8s-deploy
k8s-deploy: ## Build app image, load it into kind, and deploy manifests
	docker build -t $(APP_IMAGE) ./app
	kind load docker-image $(APP_IMAGE) --name $(KIND_CLUSTER)
	kubectl apply -f k8s/

.PHONY: k8s-forward
k8s-forward: ## Port-forward app, Grafana, Prometheus, Alertmanager, Loki and Tempo from Kubernetes
	./scripts/port-forward-k8s.sh

.PHONY: k8s-clean
k8s-clean: ## Delete Kubernetes resources and kind cluster
	-kubectl delete namespace $(K8S_NAMESPACE)
	-kind delete cluster --name $(KIND_CLUSTER)
