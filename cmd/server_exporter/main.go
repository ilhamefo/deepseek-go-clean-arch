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
			app.Post("/transaksi", exportHandler.ExportRekapTransaksi)
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
							log.Fatalf("Failed to start server: %v", err)

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
