package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/aulrizki/evermos-intern/internal/utils"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"status":  false,
				"message": "missing token",
			})
		}

		// Hapus prefix "Bearer " kalau ada
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// Parse token menggunakan helper utils.ParseJWT
		token, err := utils.ParseJWT(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"status":  false,
				"message": "invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"status":  false,
				"message": "invalid token claims",
			})
		}

		// Ambil nilai user_id dan is_admin dari claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"status":  false,
				"message": "invalid user id in token",
			})
		}

		isAdmin, _ := claims["is_admin"].(bool)

		// Simpan ke Fiber Locals
		c.Locals("user_id", uint(userIDFloat))
		c.Locals("is_admin", isAdmin)

		return c.Next()
	}
}
