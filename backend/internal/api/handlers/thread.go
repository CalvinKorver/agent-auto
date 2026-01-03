package handlers

import (
	"encoding/json"
	"net/http"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/db/models"
	"carbuyer/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ThreadHandler struct {
	threadService *services.ThreadService
}

func NewThreadHandler(threadService *services.ThreadService) *ThreadHandler {
	return &ThreadHandler{
		threadService: threadService,
	}
}

// CreateThreadRequest represents the request to create a thread
type CreateThreadRequest struct {
	SellerName string `json:"sellerName"`
	SellerType string `json:"sellerType"`
}

// UpdateThreadRequest represents the request to update a thread
type UpdateThreadRequest struct {
	SellerName *string `json:"sellerName,omitempty"`
	MarkAsRead *bool   `json:"markAsRead,omitempty"`
}

// ThreadResponse represents a thread in API responses
type ThreadResponse struct {
	ID                string  `json:"id"`
	SellerName        string  `json:"sellerName"`
	SellerType        string  `json:"sellerType"`
	Phone             string  `json:"phone,omitempty"`
	CreatedAt         string  `json:"createdAt"`
	LastMessageAt     *string `json:"lastMessageAt,omitempty"`
	MessageCount      int64   `json:"messageCount"`
	UnreadCount       int64   `json:"unreadCount"`
	LastMessagePreview string `json:"lastMessagePreview,omitempty"`
	DisplayName       string  `json:"displayName"`
}

// Helper methods to reduce code duplication

// getUserID extracts and validates the user ID from context
func (h *ThreadHandler) getUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}
	return userID, true
}

// respondError sends a JSON error response
func (h *ThreadHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// respondJSON sends a JSON success response
func (h *ThreadHandler) respondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// CreateThread creates a new thread
func (h *ThreadHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	var req CreateThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Convert seller type string to enum
	sellerType := models.SellerType(req.SellerType)

	thread, err := h.threadService.CreateThread(userID, req.SellerName, sellerType)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Calculate display name for response
	displayName := thread.SellerName
	if thread.SellerName == thread.Phone || thread.SellerName == "" {
		displayName = thread.Phone
	}

	h.respondJSON(w, http.StatusCreated, ThreadResponse{
		ID:          thread.ID.String(),
		SellerName:  thread.SellerName,
		SellerType:  string(thread.SellerType),
		Phone:       thread.Phone,
		CreatedAt:   thread.CreatedAt.Format("2006-01-02T15:04:05Z"),
		DisplayName: displayName,
		// MessageCount and UnreadCount will be 0 for new thread
		MessageCount: 0,
		UnreadCount:  0,
	})
}

// GetThreads retrieves all threads for the authenticated user
func (h *ThreadHandler) GetThreads(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	threadsWithCounts, err := h.threadService.GetUserThreads(userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to retrieve threads")
		return
	}

	var response struct {
		Threads []ThreadResponse `json:"threads"`
	}

	response.Threads = make([]ThreadResponse, len(threadsWithCounts))
	for i, threadWithCounts := range threadsWithCounts {
		threadResp := ThreadResponse{
			ID:                threadWithCounts.ID.String(),
			SellerName:        threadWithCounts.SellerName,
			SellerType:        string(threadWithCounts.SellerType),
			Phone:             threadWithCounts.Phone,
			CreatedAt:         threadWithCounts.CreatedAt.Format("2006-01-02T15:04:05Z"),
			MessageCount:      threadWithCounts.MessageCount,
			UnreadCount:       threadWithCounts.UnreadCount,
			LastMessagePreview: threadWithCounts.LastMessagePreview,
			DisplayName:       threadWithCounts.DisplayName,
		}

		if threadWithCounts.LastMessageAt != nil {
			lastMsg := threadWithCounts.LastMessageAt.Format("2006-01-02T15:04:05Z")
			threadResp.LastMessageAt = &lastMsg
		}

		response.Threads[i] = threadResp
	}

	h.respondJSON(w, http.StatusOK, response)
}

// GetThread retrieves a specific thread
func (h *ThreadHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	threadIDStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid thread ID")
		return
	}

	thread, err := h.threadService.GetThreadByID(threadID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "thread not found" {
			statusCode = http.StatusNotFound
		}
		h.respondError(w, statusCode, err.Error())
		return
	}

	// Calculate display name
	displayName := thread.SellerName
	if thread.SellerName == thread.Phone || thread.SellerName == "" {
		displayName = thread.Phone
	}

	resp := ThreadResponse{
		ID:          thread.ID.String(),
		SellerName:  thread.SellerName,
		SellerType:  string(thread.SellerType),
		Phone:       thread.Phone,
		CreatedAt:   thread.CreatedAt.Format("2006-01-02T15:04:05Z"),
		DisplayName: displayName,
		// Note: MessageCount and UnreadCount would need to be calculated here if needed
		// For now, leaving them as 0 for single thread endpoint
		MessageCount: 0,
		UnreadCount:  0,
	}

	if thread.LastMessageAt != nil {
		lastMsg := thread.LastMessageAt.Format("2006-01-02T15:04:05Z")
		resp.LastMessageAt = &lastMsg
	}

	h.respondJSON(w, http.StatusOK, resp)
}

