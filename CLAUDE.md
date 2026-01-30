# project-box-go-fiber-qwik-city

> Production-ready starter: Go Fiber + Qwik City

---

## Technologies

| Component | Technology | Version |
|-----------|------------|---------|
| Backend | Go Fiber | v2 |
| ORM | GORM | v1.25 |
| Database | SQLite / PostgreSQL | - |
| Frontend | Qwik City | ~1.14.x |
| UI | Qwik signals + component$ | - |
| Auth | JWT + Refresh Tokens | - |

---

## Structure

```
project-box-go-fiber-qwik-city/
├── backend-go-fiber/          # Go Fiber API
│   ├── cmd/server/            # Entry point
│   ├── internal/
│   │   ├── handlers/          # HTTP handlers
│   │   ├── middleware/        # Auth, CORS, Rate limit
│   │   ├── models/            # GORM models
│   │   ├── services/          # Business logic
│   │   └── utils/             # Helpers
│   ├── go.mod
│   └── Dockerfile
│
├── frontend-qwik-city/        # Qwik City app
│   ├── src/
│   │   ├── routes/            # Pages
│   │   └── lib/               # Components, stores
│   ├── package.json
│   └── Dockerfile
│
├── data/                      # Persistent data (gitignored)
│   ├── db/sqlite/
│   └── logs/
│
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Commands

```bash
# Development (without Docker)
cd backend-go-fiber && go run cmd/server/main.go
cd frontend-qwik-city && npm run dev

# Docker
docker-compose up --build

# With PostgreSQL
docker-compose --profile postgres up --build
```

---

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /health | - | Health check |
| GET | /ready | - | Readiness (DB check) |
| POST | /api/auth/register | - | Registration |
| POST | /api/auth/login | - | Login |
| POST | /api/auth/refresh | Cookie | Refresh token |
| POST | /api/auth/logout | Cookie | Logout |
| GET | /api/auth/me | Bearer | Current user |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 3001 | Backend port |
| DATABASE_URL | sqlite:./data/db/sqlite/app.db | Database |
| JWT_SECRET | dev-secret... | JWT secret (CHANGE!) |
| JWT_EXPIRES_IN | 15m | Access token TTL |
| REFRESH_TOKEN_EXPIRES_DAYS | 7 | Refresh token TTL (days) |
| CORS_ORIGINS | http://localhost:3000 | Allowed origins |

---

*Version: 1.0*
*Created: 2025-01-30*
