package services

import (
	"errors"
	"fmt"
	"log"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreferencesService struct {
	db            *gorm.DB
	modelsService *ModelsService
	dealerService *DealerService
}

func NewPreferencesService(db *gorm.DB, modelsService *ModelsService, dealerService *DealerService) *PreferencesService {
	return &PreferencesService{
		db:            db,
		modelsService: modelsService,
		dealerService: dealerService,
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
func (s *PreferencesService) CreateUserPreferences(userID uuid.UUID, year int, makeName, modelName string, trimID *uuid.UUID, zipCode string) (*models.UserPreferences, error) {
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

	// Update user's zip code if provided
	if zipCode != "" {
		if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("zip_code", zipCode).Error; err != nil {
			return nil, fmt.Errorf("failed to update user zip code: %w", err)
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

	// Trigger async dealer fetch if zip code is provided
	if zipCode != "" && s.dealerService != nil {
		go func() {
			dealers, err := s.dealerService.FetchDealersForZipCode(zipCode, makeName, modelName, year)
			if err != nil {
				log.Printf("Failed to fetch dealers for zip code %s: %v", zipCode, err)
				return
			}
			if err := s.dealerService.SaveDealersForPreferences(prefs.ID, dealers); err != nil {
				log.Printf("Failed to save dealers for preferences: %v", err)
			}
		}()
	}

	return prefs, nil
}

// UpdateUserPreferences updates existing preferences
func (s *PreferencesService) UpdateUserPreferences(userID uuid.UUID, year int, makeName, modelName string, trimID *uuid.UUID, zipCode *string) (*models.UserPreferences, error) {
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

	// Update user's zip code if provided
	if zipCode != nil && *zipCode != "" {
		if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("zip_code", *zipCode).Error; err != nil {
			return nil, fmt.Errorf("failed to update user zip code: %w", err)
		}
	}

	// Check if preferences changed (make, model, year, or zip code)
	preferencesChanged := prefs.MakeID != make.ID || prefs.ModelID != model.ID || prefs.Year != year
	zipCodeChanged := zipCode != nil && *zipCode != ""

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

	// Trigger async dealer fetch if preferences or zip code changed
	if (preferencesChanged || zipCodeChanged) && s.dealerService != nil {
		// Get current zip code
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			log.Printf("Failed to get user for dealer fetch: %v", err)
		} else if user.ZipCode != "" {
			// Delete old dealers
			if err := s.dealerService.DeleteDealersForPreferences(prefs.ID); err != nil {
				log.Printf("Failed to delete old dealers: %v", err)
			}

			// Fetch new dealers
			go func() {
				dealers, err := s.dealerService.FetchDealersForZipCode(user.ZipCode, makeName, modelName, year)
				if err != nil {
					log.Printf("Failed to fetch dealers for zip code %s: %v", user.ZipCode, err)
					return
				}
				if err := s.dealerService.SaveDealersForPreferences(prefs.ID, dealers); err != nil {
					log.Printf("Failed to save dealers for preferences: %v", err)
				}
			}()
		}
	}

	return &prefs, nil
}
