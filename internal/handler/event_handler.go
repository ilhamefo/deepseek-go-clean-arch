package handler

import (
	"event-registration/internal/core/service"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) RegisterEvent(c *fiber.Ctx) error {
	var req struct {
		EventID string `json:"event_id"`
		UserID  string `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.service.RegisterEvent(req.EventID, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Registration successful"})
}
