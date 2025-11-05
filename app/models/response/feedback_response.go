package response

import (
	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
	"time"
)

// FeedbackResponse represents the API response for feedback
type FeedbackResponse struct {
	ID             int               `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Type           enum.FeedbackType `json:"type"`
	Tags           []string          `json:"tags"`
	Sentiment      enum.Sentiment    `json:"sentiment"`
	SentimentScore *float64          `json:"sentimentScore,omitempty"`
	Votes          int               `json:"votes"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

// FromEntity converts entity.Feedback to FeedbackResponse
func FromEntity(f *entity.Feedback) *FeedbackResponse {
	return &FeedbackResponse{
		ID:             f.ID,
		Title:          f.Title,
		Description:    f.Description,
		Type:           f.Type,
		Tags:           f.Tags,
		Sentiment:      f.Sentiment,
		SentimentScore: f.SentimentScore,
		Votes:          f.Votes,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}
}