package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/core/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type GarminDashboardHandler struct {
	service *service.GarminDashboardService
	handler *common.Handler
	logger  *zap.Logger
}

// NewGarminDashboardHandler creates a new GarminDashboardHandler
func NewGarminDashboardHandler(service *service.GarminDashboardService, logger *zap.Logger, handler *common.Handler) *GarminDashboardHandler {
	return &GarminDashboardHandler{
		service: service,
		handler: handler,
		logger:  logger,
	}
}

// Get heart rate godoc
// @Summary Get heart rate
// @Description This endpoint is used to Get heart rate.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Router /dashboard/heart-rate [get]
func (h *GarminDashboardHandler) HeartRate(c *fiber.Ctx) error {
	res, err := h.service.HeartRate(c.Context())
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, res)
}
