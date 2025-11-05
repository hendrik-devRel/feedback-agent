package repository

import (
	"context"
	"database/sql"
	"feedback-agent/app/models/entity"

	"github.com/lib/pq"
)

// FeedbackRepository defines the interface for feedback data operations
type FeedbackRepository interface {
	Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error)
	GetByID(ctx context.Context, id int) (*entity.Feedback, error)
	GetAll(ctx context.Context) ([]*entity.Feedback, error)
	Update(ctx context.Context, feedback *entity.Feedback) error
}

type feedbackRepository struct {
	db *sql.DB
}

// NewFeedbackRepository creates a new feedback repository
func NewFeedbackRepository(db *sql.DB) FeedbackRepository {
	return &feedbackRepository{db: db}
}

// Create inserts a new feedback record into the database
func (r *feedbackRepository) Create(ctx context.Context, feedback *entity.Feedback) (*entity.Feedback, error) {
	query := `
		INSERT INTO feedback (title, description, type, tags, sentiment, sentiment_score, votes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, votes, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		feedback.Title,
		feedback.Description,
		int(feedback.Type),
		pq.Array(feedback.Tags),
		int(feedback.Sentiment),
		feedback.SentimentScore,
		0, // Votes always start at 0
	).Scan(&feedback.ID, &feedback.Votes, &feedback.CreatedAt, &feedback.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// GetByID retrieves a feedback record by its ID
func (r *feedbackRepository) GetByID(ctx context.Context, id int) (*entity.Feedback, error) {
	query := `
		SELECT id, title, description, type, tags, sentiment, sentiment_score, votes, created_at, updated_at
		FROM feedback
		WHERE id = $1
	`

	var feedback entity.Feedback
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&feedback.ID,
		&feedback.Title,
		&feedback.Description,
		&feedback.Type,
		pq.Array(&feedback.Tags),
		&feedback.Sentiment,
		&feedback.SentimentScore,
		&feedback.Votes,
		&feedback.CreatedAt,
		&feedback.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &feedback, nil
}

// GetAll retrieves all feedback records
func (r *feedbackRepository) GetAll(ctx context.Context) ([]*entity.Feedback, error) {
	query := `
		SELECT id, title, description, type, tags, sentiment, sentiment_score, votes, created_at, updated_at
		FROM feedback
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []*entity.Feedback
	for rows.Next() {
		var feedback entity.Feedback
		err := rows.Scan(
			&feedback.ID,
			&feedback.Title,
			&feedback.Description,
			&feedback.Type,
			pq.Array(&feedback.Tags),
			&feedback.Sentiment,
			&feedback.SentimentScore,
			&feedback.Votes,
			&feedback.CreatedAt,
			&feedback.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, &feedback)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feedbacks, nil
}

// Update updates an existing feedback record
func (r *feedbackRepository) Update(ctx context.Context, feedback *entity.Feedback) error {
	query := `
		UPDATE feedback
		SET title = $1, description = $2, type = $3, tags = $4, 
		    sentiment = $5, sentiment_score = $6, votes = $7, updated_at = NOW()
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		feedback.Title,
		feedback.Description,
		int(feedback.Type),
		pq.Array(feedback.Tags),
		int(feedback.Sentiment),
		feedback.SentimentScore,
		feedback.Votes,
		feedback.ID,
	)

	return err
}