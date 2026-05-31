# Микрозадачи — Avia Bookings API

> **Версия:** 3.1  
> **Как читать:** задачи **строго по ID**. Каждая — один каркас ниже.  
> **ТЗ:** [PROJECT.md](./PROJECT.md) · **Ревью:** `ревью S1-07`

---

## Каркас одной задачи

| Поле             | Что писать                                         |
| ---------------- | -------------------------------------------------- |
| **Задача**       | ID + короткое имя                                  |
| **Цель**         | Зачем этот шаг (1–2 предложения)                   |
| **Описание**     | Зависимости + 2–3 шага + «зачем» при необходимости |
| **Результат**    | Команда и ожидаемый выход (проверка)               |
| **Документация** | Ссылки; если не нужна — «—»                        |

**Правило:** `Результат` выполнен → следующая задача → ревью по желанию.

---

# Sprint 0 — Git и репозиторий (сделай до или сразу после S1-01)

**Итог:** локальный git, `.gitignore`, первый commit, пуш на GitHub/GitLab.

| Блок           | ID            | Задач |
| -------------- | ------------- | ----- |
| R. Репозиторий | S0-01 … S0-04 | 4     |

**Зачем отдельный блок:** не закоммитить `.env`, дамп целиком и бинарники; портфолио = публичный URL репо.

---

## R. Git и remote

### S0-01

**Задача:** S0-01 — Инициализировать git

**Цель:** Версионировать код локально и готовить пуш на GitHub.

**Описание:**

Зависимость: папка проекта `bookings/` с `go.mod`.

**Шаги:**

1. В корне проекта: `git init` (если ещё нет `.git` — проверь `ls -la .git`)
2. Убедись, что ты **не** в домашней директории — только каталог с `go.mod`

**Результат:** `git status` → «On branch main» (или master), список untracked файлов

