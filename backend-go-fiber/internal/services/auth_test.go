package services

import (
	"os"
	"testing"
	"time"

	"backend-go-fiber/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// setupTestEnv sets up environment variables for testing
func setupTestEnv() func() {
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "15m")
	os.Setenv("REFRESH_TOKEN_EXPIRES_DAYS", "7")

	return func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRES_IN")
		os.Unsetenv("REFRESH_TOKEN_EXPIRES_DAYS")
	}
}

func TestNewAuthService(t *testing.T) {
	db := setupTestDB(t)
	service := NewAuthService(db)

	if service == nil {
		t.Error("NewAuthService returned nil")
	}

	if service.db != db {
		t.Error("AuthService.db not set correctly")
	}
}

func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	input := RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	result, err := service.Register(input)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if result == nil {
		t.Fatal("Register returned nil result")
	}

	if result.User.Email != input.Email {
		t.Errorf("User email mismatch: got %s, want %s", result.User.Email, input.Email)
	}

	if result.AccessToken == "" {
		t.Error("AccessToken should not be empty")
	}

	if result.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}

	// Verify user was created in database
	var user models.User
	err = db.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		t.Errorf("User not found in database: %v", err)
	}
}

func TestRegisterWithName(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)
	name := "Test User"

	input := RegisterInput{
		Email:    "testname@example.com",
		Password: "password123",
		Name:     &name,
	}

	result, err := service.Register(input)
	if err != nil {
		t.Fatalf("Register with name failed: %v", err)
	}

	if result.User.Name == nil || *result.User.Name != name {
		t.Errorf("User name mismatch: got %v, want %s", result.User.Name, name)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	input := RegisterInput{
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	// First registration should succeed
	_, err := service.Register(input)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Second registration should fail
	_, err = service.Register(input)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "user already exists" {
		t.Errorf("Expected 'user already exists' error, got: %v", err)
	}
}

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// First register a user
	registerInput := RegisterInput{
		Email:    "login@example.com",
		Password: "password123",
	}
	_, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Now login
	loginInput := LoginInput{
		Email:    "login@example.com",
		Password: "password123",
	}

	result, err := service.Login(loginInput)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if result == nil {
		t.Fatal("Login returned nil result")
	}

	if result.User.Email != loginInput.Email {
		t.Errorf("User email mismatch: got %s, want %s", result.User.Email, loginInput.Email)
	}

	if result.AccessToken == "" {
		t.Error("AccessToken should not be empty")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// Register a user
	registerInput := RegisterInput{
		Email:    "wrongpass@example.com",
		Password: "correctpassword",
	}
	_, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Try to login with wrong password
	loginInput := LoginInput{
		Email:    "wrongpass@example.com",
		Password: "wrongpassword",
	}

	_, err = service.Login(loginInput)
	if err == nil {
		t.Error("Expected error for wrong password")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got: %v", err)
	}
}

func TestLoginNonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	loginInput := LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := service.Login(loginInput)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got: %v", err)
	}
}

func TestCreateRefreshToken(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// Create a user first
	registerInput := RegisterInput{
		Email:    "refresh@example.com",
		Password: "password123",
	}
	result, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Create refresh token
	token, err := service.CreateRefreshToken(result.User.ID)
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}

	if token == "" {
		t.Error("Refresh token should not be empty")
	}

	// Verify token was created in database
	var refreshToken models.RefreshToken
	err = db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		t.Errorf("Refresh token not found in database: %v", err)
	}

	if refreshToken.UserID != result.User.ID {
		t.Errorf("UserID mismatch: got %s, want %s", refreshToken.UserID, result.User.ID)
	}
}

func TestRefreshAccessToken(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// Create a user and refresh token
	registerInput := RegisterInput{
		Email:    "refreshaccess@example.com",
		Password: "password123",
	}
	result, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	refreshToken, err := service.CreateRefreshToken(result.User.ID)
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}

	// Refresh the access token
	newTokens, err := service.RefreshAccessToken(refreshToken)
	if err != nil {
		t.Fatalf("RefreshAccessToken failed: %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Error("New access token should not be empty")
	}

	if newTokens.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}
}

func TestRefreshAccessTokenInvalid(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	_, err := service.RefreshAccessToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid refresh token")
	}

	if err.Error() != "invalid refresh token" {
		t.Errorf("Expected 'invalid refresh token' error, got: %v", err)
	}
}

func TestRefreshAccessTokenExpired(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// Create a user
	registerInput := RegisterInput{
		Email:    "expired@example.com",
		Password: "password123",
	}
	result, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Create an expired refresh token directly in DB
	expiredToken := models.RefreshToken{
		Token:     "expired-token-123",
		UserID:    result.User.ID,
		ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired yesterday
	}
	db.Create(&expiredToken)

	// Try to use the expired token
	_, err = service.RefreshAccessToken("expired-token-123")
	if err == nil {
		t.Error("Expected error for expired refresh token")
	}

	if err.Error() != "refresh token expired" {
		t.Errorf("Expected 'refresh token expired' error, got: %v", err)
	}

	// Verify the expired token was deleted
	var count int64
	db.Model(&models.RefreshToken{}).Where("token = ?", "expired-token-123").Count(&count)
	if count != 0 {
		t.Error("Expired token should have been deleted")
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// Create a user and refresh token
	registerInput := RegisterInput{
		Email:    "revoke@example.com",
		Password: "password123",
	}
	result, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	refreshToken, err := service.CreateRefreshToken(result.User.ID)
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}

	// Revoke the token
	err = service.RevokeRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("RevokeRefreshToken failed: %v", err)
	}

	// Verify token was deleted
	var count int64
	db.Model(&models.RefreshToken{}).Where("token = ?", refreshToken).Count(&count)
	if count != 0 {
		t.Error("Refresh token should have been deleted")
	}
}

func TestGetUserByID(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	// Create a user
	registerInput := RegisterInput{
		Email:    "getuser@example.com",
		Password: "password123",
	}
	result, err := service.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Get user by ID
	user, err := service.GetUserByID(result.User.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if user.ID != result.User.ID {
		t.Errorf("User ID mismatch: got %s, want %s", user.ID, result.User.ID)
	}

	if user.Email != result.User.Email {
		t.Errorf("User email mismatch: got %s, want %s", user.Email, result.User.Email)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	cleanup := setupTestEnv()
	defer cleanup()

	service := NewAuthService(db)

	_, err := service.GetUserByID("non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent user ID")
	}
}

// Benchmark tests
func BenchmarkRegister(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db.AutoMigrate(&models.User{}, &models.RefreshToken{})

	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "15m")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	service := NewAuthService(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := RegisterInput{
			Email:    "benchmark" + string(rune(i)) + "@example.com",
			Password: "password123",
		}
		service.Register(input)
	}
}

func BenchmarkLogin(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db.AutoMigrate(&models.User{}, &models.RefreshToken{})

	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "15m")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	service := NewAuthService(db)

	// Create a user to login with
	service.Register(RegisterInput{
		Email:    "benchmark@example.com",
		Password: "password123",
	})

	loginInput := LoginInput{
		Email:    "benchmark@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Login(loginInput)
	}
}
