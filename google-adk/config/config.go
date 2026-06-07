// Package config provides application configuration management using environment variables.
package config

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration
type Config struct {
	ModelProvider string `envconfig:"MODEL_PROVIDER" default:"gemini"`
	OpenAIKey     string `envconfig:"OPENAI_API_KEY"`
	OpenAIModel   string `envconfig:"OPENAI_MODEL" default:"gpt-4o-mini"`
	GeminiKey     string `envconfig:"GEMINI_API_KEY"`
	GeminiModel   string `envconfig:"GEMINI_MODEL" default:"gemini-3.1-flash-lite"`
	LLamacppURL   string `envconfig:"LLAMACPP_URL" default:"http://localhost:12434/v1"`
	LLamacppModel string `envconfig:"LLAMACPP_MODEL" default:"granite-4.0-h-micro-UD-Q4_K_XL.gguf"`
	LogLevel      string `envconfig:"LOG_LEVEL" default:"INFO"`
	AppName       string `envconfig:"APP_NAME" default:"Google ADK Go Project"`
	AppVersion    string `envconfig:"APP_VERSION" default:"0.1.0"`
}

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	err := loadEnvFile(".env")
	if err != nil {
		// Log but don't fail if .env doesn't exist
		log.Printf("Note: .env file not found: %v", err)
	}

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
