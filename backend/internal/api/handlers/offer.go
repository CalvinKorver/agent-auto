package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/db"
	"carbuyer/internal/db/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OfferHandler struct {
	db *db.Database
}

func NewOfferHandler(database *db.Database) *OfferHandler {
	return &OfferHandler{
		db: database,
	}
}

// CreateOfferRequest represents the request to create a tracked offer
type CreateOfferRequest struct {
	OfferText string     `json:"offerText"`
	MessageID *uuid.UUID `json:"messageId,omitempty"`
}

// OfferResponse represents an offer in API responses
type OfferResponse struct {
	ID         string  `json:"id"`
	ThreadID   string  `json:"threadId"`
	MessageID  *string `json:"messageId,omitempty"`
	OfferText  string  `json:"offerText"`
	TrackedAt  string  `json:"trackedAt"`
	SellerName *string `json:"sellerName,omitempty"`
	ThreadType *string `json:"threadType,omitempty"`
}

// CreateOffer creates a new tracked offer for a thread
func (h *OfferHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	threadIDStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid thread ID"})
		return
	}

	var req CreateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.OfferText == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "offerText is required"})
		return
	}

	// Verify thread belongs to user
	var thread models.Thread
	if err := h.db.DB.Where("id = ? AND user_id = ?", threadID, userID).First(&thread).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "thread not found"})
		return
	}

	// Create the offer
	offer := models.TrackedOffer{
		ThreadID:  threadID,
		MessageID: req.MessageID,
		OfferText: req.OfferText,
		TrackedAt: time.Now(),
	}

	if err := h.db.DB.Create(&offer).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to create offer"})
		return
	}

	// Build response
	response := OfferResponse{
		ID:        offer.ID.String(),
		ThreadID:  offer.ThreadID.String(),
		OfferText: offer.OfferText,
		TrackedAt: offer.TrackedAt.Format("2006-01-02T15:04:05Z"),
	}

	if offer.MessageID != nil {
		msgID := offer.MessageID.String()
		response.MessageID = &msgID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllOffers retrieves all tracked offers for the authenticated user across all threads
func (h *OfferHandler) GetAllOffers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	var offers []models.TrackedOffer

	// Get all offers for user's threads, ordered by most recent first
	err := h.db.DB.
		Joins("JOIN threads ON threads.id = tracked_offers.thread_id").
		Where("threads.user_id = ?", userID).
		Preload("Thread").
		Order("tracked_offers.tracked_at DESC").
		Find(&offers).Error

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to fetch offers"})
		return
	}

	// Build response with thread details
	response := make([]OfferResponse, len(offers))
	for i, offer := range offers {
		response[i] = OfferResponse{
			ID:        offer.ID.String(),
			ThreadID:  offer.ThreadID.String(),
			OfferText: offer.OfferText,
			TrackedAt: offer.TrackedAt.Format("2006-01-02T15:04:05Z"),
		}

		if offer.MessageID != nil {
			msgID := offer.MessageID.String()
			response[i].MessageID = &msgID
		}

		// Include thread details if available
		if offer.Thread != nil {
			response[i].SellerName = &offer.Thread.SellerName
			threadType := string(offer.Thread.SellerType)
			response[i].ThreadType = &threadType
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Offers []OfferResponse `json:"offers"`
	}{
		Offers: response,
	})
}

// DeleteOffer deletes a tracked offer
func (h *OfferHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	offerIDStr := chi.URLParam(r, "id")
	offerID, err := uuid.Parse(offerIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid offer ID"})
		return
	}

	// Verify offer belongs to user's thread
	var offer models.TrackedOffer
	err = h.db.DB.
		Joins("JOIN threads ON threads.id = tracked_offers.thread_id").
		Where("tracked_offers.id = ? AND threads.user_id = ?", offerID, userID).
		First(&offer).Error

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "offer not found"})
		return
	}

	// Delete the offer
	if err := h.db.DB.Delete(&offer).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to delete offer"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "offer deleted successfully",
	})
}
