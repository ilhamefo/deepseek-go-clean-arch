package handler

import (
	"event-registration/internal/core/service"

	"github.com/gofiber/fiber/v2"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	service *service.EventService
}

// NewEventHandler creates a new EventHandler
func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// RegisterEvent godoc
// @Summary Register for an event
// @Description Register a user for a specific event
// @Tags events
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} domain.Event
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{event_id} [get]
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

// GetUser godoc
// @Summary Get a user by ID
// @Description get user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /users/{id} [get]
func (h *EventHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"id": id, "name": "John Doe"})
}

// User represents a user model
type User struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"John Doe"`
}

// HTTPError represents an error response
type HTTPError struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad request"`
}
