# Accounting Service — схема БД

Закупки, накладные, проводки «Склад → Филиал»; подписка на `order.delivered` для закрытия проводок; gRPC к Inventory для увеличения остатков при закупке.

## Таблицы

### suppliers (справочник поставщиков)

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор поставщика |
| name | VARCHAR, NOT NULL | Наименование |
| created_at | TIMESTAMPTZ | Время создания |

### purchases

Закупки (внешняя закупка, регистрация кладовщиком).

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор закупки |
| supplier_id | UUID (FK → suppliers), NOT NULL | Поставщик |
| document_ref | VARCHAR, NOT NULL | Номер накладной/документа |
| total_cents | BIGINT, NOT NULL | Сумма в копейках |
| created_at | TIMESTAMPTZ, NOT NULL | Время создания |

Индексы: `supplier_id`, `created_at`, `document_ref`.

### purchase_items

Позиции закупки.

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор позиции |
| purchase_id | UUID (FK → purchases), NOT NULL | Закупка |
| product_id | UUID/VARCHAR, NOT NULL | Номенклатура |
| quantity | INT, NOT NULL, > 0 | Количество |
| price_cents | BIGINT, NOT NULL | Цена за единицу, копейки |

Индексы: `purchase_id`.

### postings

Проводки «Склад → Филиал» по внутренним поставкам. Создаются при потреблении события `order.delivered` (идемпотентность по order_id).

| Колонка | Тип | Описание |
|--------|-----|----------|
| id | UUID (PK) | Идентификатор проводки |
| order_id | UUID, NOT NULL, UNIQUE | Заказ (запрет двойных проводок по одному заказу) |
| from_location_id | UUID/VARCHAR, NOT NULL | Склад (ЦС) |
| to_location_id | UUID/VARCHAR, NOT NULL | Филиал (branch_id) |
| amount_cents | BIGINT, NOT NULL | Сумма/объём (при необходимости) |
| created_at | TIMESTAMPTZ, NOT NULL | Время создания |

Индексы: `order_id`, `from_location_id`, `to_location_id`, `created_at`.

Связь с контрактами: `Purchase`, `PurchaseItem`, `RegisterPurchaseRequest` (accounting.proto).
