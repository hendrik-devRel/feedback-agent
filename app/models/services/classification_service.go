package service

import (
	"context"
	"log"
	"strings"
	"feedback-agent/app/models/enum"
	"feedback-agent/app/infrastructure/llm"
)

type ClassificationService struct {
	llmClient *llm.OpenAIClient
}

func NewClassificationService(llmClient *llm.OpenAIClient) *ClassificationService {
	return &ClassificationService{
		llmClient: llmClient,
	}
}

// ClassifyFeedback uses LLM to analyze title and description
func (s *ClassificationService) ClassifyFeedback(ctx context.Context, title, description string) enum.FeedbackType {
	// Call LLM for classification
	result, err := s.llmClient.ClassifyText(ctx, title, description)
	if err != nil {
		log.Printf("LLM classification error: %v, defaulting to general", err)
		return enum.FeedbackGeneral
	}

	// Parse LLM response (normalize to lowercase and trim whitespace)
	result = strings.ToLower(strings.TrimSpace(result))

	// Map LLM response to enum
	switch result {
	case "bug":
		return enum.FeedbackBug
	case "feature":
		return enum.FeedbackFeature
	case "general":
		return enum.FeedbackGeneral
	default:
		// If LLM returns something unexpected, default to general
		log.Printf("Unexpected classification result: %s, defaulting to general", result)
		return enum.FeedbackGeneral
	}
}