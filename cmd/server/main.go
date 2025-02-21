package main

import (
	"context"
	"event-registration/internal/config"
	"event-registration/internal/core/domain"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/infrastructure/queue"
	"event-registration/internal/repository/gorm"
	"event-registration/internal/repository/redis"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Provide dependencies
		fx.Provide(
			config.Load,
			database.NewGormDB,
			service.NewEventService,
			handler.NewEventHandler,
			fiber.New,
		),

		// Bind interfaces to implementations
		fx.Provide(
			fx.Annotate(
				gorm.NewEventRepo,
				fx.As(new(domain.EventRepository)),
			),
			fx.Annotate(
				redis.NewCacheRepo,
				fx.As(new(domain.EventCache)),
			),
			fx.Annotate(
				queue.NewEventQueue,
				fx.As(new(domain.EventQueue)),
			),
		),

		// Invoke the application
		fx.Invoke(func(app *fiber.App, eventHandler *handler.EventHandler) {
			app.Post("/register", eventHandler.RegisterEvent)
		}),

		// Lifecycle hooks
		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						log.Println("Starting server on :8080")
						if err := app.Listen(":8080"); err != nil {
							log.Fatalf("Failed to start server: %v", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Println("Stopping server")
					return app.Shutdown()
				},
			})
		}),
	)

	// Run the application
	app.Run()
}
