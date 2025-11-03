package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
	"feedback-agent/app/models/request"
)

func main() {
	// Database connection
import (
	"os"
)

func main() {
	// Database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/feedback?sslmode=disable" // fallback for local dev
	}
	db, err := sql.Open("postgres", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database!")
	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create feedback endpoint
	router.POST("/api/feedback", createFeedbackHandler(db))

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createFeedbackHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateFeedbackRequest

		// Bind JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set default sentiment if not provided
		sentiment := enum.SentimentNeutral
		if req.Sentiment != nil {
			sentiment = *req.Sentiment
		}

		// Insert into database
		// Auto-generated: id, votes (defaults to 0), created_at, updated_at
		query := `
			INSERT INTO feedback (title, description, type, tags, sentiment, sentiment_score, votes)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, votes, created_at, updated_at
		`

		var feedback entity.Feedback
		err := db.QueryRow(
			query,
			req.Title,
			req.Description,
			int(req.Type),
			pq.Array(req.Tags), // Convert []string to PostgreSQL array
			int(sentiment),
			req.SentimentScore,
			0, // Votes always start at 0
		).Scan(&feedback.ID, &feedback.Votes, &feedback.CreatedAt, &feedback.UpdatedAt)

		if err != nil {
			log.Printf("Failed to insert feedback: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
			return
		}

		// Set the input fields
		feedback.Title = req.Title
		feedback.Description = req.Description
		feedback.Type = req.Type
		feedback.Tags = req.Tags
		feedback.Sentiment = sentiment
		feedback.SentimentScore = req.SentimentScore

		c.JSON(http.StatusCreated, feedback)
	}
}