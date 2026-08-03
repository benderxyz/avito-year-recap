[![Recap Engine CI](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-ci.yml)
[![Recap Engine Release](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml/badge.svg)](https://github.com/benderxyz/avito-year-recap/actions/workflows/recap-engine-release.yml)

Запуск контейнера:
docker-compose up --build

Остановка контейнера:
docker-compose down

Проверка что сервисы запущены успешно:
curl localhost:8080/health 
curl localhost:8081/health

Переменные для бд стоят пока заглушки, добавьте нужные данные при необходимости в .env
