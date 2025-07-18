package services

import (
	"errors"
	"time"
	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/utils"

	"gorm.io/gorm"
)

type URLService struct {
	db *gorm.DB
}

func NewURLService() *URLService {
	return &URLService{
		db: database.GetDB(),
	}
}

func (s *URLService) CreateURL(req *models.CreateURLRequest, userID *uint) (*models.URL, error) {
	if !utils.IsValidURL(req.OriginalURL) {
		return nil, errors.New("invalid URL format")
	}
	
	normalizedURL := utils.NormalizeURL(req.OriginalURL)
	
	var shortCode string
	if req.CustomAlias != "" {
		if !utils.IsValidCustomAlias(req.CustomAlias) {
			return nil, errors.New("invalid custom alias")
		}
		
		var existingURL models.URL
		if err := s.db.Where("custom_alias = ? OR short_code = ?", req.CustomAlias, req.CustomAlias).First(&existingURL).Error; err == nil {
			return nil, errors.New("custom alias already exists")
		}
		
		shortCode = req.CustomAlias
	} else {
		shortCode = s.generateUniqueShortCode()
	}
	
	url := &models.URL{
		OriginalURL: normalizedURL,
		ShortCode:   shortCode,
		CustomAlias: req.CustomAlias,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    true,
	}
	
	if req.ExpiresAt != "" {
		if expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			url.ExpiresAt = &expiresAt
		}
	}
	
	if err := s.db.Create(url).Error; err != nil {
		return nil, errors.New("failed to create URL")
	}
	
	return url, nil
}

func (s *URLService) GetURLByShortCode(shortCode string) (*models.URL, error) {
	var url models.URL
	query := s.db.Where("(short_code = ? OR custom_alias = ?) AND is_active = ?", shortCode, shortCode, true)
	
	if err := query.First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("URL not found")
		}
		return nil, errors.New("database error")
	}
	
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("URL has expired")
	}
	
	return &url, nil
}

func (s *URLService) GetUserURLs(userID uint, limit, offset int) ([]models.URL, int64, error) {
	var urls []models.URL
	var total int64
	
	query := s.db.Where("user_id = ?", userID).Order("created_at DESC")
	
	if err := query.Model(&models.URL{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count URLs")
	}
	
	if err := query.Limit(limit).Offset(offset).Find(&urls).Error; err != nil {
		return nil, 0, errors.New("failed to fetch URLs")
	}
	
	return urls, total, nil
}

func (s *URLService) UpdateURL(urlID uint, userID uint, req *models.CreateURLRequest) (*models.URL, error) {
	var url models.URL
	if err := s.db.Where("id = ? AND user_id = ?", urlID, userID).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("URL not found")
		}
		return nil, errors.New("database error")
	}
	
	if req.OriginalURL != "" {
		if !utils.IsValidURL(req.OriginalURL) {
			return nil, errors.New("invalid URL format")
		}
		url.OriginalURL = utils.NormalizeURL(req.OriginalURL)
	}
	
	if req.Title != "" {
		url.Title = req.Title
	}
	
	if req.Description != "" {
		url.Description = req.Description
	}
	
	if req.ExpiresAt != "" {
		if expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			url.ExpiresAt = &expiresAt
		}
	}
	
	if err := s.db.Save(&url).Error; err != nil {
		return nil, errors.New("failed to update URL")
	}
	
	return &url, nil
}

func (s *URLService) DeleteURL(urlID uint, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", urlID, userID).Delete(&models.URL{})
	if result.Error != nil {
		return errors.New("failed to delete URL")
	}
	
	if result.RowsAffected == 0 {
		return errors.New("URL not found")
	}
	
	return nil
}

func (s *URLService) RecordClick(urlID uint, analytics *models.Analytics) error {
	analytics.URLID = urlID
	analytics.ClickedAt = time.Now()
	
	if err := s.db.Create(analytics).Error; err != nil {
		return errors.New("failed to record click")
	}
	
	return nil
}

func (s *URLService) GetURLAnalytics(urlID uint, userID uint) ([]models.Analytics, error) {
	var analytics []models.Analytics
	
	query := `
		SELECT a.* FROM analytics a
		JOIN urls u ON a.url_id = u.id
		WHERE u.id = ? AND u.user_id = ?
		ORDER BY a.clicked_at DESC
		LIMIT 1000
	`
	
	if err := s.db.Raw(query, urlID, userID).Scan(&analytics).Error; err != nil {
		return nil, errors.New("failed to fetch analytics")
	}
	
	return analytics, nil
}

func (s *URLService) GetURLStats(urlID uint, userID uint) (*models.URLStats, error) {
	var stats models.URLStats
	
	query := `
		SELECT 
			u.id as url_id,
			COUNT(a.id) as total_clicks,
			COUNT(DISTINCT a.ip_address) as unique_clicks,
			MAX(a.clicked_at) as last_clicked
		FROM urls u
		LEFT JOIN analytics a ON u.id = a.url_id
		WHERE u.id = ? AND u.user_id = ?
		GROUP BY u.id
	`
	
	if err := s.db.Raw(query, urlID, userID).Scan(&stats).Error; err != nil {
		return nil, errors.New("failed to fetch URL stats")
	}
	
	return &stats, nil
}

func (s *URLService) generateUniqueShortCode() string {
	for {
		code := utils.GenerateShortCode(6)
		var existingURL models.URL
		if err := s.db.Where("short_code = ? OR custom_alias = ?", code, code).First(&existingURL).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code
			}
		}
	}
}