package provider

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

type OpenAIProvider struct {
	client openai.Client
}

func NewOpenAIProvider() (*OpenAIProvider, error) {
	client := openai.NewClient() // Assumes API key is set in environment variable OPENAI_API_KEY
	return &OpenAIProvider{
		client: client,
	}, nil
}

func (o OpenAIProvider) Chat(ctx context.Context, model string, prompt string) (string, error) {
	res, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: shared.ChatModel(model),
	})
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}
	return res.Choices[0].Message.Content, nil
}

func (o OpenAIProvider) ValidateModel(ctx context.Context, model string) error {
	// OpenAI does not provide a direct API to validate models.
	// As a workaround, we can attempt to create a chat completion with the model.
	_, err := o.Chat(ctx, model, "Hello")
	return err
}
