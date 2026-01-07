package meilisearch

import (
	"context"
	"event-registration/internal/common/helper"
	"event-registration/internal/core/domain"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const USER_INDEX = "users"

type UserMeilisearchRepo struct {
	logger      *zap.Logger
	meilisearch meilisearch.ServiceManager
	repo        domain.UserRepository
}

func NewUserMeilisearchRepo(
	logger *zap.Logger,
	meilisearch meilisearch.ServiceManager,
	repo domain.UserRepository,
) domain.UserMeilisearchRepository {

	helper.PrettyPrint(meilisearch, "meilisearch_client ==========================")

	return &UserMeilisearchRepo{
		logger:      logger,
		meilisearch: meilisearch,
		repo:        repo,
	}
}

func (r *UserMeilisearchRepo) SetupIndex() (err error) {
	index := r.meilisearch.Index(USER_INDEX)

	taskInfo, err := index.UpdateSearchableAttributes(&[]string{
		"email",
		"username",
		"full_name",
		"jabatan",
		"nip",
		"unit_code",
		"unit_name",
	})

	if err != nil {
		r.logger.Error(
			"error_updating_meilisearch_searchable_attributes",
			zap.Error(err),
		)

		return err
	}

	r.logger.Info(
		"task_info",
		zap.Any("task_info", taskInfo),
	)

	taskInfo, err = index.UpdateFilterableAttributes(&[]interface{}{"email",
		"username",
		"full_name",
		"jabatan",
		"nip",
		"unit_code",
		"unit_name"})
	if err != nil {
		r.logger.Error(
			"error_updating_meilisearch_filterable_attributes",
			zap.Error(err),
		)

		return err
	}

	r.logger.Info(
		"task_info",
		zap.Any("task_info", taskInfo),
	)

	taskInfo, err = index.UpdateSortableAttributes(&[]string{
		"email",
		"username",
		"full_name",
		"jabatan",
		"nip",
		"unit_code",
		"unit_name",
	})
	if err != nil {
		r.logger.Error(
			"error_updating_meilisearch_sortable_attributes",
			zap.Error(err),
		)

		return err
	}

	r.logger.Info(
		"task_info",
		zap.Any("task_info", taskInfo),
	)

	return nil
}

func (r *UserMeilisearchRepo) SeedIndex() error {
	index := r.meilisearch.Index("users")

	users, err := r.repo.FindAll()
	if err != nil {
		r.logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.Error(err),
		)

		return err
	}

	taskInfo, err := index.AddDocuments(users, nil)
	if err != nil {
		r.logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.Error(err),
		)

		return err
	}

	r.logger.Info(
		"task_info",
		zap.Any("task_info", taskInfo),
	)

	return nil
}

func (r *UserMeilisearchRepo) Search(ctx context.Context, keyword string) (users []*domain.UserVCC, err error) {
	index := r.meilisearch.Index("users")

	searchRes, err := index.SearchWithContext(ctx, keyword, &meilisearch.SearchRequest{
		Limit: 10,
	})
	if err != nil {
		r.logger.Error("error_search_users_meilisearch", zap.Error(err))
		return nil, err
	}

	users = make([]*domain.UserVCC, 0, len(searchRes.Hits))

	for _, hit := range searchRes.Hits {
		user := &domain.UserVCC{}
		if err := hit.DecodeInto(user); err != nil {
			r.logger.Error("error_decode_user",
				zap.Error(err),
				zap.Any("hit", hit))
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserMeilisearchRepo) CheckHealth() error {
	if _, err := r.meilisearch.Health(); err != nil {
		r.logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.Error(err),
		)

		return err
	} else {
		r.logger.Info(
			"meilisearch_connected",
		)
	}

	return nil
}

func (r *UserMeilisearchRepo) Update(user *domain.UserVCC) {
	ctx := context.Background()
	index := r.meilisearch.Index("users")

	primaryKey := "id"
	taskInfo, err := index.AddDocumentsWithContext(
		ctx,
		[]interface{}{user},
		&meilisearch.DocumentOptions{PrimaryKey: &primaryKey},
	)
	if err != nil {
		r.logger.Error("error_updating_meilisearch_index",
			zap.String("user_id", user.ID),
			zap.Error(err))
		return
	}

	r.logger.Info("meilisearch_index_updated",
		zap.String("user_id", user.ID),
		zap.Int64("task_uid", taskInfo.TaskUID))
}
