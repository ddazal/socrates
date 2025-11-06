// Package provider defines the interface for LLM providers and their implementations.
package provider

import "context"

// Provider is an interface for LLM providers
type Provider interface {
	// Chat sends a prompt to the LLM and returns the response
	Chat(ctx context.Context, model string, prompt string) (string, error)

	// ValidateModel checks if the specified model is available
	ValidateModel(ctx context.Context, model string) error
}

type Config struct {
	Debug          bool
	MaxReflections uint
	Model          string
}
