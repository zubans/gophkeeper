# GophKeeper

GophKeeper - это безопасная клиент-серверная система для хранения приватных данных, построенная на Go с использованием принципов чистой архитектуры. Позволяет пользователям надёжно хранить и синхронизировать пароли, текстовые данные, бинарные файлы и информацию о банковских картах между несколькими устройствами.

## Возможности

### Сервер
- Регистрация, аутентификация и авторизация пользователей
- Безопасное хранение данных с шифрованием AES-256-GCM
- Синхронизация данных между несколькими клиентами с разрешением конфликтов
- RESTful API для взаимодействия с клиентами
- JWT-аутентификация с 24-часовым сроком действия
- Поддержка PostgreSQL с автоматическими миграциями
- Контроль версий (последние 10 версий для каждого элемента)

### Клиент
- Кроссплатформенное CLI-приложение (Windows, Linux, macOS)
- **Чистая архитектура** с разделением ответственности и dependency injection
- Безопасная аутентификация с удалённым сервером
- Поддержка различных типов данных:
  - Пары логин/пароль (`login_password`)
  - Произвольные текстовые данные (`text`)
  - Бинарные данные (`binary`)
  - Информация о банковских картах (`bank_card`)
- Локальное кэширование в SQLite для офлайн работы
- Автоматическая синхронизация с сервером
- Отображение информации о версии и дате сборки

## Архитектура

### Чистая архитектура клиента

Клиент построен по принципам чистой архитектуры с четким разделением слоев:

```
┌─────────────────────────────────────────┐
│              Presentation               │
│            (CLI Commands)               │
├─────────────────────────────────────────┤
│              Application                │
│              (Client)                   │
├─────────────────────────────────────────┤
│               Domain                    │
│    (AuthService, DataService,           │
│     SyncService + Interfaces)           │
├─────────────────────────────────────────┤
│            Infrastructure               │
│  (HTTPClient, Storage, TokenManager,    │
│           Encryptor)                    │
└─────────────────────────────────────────┘
```

#### Основные компоненты:
- **Интерфейсы**: Storage, HTTPClient, Encryptor, TokenManager, AuthService, DataService, SyncService
- **Сервисы**: AuthServiceImpl, DataServiceImpl, SyncServiceImpl
- **Инфраструктура**: HTTPClientImpl, ClientStorage, TokenManagerImpl
- **Приложение**: Client с dependency injection

### База данных

#### Сервер (PostgreSQL)
- `users` - пользователи с хешированными паролями
- `stored_data` - основные данные с полями `is_deleted`, `version`
- `data_history` - история версий (последние 10)
- `schema_migrations` - управление миграциями

#### Клиент (SQLite)
- Аналогичные таблицы для локального кэширования
- Поддержка офлайн работы
- Автоматическая синхронизация при подключении

## Система синхронизации

### Разрешение конфликтов
1. **По времени**: Если `updated_at` клиента > сервера → клиент побеждает
2. **По версии**: Если время одинаковое, сравнивается `version`
3. **Конфликты**: Если сервер новее, данные помечаются как конфликт

### Алгоритм синхронизации
1. Клиент отправляет данные, измененные с последней синхронизации
2. Сервер сравнивает с локальными данными
3. Применяет правила разрешения конфликтов
4. Возвращает обновленные данные сервера
5. Клиент сохраняет полученные данные локально

## Установка

### Предварительные требования
- Go 1.23.6 или новее
- PostgreSQL 12 или новее (для сервера)
- Git

### Быстрый старт с Docker

```bash
# Клонирование репозитория
git clone <repository-url>
cd Gophkeeper

# Запуск всех сервисов
make docker-up

# Проверка статуса
make docker-logs
```

### Сборка из исходного кода

```bash
# Загрузка зависимостей
make deps

# Сборка приложения
make build

# Сборка для всех платформ
make build-all
```

## Использование

### Запуск сервера

```bash
# Используя make
make run-server

# Или напрямую
go run ./cmd/server

# С пользовательскими настройками
go run ./cmd/server -db-host=localhost -db-port=5432 -db-user=gophkeeper -db-password=password -db-name=gophkeeper
```

### Использование клиента

```bash
# Регистрация нового пользователя
./bin/gophkeeper-client register username email@example.com password

# Вход в систему
./bin/gophkeeper-client login username password

# Добавление данных логин/пароль
./bin/gophkeeper-client add login_password "My Website" "username" "password" "https://example.com" "Additional notes"

# Добавление текстовых данных
./bin/gophkeeper-client add text "Important Note" "This is my important note"

# Добавление данных банковской карты
./bin/gophkeeper-client add bank_card "My Credit Card" "1234567890123456" "12/25" "123" "John Doe" "Bank Name" "Additional notes"

# Просмотр всех данных
./bin/gophkeeper-client list

# Получение конкретных данных
./bin/gophkeeper-client get <data-id>

# Удаление данных
./bin/gophkeeper-client delete <data-id>

# Синхронизация с сервером
./bin/gophkeeper-client sync

# Просмотр истории версий
./bin/gophkeeper-client history <data-id>

# Просмотр информации о версии
./bin/gophkeeper-client version
```

## Docker

### Сервисы
- **PostgreSQL Database** - Порт 5432
- **GophKeeper Server** - Порт 8080
- **GophKeeper Client** - Интерактивный режим
- **pgAdmin** - Порт 5050 (admin@gophkeeper.local / admin)

