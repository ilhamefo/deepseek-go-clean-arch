package main

import (
	"context"
	_ "event-registration/cmd/server_auth/docs"
	"event-registration/internal/common"
	"event-registration/internal/config"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/infrastructure/validator"
	"event-registration/internal/middleware"
	"event-registration/internal/repository/gorm"
	"event-registration/internal/route"
	_ "net/http/pprof"

	// swag init --generalInfo cmd/server_auth/main.go -o cmd/server_auth/docs --tags Auth

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
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
// @host 127.0.0.1:5052
// @BasePath /
func main() {

	defer tracer.Stop()
	app := fx.New(
		// Provide dependencies
		fx.Provide(
			common.Load, // config.Load should be the first to ensure config is available for other components
			config.NewLogLevel,
			config.NewZapLogger,
			config.NewZapGormLogger,
			config.NewRedisConfig,
			config.NewRedisCache,
			service.NewSessionService,
			middleware.NewMiddleware,
			validator.NewValidator,
			common.NewHandler,
			config.NewHTTPClient,
			config.NewGoogleOAuthConfig,

			// handlers, repositories, services, and routes
			fx.Annotate(database.NewGormDBAuth, fx.ResultTags(`name:"AuthDB"`)),
			fx.Annotate(gorm.NewAuthRepo, fx.ParamTags(`name:"AuthDB"`)),
			handler.NewAuthHandler,
			service.NewAuthService,
			config.NewFiberApp,
		),

		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, cfg *common.Config, logger *zap.Logger, m *middleware.Middleware) {

			// app.Use(m.SentryMiddleware(sentryOpts))
			// app.Use(m.NewZapLoggerMiddleware(logger))

			// tracer.Start(
			// 	tracer.WithService(cfg.DDService),
			// 	tracer.WithEnv(cfg.DDENV),
			// 	tracer.WithServiceVersion(cfg.DDVersion),
			// 	tracer.WithAgentAddr("localhost:8126"), // jika perlu override
			// )
			startProfilingServer()

			app.Use(
				m.PrometheusMiddleware(), // prometheus middleware
			)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						logger.Info("server_started", zap.String("server_on", cfg.ServerAuthAddress+":"+cfg.ServerAuthPort))
						if err := app.Listen(cfg.ServerAuthAddress + ":" + cfg.ServerAuthPort); err != nil {
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

func startProfilingServer() {
	// go func() {
	// 	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	// }()
}
