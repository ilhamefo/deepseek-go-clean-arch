package handler

import (
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/request"
	"event-registration/internal/core/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
	handler *common.Handler
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(service *service.AuthService, handler *common.Handler) *AuthHandler {
	return &AuthHandler{
		service: service,
		handler: handler,
	}
}

// GetUser godoc
// @Summary Get Uri
// @Description Get Uri for google oauth login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.User
// @Router /auth/google/login-url [get]
func (h *AuthHandler) GetLoginUrl(c *fiber.Ctx) error {
	url, token, err := h.service.GetLoginUrl()
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    token,
		MaxAge:   300,   // 5 mins
		HTTPOnly: false, // set it to true in production
		Secure:   true,
		SameSite: "Strict",
	})

	return h.handler.ResponseSuccess(c, fiber.Map{"url": url})
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
// @Success 200 {object} domain.User
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleHandleCallback(c *fiber.Ctx) error {
	request := new(request.GoogleCallbackRequest)

	if err := c.QueryParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseWithStatus(c, http.StatusUnprocessableEntity, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	request.StateCookie = c.Cookies("oauth_state")

	accessToken, refreshToken, err := h.service.GoogleHandleCallback(c.Context(), request)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	h.createCookies(c, accessToken, refreshToken)

	return h.handler.ResponseSuccess(c, fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}

// GetUser godoc
// @Summary Protected Route
// @Description test protected route
// @Description Requires authentication
// @Tags Auth
// @Param Cookie header string true "Cookie header: access_token=xxxx"
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.User
// @Router /protected [get]
func (h *AuthHandler) Protected(c *fiber.Ctx) error {
	return h.handler.ResponseSuccess(c, fiber.Map{"user": h.handler.ParseUser(c)})
}

// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.User
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	user := h.handler.ParseUser(c)

	accessToken, refreshToken, err := h.service.GenerateToken(&user)
	if err != nil {
		return h.handler.ResponseWithStatus(c, http.StatusInternalServerError, "Failed to generate access token", nil)
	}

	h.createCookies(c, accessToken, refreshToken)

	return h.handler.ResponseSuccess(c, fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}

// Login godoc
// @Summary Login
// @Description Login with email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body request.LoginRequest false "..."
// @Success 200 {object} domain.User
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	request := new(request.LoginRequest)

	if err := c.BodyParser(request); err != nil {
		return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
	}

	if err := h.handler.Validator.Struct(request); err != nil {
		return h.handler.ResponseWithStatus(c, http.StatusUnprocessableEntity, constant.VALIDATION_ERROR, h.handler.Validator.ValidationErrors(err))
	}

	accessToken, refreshToken, err := h.service.Login(c.Context(), request)
	if err != nil {
		return h.handler.ResponseWithStatus(c, http.StatusInternalServerError, "login_failed", err.Error())
	}

	h.createCookies(c, accessToken, refreshToken)

	return h.handler.ResponseSuccess(c, nil)
}

func (h *AuthHandler) createCookies(c *fiber.Ctx, accessToken, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   60 * 15, // 15 mins
		HTTPOnly: true,    // set it to true in production
		Secure:   true,
		SameSite: "Strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   60 * 60 * 24 * 7, // 1 week
		HTTPOnly: true,             // set it to true in production
		Secure:   true,
		SameSite: "Strict",
	})
}
