package handler

import (
	"event-registration/internal/core/domain"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func responseSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

func parseUser(c *fiber.Ctx) (user domain.User) {
	return c.Locals("user").(domain.User)
}
