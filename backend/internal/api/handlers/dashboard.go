package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/db"
	"carbuyer/internal/db/models"
	"carbuyer/internal/services"

	"github.com/google/uuid"
)

type DashboardHandler struct {
	threadService  *services.ThreadService
	messageService *services.MessageService
	db             *db.Database
}

func NewDashboardHandler(threadService *services.ThreadService, messageService *services.MessageService, database *db.Database) *DashboardHandler {
	return &DashboardHandler{
		threadService:  threadService,
		messageService: messageService,
		db:             database,
	}
}

// DashboardResponse represents the consolidated dashboard data
type DashboardResponse struct {
	Threads       []ThreadResponse       `json:"threads"`
	InboxMessages []InboxMessageResponse `json:"inboxMessages"`
	Offers        []OfferResponse        `json:"offers"`
}

// InboxMessageResponse includes additional fields for inbox messages
type InboxMessageResponse struct {
	ID                string `json:"id"`
	Sender            string `json:"sender"`
	SenderEmail       string `json:"senderEmail"`
	Subject           string `json:"subject"`
	Content           string `json:"content"`
	Timestamp         string `json:"timestamp"`
	ExternalMessageID string `json:"externalMessageId,omitempty"`
}

// GetDashboard retrieves all dashboard data (threads, inbox messages, and offers) in a single request
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Use goroutines to fetch all data in parallel
	var threads []models.Thread
	var inboxMessages []models.Message
	var offers []models.TrackedOffer
	var threadsErr, messagesErr, offersErr error

	var wg sync.WaitGroup
	wg.Add(3)

	// Fetch threads
	go func() {
		defer wg.Done()
		threads, threadsErr = h.fetchThreads(userID)
	}()

	// Fetch inbox messages
	go func() {
		defer wg.Done()
		inboxMessages, _, messagesErr = h.messageService.GetInboxMessages(userID, 50, 0)
	}()

	// Fetch offers
	go func() {
		defer wg.Done()
		offers, offersErr = h.fetchOffers(userID)
	}()

	wg.Wait()

	// Check for errors
	if threadsErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to fetch threads"})
		return
	}

	if messagesErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to fetch inbox messages"})
		return
	}

	if offersErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to fetch offers"})
		return
	}

	// Build response
	response := DashboardResponse{
		Threads:       make([]ThreadResponse, len(threads)),
		InboxMessages: make([]InboxMessageResponse, len(inboxMessages)),
		Offers:        make([]OfferResponse, len(offers)),
	}

	// Convert threads
	for i, thread := range threads {
		threadResp := ThreadResponse{
			ID:           thread.ID.String(),
			SellerName:   thread.SellerName,
			SellerType:   string(thread.SellerType),
			CreatedAt:    thread.CreatedAt.Format("2006-01-02T15:04:05Z"),
			MessageCount: thread.MessageCount,
		}

		if thread.LastMessageAt != nil {
			lastMsg := thread.LastMessageAt.Format("2006-01-02T15:04:05Z")
			threadResp.LastMessageAt = &lastMsg
		}

		response.Threads[i] = threadResp
	}

	// Convert inbox messages
	for i, msg := range inboxMessages {
		response.InboxMessages[i] = InboxMessageResponse{
			ID:                msg.ID.String(),
			Sender:            string(msg.Sender),
			SenderEmail:       msg.SenderEmail,
			Subject:           msg.Subject,
			Content:           msg.Content,
			Timestamp:         msg.Timestamp.Format("2006-01-02T15:04:05Z"),
			ExternalMessageID: msg.ExternalMessageID,
		}
	}

	// Convert offers
	for i, offer := range offers {
		offerResp := OfferResponse{
			ID:        offer.ID.String(),
			ThreadID:  offer.ThreadID.String(),
			OfferText: offer.OfferText,
			TrackedAt: offer.TrackedAt.Format("2006-01-02T15:04:05Z"),
		}

		if offer.MessageID != nil {
			msgID := offer.MessageID.String()
			offerResp.MessageID = &msgID
		}

		// Include thread details if available
		if offer.Thread != nil {
			offerResp.SellerName = &offer.Thread.SellerName
			threadType := string(offer.Thread.SellerType)
			offerResp.ThreadType = &threadType
		}

		response.Offers[i] = offerResp
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// fetchThreads is a helper to fetch threads
func (h *DashboardHandler) fetchThreads(userID uuid.UUID) ([]models.Thread, error) {
	return h.threadService.GetUserThreads(userID)
}

// fetchOffers is a helper to fetch offers
func (h *DashboardHandler) fetchOffers(userID uuid.UUID) ([]models.TrackedOffer, error) {
	var offers []models.TrackedOffer

	// Get all offers for user's threads, ordered by most recent first
	err := h.db.DB.
		Joins("JOIN threads ON threads.id = tracked_offers.thread_id").
		Where("threads.user_id = ?", userID).
		Preload("Thread").
		Order("tracked_offers.tracked_at DESC").
		Find(&offers).Error

	if err != nil {
		return nil, err
	}

	return offers, nil
}


