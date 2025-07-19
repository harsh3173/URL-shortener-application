package services

import (
	"errors"
	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}
	
	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}
	
	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	return &user, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}
	
	return &user, nil
}

func (s *AuthService) LoginOrRegisterOAuth(email, name, picture string) (*models.User, error) {
	var user models.User
	
	// Try to find existing user
	err := s.db.Where("email = ?", email).First(&user).Error
	if err == nil {
		// User exists, update profile picture if provided
		if picture != "" && user.Picture != picture {
			user.Picture = picture
			s.db.Save(&user)
		}
		return &user, nil
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("database error")
	}
	
	// User doesn't exist, create new user
	user = models.User{
		Email:    email,
		Name:     name,
		Picture:  picture,
		Password: "", // OAuth users don't have passwords
	}
	
	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}
	
	return &user, nil
}