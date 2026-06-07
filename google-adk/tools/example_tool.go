package tools

import (
	"context"
	"fmt"
)

// ExampleTool is a simple example tool that echoes input
type ExampleTool struct {
	*BaseTool
}

// NewExampleTool creates a new example tool
func NewExampleTool() *ExampleTool {
	return &ExampleTool{
		BaseTool: NewBaseTool("echo", "Echoes the input message"),
	}
}

// Execute implements the Tool interface
func (t *ExampleTool) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	message, ok := params["message"].(string)
	if !ok {
		return nil, &ToolExecutionError{
			ToolName: t.Name(),
			Err:      fmt.Errorf("missing or invalid 'message' parameter"),
		}
	}

	result := map[string]interface{}{
		"echoed_message": message,
		"timestamp":      "2026-06-07T00:00:00Z", // Placeholder timestamp
	}

	return result, nil
}

// GetSchema returns the tool schema with parameters
func (t *ExampleTool) GetSchema() map[string]interface{} {
	schema := t.BaseTool.GetSchema()
	schema["parameters"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "The message to echo",
			},
		},
		"required": []string{"message"},
	}
	return schema
}
