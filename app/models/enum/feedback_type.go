package enum

import (
	"encoding/json"
	"fmt"
)

type FeedbackType int

const (
	FeedbackBug     FeedbackType = 0
	FeedbackFeature FeedbackType = 1
	FeedbackGeneral FeedbackType = 2
)

var feedbackTypeIDs = map[FeedbackType]string{
	FeedbackBug:     "bug",
	FeedbackFeature: "feature",
	FeedbackGeneral: "general",
}

var feedbackTypeNames = map[string]FeedbackType{
	"bug":     FeedbackBug,
	"feature": FeedbackFeature,
	"general": FeedbackGeneral,
}

func (s FeedbackType) MarshalText() ([]byte, error) {
	if id, ok := feedbackTypeIDs[s]; ok {
		return []byte(id), nil
	}
	return nil, fmt.Errorf("unknown FeedbackType: %d", s)
}

func (s *FeedbackType) UnmarshalText(text []byte) error {
	if v, ok := feedbackTypeNames[string(text)]; ok {
		*s = v
		return nil
	}
	return fmt.Errorf("unknown FeedbackType: %s", string(text))
}

func (s FeedbackType) Name() string {
	if id, ok := feedbackTypeIDs[s]; ok {
		return id
	}
	return "unknown"
}

// UnmarshalJSON allows FeedbackType to be unmarshaled from JSON numbers
func (s *FeedbackType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as number first
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*s = FeedbackType(num)
		// Validate it's a valid enum value
		if _, ok := feedbackTypeIDs[*s]; !ok {
			return fmt.Errorf("unknown FeedbackType: %d", num)
		}
		return nil
	}
	// Fallback to text unmarshaling
	return s.UnmarshalText(data)
}

// MarshalJSON allows FeedbackType to be marshaled as JSON numbers
func (s FeedbackType) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(s))
}


