package main

import (
	"context"
	"event-registration/internal/common"
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
	"github.com/gofiber/swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "event-registration/exporter_docs"
	// swag init --generalInfo cmd/server_exporter/main.go -o exporter_docs --tags exporter
	// go build -o main.exe cmd/server_exporter/main.go
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:5050
// @BasePath /
func main() {
	app := fx.New(

		fx.Provide(
			common.Load,
			config.NewLogLevel,
			config.NewZapLogger,
			config.NewZapGormLogger,
			validator.NewValidator,
			fx.Annotate(database.NewGormDwhDB, fx.ResultTags(`name:"DwhDB"`)),
			fx.Annotate(database.NewGormPlnMobileDB, fx.ResultTags(`name:"PlnMobileDB"`)),
			redis.NewCacheRepo,
			fx.Annotate(gorm.NewExporterRepo, fx.ParamTags(`name:"DwhDB"`, `name:"PlnMobileDB"`)),
			service.NewExporterService,
			handler.NewExporterHandler,
			fiber.New,
		),

		fx.Invoke(func(app *fiber.App, exportHandler *handler.ExporterHandler) {

			// Register Swagger route
			app.Get("/swagger/*", swagger.New(swagger.Config{
				DeepLinking:  true,
				DocExpansion: "list",
			}))

			startProfilingServer()

			// Routes
			app.Post("/transaksi", exportHandler.ExportRekapTransaksi)
			app.Post("/transaksi-all", exportHandler.ExportAllRekapTransaksi)
			app.Post("/pelanggan", exportHandler.ExportRekapPelanggan)
			app.Get("/hello", exportHandler.HelloWorld)

			// listRoutes(app)
		}),

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, config *common.Config, logger *zap.Logger) {
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
