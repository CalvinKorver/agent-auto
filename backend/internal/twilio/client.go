package twilio

import (
	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
)

// Client wraps the Twilio client
type Client struct {
	client                *twilio.RestClient
	messagingServiceSID   string
}

// NewClient creates a new Twilio client
func NewClient(accountSID, authToken, messagingServiceSID string) *Client {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &Client{
		client:              client,
		messagingServiceSID: messagingServiceSID,
	}
}

// GetClient returns the underlying Twilio REST client
func (c *Client) GetClient() *twilio.RestClient {
	return c.client
}

// GetMessagingServiceSID returns the messaging service SID
func (c *Client) GetMessagingServiceSID() string {
	return c.messagingServiceSID
}

// ValidateWebhookSignature validates the X-Twilio-Signature header
func ValidateWebhookSignature(url string, params map[string]string, signature, authToken string) bool {
	validator := twilioClient.NewRequestValidator(authToken)
	return validator.Validate(url, params, signature)
}

