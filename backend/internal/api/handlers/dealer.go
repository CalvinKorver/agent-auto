package handlers

import (
	"encoding/json"
	"net/http"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/services"

	"github.com/google/uuid"
)

type DealerHandler struct {
	dealerService *services.DealerService
	prefsService  *services.PreferencesService
}

func NewDealerHandler(dealerService *services.DealerService, prefsService *services.PreferencesService) *DealerHandler {
	return &DealerHandler{
		dealerService: dealerService,
		prefsService:  prefsService,
	}
}

// DealerResponse represents a dealer in API responses
type DealerResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Website   *string `json:"website,omitempty"`
	Distance  float64 `json:"distance"`
	Contacted bool    `json:"contacted"`
}

// UpdateDealersRequest represents the request body for updating dealers
type UpdateDealersRequest struct {
	DealerIDs []string `json:"dealerIds"`
	Contacted bool     `json:"contacted"`
}

// GetDealers returns dealers for the user's preferences
func (h *DealerHandler) GetDealers(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get user's preferences
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

	// Get dealers for preferences
	dealers, err := h.dealerService.GetDealersForPreferences(prefs.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to retrieve dealers"})
		return
	}

	// Convert to response format
	dealerResponses := make([]DealerResponse, len(dealers))
	for i, dealer := range dealers {
		dealerResponses[i] = DealerResponse{
			ID:        dealer.ID.String(),
			Name:      dealer.Name,
			Location:  dealer.Location,
			Email:     dealer.Email,
			Phone:     dealer.Phone,
			Website:   dealer.Website,
			Distance:  dealer.Distance,
			Contacted:  dealer.Contacted,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dealerResponses)
}

// UpdateDealers updates the contacted status for dealers
func (h *DealerHandler) UpdateDealers(w http.ResponseWriter, r *http.Request) {
	// Verify user is authenticated
	_, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Parse request body
	var req UpdateDealersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Validate dealer IDs
	if len(req.DealerIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "dealerIds are required"})
		return
	}

	// Convert string IDs to UUIDs
	dealerUUIDs := make([]uuid.UUID, 0, len(req.DealerIDs))
	for _, idStr := range req.DealerIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid dealer ID format"})
			return
		}
		dealerUUIDs = append(dealerUUIDs, id)
	}

	// Update dealers
	if err := h.dealerService.UpdateDealersContacted(dealerUUIDs, req.Contacted); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to update dealers"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dealers updated successfully",
	})
}

