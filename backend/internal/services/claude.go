package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"carbuyer/internal/db/models"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// DealerInfo represents dealer information returned from Claude
type DealerInfo struct {
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Website  *string `json:"website,omitempty"`
	Distance float64 `json:"distance"`
}

type ClaudeService struct {
	client anthropic.Client
}

func NewClaudeService(apiKey string) *ClaudeService {
	return &ClaudeService{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
	}
}

// GenerateNegotiationResponse generates an AI response for car negotiation
func (s *ClaudeService) GenerateNegotiationResponse(
	userMessage string,
	year int,
	make string,
	model string,
	sellerName string,
	messageHistory []models.Message,
	trackedOffers []models.TrackedOffer,
) (string, error) {
	// Build competitive context from tracked offers
	competitiveContext := ""
	if len(trackedOffers) > 0 {
		competitiveContext = "\n\nCompetitive Context - Offers from Other Sellers:\n"
		for _, offer := range trackedOffers {
			sellerInfo := "Unknown Seller"
			if offer.Thread != nil {
				sellerInfo = offer.Thread.SellerName
			}
			competitiveContext += fmt.Sprintf("- %s: %s\n", sellerInfo, offer.OfferText)
		}
		competitiveContext += "\nUse these offers as leverage in your negotiation. Reference competing offers WITHOUT naming specific sellers (e.g., \"I have another dealer offering...\"). This creates competitive pressure."
	}

	// Build system prompt
	systemPrompt := fmt.Sprintf(`You are an expert car negotiation assistant helping a buyer communicate with car sellers.
Your goal is to secure the best possible deal while maintaining professional and respectful communication.

User's Requirements:
- Year: %d
- Make: %s
- Model: %s

Current Seller: %s%s

Guidelines:
- Always negotiate within the user's specified requirements
- Be firm but polite in negotiations
- Work towards the best price and terms for the buyer
- Be professional and concise
- Keep responses around 500 characters by default unless the user explicitly asks for something longer
- Help the user craft effective negotiation messages

- When you have competing offers, use them as leverage without naming specific sellers or offers unless it will meet the goals of the negotiation and purchase.

CRITICAL: When the user asks you to draft, write, or create a message, you must return ONLY the message content itself. Do not include any explanations, prefixes like "Here's a draft:", meta-commentary, or any other text. Return ONLY the message that should be sent to the seller, nothing else.

You should respond as best as you can to whatever the user asks. 
Sometimes they will just chat with you to understand how to best respond. 
Other times they will ask you to draft messages - in those cases, return ONLY the message content. 
You are here to serve the user. `, year, make, model, sellerName, competitiveContext)

	// Build conversation history
	messages := []anthropic.MessageParam{}

	// Add message history for context
	for _, msg := range messageHistory {
		content := msg.Content

		switch msg.Sender {
		case models.SenderTypeUser:
			content = "User's draft message: " + content
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(content)))
		case models.SenderTypeAgent:
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(content)))
		case models.SenderTypeSeller:
			content = fmt.Sprintf("Seller (%s) said: %s", sellerName, content)
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(content)))
		}
	}

	// Add current user message
	userPrompt := fmt.Sprintf("Here is the users message: \"%s\" . In this thread they are negotiating with %s. Assist the user with their request", userMessage, sellerName)
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)))

	// Log the full Claude request
	fmt.Printf("\n========== CLAUDE API REQUEST ==========\n")
	fmt.Printf("System Prompt:\n%s\n\n", systemPrompt)
	fmt.Printf("Message Count: %d\n", len(messages))
	fmt.Printf("Messages:\n")
	for i, msg := range messages {
		role := string(msg.Role)
		fmt.Printf("  [%d] %s:\n", i+1, role)
		// Extract text from content blocks
		for _, block := range msg.Content {
			// Check the union type fields
			if block.OfText != nil {
				fmt.Printf("    Text: %s\n", block.OfText.Text)
			} else if block.OfImage != nil {
				fmt.Printf("    (image block)\n")
			} else if block.OfDocument != nil {
				fmt.Printf("    (document block)\n")
			} else {
				fmt.Printf("    (other block type)\n")
			}
		}
	}
	fmt.Printf("========================================\n\n")

	// Call Claude API using the current SDK pattern
	message, err := s.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: messages,
	})

	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}

	// Extract text from response
	if len(message.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	// Access the text directly from the first content block
	if message.Content[0].Type != "text" {
		return "", fmt.Errorf("unexpected response format from Claude")
	}

	responseText := strings.TrimSpace(message.Content[0].Text)

	// Remove common prefixes that might indicate explanations or meta-commentary
	// This is a safeguard in case Claude still adds explanatory text
	// All these prefixes end with a colon, so we can find the colon position
	prefixes := []string{
		"Here's a draft:",
		"Here's the draft:",
		"Here is a draft:",
		"Here is the draft:",
		"Draft message:",
		"Message draft:",
		"Here's the message:",
		"Here is the message:",
		"Here's your message:",
		"Here is your message:",
	}
	
	responseLower := strings.ToLower(responseText)
	for _, prefix := range prefixes {
		prefixLower := strings.ToLower(prefix)
		if strings.HasPrefix(responseLower, prefixLower) {
			// Find the colon position in the original text (should be at len(prefix)-1)
			// Remove everything up to and including the colon, then trim whitespace
			colonPos := strings.Index(responseText, ":")
			if colonPos >= 0 && colonPos < len(responseText) {
				responseText = strings.TrimSpace(responseText[colonPos+1:])
			}
			break
		}
	}

	// Log the response
	fmt.Printf("\n========== CLAUDE API RESPONSE ==========\n")
	fmt.Printf("Response: %s\n", responseText)
	fmt.Printf("=========================================\n\n")

	return responseText, nil
}

