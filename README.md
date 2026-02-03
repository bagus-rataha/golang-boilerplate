# Fiber API Boilerplate

Simple REST API boilerplate with Go Fiber, GORM, JWT, and PostgreSQL.

## Features

- Go Fiber v2 - Fast HTTP framework
- GORM v2 - ORM with PostgreSQL
- JWT Authentication - Access & Refresh tokens
- Swagger Documentation - Auto-generated API docs
- Clean Architecture - Organized code structure
- Hot Reload - Air for development
- Environment Config - Configurable durations

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Air (for hot reload)

## Installation

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
# Edit .env with your database credentials
```

4. Create database
```bash
createdb fiber_api
```

5. Install Air (optional, for hot reload)
```bash
go install github.com/cosmtrek/air@latest
```

6. Install Swag (for swagger generation)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Running

### Development (with hot reload)
```bash
air
```

### Production
```bash
go run cmd/api/main.go
```

## Generate Swagger Docs
```bash
swag init -g cmd/api/main.go -o docs
```

## API Documentation

Access Swagger UI at: `http://localhost:8000/swagger/`

## Endpoints

### Auth
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token

### Users (Protected)
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `GET /api/v1/users` - List all users

## Environment Variables
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_api

# JWT (customize durations)
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRE=24h      # Access token duration
JWT_REFRESH_EXPIRE=168h    # Refresh token duration (7 days)

# Server
PORT=8000
APP_ENV=development
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
├── docs/                # Swagger documentation
├── .env.example         # Environment template
└── README.md
```

## License

MIT
