package ollama

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	client *api.Client
	model  string
}

func NewClient(host, model string) (*OllamaClient, error) {
	if host == "" || model == "" {
		log.Println("OLLAMA_HOST or OLLAMA_MODEL not set, Ollama client is disabled.")
		return &OllamaClient{client: nil}, nil
	}
	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("could not parse OLLAMA_HOST url: %w", err)
	}
	return &OllamaClient{
		client: api.NewClient(hostURL, nil),
		model:  model,
	}, nil
}

func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("Ollama client is not configured")
	}
	req := &api.GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
	}
	var fullResponse string
	err := c.client.Generate(ctx, req, func(res api.GenerateResponse) error {
		fullResponse += res.Response
		return nil
	})
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("Ollama generation failed: %w", err)
	}
	return fullResponse, nil
}
