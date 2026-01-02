package service

import (
	"context"
	"event-registration/internal/common"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"strconv"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type UserService struct {
	repo        domain.UserRepository
	meilisearch meilisearch.ServiceManager
	logger      *zap.Logger
}

func NewUserService(
	repo domain.UserRepository,
	logger *zap.Logger,
	googleConfig *oauth2.Config,
	config *common.Config,
	sessionService *SessionService,
	meilisearch meilisearch.ServiceManager,
) *UserService {
	return &UserService{
		repo:        repo,
		logger:      logger,
		meilisearch: meilisearch,
	}
}

func (s *UserService) Search(ctx context.Context, keyword string) (users []*domain.UserVCC, err error) {
	index := s.meilisearch.Index("users")

	searchRes, err := index.SearchWithContext(ctx, keyword, &meilisearch.SearchRequest{
		Limit: 10,
	})
	if err != nil {
		s.logger.Error("error_search_users_meilisearch", zap.Error(err))
		return nil, err
	}

	users = make([]*domain.UserVCC, 0, len(searchRes.Hits))

	for _, hit := range searchRes.Hits {
		user := &domain.UserVCC{}
		if err := hit.DecodeInto(user); err != nil {
			s.logger.Error("error_decode_user",
				zap.Error(err),
				zap.Any("hit", hit))
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) Roles() (roles []*domain.Role, err error) {
	roles, err = s.repo.Roles()
	if err != nil {
		s.logger.Error("error_get_roles", zap.Error(err))
		return nil, err
	}

	return roles, nil
}

func (s *UserService) GetUnits(level string) (units []*domain.UnitName, err error) {
	if level == "0" {
		units = append(units, &domain.UnitName{
			Label: "Pusat",
			Code:  "",
		})

		return units, nil
	}
	units, err = s.repo.Unit(level)
	if err != nil {
		s.logger.Error("error_get_units", zap.Error(err))
		return nil, err
	}

	return units, nil
}

func (s *UserService) Update(req *request.UpdateUserRequest) (err error) {
	level, err := strconv.Atoi(req.Level)
	if err != nil {
		s.logger.Error("error_convert_to_int", zap.Error(err))
		return err
	}
	status, err := strconv.Atoi(req.Status)
	if err != nil {
		s.logger.Error("error_convert_to_int", zap.Error(err))
		return err
	}

	user := &domain.UserVCC{
		ID:       req.ID,
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Jabatan:  req.Jabatan,
		NIP:      req.NIP,
		Level:    uint(level),
		UnitCode: &req.UnitCode,
		UnitName: &req.UnitName,
		Status:   uint(status),
	}

	for _, role := range req.Roles {
		user.Roles = append(user.Roles, &domain.Role{
			ID: role,
		})
	}

	err = s.repo.Update(user)
	if err != nil {
		s.logger.Error("error_update_user", zap.Error(err))
		return err
	}

	go func() {
		ctx := context.Background()
		index := s.meilisearch.Index("users")

		// AddDocuments with primary key "id" will automatically update existing document
		taskInfo, err := index.AddDocumentsWithContext(ctx, []interface{}{user}, &meilisearch.DocumentOptions{PrimaryKey: "id"})
		if err != nil {
			s.logger.Error("error_updating_meilisearch_index",
				zap.String("user_id", user.ID),
				zap.Error(err))
			return
		}

		s.logger.Info("meilisearch_index_updated",
			zap.String("user_id", user.ID),
			zap.Int64("task_uid", taskInfo.TaskUID))
	}()

	return nil
}

func (s *UserService) CheckHealthMeilisearch() error {
	if _, err := s.meilisearch.Health(); err != nil {
		s.logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.Error(err),
		)

		return err
	} else {
		s.logger.Info(
			"meilisearch_connected",
		)
	}

	err := s.SetupIndexUsers()
	if err != nil {
		s.logger.Error(
			"error_setup_meilisearch_index_users",
			zap.Error(err),
		)

		return err
	}

	err = s.SeedIndex()
	if err != nil {
		s.logger.Error(
			"error_seeding_meilisearch_index_users",
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (s *UserService) SetupIndexUsers() error {

	index := s.meilisearch.Index("users")

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
		s.logger.Error(
			"error_updating_meilisearch_searchable_attributes",
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
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
		s.logger.Error(
			"error_updating_meilisearch_filterable_attributes",
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
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
		s.logger.Error(
			"error_updating_meilisearch_sortable_attributes",
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
		"task_info",
		zap.Any("task_info", taskInfo),
	)

	return nil
}

func (s *UserService) SeedIndex() error {
	index := s.meilisearch.Index("users")

	users, err := s.repo.FindAll()
	if err != nil {
		s.logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.Error(err),
		)

		return err
	}

	taskInfo, err := index.AddDocuments(users, nil)
	if err != nil {
		s.logger.Error(
			"failed_to_connect_to_meilisearch",
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
		"task_info",
		zap.Any("task_info", taskInfo),
	)

	return nil
}
