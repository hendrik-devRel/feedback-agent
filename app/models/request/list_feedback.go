package request

// ListFeedbackQuery holds the query parameters accepted by GET /api/feedback.
// Type, sentiment, and tag are optional filters; page and pageSize control
// pagination.
type ListFeedbackQuery struct {
	Type      string `form:"type"`
	Sentiment string `form:"sentiment"`
	Tag       string `form:"tag"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"pageSize,default=20"`
}
