package request

import "feedback-agent/app/models/enum"

type CreateFeedbackRequest struct {
	Title          string            `json:"title" binding:"required"`
	Description    string            `json:"description"`
	Type           enum.FeedbackType `json:"type" binding:"required"`
	Tags           []string          `json:"tags"`
	Sentiment      *enum.Sentiment   `json:"sentiment,omitempty"`      // Optional, defaults to Neutral
	SentimentScore *float64          `json:"sentimentScore,omitempty"` // Optional
}

