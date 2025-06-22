# Cart Service

Сервис корзины для e-commerce платформы с интеграцией LOMS (Logistics and Order Management System).

## 🚀 Быстрый старт

### Предварительные требования
- Go 1.23.4+
- Docker и Docker Compose
- PostgreSQL (для LOMS сервиса)

### Установка и запуск

1. **Клонирование репозитория**
```bash
git clone <repository-url>
cd cart5
```

2. **Установка зависимостей**
```bash
go mod download
```

3. **Запуск с Docker Compose**
```bash
cd ../loms4_fix
docker-compose up -d
cd ../cart5
make docker-build
make docker-run
```

4. **Запуск локально**
```bash
make build-cart
make run-cart
```

## 📋 API Endpoints

### Cart Service (HTTP)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/user/{user_id}/cart/{sku_id}` | Добавить товар в корзину |
| `DELETE` | `/user/{user_id}/cart/{sku_id}` | Удалить товар из корзины |
| `DELETE` | `/user/{user_id}/cart` | Очистить корзину |
| `GET` | `/user/{user_id}/cart` | Получить содержимое корзины |
| `POST` | `/cart/checkout` | Оформить заказ |
| `GET` | `/health` | Проверка здоровья сервиса |
| `GET` | `/ready` | Проверка готовности сервиса |

### LOMS Service (gRPC)

| Метод | Описание |
|-------|----------|
| `OrderCreate` | Создание нового заказа |
| `OrderInfo` | Получение информации о заказе |
| `OrderPay` | Оплата заказа |
| `OrderCancel` | Отмена заказа |
| `StocksInfo` | Информация о наличии товаров |

## ⚙️ Конфигурация

### Cart Service (`config.yaml`)
```yaml
product_client:
  url: "http://route256.pavl.uk:8080/get_product"
  token: "testtoken"
  timeout: "30s"
  retry_attempts: 3
  retry_delay: "1s"

cart_server:
  port: "8082"
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

loms_server:
  gRPCport: "50051"

rate_limit:
  requests_per_minute: 100
```

### LOMS Service (`config.yaml`)
```yaml
loms_server:
  gRPCport: "50051"

database:
  host: "postgres"
  port: "5432"
  user: "root"
  password: "root"
  dbname: "loms_db"
  sslmode: "disable"
```

## 🛠️ Разработка

### Команды Makefile

```bash
# Сборка и запуск
make build-cart          # Сборка сервиса
make run-cart           # Запуск сервиса
make run-all            # Сборка и запуск

# Тестирование
make test               # Запуск всех тестов
make test-race          # Тесты с проверкой race conditions
make test-coverage      # Тесты с покрытием кода
make test-bench         # Бенчмарки

# Анализ кода
make lint               # Линтинг
make fmt                # Форматирование кода
make vet                # Веттинг кода
make check              # Полная проверка

# Анализ сложности
make cognitive-load     # Когнитивная нагрузка
make cyclomatic-load    # Цикломатическая сложность

# Docker
make docker-build       # Сборка Docker образа
make docker-run         # Запуск Docker контейнера

# Утилиты
make generate-mocks     # Генерация моков
make deps-check         # Проверка зависимостей
make clean              # Очистка
```

### Структура проекта

```
cart5/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа
├── internal/
│   ├── app/
│   │   └── cart/
│   │       ├── service.go       # Бизнес-логика
│   │       ├── service_test.go  # Тесты сервиса
│   │       └── mock/            # Моки для тестов
│   ├── client/
│   │   └── client.go            # HTTP клиент
│   ├── config/
│   │   └── config.go            # Конфигурация
│   ├── delivery/
│   │   ├── http_handler.go      # HTTP обработчики
│   │   ├── router.go            # Роутинг
│   │   ├── response.go          # Утилиты ответов
│   │   └── health.go            # Health check
│   ├── domain/
│   │   └── cart.go              # Доменные модели
│   ├── mw/
│   │   ├── logging.go           # Логирование
│   │   └── rate_limit.go        # Rate limiting
│   ├── repository/
│   │   ├── cart.go              # Репозиторий
│   │   └── cart_test.go         # Тесты репозитория
│   └── localErr/
│       └── err.go               # Локальные ошибки
├── e2e/
│   └── cart_test.go             # E2E тесты
├── config.yaml                  # Конфигурация
├── Dockerfile                   # Docker образ
├── Makefile                     # Команды сборки
└── README.md                    # Документация
```

## 🔧 Особенности реализации

### Безопасность
- **Rate Limiting**: Ограничение количества запросов в минуту
- **Валидация**: Проверка входных данных на всех уровнях
- **Санитизация логов**: Фильтрация чувствительных данных

### Надежность
- **Graceful Shutdown**: Корректное завершение работы сервисов
- **Retry Logic**: Повторные попытки для внешних вызовов
- **Health Checks**: Мониторинг состояния сервисов
- **Error Handling**: Централизованная обработка ошибок

### Производительность
- **Connection Pooling**: Переиспользование соединений
- **In-Memory Storage**: Быстрый доступ к данным корзины
- **Concurrent Processing**: Параллельная обработка запросов

### Мониторинг
- **Structured Logging**: Структурированные логи
- **Metrics**: Метрики производительности
- **Tracing**: Трассировка запросов

## 🧪 Тестирование

### Unit тесты
```bash
make test
```

### Интеграционные тесты
```bash
make test-coverage
```

### E2E тесты
```bash
cd e2e
go test -v
```

### Бенчмарки
```bash
make test-bench
```

## 📊 Метрики и мониторинг

### Health Check
```bash
curl http://localhost:8082/health
```

### Readiness Check
```bash
curl http://localhost:8082/ready
```

## 🐛 Отладка

### Логи
Логи выводятся в stdout в структурированном формате:
```
REQUEST: method=POST, url=/user/123/cart/456, body={"count":2}
RESPONSE: method=POST, status=200, body={"message":"Item added to cart successfully"}
```

### Отладка в Docker
```bash
docker logs cart-service
```

## 🔄 CI/CD

### GitHub Actions
Проект включает настройки для автоматической сборки и тестирования в GitHub Actions.

### Docker
```bash
# Сборка образа
docker build -t cart-service .

# Запуск контейнера
docker run -p 8082:8082 cart-service
```

## 📝 Лицензия

MIT License

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📞 Поддержка

Для вопросов и предложений создавайте Issues в репозитории.
