#!/usr/bin/env bash
set -euo pipefail

SERVICES="${1:-}"
GH_DEPLOY_TOKEN="${GH_DEPLOY_TOKEN:-}"
GITHUB_REPOSITORY="${GITHUB_REPOSITORY:-benderxyz/avito-year-recap}"
DEPLOY_DIR="${DEPLOY_DIR:-$HOME/avito-year-recap}"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.prod.yml}"
ENV_FILE="${ENV_FILE:-.env.prod}"

cd "$DEPLOY_DIR"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Missing $ENV_FILE in $DEPLOY_DIR" >&2
  exit 1
fi

if [[ -n "$GH_DEPLOY_TOKEN" ]]; then
  git fetch origin main
  git checkout main
  git pull "https://x-access-token:${GH_DEPLOY_TOKEN}@github.com/${GITHUB_REPOSITORY}.git" main
else
  git fetch origin main
  git checkout main
  git pull origin main
fi

if [[ -z "$SERVICES" ]]; then
  echo "No services to deploy"
  exit 0
fi

# shellcheck disable=SC2086
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build $SERVICES
# shellcheck disable=SC2086
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d $SERVICES

echo "Deployed services: $SERVICES"
