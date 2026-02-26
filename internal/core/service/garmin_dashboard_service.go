package service

import (
	"context"
	"event-registration/internal/common"
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
