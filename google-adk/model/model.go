package model

import (
	"context"
	"fmt"
)

// Config represents model configuration
type Config struct {
	Model   string
	APIKey  string
	BaseURL string
}

// Result represents a model generation result
type Result struct {
	Text string
}

// Model is an interface for AI models
type Model interface {
	Generate(ctx context.Context, prompt string) (*Result, error)
}

// New creates a new model based on configuration
func New(config Config) (Model, error) {
	// Determine which model to create based on the API key and base URL
	// If API key is set and looks like an OpenAI key, use OpenAI
	// If API key is set and looks like a Gemini key, use Gemini
	// Otherwise, use llama.cpp

	if config.APIKey != "" && config.BaseURL != "" && config.BaseURL != "http://localhost:12434/v1" {
		// Likely OpenAI with custom base URL
		return NewOpenAIModel(config)
	}
	if config.APIKey != "" && len(config.APIKey) > 20 && config.APIKey[:3] == "sk-" {
		// Likely OpenAI API key
		return NewOpenAIModel(config)
	}
	if config.APIKey != "" && config.APIKey != "llamacpp-dummy-key" {
		// Likely Gemini API key
		return NewGeminiModel(config)
	}
	// Use llama.cpp
	return NewLlamaCppModel(config)
}

// MockModel is a simple mock implementation for demonstration
type MockModel struct {
	model string
}

// Generate generates a response (mock implementation)
func (m *MockModel) Generate(_ context.Context, _ string) (*Result, error) {
	// Simple mock response
	response := fmt.Sprintf("Using model %s: I received your message and I'm processing it.", m.model)
	return &Result{Text: response}, nil
}
