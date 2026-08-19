package request

// CreateVoteRequest is the body for POST /api/feedback/:id/votes.
// UserID is optional: anonymous votes are allowed and are not deduplicated.
type CreateVoteRequest struct {
	UserID *int `json:"userId,omitempty"`
}
