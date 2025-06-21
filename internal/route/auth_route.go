package route

import (
	"event-registration/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterAuthRoutes(app *fiber.App, authHandler *handler.AuthHandler) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:     true,
		DocExpansion:    "list",
		WithCredentials: true,
	}))

	auth := app.Group("/auth")
	google := auth.Group("/google")
	google.Get("/login-url", authHandler.GetLoginUrl)
	google.Get("/callback", authHandler.GoogleHandleCallback)
}
