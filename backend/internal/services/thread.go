package services

import (
	"errors"
	"fmt"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadService struct {
	db *gorm.DB
}

func NewThreadService(db *gorm.DB) *ThreadService {
	return &ThreadService{
		db: db,
	}
}

// CreateThread creates a new thread for a user
func (s *ThreadService) CreateThread(userID uuid.UUID, sellerName string, sellerType models.SellerType) (*models.Thread, error) {
	// Validate input
	if sellerName == "" {
		return nil, errors.New("seller name is required")
	}

	if sellerType != models.SellerTypePrivate && sellerType != models.SellerTypeDealership && sellerType != models.SellerTypeOther {
		return nil, errors.New("invalid seller type")
	}

	// Create thread
	thread := &models.Thread{
		UserID:     userID,
		SellerName: sellerName,
		SellerType: sellerType,
	}

	if err := s.db.Create(thread).Error; err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}

	return thread, nil
}

// GetUserThreads retrieves all threads for a user
func (s *ThreadService) GetUserThreads(userID uuid.UUID) ([]models.Thread, error) {
	var threads []models.Thread
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&threads).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve threads: %w", err)
	}

	return threads, nil
}

// GetThreadByID retrieves a specific thread
func (s *ThreadService) GetThreadByID(threadID, userID uuid.UUID) (*models.Thread, error) {
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("thread not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &thread, nil
}

// DeleteThread deletes a thread (soft delete in future)
func (s *ThreadService) DeleteThread(threadID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", threadID, userID).Delete(&models.Thread{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete thread: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("thread not found")
	}

	return nil
}
