package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/services"
)

type AuthHandler struct {
	authService  *services.AuthService
	gmailService *services.GmailService
}

func NewAuthHandler(authService *services.AuthService, gmailService *services.GmailService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		gmailService: gmailService,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID             string               `json:"id"`
	Email          string               `json:"email"`
	InboxEmail     string               `json:"inboxEmail"`
	CreatedAt      string               `json:"createdAt"`
	Preferences    *PreferencesResponse `json:"preferences,omitempty"`
	GmailConnected bool                 `json:"gmailConnected"`
	GmailEmail     string               `json:"gmailEmail,omitempty"`
}

// PreferencesResponse represents user preferences in API responses
type PreferencesResponse struct {
	Year  int    `json:"year"`
	Make  string `json:"make"`
	Model string `json:"model"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Register user
	user, err := h.authService.RegisterUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Registration error for email %s: %v", req.Email, err)
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "user with this email already exists" {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to generate token"})
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		User: UserResponse{
			ID:         user.ID.String(),
			Email:      user.Email,
			InboxEmail: user.InboxEmail,
			CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Token: token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Authenticate user
	user, err := h.authService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Login error for email %s: %v", req.Email, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to generate token"})
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		User: UserResponse{
			ID:         user.ID.String(),
			Email:      user.Email,
			InboxEmail: user.InboxEmail,
			CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Token: token,
	})
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get user from database
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}

	// Build response
	userResp := UserResponse{
		ID:         user.ID.String(),
		Email:      user.Email,
		InboxEmail: user.InboxEmail,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Add preferences if they exist
	if user.Preferences != nil {
		userResp.Preferences = &PreferencesResponse{
			Year:  user.Preferences.Year,
			Make:  user.Preferences.Make,
			Model: user.Preferences.Model,
		}
	}

	// Check Gmail connection status
	userResp.GmailConnected = h.gmailService.IsConnected(userID)
	if userResp.GmailConnected {
		gmailEmail, err := h.gmailService.GetGmailEmail(userID)
		if err == nil && gmailEmail != "" {
			userResp.GmailEmail = gmailEmail
		}
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

// Logout handles user logout (client-side token removal, server just confirms)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out successfully"})
}
