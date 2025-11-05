package router

import (
	"feedback-agent/app/handler"

	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the application router
func NewRouter(feedbackHandler *handler.FeedbackHandler, healthHandler *handler.HealthHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", healthHandler.Check)

	// API routes
	api := router.Group("/api")
	{
		// Feedback routes
		api.POST("/feedback", feedbackHandler.CreateFeedback)
		api.GET("/feedback", feedbackHandler.GetAllFeedback)
		api.GET("/feedback/:id", feedbackHandler.GetFeedback)
	}

	return router
}