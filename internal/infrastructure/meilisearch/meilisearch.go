package meili

import (
	"event-registration/internal/common"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

func NewMeilisearchClient(cfg *common.Config, logger *zap.Logger) meilisearch.ServiceManager {
	client := meilisearch.New(cfg.MeilisearchHost, meilisearch.WithAPIKey(cfg.MeilisearchAPIKey))

	if _, err := client.Health(); err != nil {
		logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.String("host", cfg.MeilisearchHost),
			zap.Error(err),
		)
	} else {
		logger.Info(
			"meilisearch_connected",
			zap.String("host", cfg.MeilisearchHost),
		)
	}

	return client
}
