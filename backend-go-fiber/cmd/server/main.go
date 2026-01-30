package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"backend-go-fiber/internal/handlers"
	adminHandlers "backend-go-fiber/internal/handlers/admin"
	"backend-go-fiber/internal/middleware"
	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/services"
	adminServices "backend-go-fiber/internal/services/admin"
	"backend-go-fiber/internal/services/email"
	"backend-go-fiber/internal/services/storage"
	"backend-go-fiber/internal/services/upload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Set timezone to UTC
	time.Local = time.UTC

	// Configure logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set log level from environment (LOG_LEVEL: debug, info, warn, error)
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Use console writer in development for prettier output
	if os.Getenv("NODE_ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Get config from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// SQLite with WAL mode for better concurrency:
		// - _journal_mode=WAL: Write-Ahead Logging (readers don't block writers)
		// - _busy_timeout=5000: Wait 5s before "database is locked" error
		// - _synchronous=NORMAL: Good balance of safety and performance
		// - _cache_size=-262144: ~1GB cache (negative = KB, so 262144KB = 256MB)
		//   Note: SQLite cache_size is in pages (4KB each) or negative for KB
		// - _foreign_keys=ON: Enforce foreign key constraints
		dbURL = "file:./data/db/sqlite/app.db?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-262144&_foreign_keys=ON"
	}

	// Prefork mode: multiple processes for high-load (only with PostgreSQL!)
	// SQLite doesn't support multi-process access
	prefork := os.Getenv("PREFORK") == "true"
	if prefork {
		log.Info().Msg("Prefork mode enabled - spawning multiple processes")
	}

	// Detect database type
	isSQLite := strings.HasPrefix(dbURL, "file:") || strings.HasPrefix(dbURL, "sqlite:")
	isPostgres := strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://")

	// Connect to database with appropriate driver
	var db *gorm.DB
	var err error

	if isPostgres {
		log.Info().Msg("Connecting to PostgreSQL database")
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	} else {
		log.Info().Msg("Connecting to SQLite database with WAL mode")
		db, err = gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database connection")
	}

	if isSQLite {
		// SQLite optimizations:
		// - Single connection for writes (SQLite limitation)
		// - WAL mode allows concurrent reads
		sqlDB.SetMaxOpenConns(1) // SQLite works best with single write connection
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(0) // Keep connection open forever

		// Verify WAL mode is enabled
		var journalMode string
		db.Raw("PRAGMA journal_mode").Scan(&journalMode)
		log.Info().Str("journal_mode", journalMode).Msg("SQLite configuration")
	} else {
		// PostgreSQL: use connection pooling
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.PasswordResetToken{}, &models.AppSettings{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Seed default settings if not exist
	var settingsCount int64
	db.Model(&models.AppSettings{}).Count(&settingsCount)
	if settingsCount == 0 {
		defaultSettings := models.DefaultSettings()
		for _, setting := range defaultSettings {
			db.Create(&setting)
		}
		log.Info().Msg("Default settings seeded")
	}

	// Create Fiber app
	fiberConfig := fiber.Config{
		Prefork:      prefork, // Enable with PREFORK=true (requires PostgreSQL!)
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		BodyLimit:    10 * 1024 * 1024, // 10MB
		ErrorHandler: errorHandler,
	}

	// Trust reverse proxy headers (nginx, Cloudflare) for correct c.IP()
	// TRUSTED_PROXIES=127.0.0.1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16
	if trustedProxies := os.Getenv("TRUSTED_PROXIES"); trustedProxies != "" {
		fiberConfig.EnableTrustedProxyCheck = true
		fiberConfig.TrustedProxies = strings.Split(trustedProxies, ",")
		fiberConfig.ProxyHeader = fiber.HeaderXForwardedFor
		log.Info().Strs("proxies", fiberConfig.TrustedProxies).Msg("Trusted proxies configured")
	}

	app := fiber.New(fiberConfig)

	// Global middleware
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Fast compression for high-load
	}))
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.HelmetMiddleware())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.RateLimiterMiddleware())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${locals:requestId}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// ==========================================================================
	// Services
	// ==========================================================================

	// Auth service
	authService := services.NewAuthService(db)

	// Email service (use MockSender in development, SMTPSender in production)
	var emailSender email.Sender
	if os.Getenv("NODE_ENV") == "production" {
		emailSender = email.NewSMTPSender(email.Config{
			FromName:     os.Getenv("SMTP_FROM_NAME"),
			FromAddress:  os.Getenv("SMTP_FROM_ADDRESS"),
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     587, // TLS port
			SMTPUser:     os.Getenv("SMTP_USER"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			SMTPUseTLS:   true,
		})
		log.Info().Msg("Email service: SMTP")
	} else {
		emailSender = email.NewMockSender(email.Config{})
		log.Info().Msg("Email service: Mock (emails logged to console)")
	}

	// Password reset service
	passwordResetService := services.NewPasswordResetService(db, emailSender)

	// Storage service (local by default, S3 when configured)
	var storageService storage.Storage
	if os.Getenv("S3_BUCKET") != "" {
		s3Storage, err := storage.NewS3Storage(
			os.Getenv("S3_BUCKET"),
			os.Getenv("S3_REGION"),
			os.Getenv("S3_ENDPOINT"),
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize S3 storage")
		}
		storageService = s3Storage
		log.Info().Str("bucket", os.Getenv("S3_BUCKET")).Msg("Storage service: S3")
	} else {
		localStorage, err := storage.NewLocalStorage("./data/uploads", "/uploads")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize local storage")
		}
		storageService = localStorage
		log.Info().Msg("Storage service: Local (./data/uploads)")
	}

	// Upload service
	uploadService := upload.NewService(storageService, upload.DefaultConfig())

	// ==========================================================================
	// Handlers
	// ==========================================================================
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler(db)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	// Health routes
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// ==========================================================================
	// API Routes
	// ==========================================================================
	// Current: /api/auth/*
	// For versioning, change to: /api/v1/auth/*
	//
	// To enable versioning:
	// 1. Change "/api" to "/api/v1" below
	// 2. Update frontend API client base URL
	// 3. Optionally keep "/api" as alias to latest version
	// ==========================================================================
	api := app.Group("/api")

	// Uncomment for explicit versioning:
	// v1 := app.Group("/api/v1")
	// api := v1 // Use v1 as the main API group

	// Auth routes: /api/auth/*
	auth := api.Group("/auth")
	auth.Post("/register", middleware.RegisterRateLimiter(), authHandler.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", middleware.AuthMiddleware(), authHandler.Me)
	auth.Put("/profile", middleware.AuthMiddleware(), authHandler.UpdateProfile)
	auth.Put("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword)

	// Password reset routes: /api/auth/*
	auth.Post("/forgot-password", passwordResetHandler.ForgotPassword)
	auth.Post("/validate-reset-token", passwordResetHandler.ValidateToken)
	auth.Post("/reset-password", passwordResetHandler.ResetPassword)

	// Upload routes: /api/upload/*
	uploads := api.Group("/upload")
	uploads.Post("/", middleware.AuthMiddleware(), uploadHandler.UploadSingle)
	uploads.Post("/multiple", middleware.AuthMiddleware(), uploadHandler.UploadMultiple)
	uploads.Delete("/*", middleware.AuthMiddleware(), uploadHandler.Delete)

	// Serve uploaded files (local storage only)
	app.Static("/uploads", "./data/uploads")

	// ==========================================================================
	// Admin Routes: /api/admin/*
	// ==========================================================================
	// Admin services
	dashboardService := adminServices.NewDashboardService(db)
	usersService := adminServices.NewUsersService(db)
	settingsService := adminServices.NewSettingsService(db)

	// Admin handlers
	dashboardHandler := adminHandlers.NewDashboardHandler(dashboardService)
	usersHandler := adminHandlers.NewUsersHandler(usersService)
	filesHandler := adminHandlers.NewFilesHandler("./data/uploads")
	settingsHandler := adminHandlers.NewSettingsHandler(settingsService)

	// Admin routes group with auth + admin middleware
	adminGroup := api.Group("/admin", middleware.AuthMiddleware(), middleware.AdminOnly(db))

	// Dashboard
	adminGroup.Get("/dashboard", dashboardHandler.GetStats)

	// Users CRUD
	adminGroup.Get("/users", usersHandler.List)
	adminGroup.Get("/users/:id", usersHandler.Get)
	adminGroup.Post("/users", usersHandler.Create)
	adminGroup.Put("/users/:id", usersHandler.Update)
	adminGroup.Delete("/users/:id", usersHandler.Delete)

	// Files
	adminGroup.Get("/files", filesHandler.List)
	adminGroup.Delete("/files/*", filesHandler.Delete)

	// Settings
	adminGroup.Get("/settings", settingsHandler.GetAll)
	adminGroup.Get("/settings/:key", settingsHandler.Get)
	adminGroup.Put("/settings/:key", settingsHandler.Update)
	adminGroup.Put("/settings", settingsHandler.UpdateBatch)

	// ==========================================================================
	// Add your routes here
	// ==========================================================================
	// Example:
	// users := api.Group("/users")
	// users.Get("/", middleware.AuthMiddleware(), userHandler.List)
	// users.Get("/:id", middleware.AuthMiddleware(), userHandler.Get)
	// ==========================================================================

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Resource not found",
			},
			"meta": fiber.Map{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"requestId": c.Locals("requestId"),
			},
		})
	})

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info().Msg("Shutting down gracefully...")
		app.Shutdown()
	}()

	// Start server
	addr := host + ":" + port
	log.Info().Str("addr", addr).Msg("Server starting")

	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	if os.Getenv("NODE_ENV") == "development" {
		message = err.Error()
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": message,
		},
		"meta": fiber.Map{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"requestId": c.Locals("requestId"),
		},
	})
}
