# Bookings

Pet-проект на Go + PostgreSQL: API для работы с бронированиями.

## Requirements

- Go 1.22+
- Docker (PostgreSQL)
- GNU Make (опционально; на Windows — [GnuWin32 Make](https://gnuwin32.sourceforge.net/packages/make.htm))

## Run

Точка входа приложения — `cmd/api`:

```bash
go run ./cmd/api
```

Сборка и запуск бинарника:

```bash
make build
./bin/api          # macOS / Linux
.\bin\api.exe      # Windows (если Makefile собрал с .exe)
```

Проверка, что весь проект компилируется:

```bash
go build ./...
```

## Project layout

Hexagonal (ports & adapters) layout:

```
bookings/
├── cmd/
│   └── api/
│       └── main.go          # точка входа HTTP API
├── internal/
│   ├── config/              # загрузка конфигурации из env
│   ├── domain/              # сущности и бизнес-правила
│   ├── port/                # интерфейсы (контракты репозиториев, сервисов)
│   ├── service/             # use cases (бизнес-логика)
│   ├── adapter/
│   │   ├── http/            # HTTP handlers, middleware, DTO
│   │   └── postgres/        # реализация репозиториев через pgx
│   └── testutil/            # общие helpers для тестов
├── docs/
│   ├── PROJECT.md
│   └── TASKS.md
├── docker-compose.yml
├── Makefile
└── .env.example
```

### Зачем `internal/`

Пакеты в `internal/` доступны только внутри модуля `github.com/coffee22coder/bookings`. Внешний код не может их импортировать — это защита от случайного использования внутренних деталей как публичной библиотеки.

Правило зависимостей: `adapter` → `service` → `domain`. Домен не знает про HTTP и PostgreSQL.

| Куда класть код | Пример |
|---|---|
| HTTP handler | `internal/adapter/http/` |
| SQL-запросы | `internal/adapter/postgres/` |
| Бизнес-логика | `internal/service/` |
| Интерфейс репозитория | `internal/port/` |
| Структура `Booking` | `internal/domain/` |

## Testing

Unit-тесты (быстрые, без внешних зависимостей):

```bash
make test-unit
```

Integration-тесты (с build tag `integration`, требуют Docker/PostgreSQL):

```bash
make test-integration
```

Все тесты:

```bash
make test
```

Подробный вывод:

```bash
go test -v ./internal/...
```

На Windows, если `make` не в `PATH`:

```powershell
& "C:\Program Files (x86)\GnuWin32\bin\make.exe" test-unit
```

## Database

PostgreSQL поднимается через Docker:

```bash
docker compose up -d
```

Переменные окружения — см. `.env.example`.
