package service

import (
	"context"
	"encoding/json"
	"event-registration/internal/common"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"fmt"
	"io"
	"net/http"
	"time"

	httptrace "github.com/DataDog/dd-trace-go/contrib/net/http/v2"
	"go.uber.org/zap"
)

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type GarminService struct {
	repo       domain.GarminRepository
	logger     *zap.Logger
	config     *common.Config
	httpClient *http.Client
}

func NewGarminService(repo domain.GarminRepository, logger *zap.Logger, config *common.Config) *GarminService {
	httpClient := httptrace.WrapClient(&http.Client{
		Timeout: 30 * time.Second,
	}, httptrace.WithService("http-client"))

	return &GarminService{
		repo:       repo,
		logger:     logger,
		config:     config,
		httpClient: httpClient,
	}
}

func (s *GarminService) FetchSplits(ctx context.Context, r *request.RefreshActivitiesRequest, activityID string) (res *domain.ActivitySplitsResponse, err error) {
	url := fmt.Sprintf("https://connect.garmin.com/gc-api/activity-service/activity/%s/splits", activityID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.logger.Error("error_make_new_request", zap.Error(err))
		return nil, err
	}

	// Header
	req.Header.Set("accept", "*/*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Connect-Csrf-Token", r.GarminCsrfToken)

	// Cookie
	req.Header.Set("Cookie", r.Cookies)

	// Kirim request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Warn("request_failed_retrying",
			zap.Error(err),
		)

		s.logger.Error("error_do_request_final", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// Handle rate limiting
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("API rate limited after %d attempts", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			s.logger.Error("error_read_error_response", zap.Error(err))
			return nil, err
		}

		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(errBody))
		}

		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(errBody))
	}

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("error_read_response", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		s.logger.Error("error_unmarshal_json", zap.Error(err), zap.String("response_preview", string(body)))
		return nil, err
	}

	return res, nil
}

func (s *GarminService) Refresh(ctx context.Context, r *request.RefreshActivitiesRequest) (err error) {
	const pageSize = 20 // Increase page size for better performance
	var allActivities []*domain.Activity
	start := 0

	// s.logger.Info("starting_garmin_data_fetch", zap.String("token", r.Token[:20]+"..."))

	for {
		// Check context cancellation
		select {
		case <-ctx.Done():
			s.logger.Warn("context_cancelled_during_fetch", zap.Int("fetched_so_far", len(allActivities)))
			return ctx.Err()
		default:
		}

		pageActivities, hasMore, err := s.fetchActivitiesPage(ctx, r, start, pageSize)
		if err != nil {
			s.logger.Error("error_fetch_page",
				zap.Error(err),
				zap.Int("start", start),
				zap.Int("page_size", pageSize),
			)
			return err
		}

		// Add to total activities
		allActivities = append(allActivities, pageActivities...)

		s.logger.Info("fetched_page",
			zap.Int("page_activities", len(pageActivities)),
			zap.Int("total_activities", len(allActivities)),
			zap.Int("start", start),
		)

		// If no more data or empty response, break
		if !hasMore || len(pageActivities) == 0 {
			s.logger.Info("fetch_completed", zap.Int("total_activities", len(allActivities)))
			break
		}

		// Move to next page
		start += pageSize

		// Safety check to prevent infinite loop
		if start > 10000 { // Max 10k activities
			s.logger.Warn("max_activities_limit_reached", zap.Int("limit", 10000))
			break
		}
	}

	// Upsert all activities at once
	if len(allActivities) > 0 {
		err = s.Upsert(allActivities)
		if err != nil {
			s.logger.Error("error_upsert_all", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *GarminService) fetchActivitiesPage(ctx context.Context, r *request.RefreshActivitiesRequest, start, limit int) ([]*domain.Activity, bool, error) {
	url := fmt.Sprintf("https://connect.garmin.com/gc-api/activitylist-service/activities/search/activities?limit=%d&start=%d", limit, start)

	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			s.logger.Error("error_make_new_request", zap.Error(err))
			return nil, false, err
		}

		// Header
		req.Header.Set("accept", "*/*")
		req.Header.Set("di-backend", "connectapi.garmin.com")
		req.Header.Set("Connect-Csrf-Token", r.GarminCsrfToken)
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0")

		// Cookie
		req.Header.Set("Cookie", r.Cookies)

		// Kirim request
		resp, err := s.httpClient.Do(req)
		if err != nil {
			lastErr = err
			s.logger.Warn("request_failed_retrying",
				zap.Error(err),
				zap.Int("attempt", attempt),
				zap.Int("max_retries", maxRetries),
			)

			if attempt < maxRetries {
				// Wait before retry with exponential backoff
				waitTime := time.Duration(attempt) * time.Second
				time.Sleep(waitTime)
				continue
			}

			s.logger.Error("error_do_request_final", zap.Error(err))
			return nil, false, err
		}
		defer resp.Body.Close()

		// Handle rate limiting
		if resp.StatusCode == 429 {
			lastErr = fmt.Errorf("rate limited by API")
			s.logger.Warn("rate_limited_retrying",
				zap.Int("attempt", attempt),
				zap.Int("max_retries", maxRetries),
			)

			if attempt < maxRetries {
				// Wait longer for rate limiting
				waitTime := time.Duration(attempt*5) * time.Second
				s.logger.Info("waiting_for_rate_limit", zap.Duration("wait_time", waitTime))
				time.Sleep(waitTime)
				continue
			}

			return nil, false, fmt.Errorf("API rate limited after %d attempts", maxRetries)
		}

		if resp.StatusCode != http.StatusOK {
			errBody, err := io.ReadAll(resp.Body)
			if err != nil {
				s.logger.Error("error_read_error_response", zap.Error(err))
				return nil, false, err
			}

			lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(errBody))
			s.logger.Error("error_response",
				zap.Int("status_code", resp.StatusCode),
				zap.String("body", string(errBody)),
				zap.String("url", url),
				zap.Int("attempt", attempt),
			)

			// Don't retry for client errors (4xx except 429)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
				return nil, false, lastErr
			}

			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				time.Sleep(waitTime)
				continue
			}

			return nil, false, lastErr
		}

		// Baca response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			s.logger.Error("error_read_response", zap.Error(err))
			return nil, false, err
		}

		var pageActivities []*domain.Activity
		err = json.Unmarshal(body, &pageActivities)
		if err != nil {
			s.logger.Error("error_unmarshal_json", zap.Error(err), zap.String("response_preview", string(body[:min(500, len(body))])))
			return nil, false, err
		}

		// Success - determine if there are more pages
		// If we got less than the limit, it's likely the last page
		hasMore := len(pageActivities) == limit

		s.logger.Debug("page_fetch_success",
			zap.Int("activities_count", len(pageActivities)),
			zap.Bool("has_more", hasMore),
			zap.Int("attempt", attempt),
		)

		return pageActivities, hasMore, nil
	}

	return nil, false, lastErr
}

