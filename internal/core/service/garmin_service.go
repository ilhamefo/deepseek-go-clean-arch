package service

import (
	"context"
	"encoding/json"
	"event-registration/internal/common"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type GarminService struct {
	repo   domain.GarminRepository
	logger *zap.Logger
	config *common.Config
}

func NewGarminService(repo domain.GarminRepository, logger *zap.Logger, config *common.Config) *GarminService {
	return &GarminService{
		repo:   repo,
		logger: logger,
		config: config,
	}
}

func (s *GarminService) Refresh(ctx context.Context, r *request.RefreshActivitiesRequest) (res []*domain.Activity, err error) {
	url := "https://connect.garmin.com/activitylist-service/activities/search/activities?limit=10&start=0"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("error_make_new_request", zap.Error(err))
		return nil, err
	}

	// Header
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("authorization", "Bearer "+r.Token)
	req.Header.Set("di-backend", "connectapi.garmin.com")

	// Cookie
	req.Header.Set("Cookie", r.Cookies)

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("error_do_request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Baca response
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			s.logger.Error("error_read_request", zap.Error(err))
			return nil, err
		}

		s.logger.Error("error_response",
			zap.Any("status_code", resp.StatusCode),
			zap.Any("body", string(errBody)),
		)

		return nil, err
	}

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("error_read_request", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		s.logger.Error("error_unmarshal_json", zap.Error(err))
		return nil, err
	}

	err = s.Upsert(res)
	if err != nil {
		s.logger.Error("error_upsert", zap.Error(err))
		return nil, err
	}

	return res, err
}

func (s *GarminService) Upsert(models []*domain.Activity) (err error) {

	err = s.repo.Update(models)
	if err != nil {
		s.logger.Error("error_update_data", zap.Error(err))
		return err
	}

	return err
}
