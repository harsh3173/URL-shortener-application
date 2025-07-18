package main

import (
	"fmt"
	"log"
	"time"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/handlers"
	"url-shortener-backend/internal/middleware"
	"url-shortener-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()
	
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	authService := services.NewAuthService()
	urlService := services.NewURLService()
	
	authHandler := handlers.NewAuthHandler(authService, cfg)
	urlHandler := handlers.NewURLHandler(urlService)
	
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			
			return c.Status(code).JSON(fiber.Map{
				"error":   "internal_server_error",
				"message": err.Error(),
			})
		},
	})
	
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.CORSMiddleware(cfg))
	app.Use(middleware.RateLimitMiddleware(cfg))
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
	})
	
	apiV1 := app.Group("/api/v1")
	
	auth := apiV1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/profile", middleware.AuthMiddleware(cfg), authHandler.GetProfile)
	auth.Post("/refresh", middleware.AuthMiddleware(cfg), authHandler.RefreshToken)
	
	urls := apiV1.Group("/urls")
	urls.Post("/", middleware.OptionalAuthMiddleware(cfg), urlHandler.CreateURL)
	urls.Get("/", middleware.AuthMiddleware(cfg), urlHandler.GetUserURLs)
	urls.Put("/:id", middleware.AuthMiddleware(cfg), urlHandler.UpdateURL)
	urls.Delete("/:id", middleware.AuthMiddleware(cfg), urlHandler.DeleteURL)
	urls.Get("/:id/analytics", middleware.AuthMiddleware(cfg), urlHandler.GetURLAnalytics)
	urls.Get("/:shortCode/info", urlHandler.GetURLInfo)
	
	app.Get("/:shortCode", urlHandler.RedirectURL)
	
	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(port))
}