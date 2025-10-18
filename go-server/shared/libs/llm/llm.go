package llm

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ordo_meritum/shared/libs/llm/providers/cohere"
	"github.com/ordo_meritum/shared/libs/llm/providers/gemini"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
)

type LLMProvider interface {
	Generate(ctx context.Context, instructions string, prompt string, schema any) (string, *error_messages.ErrorBody)
}

// GetProvider returns a new LLMProvider based on the given LLM provider name.
// If the given LLM provider is not supported, it returns an error.
// Supported LLM providers are "cohere", "gemini", and "anthropic".
// The "anthropic" provider is not supported and will return an error.
func GetProvider(llm string) (LLMProvider, *error_messages.ErrorBody) {
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
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_INVALID_PROVIDER, ErrMsg: fmt.Errorf("anthropic is not yet implemented")}
	case "gemini":
		return gemini.NewClient(), nil
	default:
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_INVALID_PROVIDER, ErrMsg: fmt.Errorf("wnsupported LLM somehow request: %s", llm)}
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
