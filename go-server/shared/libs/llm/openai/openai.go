package openai

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct{}

func NewClient() *OpenAIClient {
	return &OpenAIClient{}
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string, apiKey string) (string, error) {
	log.Println(apiKey)

	requestClient := openai.NewClient(apiKey)
	resp, err := requestClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("OpenAI returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
