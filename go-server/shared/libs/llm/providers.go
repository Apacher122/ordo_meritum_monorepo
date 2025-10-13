package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ordo_meritum/shared/libs/llm/gemini"
)

type LLMProvider interface {
	Generate(ctx context.Context, instructions string, prompt string, schema any, apiKey string) (string, error)
}

func GetProvider(llm string) (LLMProvider, error) {
	switch llm {
	// case "openai":
	// 	return openai.NewClient(), nil
	// case "groq":
	// 	return groq.NewClient(), nil
	// case "ollama":
	// 	host := os.Getenv("OLLAMA_HOST")
	// 	model := os.Getenv("OLLAMA_MODEL")
	// 	return ollama.NewClient(host, model)
	// case "cohere":
	// 	return cohere.NewClient(), nil
	case "anthropic":
		return nil, fmt.Errorf("unsupported LLM provider: %s", llm)
	case "gemini":
		return gemini.NewClient(), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", llm)
	}
}

func FormatLLMResponse(raw string) string {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimSuffix(clean, "```")
	return strings.TrimSpace(clean)
}

func ParseJSON[T any](rawJSON string) (*T, error) {
	var out T
	if err := json.Unmarshal([]byte(rawJSON), &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Gemini JSON: %w", err)
	}
	return &out, nil
}
