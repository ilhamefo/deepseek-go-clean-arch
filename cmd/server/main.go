package main

import (
	"context"
	"event-registration/internal/config"
	"event-registration/internal/core/domain"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/repository/gorm"
	"event-registration/internal/repository/redis"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "event-registration/docs"

	"github.com/gofiber/swagger"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5051
// @BasePath /
func main() {
	app := fx.New(
		// Provide dependencies
		fx.Provide(
			config.Load,
			config.NewLogLevel,
			config.NewZapLogger,
			database.NewGormDB,
			service.NewEventService,
			handler.NewEventHandler,
			fiber.New,
		),

		fx.Provide(
			fx.Annotate(
				gorm.NewEventRepo,
				fx.As(new(domain.EventRepository)),
			),
			fx.Annotate(
				redis.NewCacheRepo,
				fx.As(new(domain.EventCache)),
			),
		),

		// Invoke the application
		fx.Invoke(func(app *fiber.App, eventHandler *handler.EventHandler) {
			app.Get("/swagger/*", swagger.New(swagger.Config{
				DeepLinking:  true,
				DocExpansion: "list",
			}))

			// Routes
			app.Get("/users/:id", eventHandler.GetUser)
			// app.Post("/register", eventHandler.RegisterEvent)

			// TODO: buat authentication dengan access token dan refresh token
			// Access token (JWT, short-lived, 15 mins)
			// Refresh token (UUID/random string, stored in Redis, 7 days)
			// /login dan /refresh

			// Swagger UI
			// app.Use(swagger.New(swagger.Config{
			// 	FilePath: "./docs/swagger.json",
			// }))

		}),

		// Lifecycle hooks
		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, config *config.Config, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						logger.Info(
							"server_started",
							zap.String("server_on", config.ServerAddress+":"+config.ServerPort),
						)

						if err := app.Listen(config.ServerAddress + ":" + config.ServerPort); err != nil {
							logger.Error(
								"error_listening_to_server",
								zap.Error(err),
							)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info(
						"server_stoped",
					)
					return app.Shutdown()
				},
			})
		}),
	)

	// Run the application
	app.Run()
}
