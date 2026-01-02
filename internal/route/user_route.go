package route

import (
	"event-registration/internal/handler"
	"event-registration/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterUserRoutes(app *fiber.App, userHandler *handler.UserHandler, m *middleware.Middleware) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:     true,
		DocExpansion:    "list",
		WithCredentials: true,
	}))

	app.Get("/roles", userHandler.Roles)
	app.Get("/search-user", userHandler.Search)
	app.Get("/meili-health", userHandler.CheckHealthMeilisearch)
	app.Get("/units", userHandler.GetUnits)
	app.Post("/update/:id", userHandler.Update)
}
