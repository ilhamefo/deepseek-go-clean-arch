// auth_servive
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"event-registration/internal/common"
	"event-registration/internal/core/domain"
	"event-registration/internal/request"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	config            *common.Config
}

func NewAuthService(repo domain.AuthRepository, logger *zap.Logger, googleConfig *oauth2.Config, config *common.Config) *AuthService {
	return &AuthService{
		repo:              repo,
		googleOauthConfig: googleConfig,
		logger:            logger,
		config:            config,
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

func (s *AuthService) GoogleHandleCallback(ctx context.Context, req *request.GoogleCallbackRequest) (accessToken, refreshToken string, err error) {
	var user domain.User

	// check state
	// if req.State != req.StateCookie {
	// 	s.logger.Error(
	// 		"error_invalid_state",
	// 		zap.Error(err),
	// 	)

	// 	return accessToken, refreshToken, errors.New("error_invalid_state")
	// }

	token, err := s.googleOauthConfig.Exchange(ctx, req.Code)
	if err != nil {
		s.logger.Error(
			"error_exchange_token",
			zap.Error(err),
		)
		return accessToken, refreshToken, err
	}

	client := s.googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.logger.Error(
			"error_get_client",
			zap.Error(err),
		)
		return accessToken, refreshToken, err
	}
	defer resp.Body.Close()

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error(
			"error_read_body",
			zap.Error(err),
		)
		return accessToken, refreshToken, err
	}

	if userInfo != nil {
		if err := json.Unmarshal(userInfo, &user); err != nil {
			s.logger.Error("error_unmarshal_user_info", zap.Error(err))
			return accessToken, refreshToken, err
		}

		exists, err := s.repo.IsRegistered(user.Email)
		if err != nil {
			s.logger.Error(
				"error_check_is_registered",
				zap.Error(err),
			)
			return accessToken, refreshToken, err
		}

		s.logger.Info("check_is_registered", zap.Any("exists", exists))

		if !exists {
			user.ID = ""
			err = s.register(user)
			if err != nil {
				s.logger.Error(
					"error_registered",
					zap.Error(err),
				)
				return accessToken, refreshToken, err
			}
		}

		accessToken, refreshToken, err = s.generateToken(&user)
		if err != nil {
			s.logger.Error("error_create_token", zap.Error(err))
			return accessToken, refreshToken, errors.New("invalid_credentials")
		}

		return accessToken, refreshToken, err
	}

	return accessToken, refreshToken, err
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

func (s *AuthService) login(ctx context.Context, req *request.LoginRequest) (accessToken, refreshToken string, err error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		s.logger.Error("error_get_user_by_email", zap.Error(err))
		return accessToken, refreshToken, errors.New("invalid_credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Error("error_compare_password", zap.Error(err))
		return accessToken, refreshToken, errors.New("invalid_credentials")
	}

	accessToken, refreshToken, err = s.generateToken(user)
	if err != nil {
		s.logger.Error("error_create_token", zap.Error(err))
		return accessToken, refreshToken, errors.New("invalid_credentials")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) generateToken(user *domain.User) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateAccessTokenJWT(user)
	if err != nil {
		s.logger.Error("error_generate_access_token_jwt", zap.Error(err))
		return accessToken, refreshToken, err
	}

	refreshToken, err = s.generateRefreshTokenJWT(user)
	if err != nil {
		s.logger.Error("error_generate_refresh_token_jwt", zap.Error(err))
		return accessToken, refreshToken, err
	}

	return accessToken, refreshToken, err
}

func (s *AuthService) generateAccessTokenJWT(user *domain.User) (string, error) {
	claims := map[string]any{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(s.config.AccessJwtExpiration) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(s.config.JwtSecret))
}

func (s *AuthService) generateRefreshTokenJWT(user *domain.User) (string, error) {

	claims := map[string]interface{}{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(s.config.RefreshTokenExpiration) * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(s.config.JwtSecret))
}
