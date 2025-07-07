package config

import (
	"event-registration/internal/common"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/zap"
)

func NewSentryOptions(logger *zap.Logger, cfg *common.Config) sentry.ClientOptions {
	return sentry.ClientOptions{
		Dsn: cfg.SentryDSN,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
					logger.Info("sentry_hostname", zap.String("hostname", utils.CopyString(c.Hostname())))
				}
			}

			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	}
}
