package handlers

import (
	"time"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/services"
	"url-shortener-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}
	
	user, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
	}
	
	token, err := utils.GenerateJWT(user.ID, user.Email, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication token",
		})
	}
	
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Strict",
	})
	
	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
		Message: "User registered successfully",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}
	
	user, err := h.authService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
	}
	
	token, err := utils.GenerateJWT(user.ID, user.Email, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication token",
		})
	}
	
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Strict",
	})
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
		Message: "Login successful",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Strict",
	})
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
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

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	userEmail := c.Locals("user_email").(string)
	
	token, err := utils.GenerateJWT(userID, userEmail, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication token",
		})
	}
	
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Strict",
	})
	
	return c.JSON(models.SuccessResponse{
		Success: true,
		Data: fiber.Map{
			"token": token,
		},
		Message: "Token refreshed successfully",
	})
}