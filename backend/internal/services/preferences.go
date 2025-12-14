package services

import (
	"errors"
	"fmt"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreferencesService struct {
	db *gorm.DB
}

func NewPreferencesService(db *gorm.DB) *PreferencesService {
	return &PreferencesService{
		db: db,
	}
}

// GetUserPreferences retrieves preferences for a user
func (s *PreferencesService) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	var prefs models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).First(&prefs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("preferences not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &prefs, nil
}

// CreateUserPreferences creates preferences for a user
func (s *PreferencesService) CreateUserPreferences(userID uuid.UUID, year int, make, model string) (*models.UserPreferences, error) {
	// Validate input
	if year < 1900 || year > 2100 {
		return nil, errors.New("invalid year")
	}
	if make == "" || model == "" {
		return nil, errors.New("make and model are required")
	}

	// Check if preferences already exist
	var existing models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return nil, errors.New("preferences already exist for this user")
	}

	// Create preferences
	prefs := &models.UserPreferences{
		UserID: userID,
		Year:   year,
		Make:   make,
		Model:  model,
	}

	if err := s.db.Create(prefs).Error; err != nil {
		return nil, fmt.Errorf("failed to create preferences: %w", err)
	}

	return prefs, nil
}

// UpdateUserPreferences updates existing preferences
func (s *PreferencesService) UpdateUserPreferences(userID uuid.UUID, year int, make, model string) (*models.UserPreferences, error) {
	// Validate input
	if year < 1900 || year > 2100 {
		return nil, errors.New("invalid year")
	}
	if make == "" || model == "" {
		return nil, errors.New("make and model are required")
	}

	// Find existing preferences
	var prefs models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).First(&prefs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("preferences not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Update preferences
	prefs.Year = year
	prefs.Make = make
	prefs.Model = model

	if err := s.db.Save(&prefs).Error; err != nil {
		return nil, fmt.Errorf("failed to update preferences: %w", err)
	}

	return &prefs, nil
}
