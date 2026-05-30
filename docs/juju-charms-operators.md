# Juju, charms and operators

## What is a Juju model?

A Juju model is an operational workspace. It groups applications, integrations, configuration, secrets and relations for a deployment target. You can have one model for staging, one for production, one for observability, one for data services, and so on.

In this lab, a conceptual split could be:

```text
model: workloads
  - orders-api charm
  - redis charm

model: observability
  - grafana-k8s
  - prometheus-k8s or mimir charms
  - loki-k8s
  - tempo-k8s
  - alertmanager-k8s
  - grafana-agent/alloy/otel collector
```

## What is a charm?

A charm is an operator package. It contains code that knows how to install, configure, integrate, upgrade and operate an application.

A charm is not only a deployment template. It encodes operational knowledge.

## Kubernetes operator vs Juju charm

| Concept | Kubernetes operator | Juju charm |
|---|---|---|
| Main environment | Kubernetes API/controller model | Juju model across Kubernetes, machines, clouds |
| Packaging | controller + CRDs | charm package with relations/actions/config |
| Integration | custom resources and service discovery | relation endpoints between charms |
| Operational knowledge | reconciliation logic | event handlers, relations, actions, config, lifecycle |
| Scope | usually Kubernetes-native | Kubernetes, machines, VMs, bare metal, clouds |

## How this lab maps to charms

| Lab component | Charm/operator equivalent |
|---|---|
| Go app deployment | custom `orders-api` charm |
| Redis container | redis charm or data platform charm |
| Prometheus config | prometheus-k8s + scrape relations |
| Grafana datasource provisioning | grafana-k8s relations and dashboards |
| Loki config | loki-k8s charm |
| Tempo config | tempo-k8s charm |
| Alert rules | COS configuration charm / alert rule relation |
| OpenTelemetry Collector | grafana-agent/alloy/collector charm pattern |

## Day 2 actions a charm could expose

- `scale` or application unit scaling;
- `backup` and `restore`;
- `rotate-certificates`;
- `reload-config`;
- `run-smoke-test`;
- `show-health`;
- `tune-alerts`;
- `set-log-retention`;
- `upgrade-workload`;
- relation-changed handlers for Prometheus, Grafana, Loki and Tempo.

## Technical Review sentence

> A Kubernetes manifest describes desired state. A charm goes further: it captures the operational lifecycle around that state, including relations, upgrades, configuration changes and Day 2 actions.
