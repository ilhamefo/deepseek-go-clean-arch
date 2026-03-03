package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
	"event-registration/internal/core/service"
	validate "event-registration/internal/infrastructure/validator"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type GarminDashboardHandler struct {
	service   *service.GarminDashboardService
	handler   *common.Handler
	logger    *zap.Logger
	validator *validate.Validator
}

// NewGarminDashboardHandler creates a new GarminDashboardHandler
func NewGarminDashboardHandler(service *service.GarminDashboardService, logger *zap.Logger, handler *common.Handler, validator *validate.Validator) *GarminDashboardHandler {
	return &GarminDashboardHandler{
		service:   service,
		handler:   handler,
		logger:    logger,
		validator: validator,
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
// @Success 200 {object} object{data=[]domain.ActivityVM,nextCursor=int,hasMore=bool,limit=int} "Successfully retrieved activities"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 422 {object} object{error_validations=object} "Validation error"
// @Param request query request.ActivityDashboardRequest false "..."
// @Router /dashboard/activities [get]
func (h *GarminDashboardHandler) Activities(c *fiber.Ctx) error {

	request := new(request.ActivityDashboardRequest)

	if err := c.QueryParser(request); err != nil {
		h.logger.Error("failed to parse query", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constant.INVALID_REQUEST_BODY,
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error_validations": h.validator.ValidationErrors(err),
		})
	}

	res, err := h.service.GetActivities(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseCursorPagination(c, res)
}
