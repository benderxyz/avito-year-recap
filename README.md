# Avito Recap

Generate personalized yearly recap cards from Avito activity.

[![Recap Engine CI](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml)
[![Recap Engine Release](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml)
[![@recap-engine/core NPM version](https://img.shields.io/npm/v/@recap-engine/core?style=flat-square&label=npm@recap-engine/core&color=CB3837)](https://www.npmjs.com/package/@recap-engine/core)
[![@recap-engine/react NPM version](https://img.shields.io/npm/v/@recap-engine/core?style=flat-square&label=npm@recap-engine/react&color=CB3837)](https://www.npmjs.com/package/@recap-engine/react)

## Backend:
![Static Badge](https://img.shields.io/badge/Golang--%2300ADD8?logo=go) 
![Static Badge](https://img.shields.io/badge/ClickHouse--%23FFCC01?logo=ClickHouse)

## Frontend
![Static Badge](https://img.shields.io/badge/React--%2361DAFB?logo=react)
![Static Badge](https://img.shields.io/badge/TypeScript--%233178C6?logo=typescript)


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
```

Переменные для бд стоят пока заглушки, добавьте нужные данные при необходимости в .env
