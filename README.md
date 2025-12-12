## О проекте

В репозитории опубликована **небольшая часть backend-кода на Golang**, вынесенная и обобщённая на основе моего реального рабочего опыта.

Код размещён **исключительно в демонстрационных целях** — для ознакомления со стилем моей разработки, архитектурными подходами и уровнем владения Go.  
Проект **не является самостоятельным production-решением / не был взят из production кода / и не претендует на завершённый продукт ** 

---

## Описание функциональности
Код моделирует backend-логику CRM-подобной системы, ориентированной на управление сущностями и ролями в рамках организации (например, учебного центра или внутреннего корпоративного сервиса). 
В системе реализована работа с основными доменными сущностями: 
- пользователи с разными ролями и уровнями доступа; 
- управляемые сущности (например, сотрудники, преподаватели, клиенты / учащиеся); 
- связи между сущностями и агрегированные запросы.
** Backend предоставляет REST API для: 
- аутентификации и авторизации пользователей; 
- разграничения доступа к эндпоинтам в зависимости от роли; 
- управления сущностями (CRUD-операции); 
- выполнения типовых бизнес-запросов (фильтрация, выборки, связанные данные). 
- Архитектура и структура кода отражают подход, применимый к реальным CRM-системам и внутренним бизнес-приложениям, где важны контроль доступа, безопасность, предсказуемость API и поддерживаемость кода.

## Используемый стек

- Go (Golang)
- REST API (net/http)
- SQL (MySQL / MariaDB)
- JWT-аутентификация и авторизация
- Middleware-архитектура
- Конфигурация через переменные окружения

---

## Структура проекта

- `cmd/api` — точка входа приложения
- `internal/` — основная бизнес-логика (handlers, repositories, middlewares)
- `pkg/` — переиспользуемые утилиты
- `.env.example` — шаблон переменных окружения

---


## Запуск проекта (опционально)

Проект не требует обязательного запуска для ознакомления с кодом, однако при необходимости его можно запустить локально следующими командами:

```bash
cp .env.example .env && \
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout cmd/api/key.pem \
  -out cmd/api/cert.pem \
  -subj "/CN=localhost" && \
go run ./cmd/api
```


Маршруты и доступ по роля:
В проекте реализована ролевая модель доступа (RBAC). Доступ к API-эндпоинтам ограничивается ролью пользователя.


Public routes (без проверки роли)
```bash
POST /execs/login

POST /execs/logout

POST /execs/forgotpassword

POST /execs/resetpassword/reset/{resetcode}
```
Admin-only routes
```bash
DELETE /execs/{id}

DELETE /teachers/{id}

DELETE /teachers

DELETE /students

DELETE /students/{id}
```
Admin & Manager routes
```bash
GET /execs

GET /execs/{id}

POST /execs

PATCH /execs

PATCH /execs/{id}

POST /teachers

PATCH /teachers

PATCH /teachers/{id}

PUT /teachers/{id}
```
Admin, Manager & Exec routes
```bash
GET /students

GET /students/{id}

POST /students

PATCH /students

PATCH /students/{id}

PUT /students/{id}

GET /teachers

GET /teachers/{id}

GET /teachers/{id}/students

GET /teachers/{id}/studentcount

POST /execs/{id}/updatepassword
```
Примечание: проверка ролей выполняется на уровне middleware до выполнения бизнес-логики хендлеров.
