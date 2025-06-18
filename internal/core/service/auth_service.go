package service

import (
	"context"
	"event-registration/internal/core/domain"
	"event-registration/internal/helper"
	"io"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type AuthService struct {
	repo              domain.AuthRepository
	logger            *zap.Logger
	googleOauthConfig *oauth2.Config
}

func NewAuthService(repo domain.AuthRepository, logger *zap.Logger, googleConfig *oauth2.Config) *AuthService {
	return &AuthService{
		repo:              repo,
		googleOauthConfig: googleConfig,
		logger:            logger,
	}
}

func (s *AuthService) GetLoginUrl() (url string) {
	return s.googleOauthConfig.AuthCodeURL("random-state-token", oauth2.AccessTypeOffline)
}

func (s *AuthService) GoogleHandleCallback(ctx context.Context, code string) (err error) {
	token, err := s.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return err
	}

	helper.PrettyPrint(token, "TOKENNNNNNNNN ===================")

	client := s.googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return err
	}

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if userInfo != nil {
		helper.PrettyPrint(string(userInfo), "userInfo ===================")
	}

	defer resp.Body.Close()

	return err
}
