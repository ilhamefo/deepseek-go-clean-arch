package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/core/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type GarminHandler struct {
	service *service.GarminService
	handler *common.Handler
}

// NewGarminHandler creates a new GarminHandler
func NewGarminHandler(service *service.GarminService, handler *common.Handler) *GarminHandler {
	return &GarminHandler{
		service: service,
		handler: handler,
	}
}

// Search godoc
// @Summary Search
// @Description This endpoint is used to search users by keyword.
// @Tags Garmin
// @Accept  json
// @Produce  json
// @Router /refresh [get]
func (h *GarminHandler) Refresh(c *fiber.Ctx) error {
	_, err := h.service.Refresh(c.Context())
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, true)
}
