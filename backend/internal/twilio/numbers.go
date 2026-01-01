package twilio

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	messagingApi "github.com/twilio/twilio-go/rest/messaging/v1"
)

// SearchAvailableNumbers searches for available phone numbers by area code
func SearchAvailableNumbers(client *twilio.RestClient, areaCode string, limit int) ([]twilioApi.ApiV2010AvailablePhoneNumberLocal, error) {
	params := &twilioApi.ListAvailablePhoneNumberLocalParams{}

	// Convert area code string to int
	areaCodeInt, err := strconv.Atoi(areaCode)
	if err != nil {
		return nil, fmt.Errorf("invalid area code: %w", err)
	}
	params.SetAreaCode(areaCodeInt)
	params.SetLimit(limit)

	resp, err := client.Api.ListAvailablePhoneNumberLocal("US", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search available numbers: %w", err)
	}

	return resp, nil
}

// PurchaseNumber purchases a phone number and assigns it to a messaging service
// This is a two-step process:
// 1. Purchase the number via IncomingPhoneNumbers API
// 2. Assign it to Messaging Service via Messaging Service Phone Number API (required for A2P 10DLC compliance)
func PurchaseNumber(client *twilio.RestClient, phoneNumber, messagingServiceSID string) (string, error) {
	// Step 1: Purchase the number
	params := &twilioApi.CreateIncomingPhoneNumberParams{}
	params.SetPhoneNumber(phoneNumber)

	resp, err := client.Api.CreateIncomingPhoneNumber(params)
	if err != nil {
		return "", fmt.Errorf("failed to purchase number: %w", err)
	}

	numberSID := *resp.Sid

	// Step 2: Assign the number to Messaging Service (required for A2P 10DLC compliance)
	if messagingServiceSID != "" {
		if err := AssignNumberToService(client, numberSID, messagingServiceSID); err != nil {
			// If assignment fails, we still have the number purchased
			// Log the error but return the number SID so it can be retried later
			// In production, you might want to queue this for retry
			return numberSID, fmt.Errorf("number purchased but failed to assign to messaging service: %w", err)
		}
	}

	return numberSID, nil
}

// AssignNumberToService assigns a purchased phone number to a Messaging Service
// This is required for A2P 10DLC compliance - all numbers must be in a Messaging Service
// that is linked to a registered A2P Campaign
func AssignNumberToService(client *twilio.RestClient, numberSID, serviceSID string) error {
	params := &messagingApi.CreatePhoneNumberParams{}
	params.SetPhoneNumberSid(numberSID)

	_, err := client.MessagingV1.CreatePhoneNumber(serviceSID, params)
	if err != nil {
		return fmt.Errorf("failed to link number to messaging service: %w", err)
	}
	return nil
}

// FormatPhoneNumber formats a phone number to E.164 format
func FormatPhoneNumber(phoneNumber string) string {
	// Remove all non-digit characters
	cleaned := strings.ReplaceAll(phoneNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// If it doesn't start with +, add +1 for US numbers
	if !strings.HasPrefix(cleaned, "+") {
		if strings.HasPrefix(cleaned, "1") {
			cleaned = "+" + cleaned
		} else {
			cleaned = "+1" + cleaned
		}
	}

	return cleaned
}

// ExtractAreaCode extracts area code from a phone number
func ExtractAreaCode(phoneNumber string) string {
	cleaned := FormatPhoneNumber(phoneNumber)
	// E.164 format: +1XXXXXXXXXX, area code is positions 2-4
	if len(cleaned) >= 5 && strings.HasPrefix(cleaned, "+1") {
		return cleaned[2:5]
	}
	return ""
}
