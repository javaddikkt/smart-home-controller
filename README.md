
# Smart Home Controller

**Smart Home Controller** — серверное приложение на языке Go для мониторинга и обработки данных от домашних датчиков. Предоставляет REST API и WebSocket-соединение для управления и получения информации о состоянии устройств умного дома в реальном времени.

## Стек технологий

- **Go**
- **Gin, WebSocket**
- **Swagger/OpenAPI**
- **PostgreSQL**
- **Docker, Docker Compose**
- **Migrate**
- **Make**

## Структура проекта

```
smart-home-controller/
├── api/                # OpenAPI спецификация
├── cmd/server/         # Главный исполняемый файл
├── internal/
│   ├── gateways/       # Слой взаимодействия с HTTP, WS
│   ├── repository/     # Слой доступа к данным
│   ├── domain/         # Модели
│   └── usecase/        # Бизнес-логика
├── migrations/         # Миграции для БД
├── pkg/pg_test/        # Инструменты для тестирования PostgreSQL
├── .github/workflows/  # Конфигурация CI/CD
├── docker-compose.yml  # Конфигурация Docker Compose
├── Makefile            # Запуск миграций
```

## Запуск проекта

### Требования

- [Go](https://golang.org/dl/) 1.16+
- [Docker](https://www.docker.com/get-started), [Docker Compose](https://docs.docker.com/compose/install/)

### Инструкция

```bash
git clone https://github.com/javaddikkt/smart-home-controller.git
cd smart-home-controller
docker-compose up -d
make migrate-up
export DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/db?sslmode=disable"
go run ./cmd/server
```

Приложение доступно по адресу: `http://localhost:8080`

## API

Полная спецификация лежит в [`api/swagger.yaml`](api/swagger.yaml).  
Пример запроса:

```http
POST /events
Content-Type: application/json
{
    "sensor_serial_number": "1234567890",
    "payload": 10
}
```

## Тестирование

Запуск тестов:

```bash
go test -v ./... -race
```
