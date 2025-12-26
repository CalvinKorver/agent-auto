package gmail

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// CreateGmailService creates an authenticated Gmail API service for a user
func CreateGmailService(userID uuid.UUID, tokenManager *TokenManager, oauthConfig *oauth2.Config) (*gmail.Service, error) {
	// Get and refresh token if needed
	token, err := tokenManager.RefreshTokenIfNeeded(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create authenticated HTTP client
	client := oauthConfig.Client(context.Background(), token)

	// Create Gmail service
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %w", err)
	}

	return service, nil
}
