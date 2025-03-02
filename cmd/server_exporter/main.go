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
	"net/http"

	_ "net/http/pprof"

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

			startProfilingServer()

			app.Post("/transaksi", exportHandler.ExportRekapTransaksi)
			app.Post("/transaksi-all", exportHandler.ExportAllRekapTransaksi)
			app.Post("/pelanggan", exportHandler.ExportRekapPelanggan)
			app.Get("/hello", exportHandler.HelloWorld)

			// listRoutes(app)
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

func startProfilingServer() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
