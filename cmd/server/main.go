package main

import (
	"event-registration/internal/config"
	"event-registration/internal/core/service"
	"event-registration/internal/handler"
	"event-registration/internal/infrastructure/database"
	"event-registration/internal/infrastructure/queue"
	gormRepo "event-registration/internal/repository/gorm"
	redisRepo "event-registration/internal/repository/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Application started")

	// Set up database
	db, err := database.NewGormDB(cfg.PostgresURL)
	if err != nil {
		panic(err)
	}

	// Set up Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:       cfg.RedisURL,
		PoolSize:   100,
		ClientName: "Event Registration",
	})

	// Set up RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Set up repositories
	eventRepo := gormRepo.NewEventRepo(db)
	cacheRepo := redisRepo.NewCacheRepo(rdb)
	eventQueue := queue.NewEventQueue(conn)

	// Set up services
	eventService := service.NewEventService(eventRepo, cacheRepo, eventQueue)

	// Set up Fiber app
	app := fiber.New()

	// Set up HTTP handler
	eventHandler := handler.NewEventHandler(eventService)
	app.Post("/register", eventHandler.RegisterEvent)

	// Start server
	app.Listen(":8080")
}
