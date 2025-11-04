// @title           Feedback Agent API
// @version         1.0
// @description     API for managing feedback with LLM-based auto-classification
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
	"feedback-agent/app/models/request"
	"feedback-agent/app/infrastructure/llm"
	service "feedback-agent/app/models/services"
	_ "feedback-agent/docs" // This will be created after running swag init
)

func main() {
	
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create LLM client
	llmClient := llm.NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))

	// Create classification service
	classificationService := service.NewClassificationService(llmClient)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database!")
	// Setup Gin router
	router := gin.Default()

	// Swagger documentation endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// @Summary Health check endpoint
	// @Description Returns the health status of the API
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// @Summary Create feedback
	// @Description Create a new feedback entry with optional auto-classification
	// @Tags feedback
	// @Accept json
	// @Produce json
	// @Param feedback body request.CreateFeedbackRequest true "Feedback data"
	// @Success 201 {object} entity.Feedback
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /api/feedback [post]
	router.POST("/api/feedback", createFeedbackHandler(db, classificationService))

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createFeedbackHandler(db *sql.DB, classifier *service.ClassificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateFeedbackRequest

		// Bind JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Auto-classify if Type not provided
		var feedbackType enum.FeedbackType
		if req.Type != nil {
			feedbackType = *req.Type
		} else {
			// Use LLM classification service
			ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
			defer cancel()
			feedbackType = classifier.ClassifyFeedback(ctx, req.Title, req.Description)
			log.Printf("Auto-classified feedback as: %s", feedbackType.Name())
		}
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
			int(feedbackType), // Use feedbackType instead of req.Type
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
		feedback.Type = feedbackType // Use feedbackType instead of req.Type
		feedback.Tags = req.Tags
		feedback.Sentiment = sentiment
		feedback.SentimentScore = req.SentimentScore

		c.JSON(http.StatusCreated, feedback)
	}
}