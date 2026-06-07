// Package tools provides tool implementations for AI agents.
package tools

import (
	"context"
	"fmt"
)

// Tool represents a tool that can be executed
type Tool interface {
	// Name returns the name of the tool
	Name() string
	// Description returns a description of the tool
	Description() string
	// Execute runs the tool with the given parameters
	Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error)
	// GetSchema returns the tool schema for documentation
	GetSchema() map[string]interface{}
}

// BaseTool provides a base implementation for tools
type BaseTool struct {
	name        string
	description string
}

// NewBaseTool creates a new base tool
func NewBaseTool(name, description string) *BaseTool {
	return &BaseTool{
		name:        name,
		description: description,
	}
}

// Name returns the name of the tool
func (t *BaseTool) Name() string {
	return t.name
}

// Description returns the description of the tool
func (t *BaseTool) Description() string {
	return t.description
}

// GetSchema returns the tool schema for documentation
func (t *BaseTool) GetSchema() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.name,
		"description": t.description,
	}
}

// Execute is not implemented in BaseTool - should be overridden by concrete implementations
func (t *BaseTool) Execute(_ context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

// ToolRegistry manages available tools
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register registers a tool
func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

// List returns all registered tools
func (r *ToolRegistry) List() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ExecuteTool executes a tool by name with the given parameters
func (r *ToolRegistry) ExecuteTool(ctx context.Context, name string, params map[string]interface{}) (map[string]interface{}, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, &ToolNotFoundError{Name: name}
	}

	return tool.Execute(ctx, params)
}

// ToolNotFoundError represents an error when a tool is not found
type ToolNotFoundError struct {
	Name string
}

// Error implements the error interface
func (e *ToolNotFoundError) Error() string {
	return fmt.Sprintf("tool not found: %s", e.Name)
}

// ToolExecutionError represents an error during tool execution
type ToolExecutionError struct {
	ToolName string
	Err      error
}

// Error implements the error interface
func (e *ToolExecutionError) Error() string {
	return fmt.Sprintf("tool execution failed for '%s': %v", e.ToolName, e.Err)
}

// Unwrap returns the underlying error
func (e *ToolExecutionError) Unwrap() error {
	return e.Err
}
