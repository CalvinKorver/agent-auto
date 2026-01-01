package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                     string
	Environment              string
	DatabaseURL              string
	JWTSecret                string
	JWTExpirationHours       int
	AnthropicAPIKey          string
	AllowedOrigins           []string
	RateLimitAuth            int
	RateLimitAPI             int
	MailgunAPIKey            string
	MailgunDomain            string
	MailgunWebhookSigningKey string
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURL        string
	TokenEncryptionKey       string
	TwilioAccountSID         string
	TwilioAuthToken          string
	TwilioMessagingServiceSID string
	TwilioWebhookURL         string
}

func Load() (*Config, error) {
	port := getEnv("PORT", "8080")
	environment := getEnv("ENVIRONMENT", "development")
	databaseURL := getEnv("DATABASE_URL", "")
	jwtSecret := getEnv("JWT_SECRET", "")
	jwtExpirationHours := getEnvAsInt("JWT_EXPIRATION_HOURS", 24)
	anthropicAPIKey := getEnv("ANTHROPIC_API_KEY", "")
	// Parse ALLOWED_ORIGINS - supports comma-separated list
	originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	allowedOrigins := strings.Split(originsStr, ",")
	rateLimitAuth := getEnvAsInt("RATE_LIMIT_AUTH", 5)
	rateLimitAPI := getEnvAsInt("RATE_LIMIT_API", 100)
	mailgunAPIKey := getEnv("MAILGUN_API_KEY", "")
	mailgunDomain := getEnv("MAILGUN_DOMAIN", "")
	mailgunWebhookSigningKey := getEnv("MAILGUN_WEBHOOK_SIGNING_KEY", "")
	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	googleClientSecret := getEnv("GOOGLE_CLIENT_SECRET", "")
	googleRedirectURL := getEnv("GOOGLE_REDIRECT_URL", "http://localhost:3000/oauth/callback")
	tokenEncryptionKey := getEnv("TOKEN_ENCRYPTION_KEY", "")
	twilioAccountSID := getEnv("TWILIO_ACCOUNT_SID", "")
	twilioAuthToken := getEnv("TWILIO_AUTH_TOKEN", "")
	twilioMessagingServiceSID := getEnv("TWILIO_MESSAGING_SERVICE_SID", "")
	twilioWebhookURL := getEnv("TWILIO_WEBHOOK_URL", "")

	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if anthropicAPIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	return &Config{
		Port:                      port,
		Environment:               environment,
		DatabaseURL:                databaseURL,
		JWTSecret:                 jwtSecret,
		JWTExpirationHours:        jwtExpirationHours,
		AnthropicAPIKey:           anthropicAPIKey,
		AllowedOrigins:            allowedOrigins,
		RateLimitAuth:             rateLimitAuth,
		RateLimitAPI:              rateLimitAPI,
		MailgunAPIKey:             mailgunAPIKey,
		MailgunDomain:             mailgunDomain,
		MailgunWebhookSigningKey:  mailgunWebhookSigningKey,
		GoogleClientID:            googleClientID,
		GoogleClientSecret:        googleClientSecret,
		GoogleRedirectURL:         googleRedirectURL,
		TokenEncryptionKey:        tokenEncryptionKey,
		TwilioAccountSID:         twilioAccountSID,
		TwilioAuthToken:           twilioAuthToken,
		TwilioMessagingServiceSID:  twilioMessagingServiceSID,
		TwilioWebhookURL:          twilioWebhookURL,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
