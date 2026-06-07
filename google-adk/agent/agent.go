package agent

import (
	"context"
	"fmt"

	"github.com/learning-ai-agents/google-adk/model"
)

// Config represents agent configuration
type Config struct {
	Name         string
	Model        model.Config
	SystemPrompt string
}

// Agent represents an AI agent
type Agent struct {
	config Config
	model  model.Model
}

// New creates a new agent
func New(config Config) (*Agent, error) {
	m, err := model.New(config.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	return &Agent{
		config: config,
		model:  m,
	}, nil
}

// Request represents a request to the agent
type Request struct {
	Query string
}

// Response represents a response from the agent
type Response struct {
	Text string
}

// Run processes a request and returns a response
func (a *Agent) Run(ctx context.Context, req *Request) (*Response, error) {
	// Build the prompt with system prompt
	prompt := a.config.SystemPrompt + "\n\nUser: " + req.Query

	// Generate response using the model
	result, err := a.model.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	return &Response{
		Text: result.Text,
	}, nil
}
