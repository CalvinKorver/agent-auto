package services

import (
	"errors"
	"fmt"
	"time"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageService struct {
	db            *gorm.DB
	claudeService *ClaudeService
}

func NewMessageService(db *gorm.DB, claudeService *ClaudeService) *MessageService {
	return &MessageService{
		db:            db,
		claudeService: claudeService,
	}
}

// GetThreadMessages retrieves all messages for a thread
func (s *MessageService) GetThreadMessages(threadID, userID uuid.UUID, limit, offset int) ([]models.Message, int64, error) {
	// Verify thread belongs to user
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("thread not found")
		}
		return nil, 0, fmt.Errorf("database error: %w", err)
	}

	// Get total count
	var total int64
	if err := s.db.Model(&models.Message{}).Where("thread_id = ?", threadID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// Get messages with pagination
	var messages []models.Message
	query := s.db.Where("thread_id = ?", threadID).Order("timestamp ASC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	return messages, total, nil
}

// CreateUserMessage creates a user message and generates an agent response
func (s *MessageService) CreateUserMessage(threadID, userID uuid.UUID, content string) (*models.Message, *models.Message, error) {
	// Validate input
	if content == "" {
		return nil, nil, errors.New("message content is required")
	}

	// Verify thread belongs to user and get user preferences
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("thread not found")
		}
		return nil, nil, fmt.Errorf("database error: %w", err)
	}

	// Get user preferences for context with relationships
	var prefs models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).
		Preload("Make").
		Preload("Model").
		First(&prefs).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Get recent message history for context
	var recentMessages []models.Message
	s.db.Where("thread_id = ?", threadID).Order("timestamp DESC").Limit(10).Find(&recentMessages)

	// Reverse to chronological order
	for i, j := 0, len(recentMessages)-1; i < j; i, j = i+1, j-1 {
		recentMessages[i], recentMessages[j] = recentMessages[j], recentMessages[i]
	}

	// Get tracked offers from all user's threads for competitive context
	var trackedOffers []models.TrackedOffer
	s.db.Joins("JOIN threads ON threads.id = tracked_offers.thread_id").
		Where("threads.user_id = ?", userID).
		Order("tracked_offers.tracked_at DESC").
		Limit(20).
		Preload("Thread").
		Find(&trackedOffers)

	// Extract make and model names from relationships
	makeName := ""
	modelName := ""
	if prefs.Make != nil {
		makeName = prefs.Make.Name
	}
	if prefs.Model != nil {
		modelName = prefs.Model.Name
	}

	// Log context information
	fmt.Printf("\n========== CLAUDE CONTEXT DEBUG ==========\n")
	fmt.Printf("User Message: %s\n", content)
	fmt.Printf("User Preferences: %d %s %s\n", prefs.Year, makeName, modelName)
	fmt.Printf("Seller Name: %s\n", thread.SellerName)
	fmt.Printf("Message History Count: %d\n", len(recentMessages))
	if len(recentMessages) > 0 {
		fmt.Printf("Message History:\n")
		for i, msg := range recentMessages {
			fmt.Printf("  [%d] %s: %s\n", i+1, msg.Sender, msg.Content)
		}
	} else {
		fmt.Printf("Message History: (none - this is the first message)\n")
	}
	fmt.Printf("Tracked Offers Count: %d\n", len(trackedOffers))
	if len(trackedOffers) > 0 {
		fmt.Printf("Tracked Offers:\n")
		for i, offer := range trackedOffers {
			sellerName := "Unknown"
			if offer.Thread != nil {
				sellerName = offer.Thread.SellerName
			}
			fmt.Printf("  [%d] %s: %s\n", i+1, sellerName, offer.OfferText)
		}
	}
	fmt.Printf("==========================================\n\n")

	// Create user message
	userMessage := &models.Message{
		UserID:    userID,
		ThreadID:  &threadID,
		Sender:    models.SenderTypeUser,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Generate agent response using Claude
	agentContent, err := s.claudeService.GenerateNegotiationResponse(
		content,
		prefs.Year,
		makeName,
		modelName,
		thread.SellerName,
		recentMessages,
		trackedOffers,
	)
	if err != nil {
		// Still save user message even if agent fails
		if err := s.db.Create(userMessage).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to create user message: %w", err)
		}
		return userMessage, nil, fmt.Errorf("failed to generate agent response: %w", err)
	}

	// Create agent message
	agentMessage := &models.Message{
		UserID:    userID,
		ThreadID:  &threadID,
		Sender:    models.SenderTypeAgent,
		Content:   agentContent,
		Timestamp: time.Now(),
	}

	// Save both messages in a transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(userMessage).Error; err != nil {
			return fmt.Errorf("failed to create user message: %w", err)
		}

		if err := tx.Create(agentMessage).Error; err != nil {
			return fmt.Errorf("failed to create agent message: %w", err)
		}

		// Update thread message count and last message time
		now := time.Now()
		if err := tx.Model(&models.Thread{}).Where("id = ?", threadID).Updates(map[string]interface{}{
			"message_count":   gorm.Expr("message_count + ?", 2),
			"last_message_at": now,
		}).Error; err != nil {
			return fmt.Errorf("failed to update thread: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return userMessage, agentMessage, nil
}

// CreateSellerMessage creates a seller message (for manual testing)
func (s *MessageService) CreateSellerMessage(threadID, userID uuid.UUID, content string) (*models.Message, error) {
	// Validate input
	if content == "" {
		return nil, errors.New("message content is required")
	}

	// Verify thread belongs to user
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("thread not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create seller message
	sellerMessage := &models.Message{
		UserID:    userID,
		ThreadID:  &threadID,
		Sender:    models.SenderTypeSeller,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Save message and update thread
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sellerMessage).Error; err != nil {
			return fmt.Errorf("failed to create seller message: %w", err)
		}

		// Update thread message count and last message time
		now := time.Now()
		if err := tx.Model(&models.Thread{}).Where("id = ?", threadID).Updates(map[string]interface{}{
			"message_count":   gorm.Expr("message_count + ?", 1),
			"last_message_at": now,
		}).Error; err != nil {
			return fmt.Errorf("failed to update thread: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sellerMessage, nil
}

// GetInboxMessages retrieves messages with null thread_id (inbox messages) that are not deleted
func (s *MessageService) GetInboxMessages(userID uuid.UUID, limit, offset int) ([]models.Message, int64, error) {
	// Get total count of inbox messages (excluding deleted)
	var total int64
	if err := s.db.Model(&models.Message{}).Where("user_id = ? AND thread_id IS NULL AND deleted_at IS NULL", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count inbox messages: %w", err)
	}

	// Get inbox messages with pagination, ordered by timestamp descending (newest first), excluding deleted
	var messages []models.Message
	query := s.db.Where("user_id = ? AND thread_id IS NULL AND deleted_at IS NULL", userID).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve inbox messages: %w", err)
	}

	return messages, total, nil
}

// AssignInboxMessageToThread assigns an inbox message to a thread
func (s *MessageService) AssignInboxMessageToThread(messageID, threadID, userID uuid.UUID) error {
	// Verify the thread exists and belongs to the user
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("thread not found")
		}
		return fmt.Errorf("failed to verify thread: %w", err)
	}

	// Verify the message exists, belongs to the user, and is an inbox message (thread_id is NULL)
	var message models.Message
	if err := s.db.Where("id = ? AND user_id = ? AND thread_id IS NULL", messageID, userID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("inbox message not found")
		}
		return fmt.Errorf("failed to verify message: %w", err)
	}

	// Assign the message to the thread
	message.ThreadID = &threadID
	if err := s.db.Save(&message).Error; err != nil {
		return fmt.Errorf("failed to assign message to thread: %w", err)
	}

	return nil
}

// ArchiveInboxMessage soft deletes an inbox message
func (s *MessageService) ArchiveInboxMessage(messageID, userID uuid.UUID) error {
	// Verify the message exists, belongs to the user, and is an inbox message (thread_id is NULL)
	var message models.Message
	if err := s.db.Where("id = ? AND user_id = ? AND thread_id IS NULL AND deleted_at IS NULL", messageID, userID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("inbox message not found")
		}
		return fmt.Errorf("failed to verify message: %w", err)
	}

	// Soft delete the message by setting deleted_at
	now := time.Now()
	message.DeletedAt = &now
	if err := s.db.Save(&message).Error; err != nil {
		return fmt.Errorf("failed to archive message: %w", err)
	}

	return nil
}
