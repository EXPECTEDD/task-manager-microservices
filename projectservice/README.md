# Projectservice

## О сервисе

`projectservice` - это микросервис для работы с проектами.  
Он отвечает за:

- Создание проекта
- Удаление проекта
- Получение всех проектов
- Изменение проекта

Сервис предоставляет **REST API** для взаимодействия с клиентом и **gRPC API** для взаимодействия с другими сервисами.

## Стек технологий

### Основной
- Go  
- PostgreSQL 
- REST API (клиент - userservice)  
- gRPC (сервис - сервис)  
- Docker, Docker Compose  

### Инструменты разработки
- Gin - HTTP framework  
- gRPC-Go - межсервисная коммуникация  
- Testify - тестирование  
- mockgen - генерация моков  
- slog - логирование  
- migrate / SQL миграции - управление схемой базы данных

## Конфигурация

Сервис использует:

- `.env` для переменных окружения  
- YAML-файлы (`config/local.yaml`, `config/docker.yaml`) для настройки приложения и подключения к БД/Redis  

Пример основных переменных в `.env`:
```text
DB_PASS=user_pass

MIG_NAME=migrator
MIG_PASS=migrator_pass
DB_HOST=localhost
DB_PORT=1111
DB_NAME=tasks
DB_MODE=disable
```

## API
### REST API

Основные endpoints:

#### POST /project/create - создание нового проекта

**Пример запроса:**

```text
POST /project/create
Content-Type: application/json

{
    "name": "IvanProject"
}
```

**Пример ответа:**
```json
{
    "project_id": 1
}
```

**Пример ошибок:**

- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 409 - Conflict - user already exists
- HTTP код 400 - Bad Request - массив errors с перечеслением пропущенных или невалидных полей (Name) | bad request body
- HTTP код 400 - Bad Request - invalid name
- HTTP код 400 - Bad Request - invalid owner id
- HTTP код 409 - Conflict - project already exists
- HTTP код 500 - Internal Server Error - internal server error

#### DELETE /project/delete/:project_id - удаление проекта

**Пример запроса:**

```text
DELETE /project/delete/1
Content-Type: application/json
```

**Пример ответа:**
```json
{
    "is_deleted": true
}
```

**Пример ошибок:**

- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - invalid project_id type
- HTTP код 400 - Bad Request - invalid project_id value
- HTTP код 400 - Bad Request - invalid project id
- HTTP код 404 - Not Found - project not found
- HTTP код 500 - Internal Server Error - internal server error

#### GET /project/getall - получение всех проектов

**Пример запроса:**

```text
GET /project/getall
Content-Type: application/json
```

**Пример ответа:**
```json
{
    "projects": [
        {
            "Id": 1,
            "OwnerId": 1,
            "Name": "IvanProject",
            "CreatedAt": "2026-03-25T17:35:38.50702Z"
        }
    ]
}
```

**Пример ошибок:**

- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 404 - Not Found - projects not found
- HTTP код 500 - Internal Server Error - internal server error

#### PATCH /project/update/:project_id - обновление проекта

**Пример запроса:**

```text
PATCH /project/update/1
Content-Type: application/json

{
    "new_name":"SuperBestProject"
}
```

**Пример ответа:**
```json
{
    "is_updated": true
}
```

**Пример ошибок:**

- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - invalid project_id type
- HTTP код 400 - Bad Request - invalid project_id value
- HTTP код 400 - Bad Request - массив errors с перечеслением пропущенных или невалидных полей (NewName) | bad request body
- HTTP код 400 - Bad Request - invalid new project name
- HTTP код 404 - Not Found - project not found
- HTTP код 409 - Conflict - project name already exists
- HTTP код 500 - Internal Server Error - internal server error

### gRPC API

Сервис предоставляет gRPC метод для получения user_id проекта (для taskservice).

Proto файлы находятся в proto/projectservice/user.proto.

## Тестирование

### Юнит-тесты

Можно запускать локально без поднятых контейнеров:

```bash
cd projectservice
go test ./...
```

### Интеграционные тесты

Интеграционные тесты требуют запущенных сервисов через Docker Compose, так как они взаимодействуют с базой данных и другими сервисами:

```bash
docker compose up --build
cd projectservice
go test ./...
```

Важно: текущая конфигурация использует ephemeral базы данных. При перезапуске контейнеров все данные сбрасываются.

## Структура проекта

```text
projectservice/
├─ cmd/                  # Точка входа в сервис
├─ config/               # YAML конфиги
├─ internal/             # Основной код сервиса
│  ├─ app/               # Инициализация приложения
|  ├─ config/            # Загрузка логов и переменных окружения
│  ├─ domain/            # Доменные модели
│  ├─ infrastructure/    # Взаимодействие с БД и другими сервисами
│  ├─ repository/        # Репозитории и хранилища
│  ├─ transport/         # REST и gRPC обработчики
│  └─ usecase/           # Бизнес-логика
├─ migrations/           # SQL миграции
├─ pkg/logger/           # Логирование
├─ proto/                # gRPC proto файлы
├─ tests/                # Интеграционные тесты
├─ .env
├─ Dockerfile
├─ go.mod
├─ go.sum
├─ Makefile
└─ README.md
```

## Makefile и миграции

Сервис содержит Makefile с основными командами для разработки и работы с БД.

### Локальный запуск
```bash
make local
```

Запускает сервис с переменными из .env и конфигурацией config/local.yaml.

### Применение миграций
```bash
make migrate_all_up
```
Применяет все миграции к базе данных.

Важно: перед запуском убедитесь, что переменные подключения к БД заданы в .env.

### Генерация gRPC кода
```bash
make build_proto
```
Генерирует .pb.go файлы и gRPC код из .proto файлов.