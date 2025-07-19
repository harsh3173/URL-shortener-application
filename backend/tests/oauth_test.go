package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/handlers"
	"url-shortener-backend/internal/middleware"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OAuthTestSuite struct {
	suite.Suite
	app          *fiber.App
	db           *gorm.DB
	config       *config.Config
	sessionStore *middleware.SimpleSessionStore
}

func (suite *OAuthTestSuite) SetupSuite() {
	suite.config = &config.Config{
		GoogleClientID:     "test-client-id",
		GoogleClientSecret: "test-client-secret",
		SessionSecret:      "test-session-secret",
		Environment:        "test",
		DatabaseURL:        "sqlite://test.db",
		FrontendURL:        "http://localhost:3000",
	}

	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	database.DB = suite.db
	err = database.AutoMigrate()
	suite.Require().NoError(err)

	authService := services.NewAuthService()
	urlService := services.NewURLService()
	
	suite.sessionStore = middleware.NewSimpleSessionStore(suite.config)
	oauthHandler := handlers.NewOAuthHandler(authService, suite.config, suite.sessionStore)
	urlHandler := handlers.NewURLHandler(urlService)

	suite.app = fiber.New()
	
	// OAuth routes
	auth := suite.app.Group("/auth")
	auth.Get("/login", oauthHandler.Login)
	auth.Get("/callback", oauthHandler.Callback)
	auth.Post("/logout", oauthHandler.Logout)
	auth.Get("/profile", suite.sessionStore.AuthMiddleware(), oauthHandler.GetProfile)
	
	// URL routes
	urls := suite.app.Group("/urls")
	urls.Post("/", suite.sessionStore.OptionalAuthMiddleware(), urlHandler.CreateURL)
	urls.Get("/", suite.sessionStore.AuthMiddleware(), urlHandler.GetUserURLs)
	urls.Get("/:shortCode/info", urlHandler.GetURLInfo)
	
	suite.app.Get("/:shortCode", urlHandler.RedirectURL)
}

func (suite *OAuthTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM urls")
}

func (suite *OAuthTestSuite) TestOAuthLoginGeneratesURL() {
	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	
	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)
	suite.True(response.Success)
	suite.Contains(response.Data.(map[string]interface{})["auth_url"], "accounts.google.com")
}

func (suite *OAuthTestSuite) TestCreateURLAnonymous() {
	reqBody := models.CreateURLRequest{
		OriginalURL: "https://example.com",
		Title:       "Test URL",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/urls/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var response models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)
	suite.True(response.Success)
}

func (suite *OAuthTestSuite) TestCreateURLWithSession() {
	// Create a test user first
	user := models.User{
		Name:    "Test User",
		Email:   "test@example.com",
		Picture: "https://example.com/avatar.jpg",
	}
	suite.db.Create(&user)

	// Manually create a session in the session store
	sessionID := "test-session-id"
	suite.sessionStore.Sessions[sessionID] = &middleware.SessionData{
		UserID:    user.ID,
		UserEmail: user.Email,
		CreatedAt: time.Now(),
	}
	
	reqBody := models.CreateURLRequest{
		OriginalURL: "https://example.com",
		Title:       "Test URL",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/urls/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "session_id="+sessionID)

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var response models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)
	suite.True(response.Success)
}

func (suite *OAuthTestSuite) TestGetURLInfo() {
	url := models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		Title:       "Test URL",
		IsActive:    true,
	}
	suite.db.Create(&url)

	req := httptest.NewRequest(http.MethodGet, "/urls/test123/info", nil)
	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)
	suite.True(response.Success)
}

func (suite *OAuthTestSuite) TestRedirectURL() {
	url := models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		IsActive:    true,
	}
	suite.db.Create(&url)

	req := httptest.NewRequest(http.MethodGet, "/test123", nil)
	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusFound, resp.StatusCode)
	suite.Equal("https://example.com", resp.Header.Get("Location"))
}

func (suite *OAuthTestSuite) TestUnauthorizedAccessToProtectedRoute() {
	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestOAuthTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthTestSuite))
}