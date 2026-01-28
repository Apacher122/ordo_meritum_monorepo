package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ordo_meritum/shared/contexts"
	llmErrors "github.com/ordo_meritum/shared/libs/llm/errors"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	lg "github.com/ordo_meritum/shared/utils/logger"
	"google.golang.org/genai"
)

var Temperature = float32(0.0)
var TopP = float32(0.1)
var TopK = float32(1)
var Service = "gemini"

type GeminiClient struct {
	model string
}

func NewClient() *GeminiClient {
	return &GeminiClient{
		model: "gemini-3-flash-preview",
	}
}

// Generate generates content using the Gemini AI model.
//
// It takes in a context, the instructions file path, the prompt data to format into the instructions,
// the type of schema to use, and the target to unmarshal the Gemini response into.
//
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the Gemini provider cannot be found, it will return an error with ErrorCode set to ERR_LLM_NO_CONTENT.
// If the prompt cannot be formatted, it will return an error with ErrorCode set to ERR_LLM_PROMPT_FORMATTING.
// If the instructions file cannot be read, it will return an error with ErrorCode set to ERR_LLM_INSTRUCTION_FORMATTING.
// If the Gemini response cannot be unmarshalled, it will return an error with ErrorCode set to ERR_LLM_MALFORMED_RESPONSE.
func (c *GeminiClient) Generate(
	ctx context.Context,
	instructions string,
	prompt string,
	schema any,
) (string, error) {
	userCtx, _ := contexts.FromContext(ctx)
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

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
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

	return c.callWithRetries(ctx, client, prompt, config, 2, 30*time.Second)
}

// callWithRetries makes a call to the Gemini API with the given prompt and config.
// If the call fails, it will retry up to maxRetries times with a delay of baseDelay * 2^i seconds.
// If the call still fails after maxRetries attempts, it will return an error with ErrorCode set to ERR_LLM_NO_CONTENT.
// If the Gemini response does not contain any content, it will return an error with ErrorCode set to ERR_LLM_NO_CONTENT.
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
		lg.InfoLoggerType{Service: &Service, Message: "Attempting to generate content"}.InfoLog()
		resp, err = client.Models.GenerateContent(
			ctx,
			c.model,
			genai.Text(prompt),
			config,
		)

		if err == nil {
			lg.InfoLoggerType{Service: &Service, Message: "Successfully generated content"}.InfoLog()
			break
		}

		log.Printf("Gemini API call failed (attempt %d/%d): %v", i+1, maxRetries, err)
		lg.ErrorLoggerType{Service: &Service, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: fmt.Errorf("gemini API call failed (attempt %d/%d): %v", i+1, maxRetries, err)}.ErrorLog()

		var googleErr genai.APIError
		if errors.As(err, &googleErr) {
			if googleErr.Code == http.StatusTooManyRequests || googleErr.Code >= 500 {
				delay := baseDelay * time.Duration(1<<i)
				log.Printf("Retrying in %v...", delay)
				time.Sleep(delay)
				continue
			}
		}

		return "", &llmErrors.LLMError{
			LLMProvider: "Gemini",
			Err:         llmErrors.ErrNoContent,
		}
	}

	if err != nil {
		return "", &llmErrors.LLMError{
			LLMProvider: "Gemini",
			Err:         llmErrors.ErrNoContent,
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
