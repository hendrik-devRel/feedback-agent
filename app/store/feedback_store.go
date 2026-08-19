package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"feedback-agent/app/models/entity"
	"feedback-agent/app/models/enum"
)

// ErrNotFound is returned when a requested row does not exist.
var ErrNotFound = errors.New("not found")

// ErrDuplicateVote is returned when an authenticated user votes twice
// on the same feedback item.
var ErrDuplicateVote = errors.New("duplicate vote")

// FeedbackStore wraps all database access for feedback and votes.
type FeedbackStore struct {
	db *sql.DB
}

func NewFeedbackStore(db *sql.DB) *FeedbackStore {
	return &FeedbackStore{db: db}
}

// ListFilter describes the optional filters and pagination for ListFeedback.
type ListFilter struct {
	Type      *enum.FeedbackType
	Sentiment *enum.Sentiment
	Tag       string
	Limit     int
	Offset    int
}

// Create inserts a new feedback row and fills in the generated fields
// (ID, Votes, CreatedAt, UpdatedAt) on the returned entity.
func (s *FeedbackStore) Create(ctx context.Context, f entity.Feedback) (entity.Feedback, error) {
	const query = `
		INSERT INTO feedback (title, description, type, tags, sentiment, sentiment_score, votes)
		VALUES ($1, $2, $3, $4, $5, $6, 0)
		RETURNING id, votes, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		f.Title,
		f.Description,
		int(f.Type),
		pq.Array(f.Tags),
		int(f.Sentiment),
		f.SentimentScore,
	).Scan(&f.ID, &f.Votes, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return entity.Feedback{}, fmt.Errorf("insert feedback: %w", err)
	}
	return f, nil
}

// GetByID returns a single feedback item, or ErrNotFound.
func (s *FeedbackStore) GetByID(ctx context.Context, id int) (entity.Feedback, error) {
	const query = `
		SELECT id, title, description, type, tags, sentiment, sentiment_score, votes, created_at, updated_at
		FROM feedback
		WHERE id = $1
	`
	f, err := scanFeedback(s.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Feedback{}, ErrNotFound
	}
	if err != nil {
		return entity.Feedback{}, fmt.Errorf("get feedback %d: %w", id, err)
	}
	return f, nil
}

// List returns feedback rows matching the filter, newest first, plus the
// total count of matching rows (ignoring pagination) for pagination UIs.
func (s *FeedbackStore) List(ctx context.Context, filter ListFilter) ([]entity.Feedback, int, error) {
	where, args := buildWhere(filter)

	var total int
	countQuery := "SELECT COUNT(*) FROM feedback" + where
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count feedback: %w", err)
	}

	listQuery := fmt.Sprintf(`
		SELECT id, title, description, type, tags, sentiment, sentiment_score, votes, created_at, updated_at
		FROM feedback%s
		ORDER BY created_at DESC, id DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list feedback: %w", err)
	}
	defer rows.Close()

	items := make([]entity.Feedback, 0, filter.Limit)
	for rows.Next() {
		f, err := scanFeedback(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan feedback: %w", err)
		}
		items = append(items, f)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate feedback: %w", err)
	}
	return items, total, nil
}

// AddVote records a vote for a feedback item and increments its cached
// vote counter atomically. userID may be nil for anonymous votes; the
// partial unique index on (feedback_id, user_id) rejects duplicate votes
// from the same authenticated user.
func (s *FeedbackStore) AddVote(ctx context.Context, feedbackID int, userID *int) (entity.Feedback, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Feedback{}, fmt.Errorf("begin vote tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO votes (feedback_id, user_id) VALUES ($1, $2)`,
		feedbackID, userID,
	); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return entity.Feedback{}, ErrDuplicateVote
			case "foreign_key_violation":
				return entity.Feedback{}, ErrNotFound
			}
		}
		return entity.Feedback{}, fmt.Errorf("insert vote: %w", err)
	}

	const bump = `
		UPDATE feedback
		SET votes = votes + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, description, type, tags, sentiment, sentiment_score, votes, created_at, updated_at
	`
	f, err := scanFeedback(tx.QueryRowContext(ctx, bump, feedbackID))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Feedback{}, ErrNotFound
	}
	if err != nil {
		return entity.Feedback{}, fmt.Errorf("bump vote count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return entity.Feedback{}, fmt.Errorf("commit vote tx: %w", err)
	}
	return f, nil
}

// buildWhere assembles the WHERE clause and its positional args for List.
func buildWhere(filter ListFilter) (string, []any) {
	var clauses []string
	var args []any

	if filter.Type != nil {
		args = append(args, int(*filter.Type))
		clauses = append(clauses, fmt.Sprintf("type = $%d", len(args)))
	}
	if filter.Sentiment != nil {
		args = append(args, int(*filter.Sentiment))
		clauses = append(clauses, fmt.Sprintf("sentiment = $%d", len(args)))
	}
	if filter.Tag != "" {
		args = append(args, filter.Tag)
		clauses = append(clauses, fmt.Sprintf("$%d = ANY(tags)", len(args)))
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

// scanner covers both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanFeedback(row scanner) (entity.Feedback, error) {
	var f entity.Feedback
	var typ, sentiment int
	err := row.Scan(
		&f.ID,
		&f.Title,
		&f.Description,
		&typ,
		pq.Array(&f.Tags),
		&sentiment,
		&f.SentimentScore,
		&f.Votes,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		return entity.Feedback{}, err
	}
	f.Type = enum.FeedbackType(typ)
	f.Sentiment = enum.Sentiment(sentiment)
	return f, nil
}
