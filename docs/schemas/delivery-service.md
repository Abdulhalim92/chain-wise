# Delivery Service — схема БД

**Назначение:** назначение заказа на курьера (pull/push), подтверждение забора и доставки; gRPC к Order (ValidateQR, TransitionStatus).

---

## Таблицы

### assignments

Назначение заказ ↔ курьер.

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор назначения |
| order_id | UUID, NOT NULL, UNIQUE | Заказ (один заказ — один курьер) |
| courier_id | UUID/VARCHAR, NOT NULL | Курьер |
| status | VARCHAR, NOT NULL | assigned, picked_up, delivered |
| assigned_at | TIMESTAMPTZ, NOT NULL | Время назначения |
| picked_up_at | TIMESTAMPTZ, NULL | Время забора (QR/кнопка) |
| delivered_at | TIMESTAMPTZ, NULL | Время доставки на филиал |
| created_at | TIMESTAMPTZ | Время создания записи |
| updated_at | TIMESTAMPTZ | Время обновления |

Индексы: `order_id`, `courier_id`, `status`, `assigned_at`.

### available_orders (опционально, витрина для pull)

Кэш заказов в статусе ready_for_pickup, доступных для Claim. Может заполняться по событию от Order или по запросу при открытии витрины курьером.

| Колонка | Тип | Описание |
|--------|-----|----------|
| order_id | UUID (PK) | Заказ |
| branch_id | UUID, NOT NULL | Филиал |
| created_at | TIMESTAMPTZ | Когда заказ стал доступен |
| expires_at | TIMESTAMPTZ, NULL | Опционально TTL для кэша |

Альтернатива: запрос списка заказов в статусе ready_for_pickup к Order Service по gRPC без локальной таблицы.

Связь с контрактами: `DeliveryAssignment`, `ClaimOrderRequest`/`AssignOrderRequest`, `ConfirmPickupRequest`/`ConfirmDeliveryRequest` (delivery.proto).
