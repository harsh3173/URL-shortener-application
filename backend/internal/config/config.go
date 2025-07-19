package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	HTTPSPort           string
	TLSCertFile         string
	TLSKeyFile          string
	DatabaseURL         string
	GoogleClientID      string
	GoogleClientSecret  string
	SessionSecret       string
	Environment         string
	FrontendURL         string
	RedisURL            string
	RateLimitRequests   int
	RateLimitWindow     int
	MaxURLLength        int
	CustomDomainLength  int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	rateLimitRequests, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "3600"))
	maxURLLength, _ := strconv.Atoi(getEnv("MAX_URL_LENGTH", "2048"))
	customDomainLength, _ := strconv.Atoi(getEnv("CUSTOM_DOMAIN_LENGTH", "6"))

	return &Config{
		Port:                getEnv("PORT", "8080"),
		HTTPSPort:           getEnv("HTTPS_PORT", "8443"),
		TLSCertFile:         getEnv("TLS_CERT_FILE", "/app/certs/localhost.crt"),
		TLSKeyFile:          getEnv("TLS_KEY_FILE", "/app/certs/localhost.key"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://user:password@localhost/urlshortener?sslmode=disable"),
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
		SessionSecret:       getEnv("SESSION_SECRET", "your-256-bit-session-secret"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		FrontendURL:         getEnv("FRONTEND_URL", "https://localhost:3000"),
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		RateLimitRequests:   rateLimitRequests,
		RateLimitWindow:     rateLimitWindow,
		MaxURLLength:        maxURLLength,
		CustomDomainLength:  customDomainLength,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}