// ArchiveThread archives (soft deletes) a thread
func (h *ThreadHandler) ArchiveThread(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	threadIDStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid thread ID")
		return
	}

	err = h.threadService.ArchiveThread(threadID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "thread not found" {
			statusCode = http.StatusNotFound
		}
		h.respondError(w, statusCode, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "thread archived successfully"})
}

// UpdateThread updates a thread's seller name and/or marks it as read
func (h *ThreadHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	threadIDStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid thread ID")
		return
	}

	var req UpdateThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate that at least one field is being updated
	if req.SellerName == nil && req.MarkAsRead == nil {
		h.respondError(w, http.StatusBadRequest, "at least one field must be provided")
		return
	}

	// Handle marking thread as read
	if req.MarkAsRead != nil && *req.MarkAsRead {
		if err := h.threadService.MarkThreadAsRead(threadID, userID); err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error() == "thread not found" {
				statusCode = http.StatusNotFound
			}
			h.respondError(w, statusCode, err.Error())
			return
		}
	}

	var thread *models.Thread

	// Handle updating seller name
	if req.SellerName != nil {
		// Validate seller name
		if *req.SellerName == "" {
			h.respondError(w, http.StatusBadRequest, "seller name cannot be empty")
			return
		}

		thread, err = h.threadService.UpdateThreadName(threadID, userID, *req.SellerName)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error() == "thread not found" {
				statusCode = http.StatusNotFound
			}
			h.respondError(w, statusCode, err.Error())
			return
		}
	} else {
		// If only marking as read, fetch the thread to return
		thread, err = h.threadService.GetThreadByID(threadID, userID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error() == "thread not found" {
				statusCode = http.StatusNotFound
			}
			h.respondError(w, statusCode, err.Error())
			return
		}
	}

	// Calculate display name
	displayName := thread.SellerName
	if thread.SellerName == thread.Phone || thread.SellerName == "" {
		displayName = thread.Phone
	}

	resp := ThreadResponse{
		ID:          thread.ID.String(),
		SellerName:  thread.SellerName,
		SellerType:  string(thread.SellerType),
		Phone:       thread.Phone,
		CreatedAt:   thread.CreatedAt.Format("2006-01-02T15:04:05Z"),
		DisplayName: displayName,
	}

	if thread.LastMessageAt != nil {
		lastMsg := thread.LastMessageAt.Format("2006-01-02T15:04:05Z")
		resp.LastMessageAt = &lastMsg
	}

	h.respondJSON(w, http.StatusOK, resp)
}

// MarkThreadAsRead marks a thread as read
func (h *ThreadHandler) MarkThreadAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	threadIDStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(threadIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid thread ID")
		return
	}

	err = h.threadService.MarkThreadAsRead(threadID, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "thread not found" {
			statusCode = http.StatusNotFound
		}
		h.respondError(w, statusCode, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "thread marked as read"})
}

// ConsolidateThreadsRequest represents the request to consolidate threads
type ConsolidateThreadsRequest struct {
	ThreadIDs []string `json:"threadIds"`
}

// ConsolidateThreads consolidates multiple threads into one
func (h *ThreadHandler) ConsolidateThreads(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	var req ConsolidateThreadsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate minimum thread count
	if len(req.ThreadIDs) < 2 {
		h.respondError(w, http.StatusBadRequest, "at least 2 threads required")
		return
	}

	// Convert string IDs to UUIDs
	threadUUIDs := make([]uuid.UUID, len(req.ThreadIDs))
	for i, idStr := range req.ThreadIDs {
		threadID, err := uuid.Parse(idStr)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "invalid thread ID: "+idStr)
			return
		}
		threadUUIDs[i] = threadID
	}

	// Perform consolidation
	consolidatedThread, err := h.threadService.ConsolidateThreads(userID, threadUUIDs)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "at least 2 threads required for consolidation" ||
			err.Error() == "one or more threads not found or already archived" {
			statusCode = http.StatusBadRequest
		}
		h.respondError(w, statusCode, err.Error())
		return
	}

	// Build response
	threadResp := ThreadResponse{
		ID:                consolidatedThread.ID.String(),
		SellerName:        consolidatedThread.SellerName,
		SellerType:        string(consolidatedThread.SellerType),
		Phone:             consolidatedThread.Phone,
		CreatedAt:         consolidatedThread.CreatedAt.Format("2006-01-02T15:04:05Z"),
		MessageCount:      consolidatedThread.MessageCount,
		UnreadCount:       consolidatedThread.UnreadCount,
		LastMessagePreview: consolidatedThread.LastMessagePreview,
		DisplayName:       consolidatedThread.DisplayName,
	}

	if consolidatedThread.LastMessageAt != nil {
		lastMsg := consolidatedThread.LastMessageAt.Format("2006-01-02T15:04:05Z")
		threadResp.LastMessageAt = &lastMsg
	}

	h.respondJSON(w, http.StatusOK, threadResp)
}
