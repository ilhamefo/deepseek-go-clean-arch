package handler

import (
	"event-registration/internal/core/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthService(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// GetUser godoc
// @Summary Get a user by ID
// @Description get user by ID
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /auth/{id} [get]
// func (h *AuthHandler) GetUser(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	return c.JSON(fiber.Map{"id": id, "name": "John Doe"})
// }

// GetUser godoc
// @Summary Get Uri
// @Description Get Uri for google oauth login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Router /auth/login-url [get]
func (h *AuthHandler) GetLoginUrl(c *fiber.Ctx) error {
	url := h.service.GetLoginUrl()

	return responseSuccess(c, fiber.Map{"url": url})
}

// GetUser godoc
// @Summary Get Uri
// @Description Get Uri for google oauth login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Router /auth/login-url [get]
func (h *AuthHandler) GoogleHandleCallback(c *fiber.Ctx) error {
	err := h.service.GoogleHandleCallback(c.Context(), c.Query("code"))
	if err != nil {
		return responseError(c, err.Error())
	}

	return responseSuccess(c, "OK")
}
