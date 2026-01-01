package twilio

import (
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendSMS sends an SMS message via Twilio Messaging Service
func SendSMS(client *twilio.RestClient, from, to, body, messagingServiceSID string) (string, error) {
	params := &twilioApi.CreateMessageParams{}
	
	// Use Messaging Service for A2P compliance
	if messagingServiceSID != "" {
		params.SetMessagingServiceSid(messagingServiceSID)
		// To is still required
		params.SetTo(to)
	} else {
		// Fallback to direct from number (not recommended for production)
		params.SetFrom(from)
		params.SetTo(to)
	}
	
	params.SetBody(body)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return "", fmt.Errorf("failed to send SMS: %w", err)
	}

	return *resp.Sid, nil
}

