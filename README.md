## Структура проекта:

```
    project/
    ├── cmd/
    │   └── server/
    │       └── main.go          # Точка входа
    ├── internal/
    │   ├── api/                 # Внешнее API (Tinkoff)
    │   ├── config/              # Конфигурация
    │   ├── handlers/            # Обработчики
    │   ├── mappers/             # Преобразователи между моделями
    │   ├── models/              # Модели
    |   |   ├── domain/          # Доменные модели
    |   |   ├── dto/             # DTO внешнего API
    |   |   └── entity/          # Сущности для БД (GORM)
    │   ├── repository/          # Работа с данными
    │   ├── services/            # Бизнес-логика
    │   └── storage/             # Хранилище
    ├── pkg/                     # Переиспользуемый код
    ├── scripts/                 # Скрипты
    ├── migrations/              # Миграции БД
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