# zl0y.team Billing Backend

## Описание

**zl0y.team Billing** — микросервис для MVP личного кабинета zl0y.team.  
Реализует регистрацию, аутентификацию, работу с балансом пользователя, хранение пользователей в PostgreSQL и отчетов в MongoDB.  
Сервис позволяет анонимно создавать отчеты, регистрироваться, привязывать анонимные отчеты к пользователю, а также покупать доступ к отчетам.

## Основные возможности

- Регистрация и аутентификация пользователей (JWT)
- Хранение пользователей в PostgreSQL
- Хранение и покупка отчетов в MongoDB
- Привязка анонимных отчетов к пользователю
- Мок-эндпоинт для создания анонимных отчетов (для тестирования)
- Защищённые эндпоинты (требуют JWT)

## Архитектура

- **Go (Gin)** — HTTP API
- **PostgreSQL** — пользователи
- **MongoDB** — отчеты
- **JWT** — аутентификация

### Индексы MongoDB

- `report_id` — уникальный индекс (поиск по отчету)
- `user_id` — индекс для быстрого поиска отчетов пользователя
- `client_generated_id` — индекс для привязки анонимных отчетов

## Инструкция по запуску

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/Aza-di/Zloy.git
   cd Zloy
   ```

2. Запустите сервисы:
   ```sh
   docker-compose up --build
   ```

3. API будет доступен на [http://localhost:8080](http://localhost:8080)

## Примеры curl-запросов

### Регистрация пользователя

```sh
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"testpass"}'
```

### Логин

```sh
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"testuser","password":"testpass"}'
```

### Создание анонимного отчета

```sh
curl -X POST http://localhost:8080/api/mock/create-report \
  -H "Content-Type: application/json" \
  -d '{"client_generated_id":"anon-123"}'
```

### Привязка анонимного отчета

```sh
curl -X POST http://localhost:8080/api/user/link-anonymous \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"client_generated_id":"anon-123"}'
```

### Получение отчетов пользователя

```sh
curl -X GET "http://localhost:8080/api/user/reports?limit=10&offset=0" \
  -H "Authorization: Bearer <access_token>"
```

### Покупка отчета

```sh
curl -X POST http://localhost:8080/api/reports/<report_id>/purchase \
  -H "Authorization: Bearer <access_token>"
```

---

## История коммитов

Сохранена в git.

---

## Контакты

Автор: [Aza-di](https://github.com/Aza-di) 