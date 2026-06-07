package tools_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/learning-ai-agents/google-adk/tools"
)

// mockTool is a mock implementation of the Tool interface for testing
type mockTool struct {
	name        string
	description string
	executeFunc func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error)
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Description() string {
	return m.description
}

func (m *mockTool) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, params)
	}
	return map[string]interface{}{"result": "success"}, nil
}

func (m *mockTool) GetSchema() map[string]interface{} {
	return map[string]interface{}{
		"name":        m.name,
		"description": m.description,
	}
}

func TestBaseTool_NewBaseTool(t *testing.T) {
	tool := tools.NewBaseTool("test-tool", "A test tool")

	if tool == nil {
		t.Fatal("NewBaseTool() returned nil")
	}

	if tool.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got '%s'", tool.Name())
	}

	if tool.Description() != "A test tool" {
		t.Errorf("Expected description 'A test tool', got '%s'", tool.Description())
	}
}

func TestBaseTool_Name(t *testing.T) {
	tool := tools.NewBaseTool("my-tool", "My tool description")

	if tool.Name() != "my-tool" {
		t.Errorf("Expected name 'my-tool', got '%s'", tool.Name())
	}
}

func TestBaseTool_Description(t *testing.T) {
	tool := tools.NewBaseTool("my-tool", "My tool description")

	if tool.Description() != "My tool description" {
		t.Errorf("Expected description 'My tool description', got '%s'", tool.Description())
	}
}

func TestBaseTool_GetSchema(t *testing.T) {
	tool := tools.NewBaseTool("my-tool", "My tool description")

	schema := tool.GetSchema()

	if schema["name"] != "my-tool" {
		t.Errorf("Expected schema name 'my-tool', got '%v'", schema["name"])
	}

	if schema["description"] != "My tool description" {
		t.Errorf("Expected schema description 'My tool description', got '%v'", schema["description"])
	}
}

func TestBaseTool_Execute(t *testing.T) {
	tool := tools.NewBaseTool("my-tool", "My tool description")

	ctx := context.Background()
	params := map[string]interface{}{"key": "value"}

	result, err := tool.Execute(ctx, params)

	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result from BaseTool.Execute(), got %v", result)
	}
}

func TestExampleTool_NewExampleTool(t *testing.T) {
	tool := tools.NewExampleTool()

	if tool == nil {
		t.Fatal("NewExampleTool() returned nil")
	}

	if tool.Name() != "echo" {
		t.Errorf("Expected name 'echo', got '%s'", tool.Name())
	}

	if tool.Description() != "Echoes the input message" {
		t.Errorf("Expected description 'Echoes the input message', got '%s'", tool.Description())
	}
}

func TestExampleTool_Execute_Success(t *testing.T) {
	tool := tools.NewExampleTool()

	ctx := context.Background()
	params := map[string]interface{}{
		"message": "Hello, world!",
	}

	result, err := tool.Execute(ctx, params)

	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Execute() returned nil result")
	}

	echoedMessage, ok := result["echoed_message"].(string)
	if !ok {
		t.Fatal("echoed_message not found or not a string")
	}

	if echoedMessage != "Hello, world!" {
		t.Errorf("Expected echoed_message 'Hello, world!', got '%s'", echoedMessage)
	}

	// Check timestamp exists
	if _, ok := result["timestamp"]; !ok {
		t.Error("timestamp not found in result")
	}
}

func TestExampleTool_Execute_MissingMessage(t *testing.T) {
	tool := tools.NewExampleTool()

	ctx := context.Background()
	params := map[string]interface{}{}

	_, err := tool.Execute(ctx, params)

	if err == nil {
		t.Fatal("Execute() should return error for missing message parameter")
	}

	// Check error type
	execErr, ok := err.(*tools.ToolExecutionError)
	if !ok {
		t.Fatalf("Expected ToolExecutionError, got %T", err)
	}

	if execErr.ToolName != "echo" {
		t.Errorf("Expected ToolName 'echo', got '%s'", execErr.ToolName)
	}

	if execErr.Err == nil {
		t.Fatal("Expected non-nil Err in ToolExecutionError")
	}
}

