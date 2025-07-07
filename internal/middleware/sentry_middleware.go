package middleware

import (
	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) SentryMiddleware(opts sentry.ClientOptions) fiber.Handler {

	if err := sentry.Init(opts); err != nil {
		panic("Sentry initialization failed: " + err.Error())
	}

	return sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	})
}
