package gmail

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// CreateOAuthConfig creates OAuth2 configuration for Gmail API
func CreateOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			gmail.GmailComposeScope, // https://www.googleapis.com/auth/gmail.compose - Manage drafts and send emails
		},
		Endpoint: google.Endpoint,
	}
}

// GenerateStateToken generates a random state token for CSRF protection
func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthURL generates the OAuth consent screen URL
func GetAuthURL(config *oauth2.Config, state string) string {
	return config.AuthCodeURL(state,
		oauth2.AccessTypeOffline, // Get refresh token
		oauth2.ApprovalForce,      // Force approval prompt (ensures we get refresh token)
	)
}

// ExchangeCodeForToken exchanges authorization code for tokens
func ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Verify we got a refresh token
	if token.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token received (user may have already authorized)")
	}

	return token, nil
}
