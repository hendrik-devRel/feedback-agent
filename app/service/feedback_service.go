package service

import (
	"context"
	"errors"
	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
	"feedback-agent/app/models/request"
	"feedback-agent/app/repository"
)

// FeedbackService defines the interface for feedback business logic
type FeedbackService interface {
	CreateFeedback(ctx context.Context, req *request.CreateFeedbackRequest) (*entity.Feedback, error)
	GetFeedback(ctx context.Context, id int) (*entity.Feedback, error)
	GetAllFeedback(ctx context.Context) ([]*entity.Feedback, error)
}

type feedbackService struct {
	feedbackRepo repository.FeedbackRepository
}

// NewFeedbackService creates a new feedback service
func NewFeedbackService(feedbackRepo repository.FeedbackRepository) FeedbackService {
	return &feedbackService{
		feedbackRepo: feedbackRepo,
	}
}

// CreateFeedback creates a new feedback entry
func (s *feedbackService) CreateFeedback(ctx context.Context, req *request.CreateFeedbackRequest) (*entity.Feedback, error) {
	// Validate request
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	// Set default sentiment if not provided
	sentiment := enum.SentimentNeutral
	if req.Sentiment != nil {
		sentiment = *req.Sentiment
	}

	// Create feedback entity
	feedback := &entity.Feedback{
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		Tags:           req.Tags,
		Sentiment:      sentiment,
		SentimentScore: req.SentimentScore,
		Votes:          0,
	}

	// Save to repository
	createdFeedback, err := s.feedbackRepo.Create(ctx, feedback)
	if err != nil {
		return nil, err
	}

	return createdFeedback, nil
}

// GetFeedback retrieves a feedback entry by ID
func (s *feedbackService) GetFeedback(ctx context.Context, id int) (*entity.Feedback, error) {
	if id <= 0 {
		return nil, errors.New("invalid feedback ID")
	}

	return s.feedbackRepo.GetByID(ctx, id)
}

// GetAllFeedback retrieves all feedback entries
func (s *feedbackService) GetAllFeedback(ctx context.Context) ([]*entity.Feedback, error) {
	return s.feedbackRepo.GetAll(ctx)
}