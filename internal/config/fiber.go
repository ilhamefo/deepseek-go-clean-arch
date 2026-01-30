package config

import (
	"event-registration/internal/middleware"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis"
)

func NewFiberApp(m *middleware.Middleware, redisStore *redis.Storage) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "Auth API",
		ErrorHandler:  m.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:3000",
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
	app.Use(recover.New())

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
