# Project Box: Go Fiber + Qwik City

Production-ready starter template with **Go Fiber v2** + **GORM** backend and **Qwik City** frontend.

---

## New Project from This Template

Clone the repo and start building your project in minutes.

### Step 1: Clone and rename

```bash
git clone https://github.com/simvol2030/fiber-qwik.git my-project
# or use GitHub template: "Use this template" button on GitHub
cd my-project

# Point to your own repository
git remote set-url origin https://github.com/YOUR_USER/YOUR_REPO.git
```

### Step 2: Change these values

**Quick find & replace:** Search for `YOUR_DOMAIN` and `YOUR_PROJECT` across all config files.

| What | Where | Change to |
|------|-------|-----------|
| **JWT_SECRET** | `ecosystem.config.js` or `.env` | New random string 32+ chars. Generate: `openssl rand -base64 32` |
| **Domain** | `ecosystem.config.js` -> `CORS_ORIGINS`, `FRONTEND_URL`, `ORIGIN` | Your domain (e.g. `https://my-app.com`) |
| **Domain** | `deploy/nginx.conf` -> `server_name`, `ssl_certificate`, logs | Replace all `YOUR_DOMAIN` |
| **Project path** | `ecosystem.config.js` -> `cwd`, `log_file`, `error_file` | Replace all `YOUR_PROJECT` (e.g. `my-app`) |
| **PM2 names** | `ecosystem.config.js` -> `name` | Replace `YOUR_PROJECT-backend`, `YOUR_PROJECT-frontend` |
| **App name** | `frontend-qwik-city/package.json` -> `name` | Your project name |
| **App title** | `frontend-qwik-city/src/root.tsx` -> `<title>` | Your app title |
| **Footer** | `frontend-qwik-city/src/routes/layout.tsx` | Your company/project name |
| **Admin email** | Database seed or first register | Your admin email |

### Step 3: Reset database (clean start)

```bash
# Delete existing database
rm -f data/db/sqlite/app.db

# Start backend -- GORM auto-creates tables on first run
# Register your admin user at /register, then set role in DB:
sqlite3 data/db/sqlite/app.db "UPDATE users SET role='admin' WHERE email='your@email.com';"
```

### Step 4: Choose deployment method

**Option A: PM2 (lightweight, recommended for VPS)**
```bash
# Copy PM2 config template and fill in your values
cp ecosystem.config.js.example ecosystem.config.js
# Edit ecosystem.config.js -- change JWT_SECRET, domain, paths

# Build
cd backend-go-fiber && go build -o server cmd/server/main.go
cd ../frontend-qwik-city && npm install && npm run build

# Start
pm2 start ecosystem.config.js
pm2 save
```

**Option B: Docker**
```bash
cp .env.example .env
# Edit .env -- change JWT_SECRET, domain
docker-compose up --build -d
```

### Step 5: Nginx + SSL

```bash
# Create nginx config (replace YOUR_DOMAIN)
sudo nano /etc/nginx/sites-available/YOUR_DOMAIN

# Template:
# server {
#     server_name YOUR_DOMAIN;
#     location / { proxy_pass http://127.0.0.1:3000; ... }
#     location /api/ { proxy_pass http://127.0.0.1:3001; ... }
#     location /health { proxy_pass http://127.0.0.1:3001; }
#     location /uploads/ { proxy_pass http://127.0.0.1:3001; }
#     client_max_body_size 10M;
# }

sudo ln -s /etc/nginx/sites-available/YOUR_DOMAIN /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx

# SSL
sudo certbot --nginx -d YOUR_DOMAIN
```

### What you get out of the box

- **Auth**: Register, login, logout, refresh tokens, password reset
- **Admin panel**: `/admin` -- users CRUD, file browser with upload, app settings
- **User dashboard**: `/dashboard` -- profile, name edit, password change
- **Rate limiting**: Login (5/5min), register (3/hour), global (100/min)
- **Security**: Helmet, CORS, JWT httpOnly cookies, input validation
- **File upload**: Drag-and-drop, validation (type + size), S3-ready
- **Logging**: Structured JSON (zerolog)
- **Health checks**: `/health`, `/ready`

### What to add per project

```
Your models     ->  backend-go-fiber/internal/models/
Your handlers   ->  backend-go-fiber/internal/handlers/
Your services   ->  backend-go-fiber/internal/services/
Your pages      ->  frontend-qwik-city/src/routes/
Your components ->  frontend-qwik-city/src/lib/components/
```

---

## Features

### Core
- **Authentication**: JWT access tokens + refresh tokens (httpOnly cookies)
- **Password Reset**: Forgot password flow with email tokens
- **User Roles**: Basic RBAC with user/admin roles
- **File Upload**: Storage interface (Local + S3/MinIO support)
- **Email Service**: SMTP sender with HTML templates (Mock in dev)
- **Soft Delete**: Built-in for User model (GORM DeletedAt)

