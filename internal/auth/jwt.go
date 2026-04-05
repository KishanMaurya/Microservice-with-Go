package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte("super-secret-key")

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}

// 🔥 Access Token (short-lived)
func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}

// 🔥 Refresh Token (long-lived)
func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})
}