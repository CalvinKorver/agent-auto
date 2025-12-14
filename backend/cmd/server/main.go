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
	log.Println("Starting Car Buyer Agent API Server...")

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
	authService := services.NewAuthService(database.DB, cfg.JWTSecret, cfg.JWTExpirationHours)
	preferencesService := services.NewPreferencesService(database.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	preferencesHandler := handlers.NewPreferencesHandler(preferencesService)

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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Car Buyer Agent API","version":"1.0.0"}`))
		})

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// Protected auth route
			r.With(middleware.AuthMiddleware(authService)).Get("/me", authHandler.Me)
		})

		// Preferences routes (all protected)
		r.Route("/preferences", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/", preferencesHandler.GetPreferences)
			r.Post("/", preferencesHandler.CreatePreferences)
		})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("API available at http://localhost%s/api/v1", addr)
	log.Printf("Health check at http://localhost%s/health", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
