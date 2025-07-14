package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
	"event-registration/internal/core/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type DashboardPLNMobileHandler struct {
	service *service.DashboardPLNMobileService
	handler *common.Handler
}

// NewDashboardPLNMobileHandler creates a new DashboardPLNMobileHandler
func NewDashboardPLNMobileHandler(service *service.DashboardPLNMobileService, handler *common.Handler) *DashboardPLNMobileHandler {
	return &DashboardPLNMobileHandler{
		service: service,
		handler: handler,
	}
}

// @Summary Dashboard PLN Mobile
// @Description Dashboard PLN Mobile
// @Tags Plnmobile
// @Accept json
// @Produce json
// @Param request query request.DashboardPLNMobileRequest false "..."
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string][]string
// @Router /dashboard-plnmobile/pengguna [get]
func (h *DashboardPLNMobileHandler) Summary(c *fiber.Ctx) error {
	request := new(request.DashboardPLNMobileRequest)

	if err := c.QueryParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(
			c,
			constant.VALIDATION_ERROR,
			h.handler.Validator.ValidationErrors(err),
		)
	}

	// res, lastUpdated, err := h.Service().SummaryPengguna(&req)
	// if err != nil {
	// 	ResponseError(c, err.Error(), helpers.GetLogMessage("vccFunctionFailed"))
	// 	return
	// }

	// ResponseSuccessLastUpdate(c, ApiResponse{
	// 	Data:        res,
	// 	LastUpdated: lastUpdated,
	// 	Meta:        req,
	// })

	return h.handler.ResponseSuccess(c, nil)
}
