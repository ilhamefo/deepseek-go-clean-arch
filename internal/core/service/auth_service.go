package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"event-registration/internal/core/domain"
	"event-registration/internal/request"
	"io"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

const (
	PASSWORD_LENGTH     = 12
	PASSWORD_MIN_LENGTH = 8
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

func (s *AuthService) GetLoginUrl() (url, token string, err error) {
	token, err = s.generateStateToken()
	if err != nil {
		s.logger.Error(
			"error_generate_state_token",
			zap.Error(err),
		)

		return "", "", err
	}

	return s.googleOauthConfig.AuthCodeURL(token, oauth2.AccessTypeOffline), token, nil
}

func (s *AuthService) generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		s.logger.Error(
			"error_exchange_token",
			zap.Error(err),
		)

		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *AuthService) GoogleHandleCallback(ctx context.Context, req *request.GoogleCallbackRequest) (err error) {
	var user domain.User

	// check state
	if req.State != req.StateCookie {
		s.logger.Error(
			"error_invalid_state",
			zap.Error(err),
		)

		return errors.New("error_invalid_state")
	}

	token, err := s.googleOauthConfig.Exchange(ctx, req.Code)
	if err != nil {
		s.logger.Error(
			"error_exchange_token",
			zap.Error(err),
		)
		return err
	}

	client := s.googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.logger.Error(
			"error_get_client",
			zap.Error(err),
		)
		return err
	}
	defer resp.Body.Close()

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error(
			"error_read_body",
			zap.Error(err),
		)
		return err
	}

	if userInfo != nil {
		if err := json.Unmarshal(userInfo, &user); err != nil {
			s.logger.Error("error_unmarshal_user_info", zap.Error(err))
			return err
		}

		exists, err := s.repo.IsRegistered(user.Email)
		if err != nil {
			s.logger.Error(
				"error_check_is_registered",
				zap.Error(err),
			)
			return err
		}

		s.logger.Info("check_is_registered", zap.Bool("exists", exists))

		if !exists {
			err = s.register(user)
			if err != nil {
				s.logger.Error(
					"error_registered",
					zap.Error(err),
				)
				return err
			}
		}
	}

	return err
}

func (s *AuthService) register(user domain.User) (err error) {
	var now time.Time = time.Now()

	if user.VerifiedEmail {
		user.EmailVerifiedAt = &now
	}

	password, err := s.generateSafePassword(PASSWORD_LENGTH)
	if err != nil {
		s.logger.Error(
			"error_generate_password",
			zap.Error(err),
		)
		return err
	}

	user.Password = password

	return s.repo.Register(user)
}

func (s *AuthService) generateSafePassword(length int) (string, error) {
	if length < PASSWORD_MIN_LENGTH {
		length = PASSWORD_MIN_LENGTH
	}
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		s.logger.Error("error_generate_password", zap.Error(err))
		return "", err
	}

	rawPassword := base64.RawURLEncoding.EncodeToString(bytes)[:length]

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("error_bcrypt_hash", zap.Error(err))
		return "", err
	}

	return string(hashedPassword), nil
}
