#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${NAMESPACE:-gym-pro}"

kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic gym-pro-api-secret \
  --namespace="$NAMESPACE" \
  --from-literal=POSTGRES_USER=gymadmin \
  --from-literal=POSTGRES_PASSWORD='CHANGE_ME' \
  --from-literal=POSTGRES_DB=gym_pro_db \
  --from-literal=DB_USER=gymadmin \
  --from-literal=DB_PASSWORD='CHANGE_ME' \
  --from-literal=JWT_SECRET='CHANGE_ME_LONG_RANDOM_STRING' \
  --from-literal=REDIS_PASSWORD='' \
  --from-literal=SMTP_USERNAME='' \
  --from-literal=SMTP_PASSWORD='' \
  --from-literal=SENDGRID_API_KEY='' \
  --from-literal=GOOGLE_CLIENT_ID='' \
  --from-literal=GOOGLE_CLIENT_SECRET='' \
  --from-literal=FACEBOOK_APP_ID='' \
  --from-literal=FACEBOOK_APP_SECRET='' \
  --from-literal=CLOUDINARY_URL='' \
  --from-literal=GEMINI_API_KEY='' \
  --from-literal=EXPO_ACCESS_TOKEN='' \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret docker-registry ghcr-secret \
  --namespace="$NAMESPACE" \
  --docker-server=ghcr.io \
  --docker-username=czx04 \
  --docker-password='CHANGE_ME_GHCR_TOKEN' \
  --dry-run=client -o yaml | kubectl apply -f -

echo "Secrets applied in namespace $NAMESPACE"
echo "Replace CHANGE_ME values before production use."