func TestExampleTool_Execute_InvalidMessageType(t *testing.T) {
	tool := tools.NewExampleTool()

	ctx := context.Background()
	params := map[string]interface{}{
		"message": 123, // Wrong type
	}

	_, err := tool.Execute(ctx, params)

	if err == nil {
		t.Fatal("Execute() should return error for invalid message type")
	}

	// Check error type
	execErr, ok := err.(*tools.ToolExecutionError)
	if !ok {
		t.Fatalf("Expected ToolExecutionError, got %T", err)
	}

	if execErr.ToolName != "echo" {
		t.Errorf("Expected ToolName 'echo', got '%s'", execErr.ToolName)
	}
}

func TestExampleTool_GetSchema(t *testing.T) {
	tool := tools.NewExampleTool()

	schema := tool.GetSchema()

	if schema["name"] != "echo" {
		t.Errorf("Expected schema name 'echo', got '%v'", schema["name"])
	}

	if schema["description"] != "Echoes the input message" {
		t.Errorf("Expected schema description 'Echoes the input message', got '%v'", schema["description"])
	}

	// Check parameters schema
	params, ok := schema["parameters"].(map[string]interface{})
	if !ok {
		t.Fatal("parameters not found or not a map")
	}

	if params["type"] != "object" {
		t.Errorf("Expected parameters type 'object', got '%v'", params["type"])
	}

	// Check required parameters - can be []string or []interface{}
	switch v := params["required"].(type) {
	case []string:
		if len(v) != 1 {
			t.Errorf("Expected 1 required parameter, got %d", len(v))
		}
		if v[0] != "message" {
			t.Errorf("Expected required parameter 'message', got '%s'", v[0])
		}
	case []interface{}:
		if len(v) != 1 {
			t.Errorf("Expected 1 required parameter, got %d", len(v))
		}
		if v[0] != "message" {
			t.Errorf("Expected required parameter 'message', got '%v'", v[0])
		}
	default:
		t.Fatalf("required has unexpected type: %T", params["required"])
	}
}

func TestToolRegistry_NewToolRegistry(t *testing.T) {
	registry := tools.NewToolRegistry()

	if registry == nil {
		t.Fatal("NewToolRegistry() returned nil")
	}
}

func TestToolRegistry_Register(t *testing.T) {
	registry := tools.NewToolRegistry()

	tool := &mockTool{
		name:        "test-tool",
		description: "A test tool",
	}

	registry.Register(tool)

	// Verify tool was registered
	retrievedTool, exists := registry.Get("test-tool")
	if !exists {
		t.Fatal("Tool was not registered")
	}

	if retrievedTool.Name() != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got '%s'", retrievedTool.Name())
	}
}

func TestToolRegistry_Get(t *testing.T) {
	registry := tools.NewToolRegistry()

	tool := &mockTool{
		name:        "test-tool",
		description: "A test tool",
	}

	registry.Register(tool)

	// Test getting existing tool
	retrievedTool, exists := registry.Get("test-tool")
	if !exists {
		t.Fatal("Tool should exist")
	}

	if retrievedTool != tool {
		t.Error("Retrieved tool is not the same as registered tool")
	}

	// Test getting non-existent tool
	_, exists = registry.Get("non-existent")
	if exists {
		t.Error("Non-existent tool should not exist")
	}
}

