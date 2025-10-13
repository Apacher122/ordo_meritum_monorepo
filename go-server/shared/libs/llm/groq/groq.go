package groq

import (
	"context"
	"errors"
	"fmt"

	"github.com/ordo_meritum/shared/middleware"
	"github.com/sashabaranov/go-openai"
)

type GroqClient struct{}

func NewClient() *GroqClient {
	return &GroqClient{}
}

func (c *GroqClient) Generate(ctx context.Context, prompt string) (string, error) {
	apiKey, ok := ctx.Value(middleware.APIKeyContextKey).(string)
	if !ok || apiKey == "" {
		return "", errors.New("Groq API key not found in request context")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"
	requestClient := openai.NewClientWithConfig(config)

	resp, err := requestClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "llama3-8b-8192",
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("Groq chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("Groq returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
