package handler

import (
	"event-registration/internal/constant"
	"event-registration/internal/core/service"
	validate "event-registration/internal/infrastructure/validator"
	"event-registration/internal/request"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service   *service.AuthService
	validator *validate.Validator
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(service *service.AuthService, validator *validate.Validator) *AuthHandler {
	return &AuthHandler{service: service, validator: validator}
}

// GetUser godoc
// @Summary Get Uri
// @Description Get Uri for google oauth login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Router /auth/google/login-url [get]
func (h *AuthHandler) GetLoginUrl(c *fiber.Ctx) error {
	url, token, err := h.service.GetLoginUrl()
	if err != nil {
		return responseError(c, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    token,
		MaxAge:   300,   // 5 mins
		HTTPOnly: false, // set it to true in production
		Secure:   true,
		SameSite: "Strict",
	})

	return responseSuccess(c, fiber.Map{"url": url})
}

// GetUser godoc
// @Summary Get Uri
// @Description Get Uri for google oauth login. Requires oauth_state cookie
// @Tags Auth
// @Accept  json
// @Param Cookie header string true "Cookie header: oauth_state=xxxx"
// @Param code query string true "ABCDE"
// @Param state query string true "ABCDE"
// @Produce  json
// @Success 200 {object} User
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleHandleCallback(c *fiber.Ctx) error {
	request := new(request.GoogleCallbackRequest)

	if err := c.QueryParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constant.INVALID_REQUEST_BODY,
		})
	}

	if err := h.validator.Struct(request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error_validations": h.validator.ValidationErrors(err),
		})
	}

	request.StateCookie = c.Cookies("oauth_state")

	err := h.service.GoogleHandleCallback(c.Context(), request)
	if err != nil {
		return responseError(c, err.Error())
	}

	return responseSuccess(c, "OK")
}
