package middleware

import (
	"event-registration/internal/common/helper"
	"event-registration/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Access token is required",
			})
		}

		token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid access token",
			})
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		helper.PrettyPrint(claims, "claims")

		user := domain.User{
			Email: (*claims)["email"].(string),
			// Name:  (*claims)["name"].(string),
			ID: (*claims)["sub"].(string),
		}

		// set user in context
		c.Locals("user", user)

		return c.Next()
	}
}
