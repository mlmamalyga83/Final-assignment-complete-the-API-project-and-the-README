# Todo API Service

REST API для управления задачами (Todos) и пользователями (Users), реализованное на Go с использованием фреймворка chi и Swagger-документацией.

## Возможности

- **Управление пользователями** — регистрация, просмотр, обновление и удаление пользователей
- **Управление задачами** — создание, просмотр, редактирование и удаление задач
- **Swagger UI** — интерактивная документация по адресу `/swagger/`

## Зависимости

- **Go** версии 1.22+
- **chi** v5.0.10 — веб-фреймворк
- **swaggo/http-swagger** v2.0.2 — Swagger UI
- **swaggo/swag** v1.16.2 — генерация Swagger-документации

## Установка и запуск

### 1. Установите Go

Скачайте и установите Go с [официального сайта](https://go.dev/).

### 2. Установите swaggo CLI (для генерации документации)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 3. Клонируйте репозиторий и перейдите в корень проекта

```bash
git clone <repository-url>
cd api-doc-example
```

### 4. Сгенерируйте Swagger-документацию

```bash
swag init -g main.go -o docs
```

> **Примечание:** Если `swag` не установлен, можно использовать готовый файл `docs/swagger.json`, который уже идёт в комплекте.

### 5. Запустите сервер

```bash
go run main.go
```

Сервер запустится на `http://localhost:8080`.

### 6. Откройте Swagger UI

Перейдите по адресу: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

## Структура проекта

```
api-doc-example/
├── main.go                        # Точка входа, роутинг, Swagger UI
├── internal/
│   ├── handlers/
│   │   ├── todo.go                # CRUD-обработчики для задач
│   │   └── user.go                # CRUD-обработчики для пользователей
│   └── models/
│       ├── todo.go                # Модели Todo, CreateTodoRequest, UpdateTodoRequest
│       └── user.go                # Модели User, CreateUserRequest, UpdateUserRequest, ErrorResponse, SuccessResponse
├── docs/
│   └── swagger.json               # Сгенерированная Swagger-спецификация
├── go.mod                         # Модульные зависимости
└── go.sum                         # Контрольная сумма зависимостей
```

## API

### Эндпоинты

#### Пользователи (`/api/v1/users`)

| Метод   | Путь              | Описание              |
|---------|-------------------|-----------------------|
| `GET`   | `/api/v1/users`   | Список всех пользователей |
| `GET`   | `/api/v1/users/{id}` | Получить пользователя по ID |
| `POST`  | `/api/v1/users`   | Создать нового пользователя |
| `PUT`   | `/api/v1/users/{id}` | Обновить пользователя |
| `DELETE`| `/api/v1/users/{id}` | Удалить пользователя |

#### Задачи (`/api/v1/todos`)

| Метод   | Путь              | Описание              |
|---------|-------------------|-----------------------|
| `GET`   | `/api/v1/todos`   | Список всех задач     |
| `GET`   | `/api/v1/todos/{id}` | Получить задачу по ID |
| `POST`  | `/api/v1/todos`   | Создать новую задачу  |
| `PUT`   | `/api/v1/todos/{id}` | Обновить задачу      |
| `DELETE`| `/api/v1/todos/{id}` | Удалить задачу       |

> Полная интерактивная документация доступна в [Swagger UI](http://localhost:8080/swagger/).

## Примеры запросов

### Пользователи

#### Создать пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "name": "John Doe"}'
```

**Ответ (201):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "created_at": "2026-10-01T15:00:00Z"
  }
}
```

#### Получить список пользователей

```bash
curl http://localhost:8080/api/v1/users
```

**Ответ (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "email": "john@example.com",
      "name": "John Doe",
      "created_at": "2026-10-01T15:00:00Z"
    }
  ]
}
```

#### Получить пользователя по ID

```bash
curl http://localhost:8080/api/v1/users/1
```

**Ответ (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "created_at": "2026-10-01T15:00:00Z"
  }
}
```

#### Обновить пользователя

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'
```

**Ответ (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "Jane Doe",
    "created_at": "2026-10-01T15:00:00Z"
  }
}
```

#### Удалить пользователя

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

**Ответ: 204 No Content**

#### Ошибка: пользователь не найден

```bash
curl http://localhost:8080/api/v1/users/999
```

**Ответ (404):**

```json
{
  "error": "404",
  "message": "Пользователь не найден"
}
```

### Задачи

#### Создать задачу

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "title": "Buy groceries", "description": "Milk, bread, eggs", "done": false}'
```

**Ответ (201):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Buy groceries",
    "description": "Milk, bread, eggs",
    "done": false,
    "created_at": "2026-10-01T15:00:00Z",
    "updated_at": "2026-10-01T15:00:00Z"
  }
}
```

#### Получить список задач

```bash
curl http://localhost:8080/api/v1/todos
```

**Ответ (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "title": "Buy groceries",
      "description": "Milk, bread, eggs",
      "done": false,
      "created_at": "2026-10-01T15:00:00Z",
      "updated_at": "2026-10-01T15:00:00Z"
    }
  ]
}
```

#### Получить задачу по ID

```bash
curl http://localhost:8080/api/v1/todos/1
```

**Ответ (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Buy groceries",
    "description": "Milk, bread, eggs",
    "done": false,
    "created_at": "2026-10-01T15:00:00Z",
    "updated_at": "2026-10-01T15:00:00Z"
  }
}
```

#### Обновить задачу

```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"done": true}'
```

**Ответ (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Buy groceries",
    "description": "Milk, bread, eggs",
    "done": true,
    "created_at": "2026-10-01T15:00:00Z",
    "updated_at": "2026-10-01T15:00:00Z"
  }
}
```

#### Удалить задачу

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

**Ответ: 204 No Content**

#### Ошибка: задача не найдена

```bash
curl http://localhost:8080/api/v1/todos/999
```

**Ответ (404):**

```json
{
  "error": "404",
  "message": "Задача не найдена"
}
```

## Чек-лист

- [x] Go 1.22+
- [x] Swagger-аннотации ко всем endpoint'ам (10 шт.)
- [x] Docstrings к моделям данных
- [x] Сгенерирован `docs/swagger.json`
- [x] Swagger UI доступен по адресу `/swagger/`
- [x] README.md содержит все разделы
- [x] Примеры запросов и ответов
- [x] Проект запускается без ошибок

## Лицензия

MIT

## Контакты

По вопросам обращайтесь: [support@netology.ru](mailto:support@netology.ru)