package workflow_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/learning-ai-agents/google-adk/workflow"
)

// TestState represents a simple state for testing
type TestState struct {
	Value    int
	Messages []string
}

func TestWorkflow_NewWorkflow(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	if w == nil {
		t.Fatal("NewWorkflow() returned nil")
	}
}

func TestWorkflow_AddNode(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}

	result := w.AddNode("node1", node)

	if result != w {
		t.Error("AddNode() should return the workflow for chaining")
	}
}

func TestWorkflow_AddEdge(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node1 := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}
	node2 := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}

	w.AddNode("node1", node1)
	w.AddNode("node2", node2)

	result := w.AddEdge("node1", "node2")

	if result != w {
		t.Error("AddEdge() should return the workflow for chaining")
	}
}

func TestWorkflow_SetEntryPoint(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}

	w.AddNode("start", node)

	result := w.SetEntryPoint("start")

	if result != w {
		t.Error("SetEntryPoint() should return the workflow for chaining")
	}
}

func TestWorkflow_Compile(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	if app == nil {
		t.Fatal("Compile() returned nil")
	}
}

func TestApp_Invoke_Success(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 1
		state.Messages = append(state.Messages, "processed")
		return state, nil
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	result, err := app.Invoke(context.Background(), initialState)

	if err != nil {
		t.Fatalf("Invoke() returned error: %v", err)
	}

	if result.Value != 1 {
		t.Errorf("Expected value 1, got %d", result.Value)
	}

	if len(result.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(result.Messages))
	}

	if result.Messages[0] != "processed" {
		t.Errorf("Expected message 'processed', got '%s'", result.Messages[0])
	}
}

func TestApp_Invoke_NoEntryPoint(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}

	w.AddNode("start", node)
	// Don't set entry point

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	_, err := app.Invoke(context.Background(), initialState)

	if err == nil {
		t.Fatal("Invoke() should return error when no entry point is set")
	}

	expectedError := "no entry point set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestApp_Invoke_EntryPointNotFound(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, nil
	}

	w.AddNode("start", node)
	w.SetEntryPoint("non-existent")

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	_, err := app.Invoke(context.Background(), initialState)

	if err == nil {
		t.Fatal("Invoke() should return error when entry point node is not found")
	}

	expectedError := "entry point node 'non-existent' not found"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestApp_Invoke_NodeError(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, fmt.Errorf("node execution failed")
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	_, err := app.Invoke(context.Background(), initialState)

	if err == nil {
		t.Fatal("Invoke() should return error when node execution fails")
	}

	if err.Error() == "" {
		t.Error("Invoke() returned empty error message")
	}

	// Check that error contains the node name
	if !contains(err.Error(), "start") {
		t.Errorf("Error should contain node name 'start', got '%s'", err.Error())
	}
}

func TestApp_AInvoke_Success(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 10
		return state, nil
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	initialState := TestState{
		Value:    5,
		Messages: []string{},
	}

	result, err := app.AInvoke(context.Background(), initialState)

	if err != nil {
		t.Fatalf("AInvoke() returned error: %v", err)
	}

	if result.Value != 15 {
		t.Errorf("Expected value 15, got %d", result.Value)
	}
}

func TestWorkflow_EdgeStruct(t *testing.T) {
	edge := workflow.Edge{
		From: "node1",
		To:   "node2",
	}

	if edge.From != "node1" {
		t.Errorf("Expected From 'node1', got '%s'", edge.From)
	}

	if edge.To != "node2" {
		t.Errorf("Expected To 'node2', got '%s'", edge.To)
	}
}

func TestWorkflow_MultipleNodesAndEdges(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node1 := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 1
		return state, nil
	}

	node2 := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value * 2
		return state, nil
	}

	node3 := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 10
		return state, nil
	}

	w.AddNode("node1", node1)
	w.AddNode("node2", node2)
	w.AddNode("node3", node3)
	w.AddEdge("node1", "node2")
	w.AddEdge("node2", "node3")
	w.SetEntryPoint("node1")

	app := w.Compile()

	initialState := TestState{
		Value:    5,
		Messages: []string{},
	}

	// Note: Current implementation only executes the entry point node
	// It doesn't follow edges to other nodes
	result, err := app.Invoke(context.Background(), initialState)

	if err != nil {
		t.Fatalf("Invoke() returned error: %v", err)
	}

	// Only node1 should execute
	if result.Value != 6 {
		t.Errorf("Expected value 6 (5 + 1), got %d", result.Value)
	}
}

