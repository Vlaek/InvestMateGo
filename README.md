Структура проекта:

```
    cmd/
        server/
            main.go             # Точка входа, HTTP сервер
    config/
        config.go               # Конфигурация, загрузка переменных окружения
    internal/
        api/
            dto/                # Data Transfer Objects (структуры API Tinkoff)
            mappers/            # Мапперы DTO → Domain модели
        models/                 # Domain модели бизнес-логики
        storage/                # Хранилище данных (кэш инструментов)
    pkg/
        logger/                 # Логирование
    .env                        # Переменные окружения (создать!)
    .env.example                # Пример конфигурации
    .gitignore
    go.mod                      # Зависимости Go
    go.sum
    README.md
```

Установка зависимостей:
```
    git clone https://github.com/Vlaek/InvestMateGo.git

    cd InvestMateGo

    go mod download
```

Настройка окружения:
```
    cp .env.example .env

    # Отредактируйте .env файл
    # добавьте ваш токен Tinkoff OpenAPI
```

Запуск сервера:
```
    go run cmd/server/main.go
```

Доступные эндпоинты:
| Эндпоинт  | Метод | Описание |
| ------------- | ------------- | ------------- |
| /  | GET  | Информация о сервере  |
| /health  | GET  | Проверка состояния  |
| /config  | GET  | Текущая конфигурация  |
| /bonds  | GET  | Список всех облигаций  |