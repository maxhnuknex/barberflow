<div align="center">

# BarberFlow

**Telegram-сервис онлайн-записи для барбершопов, который позволяет клиентам самостоятельно бронировать визиты, а администратору управлять записями, и отправляет напоминания.**

[![CI](https://github.com/maxhnuknex/barberflow/actions/workflows/ci.yml/badge.svg)](https://github.com/maxhnuknex/barberflow/actions/workflows/ci.yml)

Go · PostgreSQL · pgx · Telegram Bot API · Docker · golang-migrate · GitHub Actions

</div>

## О проекте

BarberFlow — Telegram-бот для записи клиентов в барбершоп. Пользователь проходит весь сценарий бронирования внутри Telegram: выбирает услугу, мастера, дату и свободное время, подтверждает запись, просматривает свои будущие визиты и при необходимости отменяет их.

Для администратора предусмотрен отдельный Telegram-интерфейс для просмотра записей по датам, работы с отдельной записью, услугами и мастерами. Фоновый worker периодически проверяет ближайшие визиты и отправляет клиентам напоминания.

## Возможности

**Клиент**
- пошаговая запись: услуга → мастер → дата → свободный слот → подтверждение;
- расчёт доступных слотов с учётом графика мастера и уже существующих записей;
- просмотр своих будущих записей;
- просмотр деталей и отмена записи;
- автоматическое Telegram-напоминание перед визитом.

**Администратор**
- просмотр записей на сегодня и выбранную дату;
- просмотр деталей записи и её отмена;
- просмотр услуг и мастеров;
- включение/отключение доступности услуг и мастеров через административное меню.


## Стек

- **Go** — приложение и бизнес-логика;
- **PostgreSQL** — хранение клиентов, услуг, мастеров, графика и записей;
- **pgx** — PostgreSQL driver/pool;
- **go-telegram/bot** — интеграция с Telegram Bot API;
- **golang-migrate** — версионирование схемы БД;
- **Docker / Docker Compose** — воспроизводимый локальный запуск;
- **GitHub Actions** — автоматические проверки после push и pull request.

## Архитектура

```mermaid
flowchart LR
    TG[Telegram User / Admin] --> H[Telegram Handlers]
    H --> S[Application Services]
    S --> R[PostgreSQL Repositories]
    R --> DB[(PostgreSQL)]

    W[Reminder Worker] --> BS[Booking Service]
    BS --> R
    W --> API[Telegram Bot API]
```

Проект разделён по ответственности:

- `internal/delivery/telegram` обрабатывает Telegram updates, callback'и и формирует UI;
- `internal/app` содержит application/business logic для booking, catalog и barber;
- `internal/repository/postgres` изолирует SQL и работу с PostgreSQL;
- `internal/domain` содержит основные модели предметной области;
- `internal/worker/reminder` выполняет периодическую фоновую обработку напоминаний;
- `cmd/barberflow` собирает зависимости и управляет жизненным циклом приложения.

## Быстрый запуск

### Требования

- Docker
- Docker Compose
- Telegram Bot Token

Клонировать репозиторий:

```bash
git clone https://github.com/maxhnuknex/barberflow.git
cd barberflow
```

Создать `.env` в корне проекта:

```env
TG_TOKEN=your_telegram_bot_token

POSTGRES_DB=barberflow
POSTGRES_USER=barberflow
POSTGRES_PASSWORD=change_me

DATABASE_URL=postgres://barberflow:change_me@postgres:5432/barberflow?sslmode=disable
```

Запустить весь стек:

```bash
docker compose up --build
```

Docker Compose последовательно запускает PostgreSQL, применяет migrations и после их успешного завершения запускает BarberFlow.


## Тестирование

Локальная проверка проекта:

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
go build ./cmd/barberflow
```
## Структура проекта

```text
cmd/
└── barberflow/
    └── main.go

internal/
├── app/
│   ├── barber/
│   ├── booking/
│   └── catalog/
├── delivery/
│   └── telegram/
│       ├── admin/
│       ├── booking/
│       ├── mybookings/
│       ├── start/
│       └── ui/
├── domain/
├── repository/
│   └── postgres/
└── worker/
    └── reminder/

migrations/
.github/
└── workflows/
    └── ci.yml

compose.yaml
Dockerfile
```

## Инженерные детали

- доступные слоты рассчитываются с учётом длительности услуги, рабочего интервала мастера и пересечений с существующими booking;
- данные сохраняются в PostgreSQL, а изменения схемы оформляются последовательными migrations;
- reminder помечается как отправленный только после успешной отправки сообщения, что предотвращает повторную обработку при следующем проходе worker;
- общий `context` используется для остановки Telegram bot и reminder worker при завершении приложения;
- CI отдельно запускает `go vet`, unit tests, race detector и сборку основного binary.
