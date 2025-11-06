package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/ollama/ollama/api"
)

// OllamaProvider implements the Provider interface using Ollama
type OllamaProvider struct {
	client *api.Client
}

// NewOllamaProvider creates a new Ollama provider using environment configuration
func NewOllamaProvider() (*OllamaProvider, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return &OllamaProvider{
		client: client,
	}, nil
}

// Chat sends a prompt to the Ollama model and returns the response
func (p *OllamaProvider) Chat(ctx context.Context, model string, prompt string) (string, error) {
	var answer string
	stream := false

	err := p.client.Chat(ctx, &api.ChatRequest{
		Model: model,
		Messages: []api.Message{
			{Role: "user", Content: prompt},
		},
		Stream: &stream,
	}, func(response api.ChatResponse) error {
		answer = response.Message.Content
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("chat call failed: %w", err)
	}

	if answer == "" {
		return "", errors.New("no response content received from model")
	}

	return answer, nil
}

// ValidateModel checks if the specified model is available in Ollama
func (p *OllamaProvider) ValidateModel(ctx context.Context, model string) error {
	results, err := p.client.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	for _, m := range results.Models {
		if m.Model == model {
			return nil
		}
	}

	return fmt.Errorf("model %q is not available - run 'ollama pull %s' to download it", model, model)
}
