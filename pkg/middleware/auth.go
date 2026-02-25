package middleware

import (
	"strings"

	"fiber.com/session-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := ""

		tokenStr = c.Cookies("auth_token")

		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authentication token")
		}

		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		c.Locals("userId", claims.UserID)
		c.Locals("userName", claims.UserName)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
