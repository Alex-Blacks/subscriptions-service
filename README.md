# Subscriptions Service

REST API сервис для управления подписками пользователей.

---

## Описание

Сервис позволяет:

- создавать подписки
- получать подписку по ID
- обновлять подписку
- удалять подписку
- получать список подписок с фильтрами
- считать сумму подписок

---

## Технологии

- Go
- Chi router
- PostgreSQL
- pgx v5
- Docker / Docker Compose
- Swagger (swaggo)
- Migration tool (golang-migrate или аналог)

---

## Запуск проекта

### 1. Клонирование

git clone <repo_url>
cd subscriptions-service


### 2. Запуск через Docker

docker compose up --build

## После запуска:

API: http://localhost:8080
Swagger UI: http://localhost:8080/swagger/index.html
PostgreSQL: localhost:5432


## API Endpoints
###  Create subscription
    POST /subscriptions
    {
    "service_name": "yandex",
    "price": 100,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "05-2026",
    "end_date": "06-2026"
    }

### Get subscription
    GET /subscriptions/{id}

### Update subscription
    PATCH /subscriptions/{id}

### Delete subscription
    DELETE /subscriptions/{id}

### List subscriptions
    GET /subscriptions?user_id=&service_name=&from=&to=&limit=&offset=

### Sum subscriptions price
    GET /subscriptions/sum?user_id=&service_name=&from=&to=

### Фильтры
    | Параметр     | Тип     | Описание               |
    | ------------ | ------- | ---------------------- |
    | user_id      | uuid    | фильтр по пользователю |
    | service_name | string  | фильтр по сервису      |
    | from         | MM-YYYY | начало периода         |
    | to           | MM-YYYY | конец периода          |
    | limit        | int     | лимит                  |
    | offset       | int     | смещение               |



## Тестирование
### Unit tests
    go test ./...

## База данных
Миграции выполняются автоматически при старте контейнера

Или вручную:
    migrate -path ./internal/migrations \
        -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

## Swagger
### Swagger генерируется через swaggo:
    Генерация документации
    swag init -g cmd/main.go -o docs

### Открыть UI
    http://localhost:8080/swagger/index.html

## Форматы данных
   - Дата
   - вход: MM-YYYY
   - внутри системы: time.Time
   - в БД: DATE / TIMESTAMP

## Архитектура
handler → service → repository → storage → postgres

## Особенности
   - Все бизнес-ошибки возвращаются через domain errors
   - Валидация входных данных на уровне handler
   - База работает через pgx pool
   - Используются транзакции для write операций

## Примеры тестовых сценариев
   - создание подписки
   - повторное создание (ErrAlreadyExists)
   - удаление несуществующей записи
   - обновление частичное (PATCH)
   - sum по фильтрам
   - list с пустыми результатами

## Требования
- Go 1.22+
- Docker + Docker Compose