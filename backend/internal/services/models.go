package services

import (
	"errors"
	"fmt"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelsService handles fetching vehicle makes/models from database
type ModelsService struct {
	db *gorm.DB
}

// NewModelsService creates a new models service
func NewModelsService(db *gorm.DB) *ModelsService {
	return &ModelsService{
		db: db,
	}
}

// GetModels returns all makes and their models in the format expected by the frontend
func (s *ModelsService) GetModels() (map[string][]string, error) {
	var makes []models.Make
	if err := s.db.Preload("Models").Order("name ASC").Find(&makes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch makes: %w", err)
	}

	result := make(map[string][]string)
	for _, makeRecord := range makes {
		modelNames := make([]string, 0, len(makeRecord.Models))
		for _, model := range makeRecord.Models {
			modelNames = append(modelNames, model.Name)
		}
		result[makeRecord.Name] = modelNames
	}

	return result, nil
}

// GetMakes returns a list of all makes
func (s *ModelsService) GetMakes() ([]string, error) {
	var makes []models.Make
	if err := s.db.Order("name ASC").Find(&makes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch makes: %w", err)
	}

	makeNames := make([]string, 0, len(makes))
	for _, make := range makes {
		makeNames = append(makeNames, make.Name)
	}

	return makeNames, nil
}

// GetModelsForMake returns the models for a specific make
func (s *ModelsService) GetModelsForMake(makeName string) ([]string, error) {
	var makeRecord models.Make
	if err := s.db.Where("name = ?", makeName).Preload("Models").First(&makeRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("make not found: %s", makeName)
		}
		return nil, fmt.Errorf("failed to fetch make: %w", err)
	}

	modelNames := make([]string, 0, len(makeRecord.Models))
	for _, model := range makeRecord.Models {
		modelNames = append(modelNames, model.Name)
	}

	return modelNames, nil
}

// GetMakeByName looks up a make by name and returns the Make model
func (s *ModelsService) GetMakeByName(name string) (*models.Make, error) {
	var makeRecord models.Make
	if err := s.db.Where("name = ?", name).First(&makeRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("make not found: %s", name)
		}
		return nil, fmt.Errorf("failed to lookup make: %w", err)
	}
	return &makeRecord, nil
}

// GetModelByName looks up a model by make ID and model name
func (s *ModelsService) GetModelByName(makeID uuid.UUID, name string) (*models.Model, error) {
	var model models.Model
	if err := s.db.Where("make_id = ? AND name = ?", makeID, name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("model not found: %s for make %s", name, makeID)
		}
		return nil, fmt.Errorf("failed to lookup model: %w", err)
	}
	return &model, nil
}

// GetTrimsForModel returns all trims for a specific model and year
func (s *ModelsService) GetTrimsForModel(modelID uuid.UUID, year int) ([]models.VehicleTrim, error) {
	var trims []models.VehicleTrim
	if err := s.db.Where("model_id = ? AND year = ?", modelID, year).
		Order("trim_name ASC").
		Find(&trims).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch trims: %w", err)
	}
	return trims, nil
}
