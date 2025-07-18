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
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite
	app    *fiber.App
	db     *gorm.DB
	config *config.Config
}

func (suite *AuthTestSuite) SetupSuite() {
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

	authService := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(authService, suite.config)

	suite.app = fiber.New()
	auth := suite.app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/profile", middleware.AuthMiddleware(suite.config), authHandler.GetProfile)
}

func (suite *AuthTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM users")
}

func (suite *AuthTestSuite) TestRegisterSuccess() {
	reqBody := models.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var response models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)
	suite.True(response.Success)
	suite.Contains(response.Message, "successfully")
}

func (suite *AuthTestSuite) TestRegisterDuplicateEmail() {
	// Create first user
	reqBody := models.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := suite.app.Test(req)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	// Try to create second user with same email
	req = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = suite.app.Test(req)
	suite.Equal(http.StatusConflict, resp.StatusCode)
}

func (suite *AuthTestSuite) TestLoginSuccess() {
	// Register user first
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
	}
	suite.db.Create(&user)

	reqBody := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&response)
	suite.True(response.Success)
}

func (suite *AuthTestSuite) TestLoginInvalidCredentials() {
	reqBody := models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}