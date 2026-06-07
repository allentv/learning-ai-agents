package main

import (
	"context"
	"log"
	"os"

	"github.com/learning-ai-agents/google-adk/config"
	"github.com/learning-ai-agents/google-adk/logging"
	"go.uber.org/zap"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
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

	// Create Gemini model using the official ADK
	if cfg.GeminiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required. Get one from https://ai.google.dev/gemini-api/docs/api-key")
	}

	model, err := gemini.NewModel(ctx, cfg.GeminiModel, &genai.ClientConfig{
		APIKey: cfg.GeminiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// Create a simple agent using the official ADK
	a, err := llmagent.New(llmagent.Config{
		Name:        "simple_agent",
		Model:       model,
		Description: "A simple agent that answers questions.",
		Instruction: "You are a helpful assistant. Provide concise and accurate responses.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
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
