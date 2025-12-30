package services

import (
	"errors"
	"fmt"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreferencesService struct {
	db            *gorm.DB
	modelsService *ModelsService
}

func NewPreferencesService(db *gorm.DB, modelsService *ModelsService) *PreferencesService {
	return &PreferencesService{
		db:            db,
		modelsService: modelsService,
	}
}

// GetUserPreferences retrieves preferences for a user with relationships loaded
func (s *PreferencesService) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	var prefs models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).
		Preload("Make").
		Preload("Model").
		Preload("Trim").
		First(&prefs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("preferences not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &prefs, nil
}

// CreateUserPreferences creates preferences for a user
func (s *PreferencesService) CreateUserPreferences(userID uuid.UUID, year int, makeName, modelName string, trimID *uuid.UUID) (*models.UserPreferences, error) {
	// Validate input
	if year < 1900 || year > 2100 {
		return nil, errors.New("invalid year")
	}
	if makeName == "" || modelName == "" {
		return nil, errors.New("make and model are required")
	}

	// Check if preferences already exist
	var existing models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return nil, errors.New("preferences already exist for this user")
	}

	// Lookup make by name
	make, err := s.modelsService.GetMakeByName(makeName)
	if err != nil {
		return nil, fmt.Errorf("make not found: %s", makeName)
	}

	// Lookup model by make ID and name
	model, err := s.modelsService.GetModelByName(make.ID, modelName)
	if err != nil {
		return nil, fmt.Errorf("model not found: %s for make %s", modelName, makeName)
	}

	// Validate trim if provided
	if trimID != nil {
		trims, err := s.modelsService.GetTrimsForModel(model.ID, year)
		if err != nil {
			return nil, fmt.Errorf("failed to validate trim: %w", err)
		}
		trimFound := false
		for _, trim := range trims {
			if trim.ID == *trimID {
				trimFound = true
				break
			}
		}
		if !trimFound {
			return nil, errors.New("trim not found for the specified model and year")
		}
	}

	// Create preferences with foreign keys
	prefs := &models.UserPreferences{
		UserID:  userID,
		MakeID:  make.ID,
		ModelID: model.ID,
		Year:    year,
		TrimID:  trimID,
	}

	if err := s.db.Create(prefs).Error; err != nil {
		return nil, fmt.Errorf("failed to create preferences: %w", err)
	}

	// Load relationships for response
	if err := s.db.Preload("Make").Preload("Model").Preload("Trim").First(prefs, prefs.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load relationships: %w", err)
	}

	return prefs, nil
}

// UpdateUserPreferences updates existing preferences
func (s *PreferencesService) UpdateUserPreferences(userID uuid.UUID, year int, makeName, modelName string, trimID *uuid.UUID) (*models.UserPreferences, error) {
	// Validate input
	if year < 1900 || year > 2100 {
		return nil, errors.New("invalid year")
	}
	if makeName == "" || modelName == "" {
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

	// Lookup make by name
	make, err := s.modelsService.GetMakeByName(makeName)
	if err != nil {
		return nil, fmt.Errorf("make not found: %s", makeName)
	}

	// Lookup model by make ID and name
	model, err := s.modelsService.GetModelByName(make.ID, modelName)
	if err != nil {
		return nil, fmt.Errorf("model not found: %s for make %s", modelName, makeName)
	}

	// Validate trim if provided
	if trimID != nil {
		trims, err := s.modelsService.GetTrimsForModel(model.ID, year)
		if err != nil {
			return nil, fmt.Errorf("failed to validate trim: %w", err)
		}
		trimFound := false
		for _, trim := range trims {
			if trim.ID == *trimID {
				trimFound = true
				break
			}
		}
		if !trimFound {
			return nil, errors.New("trim not found for the specified model and year")
		}
	}

	// Update preferences with foreign keys
	prefs.Year = year
	prefs.MakeID = make.ID
	prefs.ModelID = model.ID
	prefs.TrimID = trimID

	if err := s.db.Save(&prefs).Error; err != nil {
		return nil, fmt.Errorf("failed to update preferences: %w", err)
	}

	// Load relationships for response
	if err := s.db.Preload("Make").Preload("Model").Preload("Trim").First(&prefs, prefs.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load relationships: %w", err)
	}

	return &prefs, nil
}
