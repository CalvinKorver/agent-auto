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

	// Get user preferences for context
	var prefs models.UserPreferences
	if err := s.db.Where("user_id = ?", userID).First(&prefs).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Get recent message history for context
	var recentMessages []models.Message
	s.db.Where("thread_id = ?", threadID).Order("timestamp DESC").Limit(10).Find(&recentMessages)

	// Reverse to chronological order
	for i, j := 0, len(recentMessages)-1; i < j; i, j = i+1, j-1 {
		recentMessages[i], recentMessages[j] = recentMessages[j], recentMessages[i]
	}

	// Create user message
	userMessage := &models.Message{
		ThreadID:  threadID,
		Sender:    models.SenderTypeUser,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Generate agent response using Claude
	agentContent, err := s.claudeService.GenerateNegotiationResponse(
		content,
		prefs.Year,
		prefs.Make,
		prefs.Model,
		thread.SellerName,
		recentMessages,
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
		ThreadID:  threadID,
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
		ThreadID:  threadID,
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
