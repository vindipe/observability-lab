# Day 1 and Day 2 operations

## Day 1: initial deployment and wiring

Day 1 operations are about making the system exist in a healthy initial state.

In this repository, Day 1 includes:

1. building the Go application image;
2. starting Redis, OpenTelemetry Collector, Prometheus, Alertmanager, Loki, Tempo and Grafana;
3. wiring Grafana datasources automatically;
4. provisioning dashboards;
5. exposing ports locally;
6. generating traffic;
7. verifying `/health`, `/metrics`, Prometheus targets and Grafana dashboards.

Commands:

```bash
make up
make traffic
make test
```

## Day 2: operating after deployment

Day 2 operations are the ongoing operational tasks that keep the platform reliable.

Examples mapped to this lab:

| Operation | Local demo example | Production/Juju/operator equivalent |
|---|---|---|
| Scale | Increase app replicas in Kubernetes | `juju scale-application`, HPA, charm action |
| Upgrade | Change image tag/version | rolling upgrade encoded in charm/operator logic |
| Troubleshooting | Inspect Grafana, Prometheus targets, logs and traces | runbooks, charm actions, automated health checks |
| Alert tuning | Edit `prometheus/alert-rules.yml` | COS configuration charm / GitOps alert rules |
| Backup/recovery | Preserve volumes, Redis snapshot | charm-managed backups, object storage, restore tests |
| Retention | Prometheus `--storage.tsdb.retention.time` | storage policies, compactor, lifecycle rules |
| Certificate rotation | not needed locally | charm relation + TLS certificates operator |
| Config rotation | edit `.env` / ConfigMaps | model config, charm config, relation-driven reloads |
| Rollback | revert image/config Git commit | charm revision rollback, workload version pinning |
| Incident response | trigger `/api/error`, inspect alert flow | on-call routing, SLOs, postmortem workflow |

## Why this matters for Cloud-Native Observability

cloud-native observability work is not just about installing Grafana. It is about encoding operational knowledge into reusable automation: charms, relations, models, health checks, upgrades, dashboards and alert rules.

This lab gives you a concrete story:

> I can deploy a telemetry-producing workload, route signals to the right backends, inspect metrics/logs/traces, trigger alerts, and explain how the same architecture would become more reliable with operators and Juju models.
