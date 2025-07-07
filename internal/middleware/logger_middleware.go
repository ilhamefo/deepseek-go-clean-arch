package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (m *Middleware) NewZapLoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		ip := c.Get("CF-Connecting-IP")
		if ip == "" {
			ip = c.IP()
		}

		method := c.Method()
		path := c.OriginalURL()
		body := c.Body()

		err := c.Next()

		status := c.Response().StatusCode()
		latency := time.Since(start)
		userAgent := c.Get("User-Agent")
		referer := c.Get("Referer")

		logger.Info("incoming_request",
			zap.String("ip", ip),
			zap.String("method", method),
			zap.String("path", path),
			zap.ByteString("body", body),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("user_agent", userAgent),
			zap.String("referer", referer),
		)

		return err
	}
}
