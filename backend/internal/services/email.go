package services

import (
	"fmt"
	"regexp"
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
	gmailService  *GmailService
}

func NewEmailService(db *gorm.DB, mailgunAPIKey, mailgunDomain string, gmailService *GmailService) *EmailService {
	return &EmailService{
		db:            db,
		mailgunAPIKey: mailgunAPIKey,
		mailgunDomain: mailgunDomain,
		gmailService:  gmailService,
	}
}

// ProcessInboundEmail creates an inbox message from a forwarded email
func (s *EmailService) ProcessInboundEmail(userID uuid.UUID, from, subject, body, messageID string) (*models.Message, error) {
	// Parse and clean email body
	cleanedBody := s.cleanEmailBody(body)

	fmt.Printf("Cleaned body: %s", cleanedBody)

	// Try to extract original sender from forwarded email body
	originalSender := s.extractOriginalSenderFromBody(body)
	senderEmail := from
	if originalSender != "" {
		senderEmail = originalSender
	}

	fmt.Printf("Original sender: %s", originalSender)

	// Create inbox message (thread_id = nil)
	message := &models.Message{
		UserID:            userID,
		ThreadID:          nil, // Unassigned - goes to inbox
		Sender:            models.SenderTypeSeller,
		Content:           cleanedBody,
		Timestamp:         time.Now(),
		SenderEmail:       senderEmail,
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

// extractOriginalSenderFromBody extracts the original sender email from a forwarded email body
// Looks for pattern: "---------- Forwarded message ---------\nFrom: Name <email@domain.com>"
// Returns the email address if found, empty string otherwise
func (s *EmailService) extractOriginalSenderFromBody(body string) string {
	// Pattern to match forwarded message header and From line
	// Case-insensitive, handles variable number of dashes and whitespace variations
	// Matches: "---------- Forwarded message ---------" followed by "From: ..."
	// Uses -+ to match one or more dashes (more flexible than exact count)
	forwardedPattern := regexp.MustCompile(`(?i)-+\s*Forwarded\s+message\s*-+\s*From:\s*(.+)`)
	
	matches := forwardedPattern.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}

	fmt.Printf("Matches: %v", matches)
	
	// Extract the From line content (e.g., "Robyn Langevin <robyndl8uw@gmail.com>" or "email@domain.com")
	fromLine := strings.TrimSpace(matches[1])
	if fromLine == "" {
		return ""
	}

	fmt.Printf("From line: %s", fromLine)
	
	// Extract email address from "Name <email@domain.com>" or "email@domain.com"
	// This pattern matches email addresses with optional angle brackets
	emailPattern := regexp.MustCompile(`<?([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})>?`)
	
	emailMatches := emailPattern.FindStringSubmatch(fromLine)
	if len(emailMatches) < 2 {
		return ""
	}

	fmt.Printf("Email matches: %v", emailMatches)
	
	return strings.ToLower(emailMatches[1]) // Normalize to lowercase
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

// ReplyViaGmail sends threaded reply from user's Gmail
// This is called when user clicks "Send Email" on an AI-drafted response
func (s *EmailService) ReplyViaGmail(userID uuid.UUID, inboxMessageID uuid.UUID, replyContent string) error {
	// 1. Get original inbox message from DB by ID
	

	var message models.Message
	if err := s.db.First(&message, inboxMessageID).Error; err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	// 2. Validate message has email metadata
	if message.ExternalMessageID == "" || message.SenderEmail == "" {
		return fmt.Errorf("message was not received via email")
	}

	fmt.Printf("Sender email: %s", message.SenderEmail)

	// 3. Build reply subject
	replySubject := message.Subject
	if !strings.HasPrefix(replySubject, "Re:") {
		replySubject = "Re: " + replySubject
	}

	// 4. Send reply via Gmail
	return s.gmailService.SendReply(
		userID,
		message.SenderEmail,      // to
		replySubject,             // subject
		replyContent,             // body (from AI draft)
		message.ExternalMessageID, // In-Reply-To header
	)
}
