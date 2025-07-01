package gorm_test

import (
	"testing"
	"time"

	"event-registration/internal/core/domain"

	repo "event-registration/internal/repository/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	cleanup := func() { db.Close() }
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		cleanup()
		t.Fatalf("failed to open gorm db: %v", err)
	}
	return gormDB, mock, cleanup
}

type AuthRepoTestSuite struct {
	suite.Suite
	db      *gorm.DB
	mock    sqlmock.Sqlmock
	repo    *repo.AuthRepo
	cleanup func()
}

func (s *AuthRepoTestSuite) SetupTest() {
	db, mock, cleanup := setupMockDB(s.T())
	s.db = db
	s.mock = mock
	logger := zap.NewNop()
	s.repo = repo.NewAuthRepo(db, logger).(*repo.AuthRepo)
	s.cleanup = cleanup
}

func (s *AuthRepoTestSuite) TearDownTest() {
	s.cleanup()
}

func (s *AuthRepoTestSuite) TestIsRegistered() {
	tests := []struct {
		name      string
		email     string
		mockRows  int64
		expectReg bool
	}{
		{"user not registered", "notfound@example.com", 0, false},
		{"user registered", "found@example.com", 1, true},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			rows := sqlmock.NewRows([]string{"count"}).AddRow(tc.mockRows)
			s.mock.ExpectQuery("SELECT").WithArgs(tc.email).WillReturnRows(rows)
			reg, err := s.repo.IsRegistered(tc.email)
			require.NoError(s.T(), err)
			require.Equal(s.T(), tc.expectReg, reg)
			require.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

func (s *AuthRepoTestSuite) TestRegister() {
	s.Run("register new user", func() {
		user := domain.User{Email: "newuser@example.com"}
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.Register(user)
		require.NoError(s.T(), err)
		require.NoError(s.T(), s.mock.ExpectationsWereMet())
	})
}

func (s *AuthRepoTestSuite) TestFindByEmail() {
	s.Run("user exists", func() {
		user := domain.User{ID: "1", Email: "exists@example.com"}
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "email", "password", "name", "picture", "email_verified_at", "created_at", "updated_at"}).
			AddRow(user.ID, user.Email, user.Password, user.Name, user.Picture, nil, now, now)
		s.mock.ExpectQuery("SELECT").WithArgs(user.Email, 1).WillReturnRows(rows)

		result, err := s.repo.FindByEmail(user.Email)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), result)
		require.Equal(s.T(), user.Email, result.Email)
		require.NoError(s.T(), s.mock.ExpectationsWereMet())
	})

	s.Run("user not found", func() {
		s.mock.ExpectQuery("SELECT").WithArgs("notfound@example.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "picture", "email_verified_at", "created_at", "updated_at"}))
		result, err := s.repo.FindByEmail("notfound@example.com")
		require.Error(s.T(), err)
		if result != nil {
			require.Empty(s.T(), result.Email)
		}
		require.NoError(s.T(), s.mock.ExpectationsWereMet())
	})
}

func TestAuthRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepoTestSuite))
}
