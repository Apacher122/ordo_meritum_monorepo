package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	llmErrors "github.com/ordo_meritum/shared/libs/llm/errors"
	"google.golang.org/api/googleapi"
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
		return "", &llmErrors.LLMError{
			LLMProvider: "Gemini",
			Err:         llmErrors.ErrFailedToInit,
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	if instructions != "" {
		config.SystemInstruction = genai.NewContentFromText(instructions, genai.RoleModel)
	}

	if schema != nil {
		var finalSchema *genai.Schema
		switch v := schema.(type) {
		case *genai.Schema:
			finalSchema = v
		case map[string]any:
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", &llmErrors.LLMError{
					LLMProvider: "Gemini",
					Err:         llmErrors.ErrUnsupportedSchema,
				}
			}
			var s genai.Schema
			if err := json.Unmarshal(bytes, &s); err != nil {
				return "", &llmErrors.LLMError{
					LLMProvider:     "Gemini",
					Err:             llmErrors.ErrUnsupportedSchema,
					ProviderMessage: err.Error(),
				}
			}
			finalSchema = &s
		default:
			return "", &llmErrors.LLMError{
				LLMProvider: "Gemini",
				Err:         llmErrors.ErrUnsupportedSchema,
			}
		}
		config.ResponseSchema = finalSchema
	}

	return c.callWithRetries(ctx, client, prompt, config, 3, 1*time.Second)
}

func (c *GeminiClient) callWithRetries(
	ctx context.Context,
	client *genai.Client,
	prompt string,
	config *genai.GenerateContentConfig,
	maxRetries int,
	baseDelay time.Duration,
) (string, error) {
	var err error
	var resp *genai.GenerateContentResponse
	for i := range maxRetries {
		resp, err = client.Models.GenerateContent(
			ctx,
			c.model,
			genai.Text(prompt),
			config,
		)

		if err == nil {
			break
		}

		log.Printf("Gemini API call failed (attempt %d/%d): %v", i+1, maxRetries, err)

		var googleErr *googleapi.Error
		if errors.As(err, &googleErr) {
			if googleErr.Code == http.StatusTooManyRequests || googleErr.Code >= 500 {
				delay := baseDelay * time.Duration(1<<i)
				log.Printf("Retrying in %v...", delay)
				time.Sleep(delay)
				continue
			}
		}

		return "", &llmErrors.LLMError{
			LLMProvider:     "Gemini",
			Err:             llmErrors.ErrNoContent,
			ProviderMessage: err.Error(),
		}
	}

	if err != nil {
		return "", &llmErrors.LLMError{
			LLMProvider:     "Gemini",
			Err:             llmErrors.ErrNoContent,
			ProviderMessage: err.Error(),
		}
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", &llmErrors.LLMError{
			LLMProvider: "Gemini",
			Err:         llmErrors.ErrNoContent,
		}
	}
	return resp.Candidates[0].Content.Parts[0].Text, nil

}