### Performance
- **Gzip Compression**: Backend middleware for all responses
- **SQLite WAL Mode**: ~256MB cache, concurrent reads
- **Connection Pooling**: 100 connections for PostgreSQL

### Security
- **Helmet**: XSS, CSP, Frame options, Referrer policy
- **CORS**: Configurable allowed origins (no `*` in production)
- **Rate Limiting**: 100 req/min per IP
- **Input Validation**: go-playground/validator
- **CSP**: Relaxed in production for CDN, fonts, analytics

### Infrastructure
- **Database**: SQLite (dev) / PostgreSQL (production)
- **Docker**: Multi-stage builds, non-root users, health checks
- **Testing**: Unit tests for auth service
- **Linting**: ESLint + Prettier (frontend), go vet (backend)
- **Seeder**: Test data with `make seed`

---

## Quick Start

### 1. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Generate a secure JWT secret
make generate-secret
# Copy the output to JWT_SECRET in .env
```

### 2. Start Development

**Option A: Docker (recommended)**
```bash
make docker
# or manually:
docker-compose up --build
```

**Option B: Local Development**
```bash
# Install dependencies
make install

# Terminal 1: Start backend
make dev-backend

# Terminal 2: Start frontend
make dev-frontend
```

### 3. Access the App

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:3001 |
| Health Check | http://localhost:3001/health |

---

## Makefile Commands

```bash
make help              # Show all commands

# Development
make install           # Install all dependencies
make dev-backend       # Start Go backend
make dev-frontend      # Start Qwik City frontend

# Build
make build             # Build for production
make build-backend     # Build Go binary
make build-frontend    # Build Qwik City

# Testing
make test              # Run all tests
make test-backend      # Run Go tests
make test-backend-coverage  # Run with coverage report

# Docker
make docker            # Build and start containers
make docker-down       # Stop containers
make docker-logs       # View logs
make docker-postgres   # Start with PostgreSQL

# Database
make seed              # Seed with test data
make db-reset          # Reset SQLite database
make db-fresh          # Reset + seed

# Utilities
make lint              # Lint all code
make format            # Format all code
make clean             # Clean build artifacts
make generate-secret   # Generate JWT secret
```

---

## API Endpoints

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | - | Register new user |
| POST | `/api/auth/login` | - | Login, get tokens |
| POST | `/api/auth/refresh` | Cookie | Refresh access token |
| POST | `/api/auth/logout` | Cookie | Logout, clear tokens |
| GET | `/api/auth/me` | Bearer | Get current user |
| POST | `/api/auth/forgot-password` | - | Request password reset |
| POST | `/api/auth/validate-reset-token` | - | Validate reset token |
| POST | `/api/auth/reset-password` | - | Reset password with token |

### File Upload

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/upload` | Bearer | Upload single file |
| POST | `/api/upload/multiple` | Bearer | Upload multiple files (max 10) |
| DELETE | `/api/upload/*` | Bearer | Delete file by key |
| GET | `/uploads/*` | - | Serve uploaded files (local only) |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Basic health check |
| GET | `/ready` | Readiness (DB check) |

### Request/Response Format

```json
// Success Response
{
  "success": true,
  "data": { ... },
  "meta": {
    "timestamp": "2025-01-18T12:00:00Z",
    "requestId": "uuid"
  }
}

// Error Response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": [{"field": "email", "message": "..."}]
  },
  "meta": { ... }
}
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `NODE_ENV` | development | Environment mode |
| `PORT` | 3001 | Backend port |
| `DATABASE_URL` | sqlite:./data/... | Database connection |
| `JWT_SECRET` | - | **Required in production** (32+ chars) |
| `JWT_EXPIRES_IN` | 15m | Access token TTL |
| `REFRESH_TOKEN_EXPIRES_DAYS` | 7 | Refresh token TTL |
| `CORS_ORIGINS` | http://localhost:3000 | Allowed origins |
| `LOG_LEVEL` | info | debug/info/warn/error |
| `SMTP_HOST` | - | SMTP server for emails (production) |
| `S3_BUCKET` | - | S3 bucket for file storage |
| `FRONTEND_URL` | http://localhost:3000 | For password reset links |

See `.env.example` for complete list.

---

## Production Deployment

### Pre-deployment Checklist

```bash
make deploy-check  # Show checklist
```

- [ ] `NODE_ENV=production` in `.env`
- [ ] `JWT_SECRET` is 32+ random characters
- [ ] `CORS_ORIGINS` set to production domain(s)
- [ ] `DATABASE_URL` points to production database
- [ ] SSL/TLS configured (nginx)
- [ ] Backups configured

### Deploy with PM2

```bash
# 1. Copy PM2 config
cp ecosystem.config.js.example ecosystem.config.js
# Edit ecosystem.config.js -- change JWT_SECRET, domain, paths

