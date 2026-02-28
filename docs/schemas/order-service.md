# Order Service — схема БД

**Назначение:** жизненный цикл заказа, State Machine, генерация/валидация QR, публикация в RabbitMQ `order_events`.

---

## Таблицы

### orders

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор заказа |
| branch_id | UUID, NOT NULL | Филиал (точка), для которого заказ |
| status | VARCHAR, NOT NULL | created, reserved, in_progress, ready_for_pickup, in_transit, delivered, cancelled |
| created_at | TIMESTAMPTZ, NOT NULL | Время создания |
| updated_at | TIMESTAMPTZ | Время последнего обновления |
| courier_id | UUID, NULL | Курьер (если назначен) |

Индексы: `branch_id`, `status`, `courier_id`, `created_at`.

### order_items

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор позиции |
| order_id | UUID (FK → orders), NOT NULL | Заказ |
| product_id | UUID/VARCHAR, NOT NULL | Номенклатура |
| quantity | INT, NOT NULL, > 0 | Количество |
| unit_price_cents | BIGINT, NOT NULL | Цена за единицу на момент заказа, копейки (для проводок и отображения) |

Индексы: `order_id`.

### order_outbox (опционально, для Transactional Outbox)

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор записи |
| event_type | VARCHAR, NOT NULL | order.created, order.reserved, order.in_progress, order.ready_for_pickup, order.in_transit, order.delivered, order.cancelled |
| payload | JSONB, NOT NULL | order_id, branch_id, и др. данные события |
| created_at | TIMESTAMPTZ, NOT NULL | Время создания |
| processed_at | TIMESTAMPTZ, NULL | Время публикации в RabbitMQ (NULL = ещё не отправлено) |

Индексы: `processed_at` (для выборки необработанных).

### qr_tokens (для одноразовых QR при заборе/доставке)

Один заказ — два токена: pickup и delivery. PK составной.

| Колонка | Тип | Описание |
|--------|-----|----------|
| order_id | UUID (FK → orders), NOT NULL | Заказ |
| action | VARCHAR, NOT NULL | pickup | delivery |
| token_hash | VARCHAR, NOT NULL | Хеш одноразового токена |
| used_at | TIMESTAMPTZ, NULL | Время использования (NULL = не использован) |
| created_at | TIMESTAMPTZ, NOT NULL | Время создания токена |

PK: (order_id, action). Индексы: `order_id`, `action`.

Связь с контрактами: `Order`, `OrderItem` (product_id, quantity, unit_price_cents), `OrderStatus` (orders.proto).
