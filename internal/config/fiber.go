package config

import (
	"event-registration/internal/middleware"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis"
	"go.uber.org/zap"
)

func NewFiberApp(m *middleware.Middleware, redisStore *redis.Storage, logger *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "Auth API",
		ErrorHandler:  m.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:3000,http://127.0.0.1:4000,http://172.16.1.59:4000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
		AllowMethods:     "GET,POST",
	}))
	app.Use(compress.New())
	app.Use(helmet.New())
	// app.Use(limiter.New(limiter.Config{
	// 	Max:          100,
	// 	Expiration:   1 * time.Minute,
	// 	LimitReached: limitReachedResponse(),
	// 	Store:        redisStore,
	// }))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.Error("panic_recovered",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
				zap.Any("panic", e),
				zap.String("stack_trace", string(debug.Stack())),
			)
		},
	}))

	return app
}

func limitReachedResponse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Too many requests, please try again later.",
			"error":   "rate_limit_exceeded",
		})
	}
}
