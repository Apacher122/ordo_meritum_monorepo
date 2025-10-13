package cohere

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
)

type CohereClient struct{}

func NewClient() *CohereClient {
	return &CohereClient{}
}

func (c *CohereClient) Generate(
	ctx context.Context,
	instructions string,
	prompt string,
	schema cohere.ResponseFormatV2,
	apiKey string,
) (*request.JobPostingEvent, error) {
	if apiKey == "" {
		return nil, errors.New("error: Cohere API key not found in request context")
	}

	requestClient := cohereclient.NewClient(cohereclient.WithToken(apiKey))
	resp, _ := requestClient.V2.Chat(
		ctx,
		&cohere.V2ChatRequest{
			Model: "command-a-03-2025",
			Messages: cohere.ChatMessages{
				{
					Role: "system",
					System: &cohere.SystemMessageV2{Content: &cohere.SystemMessageV2Content{
						String: instructions,
					}},
				},
				{
					Role: "user",
					User: &cohere.UserMessageV2{Content: &cohere.UserMessageV2Content{
						String: prompt,
					}},
				},
			},
			ResponseFormat: &schema,
		},
	)

	if resp.Message == nil {
		return nil, errors.New("no assistant message returned from Cohere")
	}

	var assistantReply string
	for _, item := range resp.Message.Content {
		if item.Type == "text" && item.Text != nil {
			assistantReply += item.Text.Text
		}
	}

	if assistantReply == "" {
		return nil, errors.New("assistant message empty")
	}

	var result request.JobPostingEvent
	if err := json.Unmarshal([]byte(assistantReply), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &result, nil
}
