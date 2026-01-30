# Audit Checklist: project-box-go-fiber-qwik-city

> Audit date: 2026-01-30
> Requirement sources: VERSIONS.md, CLAUDE.md, SECURITY.md
> **Status: PENDING**

---

## 1. PROJECT STRUCTURE

### 1.1 Directories and files
- [ ] `frontend-qwik-city/` exists
- [ ] `backend-go-fiber/` exists
- [ ] `data/db/sqlite/.gitkeep` exists
- [ ] `data/db/postgres/.gitkeep` exists
- [ ] `data/logs/.gitkeep` exists
- [ ] `data/media/.gitkeep` exists
- [ ] `docker-compose.yml` exists
- [ ] `.env.example` exists
- [ ] `.gitignore` exists
- [ ] `README.md` exists

### 1.2 Workflow files (per CLAUDE.md)
- [ ] `docs/README.md` exists
- [ ] `project-doc/COMPLETED.md` exists
- [ ] `project-doc/archive/` exists
- [ ] `feedbacks/.gitkeep` exists
- [ ] `.claude/settings.json` exists
- [ ] `.claude/hooks/notify.sh` exists
- [ ] `CLAUDE.md` (codebase context) exists
- [ ] `CLAUDE.web.md` (developer workflow) exists

---

## 2. VERSION DEPENDENCIES (VERSIONS.md)

### 2.1 Backend Go Fiber
- [ ] Go version >= 1.21
- [ ] `github.com/gofiber/fiber/v2`
- [ ] `gorm.io/gorm` v1.x
- [ ] `gorm.io/driver/sqlite`
- [ ] `github.com/golang-jwt/jwt/v5`
- [ ] `golang.org/x/crypto/bcrypt`
- [ ] `github.com/rs/zerolog` (structured logging)

### 2.2 Frontend Qwik City
- [ ] Node.js >= 20.x (in Dockerfile)
- [ ] `@builder.io/qwik` ~1.14.x
- [ ] `@builder.io/qwik-city` ~1.14.x
- [ ] `typescript` ^5.x
- [ ] `vite` ^5.x or ^6.x

---

## 3. BACKEND API (CLAUDE.md)

### 3.1 Authentication Endpoints
- [ ] `POST /api/auth/register` -- registration
- [ ] `POST /api/auth/login` -- login
- [ ] `POST /api/auth/refresh` -- token refresh
- [ ] `POST /api/auth/logout` -- logout
- [ ] `GET /api/auth/me` -- current user (protected)

### 3.2 Health Endpoints
- [ ] `GET /health` -- returns `{"status":"ok"}`
- [ ] `GET /ready` -- checks DB connection

### 3.3 API Response Format
- [ ] Success: `{ success: true, data: {...}, meta: {timestamp} }`
- [ ] Error: `{ success: false, error: {code, message}, meta: {...} }`

---

## 4. SECURITY (SECURITY.md)

### 4.1 HTTP Headers (Helmet)
- [ ] Helmet middleware enabled
- [ ] X-Content-Type-Options: nosniff
- [ ] X-Frame-Options: DENY
- [ ] Strict-Transport-Security set

### 4.2 CORS
- [ ] CORS middleware enabled
- [ ] Whitelist origins (not `*`)
- [ ] AllowCredentials: true
- [ ] Allowed methods: GET, POST, PUT, DELETE

### 4.3 Rate Limiting
- [ ] Rate limiter enabled
- [ ] 100 req/min for general endpoints

### 4.4 Input Validation
- [ ] Email validation
- [ ] Password minimum 8 characters
- [ ] No raw SQL queries

### 4.5 JWT Authentication
- [ ] Access token expiry: 15 minutes
- [ ] Refresh token expiry: 7 days
- [ ] Refresh token in httpOnly cookie
- [ ] Password hashing: bcrypt (10 rounds)

### 4.6 Error Handling
- [ ] Does not show stack traces in production
- [ ] Logs errors on server
- [ ] Returns safe messages to client

---

## 5. DATABASE

### 5.1 Schema
- [ ] Table `users` exists
- [ ] Fields: id, email (unique), password_hash, name, created_at, updated_at
- [ ] Table `refresh_tokens` exists
- [ ] Fields: id, token (unique), user_id (FK), expires_at, created_at

### 5.2 Operations
- [ ] User creation works
- [ ] Email unique constraint works
- [ ] Refresh token saved in DB
- [ ] Refresh token deleted on logout

---

## 6. FRONTEND QWIK CITY

### 6.1 Pages
- [ ] `/` -- home page
- [ ] `/login` -- login page
- [ ] `/register` -- registration page
- [ ] `/dashboard` -- protected page

### 6.2 Functionality
- [ ] API client for backend requests
- [ ] Access token storage (memory/store)
- [ ] Handle 401 and redirect to /login
- [ ] Logout clears tokens

### 6.3 Qwik City
- [ ] Uses `useSignal` / `useStore` for state
- [ ] Uses `component$` for components
- [ ] Uses `routeLoader$` / `routeAction$` for data loading

---

## 7. DOCKER

### 7.1 docker-compose.yml
- [ ] Service `frontend` on port 3000
- [ ] Service `backend` on port 3001
- [ ] Volumes for `./data:/app/data`
- [ ] Environment variables configured
- [ ] Healthcheck for backend
- [ ] Profile `postgres` for PostgreSQL

### 7.2 Dockerfiles
- [ ] `backend-go-fiber/Dockerfile` exists and builds
- [ ] `frontend-qwik-city/Dockerfile` exists and builds
- [ ] Multi-stage build (builder + production)
- [ ] Non-root user in production

### 7.3 Launch
- [ ] `docker-compose up --build` starts
- [ ] Frontend accessible at http://localhost:3000
- [ ] Backend accessible at http://localhost:3001
- [ ] Health check passes

---

## 8. BROWSER TESTING

### 8.1 Basic Navigation
- [ ] Home page loads
- [ ] /login page loads
- [ ] /register page loads

### 8.2 Registration
- [ ] Registration form displays
- [ ] Field validation works
- [ ] Successful registration creates user
- [ ] After registration redirect to dashboard

### 8.3 Login
- [ ] Login form displays
- [ ] Wrong credentials show error
- [ ] Successful login saves token
- [ ] After login redirect to dashboard

### 8.4 Dashboard (Protected)
- [ ] Without token -- redirect to /login
- [ ] With token -- shows user data
- [ ] Logout works and clears session

### 8.5 Browser Console
- [ ] No JavaScript errors (except expected 401)
- [ ] No CORS errors
- [ ] No 500 errors from API

---

## 9. ENVIRONMENT VARIABLES

### 9.1 .env.example contains
- [ ] PORT
- [ ] DATABASE_URL
- [ ] JWT_SECRET
- [ ] JWT_EXPIRES_IN
- [ ] REFRESH_TOKEN_EXPIRES_DAYS
- [ ] CORS_ORIGINS

### 9.2 Validation
- [ ] DATABASE_URL has default for SQLite

---

## 10. LOGGING

### 10.1 Backend
- [ ] Structured JSON logging (zerolog)
- [ ] Levels: info, warn, error
- [ ] Timestamp in logs
- [ ] Request ID in logs

---

## AUDIT RESULTS

### Critical Issues (blockers)
*Pending audit*

### Fixed During Audit
*Pending audit*

### Notes
*Pending audit*

---

*Audit pending*
*Version: 1.0*
*Status: PENDING*
