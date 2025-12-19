package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"carbuyer/internal/db/models"
	"carbuyer/internal/services"

	"gorm.io/gorm"
)

type EmailHandler struct {
	emailService         *services.EmailService
	webhookSigningKey    string
	db                   *gorm.DB
}

func NewEmailHandler(emailService *services.EmailService, webhookSigningKey string, db *gorm.DB) *EmailHandler {
	return &EmailHandler{
		emailService:      emailService,
		webhookSigningKey: webhookSigningKey,
		db:                db,
	}
}

// MailgunWebhookPayload represents the parsed email data from Mailgun
type MailgunWebhookPayload struct {
	Signature struct {
		Timestamp string `json:"timestamp"`
		Token     string `json:"token"`
		Signature string `json:"signature"`
	} `json:"signature"`
	EventData struct {
		Sender      string `json:"sender"`
		Recipient   string `json:"recipient"`
		Subject     string `json:"subject"`
		BodyPlain   string `json:"body-plain"`
		MessageID   string `json:"Message-Id"`
		From        string `json:"from"`
		To          string `json:"to"`
	} `json:"event-data"`
}

// InboundEmail handles incoming emails from Mailgun webhook
func (h *EmailHandler) InboundEmail(w http.ResponseWriter, r *http.Request) {
	// Parse form data from Mailgun route (can be multipart or urlencoded)
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid form data"})
		return
	}

	// Extract form fields
	recipientEmail := r.FormValue("recipient")
	from := r.FormValue("from")
	sender := r.FormValue("sender")
	subject := r.FormValue("subject")
	bodyPlain := r.FormValue("stripped-text") // Use stripped-text for clean message content
	messageID := r.FormValue("Message-Id")

	// Debug: log what we received
	log.Printf("Webhook received - recipient: %s, from: %s, subject: %s, body length: %d",
		recipientEmail, from, subject, len(bodyPlain))

	// Use sender if from is empty
	if from == "" {
		from = sender
	}

	// Look up user by inbox_email
	var user models.User
	if err := h.db.Where("inbox_email = ?", recipientEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found - return 200 to prevent retries, but log
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "user not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "database error"})
		return
	}

	// Process inbound email
	message, err := h.emailService.ProcessInboundEmail(user.ID, from, subject, bodyPlain, messageID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to process email"})
		return
	}

	// Return 200 OK (critical for Mailgun)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"message_id": message.ID.String(),
	})
}

// verifySignature verifies the Mailgun webhook signature
func (h *EmailHandler) verifySignature(timestamp, token, signature string) bool {
	// Compute expected signature
	h1 := hmac.New(sha256.New, []byte(h.webhookSigningKey))
	h1.Write([]byte(timestamp))
	h1.Write([]byte(token))
	expected := hex.EncodeToString(h1.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expected))
}

// Simple test endpoint for testing without Mailgun
func (h *EmailHandler) TestInboundEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InboxEmail string `json:"inboxEmail"`
		From       string `json:"from"`
		Subject    string `json:"subject"`
		Body       string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	// Look up user
	var user models.User
	if err := h.db.Where("inbox_email = ?", req.InboxEmail).First(&user).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}

	// Process email
	message, err := h.emailService.ProcessInboundEmail(user.ID, req.From, req.Subject, req.Body, "test-"+strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
