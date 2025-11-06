// Package agent provides the core reflection-based code generation agent.
// It implements a generate-critique-refine loop using LLM providers.
package agent

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ddazal/socrates/internal/provider"
)

type agent struct {
	provider       provider.Provider
	maxReflections uint
	model          string
	debug          bool
}

func New(ctx context.Context, p provider.Provider, config provider.Config) (*agent, error) {
	// Validate that the model is available
	if err := p.ValidateModel(ctx, config.Model); err != nil {
		return nil, err
	}

	if config.Debug {
		fmt.Println("Debug mode is enabled. Detailed logs will be printed.")
	}

	return &agent{
		provider:       p,
		debug:          config.Debug,
		maxReflections: config.MaxReflections,
		model:          config.Model,
	}, nil
}

func (a *agent) Run(task string) (string, error) {
	a.logf("🚀 Generating initial code for task:\n%s\n", task)

	code, err := a.generateInitialOutput(task)
	if err != nil {
		return "", fmt.Errorf("error generating initial output: %w", err)
	}

	a.logf("🤖 Initial output:\n%s\n", code)

	for i := 0; i < int(a.maxReflections); i++ {
		a.logf("🔄 Reflection round %d...", i+1)

		critique, err := a.reflect(task, code)
		if err != nil {
			return "", fmt.Errorf("reflection step %d failed: %w", i+1, err)
		}

		a.logf("🔍 Critique:\n%s\n", critique)

		if strings.Contains(strings.ToLower(critique), "no changes needed") {
			a.logf("✅ No changes needed. Final code is ready.")
			break
		}

		code, err = a.refine(task, code, critique)
		if err != nil {
			return "", fmt.Errorf("refinement step %d failed: %w", i+1, err)
		}

		a.logf("✅ Updated code:\n%s\n", code)
	}

	return code, nil
}

func (a *agent) chat(prompt string) (string, error) {
	return a.provider.Chat(context.Background(), a.model, prompt)
}

func (a *agent) generateInitialOutput(task string) (string, error) {
	prompt := fmt.Sprintf(`
	Generate Go code that fulfills the task described between <task></task> tags. Return only the code, without explanations or comments.

    <task>%s</task>
	`, task)
	return a.chat(prompt)
}

func (a *agent) reflect(task, code string) (string, error) {
	prompt := fmt.Sprintf(`
	Analyze the Go code for correctness, potential bugs, edge cases, and adherence to idiomatic Go best practices. Be concise and specific. If the code is correct and requires no changes, respond with: "No changes needed."

	<task>%s</task>
	<code>%s</code>
	`, task, code)

	return a.chat(prompt)
}

func (a *agent) refine(task, code, critique string) (string, error) {
	prompt := fmt.Sprintf(`
	Revise the Go code between <code></code> to address the critique in <critique></critique>, while ensuring the code still fulfills the task in <task></task>. Respond with the improved code only, no explanations or comments.

	<task>%s</task>
	<code>%s</code>
	<critique>%s</critique>
	`, task, code, critique)

	return a.chat(prompt)
}

func (a *agent) logf(format string, args ...any) {
	if a.debug {
		log.Printf(format, args...)
	}
}
