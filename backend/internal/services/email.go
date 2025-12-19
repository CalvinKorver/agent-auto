package services

import (
	"fmt"
	"strings"
	"time"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailService struct {
	db            *gorm.DB
	mailgunAPIKey string
	mailgunDomain string
}

func NewEmailService(db *gorm.DB, mailgunAPIKey, mailgunDomain string) *EmailService {
	return &EmailService{
		db:            db,
		mailgunAPIKey: mailgunAPIKey,
		mailgunDomain: mailgunDomain,
	}
}

// ProcessInboundEmail creates an inbox message from a forwarded email
func (s *EmailService) ProcessInboundEmail(userID uuid.UUID, from, subject, body, messageID string) (*models.Message, error) {
	// Parse and clean email body
	cleanedBody := s.cleanEmailBody(body)

	// Create inbox message (thread_id = nil)
	message := &models.Message{
		UserID:            userID,
		ThreadID:          nil, // Unassigned - goes to inbox
		Sender:            models.SenderTypeSeller,
		Content:           cleanedBody,
		Timestamp:         time.Now(),
		SenderEmail:       from,
		ExternalMessageID: messageID,
		Subject:           subject,
		SentViaEmail:      true,
	}

	// Check for duplicate based on external_message_id
	if messageID != "" {
		var existingMessage models.Message
		if err := s.db.Where("external_message_id = ?", messageID).First(&existingMessage).Error; err == nil {
			// Message already exists, return it
			return &existingMessage, nil
		}
	}

	// Save message to database
	if err := s.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create inbox message: %w", err)
	}

	return message, nil
}

// cleanEmailBody removes quoted text and cleans up email formatting
func (s *EmailService) cleanEmailBody(body string) string {
	// Remove leading/trailing whitespace
	body = strings.TrimSpace(body)

	// Split into lines for processing
	lines := strings.Split(body, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Stop at common reply separators
		if strings.HasPrefix(trimmedLine, "On ") && strings.Contains(trimmedLine, "wrote:") {
			break
		}
		if strings.HasPrefix(trimmedLine, ">") {
			continue // Skip quoted lines
		}
		if trimmedLine == "" && len(cleanedLines) > 0 && cleanedLines[len(cleanedLines)-1] == "" {
			continue // Skip multiple blank lines
		}

		cleanedLines = append(cleanedLines, line)
	}

	cleaned := strings.Join(cleanedLines, "\n")
	return strings.TrimSpace(cleaned)
}
