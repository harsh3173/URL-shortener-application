package utils

import (
	"crypto/md5"
	"encoding/hex"
	//"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func GenerateShortCodeFromURL(originalURL string, length int) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	
	if len(hash) < length {
		return hash + GenerateShortCode(length-len(hash))
	}
	return hash[:length]
}

func IsValidURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func NormalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	
	if parsedURL.Path == "" {
		parsedURL.Path = "/"
	}
	
	return parsedURL.String()
}

func IsValidCustomAlias(alias string) bool {
	if len(alias) < 3 || len(alias) > 50 {
		return false
	}
	
	for _, char := range alias {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}
	
	reservedWords := []string{"api", "admin", "www", "app", "dashboard", "login", "register", "logout", "profile", "settings"}
	for _, word := range reservedWords {
		if strings.EqualFold(alias, word) {
			return false
		}
	}
	
	return true
}

func ParseUserAgent(userAgent string) (device, os, browser string) {
	ua := strings.ToLower(userAgent)
	
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		device = "mobile"
	} else if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		device = "tablet"
	} else {
		device = "desktop"
	}
	
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac"):
		os = "macOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	default:
		os = "Unknown"
	}
	
	switch {
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	case strings.Contains(ua, "edge"):
		browser = "Edge"
	case strings.Contains(ua, "opera"):
		browser = "Opera"
	default:
		browser = "Unknown"
	}
	
	return
}