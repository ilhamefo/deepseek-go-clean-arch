package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
	"event-registration/internal/core/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
	handler *common.Handler
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *service.UserService, handler *common.Handler) *UserHandler {
	return &UserHandler{
		service: service,
		handler: handler,
	}
}

// Search godoc
// @Summary Search
// @Description This endpoint is used to search users by keyword.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param request query request.SearchRequest false "..."
// @Success 200 {object} []domain.UserVCC
// @Router /search-user [get]
func (h *UserHandler) Search(c *fiber.Ctx) error {
	request := new(request.SearchRequest)

	if err := c.QueryParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	users, err := h.service.Search(c.Context(), request.Keyword)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, users)
}

// Roles godoc
// @Summary Roles
// @Description This endpoint is used to get all roles.
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Role
// @Router /roles [get]
func (h *UserHandler) Roles(c *fiber.Ctx) error {
	roles, err := h.service.Roles()
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, roles)
}

// Unit godoc
// @Summary Unit
// @Description This endpoint is used to get all units by level.
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.UnitName
// @Router /units [get]
func (h *UserHandler) GetUnits(c *fiber.Ctx) error {

	request := new(request.GetUnitRequest)

	if err := c.QueryParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	units, err := h.service.GetUnits(request.Level)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, units)
}

// Update godoc
// @Summary Update
// @Description This endpoint is used to update user information.
// @Tags Users
// @Param id path int true "User ID"
// @Param request body request.UpdateUserRequest false "..."
// @Accept  json
// @Produce  json
// @Router /update/{id} [post]
func (h *UserHandler) Update(c *fiber.Ctx) error {

	request := new(request.UpdateUserRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := c.ParamsParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	err := h.service.Update(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, true)
}

// Check Health godoc
// @Summary Check Health
// @Description This endpoint is used to check the health of Meilisearch.
// @Tags Users
// @Accept  json
// @Produce  json
// @Router /meili-health [get]
func (h *UserHandler) CheckHealthMeilisearch(c *fiber.Ctx) error {
	err := h.service.CheckHealthMeilisearch()
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return h.handler.ResponseSuccess(c, nil)
}
