package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"

	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

// LlamaCppLLMWrapper wraps a llama.cpp model to implement the ADK's model.LLM interface
type LlamaCppLLMWrapper struct {
	model   string
	baseURL string
	client  *http.Client
}

// NewLlamaCppLLMWrapper creates a new llama.cpp LLM wrapper
func NewLlamaCppLLMWrapper(config Config) (*LlamaCppLLMWrapper, error) {
	if config.Model == "" {
		return nil, fmt.Errorf("model name is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:12434/v1"
	}

	return &LlamaCppLLMWrapper{
		model:   config.Model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

// Name returns the model name
func (w *LlamaCppLLMWrapper) Name() string {
	return w.model
}

// GenerateContent implements the model.LLM interface for llama.cpp
func (w *LlamaCppLLMWrapper) GenerateContent(ctx context.Context, req *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		// Convert the request to llama.cpp format
		messages := []map[string]interface{}{}

		// Add messages from contents
		if len(req.Contents) > 0 {
			for _, content := range req.Contents {
				for _, part := range content.Parts {
					// Check if part has text content
					if part.Text != "" {
						// Determine role based on content
						role := "user"
						if content.Role != "" {
							role = content.Role
						}
						messages = append(messages, map[string]interface{}{
							"role":    role,
							"content": part.Text,
						})
					}
				}
			}
		}

		requestBody := map[string]interface{}{
			"model":      w.model,
			"messages":   messages,
			"max_tokens": 4096,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			yield(nil, fmt.Errorf("failed to marshal request: %w", err))
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", w.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			yield(nil, fmt.Errorf("failed to create request: %w", err))
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json")

		resp, err := w.client.Do(httpReq)
		if err != nil {
			yield(nil, fmt.Errorf("failed to send request: %w", err))
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				// Log the error but don't fail the operation
				fmt.Printf("warning: failed to close response body: %v\n", cerr)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			yield(nil, fmt.Errorf("llama.cpp API error: %s - %s", resp.Status, string(body)))
			return
		}

		var response struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			yield(nil, fmt.Errorf("failed to decode response: %w", err))
			return
		}

		if len(response.Choices) == 0 {
			yield(nil, fmt.Errorf("no choices in response"))
			return
		}

		// Create LLMResponse
		llmResponse := &model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{
						Text: response.Choices[0].Message.Content,
					},
				},
			},
		}

		yield(llmResponse, nil)
	}
}
