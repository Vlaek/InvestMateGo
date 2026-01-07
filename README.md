## Структура проекта:

```
    invest-mate/
    ├── cmd/
    │   └── server/
    │       └── main.go               # Точка входа
    ├── internal/
    │    ├── assets/                  # Модуль активов
    │    │   ├── api/                   # Внешнее API
    │    │   |   └── tinkoff              # Tinkoff API для активов
    │    │   ├── handlers/              # Обработчики
    │    │   ├── mappers/               # Преобразователи между моделями
    │    │   ├── models/                # Модели
    │    |   |   ├── domain/              # Доменные модели
    │    |   |   ├── dto/                 # DTO внешнего API
    │    |   |   └── entity/              # Сущности для БД (GORM)
    │    │   ├── repository/            # Репозиторий активов
    │    │   ├── services/              # Бизнес-логика активов
    │    │   └── storage/               # Хранилище активов
    │    ├── shared/                  # Общий модуль переиспользуемого функционала
    │    │   ├── api/                   # Внешнее API
    │    │   |   └── tinkoff              # Общий Tinkoff API клиент
    │    │   └── config/                # Конфигурация
    │    └── users/                   # Модуль пользователей
    │        ├── api/                   # Внешнее API
    │        |   └── tinkoff              # Tinkoff API для получения портфеля
    │        ├── handlers/              # Обработчики
    │        ├── mappers/               # Преобразователи между моделями
    │        ├── models/                # Модели
    │        |   ├── domain/              # Доменные модели
    │        |   ├── dto/                 # DTO внешнего API
    │        |   └── entity/              # Сущности для БД (GORM)
    │        ├── repository/            # Репозиторий
    │        └── services/              # Бизнес-логика
    ├── pkg/                          # Переиспользуемый код
    ├── scripts/                      # Скрипты
    ├── migrations/                   # Миграции БД
    └── go.mod
```

### Установка зависимостей:
```bash
    git clone https://github.com/Vlaek/InvestMateGo.git

    cd InvestMateGo

    go mod download
```

### Настройка окружения:
```bash
    cp .env.example .env

    # Отредактируйте .env файл
    # добавьте ваш токен Tinkoff OpenAPI
```

### Запуск сервера:
```bash
    go run cmd/server/main.go
```

## Доступные эндпоинты:
| Эндпоинт  | Метод | Описание |
| ------------- | ------------- | ------------- |
| /  | GET  | Информация о сервере  |
| /health  | GET  | Проверка состояния  |
| /config  | GET  | Текущая конфигурация  |
| /bonds  | GET  | Список всех облигаций  |
| /shares  | GET  | Список всех акций  |
| /etfs  | GET  | Список всех фондов  |
| /currencies  | GET  | Список всех валют  |
