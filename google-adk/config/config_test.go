package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/learning-ai-agents/google-adk/config"
)

func TestConfig_LoadDefaults(t *testing.T) {
	// Clear environment variables
	_ = os.Unsetenv("MODEL_PROVIDER")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GEMINI_MODEL")
	_ = os.Unsetenv("LLAMACPP_URL")
	_ = os.Unsetenv("LLAMACPP_MODEL")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Check default values
	if cfg.ModelProvider != "gemini" {
		t.Errorf("Expected ModelProvider to be 'gemini', got '%s'", cfg.ModelProvider)
	}
	if cfg.GeminiModel != "gemini-3.1-flash-lite" {
		t.Errorf("Expected GeminiModel to be 'gemini-3.1-flash-lite', got '%s'", cfg.GeminiModel)
	}
	if cfg.LLamacppURL != "http://localhost:12434/v1" {
		t.Errorf("Expected LLamacppURL to be 'http://localhost:12434/v1', got '%s'", cfg.LLamacppURL)
	}
	if cfg.LLamacppModel != "granite-4.0-h-micro-UD-Q4_K_XL.gguf" {
		t.Errorf("Expected LLamacppModel to be 'granite-4.0-h-micro-UD-Q4_K_XL.gguf', got '%s'", cfg.LLamacppModel)
	}
	if cfg.LogLevel != "INFO" {
		t.Errorf("Expected LogLevel to be 'INFO', got '%s'", cfg.LogLevel)
	}
	if cfg.AppName != "Google ADK Go Project" {
		t.Errorf("Expected AppName to be 'Google ADK Go Project', got '%s'", cfg.AppName)
	}
	if cfg.AppVersion != "0.1.0" {
		t.Errorf("Expected AppVersion to be '0.1.0', got '%s'", cfg.AppVersion)
	}
}

func TestConfig_LoadFromEnvironment(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("MODEL_PROVIDER", "llamacpp")
	_ = os.Setenv("GEMINI_API_KEY", "test-key")
	_ = os.Setenv("GEMINI_MODEL", "test-model")
	_ = os.Setenv("LLAMACPP_URL", "http://test:1234/v1")
	_ = os.Setenv("LLAMACPP_MODEL", "test-model.gguf")
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("APP_NAME", "Test App")
	_ = os.Setenv("APP_VERSION", "1.0.0")

	// Clear .env file if it exists
	_ = os.Remove(".env")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Check values from environment
	if cfg.ModelProvider != "llamacpp" {
		t.Errorf("Expected ModelProvider to be 'llamacpp', got '%s'", cfg.ModelProvider)
	}
	if cfg.GeminiKey != "test-key" {
		t.Errorf("Expected GeminiKey to be 'test-key', got '%s'", cfg.GeminiKey)
	}
	if cfg.GeminiModel != "test-model" {
		t.Errorf("Expected GeminiModel to be 'test-model', got '%s'", cfg.GeminiModel)
	}
	if cfg.LLamacppURL != "http://test:1234/v1" {
		t.Errorf("Expected LLamacppURL to be 'http://test:1234/v1', got '%s'", cfg.LLamacppURL)
	}
	if cfg.LLamacppModel != "test-model.gguf" {
		t.Errorf("Expected LLamacppModel to be 'test-model.gguf', got '%s'", cfg.LLamacppModel)
	}
	if cfg.LogLevel != "DEBUG" {
		t.Errorf("Expected LogLevel to be 'DEBUG', got '%s'", cfg.LogLevel)
	}
	if cfg.AppName != "Test App" {
		t.Errorf("Expected AppName to be 'Test App', got '%s'", cfg.AppName)
	}
	if cfg.AppVersion != "1.0.0" {
		t.Errorf("Expected AppVersion to be '1.0.0', got '%s'", cfg.AppVersion)
	}

	// Clean up
	_ = os.Unsetenv("MODEL_PROVIDER")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GEMINI_MODEL")
	_ = os.Unsetenv("LLAMACPP_URL")
	_ = os.Unsetenv("LLAMACPP_MODEL")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
}

