package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type GeminiClient struct {
	model string
}

func NewClient() *GeminiClient {
	return &GeminiClient{
		model: "gemini-2.5-pro",
	}
}

func (c *GeminiClient) Generate(
	ctx context.Context,
	instructions string,
	prompt string,
	schema any,
	apiKey string,
) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", fmt.Errorf("could not create request-scoped Gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType:  "application/json",
		SystemInstruction: genai.NewContentFromText(instructions, genai.RoleModel),
	}

	if schema != nil {
		var finalSchema *genai.Schema
		switch v := schema.(type) {
		case *genai.Schema:
			finalSchema = v
		case map[string]any:
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("failed to marshal map schema: %w", err)
			}
			var s genai.Schema
			if err := json.Unmarshal(bytes, &s); err != nil {
				return "", fmt.Errorf("failed to parse map schema into genai.Schema: %w", err)
			}
			finalSchema = &s
		default:
			return "", fmt.Errorf("unsupported schema type: %T", v)
		}

		config.ResponseSchema = finalSchema
	}
	resp, err := client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return "", fmt.Errorf("error: Gemini content generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("error: Gemini returned no content")
	}
	log.Println(resp.Candidates[0].Content.Parts[0].Text)
	return resp.Candidates[0].Content.Parts[0].Text, nil
}
