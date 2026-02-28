# Архитектурный обзор ChainWise

> Краткая проверка кодовой базы с точки зрения лучших практик (Clean Architecture, DDD, production readiness).

---

## 1. Критерии оценки

| Принцип | Описание | Статус в проекте |
|--------|----------|------------------|
| **Dependency Rule** | Зависимости направлены внутрь: domain не зависит от инфраструктуры | ✅ Интерфейсы репозиториев в `domain`, реализация в `repository` |
| **Single Source of Truth** | Доменные коды и маппинги в одном месте | ✅ `contracts/codes`, маппинг в `platform/errors` с использованием констант |
| **Explicit over Implicit** | Конфиг без скрытых дефолтов, явная валидация | ✅ `config.Load()` возвращает ошибку, проверка обязательных полей и диапазонов |
| **Observability** | Request ID, структурированные логи, envelope при ошибках | ✅ Middleware, slog, envelope в Gateway (включая panic) |
| **Fail-safe** | Ошибки возвращаются, не паника; единый формат ответов | ✅ REST envelope, gRPC status + details |
| **Testability** | Границы через интерфейсы, моки (gomock) | ✅ IAM app тесты с gomock, изоляция по слоям |

---

## 2. Направление зависимостей

```
cmd → internal/grpc, internal/app → internal/domain
                ↓
internal/repository → internal/domain
                ↑
        (реализует интерфейсы)

Сервисы → contracts (proto, codes)
Сервисы, Gateway → platform (config, logger, errors, middleware)
```

- **domain** не импортирует БД, HTTP, MQ.
- **app** зависит только от интерфейсов из domain.
- **contracts** и **platform** — общие модули без обратных зависимостей на сервисы.

---

## 3. Рекомендации (уже учтённые)

- Envelope для **всех** REST-ответов Gateway (включая `/`, `/health`, panic).
- Доменные коды в **contracts**; маппинг в **platform/errors** через константы (один источник правды).
- Валидация конфига: порты 1–65535, `log_level` из фиксированного набора.
- Отдельная ошибка **ErrMisconfigured** при nil-репозитории (отличие от ErrInvalidCredentials).
- Обёртка ошибок `issue jwt: %w` для трассировки.
- Моки через **gomock** (mockgen), не ручные заглушки.

---

## 4. Ссылки на план

- Структура сервисов и зависимостей: [ChainWise-Implementation-Plan.md](../ChainWise-Implementation-Plan.md) § 1.3, 1.4.
- Envelope и доменные коды: § 1.6, STEP #5.
- Схемы БД: [docs/schemas/README.md](schemas/README.md).
