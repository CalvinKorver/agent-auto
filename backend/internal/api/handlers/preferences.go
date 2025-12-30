package handlers

import (
	"encoding/json"
	"net/http"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/services"

	"github.com/google/uuid"
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
	Year   int    `json:"year"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	TrimID string `json:"trimId,omitempty"` // Optional: empty string or null means "Unspecified"
}

// PreferencesFullResponse represents preferences in API responses
type PreferencesFullResponse struct {
	Year      int     `json:"year"`
	Make      string  `json:"make"`
	Model     string  `json:"model"`
	Trim      string  `json:"trim,omitempty"` // Trim name, empty if unspecified
	TrimID    string  `json:"trimId,omitempty"` // Trim ID, empty if unspecified
	BaseMSRP  float64 `json:"baseMsrp,omitempty"` // Base MSRP from trim, 0 if not available
	CreatedAt string  `json:"createdAt"`
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

	// Return response with make/model/trim names from relationships
	makeName := ""
	modelName := ""
	trimName := ""
	trimIDStr := ""
	baseMSRP := 0.0
	if prefs.Make != nil {
		makeName = prefs.Make.Name
	}
	if prefs.Model != nil {
		modelName = prefs.Model.Name
	}
	if prefs.Trim != nil {
		trimName = prefs.Trim.TrimName
		trimIDStr = prefs.Trim.ID.String()
		if prefs.Trim.BaseMSRP.Valid {
			baseMSRP = prefs.Trim.BaseMSRP.Float64
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PreferencesFullResponse{
		Year:      prefs.Year,
		Make:      makeName,
		Model:     modelName,
		Trim:      trimName,
		TrimID:    trimIDStr,
		BaseMSRP:  baseMSRP,
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

	// Parse trimID (optional - empty string or null means "Unspecified")
	var trimID *uuid.UUID
	if req.TrimID != "" {
		parsedTrimID, err := uuid.Parse(req.TrimID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid trimId format"})
			return
		}
		trimID = &parsedTrimID
	}

	// Create preferences
	prefs, err := h.prefsService.CreateUserPreferences(userID, req.Year, req.Make, req.Model, trimID)
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

	// Return response with make/model/trim names from relationships
	makeName := ""
	modelName := ""
	trimName := ""
	trimIDStr := ""
	baseMSRP := 0.0
	if prefs.Make != nil {
		makeName = prefs.Make.Name
	}
	if prefs.Model != nil {
		modelName = prefs.Model.Name
	}
	if prefs.Trim != nil {
		trimName = prefs.Trim.TrimName
		trimIDStr = prefs.Trim.ID.String()
		if prefs.Trim.BaseMSRP.Valid {
			baseMSRP = prefs.Trim.BaseMSRP.Float64
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PreferencesFullResponse{
		Year:      prefs.Year,
		Make:      makeName,
		Model:     modelName,
		Trim:      trimName,
		TrimID:    trimIDStr,
		BaseMSRP:  baseMSRP,
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

	// Parse trimID (optional - empty string or null means "Unspecified")
	var trimID *uuid.UUID
	if req.TrimID != "" {
		parsedTrimID, err := uuid.Parse(req.TrimID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid trimId format"})
			return
		}
		trimID = &parsedTrimID
	}

	// Update preferences
	prefs, err := h.prefsService.UpdateUserPreferences(userID, req.Year, req.Make, req.Model, trimID)
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

	// Return response with make/model/trim names from relationships
	makeName := ""
	modelName := ""
	trimName := ""
	trimIDStr := ""
	baseMSRP := 0.0
	if prefs.Make != nil {
		makeName = prefs.Make.Name
	}
	if prefs.Model != nil {
		modelName = prefs.Model.Name
	}
	if prefs.Trim != nil {
		trimName = prefs.Trim.TrimName
		trimIDStr = prefs.Trim.ID.String()
		if prefs.Trim.BaseMSRP.Valid {
			baseMSRP = prefs.Trim.BaseMSRP.Float64
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PreferencesFullResponse{
		Year:      prefs.Year,
		Make:      makeName,
		Model:     modelName,
		Trim:      trimName,
		TrimID:    trimIDStr,
		BaseMSRP:  baseMSRP,
		CreatedAt: prefs.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
