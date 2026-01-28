package route

import (
	"event-registration/internal/handler"
	"event-registration/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterGarminRoutes(app *fiber.App, h *handler.GarminHandler, m *middleware.Middleware) {
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:        true,
		DocExpansion:       "list",
		WithCredentials:    false,
		DisplayOperationId: true,
	}))

	// app.Use(m.HTTPTimeoutMiddleware())

	// metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Get("/health-check", h.HealthCheck)
	app.Get("/activity/:activityID", h.GetActivity)

	app.Post("/refresh", h.Refresh)
	app.Post("/activity-types", h.GetActivityTypes)
	app.Post("/user-profile", h.GetUserProfile)
	app.Post("/heart-rate-by-date", h.GetHeartRateByDate)
	app.Post("/step-by-date", h.GetStepByDate)
	app.Post("/hrv-by-date", h.HRVByDate)
	app.Post("/body-battery-by-date", h.GetBodyBatteryByDate)
	app.Post("/splits/:activityID", h.Splits)
	app.Post("/sleep-by-date", h.GetSleepByDate)
}
