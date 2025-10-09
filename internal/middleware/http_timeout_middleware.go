package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) HTTPTimeoutMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), time.Duration(m.cfg.Timeout)*time.Second)
		defer cancel()

		c.SetUserContext(ctx)

		done := make(chan error, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					done <- fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
				}
			}()

			done <- c.Next()
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return fiber.NewError(fiber.StatusRequestTimeout, "Request timeout")
		}
	}
}
