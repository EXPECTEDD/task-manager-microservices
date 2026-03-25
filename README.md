<h1 align="center">
  <a href="https://github.com/EXPECTEDD/task-manager-microservices" target="_blank">task-manager-microservices</a>
</h1>

<h3 align="center">Микросервисное приложение для управления проектами и задачами.</h3>

## О проекте

Task Manager — это микросервисное приложение для управления проектами и задачами.  
Приложение построено на микросервисной архитектуре и состоит из отдельных сервисов для аутентификации, управления проектами и задачами.  

- Взаимодействие с клиентом осуществляется через REST API.  
- Взаимодействие между микросервисами реализовано с использованием gRPC.  
- Аутентификация реализована через серверные сессии (UUID), хранящиеся в Redis и передаваемые через cookies.  
- Авторизация основана на проверке прав доступа к ресурсам (владение проектами и задачами).  

Проект разработан в целях практики и закрепления навыков разработки микросервисных backend-приложений.

## Реализованные сервисы

- [userservice](./userservice/README.md) — сервис аутентификации и авторизации (регистрация, логин, управление сессиями).
- [projectservice](./projectservice/README.md) — сервис управления проектами.
- [taskservice](./taskservice/README.md) — сервис управления задачами.

## Стек технологий

- Go  
- PostgreSQL  
- Redis  
- REST API (клиент - сервис)  
- gRPC (сервис - сервис)  
- Docker, Docker Compose  

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/EXPECTEDD/task-manager-microservices
cd task-manager-microservices
```

### 2. Настройка переменных окружения

Создайте файл .env в корне проекта и заполните переменные окружения, такие как:
```text
#userservice
US_DB_USER=user
US_DB_PASS=user_pass
US_DB_NAME=users
US_DB_HOST=userservice_postgres
US_DB_PORT=5432
US_DB_MODE=disable

US_REDIS_HOST=redis
US_REDIS_PORT=1111
US_REDIS_PASS=redis_pass
US_REDIS_DB=0
US_REDIS_TTL_SEC=3600

#projectservice
PS_DB_USER=user
PS_DB_PASS=user_pass
PS_DB_HOST=projectservice_postgres
PS_DB_PORT=2222
PS_DB_NAME=projects
PS_DB_MODE=disable

#taskservice
TS_DB_USER=user
TS_DB_PASS=user_pass
TS_DB_HOST=taskservice_postgres
TS_DB_PORT=3333
TS_DB_NAME=tasks
TS_DB_MODE=disable
```

### 3. Сборка и запуск
```bash
docker compose up --build
```

После запуска сервисы будут доступны на соответствующих портах.

## Структура проекта
```text
task-manager-microservices/
├─ projectservice/     # Сервис управления проектами
├─ taskservice/        # Сервис управления задачами
├─ userservice/        # Сервис аутентификации и авторизации
├─ .env
├─ docker-compose.yaml
├─ LICENSE
└─ README.md
```

Подробная структура каждого сервиса указана в README соответствующего сервиса.

## Основные возможности
- Регистрация и логин пользователей
- Создание, обновление и удаление проектов
- Создание, обновление и удаление задач
- Привязка задач к проектам
- Проверка прав доступа пользователей к проектам и задачам

## Тестирование

Каждый сервис содержит юнит- и интеграционные тесты.

### Юнит-тесты

Юнит-тесты можно запускать локально без поднятых контейнеров:

```bash
cd service_folder
go test ./...
```

### Интеграционные тесты

Интеграционные тесты требуют запущенных сервисов через Docker Compose, так как они взаимодействуют с базой данных и другими сервисами:

```bash
docker compose up -build
cd service_folder
go test ./...
```

Важно: текущая конфигурация использует ephemeral базы даephemeral базы данныхнных. При перезапуске контейнеров все данные сбрасываются.