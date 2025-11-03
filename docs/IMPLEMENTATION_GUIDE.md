# Implementation Guide - Step by Step

## Phase 1: Foundation (Current → Clean Architecture)

### Step 1: Create Request/Response Models
**File**: `app/models/request/create_feedback.go`
```go
package request

import "feedback-agent/app/models/enum"

type CreateFeedbackRequest struct {
	Title          string            `json:"title" binding:"required"`
	Description    string            `json:"description"`
	Type           enum.FeedbackType `json:"type" binding:"required"`
	Tags           []string          `json:"tags"`
	Sentiment      *enum.Sentiment   `json:"sentiment,omitempty"`
	SentimentScore *float64          `json:"sentimentScore,omitempty"`
}
```

**File**: `app/models/response/feedback_response.go`
```go
package response

import "feedback-agent/app/models/entity"

// Just use entity.Feedback directly, or create this if you need different JSON structure
type FeedbackResponse entity.Feedback
```

### Step 2: Create Repository Interface
**File**: `app/repository/feedback_repository.go`
```go
package repository

import (
	"context"
	"feedback-agent/app/models/entity"
)

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error)
	GetByID(ctx context.Context, id int) (*entity.Feedback, error)
	GetAll(ctx context.Context) ([]*entity.Feedback, error)
	Update(ctx context.Context, feedback *entity.Feedback) error
}

// Implementation
type feedbackRepository struct {
	db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error) {
	// SQL implementation here
	// Return created feedback with ID, timestamps
}

// ... other methods
```

### Step 3: Create Service Layer
**File**: `app/service/feedback_service.go`
```go
package service

import (
	"context"
	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
	"feedback-agent/app/models/request"
	"feedback-agent/app/repository"
)

type FeedbackService interface {
	CreateFeedback(ctx context.Context, req *request.CreateFeedbackRequest) (*entity.Feedback, error)
	GetFeedback(ctx context.Context, id int) (*entity.Feedback, error)
	GetAllFeedback(ctx context.Context) ([]*entity.Feedback, error)
}

type feedbackService struct {
	feedbackRepo repository.FeedbackRepository
}

func NewFeedbackService(feedbackRepo repository.FeedbackRepository) FeedbackService {
	return &feedbackService{
		feedbackRepo: feedbackRepo,
	}
}

func (s *feedbackService) CreateFeedback(ctx context.Context, req *request.CreateFeedbackRequest) (*entity.Feedback, error) {
	// 1. Validate request
	// 2. Set defaults (sentiment = Neutral if not provided)
	// 3. Create Feedback entity
	// 4. Call repository.Create()
	// 5. Return created feedback
}
```

### Step 4: Create Handler
**File**: `app/handler/feedback_handler.go`
```go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"feedback-agent/app/models/request"
	"feedback-agent/app/service"
)

type FeedbackHandler struct {
	feedbackService service.FeedbackService
}

func NewFeedbackHandler(feedbackService service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	var req request.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := h.feedbackService.CreateFeedback(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}
```

### Step 5: Wire Everything in main.go
**File**: `app/main.go`
```go
package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"feedback-agent/app/handler"
	"feedback-agent/app/infrastructure/database"
	"feedback-agent/app/repository"
	"feedback-agent/app/router"
	"feedback-agent/app/service"
)

func main() {
	// 1. Setup database
	db := database.NewPostgresConnection()
	defer db.Close()

	// 2. Create repositories
	feedbackRepo := repository.NewFeedbackRepository(db)

	// 3. Create services
	feedbackService := service.NewFeedbackService(feedbackRepo)

	// 4. Create handlers
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)

	// 5. Setup router
	r := router.NewRouter(feedbackHandler)

	// 6. Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

## Building Order (Recommended)

1. **First**: Create directory structure
2. **Second**: Move database code → `repository/feedback_repository.go`
3. **Third**: Create `service/feedback_service.go` (calls repository)
4. **Fourth**: Create `handler/feedback_handler.go` (calls service)
5. **Fifth**: Create `router/router.go` (sets up routes)
6. **Sixth**: Refactor `main.go` to wire everything

This way, you build from the bottom up (repository → service → handler), testing each layer as you go.

