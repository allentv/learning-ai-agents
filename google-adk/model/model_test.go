package model_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	localmodel "github.com/learning-ai-agents/google-adk/model"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/genai"
)

func TestConfig_Struct(t *testing.T) {
	config := localmodel.Config{
		Model:   "test-model.gguf",
		APIKey:  "test-key",
		BaseURL: "http://test:1234/v1",
	}

	if config.Model != "test-model.gguf" {
		t.Errorf("Expected Model 'test-model.gguf', got '%s'", config.Model)
	}

	if config.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got '%s'", config.APIKey)
	}

	if config.BaseURL != "http://test:1234/v1" {
		t.Errorf("Expected BaseURL 'http://test:1234/v1', got '%s'", config.BaseURL)
	}
}

func TestLlamaCppLLMWrapper_New(t *testing.T) {
	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: "http://test:1234/v1",
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	if wrapper == nil {
		t.Fatal("NewLlamaCppLLMWrapper() returned nil wrapper")
	}
}

func TestLlamaCppLLMWrapper_NewWithEmptyBaseURL(t *testing.T) {
	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: "",
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	if wrapper == nil {
		t.Fatal("NewLlamaCppLLMWrapper() returned nil wrapper")
	}

	// Check that default baseURL is set
	if wrapper.Name() != "test-model.gguf" {
		t.Errorf("Expected Name 'test-model.gguf', got '%s'", wrapper.Name())
	}
}

func TestLlamaCppLLMWrapper_Name(t *testing.T) {
	config := localmodel.Config{
		Model:   "my-model.gguf",
		BaseURL: "http://test:1234/v1",
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	if wrapper.Name() != "my-model.gguf" {
		t.Errorf("Expected Name 'my-model.gguf', got '%s'", wrapper.Name())
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_Success(t *testing.T) {
	// Create a test HTTP server
	expectedResponse := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"content": "This is a test response from llama.cpp",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Decode request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Verify model in request
		if reqBody["model"] != "test-model.gguf" {
			t.Errorf("Expected model 'test-model.gguf', got '%v'", reqBody["model"])
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "Hello, world!",
					},
				},
			},
		},
	}

	var result *adkmodel.LLMResponse
	var resultErr error

	for resp, err := range wrapper.GenerateContent(ctx, req, false) {
		result = resp
		resultErr = err
		break
	}

	if resultErr != nil {
		t.Fatalf("GenerateContent() returned error: %v", resultErr)
	}

	if result == nil {
		t.Fatal("GenerateContent() returned nil response")
	}

	if result.Content == nil {
		t.Fatal("GenerateContent() returned nil content")
	}

	if len(result.Content.Parts) == 0 {
		t.Fatal("GenerateContent() returned empty parts")
	}

	expectedText := "This is a test response from llama.cpp"
	if result.Content.Parts[0].Text != expectedText {
		t.Errorf("Expected text '%s', got '%s'", expectedText, result.Content.Parts[0].Text)
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_HTTPError(t *testing.T) {
	// Create a test HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "Hello, world!",
					},
				},
			},
		},
	}

	var resultErr error
	for _, err := range wrapper.GenerateContent(ctx, req, false) {
		resultErr = err
		break
	}

	if resultErr == nil {
		t.Fatal("GenerateContent() should return error for HTTP error response")
	}

	if resultErr.Error() == "" {
		t.Error("GenerateContent() returned empty error message")
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_EmptyResponse(t *testing.T) {
	// Create a test HTTP server that returns empty choices
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "Hello, world!",
					},
				},
			},
		},
	}

	var resultErr error
	for _, err := range wrapper.GenerateContent(ctx, req, false) {
		resultErr = err
		break
	}

	if resultErr == nil {
		t.Fatal("GenerateContent() should return error for empty choices")
	}

	expectedError := "no choices in response"
	if resultErr.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, resultErr.Error())
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_InvalidJSON(t *testing.T) {
	// Create a test HTTP server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "Hello, world!",
					},
				},
			},
		},
	}

	var resultErr error
	for _, err := range wrapper.GenerateContent(ctx, req, false) {
		resultErr = err
		break
	}

	if resultErr == nil {
		t.Fatal("GenerateContent() should return error for invalid JSON")
	}

	if resultErr.Error() == "" {
		t.Error("GenerateContent() returned empty error message")
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_MultipleContents(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Verify messages array has multiple entries
		messages, ok := reqBody["messages"].([]interface{})
		if !ok {
			t.Fatal("messages field is not an array")
		}

		if len(messages) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(messages))
		}

		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "Response to multiple messages",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "First message",
					},
				},
			},
			{
				Role: "assistant",
				Parts: []*genai.Part{
					{
						Text: "Second message",
					},
				},
			},
		},
	}

	var result *adkmodel.LLMResponse
	var resultErr error

	for resp, err := range wrapper.GenerateContent(ctx, req, false) {
		result = resp
		resultErr = err
		break
	}

	if resultErr != nil {
		t.Fatalf("GenerateContent() returned error: %v", resultErr)
	}

	if result == nil {
		t.Fatal("GenerateContent() returned nil response")
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_NilContents(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Verify messages array is empty when contents is nil
		messages, ok := reqBody["messages"].([]interface{})
		if !ok {
			t.Fatal("messages field is not an array")
		}

		if len(messages) != 0 {
			t.Errorf("Expected 0 messages for nil contents, got %d", len(messages))
		}

		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "Response to empty messages",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{},
	}

	var result *adkmodel.LLMResponse
	var resultErr error

	for resp, err := range wrapper.GenerateContent(ctx, req, false) {
		result = resp
		resultErr = err
		break
	}

	if resultErr != nil {
		t.Fatalf("GenerateContent() returned error: %v", resultErr)
	}

	if result == nil {
		t.Fatal("GenerateContent() returned nil response")
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_EmptyTextPart(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Verify messages array doesn't include empty text parts
		messages, ok := reqBody["messages"].([]interface{})
		if !ok {
			t.Fatal("messages field is not an array")
		}

		// Should have 0 messages since text is empty
		if len(messages) != 0 {
			t.Errorf("Expected 0 messages for empty text, got %d", len(messages))
		}

		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "Response",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	ctx := context.Background()
	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "", // Empty text
					},
				},
			},
		},
	}

	var result *adkmodel.LLMResponse
	var resultErr error

	for resp, err := range wrapper.GenerateContent(ctx, req, false) {
		result = resp
		resultErr = err
		break
	}

	if resultErr != nil {
		t.Fatalf("GenerateContent() returned error: %v", resultErr)
	}

	if result == nil {
		t.Fatal("GenerateContent() returned nil response")
	}
}

func TestLlamaCppLLMWrapper_GenerateContent_ContextCancellation(t *testing.T) {
	// Create a test HTTP server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		<-r.Context().Done()
		// Context was cancelled
	}))
	defer server.Close()

	config := localmodel.Config{
		Model:   "test-model.gguf",
		BaseURL: server.URL,
	}

	wrapper, err := localmodel.NewLlamaCppLLMWrapper(config)
	if err != nil {
		t.Fatalf("NewLlamaCppLLMWrapper() returned error: %v", err)
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := &adkmodel.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{
						Text: "Hello, world!",
					},
				},
			},
		},
	}

	var resultErr error
	for _, err := range wrapper.GenerateContent(ctx, req, false) {
		resultErr = err
		break
	}

	// Should return an error due to cancelled context
	if resultErr == nil {
		t.Fatal("GenerateContent() should return error for cancelled context")
	}
}
