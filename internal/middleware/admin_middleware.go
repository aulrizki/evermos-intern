package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Middleware untuk memastikan user adalah admin
func IsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			return c.Status(403).JSON(fiber.Map{
				"status":  false,
				"message": "forbidden: admin only",
			})
		}
		return c.Next()
	}
}
