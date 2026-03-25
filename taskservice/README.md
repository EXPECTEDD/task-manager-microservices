# Taskservice

## О сервисе

`taskservice` - это микросервис для работы с задачами.  
Он отвечает за:

- Создание задач  
- Удаление задач  
- Получение задач
- Изменение задач  

Сервис предоставляет **REST API** для взаимодействия с клиентом.

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

#### POST /task/create/:project_id - создание новой задачи

**Пример запроса:**

```text
POST /task/create/1
Content-Type: application/json

{
    "description": "buy milk",
    "deadline": "2026-03-20T16:30:00Z"
}
```

**Пример ответа:**
```json
{
    "task_id": 1
}
```

**Пример ошибок:**

- HTTP код 504 - Gateway Timeout - project service timeout
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - ivalid project id (возвращает projectservice)
- HTTP код 404 - Not Found - project not found
- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 400 - Bad Request - invalid project id (возвращает проверка в taskservice)
- HTTP код 400 - Bad Request - массив errors с перечеслением пропущенных или невалидных полей (Description) | bad request body
- HTTP код 400 - Bad Request - invalid description
- HTTP код 500 - Internal Server Error - internal server error

#### PATCH /task/update/:task_id/:project_id - обновление задачи

**Пример запроса:**

```text
PATCH /task/update/1/1
Content-Type: application/json

{
    "new_description": "buy milk",
    "new_deadline": "2011-03-20T16:30:00Z"
}
```

**Пример ответа:**
```json
{
    "updated": true
}
```

**Пример ошибок:**

- HTTP код 504 - Gateway Timeout - project service timeout
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - ivalid project id (возвращает projectservice)
- HTTP код 404 - Not Found - project not found
- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 400 - Bad Request - invalid project id (возвращает проверка в taskservice)
- HTTP код 400 - Bad Request - invalid task id
- HTTP код 400 - Bad Request - nothing to update
- HTTP код 400 - Bad Request - invalid new description
- HTTP код 404 - Not Found - task not found
- HTTP код 500 - Internal Server Error - internal server error

#### DELETE /task/delete/:task_id/:project_id - удаление задачи

**Пример запроса:**

```text
DELETE /task/delete/1/1
Content-Type: application/json
```

**Пример ответа:**
```json
{
    "deleted": true
}
```

**Пример ошибок:**

- HTTP код 504 - Gateway Timeout - project service timeout
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - ivalid project id (возвращает projectservice)
- HTTP код 404 - Not Found - project not found
- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 400 - Bad Request - invalid project id (возвращает проверка в taskservice)
- HTTP код 400 - Bad Request - invalid task id
- HTTP код 404 - Not Found - task not found
- HTTP код 500 - Internal Server Error - internal server error

#### GET /task/getall/:project_id - получение всех задач проекта

**Пример запроса:**

```text
GET /task/getall/1
Content-Type: application/json
```

**Пример ответа:**
```json
{
    "tasks": [
        {
            "Id": 1,
            "ProjectId": 1,
            "Description": "buy milk",
            "Deadline": "2026-03-20T16:30:00Z"
        }
    ]
}
```

**Пример ошибок:**

- HTTP код 504 - Gateway Timeout - project service timeout
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - ivalid project id (возвращает projectservice)
- HTTP код 404 - Not Found - project not found
- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 400 - Bad Request - invalid project id (возвращает проверка в taskservice)
- HTTP код 404 - Not Found - task not found
- HTTP код 500 - Internal Server Error - internal server error

#### GET /task/get/:task_id/:project_id - получение определенной задачи

**Пример запроса:**

```text
GET /task/get/1/1
Content-Type: application/json
```

**Пример ответа:**
```json
{
    "task": {
        "Id": 1,
        "ProjectId": 1,
        "Description": "buy milk",
        "Deadline": "2026-03-20T16:30:00Z"
    }
}
```

**Пример ошибок:**

- HTTP код 504 - Gateway Timeout - project service timeout
- HTTP код 504 - Gateway Timeout - user service timeout
- HTTP код 400 - Bad Request - ivalid project id (возвращает projectservice)
- HTTP код 404 - Not Found - project not found
- HTTP код 404 - Not Found - user not found
- HTTP код 502 - Bad Gateway - upstream error
- HTTP код 400 - Bad Request - invalid project id (возвращает проверка в taskservice)
- HTTP код 400 - Bad Request - invalid task id
- HTTP код 404 - Not Found - task not found
- HTTP код 500 - Internal Server Error - internal server error

## gRPC

Proto файлы находятся в proto/taskservice/user.proto.

## Тестирование

### Юнит-тесты

Можно запускать локально без поднятых контейнеров:

```bash
cd taskservice
go test ./...
```

### Интеграционные тесты

Интеграционные тесты требуют запущенных сервисов через Docker Compose, так как они взаимодействуют с базой данных и другими сервисами:

```bash
docker compose up --build
cd taskservice
go test ./...
```

Важно: текущая конфигурация использует ephemeral базы данных. При перезапуске контейнеров все данные сбрасываются.

## Структура проекта

```text
taskservice/
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