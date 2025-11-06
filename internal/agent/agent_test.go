package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/ddazal/socrates/internal/provider"
	"github.com/stretchr/testify/assert"
)

// mockProvider is a mock implementation of the Provider interface for testing
type mockProvider struct {
	validateModelFunc func(ctx context.Context, model string) error
	chatFunc          func(ctx context.Context, model string, prompt string) (string, error)
}

func (m *mockProvider) ValidateModel(ctx context.Context, model string) error {
	if m.validateModelFunc != nil {
		return m.validateModelFunc(ctx, model)
	}
	return nil
}

func (m *mockProvider) Chat(ctx context.Context, model string, prompt string) (string, error) {
	if m.chatFunc != nil {
		return m.chatFunc(ctx, model, prompt)
	}
	return "mock response", nil
}

func TestNewAgent_New(t *testing.T) {
	mock := &mockProvider{
		validateModelFunc: func(ctx context.Context, model string) error {
			return nil
		},
	}

	agent, err := New(context.Background(), mock, provider.Config{
		Debug:          true,
		MaxReflections: 3,
		Model:          "test-model",
	})
	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, "test-model", agent.model)
	assert.Equal(t, uint(3), agent.maxReflections)
	assert.Equal(t, true, agent.debug)
}

func TestNewAgent_Validate(t *testing.T) {
	mock := &mockProvider{
		validateModelFunc: func(ctx context.Context, model string) error {
			return assert.AnError
		},
	}

	agent, err := New(context.Background(), mock, provider.Config{
		Model: "invalid-model",
	})
	assert.Error(t, err)
	assert.Nil(t, agent)
	assert.ErrorAs(t, err, &assert.AnError)
}

func TestAgent_Run(t *testing.T) {
	mock := &mockProvider{
		validateModelFunc: func(ctx context.Context, model string) error {
			return nil
		},
		chatFunc: func(ctx context.Context, model string, prompt string) (string, error) {
			// Return initial code
			if strings.Contains(prompt, "Generate Go code") {
				return "func Hello() string { return \"hello\" }", nil
			}
			// Return critique
			if strings.Contains(prompt, "Analyze the Go code") {
				return "No changes needed.", nil
			}
			// Return refined code
			return "func Hello() string { return \"hello\" }", nil
		},
	}

	agent, err := New(context.Background(), mock, provider.Config{
		MaxReflections: 3,
		Model:          "test-model",
	})
	assert.NoError(t, err)

	result, err := agent.Run("Create a Hello function")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Contains(t, result, "func Hello() string { return \"hello\" }")
}

func TestAgent_Run_Reflections(t *testing.T) {
	callCount := 0
	mock := &mockProvider{
		validateModelFunc: func(ctx context.Context, model string) error {
			return nil
		},
		chatFunc: func(ctx context.Context, model string, prompt string) (string, error) {
			callCount++
			// Return initial code
			if strings.Contains(prompt, "Generate Go code") {
				return "func Add(a, b int) int { return a + b }", nil
			}
			// Return critique - first round says needs changes
			if strings.Contains(prompt, "Analyze the Go code") {
				if callCount <= 2 {
					return "Add error handling for edge cases.", nil
				}
				return "No changes needed.", nil
			}
			// Return refined code
			if strings.Contains(prompt, "Revise the Go code") {
				return "func Add(a, b int) (int, error) { return a + b, nil }", nil
			}
			return "mock response", nil
		},
	}

	agent, err := New(context.Background(), mock, provider.Config{
		MaxReflections: 3,
		Model:          "test-model",
	})
	assert.NoError(t, err)

	result, err := agent.Run("Create an Add function")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// generate (1), reflect (2), refine (3), reflect(4)
	assert.Equal(t, 4, callCount, "Expected 4 calls but got %d", callCount)
}
