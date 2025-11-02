package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key untuk JWT
var jwtSecret = []byte(getSecret())

func getSecret() string {
	// Ambil dari environment variable kalau ada, kalau tidak pakai default
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "evermos_secret_key" // fallback default
	}
	return secret
}

// 🟩 GenerateJWT membuat token JWT dengan claim user_id, email, dan is_admin
func GenerateJWT(userID uint, email string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // token valid 1 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 🟩 ParseJWT untuk memvalidasi token dan membaca claim-nya
func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
