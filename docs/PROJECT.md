# Avia Bookings API — ТЗ и план разработки

> **Версия:** 3.1  
> **Проект:** pet-проект для портфолио и подготовки к Middle Go Developer  
> **База:** PostgreSQL `demo`, схема `bookings`  
> **Срок:** 4 спринта × 3–4 дня  
> **Режим:** код → тесты → self-review → ревью ментора

---

## Оглавление

1. [Видение и scope](#1-видение-и-scope)
2. [Цели обучения](#2-цели-обучения)
3. [Стек](#3-стек)
4. [Архитектура и layout](#4-архитектура-и-layout)
5. [Доменная модель БД](#5-доменная-модель-бд)
6. [Миграции и дамп](#6-миграции-и-дамп)
7. [API v1](#7-api-v1)
8. [Нефункциональные требования](#8-нефункциональные-требования)
9. [**Стратегия тестирования**](#9-стратегия-тестирования)
10. [Definition of Done](#10-definition-of-done)
11. [Процесс и ревью](#11-процесс-и-ревью)
12. [Библиография](#12-библиография-мастер-список)
13. [Микрозадачи — как работать](#13-микрозадачи-как-работать)
14. [**Микрозадачи (TASKS.md)**](./TASKS.md) ← **основной рабочий план**
15. [Карта спринтов](#15-карта-спринтов)
16. [Backlog v2](#16-backlog-v2)

---

## 1. Видение и scope

**Avia Bookings API** — HTTP-сервис: поиск рейсов, детали, создание и просмотр бронирований.

**In scope v1:** эндпоинты §7, hexagonal architecture, unit + integration + handler tests, observability.

**Out of scope v1:** JWT, оплата, фронт, отмена брони, ORM, изменение схемы `bookings.*`.

---

## 2. Цели обучения

| Область | Артефакт в репо |
|---------|-----------------|
| Go | context, `%w`, interfaces, table-driven tests, race-free code |
| HTTP | chi, middleware, `httptest`, JSON, status codes |
| PostgreSQL | pgxpool, JOIN, транзакции, integration tests |
| Архитектура | domain / port / service / adapter |
| **Тестирование** | ≥60% coverage `service`+`domain`, integration suite, Makefile targets |
| Ops | slog, graceful shutdown, `/metrics`, Docker |

---

## 3. Стек

| Компонент | Пакет / инструмент |
|-----------|-------------------|
| Go | 1.22+ |
| HTTP | `github.com/go-chi/chi/v5` |
| DB | `github.com/jackc/pgx/v5/pgxpool` |
| Unit assertions | `github.com/stretchr/testify` (`require`, `assert`) |
| Mocks | `go.uber.org/mock/mockgen` |
| Integration DB | `github.com/testcontainers/testcontainers-go` + модуль `postgres` |
| HTTP test | `net/http/httptest` (stdlib) |
| Метрики | `github.com/prometheus/client_golang` |
| Линт (опц.) | `staticcheck`, `golangci-lint` |
| Конфиг | `github.com/kelseyhightower/envconfig` |

---

## 4. Архитектура и layout

```
bookings/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── port/
│   │   └── mocks/              # mockgen output
│   ├── service/
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── dto/
│   │   └── postgres/
│   └── testutil/               # helpers для тестов (фикстуры, test DB)
├── test/
│   └── integration/            # integration + e2e (build tag integration)
├── docs/PROJECT.md
├── Makefile                    # test, test-integration, lint, run
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

**Правило:** `domain` не импортирует `adapter`. Тесты лежат рядом с кодом (`*_test.go`) + `test/integration` для тяжёлых сценариев.

---

## 5. Доменная модель БД

| Таблица | Роль |
|---------|------|
| `airports_data` | IATA справочник |
| `routes` | Маршрут (from/to, airplane, schedule) |
| `flights` | Рейс + status + timestamps |
| `bookings` | `book_ref`, `total_amount` |
| `tickets` | Пассажир |
| `segments` | ticket ↔ flight + price + class |

Статусы рейса: `Scheduled`, `On Time`, `Delayed`, `Boarding`, `Departed`, `Arrived`, `Cancelled`.  
Классы: `Economy`, `Comfort`, `Business`.

---

## 6. Миграции и дамп

- Дамп `demo-20250901-3m.sql` — единственный источник `bookings.*`.
- Свои миграции только для `app.*` (опционально, Sprint 4).
- **Запрещено** ломать учебную схему в PR.

---

## 7. API v1

Base: `/api/v1`

| Method | Path | Описание |
|--------|------|----------|
| GET | `/health` | Liveness |
| GET | `/ready` | Readiness (+ DB ping) |
| GET | `/airports` | Список аэропортов |
| GET | `/flights` | Поиск (`from`, `to`, `date`, `limit`, `offset`) |
| GET | `/flights/{flight_id}` | Детали |
| GET | `/bookings/{book_ref}` | Бронирование |
| POST | `/bookings` | Создание |

Ошибки JSON: `{ "error": { "code": "...", "message": "..." } }`  
Коды: `VALIDATION_ERROR`, `FLIGHT_NOT_FOUND`, `BOOKING_NOT_FOUND`, `FLIGHT_NOT_AVAILABLE`, `INTERNAL_ERROR`.

Детальные примеры request/response — в README (Task 4.5), контракт зафиксируй при реализации Sprint 2.

---

## 8. Нефункциональные требования

| Требование | Критерий |
|------------|----------|
| Тесты | `make test` зелёный; integration не skipped в CI-local |
| Coverage | `service` + `domain` ≥ 60% (см. §10) |
| SQL | только `$N`, без конкатенации |
| Shutdown | SIGTERM → drain HTTP → close pool |
| Логи | slog + request_id |

---

## 9. Стратегия тестирования

### 9.1 Пирамида тестов

```
        ┌─────────────┐
        │  E2E (мало) │  httptest полный router + реальная БД (integration)
        ├─────────────┤
        │ Integration │  postgres repo, testcontainers
        ├─────────────┤
        │  Unit (много)│  service + domain + mapper, моки port
        └─────────────┘
```

| Уровень | Что тестируем | Инструменты | Скорость |
|---------|---------------|-------------|----------|
| **Unit** | Бизнес-логика, валидация, маппинг ошибок | `testing`, `testify`, `mockgen` | мс |
| **Handler** | HTTP коды, JSON, middleware | `httptest.NewRecorder`, `httptest.NewRequest` | мс |
| **Integration** | SQL, транзакции, реальный PG | `testcontainers` или `TEST_DATABASE_URL` | сек–мин |
| **Manual** | Полный сценарий, DBeaver | curl, DBeaver | — |

### 9.2 Соглашения в репозитории

| Соглашение | Правило |
|------------|---------|
| Имя файла | `flight_service_test.go` рядом с `flight_service.go` |
| Table-driven | `tests := []struct{ name string; ...}{}` + `t.Run(tt.name, ...)` |
| Arrange-Act-Assert | комментарии `// arrange` допустимы |
| Моки | генерировать из `port/*.go`, не мокать pgx в unit |
| Integration tag | `//go:build integration` в `test/integration/` |
| Параллельность | `t.Parallel()` только где нет shared DB |
| Cleanup | `t.Cleanup()` для pool, containers |
| Flaky | integration не зависит от порядка; уникальные `book_ref` |
| Short mode | `if testing.Short() { t.Skip() }` в integration |

### 9.3 Makefile (создать в TASK S1-01)

```makefile
.PHONY: test test-unit test-integration test-cover lint run

test-unit:
	go test ./internal/... -count=1 -race

test-integration:
	go test -tags=integration ./test/integration/... -count=1 -timeout=5m

test: test-unit test-integration

test-cover:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -func=coverage.out
```

### 9.4 Что мокать / не мокать

| Слой | Unit | Integration |
|------|------|-------------|
| `service` | mock `port.Repository` | optional real repo |
| `handler` | mock `service` или interface use case | real router + mock service |
| `postgres` | **не unit-тестить** без БД | real PG |
| `main` | не тестить | smoke в e2e |

### 9.5 Тестовые данные

- Для integration возьми **известные** `flight_id`, `book_ref` из дампа (запиши в `testutil/constants.go` после исследования DBeaver).
- Для POST создавай **уникальные** ref (не пересекайся с прод-подобными данными дампа).
- Не коммить `.env` с паролями.

### 9.6 Матрица тестов по спринтам

| Спринт | Минимум тестов к концу спринта |
|--------|--------------------------------|
| 1 | Makefile, unit config, httptest `/health` `/ready`, unit postgres ping helper |
| 2 | unit domain errors, service Search/Get, httptest всех GET, integration ListAirports/Search |
| 3 | unit CreateBooking rules, integration tx rollback, httptest POST, e2e create+get |
| 4 | race detector в CI-local, coverage report, integration в compose profile |

---

## 10. Definition of Done (весь проект)

**Функциональность**
- [ ] Все эндпоинты §7
- [ ] README + curl

**Тестирование**
- [ ] `make test` green
- [ ] `go test -race ./internal/...` без гонок
- [ ] Coverage `domain` + `service` ≥ 60% (`make test-cover`)
- [ ] ≥ 8 unit test functions в `service`
- [ ] ≥ 4 integration test functions с реальной PG
- [ ] ≥ 6 httptest handler tests
- [ ] Документ «Как запускать тесты» в README

**Качество**
- [ ] `go vet ./...`
- [ ] Ошибки с `%w`, domain без pgx/chi
- [ ] Docker multi-stage build

---

## 11. Процесс и ревью

### 11.1 Цикл задачи

1. Прочитать **Learn** + **Документацию**
2. Реализовать **Пошаговое задание**
3. Написать тесты из блока **Тестирование** (TDD допустим: red → green)
4. **Self-review** + `make test`
5. Запрос ревью: `ревью S1-12` (ID из TASKS.md)

### 11.2 Формат запроса на ревью

```
TASK: 2.3
Реализация: ...
Тесты: go test ./internal/service -run TestFlight -v
Покрытие: 72% service (go tool cover)
Вопросы: ...
```

### 11.3 Ревьюер проверяет

- Acceptance criteria
- **Все пункты блока «Тестирование»**
- Review criteria
- Нет готового «решению задачи» — только замечания

---

## 12. Библиография (мастер-список)

### Go — основы и стиль

- https://go.dev/doc/effective_go
- https://go.dev/doc/code
- https://go.dev/blog/context
- https://go.dev/blog/go1.13-errors
- https://go.dev/wiki/TableDrivenTests
- https://go.dev/wiki/RangeLoops — for и замыкания в тестах
- https://pkg.go.dev/testing

### Go — инструменты тестирования

- https://pkg.go.dev/net/http/httptest
- https://pkg.go.dev/testing#hdr-Main — TestMain для integration
- https://go.dev/blog/subtests
- https://github.com/stretchr/testify — assert vs require
- https://github.com/uber-go/mock — mockgen
- https://go.dev/doc/tutorial/add-a-test

### Go — качество

- https://staticcheck.dev/docs/
- https://go.dev/blog/race-detector
- https://pkg.go.dev/cmd/go#hdr-Test_coverage

### HTTP / chi

- https://github.com/go-chi/chi
- https://github.com/go-chi/chi#middlewares
- https://pkg.go.dev/net/http#Server

### PostgreSQL / pgx

- https://github.com/jackc/pgx/wiki/Getting-started-with-pgx
- https://github.com/jackc/pgx/wiki/Transactions
- https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool
- https://www.postgresql.org/docs/current/tutorial.html
- https://www.postgresql.org/docs/current/sql-explain.html
- https://www.postgresql.org/docs/current/functions-datetime.html

### Testcontainers

- https://golang.testcontainers.org/
- https://golang.testcontainers.org/modules/postgres/
- https://golang.testcontainers.org/features/wait/strategies/

### Архитектура

- https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- Ports & Adapters: https://herbertograca.com/2017/09/14/ports-adapters-architecture/

### Observability

- https://pkg.go.dev/log/slog
- https://prometheus.io/docs/guides/go-application/

### DevOps

- https://docs.docker.com/compose/
- https://docs.docker.com/build/building/multi-stage/
- https://12factor.net/config

### Книга (параллельно)

- Donovan & Kernighan — глава 7 (HTTP), 9 (goroutines), 11 (testing)

---

## 13. Микрозадачи (как работать)

**Полный план:** **[docs/TASKS.md](./TASKS.md)** — ~139 задач, **2–3 шага** каждая.

### Каркас одной микрозадачи (см. TASKS.md)

| Поле | Содержание |
|------|------------|
| **Задача** | ID + название (`S1-07 — Папки internal`) |
| **Цель** | Зачем этот шаг |
| **Описание** | Зависимость + 2–3 шага + пояснение |
| **Результат** | Команда → ожидаемый выход |
| **Документация** | Ссылки или «—» |

### Правила

1. Делай задачи **строго по ID** — не перескакивай.
2. Каждая задача заканчивается **Проверкой** (команда + ожидание).
3. Задача с тестом в шагах — **сначала тест red, потом green** (где уместно).
4. Ревью: `ревью S1-12` после проверки.
5. Checkpoint спринта (`S1-CP`) — gate перед следующим спринтом.

### Большая документация

- ТЗ, API, архитектура, тест-стратегия → **этот файл (PROJECT.md)**
- Пошаговые задачи → **TASKS.md**
- Как гонять тесты (Sprint 4) → **docs/TESTING.md** (создашь в S4-08)

---

## 15. Карта спринтов

| Спринт | Дней | Микрозадач | Checkpoint | Итог |
|--------|------|------------|------------|------|
| **0** Git (опц.) | 0.5–1 | S0-01 … S0-04 | S0-CP | git + .gitignore + remote |
| **1** Foundation | 3–4 | S1-01 … S1-55 | S1-CP | airports API + test infra |
| **2** Read | 3–4 | S2-01 … S2-42 | S2-CP | все GET + mockgen |
| **3** Write | 3–4 | S3-01 … S3-26 | S3-CP | POST + tx rollback + e2e |
| **4** Prod | 3–4 | S4-01 … S4-16 | S4-CP | Docker, metrics, DoD |

### Sprint 1 — блоки (детали в TASKS.md)

| Блок | ID | Суть |
|------|-----|------|
| A | S1-01…06 | testify, Makefile, первый тест |
| B | S1-07…11 | папки, `cmd/api`, build |
| C | S1-12…19 | config, DSN, тесты env, wire main |
| D | S1-20…26 | pgxpool, ping, integration |
| E | S1-27…38 | chi, health/ready, httptest |
| F | S1-39…52 | airports: domain→repo→service→handler |
| G | S1-53…55 | testcontainers (опц.) |

### Sprint 2 — блоки

| Блок | ID | Суть |
|------|-----|------|
| A | S2-01…06 | domain + errors tests + testutil constants |
| B | S2-07…18 | flights search + SQL + httptest |
| C | S2-19…24 | flight by id |
| D | S2-25…32 | booking by ref |
| E | S2-33…37 | error mapper |
| F | S2-38…42 | mockgen + coverage |

### Sprint 3 — блоки

| Блок | ID | Суть |
|------|-----|------|
| A | S3-01…04 | ID generator tests |
| B | S3-05…10 | validation unit tests |
| C | S3-11…17 | transaction + **rollback test** |
| D | S3-18…23 | POST handler httptest |
| E | S3-24…26 | e2e + race |

### Sprint 4 — блоки

| Блок | ID | Суть |
|------|-----|------|
| A | S4-01…02 | graceful shutdown |
| B | S4-03…05 | prometheus + slog |
| C | S4-06…07 | Docker |
| D | S4-08…16 | TESTING.md, cover, lint, DoD |

---

## 16. Backlog v2

- GitHub Actions (S4-14 stub)
- OpenAPI, fuzz decode, idempotency POST
- pprof, rate limit

---

## Старт

**Первая задача (git):** [S0-01 в TASKS.md](./TASKS.md#s0-01) · **код:** [S1-01](./TASKS.md#s1-01)

После выполнения: `ревью S1-01`

---

*Версия 3.1 — микрозадачи с каркасом в [TASKS.md](./TASKS.md)*
