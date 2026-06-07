// Package main provides the entry point for the Google ADK Go application.
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/learning-ai-agents/google-adk/config"
	"github.com/learning-ai-agents/google-adk/logging"
	localmodel "github.com/learning-ai-agents/google-adk/model"
	"go.uber.org/zap"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// Setup logging
	logging.SetupLogging()
	logger := logging.GetLogger("main")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Starting application", zap.String("name", cfg.AppName), zap.String("version", cfg.AppVersion))
	logger.Info("Using model provider", zap.String("provider", cfg.ModelProvider))

	var agentModel adkmodel.LLM
	var agentTools []tool.Tool

	switch strings.ToLower(cfg.ModelProvider) {
	case "gemini":
		// Create Gemini model using the official ADK
		if cfg.GeminiKey == "" {
			log.Fatal("GEMINI_API_KEY environment variable is required. Get one from https://ai.google.dev/gemini-api/docs/api-key")
		}

		agentModel, err = gemini.NewModel(ctx, cfg.GeminiModel, &genai.ClientConfig{
			APIKey: cfg.GeminiKey,
		})
		if err != nil {
			log.Fatalf("Failed to create Gemini model: %v", err)
		}
		agentTools = []tool.Tool{
			geminitool.GoogleSearch{},
		}
		logger.Info("Using Gemini model", zap.String("model", cfg.GeminiModel))

	case "llamacpp", "llama.cpp":
		// Create llama.cpp model using our custom implementation
		logger.Info("Using llama.cpp model", zap.String("model", cfg.LLamacppModel), zap.String("url", cfg.LLamacppURL))

		// Create a model config that points to the llama.cpp server
		modelConfig := localmodel.Config{
			Model:   cfg.LLamacppModel,
			BaseURL: cfg.LLamacppURL,
			APIKey:  "llamacpp-dummy-key", // Dummy key for llama.cpp
		}

		// Create a custom model that implements the ADK's model.LLM interface
		agentModel, err = localmodel.NewLlamaCppLLMWrapper(modelConfig)
		if err != nil {
			log.Fatalf("Failed to create llama.cpp model: %v", err)
		}
		agentTools = []tool.Tool{}

	default:
		log.Fatalf("Unknown model provider: %s. Supported providers: gemini, llamacpp", cfg.ModelProvider)
	}

	// Create a simple agent using the official ADK
	a, err := llmagent.New(llmagent.Config{
		Name:        "simple_agent",
		Model:       agentModel,
		Description: "A simple agent that answers questions.",
		Instruction: "You are a helpful assistant. Provide concise and accurate responses.",
		Tools:       agentTools,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Use the launcher to run the agent
	configLauncher := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, configLauncher, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}

	logger.Info("Google ADK for Go application completed successfully")
}
