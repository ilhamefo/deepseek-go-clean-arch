package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SessionData struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	LoginAt   time.Time `json:"login_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionService struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewSessionService(redis *redis.Client, logger *zap.Logger) *SessionService {
	return &SessionService{
		redis:  redis,
		logger: logger,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, userID, email string, refreshToken string, expiration time.Duration) error {
	sessionData := SessionData{
		UserID:    userID,
		Email:     email,
		LoginAt:   time.Now(),
		ExpiresAt: time.Now().Add(expiration),
	}

	data, err := json.Marshal(sessionData)
	if err != nil {
		s.logger.Error("failed to marshal session data", zap.Error(err))
		return err
	}

	// Store with refresh token as key
	key := fmt.Sprintf("session:%s", refreshToken)
	err = s.redis.Set(ctx, key, data, expiration).Err()
	if err != nil {
		s.logger.Error("failed to store session in redis", zap.Error(err))
		return err
	}

	// Also store user active sessions (for multiple device management)
	userSessionKey := fmt.Sprintf("user_sessions:%s", userID)
	s.redis.SAdd(ctx, userSessionKey, refreshToken)
	s.redis.Expire(ctx, userSessionKey, expiration)

	return nil
}

func (s *SessionService) GetSession(ctx context.Context, refreshToken string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", refreshToken)
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		s.logger.Error("failed to get session from redis", zap.Error(err))
		return nil, err
	}

	var sessionData SessionData
	err = json.Unmarshal([]byte(data), &sessionData)
	if err != nil {
		s.logger.Error("failed to unmarshal session data", zap.Error(err))
		return nil, err
	}

	return &sessionData, nil
}

func (s *SessionService) CheckAccessToken(ctx context.Context, access string) (exist bool, err error) {
	key := fmt.Sprintf("session:%s", access)
	_, err = s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		s.logger.Error("failed to check access token in redis", zap.Error(err))
		return false, err
	}

	return true, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, refreshToken string) error {
	// Get session data first to remove from user sessions
	sessionData, err := s.GetSession(ctx, refreshToken)
	if err != nil {
		return err
	}

	// Remove from individual session
	key := fmt.Sprintf("session:%s", refreshToken)
	err = s.redis.Del(ctx, key).Err()
	if err != nil {
		s.logger.Error("failed to delete session from redis", zap.Error(err))
		return err
	}

	// Remove from user sessions set
	userSessionKey := fmt.Sprintf("user_sessions:%s", sessionData.UserID)
	s.redis.SRem(ctx, userSessionKey, refreshToken)

	return nil
}

func (s *SessionService) DeleteAllUserSessions(ctx context.Context, userID string) error {
	userSessionKey := fmt.Sprintf("user_sessions:%s", userID)

	// Get all refresh tokens for this user
	refreshTokens, err := s.redis.SMembers(ctx, userSessionKey).Result()
	if err != nil {
		s.logger.Error("failed to get user sessions", zap.Error(err))
		return err
	}

	// Delete each session
	for _, refreshToken := range refreshTokens {
		sessionKey := fmt.Sprintf("session:%s", refreshToken)
		s.redis.Del(ctx, sessionKey)
	}

	// Delete the user sessions set
	s.redis.Del(ctx, userSessionKey)

	return nil
}

func (s *SessionService) IsSessionValid(ctx context.Context, refreshToken string) bool {
	sessionData, err := s.GetSession(ctx, refreshToken)
	if err != nil {
		return false
	}

	return time.Now().Before(sessionData.ExpiresAt)
}

// BlacklistAccessToken menambahkan access token ke blacklist Redis sampai expired
func (s *SessionService) BlacklistAccessToken(ctx context.Context, accessToken string, exp time.Time) error {
	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}

	key := fmt.Sprintf("blacklist:%s", accessToken)
	return s.redis.Set(ctx, key, "1", ttl).Err()
}

// IsAccessTokenBlacklisted mengecek apakah access token ada di blacklist Redis
func (s *SessionService) IsAccessTokenBlacklisted(ctx context.Context, accessToken string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", accessToken)
	exists, err := s.redis.Exists(ctx, key).Result()
	return exists == 1, err
}
