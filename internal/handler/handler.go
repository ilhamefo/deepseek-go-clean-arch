package handler

import (
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
