package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}

	// Hash should be bcrypt format (starts with $2a$ or $2b$)
	if len(hash) < 60 {
		t.Errorf("Hash length should be at least 60, got %d", len(hash))
	}
}

func TestHashPasswordDifferentResults(t *testing.T) {
	password := "testPassword123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Same password should produce different hashes (bcrypt salts)
	if hash1 == hash2 {
		t.Error("Two hashes of the same password should be different (salt)")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Correct password should verify
	if !VerifyPassword(password, hash) {
		t.Error("VerifyPassword should return true for correct password")
	}
}

func TestVerifyPasswordWrong(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword456"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Wrong password should not verify
	if VerifyPassword(wrongPassword, hash) {
		t.Error("VerifyPassword should return false for wrong password")
	}
}

func TestVerifyPasswordEmptyHash(t *testing.T) {
	// Empty hash should not verify
	if VerifyPassword("password", "") {
		t.Error("VerifyPassword should return false for empty hash")
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	// Invalid hash should not verify
	if VerifyPassword("password", "not-a-valid-hash") {
		t.Error("VerifyPassword should return false for invalid hash")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	// Empty password should still hash (bcrypt allows it)
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword failed for empty string: %v", err)
	}

	if !VerifyPassword("", hash) {
		t.Error("Empty password should verify against its hash")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testPassword123"
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "testPassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword(password, hash)
	}
}
