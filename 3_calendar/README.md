## Календарь — запуск и тестирование через Docker Compose

Ниже — инструкции, как поднять сервис вместе с зависимостями (PostgreSQL, Redis, Mailhog), как выполнить запросы и прогнать тесты.

### Что поднимается
- **app**: Go-приложение (`cmd/calendar/main.go`), порт `8080`
- **postgres**: PostgreSQL 16, порт `5432`, БД/пользователь `calendar`, пароль `calendarpass`
- **redis**: Redis 7, порт `6379`, пароль `calendarpass`
- **mailhog**: локальный SMTP-сервер (порт `1025`) и веб-интерфейс по адресу `http://localhost:8025`

Инициализация схемы БД выполняется автоматически скриптом `docker/postgres/init.sql`.

### Быстрый старт
```bash
docker compose up --build
# Подождите, пока все сервисы станут healthy, затем приложение будет доступно на http://localhost:8080
```

Приложение читает конфиг из `configs/config.docker.yaml` и переменные окружения:
- `STORAGE_PASSWORD=calendarpass`
- `REDIS_PASSWORD=calendarpass`
- `SMTP_PASSWORD` — можно оставить пустым для Mailhog

### Полезные ссылки
- API: `http://localhost:8080`
- Здоровье: `GET http://localhost:8080/health` (ожидается 200)
- Метрики: `GET http://localhost:8080/metrics`
- Mailhog Web UI: `http://localhost:8025`

### Примеры запросов

Создать событие:
```bash
curl -sS -X POST http://localhost:8080/events \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Встреча",
    "description": "Обсуждение плана",
    "start": "2025-01-01T10:00:00Z",
    "end": "2025-01-01T11:00:00Z",
    "notify": true,
    "email": "user@example.com"
  }'
```

Список событий:
```bash
curl -sS http://localhost:8080/events | jq
```

Получить событие по ID:
```bash
curl -sS http://localhost:8080/events/<id>
```

Обновить событие:
```bash
curl -sS -X PUT http://localhost:8080/events/<id> \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Встреча (обновл.)",
    "description": "Правки",
    "start": "2025-01-01T10:30:00Z",
    "end": "2025-01-01T11:30:00Z",
    "notify": false,
    "email": "user@example.com"
  }'
```

Удалить событие:
```bash
curl -sS -X DELETE http://localhost:8080/events/<id> -i
```

Проверка здоровья:
```bash
curl -i http://localhost:8080/health
```

Метрики:
```bash
curl -sS http://localhost:8080/metrics | head -n 30
```

### Тестирование

Тесты зависят от Testcontainers для PostgreSQL. Запускайте их локально (в контейнере сервис уже работает). Нужен установленный Docker.

```bash
# Из корня проекта
go test ./...
```

При первом запуске тесты сами поднимут контейнер Postgres и применят минимальную схему.

### Конфигурация

Основной конфиг: `configs/config.yaml` (локальная разработка, хосты `localhost`).

Docker-конфиг: `configs/config.docker.yaml` (хосты: `postgres`, `redis`, `mailhog`). Переменные окружения:
- `CONFIG_PATH=/app/configs/config.docker.yaml`
- `STORAGE_PASSWORD` — пароль пользователя Postgres `calendar`
- `REDIS_PASSWORD` — пароль Redis
- `SMTP_PASSWORD` — пароль SMTP (для Mailhog можно пустым)

Файлы окружения: пример в `env.example`.

### OpenAPI

Спецификация `openapi.yaml` уже указывает сервер `http://localhost:8080`. Эндпоинты:
- `GET /health`
- `GET /metrics`
- `POST /events`
- `GET /events`
- `GET /events/{id}`
- `PUT /events/{id}`
- `DELETE /events/{id}`

Можно импортировать `openapi.yaml` в Postman/Insomnia/Swagger UI.

### Остановка и очистка
```bash
docker compose down
# Для полного удаления данных Postgres/Redis:
docker compose down -v
```


