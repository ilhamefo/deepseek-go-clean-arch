package main

import (
	"context"
	"event-registration/internal/config"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/infrastructure/validator"
	"event-registration/internal/middleware"
	"event-registration/internal/repository/gorm"
	"event-registration/internal/route"

	_ "event-registration/docs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title Auth API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:5051
// @BasePath /
func main() {
	app := fx.New(
		// Provide dependencies
		fx.Provide(
			config.Load, // config.Load should be the first to ensure config is available for other components
			config.NewLogLevel,
			config.NewZapLogger,
			config.NewZapGormLogger,
			config.NewRedisConfig,
			middleware.NewMiddleware,
			validator.NewValidator,
			config.NewGoogleOAuthConfig,
			database.NewGormDBAuth,
			gorm.NewAuthRepo,
			service.NewAuthService,
			handler.NewAuthHandler,
			config.NewFiberApp,
		),

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, logger *zap.Logger) {

			app.Use(config.NewZapLoggerMiddleware(logger))

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						logger.Info("server_started", zap.String("server_on", cfg.ServerAddress+":"+cfg.ServerPort))
						if err := app.Listen(cfg.ServerAddress + ":" + cfg.ServerPort); err != nil {
							logger.Error("error_listening_to_server", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("server_stoped")
					return app.Shutdown()
				},
			})
		}),

		fx.Invoke(route.RegisterAuthRoutes),
	)

	app.Run()
}
