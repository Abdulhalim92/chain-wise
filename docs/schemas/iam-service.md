# IAM Service — схема БД

**Назначение:** учётные записи, роли, JWT, проверка прав (в т.ч. право «создание заказа» на точке).  
**Хранилища:** PostgreSQL + Redis (сессии).

---

## Таблицы (PostgreSQL)

### users

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор пользователя |
| login | VARCHAR, UNIQUE, NOT NULL | Логин |
| password_hash | VARCHAR, NOT NULL | Хеш пароля |
| role | VARCHAR, NOT NULL | warehouse, branch_operator, courier, admin |
| branch_id | UUID, NULL | Филиал (для оператора точки/курьера) |
| created_at | TIMESTAMPTZ, NOT NULL | Время создания |
| updated_at | TIMESTAMPTZ | Время обновления |

Индексы: `login`, `role`, `branch_id`.

### permissions

Связь роль/пользователь — право. Право «создание заказа» выдаётся одному сотруднику точки.

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор |
| user_id | UUID (FK → users), NULL | Пользователь (если право персональное) |
| role | VARCHAR, NULL | Роль (если право по роли) |
| permission_key | VARCHAR, NOT NULL | create_order, assign_courier, … |
| branch_id | UUID, NULL | Ограничение по филиалу (для create_order — только свой) |

Индексы: `user_id`, `role`, `permission_key`, `branch_id`.

### user_roles (опционально, если роли хранятся отдельно)

| Колонка | Тип | Описание |
|--------|-----|----------|
| user_id | UUID (FK → users), PK | Пользователь |
| role | VARCHAR, NOT NULL | Роль |

## Сессии (Redis)

- Ключ: по session_id или token_id (из JWT jti).
- Значение: user_id, role, branch_id, courier_id (если курьер), permissions (при необходимости).
- TTL: время жизни сессии/токена.

Связь с контрактами: `LoginRequest`/`LoginResponse`, `ValidateTokenRequest`/`ValidateTokenResponse` (auth.proto): user_id, role, branch_id, courier_id, permissions.
