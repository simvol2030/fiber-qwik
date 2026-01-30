package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptCost is set to 12 for production security
// cost=10 allows ~10k hashes/sec on GPU, cost=12 is ~4x slower
// Recommended: 12-14 for production, 10 for testing
const bcryptCost = 12

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
