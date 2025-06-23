package middleware

import (
	"event-registration/internal/helper"

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

		helper.PrettyPrint(claims, "AuthMiddleware claims =================")

		return c.Next()
	}
}
