package cohere

import (
	"context"
	"fmt"
	"log"

	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	"github.com/ordo_meritum/shared/contexts"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
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
) (string, *error_messages.ErrorBody) {
	userCtx, _ := contexts.FromContext(ctx)

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
	requestClient := cohereclient.NewClient(cohereclient.WithToken(userCtx.ApiKey))

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
		return "", &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_NO_CONTENT, ErrMsg: fmt.Errorf("error: Cohere chat generation failed: %w", err)}
	}

	if resp.Message == nil || len(resp.Message.Content) == 0 {
		return "", &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_NO_CONTENT, ErrMsg: fmt.Errorf("error: Cohere returned no content")}
	}

	var assistantReply string
	for _, item := range resp.Message.Content {
		if item.Text != nil {
			assistantReply += item.Text.Text
		}
	}

	if assistantReply == "" {
		return "", &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_NO_CONTENT, ErrMsg: fmt.Errorf("error: Cohere assistant message was empty")}
	}

	log.Println(assistantReply)
	return assistantReply, nil
}
