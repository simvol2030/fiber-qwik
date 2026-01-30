package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Configure logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "file:./data/db/sqlite/app.db?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=1000000000&_foreign_keys=ON"
	}

	// Connect to database
	var db *gorm.DB
	var err error

	if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
		log.Info().Msg("Connecting to PostgreSQL database")
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	} else {
		log.Info().Msg("Connecting to SQLite database")
		db, err = gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Run migrations
	log.Info().Msg("Running migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Seed data
	log.Info().Msg("Seeding database...")
	if err := seed(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to seed database")
	}

	log.Info().Msg("✅ Database seeded successfully!")
}

func seed(db *gorm.DB) error {
	// Check if admin user already exists
	var existingAdmin models.User
	if err := db.Where("email = ?", "admin@example.com").First(&existingAdmin).Error; err == nil {
		log.Info().Msg("Admin user already exists, skipping...")
		return nil
	}

	// Create admin user
	adminPassword, _ := utils.HashPassword("admin123")
	admin := &models.User{
		Email:        "admin@example.com",
		PasswordHash: adminPassword,
		Name:         strPtr("Admin User"),
		Role:         models.RoleAdmin,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	log.Info().Str("email", admin.Email).Str("password", "admin123").Msg("Created admin user")

	// Create test user
	var existingUser models.User
	if err := db.Where("email = ?", "user@example.com").First(&existingUser).Error; err == nil {
		log.Info().Msg("Test user already exists, skipping...")
		return nil
	}

	userPassword, _ := utils.HashPassword("user1234")
	user := &models.User{
		Email:        "user@example.com",
		PasswordHash: userPassword,
		Name:         strPtr("Test User"),
		Role:         models.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}
	log.Info().Str("email", user.Email).Str("password", "user1234").Msg("Created test user")

	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📦 Seeded Users:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("  Admin: admin@example.com / admin123 (role: %s)\n", models.RoleAdmin)
	fmt.Printf("  User:  user@example.com / user1234 (role: %s)\n", models.RoleUser)
	fmt.Println(strings.Repeat("=", 50))

	return nil
}

func strPtr(s string) *string {
	return &s
}