func (s *GarminService) Upsert(models []*domain.Activity) (err error) {
	if len(models) == 0 {
		s.logger.Info("no_activities_to_upsert")
		return nil
	}

	s.logger.Info("starting_upsert", zap.Int("total_activities", len(models)))

	const batchSize = 50
	totalBatches := (len(models) + batchSize - 1) / batchSize

	for i := 0; i < len(models); i += batchSize {
		end := i + batchSize
		if end > len(models) {
			end = len(models)
		}

		batch := models[i:end]
		batchNum := (i / batchSize) + 1

		s.logger.Info("processing_batch",
			zap.Int("batch_number", batchNum),
			zap.Int("total_batches", totalBatches),
			zap.Int("batch_size", len(batch)),
		)

		err = s.repo.Update(batch)
		if err != nil {
			s.logger.Error("error_update_batch",
				zap.Error(err),
				zap.Int("batch_number", batchNum),
				zap.Int("batch_size", len(batch)),
			)
			return err
		}

		s.logger.Info("batch_completed",
			zap.Int("batch_number", batchNum),
			zap.Int("processed", end),
			zap.Int("total", len(models)),
		)
	}

	s.logger.Info("upsert_completed", zap.Int("total_activities", len(models)))

	return nil
}

func (s *GarminService) HeartRateByDate(ctx context.Context, r *request.HeartRateByDateRequest) (err error) {
	url := fmt.Sprintf("https://connect.garmin.com/gc-api/wellness-service/wellness/dailyHeartRate?date=%s", r.Date)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.logger.Error("error_make_new_request", zap.Error(err))
		return err
	}

	// Header
	req.Header.Set("accept", "*/*")
	req.Header.Set("di-backend", "connectapi.garmin.com")
	req.Header.Set("Connect-Csrf-Token", r.GarminCsrfToken)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0")

	// Cookie
	req.Header.Set("Cookie", r.Cookies)

	// Kirim request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("error_do_request_final", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("error_response_status", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("error_read_response", zap.Error(err))
		return err
	}

	var hrData *domain.HeartRate
	err = json.Unmarshal(body, &hrData)
	if err != nil {
		s.logger.Error("error_unmarshal_json", zap.Error(err), zap.String("response_preview", string(body[:min(500, len(body))])))
		return err
	}

	s.logger.Error("read_response", zap.Any("response_api", hrData))

	return s.upsertHeartRateByDate(ctx, hrData)
}

func (s *GarminService) upsertHeartRateByDate(ctx context.Context, models *domain.HeartRate) (err error) {
	if models == nil {
		s.logger.Info("no_heart_rate_data_to_upsert")
		return nil
	}

	s.logger.Info("starting_upsert_heart_rate", zap.Int64("user_profile_pk", models.UserProfilePK), zap.String("calendar_date", models.CalendarDate))

	err = s.repo.UpsertHeartRateByDate(models)
	if err != nil {
		s.logger.Error("error_upsert_heart_rate", zap.Error(err), zap.Int64("user_profile_pk", models.UserProfilePK), zap.String("calendar_date", models.CalendarDate))
		return err
	}

	return err
}
