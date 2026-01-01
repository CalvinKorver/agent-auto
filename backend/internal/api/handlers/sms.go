package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"carbuyer/internal/api/middleware"
	"carbuyer/internal/db/models"
	"carbuyer/internal/services"
	"carbuyer/internal/twilio"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SMSHandler struct {
	db           *gorm.DB
	smsService   *services.SMSService
	twilioClient *twilio.Client
	authToken    string
}

func NewSMSHandler(db *gorm.DB, smsService *services.SMSService, twilioClient *twilio.Client, authToken string) *SMSHandler {
	return &SMSHandler{
		db:           db,
		smsService:   smsService,
		twilioClient: twilioClient,
		authToken:    authToken,
	}
}

// InboundSMS handles incoming SMS webhooks from Twilio
// POST /api/v1/webhooks/sms/inbound
func (h *SMSHandler) InboundSMS(w http.ResponseWriter, r *http.Request) {
	// Parse form data (Twilio sends webhooks as form-encoded)
	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid form data"})
		return
	}

	// Validate webhook signature for security
	signature := r.Header.Get("X-Twilio-Signature")
	if signature != "" && h.authToken != "" {
		// Build URL for signature validation
		scheme := "https"
		if r.TLS == nil {
			scheme = "http"
		}
		webhookURL := scheme + "://" + r.Host + r.URL.Path
		
		// Convert form values to map for validation
		params := make(map[string]string)
		for k, v := range r.Form {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}

		// Validate signature
		if !twilio.ValidateWebhookSignature(webhookURL, params, signature, h.authToken) {
			log.Printf("Invalid Twilio signature")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid signature"})
			return
		}
	} else if signature == "" {
		log.Printf("Warning: Missing X-Twilio-Signature header (webhook may be insecure)")
	}

	// Extract webhook data
	to := r.FormValue("To")        // The Twilio number that received the SMS
	from := r.FormValue("From")    // The dealer's phone number
	body := r.FormValue("Body")    // The message content
	messageSID := r.FormValue("MessageSid")
	log.Printf("Inbound SMS - To: %s, From: %s, Body: %s", to, from, body)

	// Lookup user by Twilio phone number
	var user models.User
	if err := h.db.Where("phone_number = ?", to).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User not found - return 200 to prevent retries, but log
			log.Printf("User not found for Twilio number: %s", to)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "user not found"})
			return
		}
		log.Printf("Database error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "database error"})
		return
	}

	// Process inbound SMS
	message, err := h.smsService.ProcessInboundSMS(user.ID, from, body, messageSID)
	if err != nil {
		log.Printf("Failed to process inbound SMS: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to process SMS"})
		return
	}

	// Return 200 OK (critical for Twilio)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"message_id": message.ID.String(),
	})
}

// SendSMS handles sending SMS replies
// POST /api/v1/messages/{messageId}/sms-reply
func (h *SMSHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
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

	// Get the message to find the thread
	var message models.Message
	if err := h.db.Where("id = ? AND user_id = ?", messageID, userID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "message not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "database error"})
		return
	}

	if message.ThreadID == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "message is not assigned to a thread"})
		return
	}

	// Send SMS via service
	if err := h.smsService.SendSMS(userID, *message.ThreadID, req.Content); err != nil {
		w.Header().Set("Content-Type", "application/json")
		errMsg := err.Error()
		if errMsg == "user not found" || errMsg == "thread not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if errMsg == "user does not have a Twilio phone number allocated" {
			w.WriteHeader(http.StatusBadRequest)
		} else if errMsg == "thread does not have a phone number assigned" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "SMS sent successfully"})
}

// GetPhoneNumber returns the user's allocated phone number
// GET /api/v1/sms/phone-number
func (h *SMSHandler) GetPhoneNumber(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"phoneNumber": user.PhoneNumber,
	})
}
