package main

import (
	"context"
	"event-registration/internal/config"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/infrastructure/validator"
	"event-registration/internal/repository/gorm"
	"event-registration/internal/repository/redis"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(

		fx.Provide(
			config.Load,
			validator.NewValidator,
			database.NewGormPlnMobileDB,
			redis.NewCacheRepo,
			gorm.NewExporterRepo,
			service.NewExporterService,
			handler.NewExporterHandler,
			config.NewZapLogger,
			fiber.New,
		),

		fx.Invoke(func(app *fiber.App, exportHandler *handler.ExporterHandler) {
			// Initialize default config
			app.Use(pprof.New())

			app.Post("/transaksi", exportHandler.ExportRekapTransaksi)
			app.Get("/hello", exportHandler.HelloWorld)

			listRoutes(app)
		}),

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, config *config.Config, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						logger.Info(
							"server_started",
							zap.String("server_on", config.ServerExporterAddress+":"+config.ServerExporterPort),
						)

						if err := app.Listen(config.ServerExporterAddress + ":" + config.ServerExporterPort); err != nil {
							logger.Error(
								"error_listening_to_server",
								zap.Error(err),
							)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					// close redis connection
					// if err := redis.CloseRedisConnection(client); err != nil {
					// 	logger.Error(
					// 		"error_closing_redis",
					// 		zap.Error(err),
					// 	)
					// }

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

func listRoutes(app *fiber.App) {
	// Get the route stack
	stack := app.Stack()

	// Iterate over the stack
	for _, method := range stack {
		for _, route := range method {
			// Print the HTTP method, path, and handler name
			log.Printf("[%s] %s -> %s", method[0].Method, route.Path, route.Name)
		}
	}
}
