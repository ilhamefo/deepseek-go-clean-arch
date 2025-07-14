package route

import (
	"event-registration/internal/handler"
	"event-registration/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterDashboardPLNMobileRoutes(app *fiber.App, h *handler.DashboardPLNMobileHandler, m *middleware.Middleware) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:     true,
		DocExpansion:    "list",
		WithCredentials: true,
	}))

	// dashboard := app.Group("/dashboard/plnmobile")
	// dashboard.Get("/summary", m.AuthMiddleware(), h.SummaryPengguna)
}
