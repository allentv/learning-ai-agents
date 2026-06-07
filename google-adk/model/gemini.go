// Package model provides model implementations for various AI providers including Gemini and llama.cpp.
package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GeminiModel implements the Model interface for Google Gemini
type GeminiModel struct {
	model   string
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewGeminiModel creates a new Gemini model
func NewGeminiModel(config Config) (*GeminiModel, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	return &GeminiModel{
		model:   config.Model,
		apiKey:  config.APIKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

// Generate generates a response using Gemini
func (m *GeminiModel) Generate(ctx context.Context, prompt string) (*Result, error) {
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 1000,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", m.baseURL, m.model, m.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	return &Result{
		Text: result.Candidates[0].Content.Parts[0].Text,
	}, nil
}
