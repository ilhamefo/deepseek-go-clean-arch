package config

import (
	"event-registration/internal/common"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleOAuthConfig(cfg *common.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  cfg.GoogleRedirectUri,
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       strings.Split(cfg.GoogleOAuthScope, ","),
		Endpoint:     google.Endpoint,
	}
}
