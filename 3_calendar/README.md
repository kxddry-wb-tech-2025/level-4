# Календарное приложение (Calendar App)

Веб-приложение для управления календарными событиями с возможностью настройки уведомлений по email.

## Описание

Приложение представляет собой REST API сервис, построенный на Go, который позволяет:
- Создавать, читать, обновлять и удалять календарные события
- Настраивать уведомления по email для событий
- Автоматически архивировать старые события
- Отправлять уведомления через SMTP

## Технологии

- **Backend**: Go 1.25.0
- **Web Framework**: Echo v4
- **База данных**: PostgreSQL 17
- **Кэш**: Redis 7
- **SMTP**: Поддержка SMTP серверов
- **Контейнеризация**: Docker & Docker Compose
- **Миграции**: Migrate v4
- **Логирование**: Zap
- **Валидация**: Validator v10
- **Мониторинг**: Prometheus метрики

## Архитектура

```
cmd/calendar/          # Точка входа приложения
internal/
├── config/           # Конфигурация
├── delivery/         # HTTP handlers и SMTP клиент
├── lib/             # Вспомогательные библиотеки
├── logging/         # Система логирования
├── models/          # Модели данных
├── service/         # Бизнес-логика
└── storage/         # Работа с БД и Redis
```

## Требования

- Go 1.25.0 или выше
- Docker и Docker Compose
- PostgreSQL 17
- Redis 7

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd calendar
```

### 2. Настройка переменных окружения

Скопируйте файл с примерами переменных окружения:

```bash
cp env.example .env
```

Отредактируйте `.env` файл, указав необходимые значения:

```bash
# SMTP настройки
SMTP_PASSWORD=your_smtp_password

# PostgreSQL настройки
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_USER=user
POSTGRES_DB=calendar

# Redis настройки
REDIS_PASSWORD=your_redis_password

# Временная зона
TZ=UTC

# Пути к конфигурации
CONFIG_PATH_DOCKER=./configs/config.docker.yaml
CONFIG_PATH=configs/config.yaml
STORAGE_PASSWORD=$POSTGRES_PASSWORD
```

### 3. Запуск с Docker Compose

Самый простой способ запуска - использование Docker Compose:

```bash
docker-compose up -d
```

Это запустит:
- Приложение на порту 8080
- PostgreSQL на порту 5432
- Redis на порту 6379
- MailHog (SMTP тестирование) на портах 1025 (SMTP) и 8025 (Web UI)
- Автоматические миграции БД

### 4. Проверка работоспособности

```bash
# Проверка здоровья сервиса
curl http://localhost:8080/health

# Получение метрик
curl http://localhost:8080/metrics

# Создание события
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Тестовое событие",
    "description": "Описание события",
    "start": "2024-12-25T10:00:00Z",
    "end": "2024-12-25T11:00:00Z",
    "notify": true,
    "email": "test@example.com"
  }'
```

## Локальная разработка

### 1. Установка зависимостей

```bash
go mod download
```

### 2. Запуск зависимостей

```bash
# Запуск только БД и Redis
docker-compose up -d postgres redis mailhog

# Ожидание готовности PostgreSQL
docker-compose exec postgres pg_isready -U user -d calendar
```

### 3. Применение миграций

```bash
# Автоматически через Docker Compose
docker-compose up migrator

# Или вручную
docker run --rm -v $(pwd)/docker/postgres:/migrations \
  --network calendar_default \
  migrate/migrate:v4.17.1 \
  -path=/migrations \
  -database "postgres://user:password@localhost:5432/calendar?sslmode=disable" \
  up
```

### 4. Запуск приложения

```bash
# Установка переменных окружения
export CONFIG_PATH=configs/config.yaml
export STORAGE_PASSWORD=password
export REDIS_PASSWORD=password
export SMTP_PASSWORD=password

# Запуск
go run cmd/calendar/main.go
```

### 5. Запуск тестов

```bash
# Все тесты
go test ./...

