package common

import (
	"event-registration/internal/core/domain"
	validate "event-registration/internal/infrastructure/validator"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	validator *validate.Validator
}

func NewHandler(validator *validate.Validator) *Handler {
	return &Handler{
		validator: validator,
	}
}

func (h *Handler) ResponseSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

func (h *Handler) ParseUser(c *fiber.Ctx) (user domain.User) {
	return c.Locals("user").(domain.User)
}
