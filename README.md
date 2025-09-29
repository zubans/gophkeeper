# GophKeeper

GophKeeper - это безопасная клиент-серверная система для хранения приватных данных, построенная на Go. Позволяет пользователям надёжно хранить и синхронизировать пароли, текстовые данные, бинарные файлы и информацию о банковских картах между несколькими устройствами.

## Возможности

### Сервер
- Регистрация, аутентификация и авторизация пользователей
- Безопасное хранение данных с шифрованием
- Синхронизация данных между несколькими клиентами
- RESTful API для взаимодействия с клиентами
- JWT-аутентификация
- Поддержка PostgreSQL
- Контроль версий (последние 10 версий)
- Разрешение конфликтов по времени

### Клиент
- Кроссплатформенное CLI-приложение (Windows, Linux, macOS)
- Безопасная аутентификация с удалённым сервером
- Поддержка различных типов данных:
  - Пары логин/пароль
  - Произвольные текстовые данные
  - Бинарные данные
  - Информация о банковских картах
- Локальное кэширование и синхронизация данных
- Отображение информации о версии и дате сборки
- Локальное хранение в SQLite

## Типы данных

Система поддерживает следующие типы данных:

1. **Login/Password** - Хранение учётных данных веб-сайтов с опциональными метаданными
2. **Text** - Хранение произвольных текстовых данных
3. **Binary** - Хранение бинарных файлов и данных
4. **Bank Card** - Безопасное хранение информации о банковских картах

Все типы данных поддерживают пользовательские метаданные для дополнительного контекста.

## Архитектура

### Сервер (PostgreSQL)
- **База данных**: PostgreSQL с миграциями
- **Таблицы**:
  - `users` - пользователи
  - `stored_data` - основные данные с полями `is_deleted`, `version`
  - `data_history` - история версий (последние 10)
  - `schema_migrations` - управление миграциями

### Клиент (SQLite)
- **База данных**: SQLite с миграциями
- **Таблицы**: аналогичные серверным, но с SQLite синтаксисом
- **Локальное хранение**: все данные сохраняются локально для офлайн работы

## Механизм синхронизации

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

## Контроль версий

### История изменений
- Каждое изменение сохраняется в `data_history`
- Хранится последние 10 версий для каждого элемента данных
- Поддерживается soft delete (`is_deleted`)

### Версионирование
- Автоматическое увеличение `version` при изменениях
- Уникальные ID для истории: `{data_id}_v{version}`
- Очистка старых версий (оставляем только 10 последних)

## Установка

### Предварительные требования

- Go 1.23.6 или новее
- PostgreSQL 12 или новее (для сервера)
- Git

### Сборка из исходного кода

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd Gophkeeper
```

2. Загрузите зависимости:
```bash
make deps
```

3. Соберите приложение:
```bash
make build
```

4. Соберите для всех платформ:
```bash
make build-all
```

### Настройка базы данных

1. Установите и запустите PostgreSQL
2. Создайте базу данных:
```bash
createdb gophkeeper
```

3. Сервер автоматически создаст необходимые таблицы при первом запуске.

## Использование

### Запуск сервера

```bash
# Используя make
make run-server

# Или напрямую
go run ./cmd/server

# С пользовательскими настройками базы данных
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

## Безопасность

- Все данные шифруются с использованием AES-256-GCM перед сохранением
- Пароли хешируются с использованием SHA-256 с солью
- JWT токены используются для аутентификации
- HTTPS рекомендуется для продакшн развертываний
- Данные клиента шифруются локально перед передачей

## Docker

### Быстрый старт

```bash
# Сборка и запуск всех сервисов
make docker-up

# Или напрямую
docker-compose up -d
```

### Сервисы

- **PostgreSQL Database** - Порт 5432
- **GophKeeper Server** - Порт 8080
- **GophKeeper Client** - Интерактивный режим
- **pgAdmin** - Порт 5050 (admin@gophkeeper.local / admin)

### Команды Docker

```bash
# Просмотр логов
make docker-logs

# Остановка сервисов
make docker-down

# Сборка образов
make docker-build
```

## Разработка

### Запуск тестов

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
make test-coverage

# Запуск интеграционных тестов
go test ./tests -v
```

### Качество кода

```bash
# Форматирование кода
make fmt

# Запуск линтера
make lint
```

## Структура проекта

```
Gophkeeper/
├── cmd/
│   ├── client/          # CLI клиентское приложение
│   └── server/          # HTTP серверное приложение
├── internal/
│   ├── client/          # Реализация клиента
│   ├── config/          # Конфигурация
│   ├── crypto/          # Шифрование и хеширование
│   ├── database/        # Операции с базой данных
│   │   └── migrations/  # Миграции сервера
│   ├── migrate/         # Система миграций
│   ├── models/          # Модели данных
│   └── server/          # Реализация сервера
├── tests/               # Тесты
├── scripts/             # Скрипты сборки
├── docker-compose.yml   # Docker Compose конфигурация
├── Dockerfile.server    # Dockerfile для сервера
├── Dockerfile.client    # Dockerfile для клиента
└── Makefile            # Автоматизация сборки
```

## Модели данных

### User
```go
type User struct {
    ID           string    `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
```

### StoredData
```go
type StoredData struct {
    ID         string    `json:"id" db:"id"`
    UserID     string    `json:"user_id" db:"user_id"`
    Type       DataType  `json:"type" db:"type"`
    Title      string    `json:"title" db:"title"`
    Data       []byte    `json:"data" db:"data"`
    Metadata   string    `json:"metadata" db:"metadata"`
    Version    int       `json:"version" db:"version"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
    LastSyncAt time.Time `json:"last_sync_at" db:"last_sync_at"`
    IsDeleted  bool      `json:"is_deleted" db:"is_deleted"`
}
```

### DataHistory
```go
type DataHistory struct {
    ID        string    `json:"id" db:"id"`
    DataID    string    `json:"data_id" db:"data_id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Type      DataType  `json:"type" db:"type"`
    Title     string    `json:"title" db:"title"`
    Data      []byte    `json:"data" db:"data"`
    Metadata  string    `json:"metadata" db:"metadata"`
    Version   int       `json:"version" db:"version"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}
```

## Способы шифрования

### Шифрование данных
- **Алгоритм**: AES-256-GCM
- **Ключ**: 32-байтовый ключ шифрования
- **Режим**: GCM для аутентификации и шифрования

### Хеширование паролей
- **Алгоритм**: SHA-256
- **Соль**: 32-байтовая случайная соль
- **Итерации**: Одна итерация с солью

### JWT токены
- **Алгоритм**: HMAC-SHA256
- **Срок действия**: 24 часа
- **Поля**: user_id, username, exp, iat

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

Создайте файл `.env` в корне проекта:

```env
# Server defaults
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=gophkeeper
DB_PASSWORD=password
DB_NAME=gophkeeper
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=your-encryption-key

# Client defaults
SERVER_URL=http://localhost:8080
CLIENT_CONFIG_DIR=
```

## Тестирование

### Типы тестов

1. **Unit тесты** - Тестирование отдельных компонентов
2. **Integration тесты** - Тестирование взаимодействия компонентов
3. **End-to-end тесты** - Полное тестирование пользовательских сценариев

### Запуск тестов

```bash
# Все тесты
go test ./tests -v

# Только unit тесты
go test ./tests -v -run TestCrypto

# Только integration тесты
go test ./tests -v -run TestIntegration
```

## Лицензия

MIT License

## Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## Поддержка

Для получения поддержки создайте issue в репозитории проекта.