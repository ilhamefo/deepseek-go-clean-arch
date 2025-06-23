package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	} else if err != nil {
		msg = err.Error()
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": "failed",
		"error":   msg,
	})
}
