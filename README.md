# Socrates

**Socrates** is a command-line agent written in Go that demonstrates the reflection pattern using LLMs.

It generates Go code for a task, critiques its own output, and iteratively improves it.

The agent is **provider-agnostic**, meaning it can work with different LLM providers through a simple interface.


## Features

- Reflection Loop: Generate, critique, and Revise (up to N times)
- Go-Focused: Outputs idiomatic Go for your custom task
- CLI-First: Easy to run with task, model, and debug flags
- Debug Mode: Watch how the agent evolves its output step-by-step
- Provider-Agnostic: Easily swap between LLM providers (Ollama, OpenAI, etc.)

## Requirements

- Go 1.20+
- **For Ollama**: Ollama installed and running with a compatible model (e.g., `qwen2.5-coder:7b`)

## Usage

```bash
# Run with default task
go run cmd/socrates/main.go

# Run with custom task
go run cmd/socrates/main.go --task "Write an isPrime function in Go"

# Enable debug logs to view each step
go run cmd/socrates/main.go --task "Write a calculator CLI in Go" --debug

# Use a different model and increase reflections
go run cmd/socrates/main.go --task "Generate a REST API using Go and Gin" --model "codellama:7b" --reflections 5

# Build and install
go build -o socrates cmd/socrates/main.go
./socrates --task "Your task here"
```

## 🔧 CLI Flags
| Flag            | Description                            | Default             |
| --------------- | -------------------------------------- | ------------------- |
| `--task`        | Task for the agent to solve            | Create a function to sum any number of integers |
| `--model`       | Ollama model to use                    | `qwen2.5-coder:7b`  |
| `--reflections` | Max number of critique/revision cycles | `3`                 |
| `--debug`       | Show verbose step-by-step logs         | `false`             |