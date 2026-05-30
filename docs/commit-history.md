# Commit history

This repository was initialized with real Git commits. The generated history is:

1. `chore: initialize observability lab repository`
2. `feat(app): add Go HTTP service with demo endpoints`
3. `infra(observability): add OpenTelemetry, Prometheus and Alertmanager config`
4. `infra(grafana): add Loki Tempo Mimir configs and provisioned dashboard`
5. `infra(compose): add quickstart stack scripts and smoke tests`
6. `k8s: add local kind manifests for the observability stack`
7. `docs: add Cloud-Native technical review operations and ecosystem notes`
8. `docs: finalize professional README and usage guide`
9. `docs: align commit history with generated repository`
10. `style(app): format Go source`

Check the real history with:

```bash
git log --oneline --decorate
```

Note: `go.sum` is present but may be populated/updated by `go mod download` during the first Docker build because this artifact was generated in an environment without outbound Go module access.
