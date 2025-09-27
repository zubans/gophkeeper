# Реорганизация тестов GophKeeper

## Обзор изменений

Все тестовые файлы были перемещены в отдельные папки `tests/` внутри соответствующих пакетов для лучшей организации кода.

## Структура до изменений

```
internal/
├── crypto/
│   ├── encryption.go
│   ├── encryption_test.go
│   ├── hash.go
│   ├── hash_test.go
│   ├── jwt.go
│   └── jwt_test.go
└── models/
    ├── data.go
    ├── response.go
    └── response_test.go
```

## Структура после изменений

```
internal/
├── crypto/
│   ├── encryption.go
│   ├── hash.go
│   ├── jwt.go
│   └── tests/
│       ├── encryption_test.go
│       ├── hash_test.go
│       └── jwt_test.go
└── models/
    ├── data.go
    ├── response.go
    └── tests/
        └── response_test.go
```

## Изменения в тестовых файлах

### 1. Обновление package declaration
```go
// Было
package crypto

// Стало
package tests
```

### 2. Добавление импортов
```go
import (
    "testing"
    "gophkeeper/internal/crypto"  // Добавлен импорт основного пакета
)
```

### 3. Обновление вызовов функций
```go
// Было
hash, err := HashPassword(password)

// Стало
hash, err := crypto.HashPassword(password)
```

## Преимущества новой структуры

1. **Четкое разделение**: Тесты отделены от основного кода
2. **Лучшая организация**: Легче найти и управлять тестами
3. **Модульность**: Каждый пакет имеет свою папку с тестами
4. **Совместимость**: Все существующие команды `go test` продолжают работать

## Команды для запуска тестов

### Все тесты
```bash
go test -v ./...
make test
```

### Тесты конкретного пакета
```bash
go test -v ./internal/crypto/tests
go test -v ./internal/models/tests
```

### Тесты с покрытием
```bash
make test-coverage
```

## Проверка работоспособности

- ✅ Все тесты проходят успешно
- ✅ Приложение компилируется без ошибок
- ✅ Makefile команды работают корректно
- ✅ Структура проекта стала более организованной

## Файлы, которые были перемещены

1. `internal/crypto/encryption_test.go` → `internal/crypto/tests/encryption_test.go`
2. `internal/crypto/hash_test.go` → `internal/crypto/tests/hash_test.go`
3. `internal/crypto/jwt_test.go` → `internal/crypto/tests/jwt_test.go`
4. `internal/models/response_test.go` → `internal/models/tests/response_test.go`

## Удаленные файлы

- `simple_test.go` (временный файл)
- `debug_test.go` (временный файл)

Все изменения обратно совместимы и не влияют на функциональность приложения.
