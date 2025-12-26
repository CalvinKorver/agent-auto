package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService *services.MessageService
	emailService   *services.EmailService
}

func NewMessageHandler(messageService *services.MessageService, emailService *services.EmailService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		emailService:   emailService,
	}
}

// MessageRequest represents the request to send a message
type MessageRequest struct {
	Content string `json:"content"`
	Sender  string `json:"sender"` // "user" or "seller" (for testing)
}

// MessageResponse represents a message in API responses
type MessageResponse struct {
	ID                string `json:"id"`
	ThreadID          string `json:"threadId"`
	Sender            string `json:"sender"`
	Content           string `json:"content"`
	Timestamp         string `json:"timestamp"`
	ExternalMessageID string `json:"externalMessageId,omitempty"`
	SenderEmail       string `json:"senderEmail,omitempty"`
	Subject           string `json:"subject,omitempty"`
}

// GetMessages retrieves messages for a thread
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
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

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, total, err := h.messageService.GetThreadMessages(threadID, userID, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "thread not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	var response struct {
		Messages []MessageResponse `json:"messages"`
		Total    int64             `json:"total"`
		HasMore  bool              `json:"hasMore"`
	}

	response.Messages = make([]MessageResponse, len(messages))
	for i, msg := range messages {
		response.Messages[i] = MessageResponse{
			ID:                msg.ID.String(),
			ThreadID:          msg.ThreadID.String(),
			Sender:            string(msg.Sender),
			Content:           msg.Content,
			Timestamp:         msg.Timestamp.Format("2006-01-02T15:04:05Z"),
			ExternalMessageID: msg.ExternalMessageID,
			SenderEmail:       msg.SenderEmail,
			Subject:           msg.Subject,
		}
	}

	response.Total = total
	response.HasMore = int64(offset+len(messages)) < total

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateMessage creates a new message in a thread
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
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

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Handle different sender types
	if req.Sender == "seller" {
		// Create seller message (for testing)
		sellerMsg, err := h.messageService.CreateSellerMessage(threadID, userID, req.Content)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(struct {
			SellerMessage MessageResponse `json:"sellerMessage"`
		}{
			SellerMessage: MessageResponse{
				ID:                sellerMsg.ID.String(),
				ThreadID:          sellerMsg.ThreadID.String(),
				Sender:            string(sellerMsg.Sender),
				Content:           sellerMsg.Content,
				Timestamp:         sellerMsg.Timestamp.Format("2006-01-02T15:04:05Z"),
				ExternalMessageID: sellerMsg.ExternalMessageID,
				SenderEmail:       sellerMsg.SenderEmail,
				Subject:           sellerMsg.Subject,
			},
		})
		return
	}

	// Default: create user message + agent response
	userMsg, agentMsg, err := h.messageService.CreateUserMessage(threadID, userID, req.Content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := struct {
		UserMessage  MessageResponse  `json:"userMessage"`
		AgentMessage *MessageResponse `json:"agentMessage,omitempty"`
	}{
		UserMessage: MessageResponse{
			ID:                userMsg.ID.String(),
			ThreadID:          userMsg.ThreadID.String(),
			Sender:            string(userMsg.Sender),
			Content:           userMsg.Content,
			Timestamp:         userMsg.Timestamp.Format("2006-01-02T15:04:05Z"),
			ExternalMessageID: userMsg.ExternalMessageID,
			SenderEmail:       userMsg.SenderEmail,
			Subject:           userMsg.Subject,
		},
	}

	if agentMsg != nil {
		response.AgentMessage = &MessageResponse{
			ID:                agentMsg.ID.String(),
			ThreadID:          agentMsg.ThreadID.String(),
			Sender:            string(agentMsg.Sender),
			Content:           agentMsg.Content,
			Timestamp:         agentMsg.Timestamp.Format("2006-01-02T15:04:05Z"),
			ExternalMessageID: agentMsg.ExternalMessageID,
			SenderEmail:       agentMsg.SenderEmail,
			Subject:           agentMsg.Subject,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetInboxMessages retrieves inbox messages (messages with no thread)
func (h *MessageHandler) GetInboxMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, total, err := h.messageService.GetInboxMessages(userID, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
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

	var response struct {
		Messages []InboxMessageResponse `json:"messages"`
		Total    int64                  `json:"total"`
		HasMore  bool                   `json:"hasMore"`
	}

	response.Messages = make([]InboxMessageResponse, len(messages))
	for i, msg := range messages {
		response.Messages[i] = InboxMessageResponse{
			ID:                msg.ID.String(),
			Sender:            string(msg.Sender),
			SenderEmail:       msg.SenderEmail,
			Subject:           msg.Subject,
			Content:           msg.Content,
			Timestamp:         msg.Timestamp.Format("2006-01-02T15:04:05Z"),
			ExternalMessageID: msg.ExternalMessageID,
		}
	}

	response.Total = total
	response.HasMore = int64(offset+len(messages)) < total

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// AssignInboxMessageToThread assigns an inbox message to a thread
func (h *MessageHandler) AssignInboxMessageToThread(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	messageIDStr := chi.URLParam(r, "id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid message ID"})
		return
	}

	var req struct {
		ThreadID string `json:"threadId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	threadID, err := uuid.Parse(req.ThreadID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid thread ID"})
		return
	}

	if err := h.messageService.AssignInboxMessageToThread(messageID, threadID, userID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "thread not found" || err.Error() == "inbox message not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "assigned successfully"})
}

// ArchiveInboxMessage soft deletes an inbox message
func (h *MessageHandler) ArchiveInboxMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	messageIDStr := chi.URLParam(r, "id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid message ID"})
		return
	}

	if err := h.messageService.ArchiveInboxMessage(messageID, userID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "inbox message not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "archived successfully"})
}

// ReplyViaGmail sends an email reply via user's connected Gmail
// POST /api/v1/messages/{messageId}/reply-via-gmail
func (h *MessageHandler) ReplyViaGmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	messageIDStr := chi.URLParam(r, "messageId")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid message ID"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "content is required"})
		return
	}

	// Send reply via Gmail
	if err := h.emailService.ReplyViaGmail(userID, messageID, req.Content); err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "message not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if err.Error() == "message was not received via email" {
			w.WriteHeader(http.StatusBadRequest)
		} else if err.Error() == "gmail not connected" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "email sent successfully"})
}
