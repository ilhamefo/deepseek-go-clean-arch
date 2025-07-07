package route

import (
	"event-registration/internal/handler"
	"event-registration/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterAuthRoutes(app *fiber.App, authHandler *handler.AuthHandler, m *middleware.Middleware) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:     true,
		DocExpansion:    "list",
		WithCredentials: true,
	}))

	auth := app.Group("/auth")
	auth.Group("/login", authHandler.Login)

	google := auth.Group("/google")
	google.Get("/login-url", authHandler.GetLoginUrl)
	google.Get("/callback", authHandler.GoogleHandleCallback)

	auth.Get("/refresh-token", m.VerifyRefreshToken(), authHandler.RefreshToken)

	authenticated := app.Group("/", m.AuthMiddleware())
	authenticated.Get("/me", m.AuthMiddleware(), authHandler.Protected)
	authenticated.Post("/logout", m.VerifyRefreshToken(), m.AuthMiddleware(), authHandler.Logout)
	authenticated.Post("/logout-all", m.AuthMiddleware(), authHandler.LogoutAllDevices)
}
