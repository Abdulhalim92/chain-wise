# Inventory Service — схема БД

Остатки по локациям, асинхронный резерв по событию `order.created`, ручное списание; gRPC от Accounting для прихода при закупке.

## Таблицы

### locations

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор локации |
| type | VARCHAR, NOT NULL | warehouse (ЦС) | branch (филиал) |
| external_id | VARCHAR, NULL | Ссылка на branch_id или код склада |
| created_at | TIMESTAMPTZ | Время создания |

Индексы: `type`, `external_id`.

### stock

Остатки по локации и номенклатуре.

| Колонка | Тип | Описание |
|--------|-----|----------|
| location_id | UUID (FK → locations), NOT NULL | Локация |
| product_id | UUID/VARCHAR, NOT NULL | Номенклатура |
| quantity | INT, NOT NULL, >= 0 | Доступное количество |
| reserved | INT, NOT NULL, >= 0 | Зарезервировано |
| updated_at | TIMESTAMPTZ | Время обновления |

PK: (location_id, product_id). Индексы: `location_id`, `product_id`. Ограничение: quantity >= reserved.

### reservations

Резерв по заказу: по одной строке на каждую номенклатуру заказа. Идемпотентность по (order_id, product_id).

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор резерва |
| order_id | UUID, NOT NULL | Заказ |
| location_id | UUID (FK → locations), NOT NULL | Локация (обычно ЦС) |
| product_id | UUID/VARCHAR, NOT NULL | Номенклатура |
| quantity | INT, NOT NULL, > 0 | Зарезервировано |
| created_at | TIMESTAMPTZ | Время резерва |

UNIQUE(order_id, product_id). Индексы: `order_id`, `location_id`, `product_id`.

### stock_movements (опционально, для аудита операций прихода/списания)

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор операции |
| location_id | UUID (FK → locations), NOT NULL | Локация |
| product_id | UUID/VARCHAR, NOT NULL | Номенклатура |
| quantity_delta | INT, NOT NULL | Приход (>0) или списание (<0) |
| reason | VARCHAR | purchase, reserve, cancel_reserve, adjustment, transfer |
| reference_id | UUID/VARCHAR, NULL | order_id, purchase_id и т.д. |
| created_at | TIMESTAMPTZ | Время операции |

Индексы: `location_id`, `product_id`, `created_at`, `reference_id`.

Связь с контрактами: `StockItem`, `IncreaseStockRequest`/`GetStockRequest`, `CancelReservationRequest` (inventory.proto).
