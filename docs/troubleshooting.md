# Troubleshooting

## Docker Compose services do not start

```bash
docker compose ps
docker compose logs -f --tail=200
```

Common causes:

- Docker Desktop not running;
- ports `3000`, `8080`, `9090`, `9093`, `3100`, `3200`, `4317` already used;
- first image pulls are slow;
- WSL/Docker Desktop file sharing issues.

## App works but no metrics in Prometheus

Check:

```bash
curl http://localhost:8080/metrics
open http://localhost:9090/targets
```

Prometheus target `demo-app` should be UP.

## No traces in Tempo/Grafana

Generate traffic first:

```bash
make traffic
```

Then check collector and Tempo logs:

```bash
docker compose logs otel-collector tempo
```

The app sends OTLP traces to `otel-collector:4317`; the collector exports them to `tempo:4317`.

## No logs in Loki

The default path is:

```text
orders-api -> /var/log/observability-lab/app.log -> OTel filelog receiver -> Loki OTLP endpoint
```

Check:

```bash
docker compose exec app sh -lc 'ls -l /var/log/observability-lab && tail /var/log/observability-lab/app.log'
docker compose logs otel-collector loki
```

If the OTLP logs path is problematic in your environment, run the optional Promtail profile:

```bash
make promtail-up
```

## Alerts do not fire

Run:

```bash
make errors
```

Then open:

- Prometheus alerts: <http://localhost:9090/alerts>
- Alertmanager: <http://localhost:9093>

Alerts have short demo windows but still need around 30-60 seconds.

## Kubernetes image pull error

Use:

```bash
make k8s-deploy
```

This builds `observability-lab-app:local` and loads it into kind.
