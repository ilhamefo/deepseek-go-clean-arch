package middleware

import (
	"event-registration/internal/common"
	"event-registration/internal/core/service"
)

type Middleware struct {
	cfg            *common.Config
	handler        *common.Handler
	sessionService *service.SessionService
}

func NewMiddleware(config *common.Config, handler *common.Handler, sessionService *service.SessionService) *Middleware {
	return &Middleware{
		cfg:            config,
		handler:        handler,
		sessionService: sessionService,
	}
}
