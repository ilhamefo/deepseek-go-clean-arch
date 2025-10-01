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
// @Param request body request.GarminBasicRequest false "..."
// @Produce  json
// @Router /refresh [post]
func (h *GarminHandler) Refresh(c *fiber.Ctx) error {
	request := new(request.GarminBasicRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	err := h.service.Refresh(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, nil)
}

// Splits godoc
// @Summary Splits
// @Description This endpoint is used to Splits Garmin activities.
// @Tags Garmin
// @Accept  json
// @Param request body request.GarminBasicRequest false "..."
// @Param activityID path int true "Activity ID"
// @Produce  json
// @Router /splits/{activityID} [post]
func (h *GarminHandler) Splits(c *fiber.Ctx) error {
	requestActivityID := new(request.ActivityRequest)
	request := new(request.GarminBasicRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	if err := c.ParamsParser(requestActivityID); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	res, err := h.service.FetchSplits(c.Context(), request, requestActivityID.ActivityID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, res)
}

// Heart Rate By Date godoc
// @Summary Heart Rate By Date
// @Description This endpoint is used to get Garmin heart rate by date.
// @Tags Garmin
// @Accept  json
// @Param request body request.GarminByDateRequest false "..."
// @Produce  json
// @Router /heart-rate-by-date [post]
func (h *GarminHandler) GetHeartRateByDate(c *fiber.Ctx) error {
	request := new(request.GarminByDateRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	err := h.service.HeartRateByDate(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, nil)
}

// User Profile godoc
// @Summary User Profile
// @Description This endpoint is used to get Garmin user profile.
// @Tags Garmin
// @Accept  json
// @Param request body request.GarminBasicRequest false "..."
// @Produce  json
// @Router /user-profile [post]
func (h *GarminHandler) GetUserProfile(c *fiber.Ctx) error {
	request := new(request.GarminBasicRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	response, err := h.service.GetUserSettings(c.Context(), request, true)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, response)
}

// GetStepByDate godoc
// @Summary GetStepByDate
// @Description This endpoint is used to get Garmin step by date.
// @Tags Garmin
// @Accept  json
// @Param request body request.GarminByDateRequest false "..."
// @Produce  json
// @Router /step-by-date [post]
func (h *GarminHandler) GetStepByDate(c *fiber.Ctx) error {
	request := new(request.GarminByDateRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	err := h.service.StepByDate(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, nil)
}

// HRVByDate godoc
// @Summary HRVByDate
// @Description This endpoint is used to get Garmin HRV by date.
// @Tags Garmin
// @Accept  json
// @Param request body request.GarminByDateRequest false "..."
// @Produce  json
// @Router /hrv-by-date [post]
func (h *GarminHandler) HRVByDate(c *fiber.Ctx) error {
	request := new(request.GarminByDateRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	err := h.service.HRVByDate(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, nil)
}

// GetActivityTypes godoc
// @Summary GetActivityTypes
// @Description This endpoint is used to get Garmin activity types.
// @Tags Garmin
// @Accept  json
// @Param request body request.GarminBasicRequest false "..."
// @Produce  json
// @Router /activity-types [post]
func (h *GarminHandler) GetActivityTypes(c *fiber.Ctx) error {
	request := new(request.GarminBasicRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	err := h.service.GetActivityTypes(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, nil)
}
