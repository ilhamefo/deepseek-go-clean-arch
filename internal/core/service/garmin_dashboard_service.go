package service

import (
	"context"
	"event-registration/internal/common"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type GarminDashboardService struct {
	repo        domain.GarminDashboardRepository
	logger      *zap.Logger
	config      *common.Config
	httpClient  *http.Client
	redisClient *redis.Client
}

func NewGarminDashboardService(repo domain.GarminDashboardRepository, logger *zap.Logger, config *common.Config, redisClient *redis.Client, httpClient *http.Client) *GarminDashboardService {
	return &GarminDashboardService{
		repo:        repo,
		logger:      logger,
		config:      config,
		httpClient:  httpClient,
		redisClient: redisClient,
	}
}

func (s *GarminDashboardService) HeartRate(ctx context.Context) (res domain.HeartRate, err error) {
	res, err = s.repo.GetHeartRate(ctx, time.Now())
	if err != nil {
		s.logger.Error("error_get_heart_rate", zap.Error(err))
		return res, err
	}

	return res, nil
}

func (s *GarminDashboardService) GetActivities(ctx context.Context, payload *request.ActivityDashboardRequest) (res domain.PaginatedResponse[domain.ActivityVM], err error) {
	if payload.Limit <= 0 {
		payload.Limit = 10 // default limit
	}
	if payload.Limit > 100 {
		payload.Limit = 100 // max limit
	}

	data, nextCursor, hasMore, err := s.repo.GetActivities(ctx, payload)
	if err != nil {
		s.logger.Error("error_get_activities", zap.Error(err))
		return res, err
	}

	res = domain.PaginatedResponse[domain.ActivityVM]{
		Data:       data,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      payload.Limit,
	}

	return res, nil
}
