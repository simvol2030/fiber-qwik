package utils

import (
	"os"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	// Set test environment
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "15m")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	payload := JWTPayload{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	token, err := GenerateAccessToken(payload)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateAccessToken returned empty token")
	}

	// Token should have 3 parts separated by dots (header.payload.signature)
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("Token should have 3 parts (2 dots), got %d dots", parts)
	}
}

func TestVerifyAccessToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "15m")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	payload := JWTPayload{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	token, err := GenerateAccessToken(payload)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	verified, err := VerifyAccessToken(token)
	if err != nil {
		t.Fatalf("VerifyAccessToken failed: %v", err)
	}

	if verified.UserID != payload.UserID {
		t.Errorf("UserID mismatch: got %s, want %s", verified.UserID, payload.UserID)
	}

	if verified.Email != payload.Email {
		t.Errorf("Email mismatch: got %s, want %s", verified.Email, payload.Email)
	}
}

func TestVerifyAccessTokenInvalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	defer os.Unsetenv("JWT_SECRET")

	_, err := VerifyAccessToken("invalid-token")
	if err == nil {
		t.Error("VerifyAccessToken should fail for invalid token")
	}
}

func TestVerifyAccessTokenWrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "original-secret-that-is-32-chars!")
	os.Setenv("JWT_EXPIRES_IN", "15m")

	payload := JWTPayload{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	token, err := GenerateAccessToken(payload)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	// Change secret
	os.Setenv("JWT_SECRET", "different-secret-also-32-chars!!")

	_, err = VerifyAccessToken(token)
	if err == nil {
		t.Error("VerifyAccessToken should fail with wrong secret")
	}

	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRES_IN")
}

func TestVerifyAccessTokenExpired(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "1ms") // Very short expiration
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	payload := JWTPayload{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	token, err := GenerateAccessToken(payload)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err = VerifyAccessToken(token)
	if err == nil {
		t.Error("VerifyAccessToken should fail for expired token")
	}
}

func TestGetExpiresInSeconds(t *testing.T) {
	os.Setenv("JWT_EXPIRES_IN", "15m")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	seconds := GetExpiresInSeconds()
	expected := 15 * 60 // 15 minutes = 900 seconds

	if seconds != expected {
		t.Errorf("GetExpiresInSeconds: got %d, want %d", seconds, expected)
	}
}

func TestGetExpiresInSecondsDefault(t *testing.T) {
	os.Unsetenv("JWT_EXPIRES_IN")

	seconds := GetExpiresInSeconds()
	expected := 15 * 60 // Default 15 minutes

	if seconds != expected {
		t.Errorf("GetExpiresInSeconds default: got %d, want %d", seconds, expected)
	}
}

func TestGetRefreshTokenExpiresDays(t *testing.T) {
	os.Setenv("REFRESH_TOKEN_EXPIRES_DAYS", "14")
	defer os.Unsetenv("REFRESH_TOKEN_EXPIRES_DAYS")

	days := GetRefreshTokenExpiresDays()
	if days != 14 {
		t.Errorf("GetRefreshTokenExpiresDays: got %d, want 14", days)
	}
}

func TestGetRefreshTokenExpiresDaysDefault(t *testing.T) {
	os.Unsetenv("REFRESH_TOKEN_EXPIRES_DAYS")

	days := GetRefreshTokenExpiresDays()
	if days != 7 {
		t.Errorf("GetRefreshTokenExpiresDays default: got %d, want 7", days)
	}
}

func TestGetRefreshTokenExpiresDaysInvalid(t *testing.T) {
	os.Setenv("REFRESH_TOKEN_EXPIRES_DAYS", "not-a-number")
	defer os.Unsetenv("REFRESH_TOKEN_EXPIRES_DAYS")

	days := GetRefreshTokenExpiresDays()
	if days != 7 {
		t.Errorf("GetRefreshTokenExpiresDays with invalid value: got %d, want 7", days)
	}
}

func BenchmarkGenerateAccessToken(b *testing.B) {
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "15m")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	payload := JWTPayload{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateAccessToken(payload)
	}
}

func BenchmarkVerifyAccessToken(b *testing.B) {
	os.Setenv("JWT_SECRET", "test-secret-that-is-at-least-32-characters-long")
	os.Setenv("JWT_EXPIRES_IN", "1h")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRES_IN")

	payload := JWTPayload{
		UserID: "user-123",
		Email:  "test@example.com",
	}
	token, _ := GenerateAccessToken(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyAccessToken(token)
	}
}
