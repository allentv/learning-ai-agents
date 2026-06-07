// Package tools provides tool implementations using the official Google ADK.
package tools

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// EchoArgs represents the arguments for the echo tool
type EchoArgs struct {
	Message string `json:"message"`
}

// EchoResult represents the result of the echo tool
type EchoResult struct {
	EchoedMessage string `json:"echoed_message"`
	Timestamp     string `json:"timestamp"`
}

// NewEchoTool creates a new echo tool using the official ADK
func NewEchoTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "echo",
		Description: "Echoes the input message",
	}, func(_ agent.ToolContext, args EchoArgs) (EchoResult, error) {
		return EchoResult{
			EchoedMessage: args.Message,
			Timestamp:     "2026-06-07T00:00:00Z", // Placeholder timestamp
		}, nil
	})
}