# Тесты с покрытием
go test -cover ./...

# Тесты конкретного пакета
go test ./internal/service/events/...
```

## Конфигурация

### Основные настройки

Файл `configs/config.yaml` содержит основные настройки:

```yaml
env: dev

server:
  port: 8080
  timeout: 30s
  idle_timeout: 60s

smtp:
  host: localhost
  port: 1025
  username: "email@example.com"

redis:
  host: localhost
  port: 6379
  username: default
  database: 0

storage:
  host: localhost
  port: 5432
  username: user
  database: calendar
  sslmode: disable

worker:
  interval: 1s
  limit: 100

archiver:
  interval: 30s
  older_than: 30d
  batch_size: 100
```

### Docker конфигурация

Для Docker используется `configs/config.docker.yaml` с настройками для контейнеров.

## API Endpoints

### Основные эндпоинты

- `GET /health` - Проверка здоровья сервиса
- `GET /metrics` - Prometheus метрики
- `POST /events` - Создание события
- `GET /events` - Получение списка событий
- `GET /events/{id}` - Получение события по ID
- `PUT /events/{id}` - Обновление события
- `DELETE /events/{id}` - Удаление события

### Модель события

```json
{
  "id": "uuid",
  "title": "Название события",
  "description": "Описание события",
  "start": "2024-12-25T10:00:00Z",
  "end": "2024-12-25T11:00:00Z",
  "notify": true,
  "email": "user@example.com"
}
```

## Развертывание

### Development окружение

```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f app

# Остановка
docker-compose down
```

### Production окружение

1. **Подготовка переменных окружения**
   ```bash
   # Создайте .env.prod с production значениями
   cp env.example .env.prod
   # Отредактируйте .env.prod
   ```

2. **Сборка и запуск**
   ```bash
   # Сборка production образа
   docker build -t calendar:prod .

   # Запуск с production конфигурацией
   docker run -d \
     --name calendar-prod \
     --env-file .env.prod \
     -p 8080:8080 \
     calendar:prod
   ```

3. **Использование внешних сервисов**
   - Замените MailHog на реальный SMTP сервер
   - Настройте внешний PostgreSQL и Redis
   - Настройте reverse proxy (nginx)
   - Настройте SSL/TLS

## Мониторинг и логирование

### Метрики

Приложение предоставляет Prometheus метрики по адресу `/metrics`.

### Логирование

Логи собираются через Zap logger и выводятся в stdout/stderr.

### Health Check

Эндпоинт `/health` для проверки состояния сервиса.

## Устранение неполадок

### Частые проблемы

1. **PostgreSQL не доступен**
   ```bash
   docker-compose logs postgres
   docker-compose exec postgres pg_isready -U user -d calendar
   ```

2. **Redis не доступен**
   ```bash
   docker-compose logs redis
   docker-compose exec redis redis-cli -a password ping
   ```

3. **SMTP ошибки**
   - Проверьте настройки SMTP в конфигурации
   - Для тестирования используйте MailHog (http://localhost:8025)

4. **Миграции не применяются**
   ```bash
   docker-compose logs migrator
   docker-compose up migrator
   ```

### Просмотр логов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis
```

## Разработка

### Структура проекта

- `cmd/` - Точки входа приложений
- `internal/` - Внутренний код приложения
- `configs/` - Конфигурационные файлы
- `docker/` - Docker файлы и миграции
- `testutil/` - Утилиты для тестирования

### Добавление новых функций

1. Создайте модель в `internal/models/`
2. Добавьте репозиторий в `internal/storage/repos/`
3. Реализуйте сервис в `internal/service/`
4. Создайте HTTP handler в `internal/delivery/http/`
5. Добавьте тесты

### Code Style

- Используйте `go fmt` для форматирования
- Запускайте `go vet` для проверки кода
- Следуйте стандартным Go conventions

## Лицензия

[Укажите лицензию проекта]

## Поддержка

[Укажите контактную информацию для поддержки]


