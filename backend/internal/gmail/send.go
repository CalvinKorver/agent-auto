package gmail

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
)

// SendReply sends an email reply maintaining proper threading
// externalMessageID is the Message-ID from the original email (stored in database)
func SendReply(service *gmail.Service, to, subject, htmlBody, externalMessageID string) error {
	// Build email message with proper threading headers
	// Format: RFC 2822
	message := buildReplyMessage(to, subject, htmlBody, externalMessageID)

	fmt.Printf("=== SENDING EMAIL VIA GMAIL ===\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("In-Reply-To: %s\n", externalMessageID)
	fmt.Printf("Full message:\n%s\n", message)
	fmt.Printf("================================\n")

	// Encode message to base64url (required by Gmail API)
	encoded := base64.URLEncoding.EncodeToString([]byte(message))
	// Gmail API requires base64url (no padding)
	encoded = strings.TrimRight(encoded, "=")
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")

	// Create Gmail message
	gmailMessage := &gmail.Message{
		Raw: encoded,
	}

	// Send the message
	// Gmail will automatically thread based on In-Reply-To and References headers
	_, err := service.Users.Messages.Send("me", gmailMessage).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildReplyMessage constructs the MIME message with threading headers
func buildReplyMessage(to, subject, htmlBody, externalMessageID string) string {
	// Ensure subject has "Re:" prefix
	if !strings.HasPrefix(subject, "Re:") {
		subject = "Re: " + subject
	}

	// Build message with threading headers
	var sb strings.Builder

	// Standard headers
	sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	// Threading headers - these ensure the reply is threaded with the original
	if externalMessageID != "" {
		// Ensure Message-ID is wrapped in angle brackets
		messageID := externalMessageID
		if !strings.HasPrefix(messageID, "<") {
			messageID = "<" + messageID
		}
		if !strings.HasSuffix(messageID, ">") {
			messageID = messageID + ">"
		}

		// In-Reply-To: references the message we're replying to
		sb.WriteString(fmt.Sprintf("In-Reply-To: %s\r\n", messageID))

		// References: includes the full chain (in this case, just the original)
		sb.WriteString(fmt.Sprintf("References: %s\r\n", messageID))
	}

	// Content-Type for HTML email
	sb.WriteString("Content-Type: text/html; charset=utf-8\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")

	// Empty line to separate headers from body
	sb.WriteString("\r\n")

	// Email body
	sb.WriteString(htmlBody)

	return sb.String()
}
