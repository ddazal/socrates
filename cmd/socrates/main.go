package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ddazal/socrates/internal/agent"
	llmProvider "github.com/ddazal/socrates/internal/provider"
)

func main() {
	var model string
	var task string
	var reflections uint
	var debug bool
	var provider string

	const defaultModel = "qwen2.5-coder:7b"

	flag.StringVar(&model, "model", defaultModel, "Model to use for chat")
	flag.StringVar(&task, "task", "Create a function to sum any number of integers", "The task for which the agent should generate code")
	flag.UintVar(&reflections, "reflections", 3, "Max number of reflections")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.StringVar(&provider, "provider", "ollama", "LLM provider to use (ollama, openai)")

	flag.Parse()

	if strings.TrimSpace(task) == "" {
		log.Fatal("Task description cannot be empty.")
	}

	if reflections == 0 {
		log.Fatal("Reflections must be greater than zero.")
	}

	llm, err := getProvider(provider)
	if err != nil {
		log.Fatal("Error getting provider: ", err)
	}

	agentInstance, err := agent.New(context.Background(), llm, llmProvider.Config{
		Debug:          debug,
		MaxReflections: reflections,
		Model:          model,
	})
	if err != nil {
		log.Fatal("Error creating agent: ", err)
	}

	answer, err := agentInstance.Run(task)
	if err != nil {
		log.Fatal("Error during chat: ", err)
	}

	fmt.Println(answer)
}

func getProvider(provider string) (llmProvider.Provider, error) {
	var llm llmProvider.Provider
	var err error

	switch strings.ToLower(provider) {
	case "ollama":
		llm, err = llmProvider.NewOllamaProvider()
		if err != nil {
			return nil, fmt.Errorf("error creating Ollama provider: %w", err)
		}
	case "openai":
		llm, err = llmProvider.NewOpenAIProvider()
		if err != nil {
			return nil, fmt.Errorf("error creating OpenAI provider: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported provider: %q", provider)
	}

	return llm, nil
}
