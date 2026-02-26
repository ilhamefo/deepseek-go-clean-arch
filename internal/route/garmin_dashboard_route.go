package route

import (
	"event-registration/internal/handler"
	"event-registration/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterGarminDashboardRoutes(app *fiber.App, h *handler.GarminDashboardHandler, m *middleware.Middleware) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:        true,
		DocExpansion:       "list",
		WithCredentials:    false,
		DisplayOperationId: true,
	}))

	dashboard := app.Group("/dashboard")
	dashboard.Get("/heart-rate", h.HeartRate)
	dashboard.Get("/activities", h.Activities)
}
