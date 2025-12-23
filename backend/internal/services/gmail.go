package services

import (
	"fmt"

	"carbuyer/internal/gmail"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// GmailService handles Gmail OAuth and sending operations
type GmailService struct {
	tokenManager *gmail.TokenManager
	oauthConfig  *oauth2.Config
}

// NewGmailService creates a new Gmail service
func NewGmailService(db *gorm.DB, clientID, clientSecret, redirectURL, encryptionKey string) (*GmailService, error) {
	// Create OAuth config
	oauthConfig := gmail.CreateOAuthConfig(clientID, clientSecret, redirectURL)

	// Create token manager
	tokenManager, err := gmail.NewTokenManager(db, encryptionKey, oauthConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	return &GmailService{
		tokenManager: tokenManager,
		oauthConfig:  oauthConfig,
	}, nil
}

// GetAuthURL generates OAuth authorization URL
func (s *GmailService) GetAuthURL(state string) string {
	return gmail.GetAuthURL(s.oauthConfig, state)
}

// ExchangeCodeForToken exchanges authorization code for tokens
func (s *GmailService) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	return gmail.ExchangeCodeForToken(s.oauthConfig, code)
}

// StoreToken stores encrypted token in database
func (s *GmailService) StoreToken(userID uuid.UUID, token *oauth2.Token, gmailEmail string) error {
	return s.tokenManager.StoreToken(userID, token, gmailEmail)
}

// IsConnected checks if user has Gmail connected
func (s *GmailService) IsConnected(userID uuid.UUID) bool {
	return s.tokenManager.IsConnected(userID)
}

// GetGmailEmail returns connected Gmail email
func (s *GmailService) GetGmailEmail(userID uuid.UUID) (string, error) {
	return s.tokenManager.GetGmailEmail(userID)
}

// DisconnectGmail revokes and deletes user's Gmail connection
func (s *GmailService) DisconnectGmail(userID uuid.UUID) error {
	return s.tokenManager.RevokeToken(userID)
}

// SendReply sends an email reply via user's Gmail
func (s *GmailService) SendReply(userID uuid.UUID, to, subject, htmlBody, externalMessageID string) error {
	// Create Gmail service for this user
	service, err := gmail.CreateGmailService(userID, s.tokenManager, s.oauthConfig)
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	// Send reply
	return gmail.SendReply(service, to, subject, htmlBody, externalMessageID)
}
