# Тестовое задание Avito 2025

[Задача](./task.md)

## Запуск проекта

1. Клонируйте репозиторий:

```bash
git clone https://github.com/danyarmarkin/AvitoTask2025.git
```

2. Перейдите в директорию проекта:

```bash
cd AvitoTask2025
```

3. Запустите Docker контейнер:

```bash
docker-compose up -d
```

## Технологии
1. Кодогеренация: OpenAPI Generator https://github.com/oapi-codegen/oapi-codegen
    - Конфигурация генерации находится в файле [oapi-codegen.yaml](./oapi-codegen.yaml)
    - Генерация кода выполняется командой:
    ```bash
    make generate
    ```
2. Миграции: Goose https://github.com/pressly/goose
    - Миграции находятся в директории [db/migrations](./db/migrations)