### Команды Docker

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

## Безопасность

### Шифрование данных
- **Алгоритм**: AES-256-GCM
- **Ключ**: 32-байтовый ключ шифрования
- **Режим**: GCM для аутентификации и шифрования

### Хеширование паролей
- **Алгоритм**: SHA-256
- **Соль**: 32-байтовая случайная соль

### JWT токены
- **Алгоритм**: HMAC-SHA256
- **Срок действия**: 24 часа
- **Поля**: user_id, username, exp, iat

## Тестирование

Проект имеет хорошо организованную структуру тестов с разделением по типам:

### Структура тестов
```
internal/client/tests/
├── mocks/
│   └── mocks.go              # Все моки для тестирования
├── auth_service_test.go      # Тесты сервиса аутентификации
├── data_service_test.go      # Тесты сервиса данных
├── sync_service_test.go      # Тесты сервиса синхронизации
└── integration_test.go       # Интеграционные тесты
```

### Запуск тестов

```bash
# Все тесты
make test

# Тесты с покрытием
make test-coverage

# Тесты конкретного модуля
go test ./internal/client/tests -v

# Интеграционные тесты
go test ./tests -v
```

### Типы тестов
1. **Unit тесты** - Тестирование отдельных сервисов с моками
2. **Integration тесты** - Тестирование взаимодействия компонентов
3. **End-to-end тесты** - Полное тестирование пользовательских сценариев

## API Endpoints

### Аутентификация
- `POST /api/v1/register` - Регистрация нового пользователя
- `POST /api/v1/login` - Аутентификация пользователя

### Управление данными
- `GET /api/v1/data` - Получение всех данных пользователя
- `POST /api/v1/data` - Создание новых данных
- `PUT /api/v1/data` - Обновление существующих данных
- `DELETE /api/v1/data?id=<id>` - Удаление данных

### Синхронизация
- `POST /api/v1/sync` - Синхронизация данных с сервером

## Конфигурация

### Переменные окружения

#### Сервер
- `PORT` - Порт сервера (по умолчанию: 8080)
- `DB_HOST` - Хост базы данных (по умолчанию: localhost)
- `DB_PORT` - Порт базы данных (по умолчанию: 5432)
- `DB_USER` - Пользователь базы данных (по умолчанию: gophkeeper)
- `DB_PASSWORD` - Пароль базы данных (по умолчанию: password)
- `DB_NAME` - Имя базы данных (по умолчанию: gophkeeper)
- `JWT_SECRET` - Секретный ключ JWT
- `ENCRYPTION_KEY` - Ключ шифрования данных

#### Клиент
- `SERVER_URL` - URL сервера (по умолчанию: http://localhost:8080)
- `CLIENT_CONFIG_DIR` - Директория конфигурации клиента (по умолчанию: ~/.gophkeeper)

### Файл .env

```env
# Server
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=gophkeeper
DB_PASSWORD=password
DB_NAME=gophkeeper
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=your-encryption-key

# Client
SERVER_URL=http://localhost:8080
CLIENT_CONFIG_DIR=
```

## Структура проекта

```
Gophkeeper/
├── cmd/
│   ├── client/                    # CLI клиентское приложение
│   └── server/                    # HTTP серверное приложение
├── internal/
│   ├── app/                       # Приложения
│   │   ├── client/                # Клиентское приложение
│   │   └── server/                # Серверное приложение
│   ├── client/                    # Клиентская логика
│   │   ├── tests/                 # Тесты клиента
│   │   │   ├── mocks/             # Моки для тестирования
│   │   │   ├── auth_service_test.go
│   │   │   ├── data_service_test.go
│   │   │   ├── sync_service_test.go
│   │   │   └── integration_test.go
│   │   ├── cli/                   # CLI команды
│   │   │   └── tests/             # Тесты CLI
│   │   ├── migrations/            # Миграции клиента
│   │   ├── auth_service.go        # Сервис аутентификации
│   │   ├── data_service.go        # Сервис данных
│   │   ├── sync_service.go        # Сервис синхронизации
│   │   ├── http_client.go         # HTTP клиент
│   │   ├── token_manager.go       # Менеджер токенов
│   │   ├── interfaces.go          # Интерфейсы
│   │   ├── client.go              # Главный клиент
│   │   └── storage.go             # Локальное хранилище
│   ├── config/                    # Конфигурация
│   ├── crypto/                    # Шифрование и хеширование
│   │   └── tests/                 # Тесты криптографии
│   ├── database/                  # Операции с базой данных
│   │   └── migrations/            # Миграции сервера
│   ├── migrate/                   # Система миграций
│   ├── models/                    # Модели данных
│   │   └── tests/                 # Тесты моделей
│   └── server/                    # Реализация сервера
├── tests/                         # E2E тесты
├── scripts/                       # Скрипты сборки
├── docker-compose.yml             # Docker Compose конфигурация
├── Dockerfile.server              # Dockerfile для сервера
├── Dockerfile.client              # Dockerfile для клиента
├── Makefile                       # Автоматизация сборки
└── README.md                      # Этот файл
```

## Разработка

### Качество кода

```bash
# Форматирование кода
make fmt

# Запуск линтера
make lint

# Проверка зависимостей
make deps
```

### Миграции

```bash
# Сервер (PostgreSQL)
make migrate-srv

# Клиент (SQLite)
make migrate-cli
```

## Лицензия

MIT License
