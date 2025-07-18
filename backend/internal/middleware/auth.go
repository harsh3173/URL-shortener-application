package middleware

import (
	"strings"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		
		if cookie := c.Cookies("token"); cookie != "" {
			token = cookie
		} else if auth := c.Get("Authorization"); auth != "" {
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "No token provided",
			})
		}
		
		claims, err := utils.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid token",
			})
		}
		
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		
		return c.Next()
	}
}

func OptionalAuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		
		if cookie := c.Cookies("token"); cookie != "" {
			token = cookie
		} else if auth := c.Get("Authorization"); auth != "" {
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		
		if token != "" {
			if claims, err := utils.ValidateJWT(token, cfg.JWTSecret); err == nil {
				c.Locals("user_id", claims.UserID)
				c.Locals("user_email", claims.Email)
			}
		}
		
		return c.Next()
	}
}