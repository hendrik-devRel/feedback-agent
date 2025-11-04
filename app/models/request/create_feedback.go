package request

import "feedback-agent/app/models/enum"

type CreateFeedbackRequest struct {
	Title          string            `json:"title" binding:"required"`
	Description    string            `json:"description"`
	Type           *enum.FeedbackType `json:"type,omitempty"`
	Tags           []string          `json:"tags,omitempty"`           // Optional
	Sentiment      *enum.Sentiment   `json:"sentiment,omitempty"`      // Optional, defaults to Neutral
	SentimentScore *float64          `json:"sentimentScore,omitempty"` // Optional
}