**Документация:** [https://git-scm.com/docs/git-init](https://git-scm.com/docs/git-init)

---

### S0-02

**Задача:** S0-02 — Создать `.gitignore`

**Цель:** Исключить секреты, артефакты и тяжёлые файлы до первого commit.

**Описание:**

Зависимость: **S0-01**.

**Шаги:**

1. Создай `.gitignore` в корне (шаблон Go + проект):

- `.env` (пароли БД)
- `coverage.out`, `*.test`
- `bin/`, `api`, `dist/`
- `.idea/`, `.vscode/` (опционально)
- **не** коммить `demo-20250901-3m.sql` если файл >100MB — см. шаг 2

2. Реши про дамп: **вариант A** — добавить `*.sql` в gitignore (дамп только локально); **вариант B** — Git LFS (позже). Для портфолио обычно **A** + в README «импорт дампа по ссылке»

**Результат:** `git status` не показывает `.env` после `touch .env`

**Документация:** [https://github.com/github/gitignore/blob/main/Go.gitignore](https://github.com/github/gitignore/blob/main/Go.gitignore)

---

### S0-03

**Задача:** S0-03 — Первый commit

**Цель:** Зафиксировать стартовое состояние (compose, go.mod, docs, dbcheck).

**Описание:**

Зависимость: **S0-02**. **До commit** — никаких паролей в отслеживаемых файлах.

**Шаги:**

1. `git add .` → проверь `git status` (нет `.env`, нет гигантского sql если в ignore)
2. `git commit -m "chore: initial project setup (docker, go.mod, docs)"`

**Результат:** `git log -1` показывает commit; `git status` clean

**Документация:** [https://git-scm.com/docs/git-commit](https://git-scm.com/docs/git-commit)

---

### S0-04

**Задача:** S0-04 — Подключить удалённый репозиторий

**Цель:** Публичный URL для портфолио и ревью.

**Описание:**

Зависимость: **S0-03**. Аккаунт GitHub (или GitLab).

**Шаги:**

1. На GitHub: **New repository** → имя `bookings` (или как module path) → **без** README/license (уже есть локально)
2. Локально: `git remote add origin git@github.com:<username>/bookings.git` (или HTTPS URL)
3. `git push -u origin main` (ветка может называться `master` — приведи к одному имени)

**Результат:** в браузере видны файлы; `git remote -v` показывает origin

**Документация:** [https://docs.github.com/en/repositories/creating-and-managing-repositories](https://docs.github.com/en/repositories/creating-and-managing-repositories)

**Ревью:** `ревью S0-CP` — скинь URL репо (без секретов в истории).

---

## S0-CP

**Задача:** S0-CP — Git готов

**Цель:** Репозиторий готов к работе по TASKS.

**Описание:**

**Шаги:**

1. Проверь на GitHub: есть `docs/`, `go.mod`, `docker-compose.yml`, нет `.env`
2. В README одна строка: «clone → docker compose → …»

**Результат (gate):** публичный/приватный URL + первый commit + `.gitignore` работает

**Документация:** —

---

# Sprint 1 — Foundation

**Итог спринта:** `GET /api/v1/airports`, health/ready, `make test-unit` green.

| Блок                     | ID            | Задач |
| ------------------------ | ------------- | ----- |
| A. Тест-инфра            | S1-01 … S1-06 | 6     |
| B. Структура             | S1-07 … S1-11 | 5     |
| C. Config                | S1-12 … S1-19 | 8     |
| D. Postgres pool         | S1-20 … S1-26 | 7     |
| E. HTTP health           | S1-27 … S1-38 | 12    |
| F. Airports slice        | S1-39 … S1-52 | 14    |
| G. Testcontainers (опц.) | S1-53 … S1-55 | 3     |

**Checkpoint:** [S1-CP](#s1-cp)

---

## A. Тест-инфраструктура

### S1-01

**Задача:** S1-01 — Подключить testify

**Цель:** Подключить библиотеку для удобных assert/require в unit-тестах.

**Описание:**

Зависимость: инициализированный Go-модуль в корне репо.

**Шаги:**

1. В терминале: `go get github.com/stretchr/testify/require`
2. Убедись, что в `go.mod` появилась зависимость testify

**Результат:** `go mod verify` → exit 0

**Документация:** [https://github.com/stretchr/testify#assert-and-require](https://github.com/stretchr/testify#assert-and-require)

---

### S1-02

**Задача:** S1-02 — Создать Makefile (unit)

**Цель:** Единые команды `make test-unit` / `test-integration` для всего репо.

**Описание:**

Зависимость: сначала закрой задачу **S1-01**.

**Шаги:**

1. Создай файл `Makefile` в корне репо
2. Добавь target `test-unit` из PROJECT.md §9.3 (только `go test ./internal/...`)

**Результат:** `make test-unit` → может быть «no test files» — это ok

**Документация:** [https://pkg.go.dev/cmd/go#hdr-Test_packages](https://pkg.go.dev/cmd/go#hdr-Test_packages)

---

### S1-03

**Задача:** S1-03 — Target test-integration

**Цель:** Отделить быстрые unit-тесты от медленных integration (`-tags=integration`).

**Описание:**

Зависимость: сначала закрой задачу **S1-02**.

**Шаги:**

1. В `Makefile` добавь `test-integration` с `-tags=integration`
2. Добавь `test: test-unit test-integration`

**Результат:** `make test-integration` → ok или «no packages» (пока нет integration)

**Документация:** —

---

### S1-04

**Задача:** S1-04 — Дополнить `.gitignore`

**Цель:** После новых артефактов сборки ничего лишнего не попадает в commit.

**Описание:**

Зависимость: **S0-02** (базовый gitignore) + **S1-02**.

**Шаги:**

1. Добавь в `.gitignore` то, что появилось позже: `bin/`, `coverage.out` (если ещё нет)
2. `git status` после `make test-unit` / `go build` — нет мусора в untracked

**Результат:** бинарник `bin/api` не в `git status` как файл для коммита

**Документация:** —

---

### S1-05

**Задача:** S1-05 — Пакет testutil

**Цель:** Завести пакет для общих test helpers.

**Описание:**

Зависимость: сначала закрой задачу **S1-02**.

**Шаги:**

1. Создай `internal/testutil/doc.go` с `package testutil` и комментарием пакета
2. Создай пустой `internal/config/config_test.go` с `package config_test` или `package config`

**Результат:** `go build ./internal/...` → exit 0

**Документация:** —

---

### S1-06

**Задача:** S1-06 — Первый тест + README Testing

**Цель:** Убедиться, что инфраструктура тестов реально работает.

**Описание:**

Зависимость: сначала закрой задачу **S1-05**.

**Шаги:**

1. В `config_test.go` напиши `TestPlaceholder` с `require.True(t, true)`
2. В `README.md` добавь секцию **Testing** с командой `make test-unit`

**Результат:** `make test-unit` → **PASS**, 1 test

**Документация:** [https://go.dev/doc/tutorial/add-a-test](https://go.dev/doc/tutorial/add-a-test)

---

## B. Структура проекта

### S1-07

**Задача:** S1-07 — Папки internal

**Цель:** Разложить код по hexagonal layout до написания фич.

**Описание:**

Зависимость: сначала закрой задачу **S1-06**.

**Шаги:**

1. Создай каталоги: `internal/config`, `internal/domain`, `internal/port`, `internal/service`, `internal/adapter/http`, `internal/adapter/postgres`
2. Создай `cmd/api/`

**Результат:** `ls internal` показывает все папки

**Документация:** PROJECT.md §4

---

### S1-08

**Задача:** S1-08 — cmd/api/main.go (заглушка)

**Цель:** Отделить точку входа от бизнес-кода (`cmd/api`).

**Описание:**

Зависимость: сначала закрой задачу **S1-07**.

**Шаги:**

1. Создай `cmd/api/main.go`, `package main`, `slog.Info("api starting")`
2. Не подключай пока БД

**Результат:** `go run ./cmd/api` → в консоли лог «starting», exit 0

**Документация:** —

---

### S1-09

**Задача:** S1-09 — Сборка бинарника

**Цель:** Отделить точку входа от бизнес-кода (`cmd/api`).

**Описание:**

Зависимость: сначала закрой задачу **S1-08**.

**Шаги:**

1. Выполни `go build -o bin/api ./cmd/api`
2. Добавь `bin/` в `.gitignore`

**Результат:** `./bin/api` запускается

**Документация:** —

---

### S1-10

**Задача:** S1-10 — Удалить корневой main.go

**Цель:** Отделить точку входа от бизнес-кода (`cmd/api`).

**Описание:**

Зависимость: сначала закрой задачу **S1-09**.

**Шаги:**

1. Удали `main.go` из корня (старый dbcheck)
2. Обнови README: точка входа `cmd/api`

**Результат:** `go build ./...` — нет ошибок; в корне нет `main.go`

**Документация:** —

---

### S1-11

**Задача:** S1-11 — README: дерево проекта

**Цель:** Онбординг: другой разработчик поднимает проект за 10 минут.

**Описание:**

Зависимость: сначала закрой задачу **S1-10**.

**Шаги:**

1. В README вставь дерево каталогов (скопируй из PROJECT.md §4, подправь под факт)
2. Одним абзацем: зачем `internal/`

**Результат:** ревью-чек: новичок понимает, куда класть handler

**Документация:** —

---

## C. Config

### S1-12

**Задача:** S1-12 — Зависимость envconfig

**Цель:** Конфигурация только из env — готово к Docker и 12-factor.

**Описание:**

Зависимость: сначала закрой задачу **S1-11**.

**Шаги:**

1. `go get github.com/kelseyhightower/envconfig`
2. Создай `internal/config/config.go` — пустой `package config`

**Результат:** `go mod tidy` → ok

**Документация:** —

---

### S1-13

**Задача:** S1-13 — Struct Config

**Цель:** Конфигурация только из env — готово к Docker и 12-factor.

**Описание:**

Зависимость: сначала закрой задачу **S1-12**.

**Шаги:**

1. В `config.go` объяви `type Config struct` с полями: `HTTPPort`, `DBHost`, `DBPort`, `DBUser`, `DBPassword`, `DBName`, `DBSSLMode`
2. Добавь struct tags `envconfig:"..."`

**Результат:** `go build ./internal/config` → ok

**Документация:** [https://github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)

---

### S1-14

**Задача:** S1-14 — Load() без валидации

**Цель:** Шаг к MVP: Load() без валидации.

**Описание:**

Зависимость: сначала закрой задачу **S1-13**.

**Шаги:**

1. Создай `load.go` с `func Load() (Config, error)` — вызов `envconfig.Process("", &cfg)`
2. Пока без проверок портов

**Результат:** временный `TestLoad_Smoke` в `config_test.go` с `t.Setenv` для DB_HOST → Load не паникует

**Документация:** —

---

### S1-15

**Задача:** S1-15 — .env.example

**Цель:** Конфигурация только из env — готово к Docker и 12-factor.

**Описание:**

Зависимость: сначала закрой задачу **S1-14**.

**Шаги:**

1. Создай `.env.example` со всеми ключами, `DB_PORT=5433`, `DB_NAME=demo`
2. Скопируй у себя в `.env` (локально, не в git)

**Результат:** файл есть; `.env` в gitignore

**Документация:** —

---

### S1-16

**Задача:** S1-16 — Метод DSN()

**Цель:** Конфигурация только из env — готово к Docker и 12-factor.

**Описание:**

Зависимость: сначала закрой задачу **S1-15**.

**Шаги:**

1. Добавь `(c Config) DSN() string` — URI `postgres://...`
2. Не логируй DSN в main

**Результат:** unit-тест: при заданных env DSN содержит `demo` и `5433`

**Документация:** —

---

### S1-17

**Задача:** S1-17 — Валидация Load

**Цель:** Шаг к MVP: Валидация Load.

**Описание:**

Зависимость: сначала закрой задачу **S1-16**.

**Шаги:**

1. В `Load()` добавь проверки: пустой `DBHost` → error; порт HTTP 1–65535
2. Ошибки оберни: `fmt.Errorf("load config: %w", err)`

**Результат:** `TestLoad_MissingDBHost` → error; `TestLoad_InvalidHTTPPort` → error

**Документация:** —

---

### S1-18

**Задача:** S1-18 — Подключить config в main

**Цель:** Отделить точку входа от бизнес-кода (`cmd/api`).

**Описание:**

Зависимость: сначала закрой задачу **S1-17**.

**Шаги:**

1. В `cmd/api/main.go`: `cfg, err := config.Load()`; при err — `slog.Error` + `os.Exit(1)`
2. Залогируй только `HTTPPort`, не пароль

**Результат:** запуск без env → exit 1; с env из `.env` → старт ok

**Документация:** —

---

### S1-19

**Задача:** S1-19 — make run

**Цель:** Отделить точку входа от бизнес-кода (`cmd/api`).

**Описание:**

Зависимость: сначала закрой задачу **S1-18**.

**Шаги:**

1. В `Makefile` добавь `run: go run ./cmd/api`
2. Документируй в README: `export $(cat .env | xargs)` или direnv

**Результат:** `make run` → api starting

**Документация:** —

---

## D. Postgres pool

### S1-20

**Задача:** S1-20 — Файл pool.go

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-19, Docker PG up**.

**Шаги:**

1. Создай `internal/adapter/postgres/pool.go`, `package postgres`
2. Объяви сигнатуру `func NewPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error)` — пока заглушка `return nil, nil` не надо — сразу реализуй New

_Зачем:_ одно соединение на все горутины — антипаттерн; нужен `pgxpool`.

**Результат:** `go build ./internal/adapter/postgres` — компиляция (если реализовал New с pgxpool)

**Документация:** [https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)

---

### S1-21

**Задача:** S1-21 — NewPool + Ping

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-20**.

**Шаги:**

1. В `NewPool`: `pgxpool.New(ctx, cfg.DSN())`, затем `Ping` с `context.WithTimeout` 5s
2. При ошибке — `fmt.Errorf("ping db: %w", err)`

_Зачем:_ одно соединение на все горутины — антипаттерн; нужен `pgxpool`.

**Результат:** `make run` с валидным `.env` → старт без ошибки БД

**Документация:** —

---

### S1-22

**Задача:** S1-22 — main: создать и закрыть pool

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-21**.

**Шаги:**

1. В `main`: `pool, err := postgres.NewPool(ctx, cfg)`; `defer pool.Close()`
2. При ошибке pool — exit 1

_Зачем:_ одно соединение на все горутины — антипаттерн; нужен `pgxpool`.

**Результат:** останови Docker → `make run` → fail fast с ошибкой ping

**Документация:** —

---

### S1-23

**Задача:** S1-23 — Тест: невалидный DSN

**Цель:** Конфигурация только из env — готово к Docker и 12-factor.

**Описание:**

Зависимость: сначала закрой задачу **S1-21**.

**Шаги:**

1. Создай `pool_test.go`, `TestNewPool_InvalidHost` — cfg с host `invalid-host-xyz`
2. Ожидай error, не panic

**Результат:** `go test ./internal/adapter/postgres -run TestNewPool_Invalid -v` → PASS

**Документация:** —

---

### S1-24

**Задача:** S1-24 — integration: ping реальной БД

**Цель:** Шаг к MVP: integration: ping реальной БД.

**Описание:**

Зависимость: сначала закрой задачу **S1-23**.

**Шаги:**

1. Создай `test/integration/pool_test.go` с `//go:build integration`
2. Тест `TestPool_Ping`: Load config из env, NewPool, Ping; `t.Skip` если `testing.Short()`

**Результат:** `make test-integration` → PASS (нужен postgres на 5433)

**Документация:** —

---

### S1-25

**Задача:** S1-25 — README TEST_DATABASE_URL

**Цель:** Онбординг: другой разработчик поднимает проект за 10 минут.

**Описание:**

Зависимость: сначала закрой задачу **S1-24**.

**Шаги:**

1. В README добавь: для integration нужен Docker и `DB_`\* env
2. Пример export переменных

**Результат:** — (документация)

**Документация:** —

---

### S1-26

**Задача:** S1-26 — Тест Close pool

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-24**.

**Шаги:**

1. В `pool_test.go` (unit или integration): создай pool, `Close()`, второй `Close()` не паникует
2. Используй `t.Cleanup` для Close

_Зачем:_ одно соединение на все горутины — антипаттерн; нужен `pgxpool`.

**Результат:** `go test ./internal/adapter/postgres -v` → PASS

**Документация:** —

---

## E. HTTP + health/ready

### S1-27

**Задача:** S1-27 — Зависимость chi

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-22**.

**Шаги:**

1. `go get github.com/go-chi/chi/v5`
2. `go get github.com/go-chi/chi/v5/middleware`

**Результат:** `go mod tidy` → ok

**Документация:** —

---

### S1-28

**Задача:** S1-28 — server.go каркас

**Цель:** Шаг к MVP: server.go каркас.

**Описание:**

Зависимость: сначала закрой задачу **S1-27**.

**Шаги:**

1. Создай `internal/adapter/http/server.go`
2. Struct `Server` с полями: `cfg`, `logger`, `router chi.Router`, зависимость для ping (интерфейс позже)

**Результат:** `go build ./internal/adapter/http` → ok

**Документация:** —

---

### S1-29

**Задача:** S1-29 — Router + GET /health

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-28**.

**Шаги:**

1. Функция `NewServer(...)` строит `chi.NewRouter()`
2. Handler `GET /health` → JSON `{"status":"ok"}`

**Результат:** пока без main — напиши временный тест в следующей задаче

**Документация:** —

---

### S1-30

**Задача:** S1-30 — httptest /health

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-29**.

**Шаги:**

1. Создай `server_test.go`, `TestHealth_Returns200`
2. `httptest.NewRecorder` + `NewRequest("GET","/health",nil)`

**Результат:** `go test ./internal/adapter/http -run TestHealth -v` → PASS

**Документация:** [https://pkg.go.dev/net/http/httptest](https://pkg.go.dev/net/http/httptest)

---

### S1-31

**Задача:** S1-31 — Middleware RequestID

**Цель:** Шаг к MVP: Middleware RequestID.

**Описание:**

Зависимость: сначала закрой задачу **S1-30**.

**Шаги:**

1. Подключи `middleware.RequestID` в router
2. В access log (или handler) логируй request id из context

**Результат:** расширь тест: в ответе есть заголовок `X-Request-Id` (если chi отдаёт)

**Документация:** —

---

### S1-32

**Задача:** S1-32 — Интерфейс DBPing

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-29**.

**Шаги:**

1. В `internal/port` создай `health.go`: `type DBPinger interface { Ping(ctx context.Context) error }`
2. Адаптер-обёртка над `*pgxpool.Pool` в postgres-пакете (метод Ping)

**Результат:** `go build ./...` — pool удовлетворяет интерфейсу

**Документация:** —

---

### S1-33

**Задача:** S1-33 — GET /ready

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-32**.

**Шаги:**

1. Handler `GET /ready`: вызывает `Ping`; ok → 200, err → 503 JSON
2. Подключи route в router

_Зачем:_ readiness отличается от liveness — проверяет БД, не только процесс.

**Результат:** допиши httptest с fake pinger (см. S1-34)

**Документация:** —

---

### S1-34

**Задача:** S1-34 — httptest /ready 503 (fake)

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-33**.

**Шаги:**

1. В тесте передай stub `DBPinger` с `return errors.New("down")`
2. `TestReady_DBDown_Returns503`

_Зачем:_ readiness отличается от liveness — проверяет БД, не только процесс.

**Результат:** `go test ./internal/adapter/http -run TestReady_DBDown -v` → PASS

**Документация:** —

---

### S1-35

**Задача:** S1-35 — httptest /ready 200 (real pool)

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-34, S1-24**.

**Шаги:**

1. Integration-тест: Server с реальным pool → `/ready` 200
2. Файл в `test/integration/health_test.go` + build tag

_Зачем:_ одно соединение на все горутины — антипаттерн; нужен `pgxpool`.

**Результат:** `make test-integration` → PASS

**Документация:** —

---

### S1-36

**Задача:** S1-36 — Recover middleware

**Цель:** Шаг к MVP: Recover middleware.

**Описание:**

Зависимость: сначала закрой задачу **S1-31**.

**Шаги:**

1. Добавь `middleware.Recoverer`
2. Тестовый route `/panic` (только в тестовом router) → 500

**Результат:** `TestRecover_PanicReturns500` → PASS

**Документация:** —

---

### S1-37

**Задача:** S1-37 — main: ListenAndServe

**Цель:** Подключить модуль к `cmd/api/main.go`.

**Описание:**

Зависимость: сначала закрой задачу **S1-35**.

**Шаги:**

1. В `main`: `http.Server{ Addr: ":"+cfg.HTTPPort, Handler: srv.Router() }`
2. Запуск в goroutine; `ReadHeaderTimeout` задай

**Результат:** `make run` + `curl localhost:8080/health` → 200

**Документация:** —

---

### S1-38

**Задача:** S1-38 — curl ready

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-37**.

**Шаги:**

1. `curl -i localhost:8080/ready` — 200 при живой БД
2. `docker compose stop` postgres → ready 503 → `docker compose start`

_Зачем:_ readiness отличается от liveness — проверяет БД, не только процесс.

**Результат:** статусы соответствуют

**Документация:** —

---

## F. Airports (вертикальный slice)

### S1-39

**Задача:** S1-39 — domain.Airport

**Цель:** Первый вертикальный slice: DB → service → HTTP.

**Описание:**

Зависимость: сначала закрой задачу **S1-38**.

**Шаги:**

1. Создай `internal/domain/airport.go` — поля по колонкам `airports_data` (изучи в DBeaver)
2. Без json/db тегов в domain

_Зачем:_ domain не знает про HTTP/SQL — только бизнес-данные.

**Результат:** `go build ./internal/domain` → ok

**Документация:** —

---

### S1-40

**Задача:** S1-40 — port.AirportRepository

**Цель:** Первый вертикальный slice: DB → service → HTTP.

**Описание:**

Зависимость: сначала закрой задачу **S1-39**.

**Шаги:**

1. Создай `internal/port/airport_repository.go` с `List(ctx, limit, offset)` и `Count(ctx)`
2. Только интерфейс

**Результат:** `go build ./internal/port` → ok

**Документация:** —

---

### S1-41

**Задача:** S1-41 — SQL в DBeaver

**Цель:** Проверить SQL на дампе до Go-кода.

**Описание:**

Зависимость: сначала закрой задачу **S1-40**.

**Шаги:**

1. Напиши `SELECT ... FROM bookings.airports_data LIMIT 5 OFFSET 0`
2. Напиши `SELECT COUNT(*) FROM bookings.airports_data`
3. Сохрани SQL в комментарии в `airport_repo.go` (не выполняй из кода строкой из файла без $1)

**Результат:** оба запроса возвращают данные в DBeaver

**Документация:** —

---

### S1-42

**Задача:** S1-42 — postgres List

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S1-41**.

**Шаги:**

1. Создай `internal/adapter/postgres/airport_repo.go`, struct repo с pool
2. Реализуй `List` — Scan в `[]domain.Airport`

**Результат:** временный main или integration в S1-44

**Документация:** —

---

### S1-43

**Задача:** S1-43 — postgres Count

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S1-42**.

**Шаги:**

1. Реализуй `Count` на том же repo
2. Проверь `rows.Err()` в List после цикла

**Результат:** `go build ./internal/adapter/postgres` → ok

**Документация:** —

---

### S1-44

**Задача:** S1-44 — integration airport repo

**Цель:** Первый вертикальный slice: DB → service → HTTP.

**Описание:**

Зависимость: сначала закрой задачу **S1-43**.

**Шаги:**

1. `test/integration/airport_repo_test.go`: List limit=3 → len 3; Count > 0
2. build tag integration

**Результат:** `make test-integration -run Airport` → PASS

**Документация:** —

---

### S1-45

**Задача:** S1-45 — service AirportService

**Цель:** Первый вертикальный slice: DB → service → HTTP.

**Описание:**

Зависимость: сначала закрой задачу **S1-44**.

**Шаги:**

1. `internal/service/airport_service.go` — принимает `port.AirportRepository`
2. Метод `List(ctx, limit, offset)` — дефолт limit=20 если 0

**Результат:** `go build ./internal/service` → ok

**Документация:** —

---

### S1-46

**Задача:** S1-46 — unit test service (ручной fake)

**Цель:** Автотест: unit test service (ручной fake).

**Описание:**

Зависимость: сначала закрой задачу **S1-45**.

**Шаги:**

1. В `airport_service_test.go` создай stub repo (struct с func полями), без mockgen пока
2. `TestAirportService_List_DefaultLimit` — limit 0 → repo получает 20

**Результат:** `go test ./internal/service -run Airport -v` → PASS

**Документация:** —

---

### S1-47

**Задача:** S1-47 — dto + handler list

**Цель:** Шаг к MVP: dto + handler list.

**Описание:**

Зависимость: сначала закрой задачу **S1-46**.

**Шаги:**

1. `internal/adapter/http/dto/airport.go` — response JSON struct
2. `handler/airport_handler.go` — `List(w,r)` парсит query limit/offset

**Результат:** `go build ./internal/adapter/http` → ok

**Документация:** —

---

### S1-48

**Задача:** S1-48 — Зарегистрировать route

**Цель:** Шаг к MVP: Зарегистрировать route.

**Описание:**

Зависимость: сначала закрой задачу **S1-47**.

**Шаги:**

1. В `NewServer` добавь `r.Get("/api/v1/airports", handler.List)`
2. Передай service/handler через конструктор Server

**Результат:** `make run` + curl airports → JSON с items

**Документация:** —

---

### S1-49

**Задача:** S1-49 — httptest airports 200

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-48**.

**Шаги:**

1. Подставь mock service или repo в handler для теста
2. `TestListAirports_200` — status 200, body содержит `"items"`

**Результат:** `go test ./internal/adapter/http -run TestListAirports -v` → PASS

**Документация:** —

---

### S1-50

**Задача:** S1-50 — httptest invalid limit

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S1-49**.

**Шаги:**

1. `TestListAirports_InvalidLimit_400` — query `limit=-1` → 400

**Результат:** PASS

**Документация:** —

---

### S1-51

**Задача:** S1-51 — unit test repo error path

**Цель:** Автотест: unit test repo error path.

**Описание:**

Зависимость: сначала закрой задачу **S1-46**.

**Шаги:**

1. `TestAirportService_List_RepoError` — stub возвращает err → service возвращает err

**Результат:** `make test-unit` → все PASS

**Документация:** —

---

### S1-52

**Задача:** S1-52 — curl + checkpoint doc

**Цель:** Ручная проверка сценария: curl + checkpoint doc.

**Описание:**

Зависимость: сначала закрой задачу **S1-51**.

**Шаги:**

1. `curl "localhost:8080/api/v1/airports?limit=5&offset=0"` — сохрани пример ответа в README
2. Отметь в README выполненные эндпоинты Sprint 1

**Результат:** JSON с 5 аэропортами и total

**Документация:** —

---

## G. Testcontainers (опционально)

### S1-53

**Задача:** S1-53 — go get testcontainers

**Цель:** Изолированная PG для integration без ручного compose.

**Описание:**

Зависимость: сначала закрой задачу **S1-52**.

**Шаги:**

1. `go get github.com/testcontainers/testcontainers-go/modules/postgres`
2. Создай `internal/testutil/postgres.go` — заготовка функции `SetupDB(t *testing.T)`

**Результат:** `go build ./internal/testutil` → ok

**Документация:** —

---

### S1-54

**Задача:** S1-54 — SetupDB реализация

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S1-53**.

**Шаги:**

1. В `SetupDB` подними container, верни `*pgxpool.Pool`, `t.Cleanup` terminate
2. Документируй: нужен Docker Desktop

**Результат:** один тест `TestSetupDB` в testutil или integration PASS

**Документация:** [https://golang.testcontainers.org/modules/postgres/](https://golang.testcontainers.org/modules/postgres/)

---

### S1-55

**Задача:** S1-55 — Перевести integration на harness

**Цель:** Зафиксировать стабильные id из дампа для тестов.

**Описание:**

Зависимость: сначала закрой задачу **S1-54**.

**Шаги:**

1. В `airport_repo` integration используй `testutil.SetupDB` вместо localhost (или оставь оба режима через env `TEST_DB_MODE`)
2. README обнови

**Результат:** `make test-integration` → PASS

**Документация:** —

---

## S1-CP

**Задача:** S1-CP — Конец спринта 1

**Цель:** Убедиться, что весь спринт собран: тесты, curl, gate-критерии.

**Описание:**

Прогони команды по порядку. Все должны завершиться без ошибок.

```bash
make test-unit
make test-integration
docker compose up -d
make run
curl localhost:8080/health
curl localhost:8080/ready
curl "localhost:8080/api/v1/airports?limit=5"
```

**Результат (gate):** ≥5 unit + ≥3 httptest + ≥1 integration + airports curl.

**Ревью:** `ревью S1-CP`

**Документация:** PROJECT.md §10 Definition of Done

---

# Sprint 2 — Read API

**Итог:** все GET из §7 (кроме уже готовых airports), mockgen, error mapper.

| Блок              | ID            | Задач |
| ----------------- | ------------- | ----- |
| A. Domain         | S2-01 … S2-06 | 6     |
| B. Flights search | S2-07 … S2-18 | 12    |
| C. Flight by ID   | S2-19 … S2-24 | 6     |
| D. Booking by ref | S2-25 … S2-32 | 8     |
| E. Errors HTTP    | S2-33 … S2-37 | 5     |
| F. mockgen        | S2-38 … S2-42 | 5     |

**Checkpoint:** [S2-CP](#s2-cp)

---

## A. Domain (Sprint 2)

### S2-01

**Задача:** S2-01 — domain/flight.go

**Цель:** Подключить модуль к `cmd/api/main.go`.

**Описание:**

Зависимость: завершён предыдущий спринт (`S1-CP`).

**Шаги:**

1. Создай `internal/domain/flight.go` с полями по `bookings.flights` + route airports
2. `go build ./internal/domain`

**Результат:** ok

**Документация:** —

---

### S2-02

**Задача:** S2-02 — domain/booking.go

**Цель:** Подключить модуль к `cmd/api/main.go`.

**Описание:**

Зависимость: сначала закрой задачу **S2-01**.

**Шаги:**

1. `booking.go`, `ticket.go`, `segment.go` — nested структура для read
2. Без тегов json в domain

**Результат:** ok

**Документация:** —

---

### S2-03

**Задача:** S2-03 — errors.go sentinels

**Цель:** Шаг к MVP: errors.go sentinels.

**Описание:**

Зависимость: сначала закрой задачу **S2-02**.

**Шаги:**

1. `var ErrNotFound = errors.New("not found")`, `ErrValidation`
2. Комментарий: HTTP mapping в adapter, не здесь

**Результат:** ok

**Документация:** —

---

### S2-04

**Задача:** S2-04 — Test errors.Is

**Цель:** Автотест: Test errors.Is.

**Описание:**

Зависимость: сначала закрой задачу **S2-03**.

**Шаги:**

1. `errors_test.go`: wrap ErrNotFound → `errors.Is` true
2. `TestValidationError` для typed error (если добавишь struct)

**Результат:** `go test ./internal/domain -v` → PASS

**Документация:** [https://go.dev/blog/go1.13-errors](https://go.dev/blog/go1.13-errors)

---

### S2-05

**Задача:** S2-05 — testutil/constants.go

**Цель:** Завести пакет для общих test helpers.

**Описание:**

Зависимость: сначала закрой задачу **S2-04**.

**Шаги:**

1. В DBeaver найди: реальный `flight_id`, `book_ref`, рабочие `SVO`→`LED` + date
2. Запиши в `internal/testutil/constants.go` как константы

**Результат:** SQL по константам возвращает строки

**Документация:** —

---

### S2-06

**Задача:** S2-06 — port FlightRepository

**Цель:** Шаг к MVP: port FlightRepository.

**Описание:**

Зависимость: сначала закрой задачу **S2-05**.

**Шаги:**

1. `port/flight_repository.go`: `Search`, `GetByID`, `CountSearch`
2. `port/booking_repository.go`: `GetByRef`

**Результат:** `go build ./internal/port` → ok

**Документация:** —

---

## B. Flights search

### S2-07

**Задача:** S2-07 — SQL Search в DBeaver

**Цель:** Проверить SQL на дампе до Go-кода.

**Описание:**

Зависимость: сначала закрой задачу **S2-06**.

**Шаги:**

1. JOIN flights + routes, фильтр from/to/date
2. `EXPLAIN ANALYZE` — скрин или время в комментарии

**Результат:** ≥1 row на тестовых константах

**Документация:** —

---

### S2-08

**Задача:** S2-08 — FlightSearchFilter type

**Цель:** Главный read-сценарий: поиск рейсов по маршруту и дате.

**Описание:**

Зависимость: сначала закрой задачу **S2-07**.

**Шаги:**

1. В `port` или `service` объяви filter struct: From, To, Date, Limit, Offset
2. Документируй политику timezone в комментарии

**Результат:** компиляция

**Документация:** —

---

### S2-09

**Задача:** S2-09 — postgres Search

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S2-08**.

**Шаги:**

1. `flight_repo.go` — метод `Search(ctx, filter)`
2. Параметры только `$1..$n`

**Результат:** integration `TestFlightRepo_Search` len>0

**Документация:** —

---

### S2-10

**Задача:** S2-10 — postgres CountSearch

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S2-09**.

**Шаги:**

1. `CountSearch` с теми же фильтрами
2. integration: count >= len(items)

**Результат:** PASS

**Документация:** —

---

### S2-11

**Задача:** S2-11 — FlightService.Search validate

**Цель:** Главный read-сценарий: поиск рейсов по маршруту и дате.

**Описание:**

Зависимость: сначала закрой задачу **S2-10**.

**Шаги:**

1. `flight_service.go` — валидация IATA len=3, date parse `2006-01-02`
2. limit default 20, max 100

**Результат:** unit `TestSearch_EmptyFrom` → ErrValidation

**Документация:** —

---

### S2-12

**Задача:** S2-12 — unit Search happy path

**Цель:** Шаг к MVP: unit Search happy path.

**Описание:**

Зависимость: сначала закрой задачу **S2-11**.

**Шаги:**

1. stub repo возвращает 2 flights
2. `TestSearch_HappyPath` → 2 items

**Результат:** PASS

**Документация:** —

---

### S2-13

**Задача:** S2-13 — unit Search repo error

**Цель:** Шаг к MVP: unit Search repo error.

**Описание:**

Зависимость: сначала закрой задачу **S2-12**.

**Шаги:**

1. stub repo → error
2. service возвращает wrapped err

**Результат:** PASS

**Документация:** —

---

### S2-14

**Задача:** S2-14 — dto flights response

**Цель:** Шаг к MVP: dto flights response.

**Описание:**

Зависимость: сначала закрой задачу **S2-12**.

**Шаги:**

1. `dto/flight.go` — items, total, limit, offset
2. mapper domain → dto

**Результат:** компиляция

**Документация:** —

---

### S2-15

**Задача:** S2-15 — handler GET /flights

**Цель:** Шаг к MVP: handler GET /flights.

**Описание:**

Зависимость: сначала закрой задачу **S2-14**.

**Шаги:**

1. Парсинг query: from, to, date, limit, offset
2. Вызов service.Search

**Результат:** register route

**Документация:** —

---

### S2-16

**Задача:** S2-16 — httptest flights 200

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S2-15**.

**Шаги:**

1. mock service → 200 + JSON items
2. subtest с query string

**Результат:** PASS

**Документация:** —

---

### S2-17

**Задача:** S2-17 — httptest missing date 400

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S2-16**.

**Шаги:**

1. запрос без `date` → 400 + error code

**Результат:** PASS

**Документация:** —

---

### S2-18

**Задача:** S2-18 — curl flights

**Цель:** Зафиксировать стабильные id из дампа для тестов.

**Описание:**

Зависимость: сначала закрой задачу **S2-17**.

**Шаги:**

1. curl с константами из testutil
2. README пример

**Результат:** 200 + items

**Документация:** —

---

## C. Flight by ID

### S2-19

**Задача:** S2-19 — postgres GetByID

**Цель:** Один пул соединений к PostgreSQL на всё приложение.

**Описание:**

Зависимость: сначала закрой задачу **S2-18**.

**Шаги:**

1. `GetByID(ctx, id)` — pgx.ErrNoRows → domain.ErrNotFound
2. integration с `testutil.FlightID`

**Результат:** PASS

**Документация:** —

---

### S2-20

**Задача:** S2-20 — FlightService.GetByID

**Цель:** Шаг к MVP: FlightService.GetByID.

**Описание:**

Зависимость: сначала закрой задачу **S2-19**.

**Шаги:**

1. invalid id ≤0 → validation
2. unit: not found, happy

**Результат:** `go test ./internal/service -run GetByID -v`

**Документация:** —

---

### S2-21

**Задача:** S2-21 — handler GET /flights/{id}

**Цель:** Шаг к MVP: handler GET /flights/{id}.

**Описание:**

Зависимость: сначала закрой задачу **S2-20**.

**Шаги:**

1. chi URL param `flight_id`
2. map 404/400

**Результат:** route registered

**Документация:** —

---

### S2-22

**Задача:** S2-22 — httptest 200/404/400

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S2-21**.

**Шаги:**

1. три subtests: ok id, missing, `abc`

**Результат:** PASS

**Документация:** —

---

### S2-23

**Задача:** S2-23 — curl flight by id

**Цель:** Ручная проверка сценария: curl flight by id.

**Описание:**

Зависимость: сначала закрой задачу **S2-22**.

**Шаги:**

1. curl существующий id
2. curl 99999999 → 404

**Результат:** статусы верные

**Документация:** —

---

### S2-24

**Задача:** S2-24 — integration empty search

**Цель:** Шаг к MVP: integration empty search.

**Описание:**

Зависимость: сначала закрой задачу **S2-10**.

**Шаги:**

1. Search с несуществующим маршрутом → 0 rows, не error

**Результат:** PASS

**Документация:** —

---

## D. Booking by ref

### S2-25

**Задача:** S2-25 — README SQL strategy

**Цель:** Онбординг: другой разработчик поднимает проект за 10 минут.

**Описание:**

Зависимость: сначала закрой задачу **S2-24**.

**Шаги:**

1. В README опиши: 2 SQL vs 1 json_agg (1 абзац trade-off)
2. Выбери стратегию для реализации

**Результат:** текст в README

**Документация:** —

---

### S2-26

**Задача:** S2-26 — postgres GetByRef header

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S2-25**.

**Шаги:**

1. SQL booking by `book_ref`
2. integration: known ref exists

**Результат:** PASS

**Документация:** —

---

### S2-27

**Задача:** S2-27 — postgres tickets+segments

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S2-26**.

**Шаги:**

1. Второй запрос(и) tickets + segments для ref
2. Собери `domain.Booking` в Go

**Результат:** integration полная структура

**Документация:** —

---

### S2-28

**Задача:** S2-28 — BookingService.GetByRef

**Цель:** Чтение агрегата бронирования по book_ref.

**Описание:**

Зависимость: сначала закрой задачу **S2-27**.

**Шаги:**

1. service wrapper
2. unit: not found

**Результат:** PASS

**Документация:** —

---

### S2-29

**Задача:** S2-29 — handler GET /bookings/{ref}

**Цель:** Чтение агрегата бронирования по book_ref.

**Описание:**

Зависимость: сначала закрой задачу **S2-28**.

**Шаги:**

1. URL param book_ref length 6
2. dto response

**Результат:** route ok

**Документация:** —

---

### S2-30

**Задача:** S2-30 — httptest booking 200

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S2-29**.

**Шаги:**

1. mock service → JSON с tickets array

**Результат:** PASS

**Документация:** —

---

### S2-31

**Задача:** S2-31 — httptest booking 404

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S2-30**.

**Шаги:**

1. not found → 404 code BOOKING_NOT_FOUND

**Результат:** PASS

**Документация:** —

---

### S2-32

**Задача:** S2-32 — curl booking

**Цель:** Зафиксировать стабильные id из дампа для тестов.

**Описание:**

Зависимость: сначала закрой задачу **S2-31**.

**Шаги:**

1. curl real book_ref из testutil

**Результат:** 200

**Документация:** —

---

## E. Error mapper

### S2-33

**Задача:** S2-33 — http/errors.go mapper

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S2-32**.

**Шаги:**

1. `MapError(err) (status, code, message)`
2. table in comment

**Результат:** unit tests mapper 3 cases

**Документация:** —

---

### S2-34

**Задача:** S2-34 — writeJSON writeError

**Цель:** Шаг к MVP: writeJSON writeError.

**Описание:**

Зависимость: сначала закрой задачу **S2-33**.

**Шаги:**

1. helpers в `response.go`
2. Content-Type application/json

**Результат:** unit: writeError NotFound → body code

**Документация:** —

---

### S2-35

**Задача:** S2-35 — подключить mapper ко всем handlers

**Цель:** Единый маппинг domain error → HTTP JSON.

**Описание:**

Зависимость: сначала закрой задачу **S2-34**.

**Шаги:**

1. замени прямые http.Error на writeError
2. пройдись по airports/flights/bookings handlers

**Результат:** `make test-unit` green

**Документация:** —

---

### S2-36

**Задача:** S2-36 — Test internal error no leak

**Цель:** Автотест: Test internal error no leak.

**Описание:**

Зависимость: сначала закрой задачу **S2-35**.

**Шаги:**

1. handler test: unknown error → 500, body без подстроки `pq:`

**Результат:** PASS

**Документация:** —

---

### S2-37

**Задача:** S2-37 — Test validation body

**Цель:** Автотест: Test validation body.

**Описание:**

Зависимость: сначала закрой задачу **S2-36**.

**Шаги:**

1. validation → 400 + VALIDATION_ERROR

**Результат:** PASS

**Документация:** —

---

## F. mockgen

### S2-38

**Задача:** S2-38 — go:generate mockgen

**Цель:** Генерировать моки репозиториев для unit service.

**Описание:**

Зависимость: сначала закрой задачу **S2-37**.

**Шаги:**

1. `//go:generate` в `port` для Flight + Booking repos
2. `make generate` в Makefile

**Результат:** `make generate` создаёт mocks

**Документация:** [https://github.com/uber-go/mock](https://github.com/uber-go/mock)

---

### S2-39

**Задача:** S2-39 — refactor service tests gomock

**Цель:** Шаг к MVP: refactor service tests gomock.

**Описание:**

Зависимость: сначала закрой задачу **S2-38**.

**Шаги:**

1. перепиши один тест Search на gomock EXPECT
2. `go test ./internal/service -v`

**Результат:** PASS

**Документация:** —

---

### S2-40

**Задача:** S2-40 — покрыть все методы FlightService

**Цель:** Шаг к MVP: покрыть все методы FlightService.

**Описание:**

Зависимость: сначала закрой задачу **S2-39**.

**Шаги:**

1. тест на каждый публичный метод (Search, GetByID)

**Результат:** `make test-cover` service ≥ 60%

**Документация:** —

---

### S2-41

**Задача:** S2-41 — покрыть BookingService read

**Цель:** Шаг к MVP: покрыть BookingService read.

**Описание:**

Зависимость: сначала закрой задачу **S2-40**.

**Шаги:**

1. GetByRef: ok + not found

**Результат:** cover растёт

**Документация:** —

---

### S2-42

**Задача:** S2-42 — README coverage

**Цель:** Онбординг: другой разработчик поднимает проект за 10 минут.

**Описание:**

Зависимость: сначала закрой задачу **S2-41**.

**Шаги:**

1. документируй `make test-cover`

**Результат:** —

**Документация:** —

---

## S2-CP

**Задача:** S2-CP — Конец спринта 2

**Цель:** Убедиться, что весь спринт собран: тесты, curl, gate-критерии.

**Описание:**

Прогони команды по порядку. Все должны завершиться без ошибок.

```bash
make test-cover
curl "localhost:8080/api/v1/flights?from=SVO&to=LED&date=2017-07-14"
curl localhost:8080/api/v1/flights/<FLIGHT_ID>
curl localhost:8080/api/v1/bookings/<BOOK_REF>
```

**Результат (gate):** service ≥60%, ≥6 httptest на GET, 3 integration repo tests.

**Документация:** PROJECT.md §10 Definition of Done

---

# Sprint 3 — Write path

**Итог:** `POST /bookings`, транзакция, rollback test, e2e.

| Блок           | ID            |
| -------------- | ------------- |
| A. ID gen      | S3-01 … S3-04 |
| B. Validation  | S3-05 … S3-10 |
| C. Transaction | S3-11 … S3-17 |
| D. HTTP POST   | S3-18 … S3-23 |
| E. E2E         | S3-24 … S3-26 |

---

### S3-01

**Задача:** S3-01 — port IDGenerator

**Цель:** Шаг к MVP: port IDGenerator.

**Описание:**

Зависимость: завершён предыдущий спринт (`S2-CP`).

**Шаги:**

1. `port/id_generator.go` — BookRef(), TicketNo()
2. Изучи в DBeaver формат существующих ref (длина, символы)

**Результат:** ok

**Документация:** —

---

### S3-02

**Задача:** S3-02 — adapter generator

**Цель:** Шаг к MVP: adapter generator.

**Описание:**

Зависимость: сначала закрой задачу **S3-01**.

**Шаги:**

1. `internal/adapter/idgen/generator.go` на `crypto/rand`
2. Не используй math/rand

**Результат:** ok

**Документация:** [https://pkg.go.dev/crypto/rand](https://pkg.go.dev/crypto/rand)

---

### S3-03

**Задача:** S3-03 — unit TestBookRef format

**Цель:** Шаг к MVP: unit TestBookRef format.

**Описание:**

Зависимость: сначала закрой задачу **S3-02**.

**Шаги:**

1. Test: len=6, alphanumeric
2. 100 итераций — все уникальны (map)

**Результат:** PASS

**Документация:** —

---

### S3-04

**Задача:** S3-04 — unit TicketNo format

**Цель:** Шаг к MVP: unit TicketNo format.

**Описание:**

Зависимость: сначала закрой задачу **S3-03**.

**Шаги:**

1. Test длины/формата ticket_no по дампу

**Результат:** PASS

**Документация:** —

---

### S3-05

**Задача:** S3-05 — CreateBookingCommand struct

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-04**.

**Шаги:**

1. `service/booking_command.go` — input struct + `Validate()` method
2. Правила в комментарии

**Результат:** unit Validate empty passengers

**Документация:** —

---

### S3-06

**Задача:** S3-06 — Validate fare + price

**Цель:** Шаг к MVP: Validate fare + price.

**Описание:**

Зависимость: сначала закрой задачу **S3-05**.

**Шаги:**

1. tests: invalid fare, negative price

**Результат:** PASS

**Документация:** —

---

### S3-07

**Задача:** S3-07 — port ExistsFlight / GetFlightForBooking

**Цель:** Шаг к MVP: port ExistsFlight / GetFlightForBooking.

**Описание:**

Зависимость: сначала закрой задачу **S3-06**.

**Шаги:**

1. метод repo для проверки flight exists + status
2. integration: cancelled flight id → статус Cancelled

**Результат:** PASS

**Документация:** —

---

### S3-08

**Задача:** S3-08 — service rule flight not found

**Цель:** Шаг к MVP: service rule flight not found.

**Описание:**

Зависимость: сначала закрой задачу **S3-07**.

**Шаги:**

1. unit: repo → ErrNotFound → domain ErrNotFound

**Результат:** PASS

**Документация:** —

---

### S3-09

**Задача:** S3-09 — service rule cancelled

**Цель:** Шаг к MVP: service rule cancelled.

**Описание:**

Зависимость: сначала закрой задачу **S3-08**.

**Шаги:**

1. unit: repo returns cancelled → ErrFlightNotAvailable (добавь sentinel)

**Результат:** PASS

**Документация:** —

---

### S3-10

**Задача:** S3-10 — service total_amount

**Цель:** Шаг к MVP: service total_amount.

**Описание:**

Зависимость: сначала закрой задачу **S3-09**.

**Шаги:**

1. Validate sum(segment prices) == total
2. unit mismatch → validation

**Результат:** PASS

**Документация:** —

---

### S3-11

**Задача:** S3-11 — port CreateBooking method

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-10**.

**Шаги:**

1. `BookingRepository.Create(ctx, booking)` в port
2. domain booking aggregate для write

**Результат:** ok

**Документация:** —

---

### S3-12

**Задача:** S3-12 — postgres BeginTx skeleton

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-11**.

**Шаги:**

1. `booking_repo.go` Create: Begin, defer Rollback
2. пустые Exec пока — компилируется

**Результат:** ok

**Документация:** [https://github.com/jackc/pgx/wiki/Transactions](https://github.com/jackc/pgx/wiki/Transactions)

---

### S3-13

**Задача:** S3-13 — INSERT bookings

**Цель:** Шаг к MVP: INSERT bookings.

**Описание:**

Зависимость: сначала закрой задачу **S3-12**.

**Шаги:**

1. первый Exec INSERT bookings
2. integration: только bookings row (временно)

**Результат:** row in DBeaver

**Документация:** —

---

### S3-14

**Задача:** S3-14 — INSERT tickets loop

**Цель:** Шаг к MVP: INSERT tickets loop.

**Описание:**

Зависимость: сначала закрой задачу **S3-13**.

**Шаги:**

1. batch/loop INSERT tickets
2. integration: bookings + tickets

**Результат:** PASS

**Документация:** —

---

### S3-15

**Задача:** S3-15 — INSERT segments

**Цель:** Шаг к MVP: INSERT segments.

**Описание:**

Зависимость: сначала закрой задачу **S3-14**.

**Шаги:**

1. INSERT segments
2. integration: full create

**Результат:** PASS

**Документация:** —

---

### S3-16

**Задача:** S3-16 — Commit + return

**Цель:** Шаг к MVP: Commit + return.

**Описание:**

Зависимость: сначала закрой задачу **S3-15**.

**Шаги:**

1. Commit, return domain booking
2. unit integration happy POST data

**Результат:** PASS

**Документация:** —

---

### S3-17

**Задача:** S3-17 — Test rollback

**Цель:** Запись в несколько таблиц атомарно (ACID).

**Описание:**

Зависимость: сначала закрой задачу **S3-16**.

**Шаги:**

1. тест: намеренно broken segment (invalid flight_id) → error
2. assert: нет bookings row с этим book_ref

**Результат:** PASS — **критичный тест спринта**

**Документация:** —

---

### S3-18

**Задача:** S3-18 — dto POST request

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-17**.

**Шаги:**

1. `dto/booking_create.go` request/response
2. mapper command

**Результат:** ok

**Документация:** —

---

### S3-19

**Задача:** S3-19 — handler POST /bookings

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-18**.

**Шаги:**

1. decode JSON, call service.Create
2. 201 + body

**Результат:** compile

**Документация:** —

---

### S3-20

**Задача:** S3-20 — httptest 201

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S3-19**.

**Шаги:**

1. mock service → 201

**Результат:** PASS

**Документация:** —

---

### S3-21

**Задача:** S3-21 — httptest bad JSON 400

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S3-20**.

**Шаги:**

1. body `{` → 400

**Результат:** PASS

**Документация:** —

---

### S3-22

**Задача:** S3-22 — httptest MaxBytesReader

**Цель:** Поднять HTTP-слой и probes для оркестратора.

**Описание:**

Зависимость: сначала закрой задачу **S3-21**.

**Шаги:**

1. body > 1MB → 413

**Результат:** PASS

**Документация:** —

---

### S3-23

**Задача:** S3-23 — curl POST manual

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-22**.

**Шаги:**

1. `testdata/booking_valid.json`
2. curl POST → 201, проверь в DBeaver

**Результат:** row exists

**Документация:** —

---

### S3-24

**Задача:** S3-24 — testdata + constants flight_id

**Цель:** Зафиксировать стабильные id из дампа для тестов.

**Описание:**

Зависимость: сначала закрой задачу **S3-23**.

**Шаги:**

1. в JSON используй flight_id из testutil
2. document in README

**Результат:** ok

**Документация:** —

---

### S3-25

**Задача:** S3-25 — e2e CreateAndGet

**Цель:** Write-path: создание бронирования через API.

**Описание:**

Зависимость: сначала закрой задачу **S3-24**.

**Шаги:**

1. `test/integration/e2e_booking_test.go`: POST router + real DB + GET same ref

**Результат:** `make test-integration` PASS

**Документация:** —

---

### S3-26

**Задача:** S3-26 — go test -race

**Цель:** Проверка качества перед финальным ревью.

**Описание:**

Зависимость: сначала закрой задачу **S3-25**.

**Шаги:**

1. `go test -race ./internal/...`
2. исправь гонки если есть

**Результат:** race-free

**Документация:** —

---

## S3-CP

**Задача:** S3-CP — Конец спринта 3

**Цель:** Убедиться, что весь спринт собран: тесты, curl, gate-критерии.

**Описание:**

Прогони команды по порядку. Все должны завершиться без ошибок.

```bash
make test
curl -X POST localhost:8080/api/v1/bookings -H 'Content-Type: application/json' -d @testdata/booking_valid.json
```

**Результат (gate):** rollback test S3-17 + e2e S3-25 + 4 httptest POST.

**Документация:** PROJECT.md §10 Definition of Done

---

# Sprint 3 — Write path

**Итог:** POST /bookings, транзакция, rollback test, e2e.

| Блок           | ID            |
| -------------- | ------------- |
| A. ID gen      | S3-01 … S3-04 |
| B. Validation  | S3-05 … S3-10 |
| C. Transaction | S3-11 … S3-17 |
| D. HTTP POST   | S3-18 … S3-23 |
| E. E2E         | S3-24 … S3-26 |

**Checkpoint:** [S3-CP](#s3-cp)

---

# Sprint 4 — Production-ready

| Блок        | ID            |
| ----------- | ------------- |
| A. Shutdown | S4-01 … S4-04 |
| B. Metrics  | S4-05 … S4-08 |
| C. Docker   | S4-09 … S4-12 |
| D. Docs/CI  | S4-13 … S4-18 |

---

### S4-01

**Задача:** S4-01 — signal.NotifyContext в main

**Цель:** Корректное завершение при SIGTERM.

**Описание:**

Зависимость: завершён предыдущий спринт (`S3-CP`).

**Шаги:**

1. замени Background на signal context SIGINT/SIGTERM
2. лог "shutdown started"

**Результат:** `make run`, Ctrl+C — чистый exit

**Документация:** [https://pkg.go.dev/os/signal](https://pkg.go.dev/os/signal)

---

### S4-02

**Задача:** S4-02 — Server.Shutdown

**Цель:** Корректное завершение при SIGTERM.

**Описание:**

Зависимость: сначала закрой задачу **S4-01**.

**Шаги:**

1. Shutdown с timeout 10s
2. затем pool.Close()

**Результат:** manual: sleep handler 5s + Ctrl+C — дождаться завершения

**Документация:** —

---

### S4-03

**Задача:** S4-03 — prometheus dep

**Цель:** Наблюдаемость: метрики для Prometheus.

**Описание:**

Зависимость: сначала закрой задачу **S4-02**.

**Шаги:**

1. `go get github.com/prometheus/client_golang/prometheus/promhttp`
2. пустой handler `/metrics`

**Результат:** curl /metrics → text

**Документация:** —

---

### S4-04

**Задача:** S4-04 — middleware metrics

**Цель:** Наблюдаемость: метрики для Prometheus.

**Описание:**

Зависимость: сначала закрой задачу **S4-03**.

**Шаги:**

1. Counter + Histogram duration
2. labels: method, pattern (не raw path)

**Результат:** httptest metrics after 3 requests

**Документация:** —

---

### S4-05

**Задача:** S4-05 — slog JSON env

**Цель:** Конфигурация только из env — готово к Docker и 12-factor.

**Описание:**

Зависимость: сначала закрой задачу **S4-02**.

**Шаги:**

1. `LOG_FORMAT=json` → json handler
2. test: default text

**Результат:** log line parseable

**Документация:** —

---

### S4-06

**Задача:** S4-06 — Dockerfile build stage

**Цель:** Запуск API в контейнере рядом с PostgreSQL.

**Описание:**

Зависимость: сначала закрой задачу **S4-05**.

**Шаги:**

1. Dockerfile multi-stage
2. `docker build -t bookings-api .`

**Результат:** image builds

**Документация:** —

---

### S4-07

**Задача:** S4-07 — compose service api

**Цель:** Шаг к MVP: compose service api.

**Описание:**

Зависимость: сначала закрой задачу **S4-06**.

**Шаги:**

1. `docker-compose.yml` сервис `api`, depends_on postgres healthy
2. env DB_HOST=postgres

**Результат:** `docker compose up --build` + curl health

**Документация:** —

---

### S4-08

**Задача:** S4-08 — TESTING.md

**Цель:** Зафиксировать стабильные id из дампа для тестов.

**Описание:**

Зависимость: сначала закрой задачу **S4-07**.

**Шаги:**

1. создай `docs/TESTING.md`: pyramid, make targets, tags, testutil constants
2. ссылка из README

**Результат:** файл есть

**Документация:** —

---

### S4-09

**Задача:** S4-09 — make test-cover в CI-local

**Цель:** Автотест: make test-cover в CI-local.

**Описание:**

Зависимость: сначала закрой задачу **S4-08**.

**Шаги:**

1. `make test-cover` ≥60% domain+service
2. запиши % в README

**Результат:** cover OK

**Документация:** —

---

### S4-10

**Задача:** S4-10 — staticcheck

**Цель:** Проверка качества перед финальным ревью.

**Описание:**

Зависимость: сначала закрой задачу **S4-09**.

**Шаги:**

1. `go install staticcheck@latest`
2. `staticcheck ./...` — 0 critical

**Результат:** clean

**Документация:** —

---

### S4-11

**Задача:** S4-11 — README все curl

**Цель:** Онбординг: другой разработчик поднимает проект за 10 минут.

**Описание:**

Зависимость: сначала закрой задачу **S4-10**.

**Шаги:**

1. все эндпоинты §7 с примерами
2. mermaid/architecture diagram

**Результат:** copy-paste works

**Документация:** —

---

### S4-12

**Задача:** S4-12 — vet + lint Makefile

**Цель:** Единые команды `make test-unit` / `test-integration` для всего репо.

**Описание:**

Зависимость: сначала закрой задачу **S4-11**.

**Шаги:**

1. `make lint` → go vet + staticcheck

**Результат:** make lint green

**Документация:** —

---

### S4-13

**Задача:** S4-13 — integration в compose profile

**Цель:** Шаг к MVP: integration в compose profile.

**Описание:**

Зависимость: сначала закрой задачу **S4-07 (опц.)**.

**Шаги:**

1. profile `test` в compose для CI-like run

**Результат:** documented

**Документация:** —

---

### S4-14

**Задача:** S4-14 — .github/workflows stub (опц.)

**Цель:** Шаг к MVP: .github/workflows stub (опц.).

**Описание:**

Зависимость: сначала закрой задачу **S4-12**.

**Шаги:**

1. workflow: checkout, go test unit, integration with service container

**Результат:** yaml valid

**Документация:** —

---

### S4-15

**Задача:** S4-15 — финальный self-review DoD

**Цель:** Закрыть DoD всего проекта.

**Описание:**

Зависимость: сначала закрой задачу **S4-12**.

**Шаги:**

1. пройди PROJECT.md §10 чек-лист
2. отметь [x] в копии или issue

**Результат:** все пункты

**Документация:** —

---

### S4-16

**Задача:** S4-16 — ревью S4-CP

**Цель:** Закрыть DoD всего проекта.

**Описание:**

Зависимость: сначала закрой задачу **S4-15**.

**Шаги:**

1. запрос `ревью S4-CP`
2. 5 тезисов для собеса запиши в README

**Результат:** ментор approve

**Документация:** —

---

## S4-CP

```bash
make lint
make test
make test-cover
docker compose up --build
# все curl из README
```

---

# Сводная таблица (все ID)

| Sprint    | Задач    | Первый ID | Checkpoint |
| --------- | -------- | --------- | ---------- |
| 1         | 55       | S1-01     | S1-CP      |
| 2         | 42       | S2-01     | S2-CP      |
| 3         | 26       | S3-01     | S3-CP      |
| 4         | 16       | S4-01     | S4-CP      |
| **Всего** | **~139** |           |            |

---

_Рекомендуемый порядок: **S0-01…S0-04** (git) → **S1-01** (testify) → …_
