# Local Kubernetes mode

The Kubernetes mode is intentionally more explicit than the Docker Compose mode. It is useful for technical review discussion because it shows how the same observability topology maps to deployments, services, config maps and health probes.

## Prerequisites

- Docker Desktop running
- `kind`
- `kubectl`

## Run

```bash
make k8s-up
make k8s-deploy
make k8s-forward
```

Then open:

- App: <http://localhost:8080/health>
- Grafana: <http://localhost:3000>
- Prometheus: <http://localhost:9090>
- Alertmanager: <http://localhost:9093>

## Notes

This is not meant to replace production-grade Helm charts or Juju charms. It is a local manifest-based translation of the Docker Compose architecture so you can explain:

- Deployments vs Services;
- ConfigMaps for configuration injection;
- probes as health checks;
- Prometheus scraping Kubernetes services;
- why operators/charmed operators would reduce manual Day 2 work.
