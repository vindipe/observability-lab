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
11. `fix(app): preserve generated Go module checksums during Docker build`
12. `fix(app): resolve Go module checksums during Docker build`
13. `docs: add troubleshooting note for Go checksum builds`

Check the real history with:

```bash
git log --oneline --decorate
```

Note: this repository keeps Docker as the primary reproducible path. The Dockerfile resolves modules with `go mod tidy` inside the build stage after copying the actual source imports. This avoids requiring Go on the host while still producing a valid checksum set inside the image build. If you want to develop the Go app directly on the host, run `cd app && go mod tidy` once to populate the local checksum file.
