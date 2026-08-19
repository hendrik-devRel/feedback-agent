package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
	"feedback-agent/app/models/request"
	"feedback-agent/app/store"
)

const maxPageSize = 100

// FeedbackHandler exposes the HTTP endpoints for feedback and votes.
type FeedbackHandler struct {
	store *store.FeedbackStore
}

func NewFeedbackHandler(s *store.FeedbackStore) *FeedbackHandler {
	return &FeedbackHandler{store: s}
}

// Register mounts all feedback routes on the router.
func (h *FeedbackHandler) Register(router *gin.Engine) {
	router.POST("/api/feedback", h.Create)
	router.GET("/api/feedback", h.List)
	router.GET("/api/feedback/:id", h.Get)
	router.POST("/api/feedback/:id/votes", h.Vote)
}

// Create handles POST /api/feedback.
func (h *FeedbackHandler) Create(c *gin.Context) {
	var req request.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sentiment := enum.SentimentNeutral
	if req.Sentiment != nil {
		sentiment = *req.Sentiment
	}

	feedback, err := h.store.Create(c.Request.Context(), entity.Feedback{
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		Tags:           req.Tags,
		Sentiment:      sentiment,
		SentimentScore: req.SentimentScore,
	})
	if err != nil {
		log.Printf("Failed to insert feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// List handles GET /api/feedback with optional type/sentiment/tag filters
// and page/pageSize pagination.
func (h *FeedbackHandler) List(c *gin.Context) {
	var q request.ListFeedbackQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter, err := buildListFilter(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, total, err := h.store.List(c.Request.Context(), filter)
	if err != nil {
		log.Printf("Failed to list feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    items,
		"total":    total,
		"page":     q.Page,
		"pageSize": filter.Limit,
	})
}

// Get handles GET /api/feedback/:id.
func (h *FeedbackHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	feedback, err := h.store.GetByID(c.Request.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}
	if err != nil {
		log.Printf("Failed to get feedback %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feedback"})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// Vote handles POST /api/feedback/:id/votes and returns the updated
// feedback item with its new vote count.
func (h *FeedbackHandler) Vote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	// The body is optional: an empty body counts as an anonymous vote.
	var req request.CreateVoteRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	feedback, err := h.store.AddVote(c.Request.Context(), id, req.UserID)
	switch {
	case errors.Is(err, store.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
	case errors.Is(err, store.ErrDuplicateVote):
		c.JSON(http.StatusConflict, gin.H{"error": "User has already voted on this feedback"})
	case err != nil:
		log.Printf("Failed to add vote for feedback %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add vote"})
	default:
		c.JSON(http.StatusCreated, feedback)
	}
}

// buildListFilter validates and converts the raw query into a store filter.
func buildListFilter(q request.ListFeedbackQuery) (store.ListFilter, error) {
	filter := store.ListFilter{Tag: q.Tag}

	if q.Type != "" {
		var t enum.FeedbackType
		if err := t.UnmarshalText([]byte(q.Type)); err != nil {
			return store.ListFilter{}, err
		}
		filter.Type = &t
	}
	if q.Sentiment != "" {
		var s enum.Sentiment
		if err := s.UnmarshalText([]byte(q.Sentiment)); err != nil {
			return store.ListFilter{}, err
		}
		filter.Sentiment = &s
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > maxPageSize {
		q.PageSize = maxPageSize
	}
	filter.Limit = q.PageSize
	filter.Offset = (q.Page - 1) * q.PageSize

	return filter, nil
}
