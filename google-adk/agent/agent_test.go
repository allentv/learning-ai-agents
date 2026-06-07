package agent_test

import (
	"context"
	"fmt"
	"iter"
	"testing"

	"github.com/learning-ai-agents/google-adk/agent"
	localmodel "github.com/learning-ai-agents/google-adk/model"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/genai"
)

// mockLLM is a mock implementation of the LLM interface for testing
type mockLLM struct {
	responseText string
	shouldError  bool
	errorMsg     string
}

func (m *mockLLM) Name() string {
	return "mock-llm"
}

func (m *mockLLM) GenerateContent(ctx context.Context, req *adkmodel.LLMRequest, _ bool) iter.Seq2[*adkmodel.LLMResponse, error] {
	return func(yield func(*adkmodel.LLMResponse, error) bool) {
		if m.shouldError {
			yield(nil, fmt.Errorf("%s", m.errorMsg))
			return
		}

		response := &adkmodel.LLMResponse{
			Content: &genai.Content{
				Role: "model",
				Parts: []*genai.Part{
					{
						Text: m.responseText,
					},
				},
			},
		}
		yield(response, nil)
	}
}

func TestAgent_New(t *testing.T) {
	config := agent.Config{
		Name: "test-agent",
		ModelConfig: localmodel.Config{
			Model:   "test-model.gguf",
			BaseURL: "http://localhost:12434/v1",
		},
		SystemPrompt: "You are a helpful assistant.",
	}

	a, err := agent.New(config)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if a == nil {
		t.Fatal("New() returned nil agent")
	}
}

func TestAgent_NewWithEmptyModel(t *testing.T) {
	config := agent.Config{
		Name: "test-agent",
		ModelConfig: localmodel.Config{
			Model:   "",
			BaseURL: "http://localhost:12434/v1",
		},
		SystemPrompt: "You are a helpful assistant.",
	}

	_, err := agent.New(config)
	if err == nil {
		t.Fatal("New() should return error for empty model")
	}
}

func TestAgent_Run(t *testing.T) {
	// Create a mock LLM
	mockModel := &mockLLM{
		responseText: "This is a test response",
		shouldError:  false,
	}

	// Create agent with mock model
	config := agent.Config{
		Name:         "test-agent",
		SystemPrompt: "You are a helpful assistant.",
	}
	a, err := agent.NewWithModel(config, mockModel)
	if err != nil {
		t.Fatalf("NewWithModel() returned error: %v", err)
	}

	ctx := context.Background()
	req := &agent.Request{
		Query: "Hello, world!",
	}

	resp, err := a.Run(ctx, req)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Run() returned nil response")
	}

	expectedText := "This is a test response"
	if resp.Text != expectedText {
		t.Errorf("Expected response text '%s', got '%s'", expectedText, resp.Text)
	}
}

func TestAgent_RunModelError(t *testing.T) {
	// Create a mock LLM that returns an error
	mockModel := &mockLLM{
		shouldError: true,
		errorMsg:    "model error",
	}

	// Create agent with mock model
	config := agent.Config{
		Name:         "test-agent",
		SystemPrompt: "You are a helpful assistant.",
	}
	a, err := agent.NewWithModel(config, mockModel)
	if err != nil {
		t.Fatalf("NewWithModel() returned error: %v", err)
	}

	ctx := context.Background()
	req := &agent.Request{
		Query: "Hello, world!",
	}

	_, err = a.Run(ctx, req)
	if err == nil {
		t.Fatal("Run() should return error when model fails")
	}

	if err.Error() != "failed to generate response: model error" {
		t.Errorf("Expected error message 'failed to generate response: model error', got '%s'", err.Error())
	}
}

func TestAgent_RunEmptyResponse(t *testing.T) {
	// Create a mock LLM that returns empty response
	mockModel := &mockLLM{
		responseText: "",
		shouldError:  false,
	}

	// Create agent with mock model
	config := agent.Config{
		Name:         "test-agent",
		SystemPrompt: "You are a helpful assistant.",
	}
	a, err := agent.NewWithModel(config, mockModel)
	if err != nil {
		t.Fatalf("NewWithModel() returned error: %v", err)
	}

	ctx := context.Background()
	req := &agent.Request{
		Query: "Hello, world!",
	}

	resp, err := a.Run(ctx, req)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Run() returned nil response")
	}

	// Empty response should still be returned
	if resp.Text != "" {
		t.Errorf("Expected empty response text, got '%s'", resp.Text)
	}
}

func TestAgent_RequestStruct(t *testing.T) {
	req := &agent.Request{
		Query: "Test query",
	}

	if req.Query != "Test query" {
		t.Errorf("Expected query 'Test query', got '%s'", req.Query)
	}
}

func TestAgent_ResponseStruct(t *testing.T) {
	resp := &agent.Response{
		Text: "Test response",
	}

	if resp.Text != "Test response" {
		t.Errorf("Expected text 'Test response', got '%s'", resp.Text)
	}
}

func TestAgent_ConfigStruct(t *testing.T) {
	config := agent.Config{
		Name: "test-agent",
		ModelConfig: localmodel.Config{
			Model:   "test-model.gguf",
			BaseURL: "http://localhost:12434/v1",
		},
		SystemPrompt: "You are a helpful assistant.",
	}

	if config.Name != "test-agent" {
		t.Errorf("Expected name 'test-agent', got '%s'", config.Name)
	}

	if config.ModelConfig.Model != "test-model.gguf" {
		t.Errorf("Expected model 'test-model.gguf', got '%s'", config.ModelConfig.Model)
	}

	if config.SystemPrompt != "You are a helpful assistant." {
		t.Errorf("Expected system prompt 'You are a helpful assistant.', got '%s'", config.SystemPrompt)
	}
}
