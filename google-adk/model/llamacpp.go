package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LlamaCppModel implements the Model interface for llama.cpp
type LlamaCppModel struct {
	model   string
	baseURL string
	client  *http.Client
}

// NewLlamaCppModel creates a new llama.cpp model
func NewLlamaCppModel(config Config) (*LlamaCppModel, error) {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:12434/v1"
	}

	return &LlamaCppModel{
		model:   config.Model,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

// Generate generates a response using llama.cpp
func (m *LlamaCppModel) Generate(ctx context.Context, prompt string) (*Result, error) {
	requestBody := map[string]interface{}{
		"model": m.model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 4096,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			// Log the error but don't fail the operation
			fmt.Printf("warning: failed to close response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp API error: %s - %s", resp.Status, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &Result{
		Text: response.Choices[0].Message.Content,
	}, nil
}
