package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

type Claims struct {
	JWTPayload
	jwt.RegisteredClaims
}

const minJWTSecretLength = 32

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")

	// In production, JWT_SECRET must be set and at least 32 characters
	if os.Getenv("NODE_ENV") == "production" {
		if secret == "" {
			panic("JWT_SECRET environment variable is required in production")
		}
		if len(secret) < minJWTSecretLength {
			panic("JWT_SECRET must be at least 32 characters in production")
		}
	}

	// Fallback for development only
	if secret == "" {
		secret = "dev-secret-change-in-production-min-32-chars"
	}

	return []byte(secret)
}

func getJWTExpiresIn() time.Duration {
	expiresIn := os.Getenv("JWT_EXPIRES_IN")
	if expiresIn == "" {
		return 15 * time.Minute
	}

	// Parse duration like "15m", "1h", "7d"
	duration, err := time.ParseDuration(expiresIn)
	if err == nil {
		return duration
	}

	return 15 * time.Minute
}

func GenerateAccessToken(payload JWTPayload) (string, error) {
	expiresIn := getJWTExpiresIn()

	claims := Claims{
		JWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func VerifyAccessToken(tokenString string) (*JWTPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.JWTPayload, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func GetExpiresInSeconds() int {
	return int(getJWTExpiresIn().Seconds())
}

func GetRefreshTokenExpiresDays() int {
	daysStr := os.Getenv("REFRESH_TOKEN_EXPIRES_DAYS")
	if daysStr == "" {
		return 7
	}
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 7
	}
	return days
}
