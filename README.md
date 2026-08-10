# Avito Recap

Generate personalized yearly recap cards from Avito activity.

### Always Up-Time

Recap Engine Library

[![Recap Engine CI](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml/badge.svg?v=1)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml)
[![Recap Engine Release](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml/badge.svg?v=1)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml)
[![Recap Engine Release](https://github.com/benderxyz/avito-year-recap/actions/workflows/deploy.yml/badge.svg?v=1)](https://github.com/benderxyz/avito-year-recap/actions/workflows/deploy.yml)

### Recap Engine NPM packages


[![@recap-engine/core NPM version](https://img.shields.io/npm/v/@recap-engine/core?style=flat-square&label=@recap-engine/core&color=CB3837&logo=npm&v=1)](https://www.npmjs.com/package/@recap-engine/core)
[![@recap-engine/react NPM version](https://img.shields.io/npm/v/@recap-engine/core?style=flat-square&label=@recap-engine/react&color=CB3837&logo=npm&v=1)](https://www.npmjs.com/package/@recap-engine/react)

### Links

[![Recap Engine Documentation Site](https://img.shields.io/badge/Recap%20Engine-View%20Documentation-3ECC5F?logo=Docusaurus&v=1)](https://recaps.hakolr.dev/docs)
[![Recap Engine Documentation Site](https://img.shields.io/badge/Recaps%20Engine-View%20Demo-3ECC5F&v=1)](https://recaps.hakolr.dev)

# Структура проекта

Монорепа: три Go-сервиса, React-фронтенд и npm-библиотека Recap Engine.

```
.
├── analytics-service/     Go, порт 8080. События и метрики, хранилище ClickHouse
├── cards-service/         Go, порт 8081. Сборка recap, бейджи, шеринг, LLM-обогащение
├── user-service/          Go, порт 8082. Профили пользователей, база users
├── frontend/              React + Vite, порт 3000. Витрина и модалка recap
├── recap-engine/          pnpm-воркспейс: npm-пакеты и сайт документации
│   ├── packages/core/     @recap-engine/core, движок сцен
│   ├── packages/react/    @recap-engine/react, React-обвязка
│   └── docs/              Docusaurus, деплоится как сервис docs
├── seed-data/             Тестовые данные
│   ├── seed-script/       Go-скрипт, заливает профили и события через API
│   └── users/             JSON-профили для сидов
├── infra/                 Конфигурация окружения
│   ├── nginx/             Реверс-прокси для прода
│   └── postgres/init/     Создание баз cards и users при первом старте
├── scripts/deploy.sh      Деплой на VPS, запускается из GitHub Actions
├── docs/                  Архитектура и заметки по Recap Engine (не Docusaurus)
├── postman/               Коллекция и окружение для ручных запросов
└── .github/workflows/     CI сервисов и фронта, релиз npm, деплой
```

Каждый Go-сервис устроен одинаково: `cmd/` — точка входа, `internal/` — логика,
`migrations/` — SQL-миграции, свой `Dockerfile` и `README.md`. Модули связаны
через `go.work` в корне.

В корне лежат общие конфиги: `docker-compose.yml` для локального запуска,
`docker-compose.prod.yml` для VPS, `biome.json` — один на `frontend` и
`recap-engine`, `.golangci.yaml` — один на все Go-модули.

# Команда

- **Максим Шудегов** — `analytics-service` (ClickHouse, ingest событий, агрегации, глобальная аналитика и кэш запросов), `seed-data`, рефакторинг бейджей и правил в `cards-service` с вынесением story, badges и metric definitions в Postgres, список пользователей в `user-service`, prod-конфиг LLM и proxy.

- **Максим Гайворонский** – `recap-engine`: @recap-engine/core и @recap-engine/react. На базе библиотеки собрал `frontend` на React + Vite. Для Recap Engine поднял сайт документации на Docusaurus https://recaps.hakolr.dev/docs/. Настроил деплой и CI/CD для фронтенд проектов: Link, Typecheck, Tests frontend и recap-engine, релиз npm-пакетов и деплой на VPS через Docker Compose с nginx как reverse proxy. Деплой идёт по push в main с детектом изменений — пересобираются только затронутые сервисы (frontend, backend, docs).

- **Никита Макаров** – репозиторий и базовая структура монорепы, общий линтер на все сервисы проекта. Реализовал `user-service` на Postgres: схема и миграции, профили, тесты, разделение логики по слоям. Сделал рекомендации после завершения рекапа. Проработал data-driven правила для карточек и бейджей: условия бейджей и сцен story уехали из `if`-лестниц в `generate` в правила с предикатом в Postgres, пороги и тексты меняются без релиза сервиса. Отрефакторил `cards-service` под общение с фронтом, добавил CORS. Перевёл Docker-образы всех сервисов на non-root пользователей, почистил контекст сборки. Разработка описания командного решения и согласование его в команде.

- **Дарина Ионова** – `cards-service`: генерация recap-пейлоада для фронтового виджета, интеграция с `analytics-service`, поддержка диапазона лет в запросах метрик. Интегрировала LLM: OpenAI-клиент с ретраями через `/responses`, провайдеры через интерфейс, систeмный промпт копирайтера на русском для переписывания бейджей и генерации инсайта без цифр в тексте, санитизация и валидация ответа по длинам полей, кэш в Postgres с версионированием промпта (хэш промпта и модели в ключе кэша), fallback на базовый recap при недоступности LLM или невалидном ответе, `LLM_ENABLED`-флаг для отключения. Добавила primary key и публичные метрики для LLM. Покрыла `llm`-пакет тестами.

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

Переменные для БД смотри в `.env.example`.

Postgres один контейнер, порт `5432`. База `cards` для cards-service. База `users` для user-service. Init скрипт лежит в `infra/postgres/init/`.

Локальный пароль ClickHouse по умолчанию `recap`. Пароль Postgres тоже `recap`.

## Сиды

Rules для cards берутся из `cards-service/seeds/`, если `SEED_ON_START=true`. В compose это уже включено. Миграции и сиды cards пишутся в базу `cards`.

Пользователей и события заливает `seed-data/seed-script` через API. Профили уходят в user-service (база `users`). События уходят в analytics.

Smoke после `docker compose up --build`

```sh
curl -s localhost:8080/health
curl -s localhost:8081/health
curl -s localhost:8082/health
go -C seed-data/seed-script run . -user 42 -year 2026
curl -s 'localhost:8080/users/42/metrics?from=2026-01-01T00:00:00Z&to=2027-01-01T00:00:00Z'
```

Поднять стек для локальных сидов (Postgres, user, analytics, clickhouse, cards):

```sh
docker compose up -d postgres clickhouse user analytics cards
```

cards накатит schema и seed rules в `cards`. user накатит schema в `users`. analytics нужен seed-script для событий. После health можно запускать seed-script.

Если volume `postgres-data` уже есть, а баз `cards` или `users` в нём нет, init не повторится. Тогда сбрось volume проекта и подними заново.

```sh
docker compose down -v
docker volume rm avito-year-recap_cards-postgres-data avito-year-recap_user-postgres-data 2>/dev/null || true
docker compose up --build
```

Go-линтер для сервисов запускается из корня проекта:

```sh
golangci-lint run -c .golangci.yaml ./cards-service/... ./analytics-service/... ./user-service/...
```
