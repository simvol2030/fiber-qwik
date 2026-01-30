# Backend Go Fiber Template

Production-ready Go Fiber v2 backend template with JWT authentication, GORM, and security best practices.

## Features

- Go Fiber v2 (Express-inspired, fastest Go web framework)
- JWT Authentication (access + refresh tokens)
- GORM with SQLite (easy switch to PostgreSQL)
- Security middleware (Helmet, CORS, Rate Limiting)
- Structured logging with zerolog
- Health check endpoints
- Graceful shutdown
- Connection pooling
- Docker-ready

## Project Structure

```
cmd/
└── server/
    └── main.go              # Entry point

internal/
├── handlers/
│   ├── auth.go              # Auth endpoints
│   └── health.go            # Health check endpoints
├── middleware/
│   ├── auth.go              # JWT authentication
│   └── security.go          # Helmet, CORS, Rate limiting
├── models/
│   └── user.go              # Database models
├── services/
│   └── auth.go              # Auth business logic
└── utils/
    ├── jwt.go               # JWT utilities
    ├── password.go          # Password hashing
    └── response.go          # API response helpers
```

## Quick Start

### Development

```bash
# Install dependencies
go mod download

# Create database directory
mkdir -p data/db/sqlite

# Run development server
go run ./cmd/server

# Or with hot reload (using air)
air
```

Server runs at http://localhost:3001

### Production Build

```bash
# Build binary
CGO_ENABLED=1 go build -o server ./cmd/server

# Run production server
./server
```

### Docker

```bash
# Build image
docker build -t backend-go-fiber .

# Run container
docker run -p 3001:3001 \
  -e JWT_SECRET=your-secret-key \
  -v ./data:/app/data \
  backend-go-fiber
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login user |
| POST | `/api/auth/refresh` | Refresh access token |
| POST | `/api/auth/logout` | Logout user |
| GET | `/api/auth/me` | Get current user (protected) |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Basic health check |
| GET | `/ready` | Readiness check (DB connection) |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3001` | Server port |
| `HOST` | `0.0.0.0` | Server host |
| `NODE_ENV` | `development` | Environment |
| `DATABASE_URL` | `file:./data/db/sqlite/app.db...` | Database connection |
| `JWT_SECRET` | `dev-secret...` | JWT signing secret |
| `JWT_EXPIRES_IN` | `15m` | Access token expiry |
| `REFRESH_TOKEN_EXPIRES_DAYS` | `7` | Refresh token expiry |
| `CORS_ORIGINS` | `http://localhost:3000` | Allowed CORS origins |

## Performance

Go Fiber is built for high performance:
- ~600k requests/sec
- Minimal memory footprint
- Goroutines for concurrent requests
- Connection pooling configured

### Configuration in main.go:
```go
app := fiber.New(fiber.Config{
    Prefork:       false,        // Multi-process (set true for production)
    ReadTimeout:   10 * time.Second,
    WriteTimeout:  10 * time.Second,
    IdleTimeout:   30 * time.Second,
    BodyLimit:     10 * 1024 * 1024,
})
```

## Switch to PostgreSQL

1. Install PostgreSQL driver:
```bash
go get gorm.io/driver/postgres
```

2. Update connection in main.go:
```go
import "gorm.io/driver/postgres"

dsn := os.Getenv("DATABASE_URL")
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

3. Update DATABASE_URL:
```env
DATABASE_URL=postgresql://user:password@localhost:5432/app
```
