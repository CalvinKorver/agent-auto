package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/gmail"
	"carbuyer/internal/services"

	"github.com/google/uuid"
)

type GmailHandler struct {
	gmailService *services.GmailService
	frontendURL  string
}

func NewGmailHandler(gmailService *services.GmailService, frontendURL string) *GmailHandler {
	return &GmailHandler{
		gmailService: gmailService,
		frontendURL:  frontendURL,
	}
}

// GetAuthURL generates and returns OAuth authorization URL
// GET /api/v1/gmail/connect
func (h *GmailHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Generate state token (include user ID for verification in callback)
	state, err := gmail.GenerateStateToken()
	if err != nil {
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
		return
	}

	// Store state in session/cache (for now, we'll encode userID in state)
	// In production, you'd want to store this in Redis or a session store
	stateWithUser := fmt.Sprintf("%s:%s", state, userID.String())

	// Generate OAuth URL
	authURL := h.gmailService.GetAuthURL(stateWithUser)

	// Return URL to frontend
	response := map[string]string{
		"authUrl": authURL,
		"state":   stateWithUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OAuthCallback handles OAuth callback from Google
// GET /oauth/callback?code=...&state=...
func (h *GmailHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Get code and state from query params
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	log.Printf("OAuth callback received - code: %s, state: %s", code, state)

	if code == "" {
		http.Redirect(w, r, h.frontendURL+"/gmail-error?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	// Parse state to get user ID
	// Format: "stateToken:userID"
	parts := strings.Split(state, ":")
	log.Printf("State parts: %v (len=%d)", parts, len(parts))

	if len(parts) != 2 {
		log.Printf("Invalid state format: expected 2 parts, got %d", len(parts))
		http.Redirect(w, r, h.frontendURL+"/gmail-error?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	userID, err := uuid.Parse(parts[1])
	if err != nil {
		log.Printf("Failed to parse UUID from state: %v", err)
		http.Redirect(w, r, h.frontendURL+"/gmail-error?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Successfully parsed userID: %s", userID)

	// Exchange code for token
	token, err := h.gmailService.ExchangeCodeForToken(code)
	if err != nil {
		http.Redirect(w, r, h.frontendURL+"/gmail-error?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get user's Gmail email from token (we need to call Gmail API to get this)
	// For now, we'll store it as empty and update it later
	// TODO: Call Gmail API to get user's email address
	gmailEmail := "" // Will be populated when first email is sent

	// Store token
	if err := h.gmailService.StoreToken(userID, token, gmailEmail); err != nil {
		http.Redirect(w, r, h.frontendURL+"/gmail-error?error=store_failed", http.StatusTemporaryRedirect)
		return
	}

	// Redirect to dashboard with success parameter
	http.Redirect(w, r, h.frontendURL+"/dashboard?gmail_connected=true", http.StatusTemporaryRedirect)
}

// GetGmailStatus returns whether user has Gmail connected
// GET /api/v1/gmail/status
func (h *GmailHandler) GetGmailStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	connected := h.gmailService.IsConnected(userID)

	response := map[string]interface{}{
		"connected": connected,
	}

	if connected {
		email, err := h.gmailService.GetGmailEmail(userID)
		if err == nil && email != "" {
			response["gmailEmail"] = email
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DisconnectGmail revokes Gmail access
// POST /api/v1/gmail/disconnect
func (h *GmailHandler) DisconnectGmail(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Revoke and delete token
	if err := h.gmailService.DisconnectGmail(userID); err != nil {
		http.Error(w, "Failed to disconnect Gmail", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Gmail disconnected successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
