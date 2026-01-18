# Simple Broker

## О проекте
Легковесный брокер сообщений на Go для маршрутизации и батчинга событий между микросервисами.

## Основные возможности
- ✅ Прием HTTP-событий по группам
- ✅ Батчинг событий с настраиваемым размером и таймаутом
- ✅ Маршрутизация на multiple endpoints
- ✅ Асинхронная обработка без блокировок
- ✅ Graceful shutdown
- ✅ Поддержка JSON формата

## Архитектура
```
HTTP → Router → Input Channel → Dispatch Queues → Output Channel → HTTP Client
      ↑                                     ↓
    Config ←-------------------------- Group Configuration
```

## Быстрый старт

### 1. Установка
```bash
git clone <repository>
cd simple_brocker
go mod download
```

### 2. Конфигурация
Создайте `config/config.yaml`:
```yaml
address: "0.0.0.0:8080"

group:
  users:
    address: ["http://user-service:8081", "http://backup-user:8082"]
    cooldown: 1s
    batch_size: 10
  
  payments:
    address: ["http://payment-service:8083"]
    cooldown: 500ms
    batch_size: 5
```

### 3. Запуск
```bash
go run cmd/app/main.go
```

## Использование

### Отправка события
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"user_id": 123, "action": "login"}'
```

### Группы событий
- `/users` → events отправляются в user-service
- `/payments` → events отправляются в payment-service

## Конфигурация групп

| Параметр | Описание | Пример |
|----------|----------|---------|
| `address` | Список endpoint'ов для отправки | `["service:8080", "backup:8081"]` |
| `cooldown` | Таймаут батчинга | `1s`, `500ms` |
| `batch_size` | Максимальный размер батча | `10`, `50`, `100` |

## Разработка

### Структура проекта
```
simple_brocker/
├── cmd/app/main.go          # Точка входа
├── config/                  # Конфигурация
├── internal/
│   ├── server/             # HTTP сервер
│   ├── web/               # Роутинг и обработчики
│   ├── service/           # Бизнес-логика
│   │   ├── event/         # Модель события
│   │   ├── queue/         # Очередь и батчинг
│   │   └── dispatch/      # Маршрутизация по группам
│   └── config/            # Загрузка конфигурации
└── README.md
```
