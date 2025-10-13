package types

import "context"

type LLMProvider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type LLMResponse struct {
	Text string `json:"text"`
}
