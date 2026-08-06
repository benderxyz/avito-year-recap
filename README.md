# Avito Recap

Generate personalized yearly recap cards from Avito activity.

### Always Up-Time

Recap Engine Library

[![Recap Engine CI](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml)
[![Recap Engine Release](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml)
[![Recap Engine Release](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-docs-deploy.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-docs-deploy.yml)

### Recap Engine NPM packages

[![@recap-engine/core NPM version](https://img.shields.io/npm/v/@recap-engine/core?style=flat-square&label=@recap-engine/core&color=CB3837&logo=npm)](https://www.npmjs.com/package/@recap-engine/core)
[![@recap-engine/react NPM version](https://img.shields.io/npm/v/@recap-engine/core?style=flat-square&label=@recap-engine/react&color=CB3837&logo=npm)](https://www.npmjs.com/package/@recap-engine/react)

### Documentation

[![Recap Engine Documentation Site](https://img.shields.io/badge/Recap%20Engine-View%20Documentation-3ECC5F?logo=Docusaurus)](https://recap.hakolr.dev)

# Быстрый старт

Запуск контейнера:
```sh
docker-compose up --build
```

Остановка контейнера:
```sh
docker-compose down
```

Проверка что сервисы запущены успешно:
```sh
curl localhost:8080/health 
curl localhost:8081/health
curl localhost:8082/health
```

Переменные для БД смотри в `.env.example`. Локальный пароль ClickHouse по умолчанию `recap`.

Smoke после `docker compose up --build`

```sh
curl -s localhost:8080/health
curl -s localhost:8082/health
go -C seed-data/seed-avito run . -user 42 -year 2026
curl -s 'localhost:8080/users/42/metrics?from=2026-01-01T00:00:00Z&to=2027-01-01T00:00:00Z'
```