func TestConfig_LoadFromEnvFile(t *testing.T) {
	// Create a temporary .env file
	envContent := `MODEL_PROVIDER=llamacpp
GEMINI_API_KEY=env-file-key
GEMINI_MODEL=env-file-model
LLAMACPP_URL=http://env-test:5678/v1
LLAMACPP_MODEL=env-file-model.gguf
LOG_LEVEL=DEBUG
APP_NAME=Env File App
APP_VERSION=2.0.0`

	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envFilePath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// Clear environment variables
	_ = os.Unsetenv("MODEL_PROVIDER")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GEMINI_MODEL")
	_ = os.Unsetenv("LLAMACPP_URL")
	_ = os.Unsetenv("LLAMACPP_MODEL")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Check values from .env file
	if cfg.ModelProvider != "llamacpp" {
		t.Errorf("Expected ModelProvider to be 'llamacpp', got '%s'", cfg.ModelProvider)
	}
	if cfg.GeminiKey != "env-file-key" {
		t.Errorf("Expected GeminiKey to be 'env-file-key', got '%s'", cfg.GeminiKey)
	}
	if cfg.GeminiModel != "env-file-model" {
		t.Errorf("Expected GeminiModel to be 'env-file-model', got '%s'", cfg.GeminiModel)
	}
	if cfg.LLamacppURL != "http://env-test:5678/v1" {
		t.Errorf("Expected LLamacppURL to be 'http://env-test:5678/v1', got '%s'", cfg.LLamacppURL)
	}
	if cfg.LLamacppModel != "env-file-model.gguf" {
		t.Errorf("Expected LLamacppModel to be 'env-file-model.gguf', got '%s'", cfg.LLamacppModel)
	}
	if cfg.LogLevel != "DEBUG" {
		t.Errorf("Expected LogLevel to be 'DEBUG', got '%s'", cfg.LogLevel)
	}
	if cfg.AppName != "Env File App" {
		t.Errorf("Expected AppName to be 'Env File App', got '%s'", cfg.AppName)
	}
	if cfg.AppVersion != "2.0.0" {
		t.Errorf("Expected AppVersion to be '2.0.0', got '%s'", cfg.AppVersion)
	}
}

func TestConfig_LoadWithQuotedValues(t *testing.T) {
	// Create a temporary .env file with quoted values
	envContent := `MODEL_PROVIDER="llamacpp"
GEMINI_API_KEY='quoted-key'
GEMINI_MODEL="quoted-model"
LLAMACPP_URL='http://quoted:9999/v1'`

	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envFilePath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// Clear environment variables
	_ = os.Unsetenv("MODEL_PROVIDER")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GEMINI_MODEL")
	_ = os.Unsetenv("LLAMACPP_URL")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Check that quotes are removed
	if cfg.ModelProvider != "llamacpp" {
		t.Errorf("Expected ModelProvider to be 'llamacpp', got '%s'", cfg.ModelProvider)
	}
	if cfg.GeminiKey != "quoted-key" {
		t.Errorf("Expected GeminiKey to be 'quoted-key', got '%s'", cfg.GeminiKey)
	}
	if cfg.GeminiModel != "quoted-model" {
		t.Errorf("Expected GeminiModel to be 'quoted-model', got '%s'", cfg.GeminiModel)
	}
	if cfg.LLamacppURL != "http://quoted:9999/v1" {
		t.Errorf("Expected LLamacppURL to be 'http://quoted:9999/v1', got '%s'", cfg.LLamacppURL)
	}
}

func TestConfig_LoadWithCommentsAndEmptyLines(t *testing.T) {
	// Create a temporary .env file with comments and empty lines
	envContent := `# This is a comment
MODEL_PROVIDER=llamacpp

# Another comment
GEMINI_API_KEY=test-key

GEMINI_MODEL=test-model
# Final comment`

	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envFilePath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// Clear environment variables
	_ = os.Unsetenv("MODEL_PROVIDER")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GEMINI_MODEL")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Check values are loaded correctly
	if cfg.ModelProvider != "llamacpp" {
		t.Errorf("Expected ModelProvider to be 'llamacpp', got '%s'", cfg.ModelProvider)
	}
	if cfg.GeminiKey != "test-key" {
		t.Errorf("Expected GeminiKey to be 'test-key', got '%s'", cfg.GeminiKey)
	}
	if cfg.GeminiModel != "test-model" {
		t.Errorf("Expected GeminiModel to be 'test-model', got '%s'", cfg.GeminiModel)
	}
}

func TestConfig_LoadNonExistentEnvFile(t *testing.T) {
	// Change to a directory without .env file
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// Clear environment variables
	_ = os.Unsetenv("MODEL_PROVIDER")
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GEMINI_MODEL")
	_ = os.Unsetenv("LLAMACPP_URL")
	_ = os.Unsetenv("LLAMACPP_MODEL")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")

	// Should not fail even if .env doesn't exist
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error for non-existent .env: %v", err)
	}

	// Should still have default values
	if cfg.ModelProvider != "gemini" {
		t.Errorf("Expected default ModelProvider to be 'gemini', got '%s'", cfg.ModelProvider)
	}
}
