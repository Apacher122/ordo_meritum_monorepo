package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ordo_meritum/shared/contexts"
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

// Generate generates content based on the given prompt and instructions.
//
// It will make a single request to the Gemini API with the given
// configuration and return the response if successful. If the
// request fails, it will retry up to 3 times with an exponential
// backoff starting at 1 second. If all retries fail, it will
// return an error.
//
// The response will be unmarshalled into a string and returned. If the
// response is not of type "application/json", an error will be returned.
//
// The schema parameter can be used to specify the expected response schema.
// If the schema is not provided, the response will be unmarshalled into a
// map[string]interface{}.
//
// The apiKey parameter is used to authenticate the request. If the key is
// invalid, an error will be returned.
func (c *GeminiClient) Generate(
	ctx context.Context,
	instructions string,
	prompt string,
	schema any,
) (string, error) {
	userCtx, _ := contexts.FromContext(ctx)
	log.Printf("User context: %+v", userCtx)
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: userCtx.ApiKey,
	})
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
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
