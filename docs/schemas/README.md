# Схемы БД сервисов ChainWise

Документация схем БД для реализации миграций и репозиториев (STEP #8, #10, #11, #5, #13). Схемы согласованы с доменной моделью и контрактами (contracts/proto).


| Сервис               | Файл                                                 | Назначение                                 |
| -------------------- | ---------------------------------------------------- | ------------------------------------------ |
| Order Service        | [order-service.md](./order-service.md)               | Заказы, позиции, outbox, QR-токены         |
| IAM Service          | [iam-service.md](./iam-service.md)                   | Пользователи, права, сессии (Redis)        |
| Inventory Service    | [inventory-service.md](./inventory-service.md)       | Локации, остатки, резервы, движения        |
| Accounting Service   | [accounting-service.md](./accounting-service.md)     | Поставщики, закупки, проводки склад↔филиал |
| Delivery Service     | [delivery-service.md](./delivery-service.md)         | Назначения заказ↔курьер, витрина           |
| Notification Service | [notification-service.md](./notification-service.md) | Опционально: подписки для WebSocket        |


