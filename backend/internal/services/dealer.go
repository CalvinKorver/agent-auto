package services

import (
	"errors"
	"fmt"
	"log"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DealerService struct {
	db            *gorm.DB
	claudeService *ClaudeService
}

func NewDealerService(db *gorm.DB, claudeService *ClaudeService) *DealerService {
	return &DealerService{
		db:            db,
		claudeService: claudeService,
	}
}

// FetchDealersForZipCode fetches dealers using Claude based on zip code and vehicle info
func (s *DealerService) FetchDealersForZipCode(zipCode string, make string, model string, year int) ([]DealerInfo, error) {
	if zipCode == "" {
		return nil, errors.New("zip code is required")
	}
	if make == "" || model == "" {
		return nil, errors.New("make and model are required")
	}

	return s.claudeService.FetchNearbyDealers(zipCode, make, model, year)
}

// SaveDealersForPreferences saves dealers to the database for a given preference
func (s *DealerService) SaveDealersForPreferences(preferenceID uuid.UUID, dealers []DealerInfo) error {
	if len(dealers) == 0 {
		return nil // Nothing to save
	}

	// Delete existing dealers for this preference
	if err := s.DeleteDealersForPreferences(preferenceID); err != nil {
		return fmt.Errorf("failed to delete existing dealers: %w", err)
	}

	// Create dealer records (contacted defaults to false)
	for _, dealerInfo := range dealers {
		dealer := &models.Dealer{
			Name:             dealerInfo.Name,
			Location:         dealerInfo.Location,
			Email:            dealerInfo.Email,
			Phone:            dealerInfo.Phone,
			Website:          dealerInfo.Website,
			Distance:         dealerInfo.Distance,
			UserPreferenceID: preferenceID,
		}

		if err := s.db.Create(dealer).Error; err != nil {
			return fmt.Errorf("failed to create dealer %s: %w", dealerInfo.Name, err)
		}
	}

	log.Printf("Saved %d dealers for preference %s", len(dealers), preferenceID)
	return nil
}

// GetDealersForPreferences retrieves dealers for a given preference
func (s *DealerService) GetDealersForPreferences(preferenceID uuid.UUID) ([]models.Dealer, error) {
	var dealers []models.Dealer
	if err := s.db.Where("user_preference_id = ?", preferenceID).
		Order("distance ASC"). // Sort by distance, closest first
		Find(&dealers).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve dealers: %w", err)
	}

	return dealers, nil
}

// DeleteDealersForPreferences removes all dealers for a given preference
func (s *DealerService) DeleteDealersForPreferences(preferenceID uuid.UUID) error {
	result := s.db.Where("user_preference_id = ?", preferenceID).Delete(&models.Dealer{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete dealers: %w", result.Error)
	}
	return nil
}

// UpdateDealersContacted updates the contacted status for multiple dealers
func (s *DealerService) UpdateDealersContacted(dealerIDs []uuid.UUID, contacted bool) error {
	if len(dealerIDs) == 0 {
		return nil // Nothing to update
	}

	result := s.db.Model(&models.Dealer{}).
		Where("id IN ?", dealerIDs).
		Update("contacted", contacted)

	if result.Error != nil {
		return fmt.Errorf("failed to update dealers: %w", result.Error)
	}

	log.Printf("Updated contacted status for %d dealers to %v (rows affected: %d)", len(dealerIDs), contacted, result.RowsAffected)
	return nil
}