func TestToolRegistry_List(t *testing.T) {
	registry := tools.NewToolRegistry()

	tool1 := &mockTool{
		name:        "tool1",
		description: "Tool 1",
	}
	tool2 := &mockTool{
		name:        "tool2",
		description: "Tool 2",
	}

	registry.Register(tool1)
	registry.Register(tool2)

	tools := registry.List()

	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}

	// Check that both tools are present
	found1 := false
	found2 := false
	for _, t := range tools {
		if t.Name() == "tool1" {
			found1 = true
		}
		if t.Name() == "tool2" {
			found2 = true
		}
	}

	if !found1 {
		t.Error("tool1 not found in list")
	}

	if !found2 {
		t.Error("tool2 not found in list")
	}
}

func TestToolRegistry_ExecuteTool_Success(t *testing.T) {
	registry := tools.NewToolRegistry()

	tool := &mockTool{
		name:        "test-tool",
		description: "A test tool",
		executeFunc: func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{
				"result":  "success",
				"message": params["message"],
			}, nil
		},
	}

	registry.Register(tool)

	ctx := context.Background()
	params := map[string]interface{}{
		"message": "test message",
	}

	result, err := registry.ExecuteTool(ctx, "test-tool", params)

	if err != nil {
		t.Fatalf("ExecuteTool() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("ExecuteTool() returned nil result")
	}

	resultStr, ok := result["result"].(string)
	if !ok || resultStr != "success" {
		t.Errorf("Expected result 'success', got '%v'", result["result"])
	}
}

func TestToolRegistry_ExecuteTool_NotFound(t *testing.T) {
	registry := tools.NewToolRegistry()

	ctx := context.Background()
	params := map[string]interface{}{}

	_, err := registry.ExecuteTool(ctx, "non-existent", params)

	if err == nil {
		t.Fatal("ExecuteTool() should return error for non-existent tool")
	}

	// Check error type
	notFoundErr, ok := err.(*tools.ToolNotFoundError)
	if !ok {
		t.Fatalf("Expected ToolNotFoundError, got %T", err)
	}

	if notFoundErr.Name != "non-existent" {
		t.Errorf("Expected Name 'non-existent', got '%s'", notFoundErr.Name)
	}
}

func TestToolRegistry_ExecuteTool_WithParams(t *testing.T) {
	registry := tools.NewToolRegistry()

	tool := &mockTool{
		name:        "test-tool",
		description: "A test tool",
		executeFunc: func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			// Verify params are passed correctly
			if params["key1"] != "value1" {
				return nil, fmt.Errorf("expected key1=value1, got %v", params["key1"])
			}
			if params["key2"] != 42 {
				return nil, fmt.Errorf("expected key2=42, got %v", params["key2"])
			}
			return map[string]interface{}{"status": "ok"}, nil
		},
	}

	registry.Register(tool)

	ctx := context.Background()
	params := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	result, err := registry.ExecuteTool(ctx, "test-tool", params)

	if err != nil {
		t.Fatalf("ExecuteTool() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("ExecuteTool() returned nil result")
	}

	status, ok := result["status"].(string)
	if !ok || status != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", result["status"])
	}
}

func TestToolNotFoundError_Error(t *testing.T) {
	err := &tools.ToolNotFoundError{
		Name: "missing-tool",
	}

	expected := "tool not found: missing-tool"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestToolExecutionError_Error(t *testing.T) {
	innerErr := fmt.Errorf("inner error")
	err := &tools.ToolExecutionError{
		ToolName: "failing-tool",
		Err:      innerErr,
	}

	expected := "tool execution failed for 'failing-tool': inner error"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestToolExecutionError_Unwrap(t *testing.T) {
	innerErr := fmt.Errorf("inner error")
	err := &tools.ToolExecutionError{
		ToolName: "failing-tool",
		Err:      innerErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != innerErr {
		t.Error("Unwrap() did not return the inner error")
	}
}

func TestToolExecutionError_NilUnwrap(t *testing.T) {
	err := &tools.ToolExecutionError{
		ToolName: "failing-tool",
		Err:      nil,
	}

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Errorf("Expected nil from Unwrap(), got %v", unwrapped)
	}
}
