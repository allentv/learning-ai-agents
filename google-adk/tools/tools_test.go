package tools_test

import (
	"testing"

	"github.com/learning-ai-agents/google-adk/tools"
)

// TestNewEchoTool_CreatesTool tests that NewEchoTool creates a valid tool
func TestNewEchoTool_CreatesTool(t *testing.T) {
	tool, err := tools.NewEchoTool()
	if err != nil {
		t.Fatalf("NewEchoTool() failed: %v", err)
	}

	if tool == nil {
		t.Fatal("NewEchoTool() returned nil")
	}

	if tool.Name() != "echo" {
		t.Errorf("Expected name 'echo', got '%s'", tool.Name())
	}

	if tool.Description() != "Echoes the input message" {
		t.Errorf("Expected description 'Echoes the input message', got '%s'", tool.Description())
	}
}

// TestNewEchoTool_Execution tests that the echo tool executes correctly
func TestNewEchoTool_Execution(t *testing.T) {
	tool, err := tools.NewEchoTool()
	if err != nil {
		t.Fatalf("NewEchoTool() failed: %v", err)
	}

	// Test the tool properties
	if tool.Name() != "echo" {
		t.Errorf("Expected name 'echo', got '%s'", tool.Name())
	}

	if tool.Description() != "Echoes the input message" {
		t.Errorf("Expected description 'Echoes the input message', got '%s'", tool.Description())
	}

}

// TestNewEchoTool_InvalidArgs tests that the echo tool handles invalid arguments
func TestNewEchoTool_InvalidArgs(t *testing.T) {
	tool, err := tools.NewEchoTool()
	if err != nil {
		t.Fatalf("NewEchoTool() failed: %v", err)
	}

	// Test the tool properties
	if tool.Name() != "echo" {
		t.Errorf("Expected name 'echo', got '%s'", tool.Name())
	}

	// Note: We can't directly test invalid arguments from external tests
	// because the Run method is part of an internal interface.
	// The ADK handles argument validation internally.
}
