# Go Metrics Server

Простой HTTP сервер на Go с встроенными метриками Prometheus. **Свои метрики не написаны** - используются стандартные метрики Go runtime, которые предоставляются сами по себе.

## Структура проекта

```
.
├── cmd/app/main.go          # Точка входа приложения
├── internal/delivery/http/  # HTTP сервер
├── docker-compose.yml       # Docker Compose для запуска
├── prometheus/              # Конфигурация Prometheus
└── grafana/                 # Конфигурация Grafana + дашборды
```

## Запуск

### Docker Compose (рекомендуется)

```bash
# Запуск всех сервисов
docker compose up -d

# Проверка статуса
docker compose ps
```

### Ручной запуск

```bash
# Запуск приложения
go run cmd/app/main.go

# В отдельном терминале - Prometheus
docker run -p 9090:9090 -v $(pwd)/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

# В отдельном терминале - Grafana
docker run -p 3000:3000 -e GF_SECURITY_ADMIN_PASSWORD=admin grafana/grafana
```

## Доступные сервисы

- **Приложение**: http://localhost:2112
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

## Метрики

Сервер автоматически предоставляет стандартные метрики Go runtime:

- **Scheduler**: `go_goroutines`, `go_threads`, `go_sched_gomaxprocs_threads`
- **Garbage Collector**: `go_gc_duration_seconds`, `go_memstats_*`
- **Memory**: `go_memstats_heap_*`, `go_memstats_alloc_*`
- **Process**: `process_cpu_seconds_total`, `process_resident_memory_bytes`

## Примеры запросов

### Получение метрик
```bash
curl http://localhost:2112/metrics
```

### Prometheus запросы
```bash
# Количество горутин
curl "http://localhost:9090/api/v1/query?query=go_goroutines"

# Время GC
curl "http://localhost:9090/api/v1/query?query=go_gc_duration_seconds"

# Использование памяти
curl "http://localhost:9090/api/v1/query?query=go_memstats_heap_alloc_bytes"
```

## Дашборд Grafana

Автоматически загружается дашборд "Go Runtime" с панелями:
- Scheduler (горутины, потоки, GOMAXPROCS)
- Garbage Collector (время паузы по квантилям)

## Остановка

```bash
docker compose down
```
