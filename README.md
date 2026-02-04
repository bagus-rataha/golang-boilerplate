# Fiber API Boilerplate v2

Production-ready REST API boilerplate with Go Fiber, GORM, JWT, and PostgreSQL.

## Features

### Core
- **Go Fiber v2** - Fast HTTP framework
- **GORM v2** - ORM with PostgreSQL
- **JWT Authentication** - Access & Refresh tokens
- **Swagger Documentation** - Auto-generated API docs
- **Clean Architecture** - Organized code structure
- **Hot Reload** - Air for development

### v2 Features
- **Request Validation** - go-playground/validator with user-friendly messages
- **Environment Config** - Development vs Production modes
- **Rate Limiting** - Production-only rate limiting
- **CORS Protection** - Environment-specific CORS
- **Conditional Swagger** - Development-only API docs
- **Migration Strategy** - golang-migrate for production

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

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Air (optional, for hot reload)

### Installation

1. Clone repository
```bash
git clone <repository-url>
cd fiber-api-boilerplate
```

2. Install dependencies
```bash
go mod download
```

3. Setup environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Create database
```bash
createdb fiber_api
```

5. Install tools (optional)
```bash
# Hot reload
go install github.com/cosmtrek/air@latest

# Swagger generation
go install github.com/swaggo/swag/cmd/swag@latest
```

### Running

#### Development
```bash
# With hot reload
air

# Without hot reload
go run cmd/api/main.go
```

#### Production
```bash
# Build
go build -o api cmd/api/main.go

# Run with production environment
APP_ENV=production ./api
```

## Generate Swagger Docs

```bash
swag init -g cmd/api/main.go -o docs
```

Access Swagger UI at: `http://localhost:8000/swagger/` (development only)

## API Endpoints

### Auth (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/register | Register new user |
| POST | /api/v1/auth/login | Login user |
| POST | /api/v1/auth/refresh | Refresh access token |

### Users (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/users/me | Get current user profile |
| PUT | /api/v1/users/me | Update current user profile |
| GET | /api/v1/users | List all users |

## Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_api

# JWT
JWT_SECRET=your-super-secret-key
JWT_ACCESS_EXPIRE=24h
JWT_REFRESH_EXPIRE=168h

# Server
PORT=8000
APP_ENV=development  # or production

# CORS (Production)
ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com

# Rate Limit (Production)
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m
```

## Project Structure

```
fiber-api-boilerplate/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── models/          # Database models
│   ├── repository/      # Database operations
│   ├── services/        # Business logic
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   └── utils/           # Helper functions
├── migrations/          # Database migrations
├── docs/                # Swagger documentation
├── .env.example         # Environment template
└── README.md
```

## Docker Deployment

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
      - JWT_SECRET=your-production-secret
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

- [ ] Set `APP_ENV=production`
- [ ] Use strong `JWT_SECRET`
- [ ] Configure `ALLOWED_ORIGINS`
- [ ] Set appropriate rate limits
- [ ] Run database migrations with golang-migrate
- [ ] Enable HTTPS (via reverse proxy)
- [ ] Setup logging/monitoring
- [ ] Configure database connection pooling

## Troubleshooting

### Swagger not loading
- Ensure you're in development mode (`APP_ENV=development`)
- Regenerate docs: `swag init -g cmd/api/main.go -o docs`

### Database connection failed
- Check PostgreSQL is running
- Verify credentials in `.env`
- Ensure database exists: `createdb fiber_api`

### Rate limit errors
- Rate limiting only active in production
- Adjust `RATE_LIMIT_MAX` and `RATE_LIMIT_WINDOW`

### CORS errors
- In development, all origins allowed
- In production, configure `ALLOWED_ORIGINS`

## License

MIT
