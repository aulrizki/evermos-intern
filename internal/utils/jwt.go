package utils

import (
    "time"
    "os"
    "github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
    UserID  uint   `json:"user_id"`
    IsAdmin bool   `json:"is_admin"`
    Email   string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID uint, email string, isAdmin bool) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    claims := &JWTClaim{
        UserID:  userID,
        Email:   email,
        IsAdmin: isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
