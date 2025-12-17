package services

import (
	"context"
	"fmt"
	"strings"

	"carbuyer/internal/db/models"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

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
) (string, error) {
	// Build system prompt
	systemPrompt := fmt.Sprintf(`You are an expert car negotiation assistant helping a buyer communicate with car sellers.
Your goal is to secure the best possible deal while maintaining professional and respectful communication.

User's Requirements:
- Year: %d
- Make: %s
- Model: %s

Current Seller: %s

Guidelines:
- Always negotiate within the user's specified requirements
- Be firm but polite in negotiations
- Ask relevant questions about vehicle condition, history, and pricing
- Highlight any concerns or red flags
- Work towards the best price and terms for the buyer
- Never deviate from the specified year, make, and model
- Be professional and concise
- Help the user craft effective negotiation messages

When the user provides a message, enhance it to be more effective for negotiation while keeping their intent.`, year, make, model, sellerName)

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
	userPrompt := fmt.Sprintf("The user wants to send this message: \"%s\"\n\nEnhance this message to be more effective for negotiating with %s. Return ONLY the enhanced message, nothing else.", userMessage, sellerName)
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)))

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

	return strings.TrimSpace(message.Content[0].Text), nil
}
