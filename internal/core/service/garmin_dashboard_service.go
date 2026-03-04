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

func (s *GarminDashboardService) GetActivityDetails(ctx context.Context, payload *request.ActivityDetailsDashboardRequest) (res domain.ActivityDetailsResponse, err error) {
	metrics, activitySmmary, err := s.repo.GetActivityDetails(ctx, payload.ActivityID)
	if err != nil {
		s.logger.Error("error_get_activities", zap.Error(err))
		return res, err
	}

	res = s.transformToCompactFormat(payload.ActivityID, metrics, activitySmmary)

	return res, nil
}

func (s *GarminDashboardService) transformToCompactFormat(activityID int64, metrics []*domain.ActivityMetricsTimeseries, activitySmmary *domain.ActivityVM) domain.ActivityDetailsResponse {
	descriptors := []domain.MetricDescriptorDTO{
		{MetricsIndex: 0, Key: "directDoubleCadence", Unit: domain.UnitDTO{ID: 92, Key: "stepsPerMinute", Factor: 1.0}},
		{MetricsIndex: 1, Key: "directAirTemperature", Unit: domain.UnitDTO{ID: 227, Key: "celcius", Factor: 1.0}},
		{MetricsIndex: 2, Key: "directFractionalCadence", Unit: domain.UnitDTO{ID: 92, Key: "stepsPerMinute", Factor: 1.0}},
		{MetricsIndex: 3, Key: "directSpeed", Unit: domain.UnitDTO{ID: 20, Key: "mps", Factor: 0.1}},
		{MetricsIndex: 4, Key: "sumMovingDuration", Unit: domain.UnitDTO{ID: 40, Key: "second", Factor: 1000.0}},
		{MetricsIndex: 5, Key: "sumDuration", Unit: domain.UnitDTO{ID: 40, Key: "second", Factor: 1000.0}},
		{MetricsIndex: 6, Key: "directPower", Unit: domain.UnitDTO{ID: 10, Key: "watt", Factor: 1.0}},
		{MetricsIndex: 7, Key: "directLongitude", Unit: domain.UnitDTO{ID: 60, Key: "dd", Factor: 1.0}},
		{MetricsIndex: 8, Key: "sumAccumulatedPower", Unit: domain.UnitDTO{ID: 10, Key: "watt", Factor: 1.0}},
		{MetricsIndex: 9, Key: "directTimestamp", Unit: domain.UnitDTO{ID: 120, Key: "gmt", Factor: 0.0}},
		{MetricsIndex: 10, Key: "directRunCadence", Unit: domain.UnitDTO{ID: 92, Key: "stepsPerMinute", Factor: 1.0}},
		{MetricsIndex: 11, Key: "sumDistance", Unit: domain.UnitDTO{ID: 1, Key: "meter", Factor: 100.0}},
		{MetricsIndex: 12, Key: "directVerticalOscillation", Unit: domain.UnitDTO{ID: 5, Key: "centimeter", Factor: 1.0}},
		{MetricsIndex: 13, Key: "directGradeAdjustedSpeed", Unit: domain.UnitDTO{ID: 20, Key: "mps", Factor: 0.1}},
		{MetricsIndex: 14, Key: "directBodyBattery", Unit: domain.UnitDTO{ID: 6, Key: "dimensionless", Factor: 1.0}},
		{MetricsIndex: 15, Key: "directHeartRate", Unit: domain.UnitDTO{ID: 100, Key: "bpm", Factor: 1.0}},
		{MetricsIndex: 16, Key: "sumElapsedDuration", Unit: domain.UnitDTO{ID: 40, Key: "second", Factor: 1000.0}},
		{MetricsIndex: 17, Key: "directLatitude", Unit: domain.UnitDTO{ID: 60, Key: "dd", Factor: 1.0}},
		{MetricsIndex: 18, Key: "directElevation", Unit: domain.UnitDTO{ID: 1, Key: "meter", Factor: 100.0}},
		{MetricsIndex: 19, Key: "directStrideLength", Unit: domain.UnitDTO{ID: 5, Key: "centimeter", Factor: 1.0}},
		{MetricsIndex: 20, Key: "directVerticalRatio", Unit: domain.UnitDTO{ID: 6, Key: "dimensionless", Factor: 1.0}},
		{MetricsIndex: 21, Key: "directVerticalSpeed", Unit: domain.UnitDTO{ID: 20, Key: "mps", Factor: 0.1}},
		{MetricsIndex: 22, Key: "directGroundContactTime", Unit: domain.UnitDTO{ID: 41, Key: "ms", Factor: 1.0}},
	}

	// Transform metrics to array format
	detailMetrics := make([]domain.Metrics, 0, len(metrics))
	for _, m := range metrics {
		metricValues := make([]float64, 23)

		// Map each field to its index position (use 0 for nil values)
		if m.DirectDoubleCadence != nil {
			metricValues[0] = float64(*m.DirectDoubleCadence)
		}
		if m.DirectAirTemperature != nil {
			metricValues[1] = *m.DirectAirTemperature
		}
		if m.DirectFractionalCadence != nil {
			metricValues[2] = *m.DirectFractionalCadence
		}
		if m.DirectSpeed != nil {
			metricValues[3] = *m.DirectSpeed
		}
		if m.SumMovingDuration != nil {
			metricValues[4] = *m.SumMovingDuration
		}
		if m.SumDuration != nil {
			metricValues[5] = *m.SumDuration
		}
		if m.DirectPower != nil {
			metricValues[6] = *m.DirectPower
		}
		if m.DirectLongitude != nil {
			metricValues[7] = *m.DirectLongitude
		}
		if m.SumAccumulatedPower != nil {
			metricValues[8] = *m.SumAccumulatedPower
		}
		if m.DirectTimestamp != nil {
			metricValues[9] = float64(*m.DirectTimestamp)
		}
		if m.DirectRunCadence != nil {
			metricValues[10] = float64(*m.DirectRunCadence)
		}
		if m.SumDistance != nil {
			metricValues[11] = *m.SumDistance
		}
		if m.DirectVerticalOscillation != nil {
			metricValues[12] = *m.DirectVerticalOscillation
		}
		if m.DirectGradeAdjustedSpeed != nil {
			metricValues[13] = *m.DirectGradeAdjustedSpeed
		}
		if m.DirectBodyBattery != nil {
			metricValues[14] = float64(*m.DirectBodyBattery)
		}
		if m.DirectHeartRate != nil {
			metricValues[15] = float64(*m.DirectHeartRate)
		}
		if m.SumElapsedDuration != nil {
			metricValues[16] = *m.SumElapsedDuration
		}
		if m.DirectLatitude != nil {
			metricValues[17] = *m.DirectLatitude
		}
		if m.DirectElevation != nil {
			metricValues[18] = *m.DirectElevation
		}
		if m.DirectStrideLength != nil {
			metricValues[19] = *m.DirectStrideLength
		}
		if m.DirectVerticalRatio != nil {
			metricValues[20] = *m.DirectVerticalRatio
		}
		if m.DirectVerticalSpeed != nil {
			metricValues[21] = *m.DirectVerticalSpeed
		}
		if m.DirectGroundContactTime != nil {
			metricValues[22] = *m.DirectGroundContactTime
		}

		detailMetrics = append(detailMetrics, domain.Metrics{
			Metrics: metricValues,
		})
	}

	return domain.ActivityDetailsResponse{
		ActivityID:            activityID,
		MeasurementCount:      23,
		MetricsCount:          len(metrics),
		TotalMetricsCount:     len(metrics),
		MetricDescriptors:     descriptors,
		ActivityDetailMetrics: detailMetrics,
		GeoPolylineDTO:        domain.GeoPolylineDTO{Polyline: []interface{}{}},
		HeartRateDTOs:         nil,
		PendingData:           nil,
		DetailsAvailable:      len(metrics) > 0,
		ActivitySummary:       activitySmmary,
	}
}
