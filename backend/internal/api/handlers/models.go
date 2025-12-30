package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"carbuyer/internal/services"
)

type ModelsHandler struct {
	modelsService *services.ModelsService
}

func NewModelsHandler(modelsService *services.ModelsService) *ModelsHandler {
	return &ModelsHandler{
		modelsService: modelsService,
	}
}

// GetModels returns all vehicle makes and models
func (h *ModelsHandler) GetModels(w http.ResponseWriter, r *http.Request) {
	models, err := h.modelsService.GetModels()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Models data not yet available"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}

// TrimResponse represents a trim in the API response
type TrimResponse struct {
	ID        string `json:"id"`
	TrimName  string `json:"trimName"`
}

// GetTrims returns trims for a specific make, model, and year
func (h *ModelsHandler) GetTrims(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	makeName := r.URL.Query().Get("make")
	modelName := r.URL.Query().Get("model")
	yearStr := r.URL.Query().Get("year")

	if makeName == "" || modelName == "" || yearStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "make, model, and year query parameters are required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid year parameter"})
		return
	}

	// Lookup make
	makeRecord, err := h.modelsService.GetMakeByName(makeName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "make not found"})
		return
	}

	// Lookup model
	model, err := h.modelsService.GetModelByName(makeRecord.ID, modelName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "model not found"})
		return
	}

	// Get trims for model and year
	trims, err := h.modelsService.GetTrimsForModel(model.ID, year)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to fetch trims"})
		return
	}

	// Convert to response format
	trimResponses := make([]TrimResponse, len(trims))
	for i, trim := range trims {
		trimResponses[i] = TrimResponse{
			ID:       trim.ID.String(),
			TrimName: trim.TrimName,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trimResponses)
}


