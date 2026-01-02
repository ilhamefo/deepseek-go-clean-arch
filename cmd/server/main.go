package main

import (
	"context"
	"event-registration/internal/common"
	"event-registration/internal/config"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/infrastructure/validator"
	"event-registration/internal/middleware"
	"event-registration/internal/repository/gorm"
	"event-registration/internal/route"

	_ "event-registration/docs"

	"github.com/getsentry/sentry-go"
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
			common.Load, // config.Load should be the first to ensure config is available for other components
			config.NewLogLevel,
			config.NewZapLogger,
			config.NewSentryOptions,
			config.NewZapGormLogger,
			config.NewRedisConfig,
			config.NewMeilisearchClient,
			config.NewRedisCache,
			service.NewSessionService,
			middleware.NewMiddleware,
			validator.NewValidator,
			common.NewHandler,
			config.NewGoogleOAuthConfig,
			fx.Annotate(database.NewGormDBAuth, fx.ResultTags(`name:"authDB"`)),
			fx.Annotate(database.NewGormDBVCC, fx.ResultTags(`name:"VCCDB"`)),
			fx.Annotate(gorm.NewAuthRepo, fx.ParamTags(`name:"authDB"`)),
			fx.Annotate(gorm.NewUserRepo, fx.ParamTags(`name:"VCCDB"`)),
			service.NewUserService,
			handler.NewUserHandler,
			service.NewAuthService,
			handler.NewAuthHandler,
			config.NewFiberApp,
		),

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, cfg *common.Config, logger *zap.Logger, m *middleware.Middleware, sentryOpts sentry.ClientOptions) {

			// app.Use(m.SentryMiddleware(sentryOpts))
			app.Use(m.NewZapLoggerMiddleware(logger))

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

		fx.Invoke(route.RegisterUserRoutes),
		// fx.Invoke(route.RegisterAuthRoutes),
	)

	app.Run()
}
