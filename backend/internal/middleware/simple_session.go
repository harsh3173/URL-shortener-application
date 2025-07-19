package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type SessionData struct {
	UserID    uint
	UserEmail string
	CreatedAt time.Time
}

type SimpleSessionStore struct {
	Sessions map[string]*SessionData // Made public for testing
	mu       sync.RWMutex
	config   *config.Config
	done     chan struct{} // For graceful shutdown
}

func NewSimpleSessionStore(config *config.Config) *SimpleSessionStore {
	store := &SimpleSessionStore{
		Sessions: make(map[string]*SessionData),
		config:   config,
		done:     make(chan struct{}),
	}
	
	// Start cleanup routine with graceful shutdown
	go store.cleanup()
	
	return store
}

func (s *SimpleSessionStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			expired := 0
			for sessionID, data := range s.Sessions {
				if time.Since(data.CreatedAt) > 24*time.Hour {
					delete(s.Sessions, sessionID)
					expired++
				}
			}
			s.mu.Unlock()
			// Log cleanup activity for monitoring
			if expired > 0 {
				fmt.Printf("Session cleanup: removed %d expired sessions\n", expired)
			}
		case <-s.done:
			return // Graceful shutdown
		}
	}
}

// Close gracefully shuts down the session store
func (s *SimpleSessionStore) Close() {
	close(s.done)
}

func (s *SimpleSessionStore) generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate secure session ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *SimpleSessionStore) CreateSession(c *fiber.Ctx, userID uint, userEmail string) (string, error) {
	sessionID, err := s.generateSessionID()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	
	s.mu.Lock()
	s.Sessions[sessionID] = &SessionData{
		UserID:    userID,
		UserEmail: userEmail,
		CreatedAt: time.Now(),
	}
	s.mu.Unlock()
	
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   s.config.Environment == "production",
		SameSite: "Lax",
	})
	
	return sessionID, nil
}

func (s *SimpleSessionStore) GetSession(sessionID string) *SessionData {
	s.mu.Lock() // Use write lock to allow safe deletion
	defer s.mu.Unlock()
	
	if data, exists := s.Sessions[sessionID]; exists {
		if time.Since(data.CreatedAt) < 24*time.Hour {
			return data
		}
		// Session expired, delete it safely
		delete(s.Sessions, sessionID)
	}
	return nil
}

func (s *SimpleSessionStore) DestroySession(c *fiber.Ctx, sessionID string) {
	s.mu.Lock()
	delete(s.Sessions, sessionID)
	s.mu.Unlock()
	
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   s.config.Environment == "production",
		SameSite: "Lax",
	})
}

func (s *SimpleSessionStore) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "No session found",
			})
		}
		
		sessionData := s.GetSession(sessionID)
		if sessionData == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired session",
			})
		}
		
		c.Locals("user_id", sessionData.UserID)
		c.Locals("user_email", sessionData.UserEmail)
		c.Locals("session_id", sessionID)
		
		return c.Next()
	}
}

func (s *SimpleSessionStore) OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID != "" {
			sessionData := s.GetSession(sessionID)
			if sessionData != nil {
				c.Locals("user_id", sessionData.UserID)
				c.Locals("user_email", sessionData.UserEmail)
				c.Locals("session_id", sessionID)
			}
		}
		
		return c.Next()
	}
}