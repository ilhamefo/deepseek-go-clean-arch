package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
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

// Refresh godoc
// @Summary Refresh
// @Description This endpoint is used to Refresh Garmin activities.
// @Tags Garmin
// @Accept  json
// @Param request body request.RefreshActivitiesRequest false "..."
// @Produce  json
// @Router /refresh [post]
func (h *GarminHandler) Refresh(c *fiber.Ctx) error {
	request := new(request.RefreshActivitiesRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	res, err := h.service.Refresh(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, res)
}
