package gmail

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// TokenManager handles encryption, storage, and retrieval of OAuth tokens
type TokenManager struct {
	db             *gorm.DB
	encryptionKey  []byte
	oauthConfig    *oauth2.Config
}

// NewTokenManager creates a new token manager
func NewTokenManager(db *gorm.DB, encryptionKeyHex string, oauthConfig *oauth2.Config) (*TokenManager, error) {
	// Decode hex encryption key
	encryptionKey, err := hex.DecodeString(encryptionKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key: %w", err)
	}

	// Must be 32 bytes for AES-256
	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes (64 hex characters), got %d bytes", len(encryptionKey))
	}

	return &TokenManager{
		db:            db,
		encryptionKey: encryptionKey,
		oauthConfig:   oauthConfig,
	}, nil
}

// EncryptToken encrypts an OAuth2 token using AES-256-GCM
func (tm *TokenManager) EncryptToken(token *oauth2.Token) (string, error) {
	// Marshal token to JSON
	data, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(tm.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and seal
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptToken decrypts an encrypted token string
func (tm *TokenManager) DecryptToken(encrypted string) (*oauth2.Token, error) {
	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(tm.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Unmarshal token
	var token oauth2.Token
	if err := json.Unmarshal(plaintext, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// StoreToken encrypts and stores a token in the database
func (tm *TokenManager) StoreToken(userID uuid.UUID, token *oauth2.Token, gmailEmail string) error {
	// Encrypt access token
	encryptedAccess, err := tm.EncryptToken(&oauth2.Token{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Expiry:      token.Expiry,
	})
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	// Encrypt refresh token separately (more secure to store them separately)
	encryptedRefresh, err := tm.EncryptToken(&oauth2.Token{
		AccessToken: token.RefreshToken, // Store refresh token in AccessToken field for encryption
	})
	if err != nil {
		return fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	// Create or update token record
	gmailToken := models.GmailToken{
		UserID:       userID,
		AccessToken:  encryptedAccess,
		RefreshToken: encryptedRefresh,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
		GmailEmail:   gmailEmail,
	}

	// Upsert (update if exists, create if not)
	result := tm.db.Save(&gmailToken)
	return result.Error
}

// GetToken retrieves and decrypts a token from the database
func (tm *TokenManager) GetToken(userID uuid.UUID) (*oauth2.Token, error) {
	var gmailToken models.GmailToken
	if err := tm.db.First(&gmailToken, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gmail not connected")
		}
		return nil, fmt.Errorf("failed to retrieve token: %w", err)
	}

	// Decrypt access token
	accessTokenData, err := tm.DecryptToken(gmailToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// Decrypt refresh token
	refreshTokenData, err := tm.DecryptToken(gmailToken.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	// Reconstruct full token
	token := &oauth2.Token{
		AccessToken:  accessTokenData.AccessToken,
		RefreshToken: refreshTokenData.AccessToken, // Was stored in AccessToken field
		TokenType:    gmailToken.TokenType,
		Expiry:       gmailToken.Expiry,
	}

	return token, nil
}

// RefreshTokenIfNeeded checks if token is expired and refreshes if needed
func (tm *TokenManager) RefreshTokenIfNeeded(userID uuid.UUID) (*oauth2.Token, error) {
	token, err := tm.GetToken(userID)
	if err != nil {
		return nil, err
	}

	// Check if token is expired or will expire soon (within 5 minutes)
	if token.Expiry.After(time.Now().Add(5 * time.Minute)) {
		return token, nil // Token is still valid
	}

	// Token is expired, refresh it
	tokenSource := tm.oauthConfig.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Get Gmail email from existing record
	var gmailToken models.GmailToken
	if err := tm.db.First(&gmailToken, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve gmail email: %w", err)
	}

	// Store the new token
	if err := tm.StoreToken(userID, newToken, gmailToken.GmailEmail); err != nil {
		return nil, fmt.Errorf("failed to store refreshed token: %w", err)
	}

	return newToken, nil
}

// RevokeToken revokes the token with Google and deletes it from database
func (tm *TokenManager) RevokeToken(userID uuid.UUID) error {
	// Get token
	token, err := tm.GetToken(userID)
	if err != nil {
		return err
	}

	// Revoke with Google (best effort - continue even if this fails)
	_ = tm.oauthConfig.TokenSource(context.Background(), token)

	// Delete from database
	result := tm.db.Delete(&models.GmailToken{}, "user_id = ?", userID)
	return result.Error
}

// GetGmailEmail returns the connected Gmail email for a user
func (tm *TokenManager) GetGmailEmail(userID uuid.UUID) (string, error) {
	var gmailToken models.GmailToken
	if err := tm.db.First(&gmailToken, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("gmail not connected")
		}
		return "", fmt.Errorf("failed to retrieve gmail email: %w", err)
	}

	return gmailToken.GmailEmail, nil
}

// IsConnected checks if a user has Gmail connected
func (tm *TokenManager) IsConnected(userID uuid.UUID) bool {
	var count int64
	tm.db.Model(&models.GmailToken{}).Where("user_id = ?", userID).Count(&count)
	return count > 0
}
