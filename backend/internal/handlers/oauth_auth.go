package handlers

import (
	"context"
	"time"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/middleware"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/services"
	"url-shortener-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	authService   *services.AuthService
	config        *config.Config
	oauthConfig   *oauth2.Config
	sessionStore  *middleware.SimpleSessionStore
}

func NewOAuthHandler(authService *services.AuthService, config *config.Config, sessionStore *middleware.SimpleSessionStore) *OAuthHandler {
	oauthConfig := utils.NewGoogleOAuthConfig(
		config.GoogleClientID,
		config.GoogleClientSecret,
		config.FrontendURL+"/auth/callback",
	)

	return &OAuthHandler{
		authService:  authService,
		config:       config,
		oauthConfig:  oauthConfig,
		sessionStore: sessionStore,
	}
}

func (h *OAuthHandler) Login(c *fiber.Ctx) error {
	state, err := utils.GenerateStateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "state_generation_failed",
			Message: "Failed to generate state token",
		})
	}

	// Store state in a temporary session-like cookie
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute), // Short-lived
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Lax",
	})

	authURL := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	return c.JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"auth_url": authURL,
		},
		Message: "OAuth login URL generated",
	})
}

func (h *OAuthHandler) Callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_callback",
			Message: "Missing code or state parameter",
		})
	}

	// Verify state token
	storedState := c.Cookies("oauth_state")
	if storedState == "" || storedState != state {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_state",
			Message: "Invalid state parameter",
		})
	}

	// Clear the state cookie
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Lax",
	})

	// Exchange code for token
	ctx := context.Background()
	token, err := h.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "token_exchange_failed",
			Message: "Failed to exchange code for token",
		})
	}

	// Get user info from Google
	userInfo, err := utils.GetGoogleUserInfo(ctx, h.oauthConfig, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "userinfo_failed",
			Message: "Failed to get user information",
		})
	}

	// Register or login user
	user, err := h.authService.LoginOrRegisterOAuth(userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
	}

	// Create session
	h.sessionStore.CreateSession(c, user.ID, user.Email)

	return c.JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"user": user,
		},
		Message: "Login successful",
	})
}

func (h *OAuthHandler) Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		h.sessionStore.DestroySession(c, sessionID)
	}

	return c.JSON(models.SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

func (h *OAuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "user_not_found",
			Message: err.Error(),
		})
	}

	return c.JSON(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}