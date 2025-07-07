package middleware

import (
	"event-registration/internal/common/constant"
	"event-registration/internal/core/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return m.handler.ResponseWithStatus(c, fiber.StatusUnauthorized, "access_token_is_required", nil)
		}

		token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			return m.handler.ResponseWithStatus(c, fiber.StatusUnauthorized, "invalid_access_token", nil)
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return m.handler.ResponseWithStatus(c, fiber.StatusUnauthorized, "invalid_token_claims", nil)
		}

		tokenType, ok := (*claims)["type"].(string)
		if !ok || tokenType != constant.ACCESS_TOKEN {
			return m.handler.ResponseWithStatus(c, fiber.StatusUnauthorized, "invalid_token_type", nil)
		}

		isBlacklisted, err := m.sessionService.IsAccessTokenBlacklisted(c.Context(), accessToken)
		if err != nil {
			return m.handler.ResponseWithStatus(c, fiber.StatusInternalServerError, "error_checking_blacklist", nil)
		}

		if isBlacklisted {
			return m.handler.ResponseWithStatus(c, fiber.StatusUnauthorized, "access_token_blacklisted", nil)
		}

		user := domain.User{
			Email: (*claims)["email"].(string),
			ID:    (*claims)["sub"].(string),
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func (m *Middleware) VerifyRefreshToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token")
		if refreshToken == "" {
			return m.handler.ResponseWithStatus(c, http.StatusUnauthorized, "refresh_token_is_required", nil)
		}

		token, err := jwt.ParseWithClaims(refreshToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			return m.handler.ResponseWithStatus(c, http.StatusUnauthorized, "invalid_refresh_token", nil)
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return m.handler.ResponseWithStatus(c, http.StatusUnauthorized, "invalid_token_claims", nil)
		}

		tokenType, ok := (*claims)["type"].(string)
		if !ok || tokenType != constant.REFRESH_TOKEN {
			return m.handler.ResponseWithStatus(c, http.StatusUnauthorized, "invalid_token_type", nil)
		}

		if !m.sessionService.IsSessionValid(c.Context(), refreshToken) {
			return m.handler.ResponseWithStatus(c, http.StatusUnauthorized, "session_expired_or_invalid", nil)
		}

		user := domain.User{
			Email: (*claims)["email"].(string),
			ID:    (*claims)["sub"].(string),
		}

		c.Locals("user", user)
		return c.Next()
	}
}
