package services

import (
	"errors"
	"fmt"
	"time"

	"carbuyer/internal/db/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db                 *gorm.DB
	jwtSecret          string
	jwtExpirationHours int
	mailgunDomain      string
}

func NewAuthService(db *gorm.DB, jwtSecret string, jwtExpirationHours int, mailgunDomain string) *AuthService {
	return &AuthService{
		db:                 db,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
		mailgunDomain:      mailgunDomain,
	}
}

// RegisterUser creates a new user with hashed password
func (s *AuthService) RegisterUser(email, password string) (*models.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with generated inbox email
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	// Create user first to get the ID
	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate inbox email using user ID and Mailgun domain
	user.InboxEmail = fmt.Sprintf("%s@%s", user.ID.String(), s.mailgunDomain)

	// Update user with inbox email
	if err := s.db.Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to set inbox email: %w", err)
	}

	return user, nil
}

// AuthenticateUser validates credentials and returns user
func (s *AuthService) AuthenticateUser(email, password string) (*models.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// GenerateToken creates a JWT token for the user
func (s *AuthService) GenerateToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(s.jwtExpirationHours))

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	return userID, nil
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Preferences.Make").
		Preload("Preferences.Model").
		Preload("Preferences.Trim").
		Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}
