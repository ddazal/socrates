package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ddazal/socrates/internal/agent"
	"github.com/ddazal/socrates/internal/provider"
)

func main() {
	var model string
	var task string
	var reflections uint
	var debug bool

	flag.StringVar(&model, "model", "qwen2.5-coder:7b", "Model to use for chat")
	flag.StringVar(&task, "task", "Create a function to sum any number of integers", "The task for which the agent should generate code")
	flag.UintVar(&reflections, "reflections", 3, "Max number of reflections")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")

	flag.Parse()

	if strings.TrimSpace(task) == "" {
		log.Fatal("Task description cannot be empty.")
	}

	if reflections == 0 {
		log.Fatal("Reflections must be greater than zero.")
	}

	// Create the Ollama provider
	llm, err := provider.NewOllamaProvider()
	if err != nil {
		log.Fatal("Error creating Ollama provider: ", err)
	}

	agentInstance, err := agent.New(context.Background(), llm, provider.Config{
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
