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


## Docker build fails with missing go.sum entries

Symptom:

```text
missing go.sum entry for module providing package ...
```

The app image intentionally resolves Go modules inside Docker, so a fresh clone can start without requiring Go on the host. The Dockerfile copies `go.mod` and `main.go`, runs `go mod tidy` inside the build stage, and then compiles the binary.

If Docker is using a stale cached layer from an older generated version, rebuild the app image without cache:

```bash
docker compose build --no-cache app
make up
```

For local Go development outside Docker, populate the checksum file once:

```bash
cd app
go mod tidy
```

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

## App keeps restarting with `ServeMux` pattern conflict

If `docker compose ps` shows the app as restarting and the logs contain:

```text
panic: pattern "/metrics" conflicts with pattern "GET /"
```

then the root route is being registered as a subtree match. With the modern Go `net/http.ServeMux`, `GET /` can overlap with more specific paths such as `/metrics` when the latter is registered without the same method specificity.

The app fixes this by registering the root route as an exact match and by making the metrics route method-specific:

```go
mux.HandleFunc("GET /{$}", instrumentRoute("/", rootHandler))
mux.Handle("GET /metrics", promhttp.Handler())
```

After applying the fix, rebuild the app image and recreate the container:

```bash
docker compose build --no-cache app
make up
make test
```

## `make traffic` prints HTTP 500 for `/api/error`

That endpoint intentionally returns failures to generate error-rate metrics, logs,
traces, and alert traffic. In earlier versions of the lab the traffic script used
`curl -f`, which made those intentional HTTP 500 responses look like script
errors. The script now treats `/api/error` failures as expected demo traffic.

## Loki `/ready` returns HTTP 503

On some local Docker/WSL runs, Loki's `/ready` endpoint can remain stricter than
necessary for this lab even though the HTTP API is available. For smoke testing,
prefer:

```bash
curl -fsS http://localhost:3100/loki/api/v1/status/buildinfo
```

For deeper troubleshooting, inspect Loki logs:

```bash
docker compose logs --tail=200 loki
```

Then verify labels after generating traffic:

```bash
curl -G -s "http://localhost:3100/loki/api/v1/labels"
```
