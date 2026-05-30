# Juju concept mapping

This directory contains a small conceptual charm skeleton. It is not meant to be a production charm; it exists to demonstrate that the lab can be reasoned about through Juju's model/charm/operator vocabulary.

## Conceptual deployment

```bash
juju add-model observability-lab
juju deploy ./charm-skeleton orders-api
juju deploy redis-k8s redis
juju deploy cos-lite --trust
juju relate orders-api prometheus:metrics-endpoint
juju relate orders-api grafana:grafana-dashboard
juju relate orders-api loki:logging
```

The exact relation names depend on real charms. The point for the technical review is the model: charms express operational knowledge and relations wire applications together.

## Why this matters

an operator-driven observability story is not just a set of containers. It is an operator-driven stack where dashboards, scrape jobs, logs, traces, alerts and lifecycle actions can be encoded and reused.
