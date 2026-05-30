# Optional Mimir mode

The base demo uses Prometheus directly because it is easier to run locally and easier to debug in an technical review.

Mimir is included as an optional profile to discuss scalable, horizontally sharded, multi-tenant metrics storage.

Run:

```bash
docker compose --profile advanced up -d --build
```

Then switch Prometheus to `prometheus/prometheus-mimir.yml` if you want to demonstrate `remote_write`.
For a real production setup, this would normally be deployed with Helm/Juju/charmed operators, object storage, replication, proper tenants, TLS, auth and retention policies.
