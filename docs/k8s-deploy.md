# Kubernetes production deploy — Gym Pro API

Deploy stack: **API** + **PostgreSQL (pgvector)** + **Redis** + **Nginx Ingress** on a VPS running Kubernetes (k3s, kubeadm, minikube, etc.).

CI/CD (GitHub Actions) applies manifests with `kubectl apply -k k8s/` after image push to GHCR.

## Prerequisites

- VPS with Kubernetes cluster and `kubectl` configured for that cluster
- GitHub Actions secrets: `VPS_HOST`, `VPS_USER`, `VPS_SSH_KEY`, `VPS_PORT`
- Docker image published: `ghcr.io/czx04/gym-pro-be:master`
- Ingress NGINX controller installed on the cluster
- Secrets created on the cluster (not in git):
  - `gym-pro-api-secret`
  - `ghcr-secret`

## Architecture

```
Internet
    → Ingress (api.gympro.example.com)
    → Service gym-pro-api-svc:80
    → Deployment gym-pro-api (x2) :8080
    → postgres-svc:5432 / redis-svc:6379
```

## First-time setup on VPS

### 1. Namespace

```bash
kubectl apply -f k8s/namespace.yaml
```

### 2. Application secrets

Edit placeholders, then run:

```bash
chmod +x k8s/scripts/create-secrets.example.sh
./k8s/scripts/create-secrets.example.sh
```

Or create manually:

```bash
kubectl create secret generic gym-pro-api-secret \
  -n gym-pro \
  --from-literal=POSTGRES_USER=gymadmin \
  --from-literal=POSTGRES_PASSWORD='<strong-password>' \
  --from-literal=POSTGRES_DB=gym_pro_db \
  --from-literal=DB_USER=gymadmin \
  --from-literal=DB_PASSWORD='<strong-password>' \
  --from-literal=JWT_SECRET='<long-random-string>' \
  --from-literal=REDIS_PASSWORD=''
```

`POSTGRES_*` and `DB_*` passwords must match. `REDIS_PASSWORD` can be empty when Redis has no auth.

### 3. GHCR pull secret

```bash
kubectl create secret docker-registry ghcr-secret \
  -n gym-pro \
  --docker-server=ghcr.io \
  --docker-username=czx04 \
  --docker-password='<GHCR_TOKEN_WITH_read:packages>'
```

Token: GitHub → Settings → Developer settings → Personal access tokens → `read:packages`.

### 4. Ingress NGINX controller

```bash
chmod +x k8s/scripts/install-ingress-nginx.sh
./k8s/scripts/install-ingress-nginx.sh
```

Helm is used when available; otherwise the official static manifest is applied.

### 5. Deploy manifests

```bash
kubectl apply -k k8s/
```

### 6. DNS

Point `api.gympro.example.com` to the Ingress external IP:

```bash
kubectl get svc -n ingress-nginx
```

Update `k8s/ingress.yaml` host and `k8s/configmap.yaml` `ALLOWED_ORIGINS` if you use a different domain.

## Verify

```bash
kubectl get pods -n gym-pro
kubectl get svc -n gym-pro
kubectl get ingress -n gym-pro
kubectl logs -n gym-pro deployment/gym-pro-api --tail=50
```

Health check via Ingress IP:

```bash
INGRESS_IP=$(kubectl get ingress gym-pro-api-ingress -n gym-pro -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl -H "Host: api.gympro.example.com" "http://${INGRESS_IP}/health"
```

Expected: `{"status":"ok","message":"Service is healthy"}`

## CI/CD deploy flow

On push to `master`:

1. Lint, test, build and push image to GHCR
2. Copy `k8s/` to VPS `/opt/gym-pro/k8s`
3. SSH: verify `kubectl`, required secrets, then:

```bash
kubectl apply -k /opt/gym-pro/k8s
kubectl rollout status statefulset/postgres -n gym-pro --timeout=180s
kubectl rollout status deployment/redis -n gym-pro --timeout=180s
kubectl rollout status deployment/gym-pro-api -n gym-pro --timeout=180s
```

If secrets are missing, the job fails with:

`Create gym-pro-api-secret and ghcr-secret on VPS before deploy.`

## Rollback API

```bash
kubectl rollout history deployment/gym-pro-api -n gym-pro
kubectl rollout undo deployment/gym-pro-api -n gym-pro
kubectl rollout status deployment/gym-pro-api -n gym-pro
```

## Troubleshooting

| Symptom | Check |
|---------|--------|
| `ImagePullBackOff` | `ghcr-secret`, token scope `read:packages` |
| API `CrashLoopBackOff` | `kubectl logs deployment/gym-pro-api -n gym-pro` — JWT, DB creds |
| DB connection refused | `kubectl get pods -n gym-pro -l app=postgres` |
| Ingress 404 | Host header, `kubectl describe ingress -n gym-pro` |
| PVC pending | StorageClass on cluster (`kubectl get sc`) |

## Docker Compose fallback

For manual / non-K8s VPS deploy:

```bash
docker compose -f docker-compose.prod.yml up -d
```

See `docker-compose.prod.yml` at repository root.

## Files not applied by kustomize

- `secrets.example.yaml`
- `ghcr-secret.example.yaml`
- `scripts/*`

These are documentation and examples only.
