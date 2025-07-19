package middleware

import (
	"crypto/rand"
	"encoding/base64"
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
	sessions map[string]*SessionData
	mu       sync.RWMutex
	config   *config.Config
}

func NewSimpleSessionStore(config *config.Config) *SimpleSessionStore {
	store := &SimpleSessionStore{
		sessions: make(map[string]*SessionData),
		config:   config,
	}
	
	// Start cleanup routine
	go store.cleanup()
	
	return store
}

func (s *SimpleSessionStore) cleanup() {
	for {
		time.Sleep(1 * time.Hour)
		s.mu.Lock()
		for sessionID, data := range s.sessions {
			if time.Since(data.CreatedAt) > 24*time.Hour {
				delete(s.sessions, sessionID)
			}
		}
		s.mu.Unlock()
	}
}

func (s *SimpleSessionStore) generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *SimpleSessionStore) CreateSession(c *fiber.Ctx, userID uint, userEmail string) string {
	sessionID := s.generateSessionID()
	
	s.mu.Lock()
	s.sessions[sessionID] = &SessionData{
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
	
	return sessionID
}

func (s *SimpleSessionStore) GetSession(sessionID string) *SessionData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if data, exists := s.sessions[sessionID]; exists {
		if time.Since(data.CreatedAt) < 24*time.Hour {
			return data
		}
		// Session expired, delete it
		delete(s.sessions, sessionID)
	}
	return nil
}

func (s *SimpleSessionStore) DestroySession(c *fiber.Ctx, sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
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