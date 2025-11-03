package enum

import (
	"encoding/json"
	"fmt"
)

type Sentiment int

const (
	SentimentNeutral  Sentiment = 0
	SentimentPositive Sentiment = 1
	SentimentNegative Sentiment = 2
)

var sentimentIDs = map[Sentiment]string{
	SentimentNeutral:  "neutral",
	SentimentPositive: "positive",
	SentimentNegative: "negative",
}

var sentimentNames = map[string]Sentiment{
	"neutral":  SentimentNeutral,
	"positive": SentimentPositive,
	"negative": SentimentNegative,
}

func (s Sentiment) MarshalText() ([]byte, error) {
	if id, ok := sentimentIDs[s]; ok {
		return []byte(id), nil
	}
	return nil, fmt.Errorf("unknown Sentiment: %d", s)
}

func (s *Sentiment) UnmarshalText(text []byte) error {
	if v, ok := sentimentNames[string(text)]; ok {
		*s = v
		return nil
	}
	return fmt.Errorf("unknown Sentiment: %s", string(text))
}

func (s Sentiment) Name() string {
	if id, ok := sentimentIDs[s]; ok {
		return id
	}
	return "unknown"
}

func (s Sentiment) Category() string {
	switch s {
	case SentimentPositive:
		return "Kudos & Carrots"
	case SentimentNegative:
		return "Critiques"
	default:
		return "Context"
	}
}

// UnmarshalJSON allows Sentiment to be unmarshaled from JSON numbers
func (s *Sentiment) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as number first
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*s = Sentiment(num)
		// Validate it's a valid enum value
		if _, ok := sentimentIDs[*s]; !ok {
			return fmt.Errorf("unknown Sentiment: %d", num)
		}
		return nil
	}
	// Fallback to text unmarshaling
	return s.UnmarshalText(data)
}

// MarshalJSON allows Sentiment to be marshaled as JSON numbers
func (s Sentiment) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(s))
}


