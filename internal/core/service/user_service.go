package service

import (
	"context"
	"event-registration/internal/common"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type UserService struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

func NewUserService(repo domain.UserRepository, logger *zap.Logger, googleConfig *oauth2.Config, config *common.Config, sessionService *SessionService) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) Search(ctx context.Context, keyword string) (users []*domain.UserVCC, err error) {
	users, err = s.repo.Search(keyword)
	if err != nil {
		s.logger.Error("error_search_users", zap.Error(err))
		return nil, err
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

	return nil
}
