package redis

import (
	"context"
	"encoding/json"
	"event-registration/internal/core/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	client *redis.Client
}

func NewCacheRepo(client *redis.Client) *CacheRepo {
	return &CacheRepo{client: client}
}

func (r *CacheRepo) Get(key string) (*domain.Event, error) {
	data, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	var event domain.Event
	err = json.Unmarshal(data, &event)
	return &event, err
}

func (r *CacheRepo) Set(key string, event *domain.Event, expiration time.Duration) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), key, data, expiration).Err()
}
