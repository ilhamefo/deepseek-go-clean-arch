package middleware

import "event-registration/internal/common"

type Middleware struct {
	cfg     *common.Config
	handler *common.Handler
}

func NewMiddleware(config *common.Config, handler *common.Handler) *Middleware {
	return &Middleware{
		cfg:     config,
		handler: handler,
	}
}
