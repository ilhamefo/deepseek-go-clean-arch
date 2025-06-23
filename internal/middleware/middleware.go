package middleware

import "event-registration/internal/common"

type Middleware struct {
	cfg *common.Config
}

func NewMiddleware(config *common.Config) *Middleware {
	return &Middleware{
		cfg: config,
	}
}
