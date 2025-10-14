package cohere

import (
	"context"
	"errors"
	"fmt"
	"log"

	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
)

type CohereClient struct {
	model string
}

func NewClient() *CohereClient {
	return &CohereClient{
		model: "command-a-03-2025",
	}
}

func (c *CohereClient) Generate(
	ctx context.Context,
	instructions string,
	prompt string,
	schema any,
	apiKey string,
) (string, error) {
	if apiKey == "" {
		return "", errors.New("error: Cohere API key not provided")
	}

	var response *cohere.ResponseFormatV2
	if schema != nil {
		var schemaMap *cohere.JsonResponseFormatV2
		log.Printf("schema type: %s", schema)
		switch v := schema.(type) {
		case *cohere.JsonResponseFormatV2:
			schemaMap = v
			response = &cohere.ResponseFormatV2{
				Type:       "josn_object",
				JsonObject: schemaMap,
			}
		}
	}
	requestClient := cohereclient.NewClient(cohereclient.WithToken(apiKey))

	resp, err := requestClient.V2.Chat(
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
			ResponseFormat: response,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error: Cohere chat generation failed: %w", err)
	}

	if resp.Message == nil || len(resp.Message.Content) == 0 {
		return "", errors.New("error: Cohere returned no content")
	}

	var assistantReply string
	for _, item := range resp.Message.Content {
		if item.Text != nil {
			assistantReply += item.Text.Text
		}
	}

	if assistantReply == "" {
		return "", errors.New("error: Cohere assistant message was empty")
	}

	log.Println(assistantReply)
	return assistantReply, nil
}
