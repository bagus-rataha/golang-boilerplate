# Fiber API Boilerplate v3

Ultra-clean REST API with container pattern and optimized architecture.

## What's New in v3

- **Ultra-clean main.go** - 30 lines, pure orchestration
- **Container pattern** - DI in `internal/container/`
- **App module** - Fiber setup in `internal/app/`
- **Routes module** - Organized in `internal/routes/`
- **Connection pooling** - Optimized DB performance

## Project Structure

```
fiber-api-boilerplate/
├── cmd/api/main.go           # Ultra-clean entry point
├── internal/
│   ├── app/app.go           # Fiber setup
│   ├── container/           # Dependency injection
│   ├── routes/              # Route definitions
│   ├── config/              # Configuration
│   ├── middleware/          # Middleware
│   ├── models/              # Database models
│   ├── repository/          # Data access
│   ├── services/            # Business logic
│   ├── handlers/            # HTTP handlers
│   └── utils/               # Helpers
└── migrations/              # SQL migrations
```

## Quick Start

```bash
go mod download
cp .env.example .env
createdb fiber_api
swag init -g cmd/api/main.go -o docs
air
```

## API Endpoints

### Public
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`

### Protected (JWT)
- `GET /api/v1/users/me`
- `PUT /api/v1/users/me`
- `GET /api/v1/users`

### Health
- `GET /health`

## Adding New Module

### 1. Add to Container
```go
// internal/container/container.go
type Container struct {
    AuthHandler    *handlers.AuthHandler
    UserHandler    *handlers.UserHandler
    ProductHandler *handlers.ProductHandler  // NEW
}
```

### 2. Initialize in Container
```go
productRepo := repository.NewProductRepository(db)
productService := services.NewProductService(productRepo)
productHandler := handlers.NewProductHandler(productService)
```

### 3. Add Routes
```go
// internal/routes/api.go
func setupProductRoutes(api fiber.Router, c *container.Container, cfg *config.Config) {
    products := api.Group("/products")
    products.Use(middleware.JWTProtected(cfg.JWTAccessSecret))
    products.Get("/", c.ProductHandler.List)
    products.Post("/", c.ProductHandler.Create)
}
```

### 4. Register Routes
```go
// internal/routes/api.go - SetupAPIRoutes
setupProductRoutes(api, c, cfg)
```

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_api

JWT_ACCESS_SECRET=your-access-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key
JWT_ACCESS_EXPIRE=24h
JWT_REFRESH_EXPIRE=168h

PORT=8000
APP_ENV=development

ALLOWED_ORIGINS=https://yourdomain.com
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m
```

## Connection Pool Settings

Default values (hardcoded):
- MaxIdleConns: 10
- MaxOpenConns: 100
- ConnMaxLifetime: 1 hour
- ConnMaxIdleTime: 10 minutes

## Features

- Go Fiber v2
- GORM + PostgreSQL with connection pooling
- JWT Authentication
- Request Validation
- Swagger (dev only)
- Environment-based config
- Rate Limiting (prod only)
- CORS Protection
- Clean Architecture
- Unit & integration tests (testify, golang-migrate)
- Hot Reload (Air)

## Development vs Production

| Feature | Development | Production |
|---------|-------------|------------|
| Auto Migration | Enabled | Disabled |
| Swagger UI | Enabled | Disabled |
| CORS | Allow all | Restricted |
| Rate Limiting | Disabled | Enabled |
| SQL Logging | Enabled | Disabled |
| Stack Trace | Enabled | Disabled |
| Prefork | Disabled | Enabled |
| Compression | Disabled | Enabled |

## Testing

Tests are split into two kinds:

- **Unit tests** — fast, no database. Cover `utils`, `services` (with mocked
  repositories), and `handlers` (with mocked services via `fiber.App.Test()`).
  Each layer is tested in isolation; only the boundary to external resources
  (the DB) is mocked, using `testify`.
- **Integration tests** — exercise the `repository` layer against a real
  PostgreSQL database. Gated behind the `integration` build tag so they never
  run by accident.

### Unit tests (no DB)

```bash
go test ./...
go test -cover ./...
```

### Integration tests (real DB)

Requires a dedicated test database and `TEST_DATABASE_URL`:

```bash
createdb fiber_api_test
export TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/fiber_api_test?sslmode=disable"
go test -tags=integration ./internal/repository/...
```

- `TEST_DATABASE_URL` is read by the test helper (URL form, used for both GORM
  and golang-migrate); it can also be placed in `.env`. URL-encode special
  characters in the password (e.g. `@` → `%40`). On PowerShell, set it with
  `$env:TEST_DATABASE_URL="..."`.
- The schema is applied via golang-migrate, and every table is truncated before
  each test so cases stay isolated.

### Run everything (unit + integration) with coverage

```bash
go test -tags=integration -cover ./...

# HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Docker

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/.env.example .env

EXPOSE 8000
CMD ["./api"]
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8000:8000"
    environment:
      - APP_ENV=production
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=fiber_api
      - JWT_ACCESS_SECRET=your-production-access-secret
      - JWT_REFRESH_SECRET=your-production-refresh-secret
      - ALLOWED_ORIGINS=https://yourdomain.com
      - RATE_LIMIT_MAX=100
      - RATE_LIMIT_WINDOW=1m
    depends_on:
      - db

  db:
    image: postgres:14-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=fiber_api
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

## Production Checklist

- [ ] APP_ENV=production
- [ ] Strong JWT_ACCESS_SECRET & JWT_REFRESH_SECRET (32+ chars each, distinct values)
- [ ] Configure ALLOWED_ORIGINS
- [ ] Run migrations manually
- [ ] Setup SSL/TLS
- [ ] Configure reverse proxy
- [ ] Enable monitoring

## License

MIT
