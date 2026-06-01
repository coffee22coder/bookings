# Bookings

Pet-проект на Go + PostgreSQL: API для работы с бронированиями.

## Requirements

- Go 1.22+
- Docker (PostgreSQL)
- GNU Make (опционально; на Windows — [GnuWin32 Make](https://gnuwin32.sourceforge.net/packages/make.htm))

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
