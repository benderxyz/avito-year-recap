Запуск контейнера:
docker-compose up --build

Остановка контейнера:
docker-compose down

Проверка что сервисы запущены успешно:
curl localhost:8080/health 
curl localhost:8081/health

Переменные для бд стоят пока заглушки, добавьте нужные данные при необходимости в .env
