// Package agent provides AI agent functionality for processing requests and generating responses.
package agent

import (
	"context"
	"fmt"

	localmodel "github.com/learning-ai-agents/google-adk/model"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/genai"
)

// Config represents agent configuration
type Config struct {
	Name         string
	ModelConfig  localmodel.Config
	SystemPrompt string
}

// Agent represents an AI agent
type Agent struct {
	config Config
	model  adkmodel.LLM
}

// New creates a new agent
func New(config Config) (*Agent, error) {
	m, err := localmodel.NewLlamaCppLLMWrapper(config.ModelConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	return &Agent{
		config: config,
		model:  m,
	}, nil
}

// NewWithModel creates a new agent with a custom model (for testing)
func NewWithModel(config Config, model adkmodel.LLM) (*Agent, error) {
	if model == nil {
		return nil, fmt.Errorf("model cannot be nil")
	}

	return &Agent{
		config: config,
		model:  model,
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

	// Create ADK LLM request
	llmReq := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// Generate response using the model
	var result *adkmodel.LLMResponse

	// Use the GenerateContent method which returns an iterator
	for resp, genErr := range a.model.GenerateContent(ctx, llmReq, false) {
		if genErr != nil {
			return nil, fmt.Errorf("failed to generate response: %w", genErr)
		}
		result = resp
		break // Take the first response
	}

	if result == nil || result.Content == nil || len(result.Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	return &Response{
		Text: result.Content.Parts[0].Text,
	}, nil
}
