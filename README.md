# ChainWise

ERP-платформа для автоматизации цепочки поставок: внешняя закупка на центральный склад и внутренние поставки со склада в филиалы (точки). Операции выполняются с планшетов; интерфейс зависит от роли; списки заказов обновляются в реальном времени по WebSocket.

## Возможности

- **Внешняя закупка** — учёт поступлений от поставщиков, накладные, цены, увеличение остатков на ЦС.
- **Внутренние поставки** — заявки филиалов, резервирование, сборка, доставка курьером, подтверждение прихода по QR.
- **Роли:** оператор точки (создание заказа, назначение курьера), кладовщик (сборка, QR, учёт закупок), курьер (забор/доставка по QR), администратор.

## Стек

| Компонент | Технология |
|-----------|------------|
| Backend | Go 1.26 |
| Межсервисное API | gRPC (Protocol Buffers v3) |
| Внешний API | REST/JSON через API Gateway |
| Очереди | RabbitMQ (topic `order_events`) |
| БД | PostgreSQL 16 (отдельная на сервис) |
| Сессии | Redis (IAM) |
| Real-time | WebSocket / SSE |

## Структура репозитория

```
├── contracts/          # Контракты gRPC (.proto), генерация Go-кода (buf)
├── gateway/            # API Gateway — единая точка входа
├── services/
│   ├── iam-service/        # Аутентификация, JWT, роли, права
│   ├── order-service/      # Заказы, State Machine, QR, события
│   ├── inventory-service/   # Остатки, резерв по событиям
│   ├── accounting-service/ # Закупки, проводки склад↔филиал
│   ├── delivery-service/   # Назначение курьера, подтверждение забора/доставки
│   └── notification-service/ # WebSocket/SSE по событиям заказов
├── docs/
│   └── schemas/        # Схемы БД сервисов (документация)
├── go.work             # Go workspace
├── ChainWise-Implementation-Plan.md
└── ChainWise-Project-Description.md
```

## Требования

- Go 1.26+
- [buf](https://buf.build/docs/installation) — для генерации кода из `.proto`
- PostgreSQL 16, RabbitMQ, Redis (для локального запуска сервисов)

## Быстрый старт

1. Клонировать репозиторий и перейти в каталог:
   ```bash
   git clone https://github.com/Abdulhalim92/chain-wise.git
   cd chain-wise
   ```

2. Сгенерировать Go-код из контрактов (нужен установленный buf):
   ```bash
   cd contracts && make generate && cd ..
   ```

3. Собрать и запустить нужный сервис, например:
   ```bash
   go build -o bin/gateway ./gateway/cmd/gateway
   ./bin/gateway
   ```
   Или из каталога сервиса:
   ```bash
   cd services/iam-service && go build -o ../../bin/iam-service ./cmd/iam-service && ../../bin/iam-service
   ```

### Конфигурация (platform)

Конфигурация загружается через **Viper** без дефолтов: сначала ищется файл `.env`, при отсутствии — `config.yml` (путь к YAML задаётся через `CONFIG_FILE`). Если ни один файл не найден, сервис завершается с ошибкой. Переменные окружения (PORT, GRPC_PORT, ENV, LOG_LEVEL, LOG_FORMAT, LOG_ADD_SOURCE) переопределяют значения из файла.

- **PORT**, **GRPC_PORT** — порты HTTP и gRPC (по умолчанию 8080, 9090).
- **ENV**, **LOG_LEVEL** — окружение и уровень логов (info, debug, warn, error).
- **LOG_FORMAT** — `json` (по умолчанию) или `text`.
- **LOG_ADD_SOURCE** — добавлять в логи файл и строку вызова (true/false).

Пример `config.yml`: см. `platform/config/config.example.yml`.

**Где лежат конфиги:** файлы `.env` и `config.yml` ищутся в **рабочем каталоге процесса**. Для запуска из корня репо положите туда `.env` или `config.yml`; при запуске из каталога сервиса — в каталог сервиса. Путь к YAML (если не используете `.env`) задаётся через `CONFIG_FILE`.

### Go workspace (go.work, go.work.sum)

- **go.work** — описание workspace: список модулей (contracts, platform, gateway, services/...) и версия Go. Нужен, чтобы из корня репо собирать и запускать все модули без `replace` в каждом go.mod (зависимости между модулями разрешаются через workspace).
- **go.work.sum** — хеши модулей, на которые ссылается go.work (аналог go.sum для workspace). Обычно коммитится в репозиторий.

**Как с ними работать:** не редактировать go.work вручную без необходимости. Добавить новый модуль в workspace: `go work use ./path/to/module`. Синхронизировать зависимости после изменений в go.mod: из корня выполнить `go work sync` (обновит go.work.sum и при необходимости go.work). Сборка из корня: `go build ./platform/... ./gateway/... ./services/iam-service/...` и т.д.

### go mod tidy: когда и где

- **Из корня репо:** `go work sync` обновляет go.work.sum; зависимости каждого модуля при этом не пересчитываются. Чтобы подтянуть/очистить зависимости **всех** модулей, можно выполнить `go mod tidy` в каждом каталоге модуля (см. ниже).
- **Для каждого модуля:** зайти в каталог модуля (`platform`, `gateway`, `services/iam-service`, …) и выполнить `go mod tidy`. Так обновляются go.mod и go.sum этого модуля. Удобно делать после добавления импортов или смены версий. В монорепозитории с go.work достаточно выполнить tidy в тех модулях, которые вы меняли; при сборке из корня Go подставит локальные модули из workspace.

### Логирование (platform)

Логгер строится на **slog** (Go 1.21+): структурированные логи с уровнем, форматом (JSON/текст), опциональным указанием источника (файл:строка) и именем сервиса. Поддержка дочерних логгеров: `WithService(name)`, `With(args...)` для добавления постоянных полей. Все HTTP-запросы и gRPC-вызовы проходят через middleware/interceptors с логированием метода, пути, статуса и длительности.

## Документация

- [План реализации](ChainWise-Implementation-Plan.md) — архитектура и шаги (STEP #1–#26).
- [Описание проекта](ChainWise-Project-Description.md) — сценарии, цепочки запросов, тестовые кейсы.
- [Схемы БД](docs/schemas/README.md) — таблицы и поля для Order, IAM, Inventory, Accounting, Delivery, Notification.

## Лицензия

Проприетарный проект.
