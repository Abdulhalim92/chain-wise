# Схемы БД сервисов ChainWise

Документация схем БД для реализации миграций и репозиториев. Схемы согласованы с доменной моделью и контрактами (`contracts/proto`).

---

## Сервисы и файлы

| Сервис | Файл | Назначение |
|--------|------|------------|
| **Order Service** | [order-service.md](order-service.md) | Заказы, позиции, outbox, QR-токены |
| **IAM Service** | [iam-service.md](iam-service.md) | Пользователи, права, сессии (Redis) |
| **Inventory Service** | [inventory-service.md](inventory-service.md) | Локации, остатки, резервы, движения |
| **Accounting Service** | [accounting-service.md](accounting-service.md) | Поставщики, закупки, проводки склад↔филиал |
| **Delivery Service** | [delivery-service.md](delivery-service.md) | Назначения заказ↔курьер, витрина |
| **Notification Service** | [notification-service.md](notification-service.md) | Подписки для WebSocket (опционально) |

---

## Использование

- Ссылки на шаги плана: **STEP #3** (проектирование), **STEP #6, #9, #11, #12, #14** (миграции и репозитории).
- При изменении схем обновлять соответствующий `.md` и миграции в каталоге `migrations/` сервиса.
