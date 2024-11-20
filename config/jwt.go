package config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GenerateJWT generates a JWT token with user ID and expiration time
func GenerateJWT(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}

// ValidateJWT validates the provided JWT token and extracts claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
