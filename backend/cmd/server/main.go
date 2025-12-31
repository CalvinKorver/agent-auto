package main

import (
	"fmt"
	"log"
	"net/http"

	"carbuyer/internal/api/handlers"
	"carbuyer/internal/api/middleware"
	"carbuyer/internal/config"
	"carbuyer/internal/db"
	"carbuyer/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Otto API Server...")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Port: %s", cfg.Port)

	// Initialize database
	database, err := db.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService(database.DB, cfg.JWTSecret, cfg.JWTExpirationHours, cfg.MailgunDomain)
	modelsService := services.NewModelsService(database.DB)
	claudeService := services.NewClaudeService(cfg.AnthropicAPIKey)
	dealerService := services.NewDealerService(database.DB, claudeService)
	preferencesService := services.NewPreferencesService(database.DB, modelsService, dealerService)
	threadService := services.NewThreadService(database.DB)
	messageService := services.NewMessageService(database.DB, claudeService)

	// Initialize Gmail service (for sending emails via user's Gmail)
	gmailService, err := services.NewGmailService(
		database.DB,
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		cfg.TokenEncryptionKey,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Gmail service: %v", err)
	}

	emailService := services.NewEmailService(database.DB, cfg.MailgunAPIKey, cfg.MailgunDomain, gmailService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, gmailService)
	preferencesHandler := handlers.NewPreferencesHandler(preferencesService)
	dealerHandler := handlers.NewDealerHandler(dealerService, preferencesService)
	threadHandler := handlers.NewThreadHandler(threadService)
	messageHandler := handlers.NewMessageHandler(messageService, emailService)
	emailHandler := handlers.NewEmailHandler(emailService, cfg.MailgunWebhookSigningKey, database.DB)
	gmailHandler := handlers.NewGmailHandler(gmailService, cfg.AllowedOrigins[0]) // Use first allowed origin as frontend URL
	offerHandler := handlers.NewOfferHandler(database)
	dashboardHandler := handlers.NewDashboardHandler(threadService, messageService, database)
	modelsHandler := handlers.NewModelsHandler(modelsService)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		sqlDB, err := database.DB.DB()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database connection error"))
			return
		}

		if err := sqlDB.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database ping failed"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","database":"connected"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Otto API","version":"1.0.0"}`))
		})

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// Protected auth routes
			r.With(middleware.AuthMiddleware(authService)).Get("/me", authHandler.Me)
			r.With(middleware.AuthMiddleware(authService)).Post("/logout", authHandler.Logout)
		})

		// Preferences routes (all protected)
		r.Route("/preferences", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/", preferencesHandler.GetPreferences)
			r.Post("/", preferencesHandler.CreatePreferences)
			r.Put("/", preferencesHandler.UpdatePreferences)
		})

		// Dealer routes (all protected)
		r.Route("/dealers", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/", dealerHandler.GetDealers)
			r.Put("/", dealerHandler.UpdateDealers)
		})

		// Dashboard route (protected) - consolidated endpoint for threads, inbox messages, and offers
		r.Route("/dashboard", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/", dashboardHandler.GetDashboard)
		})

		// Thread routes (all protected)
		r.Route("/threads", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/", threadHandler.GetThreads)
			r.Post("/", threadHandler.CreateThread)
			r.Get("/{id}", threadHandler.GetThread)
			r.Delete("/{id}", threadHandler.ArchiveThread)

			// Message routes nested under threads
			r.Get("/{id}/messages", messageHandler.GetMessages)
			r.Post("/{id}/messages", messageHandler.CreateMessage)

			// Offer routes nested under threads
			r.Post("/{id}/offers", offerHandler.CreateOffer)
		})

		// Offer routes (all protected)
		r.Route("/offers", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/", offerHandler.GetAllOffers)
			r.Delete("/{id}", offerHandler.DeleteOffer)
		})

		// Inbox message routes (all protected)
		r.Route("/inbox", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/messages", messageHandler.GetInboxMessages)
			r.Put("/messages/{id}/assign", messageHandler.AssignInboxMessageToThread)
			r.Delete("/messages/{id}", messageHandler.ArchiveInboxMessage)
		})

		// Gmail OAuth routes (all protected)
		r.Route("/gmail", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/connect", gmailHandler.GetAuthURL)
			r.Get("/status", gmailHandler.GetGmailStatus)
			r.Post("/disconnect", gmailHandler.DisconnectGmail)
		})

		// Message reply route (protected)
		r.Route("/messages", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Post("/{messageId}/reply-via-gmail", messageHandler.ReplyViaGmail)
			r.Post("/{messageId}/draft", messageHandler.CreateDraftViaGmail)
		})

		// Webhook routes (public - no auth)
		r.Route("/webhooks", func(r chi.Router) {
			r.Post("/email/inbound", emailHandler.InboundEmail)
			r.Post("/email/test", emailHandler.TestInboundEmail) // For testing without Mailgun
		})

		// Models route (public - no auth)
		r.Get("/models", modelsHandler.GetModels)
		r.Get("/trims", modelsHandler.GetTrims)
	})

	// OAuth callback route (public - outside /api/v1)
	r.Get("/oauth/callback", gmailHandler.OAuthCallback)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("API available at http://localhost%s/api/v1", addr)
	log.Printf("Health check at http://localhost%s/health", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
