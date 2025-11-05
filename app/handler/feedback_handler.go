package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"feedback-agent/app/models/request"
	"feedback-agent/app/models/response"
	"feedback-agent/app/service"

	"github.com/gin-gonic/gin"
)

// FeedbackHandler handles HTTP requests for feedback operations
type FeedbackHandler struct {
	feedbackService service.FeedbackService
}

// NewFeedbackHandler creates a new feedback handler
func NewFeedbackHandler(feedbackService service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

// CreateFeedback handles POST /api/feedback
func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	var req request.CreateFeedbackRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	feedback, err := h.feedbackService.CreateFeedback(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, response.FromEntity(feedback))
}

// GetFeedback handles GET /api/feedback/:id
func (h *FeedbackHandler) GetFeedback(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feedback ID"})
		return
	}

	feedback, err := h.feedbackService.GetFeedback(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "feedback not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEntity(feedback))
}

// GetAllFeedback handles GET /api/feedback
func (h *FeedbackHandler) GetAllFeedback(c *gin.Context) {
	feedbacks, err := h.feedbackService.GetAllFeedback(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	responses := make([]*response.FeedbackResponse, len(feedbacks))
	for i, f := range feedbacks {
		responses[i] = response.FromEntity(f)
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}