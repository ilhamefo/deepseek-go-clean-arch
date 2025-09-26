package route

import (
	"event-registration/internal/handler"
	"event-registration/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterGarminRoutes(app *fiber.App, h *handler.GarminHandler, m *middleware.Middleware) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:     true,
		DocExpansion:    "list",
		WithCredentials: true,
	}))

	app.Post("/refresh", h.Refresh)
	app.Post("/heart-rate-by-date", h.GetHeartRateByDate)
	app.Post("/splits/:activityID", h.Splits)
}
