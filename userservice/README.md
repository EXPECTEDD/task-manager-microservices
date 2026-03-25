# Userservice

## О сервисе

`userservice` - это микросервис для работы с пользователями.  
Он отвечает за:

- Регистрацию пользователей  
- Логин и аутентификацию  
- Управление сессиями через UUID, хранящиеся в Redis и передаваемые через cookies  
- Авторизацию и проверку прав доступа к ресурсам (через сессии)  

Сервис предоставляет **REST API** для взаимодействия с клиентом и **gRPC API** для взаимодействия с другими сервисами.

## Стек технологий

### Основной
- Go  
- PostgreSQL  
- Redis  
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
- uuid - генерация уникальных id сессий

## Конфигурация

Сервис использует:

- `.env` для переменных окружения  
- YAML-файлы (`config/local.yaml`, `config/docker.yaml`) для настройки приложения и подключения к БД/Redis  

Пример основных переменных в `.env`:
```text
DB_PASS=user_pass
REDIS_PASS=redis_pass

MIG_NAME=migrator
MIG_PASS=migrator_pass
DB_HOST=localhost
DB_PORT=1111
DB_NAME=users
DB_MODE=disable
```

## API
### REST API

Основные endpoints:

#### POST /user/registration - регистрация нового пользователя

**Пример запроса:**

```text
POST /user/registration
Content-Type: application/json

{
    "first_name": "Иван",
    "middle_name": "Иванович",
    "last_name": "Иванов",
    "password": "123",
    "email": "ivan@gmail.com"
}
```

**Пример ответа:**
```json
{
    "user_id": 1
}
```

**Пример ошибок:**

- HTTP код 409 - Conflict - user already exists
- HTTP код 400 - Bad Request - массив errors с перечеслением пропущенных или невалидных полей (FirstName, LastName, Password, Email) | bad request body
- HTTP код 500 - Internal Server Error - internal server error

#### POST /user/login - вход пользователя

**Пример запроса:**

```text
POST /user/login
Content-Type: application/json

{
	"email": "ivan@gmail.com",
    "password": "123"
}
```

**Пример ответа:**
```json
{
    "user": {
        "first_name": "Иван",
        "middle_name": "Иванович",
        "last_name": "Иванов"
    }
}
```

**Пример ошибок:**

- HTTP код 404 - Not found - user not found
- HTTP код 400 - Bad Request - массив errors с перечеслением пропущенных или невалидных полей (Email, Password) | bad request body
- HTTP код 401 - Unauthorized - wrong password 
- HTTP код 500 - Internal Server Error - internal server error

### gRPC API

Сервис предоставляет gRPC методы для проверки сессий и аутентификации пользователей для других сервисов (projectservice, taskservice).

Proto файлы находятся в proto/userservice/user.proto.

## Тестирование

### Юнит-тесты

Можно запускать локально без поднятых контейнеров:

```bash
cd userservice
go test ./...
```

### Интеграционные тесты

Интеграционные тесты требуют запущенных сервисов через Docker Compose, так как они взаимодействуют с базой данных и другими сервисами:

```bash
docker compose up --build
cd userservice
go test ./...
```

Важно: текущая конфигурация использует ephemeral базы данных. При перезапуске контейнеров все данные сбрасываются.

## Структура проекта

```text
userservice/
├─ cmd/                  # Точка входа в сервис
├─ config/               # YAML конфиги
├─ internal/             # Основной код сервиса
│  ├─ app/               # Инициализация приложения
|  ├─ config/            # Загрузка логов и переменных окружения
│  ├─ domain/            # Доменные модели
│  ├─ infrastructure/    # Взаимодействие с БД, Redis, bcrypt, uuid
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