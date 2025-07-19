package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Name      string         `json:"name" gorm:"not null"`
	Picture   string         `json:"picture,omitempty" gorm:"size:500"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	URLs      []URL          `json:"urls,omitempty" gorm:"foreignKey:UserID"`
}

type URL struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OriginalURL string         `json:"original_url" gorm:"not null;type:text"`
	ShortCode   string         `json:"short_code" gorm:"unique;not null;size:10"`
	CustomAlias string         `json:"custom_alias,omitempty" gorm:"unique;size:50"`
	UserID      *uint          `json:"user_id,omitempty" gorm:"index"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string         `json:"title,omitempty" gorm:"size:200"`
	Description string         `json:"description,omitempty" gorm:"size:500"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Analytics   []Analytics    `json:"analytics,omitempty" gorm:"foreignKey:URLID"`
}

type Analytics struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	URLID     uint           `json:"url_id" gorm:"not null;index"`
	URL       URL            `json:"url" gorm:"foreignKey:URLID"`
	IPAddress string         `json:"ip_address" gorm:"size:45"`
	UserAgent string         `json:"user_agent" gorm:"size:500"`
	Referrer  string         `json:"referrer" gorm:"size:500"`
	Country   string         `json:"country,omitempty" gorm:"size:100"`
	City      string         `json:"city,omitempty" gorm:"size:100"`
	Device    string         `json:"device,omitempty" gorm:"size:100"`
	OS        string         `json:"os,omitempty" gorm:"size:100"`
	Browser   string         `json:"browser,omitempty" gorm:"size:100"`
	ClickedAt time.Time      `json:"clicked_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type URLStats struct {
	URLID      uint      `json:"url_id"`
	TotalClicks int64     `json:"total_clicks"`
	UniqueClicks int64    `json:"unique_clicks"`
	LastClicked *time.Time `json:"last_clicked"`
}

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Title       string `json:"title,omitempty" validate:"omitempty,max=200"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}


type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}