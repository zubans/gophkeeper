# Docker Setup для GophKeeper

## Обзор

Docker Compose конфигурация для запуска GophKeeper с PostgreSQL базой данных, сервером, клиентом и pgAdmin для управления базой данных.

## Быстрый старт

### 1. Запуск всех сервисов
```bash
# Сборка и запуск всех сервисов
make docker-up

# Или напрямую
docker-compose up -d
```

### 2. Проверка статуса
```bash
# Просмотр логов
make docker-logs

# Или
docker-compose logs -f
```

### 3. Остановка сервисов
```bash
# Остановка всех сервисов
make docker-down

# Или
docker-compose down
```

## Сервисы

### PostgreSQL Database
- **Порт**: 5432
- **База данных**: gophkeeper
- **Пользователь**: gophkeeper
- **Пароль**: password
- **Данные**: сохраняются в volume `postgres_data`

### GophKeeper Server
- **Порт**: 8080
- **URL**: http://localhost:8080
- **Автоматически применяет миграции при запуске**
- **Зависит от PostgreSQL**

### GophKeeper Client
- **Интерактивный режим**
- **Подключается к серверу через внутреннюю сеть**
- **Локальная SQLite база для кэширования**

### pgAdmin (Web UI для PostgreSQL)
- **Порт**: 5050
- **URL**: http://localhost:5050
- **Email**: admin@gophkeeper.local
- **Пароль**: admin

## Режимы работы

### Production режим
```bash
# Полный запуск всех сервисов
make docker-up
```

### Development режим
```bash
# Запуск только базы данных и pgAdmin
make docker-dev

# Затем запуск сервера локально
make run-server
```

## Команды Docker

### Основные команды
```bash
# Сборка образов
make docker-build

# Запуск сервисов
make docker-up

# Остановка сервисов
make docker-down

# Просмотр логов
make docker-logs

# Очистка ресурсов
make docker-clean
```

### Прямые команды Docker Compose
```bash
# Запуск в фоне
docker-compose up -d

# Запуск с пересборкой
docker-compose up --build

# Остановка и удаление volumes
docker-compose down -v

# Просмотр статуса
docker-compose ps

# Выполнение команд в контейнере
docker-compose exec server /bin/sh
docker-compose exec postgres psql -U gophkeeper -d gophkeeper
```

## Структура файлов

```
.
├── docker-compose.yml          # Production конфигурация
├── docker-compose.dev.yml      # Development конфигурация
├── Dockerfile.server           # Dockerfile для сервера
├── Dockerfile.client           # Dockerfile для клиента
├── .dockerignore              # Исключения для Docker
└── scripts/
    └── init-db.sql            # Скрипт инициализации БД
```

## Переменные окружения

### Server
- `DB_HOST` - хост базы данных (postgres)
- `DB_PORT` - порт базы данных (5432)
- `DB_USER` - пользователь БД (gophkeeper)
- `DB_PASSWORD` - пароль БД (password)
- `DB_NAME` - имя БД (gophkeeper)
- `JWT_SECRET` - секретный ключ для JWT
- `ENCRYPTION_KEY` - ключ шифрования

### Client
- `SERVER_URL` - URL сервера (http://server:8080)

## Volumes

- `postgres_data` - данные PostgreSQL
- `pgadmin_data` - настройки pgAdmin
- `go_mod_cache` - кэш Go модулей (dev режим)

## Сети

- `gophkeeper-network` - внутренняя сеть для production
- `gophkeeper-dev-network` - внутренняя сеть для development

## Troubleshooting

### Проблемы с подключением к БД
```bash
# Проверка статуса PostgreSQL
docker-compose exec postgres pg_isready -U gophkeeper -d gophkeeper

# Подключение к БД
docker-compose exec postgres psql -U gophkeeper -d gophkeeper
```

### Проблемы с миграциями
```bash
# Просмотр логов сервера
docker-compose logs server

# Ручной запуск миграций
docker-compose exec server ./gophkeeper-server -migrate-only
```

### Очистка и перезапуск
```bash
# Полная очистка
make docker-clean

# Пересборка и запуск
make docker-build
make docker-up
```

## Мониторинг

### Проверка здоровья сервисов
```bash
# Статус всех контейнеров
docker-compose ps

# Логи конкретного сервиса
docker-compose logs server
docker-compose logs postgres
```

### Подключение к контейнерам
```bash
# Shell в сервер
docker-compose exec server /bin/sh

# Shell в PostgreSQL
docker-compose exec postgres /bin/sh

# psql в PostgreSQL
docker-compose exec postgres psql -U gophkeeper -d gophkeeper
```