func TestWorkflow_GenericStateTypes(t *testing.T) {
	// Test with string state
	wString := workflow.NewWorkflow[string]()

	node := func(ctx context.Context, state string) (string, error) {
		return state + " processed", nil
	}

	wString.AddNode("start", node)
	wString.SetEntryPoint("start")

	app := wString.Compile()

	result, err := app.Invoke(context.Background(), "input")

	if err != nil {
		t.Fatalf("Invoke() with string state returned error: %v", err)
	}

	expected := "input processed"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with int state
	wInt := workflow.NewWorkflow[int]()

	nodeInt := func(ctx context.Context, state int) (int, error) {
		return state * 2, nil
	}

	wInt.AddNode("start", nodeInt)
	wInt.SetEntryPoint("start")

	appInt := wInt.Compile()

	resultInt, err := appInt.Invoke(context.Background(), 5)

	if err != nil {
		t.Fatalf("Invoke() with int state returned error: %v", err)
	}

	if resultInt != 10 {
		t.Errorf("Expected 10, got %d", resultInt)
	}
}

func TestWorkflow_ChainedOperations(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node1 := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 1
		state.Messages = append(state.Messages, "step1")
		return state, nil
	}

	node2 := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 2
		state.Messages = append(state.Messages, "step2")
		return state, nil
	}

	// Chain operations in a single node
	node3 := func(ctx context.Context, state TestState) (TestState, error) {
		state.Value = state.Value + 3
		state.Messages = append(state.Messages, "step3")
		return state, nil
	}

	w.AddNode("node1", node1)
	w.AddNode("node2", node2)
	w.AddNode("node3", node3)
	w.SetEntryPoint("node1")

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	result, err := app.Invoke(context.Background(), initialState)

	if err != nil {
		t.Fatalf("Invoke() returned error: %v", err)
	}

	// Only node1 should execute
	if result.Value != 1 {
		t.Errorf("Expected value 1, got %d", result.Value)
	}

	if len(result.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(result.Messages))
	}
}

func TestWorkflow_ContextCancellation(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return state, ctx.Err()
		default:
			return state, nil
		}
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	_, err := app.Invoke(ctx, initialState)

	if err == nil {
		t.Fatal("Invoke() should return error when context is cancelled")
	}
}

func TestWorkflow_EmptyWorkflow(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	_, err := app.Invoke(context.Background(), initialState)

	if err == nil {
		t.Fatal("Invoke() should return error for empty workflow (no entry point)")
	}
}

func TestWorkflow_NodeReturningError(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		return state, fmt.Errorf("custom error: invalid state")
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	initialState := TestState{
		Value:    0,
		Messages: []string{},
	}

	_, err := app.Invoke(context.Background(), initialState)

	if err == nil {
		t.Fatal("Invoke() should return error when node returns error")
	}

	if !contains(err.Error(), "custom error") {
		t.Errorf("Error should contain 'custom error', got '%s'", err.Error())
	}
}

func TestWorkflow_StateModification(t *testing.T) {
	w := workflow.NewWorkflow[TestState]()

	node := func(ctx context.Context, state TestState) (TestState, error) {
		// Modify state in various ways
		state.Value = state.Value * 2
		state.Messages = append(state.Messages, "modified")
		return state, nil
	}

	w.AddNode("start", node)
	w.SetEntryPoint("start")

	app := w.Compile()

	initialState := TestState{
		Value:    10,
		Messages: []string{"original"},
	}

	result, err := app.Invoke(context.Background(), initialState)

	if err != nil {
		t.Fatalf("Invoke() returned error: %v", err)
	}

	if result.Value != 20 {
		t.Errorf("Expected value 20, got %d", result.Value)
	}

	if len(result.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(result.Messages))
	}

	if result.Messages[0] != "original" {
		t.Errorf("Expected first message 'original', got '%s'", result.Messages[0])
	}

	if result.Messages[1] != "modified" {
		t.Errorf("Expected second message 'modified', got '%s'", result.Messages[1])
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
