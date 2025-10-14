package llm

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ordo_meritum/shared/libs/llm/cohere"
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
	case "cohere":
		return cohere.NewClient(), nil
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

	re := regexp.MustCompile("(?s)```json\\s*(\\{.*?\\})\\s*```")
	match := re.FindStringSubmatch(clean)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	clean = strings.ReplaceAll(clean, "```json", "")
	clean = strings.ReplaceAll(clean, "```", "")
	clean = regexp.MustCompile(`(?i)^here\s+is.*?schema[:\s]*`).ReplaceAllString(clean, "")
	return strings.TrimSpace(clean)
}
