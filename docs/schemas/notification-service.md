# Notification Service — схема БД

Постоянная БД не обязательна. Сервис подписан на RabbitMQ (notification_events / order_events), рассылает обновления по WebSocket/SSE по роли, branch_id, courier_id.

## Опциональная схема (при необходимости)

Если нужна персистентность подписок (например, для переподключения или маршрутизации):

### subscriptions

Подписки клиентов для маршрутизации событий.

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор подписки |
| client_id | VARCHAR, NOT NULL | Идентификатор WebSocket-клиента/сессии |
| user_id | UUID, NULL | Пользователь (из JWT) |
| role | VARCHAR, NULL | warehouse, branch_operator, courier, admin |
| branch_id | UUID, NULL | Филиал (для фильтрации событий по заказам точки) |
| courier_id | UUID/VARCHAR, NULL | Курьер (для фильтрации «мои заказы») |
| connected_at | TIMESTAMPTZ | Время подключения |
| disconnected_at | TIMESTAMPTZ, NULL | Время отключения |

Индексы: `client_id`, `role`, `branch_id`, `courier_id`.

В минимальном варианте подписки хранятся только в памяти (map по client_id с контекстом role/branch_id/courier_id); при перезапуске сервиса клиенты переподключаются по WebSocket.

Связь с контрактами: события заказов (OrderEventType, OrderEvent в notifications.proto); маршрутизация по контексту из JWT (role, branch_id, courier_id).
