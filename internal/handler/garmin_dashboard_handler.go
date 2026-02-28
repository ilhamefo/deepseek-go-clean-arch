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

// Get activities godoc
// @Summary Get activities with cursor pagination
// @Description This endpoint is used to Get activities with cursor-based pagination.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param cursor query int false "Cursor for pagination (activity_id)"
// @Param limit query int false "Limit per page (default: 10, max: 100)"
// @Router /dashboard/activities [get]
func (h *GarminDashboardHandler) Activities(c *fiber.Ctx) error {
	cursor := c.QueryInt("cursor", 0)
	limit := c.QueryInt("limit", 10)

	res, err := h.service.GetActivities(c.Context(), int64(cursor), limit)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseCursorPagination(c, res)
}
