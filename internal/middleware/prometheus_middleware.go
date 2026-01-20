package middleware

import (
	"event-registration/internal/infrastructure/metric"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		metric.ActiveConnections.Inc()
		defer metric.ActiveConnections.Dec()

		err := c.Next()

		duration := time.Since(start)

		metric.RecordHTTPRequest(
			c.Method(),
			c.Route().Path,
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