// FetchNearbyDealers fetches nearby dealers using Claude AI
func (s *ClaudeService) FetchNearbyDealers(zipCode string, make string, model string, year int) ([]DealerInfo, error) {
	systemPrompt := `You are a helpful assistant that finds car dealerships. When given a zip code, vehicle make, model, and year, you should return a JSON array of up to 6 nearest dealerships that sell that brand.

For each dealer, provide:
- name: The dealership name
- location: Full address or city, state
- email: Email address if available (can be null)
- phone: Phone number if available (can be null)
- website: Website URL if available (can be null)
- distance: Estimated distance in miles from the zip code to the dealer location

Return ONLY valid JSON, no other text. The JSON should be an array of objects.`

	userPrompt := fmt.Sprintf(`Find up to 6 nearest dealerships for a %d %s %s near zip code %s. Return the results as a JSON array with the structure: [{"name": "Dealer Name", "location": "Address", "email": "email@example.com" or null, "phone": "123-456-7890" or null, "website": "https://example.com" or null, "distance": 5.2}]. Only return the JSON array, no other text.`, year, make, model, zipCode)

	// Call Claude API
	message, err := s.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	// Extract text from response
	if len(message.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	if message.Content[0].Type != "text" {
		return nil, fmt.Errorf("unexpected response format from Claude")
	}

	responseText := strings.TrimSpace(message.Content[0].Text)

	// Try to extract JSON from the response (in case Claude adds extra text)
	// Look for array start
	startIdx := strings.Index(responseText, "[")
	if startIdx == -1 {
		return nil, fmt.Errorf("no JSON array found in Claude response")
	}

	// Look for array end
	endIdx := strings.LastIndex(responseText, "]")
	if endIdx == -1 || endIdx <= startIdx {
		return nil, fmt.Errorf("invalid JSON array in Claude response")
	}

	jsonText := responseText[startIdx : endIdx+1]

	// Parse JSON
	var dealers []DealerInfo
	if err := json.Unmarshal([]byte(jsonText), &dealers); err != nil {
		return nil, fmt.Errorf("failed to parse Claude JSON response: %w", err)
	}

	// Validate we have at least some dealers
	if len(dealers) == 0 {
		return nil, fmt.Errorf("no dealers found in Claude response")
	}

	// Limit to 6 dealers
	if len(dealers) > 6 {
		dealers = dealers[:6]
	}

	return dealers, nil
}
