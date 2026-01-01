package services

import (
	"fmt"
	"time"

	"carbuyer/internal/db/models"
	"carbuyer/internal/twilio"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SMSService struct {
	db            *gorm.DB
	twilioClient  *twilio.Client
}

func NewSMSService(db *gorm.DB, twilioClient *twilio.Client) *SMSService {
	return &SMSService{
		db:           db,
		twilioClient: twilioClient,
	}
}

// ProcessInboundSMS creates an inbox message from an incoming SMS
func (s *SMSService) ProcessInboundSMS(userID uuid.UUID, from, body, messageSID string) (*models.Message, error) {
	// Get phone message type
	var phoneType models.MessageType
	if err := s.db.Where("type = ?", models.MessageTypePhone).First(&phoneType).Error; err != nil {
		return nil, fmt.Errorf("failed to get phone message type: %w", err)
	}

	// Check for duplicate based on external_message_id (Twilio message SID)
	if messageSID != "" {
		var existingMessage models.Message
		if err := s.db.Where("external_message_id = ?", messageSID).First(&existingMessage).Error; err == nil {
			// Message already exists, return it
			return &existingMessage, nil
		}
	}

	// Create inbox message (thread_id = nil)
	message := &models.Message{
		UserID:            userID,
		ThreadID:          nil, // Unassigned - goes to inbox
		MessageTypeID:     &phoneType.ID,
		Sender:            models.SenderTypeSeller,
		Content:           body,
		Timestamp:      time.Now(),
		SenderPhone:       from,
		ExternalMessageID: messageSID,
	}

	// Save message to database
	if err := s.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create inbox message: %w", err)
	}

	return message, nil
}

// SendSMS sends an SMS reply via Twilio
func (s *SMSService) SendSMS(userID uuid.UUID, threadID uuid.UUID, content string) error {
	// Get user's Twilio phone number
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.PhoneNumber == "" {
		return fmt.Errorf("user does not have a Twilio phone number allocated")
	}

	// Get thread to find dealer phone number
	var thread models.Thread
	if err := s.db.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		return fmt.Errorf("thread not found: %w", err)
	}

	if thread.Phone == "" {
		return fmt.Errorf("thread does not have a phone number assigned")
	}

	// Send SMS via Twilio
	client := s.twilioClient.GetClient()
	messagingServiceSID := s.twilioClient.GetMessagingServiceSID()
	
	_, err := twilio.SendSMS(
		client,
		user.PhoneNumber, // from
		thread.Phone,     // to (dealer)
		content,
		messagingServiceSID,
	)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// Get phone message type
	var phoneType models.MessageType
	if err := s.db.Where("type = ?", models.MessageTypePhone).First(&phoneType).Error; err != nil {
		return fmt.Errorf("failed to get phone message type: %w", err)
	}

	// Save outbound message to database
	message := &models.Message{
		UserID:        userID,
		ThreadID:      &threadID,
		MessageTypeID: &phoneType.ID,
		Sender:        models.SenderTypeAgent,
		Content:       content,
		Timestamp:      time.Now(),
		SenderPhone:   thread.Phone,
		SentViaSMS:    true,
	}

	if err := s.db.Create(message).Error; err != nil {
		return fmt.Errorf("failed to save outbound message: %w", err)
	}

	// Update thread message count and last message time
	now := time.Now()
	if err := s.db.Model(&models.Thread{}).Where("id = ?", threadID).Updates(map[string]interface{}{
		"message_count":   gorm.Expr("message_count + ?", 1),
		"last_message_at": now,
	}).Error; err != nil {
		return fmt.Errorf("failed to update thread: %w", err)
	}

	return nil
}

// AllocatePhoneNumber allocates a Twilio phone number for a user
func (s *SMSService) AllocatePhoneNumber(userID uuid.UUID) error {
	// Check if user already has a phone number
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.PhoneNumber != "" {
		// User already has a phone number
		return nil
	}

	// Try to get area code from user's zip code if available
	areaCode := ""
	if user.ZipCode != "" {
		// For now, we'll use a default area code or try to extract from zip
		// In a real implementation, you might want to map zip codes to area codes
		// For simplicity, we'll search for numbers in a common area code
		areaCode = "415" // Default to San Francisco area code
	} else {
		areaCode = "415" // Default area code
	}

	// Search for available numbers
	client := s.twilioClient.GetClient()
	availableNumbers, err := twilio.SearchAvailableNumbers(client, areaCode, 1)
	if err != nil {
		return fmt.Errorf("failed to search available numbers: %w", err)
	}
	if len(availableNumbers) == 0 {
		return fmt.Errorf("no available phone numbers found for area code %s", areaCode)
	}

	// Purchase the first available number
	phoneNumber := *availableNumbers[0].PhoneNumber
	messagingServiceSID := s.twilioClient.GetMessagingServiceSID()
	
	twilioSID, err := twilio.PurchaseNumber(client, phoneNumber, messagingServiceSID)
	if err != nil {
		return fmt.Errorf("failed to purchase number: %w", err)
	}
	// Update user with phone number and Twilio SID
	user.PhoneNumber = phoneNumber
	user.TwilioSID = twilioSID
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to save phone number to user: %w", err)
	}
	return nil
}

