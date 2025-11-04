package llm

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			panic("OPENAI_API_KEY environment variable is required")
		}
	}

	return &OpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

// ClassifyText uses OpenAI to classify feedback text
func (c *OpenAIClient) ClassifyText(ctx context.Context, title, description string) (string, error) {
	// Combine title and description
	fullText := title
	if description != "" {
		fullText += "\n\n" + description
	}

	// Create a structured prompt for classification
	prompt := `Analyze the following feedback and classify it into ONE of these categories:
- "bug" - if it describes a bug, error, issue, or broken functionality
- "feature" - if it's a feature request, enhancement, or new functionality suggestion
- "general" - if it's general feedback, praise, or doesn't fit the above categories

Respond with ONLY the category name (bug, feature, or general), nothing else.

Feedback:
` + fullText

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo, // or GPT4 for better accuracy
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.3, // Lower temperature for more consistent classification
			MaxTokens:   10,  // Only need the category name
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "general", nil // Default fallback
	}

	return resp.Choices[0].Message.Content, nil
}
