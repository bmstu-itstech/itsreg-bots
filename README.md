# ITS Reg: сервис для создания регистраций на мероприятия

## Как запускать

### Локальный запуск

Для локального запуска сервера и для проведения интеграционных тестов используется _локальное окружение_.

Поднимать при помощи `docker compose`:

```shell
docker compose -f deployment/docker-compose.local.yaml up -d
```

Запуск тестов:
- `go test -short ./internal/domain/...` - юнит-тесты домена;
- `go test -count=1 ./internal/infra/...` - интеграционные (инфраструктурные) тесты; также требуется задать переменную
  окружения `DATABASE_URL` (для local - `postgres://test-user:test-pass@localhost:30001/test-db?sslmode=disable`);
- `go test -count=1 ./internal/ports/...` - компонентные (функциональные) тесты; также требуется задать переменную 
  окружения `PORT`.

Запуск серверов:
- `go run ./cmd/http/http.go` - HTTP API сервиса, требуется задать переменные окружения `PORT` и `DATABASE_URL`;
- `go run ./cmd/telegram/telegram.go` - сервер для взаимодействия с telegram API, требуется задать переменную окружения `DATABASE_URL`.

### Дев

Для запуска полной, рабочей копии, проекта, но для локального запуска используется _dev окружение_.

Поднимать при помощи `docker compose`:

```shell
docker compose -f deployment/docker-compose.dev.yaml up -d
```

Загружает окружение из `.env` файла. Пример содержимого `.env` представлен в файле 

### Продакшен

Поднимать при помощи `docker compose`:

```shell
docker compose -f deployment/docker-compose.prod.yaml up -d
```

Загружает окружение из `.env` файла. Пример содержимого `.env` представлен в файле 