# 2. Build backend
cd backend-go-fiber && go build -o server cmd/server/main.go

# 3. Build frontend
cd frontend-qwik-city && npm install && npm run build

# 4. Start services
pm2 start ecosystem.config.js
pm2 save
```

### Deploy with Docker

```bash
# 1. Build production images
docker-compose build

# 2. Start containers
docker-compose up -d

# 3. Check health
curl http://localhost:3001/health
```

### Deploy with Nginx (Reverse Proxy)

```bash
# 1. Copy nginx config (replace YOUR_DOMAIN with your domain)
sudo cp deploy/nginx.conf /etc/nginx/sites-available/YOUR_DOMAIN.conf

# 2. Replace all YOUR_DOMAIN placeholders
sudo sed -i 's/YOUR_DOMAIN/my-app.example.com/g' /etc/nginx/sites-available/YOUR_DOMAIN.conf

# 3. Enable site
sudo ln -s /etc/nginx/sites-available/YOUR_DOMAIN.conf /etc/nginx/sites-enabled/

# 4. Test and reload
sudo nginx -t && sudo systemctl reload nginx
```

### SSL with Let's Encrypt

```bash
# Install certbot
apt install certbot python3-certbot-nginx

# Get certificate (replace with your domain)
certbot --nginx -d YOUR_DOMAIN

# Auto-renewal (cron)
0 0 * * * /usr/bin/certbot renew --quiet
```

---

## Database: SQLite → PostgreSQL

The project starts with **SQLite** (zero setup) and switches to **PostgreSQL** when you need production scale. GORM auto-migrates tables on startup — no manual migrations needed.

### When to switch

| SQLite | PostgreSQL |
|--------|-----------|
| Development, prototyping | Production with traffic |
| Single server, low load | Multiple connections, high load |
| 1 process (Prefork=false) | Multi-process (Prefork=true) |

### How to switch (1 line change)

```bash
# .env — development (SQLite, default)
DATABASE_URL=sqlite:./data/db/sqlite/app.db

# .env — production (PostgreSQL, just uncomment and fill)
DATABASE_URL=postgresql://user:password@localhost:5432/app?sslmode=disable
```

That's it. Restart the backend — GORM creates all tables automatically.

### PostgreSQL with Docker (local testing)

```bash
docker-compose --profile postgres up -d

# Backend will use:
DATABASE_URL=postgresql://app:secret@localhost:5432/app
```

### PostgreSQL in production

```bash
# 1. Create database
createdb -U postgres app

# 2. Update ecosystem.config.js
DATABASE_URL=postgresql://user:password@host:5432/app?sslmode=require

# 3. Enable multi-process mode (optional, for high load)
PREFORK=true

# 4. Restart
pm2 restart all
```

---

## Project Structure

```
project-box-go-fiber-qwik-city/
├── backend-go-fiber/
│   ├── cmd/server/main.go      # Entry point
│   ├── internal/
│   │   ├── handlers/           # HTTP handlers
│   │   ├── middleware/         # Auth, CORS, Security
│   │   ├── models/             # GORM models
│   │   ├── services/           # Business logic
│   │   └── utils/              # JWT, validation, helpers
│   ├── go.mod
│   └── Dockerfile
│
├── frontend-qwik-city/
│   ├── src/
│   │   ├── routes/             # Pages
│   │   └── lib/
│   │       ├── api/client.ts   # API client
│   │       └── stores/         # Auth store
│   ├── package.json
│   └── Dockerfile
│
├── deploy/
│   └── nginx.conf              # Nginx reverse proxy config
│
├── data/                       # Persistent data (gitignored)
│   ├── db/sqlite/
│   ├── db/postgres/
│   └── logs/
│
├── docker-compose.yml
├── Makefile
├── .env.example
└── README.md
```

---

## Troubleshooting

### Docker: "JWT_SECRET is required"

```bash
# Create .env file
cp .env.example .env
# Edit .env and set JWT_SECRET
```

### Backend: "Failed to connect to database"

```bash
# Ensure data directory exists
mkdir -p data/db/sqlite

# Check permissions
chmod 755 data/db/sqlite
```

### Frontend: Build fails

```bash
# Clear cache and reinstall
rm -rf frontend-qwik-city/node_modules frontend-qwik-city/.qwik frontend-qwik-city/dist
cd frontend-qwik-city && npm install
```

---

## License

MIT

---

*Built with [Box-App](https://github.com/your-repo/box-app) templates v2.0*
