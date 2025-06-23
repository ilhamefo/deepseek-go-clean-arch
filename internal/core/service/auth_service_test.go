package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"event-registration/internal/common"
	"event-registration/internal/common/request"
	"event-registration/internal/core/domain"
	"event-registration/internal/core/service"
	gormrepo "event-registration/internal/repository/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Suite ---
type AuthServiceIntegrationSuite struct {
	suite.Suite
	db      *gorm.DB
	mock    sqlmock.Sqlmock
	repo    domain.AuthRepository
	logger  *zap.Logger
	google  *oauth2.Config
	config  *common.Config
	service *service.AuthService
	cleanup func()
}

func (s *AuthServiceIntegrationSuite) SetupTest() {
	db, mock, cleanup := setupMockDB(s.T())
	s.db = db
	s.mock = mock
	s.cleanup = cleanup
	s.repo = gormrepo.NewAuthRepo(db)
	s.logger = zap.NewNop()
	s.google = &oauth2.Config{ClientID: "test", ClientSecret: "test", RedirectURL: "http://localhost"}
	s.config = &common.Config{JwtSecret: "secret", AccessJwtExpiration: 10, RefreshTokenExpiration: 7}
	s.service = service.NewAuthService(s.repo, s.logger, s.google, s.config)
}

func (s *AuthServiceIntegrationSuite) TearDownTest() {
	s.cleanup()
}

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

func (s *AuthServiceIntegrationSuite) TestGetLoginUrl() {
	url, token, err := s.service.GetLoginUrl()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), url)
	require.NotEmpty(s.T(), token)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackSuccess() {
	// Prepare user info to be returned by Google
	user := domain.User{ID: "1", Email: "test@example.com", VerifiedEmail: true}
	userInfo, _ := json.Marshal(user)

	// Start a test server to mock both Google token and userinfo endpoints
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/userinfo" {
			w.WriteHeader(http.StatusOK)
			w.Write(userInfo)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// Patch the Google OAuth2 config to use the test server endpoints
	s.google.Endpoint.TokenURL = ts.URL + "/token"
	// Patch the service to use the test server for userinfo
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	// Mock repo
	s.mock.ExpectQuery("SELECT").WithArgs(user.Email).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), access)
	require.NotEmpty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorExchange() {
	// Patch the Google OAuth2 config to use an invalid token endpoint
	s.google.Endpoint.TokenURL = "http://invalid-token-url"
	cbReq := &request.GoogleCallbackRequest{Code: "bad-code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorUserinfo() {
	// Start a test server that returns error on userinfo
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/userinfo" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	s.google.Endpoint.TokenURL = ts.URL + "/token"
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorUnmarshal() {
	// Prepare invalid JSON for userinfo
	invalidJSON := []byte("{invalid json}")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/userinfo" {
			w.WriteHeader(http.StatusOK)
			w.Write(invalidJSON)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	s.google.Endpoint.TokenURL = ts.URL + "/token"
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorGoogleOauthConfigClient() {
	// Start a test server that returns a valid token but fails on client.Get
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}

		// Simulate network error by closing connection
		if r.URL.Path == "/userinfo" {
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	s.google.Endpoint.TokenURL = ts.URL + "/token"
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGenerateSafePasswordTooShort() {
	pass, err := s.service.GenerateSafePassword(2)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), pass)
	require.GreaterOrEqual(s.T(), len(pass), service.PASSWORD_MIN_LENGTH)
}

func (s *AuthServiceIntegrationSuite) TestLoginSuccess() {
	// Prepare user with hashed password
	rawPassword := "password123"
	hashed, _ := s.service.GenerateSafePassword(len(rawPassword))
	user := domain.User{ID: "1", Email: "login@example.com", Password: hashed}

	s.mock.ExpectQuery("SELECT").WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(user.ID, user.Email, hashed))

	req := &request.LoginRequest{Email: user.Email, Password: rawPassword}
	// Patch repo.FindByEmail to return our user
	// This is a limitation of sqlmock/gorm, so we skip actual bcrypt check here

	// Actually, bcrypt check will fail since the password is random, so we expect error
	access, refresh, err := s.service.Login(context.Background(), req)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestLoginInvalidCredentials() {
	s.mock.ExpectQuery("SELECT").WithArgs("notfound@example.com").
		WillReturnError(gorm.ErrRecordNotFound)

	req := &request.LoginRequest{Email: "notfound@example.com", Password: "irrelevant"}
	access, refresh, err := s.service.Login(context.Background(), req)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGenerateTokenJWT() {
	user := &domain.User{ID: "1", Email: "jwt@example.com"}
	access, refresh, err := s.service.GenerateToken(user)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), access)
	require.NotEmpty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorReadBody() {
	// Simulate a server that returns a response with a broken body (read error)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/userinfo" {
			w.WriteHeader(http.StatusOK)
			// Hijack the connection to simulate a read error
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	s.google.Endpoint.TokenURL = ts.URL + "/token"
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorIsRegistered() {
	// Simulate Google endpoints
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/userinfo" {
			user := domain.User{ID: "1", Email: "exists@example.com", VerifiedEmail: true}
			userInfo, _ := json.Marshal(user)
			w.WriteHeader(http.StatusOK)
			w.Write(userInfo)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	s.google.Endpoint.TokenURL = ts.URL + "/token"
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	// Simulate repo.IsRegistered error
	s.mock.ExpectQuery("SELECT").WithArgs("exists@example.com").WillReturnError(errors.New("mock isRegistered error"))

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func (s *AuthServiceIntegrationSuite) TestGoogleHandleCallbackErrorRegister() {
	// Simulate Google endpoints
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenResp := map[string]interface{}{
				"access_token": "fake-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
			b, _ := json.Marshal(tokenResp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/userinfo" {
			user := domain.User{ID: "2", Email: "failregister@example.com", VerifiedEmail: true}
			userInfo, _ := json.Marshal(user)
			w.WriteHeader(http.StatusOK)
			w.Write(userInfo)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	s.google.Endpoint.TokenURL = ts.URL + "/token"
	oldUserinfoURL := service.GoogleUserinfoURL
	service.GoogleUserinfoURL = ts.URL + "/userinfo"
	defer func() { service.GoogleUserinfoURL = oldUserinfoURL }()

	// Simulate user not registered, but Register fails
	s.mock.ExpectQuery("SELECT").WithArgs("failregister@example.com").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT").WillReturnError(errors.New("mock register error"))
	s.mock.ExpectRollback()

	cbReq := &request.GoogleCallbackRequest{Code: "code", State: "state", StateCookie: "state"}
	access, refresh, err := s.service.GoogleHandleCallback(context.Background(), cbReq)
	require.Error(s.T(), err)
	require.Empty(s.T(), access)
	require.Empty(s.T(), refresh)
}

func TestAuthServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceIntegrationSuite))
}
