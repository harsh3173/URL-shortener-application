package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

type URLTestSuite struct {
	suite.Suite
	app    *fiber.App
	db     *gorm.DB
	config *config.Config
}

func (suite *URLTestSuite) SetupSuite() {
	suite.config = &config.Config{
		JWTSecret:   "test-secret",
		Environment: "test",
		DatabaseURL: "sqlite://test.db",
	}

	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	database.DB = suite.db
	err = database.AutoMigrate()
	suite.Require().NoError(err)

	urlService := services.NewURLService()
	urlHandler := handlers.NewURLHandler(urlService)

	suite.app = fiber.New()
	urls := suite.app.Group("/urls")
	urls.Post("/", middleware.OptionalAuthMiddleware(suite.config), urlHandler.CreateURL)
	urls.Get("/:shortCode/info", urlHandler.GetURLInfo)
	suite.app.Get("/:shortCode", urlHandler.RedirectURL)
}

func (suite *URLTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM urls")
}

func (suite *URLTestSuite) TestCreateURLSuccess() {
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

func (suite *URLTestSuite) TestCreateURLInvalidURL() {
	reqBody := models.CreateURLRequest{
		OriginalURL: "invalid-url",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/urls/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *URLTestSuite) TestCreateURLWithCustomAlias() {
	reqBody := models.CreateURLRequest{
		OriginalURL: "https://example.com",
		CustomAlias: "my-custom-link",
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

func (suite *URLTestSuite) TestGetURLInfo() {
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

func (suite *URLTestSuite) TestRedirectURL() {
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

func (suite *URLTestSuite) TestRedirectURLNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestURLTestSuite(t *testing.T) {
	suite.Run(t, new(URLTestSuite))
}