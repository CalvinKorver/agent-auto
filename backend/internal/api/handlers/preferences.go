package handlers

import (
	"encoding/json"
	"net/http"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/services"
)

type PreferencesHandler struct {
	prefsService *services.PreferencesService
}

func NewPreferencesHandler(prefsService *services.PreferencesService) *PreferencesHandler {
	return &PreferencesHandler{
		prefsService: prefsService,
	}
}

// PreferencesRequest represents the request body for creating/updating preferences
type PreferencesRequest struct {
	Year  int    `json:"year"`
	Make  string `json:"make"`
	Model string `json:"model"`
}

// PreferencesFullResponse represents preferences in API responses
type PreferencesFullResponse struct {
	Year      int    `json:"year"`
	Make      string `json:"make"`
	Model     string `json:"model"`
	CreatedAt string `json:"createdAt"`
}

// GetPreferences returns the user's preferences
func (h *PreferencesHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get preferences
	prefs, err := h.prefsService.GetUserPreferences(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "preferences not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "No preferences set"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to retrieve preferences"})
		}
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PreferencesFullResponse{
		Year:      prefs.Year,
		Make:      prefs.Make,
		Model:     prefs.Model,
		CreatedAt: prefs.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// CreatePreferences creates user preferences
func (h *PreferencesHandler) CreatePreferences(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Parse request body
	var req PreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Create preferences
	prefs, err := h.prefsService.CreateUserPreferences(userID, req.Year, req.Make, req.Model)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "preferences already exist for this user" {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PreferencesFullResponse{
		Year:      prefs.Year,
		Make:      prefs.Make,
		Model:     prefs.Model,
		CreatedAt: prefs.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdatePreferences updates user preferences
func (h *PreferencesHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Parse request body
	var req PreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Update preferences
	prefs, err := h.prefsService.UpdateUserPreferences(userID, req.Year, req.Make, req.Model)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "preferences not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PreferencesFullResponse{
		Year:      prefs.Year,
		Make:      prefs.Make,
		Model:     prefs.Model,
		CreatedAt: prefs.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
