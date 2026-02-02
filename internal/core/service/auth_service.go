// auth_servive
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"event-registration/internal/common"
	"event-registration/internal/common/constant"
	"event-registration/internal/common/helper"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

const (
	PASSWORD_LENGTH     = 12
	PASSWORD_MIN_LENGTH = 8
	// Argon2 parameters
	ARGON2_TIME    = 3
	ARGON2_MEMORY  = 64 * 1024 // 64 MB
	ARGON2_THREADS = 4
	ARGON2_SALT    = 16
)

var GoogleUserinfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type AuthService struct {
	repo              domain.AuthRepository
	logger            *zap.Logger
	GoogleOauthConfig *oauth2.Config
	config            *common.Config
	sessionService    *SessionService
}

func NewAuthService(repo domain.AuthRepository, logger *zap.Logger, googleConfig *oauth2.Config, config *common.Config, sessionService *SessionService) *AuthService {
	return &AuthService{
		repo:              repo,
		GoogleOauthConfig: googleConfig,
		logger:            logger,
		config:            config,
		sessionService:    sessionService,
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

	return s.GoogleOauthConfig.AuthCodeURL(token, oauth2.AccessTypeOffline), token, nil
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
	var user *domain.User

	// check state
	if req.State != req.StateCookie {
		s.logger.Error(
			"error_invalid_state",
			zap.Error(err),
		)

		return accessToken, refreshToken, errors.New("error_invalid_state")
	}

	token, err := s.GoogleOauthConfig.Exchange(ctx, req.Code)
	if err != nil {
		s.logger.Error(
			"error_exchange_token",
			zap.Error(err),
		)
		return accessToken, refreshToken, err
	}

	client := s.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get(GoogleUserinfoURL)
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

	if userInfo == nil {
		return accessToken, refreshToken, errors.New("error_get_user_info")
	}

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
		user.ID = helper.GenerateUUID()
		err = s.Register(*user)
		if err != nil {
			s.logger.Error(
				"error_registered",
				zap.Error(err),
			)
			return accessToken, refreshToken, err
		}
	} else {
		user, err = s.repo.FindByEmail(user.Email)
		if err != nil {
			s.logger.Error(
				"error_get_user_by_email",
				zap.Error(err),
			)
			return accessToken, refreshToken, err
		}
	}

	accessToken, refreshToken, err = s.GenerateToken(user)
	if err != nil {
		s.logger.Error("error_create_token", zap.Error(err))
		return accessToken, refreshToken, errors.New("invalid_credentials")
	}

	return accessToken, refreshToken, err
}

func (s *AuthService) Register(user domain.User) (err error) {
	var now time.Time = time.Now()

	if user.VerifiedEmail {
		user.EmailVerifiedAt = &now
	}

	password, err := s.GenerateSafePasswordV2(PASSWORD_LENGTH)
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

func (s *AuthService) GenerateSafePassword(length int) (string, error) {
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

func (s *AuthService) GenerateSafePasswordV2(length int) (string, error) {
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

	helper.PrettyPrint(
		rawPassword,
		"RAW PASSWORD ===================",
	)

	salt := make([]byte, ARGON2_SALT)
	_, err = rand.Read(salt)
	if err != nil {
		s.logger.Error("error_generate_salt", zap.Error(err))
		return "", err
	}

	hashedPassword := argon2.IDKey(
		[]byte(rawPassword),
		salt,
		ARGON2_TIME,
		ARGON2_MEMORY,
		ARGON2_THREADS,
		32,
	)

	hashedWithSalt := append(salt, hashedPassword...)
	encodedHash := base64.RawURLEncoding.EncodeToString(hashedWithSalt)

	return encodedHash, nil
}

func (s *AuthService) VerifyArgon2Password(hashedPassword string, rawPassword string) (bool, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(hashedPassword)
	if err != nil {
		s.logger.Error("error_decode_hash", zap.Error(err))
		return false, err
	}

	salt := decoded[:ARGON2_SALT]
	storedHash := decoded[ARGON2_SALT:]

	computedHash := argon2.IDKey(
		[]byte(rawPassword),
		salt,
		ARGON2_TIME,
		ARGON2_MEMORY,
		ARGON2_THREADS,
		32,
	)

	return string(storedHash) == string(computedHash), nil
}

func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest) (accessToken, refreshToken string, err error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		s.logger.Error("error_get_user_by_email", zap.Error(err))
		return accessToken, refreshToken, errors.New(constant.ACCESS_TOKEN)
	}

	isValid, err := s.VerifyArgon2Password(user.Password, req.Password)
	if err != nil || !isValid {
		s.logger.Error("error_verify_password", zap.Error(err))
		return accessToken, refreshToken, errors.New(constant.ACCESS_TOKEN)
	}

	accessToken, refreshToken, err = s.GenerateToken(user)
	if err != nil {
		s.logger.Error("error_create_token", zap.Error(err))
		return accessToken, refreshToken, errors.New(constant.ACCESS_TOKEN)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GenerateToken(user *domain.User) (accessToken, refreshToken string, err error) {
	accessToken, err = s.GenerateAccessTokenJWT(user)
	if err != nil {
		s.logger.Error("error_generate_access_token_jwt", zap.Error(err))
		return accessToken, refreshToken, err
	}

	refreshToken, err = s.GenerateRefreshTokenJWT(user)
	if err != nil {
		s.logger.Error("error_generate_refresh_token_jwt", zap.Error(err))
		return accessToken, refreshToken, err
	}

	expiration := time.Duration(s.config.RefreshTokenExpiration) * 24 * time.Hour
	err = s.sessionService.CreateSession(context.Background(), user.ID, user.Email, refreshToken, expiration)
	if err != nil {
		s.logger.Error("error_create_session", zap.Error(err))
		return accessToken, refreshToken, err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GenerateAccessTokenJWT(user *domain.User) (string, error) {
	claims := map[string]any{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(s.config.AccessJwtExpiration) * time.Minute).Unix(),
		"type":  constant.ACCESS_TOKEN,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(s.config.JwtSecret))
}

func (s *AuthService) GenerateRefreshTokenJWT(user *domain.User) (string, error) {

	claims := map[string]any{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(s.config.RefreshTokenExpiration) * 24 * time.Hour).Unix(),
		"type":  constant.REFRESH_TOKEN,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(s.config.JwtSecret))
}

func (s *AuthService) Logout(ctx context.Context, refreshToken, accessToken string) error {
	err := s.sessionService.BlacklistAccessToken(ctx, accessToken, time.Now().Add(time.Duration(s.config.AccessJwtExpiration)*time.Minute))
	if err != nil {
		s.logger.Error("error_blacklist_access_token", zap.Error(err))
		return err
	}

	err = s.sessionService.DeleteSession(ctx, refreshToken)
	if err != nil {
		s.logger.Error("error_delete_session", zap.Error(err))
		return err
	}

	return nil
}

func (s *AuthService) LogoutAllDevices(ctx context.Context, userID string) error {
	return s.sessionService.DeleteAllUserSessions(ctx, userID)
}
