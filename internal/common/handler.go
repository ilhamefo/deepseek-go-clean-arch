package common

import (
	"event-registration/internal/core/domain"
	validate "event-registration/internal/infrastructure/validator"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	Validator *validate.Validator
	logger    *zap.Logger
}

func NewHandler(validator *validate.Validator, logger *zap.Logger) *Handler {
	return &Handler{
		Validator: validator,
		logger:    logger,
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

func (h *Handler) ResponseError(c *fiber.Ctx, status int, message string, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
		"error":   err.Error(),
	})
}

func (h *Handler) ResponseWithStatus(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	})
}

func (h *Handler) ResponsePaginated(c *fiber.Ctx, data interface{}, page, pageSize, total int) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":    http.StatusOK,
		"message":   "success",
		"data":      data,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}
