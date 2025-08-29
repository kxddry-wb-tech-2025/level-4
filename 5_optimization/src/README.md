# Level-0 - Go Microservice with Kafka, Cache, and PostgreSQL

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-000?style=for-the-badge&logo=apachekafka)
![PostgreSQL](https://img.shields.io/badge/postgresql-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

## 📋 Описание проекта

Level-0 представляет собой современный микросервис, разработанный на языке Go, который демонстрирует архитектуру распределенной системы с использованием ключевых технологий:

- **Apache Kafka** - для асинхронной обработки сообщений и построения event-driven архитектуры
- **PostgreSQL** - как основная база данных для персистентного хранения
- **Redis/Cache** - для кэширования данных и повышения производительности
- **TypeScript Frontend** - современный веб-интерфейс для взаимодействия с API
- **Docker** - для контейнеризации и упрощения развертывания

Проект реализует принципы чистой архитектуры и следует лучшим практикам разработки микросервисов.

## 🏗️ Архитектура проекта

```
level-0/
├── cmd/                    # Точки входа приложения
├── config/                 # Конфигурационные файлы
├── frontend/              # TypeScript/React фронтенд
├── internal/              # Внутренняя логика приложения
│   ├── handlers/          # HTTP обработчики
│   ├── services/          # Бизнес-логика
│   ├── repository/        # Слой работы с данными
│   ├── models/            # Модели данных
│   └── kafka/             # Kafka продюсеры и консьюмеры
├── migrations/            # SQL миграции базы данных
├── nginx/                 # Конфигурация nginx
├── docker-compose.yml     # Orchestration для локальной разработки
├── Dockerfile            # Образ для контейнеризации
├── go.mod                # Go модули и зависимости
└── model.json            # Модель данных/API спецификация
```

## 🚀 Быстрый старт

### Предварительные требования

Убедитесь, что у вас установлены

- [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)

### Установка и запуск

1. **Клонирование репозитория:**
   ```bash
   git clone https://github.com/kxddry-wb-tech-2025/level-0.git
   cd level-0
   ```

2. **Запуск инфраструктуры через Docker Compose:**
   ```bash
   # Запуск всех сервисов (PostgreSQL, Kafka, Nginx, Go, Sender, ...)
   docker-compose up -d
   ```

### Конфигурация

Пример конфигурации расположен в папке `config/`. Убедитесь, что создан файл .env с POSTGRES_PASSWORD.

## 🛠️ Разработка

### Структура кода

- **cmd/** - содержит точки входа для различных компонентов (основное приложение, миграции, воркеры)
- **internal/** - внутренняя логика приложения, недоступная для импорта извне
- **config/** - конфигурационные файлы и структуры
- **migrations/** - SQL скрипты для создания и изменения схемы БД

### Основные компоненты

1. **HTTP API** - RESTful API для взаимодействия с клиентами
2. **Kafka Producer/Consumer** - для обработки асинхронных событий
3. **Repository Layer** - абстракция для работы с базой данных
4. **Cache Layer** - Redis для кэширования частых запросов
5. **Business Logic** - сервисный слой с основной логикой
