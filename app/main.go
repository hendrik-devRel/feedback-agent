package main

import (
	"log"

	"feedback-agent/app/config"
	"feedback-agent/app/handler"
	"feedback-agent/app/infrastructure/database"
	"feedback-agent/app/repository"
	"feedback-agent/app/router"
	"feedback-agent/app/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup database connection
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	feedbackRepo := repository.NewFeedbackRepository(db)

	// Initialize services
	feedbackService := service.NewFeedbackService(feedbackRepo)

	// Initialize handlers
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)
	healthHandler := handler.NewHealthHandler()

	// Setup router
	r := router.NewRouter(feedbackHandler, healthHandler)

	// Start server
	log.Printf("Starting server on %s", cfg.ServerPort)
	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}