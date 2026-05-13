# Gym Pro — Kubernetes

Production manifests for namespace `gym-pro`: API, PostgreSQL, Redis, Ingress.

Full guide: [docs/k8s-deploy.md](../docs/k8s-deploy.md)

## Layout

| File | Purpose |
|------|---------|
| `namespace.yaml` | Namespace `gym-pro` |
| `configmap.yaml` | Non-sensitive app config |
| `postgres.yaml` | StatefulSet + PVC + `postgres-svc` |
| `redis.yaml` | Deployment + `redis-svc` |
| `deployment.yaml` | API (2 replicas, GHCR image) |
| `service.yaml` | `gym-pro-api-svc` :80 → :8080 |
| `ingress.yaml` | Nginx Ingress host `api.gympro.example.com` |
| `secrets.example.yaml` | Template only — do not apply as-is |
| `ghcr-secret.example.yaml` | Template only — use `kubectl create` |
| `scripts/create-secrets.example.sh` | Example secret creation |
| `scripts/install-ingress-nginx.sh` | Install Ingress NGINX controller |

## Quick start

```bash
kubectl apply -f k8s/namespace.yaml
./k8s/scripts/create-secrets.example.sh
./k8s/scripts/install-ingress-nginx.sh
kubectl apply -k k8s/
```

## Verify

```bash
kubectl get pods,svc,ingress -n gym-pro
curl -H "Host: api.gympro.example.com" http://<INGRESS_IP>/health
```

## Fallback

Docker Compose deploy: `docker-compose.prod.yml` at repo root.
