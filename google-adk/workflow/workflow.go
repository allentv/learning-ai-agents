// Package workflow provides workflow graph functionality for orchestrating AI agent tasks.
package workflow

import (
	"context"
	"fmt"
	"sync"
)

// Node represents a workflow node
type Node[T any] func(ctx context.Context, state T) (T, error)

// Edge represents a workflow edge
type Edge struct {
	From string
	To   string
}

// Workflow represents a workflow graph
type Workflow[T any] struct {
	nodes      map[string]Node[T]
	edges      []Edge
	entryPoint string
	mu         sync.RWMutex
}

// NewWorkflow creates a new workflow
func NewWorkflow[T any]() *Workflow[T] {
	return &Workflow[T]{
		nodes: make(map[string]Node[T]),
		edges: make([]Edge, 0),
	}
}

// AddNode adds a node to the workflow
func (w *Workflow[T]) AddNode(name string, node Node[T]) *Workflow[T] {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.nodes[name] = node
	return w
}

// AddEdge adds an edge between nodes
func (w *Workflow[T]) AddEdge(from, to string) *Workflow[T] {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.edges = append(w.edges, Edge{From: from, To: to})
	return w
}

// SetEntryPoint sets the entry point for the workflow
func (w *Workflow[T]) SetEntryPoint(name string) *Workflow[T] {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.entryPoint = name
	return w
}

// Compile compiles the workflow into an executable app
func (w *Workflow[T]) Compile() *App[T] {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return &App[T]{
		workflow: w,
	}
}

// App represents a compiled workflow
type App[T any] struct {
	workflow *Workflow[T]
}

// Invoke runs the workflow synchronously
func (a *App[T]) Invoke(ctx context.Context, initialState T) (T, error) {
	// Execute entry point node
	if a.workflow.entryPoint == "" {
		var zero T
		return zero, fmt.Errorf("no entry point set")
	}

	node, exists := a.workflow.nodes[a.workflow.entryPoint]
	if !exists {
		var zero T
		return zero, fmt.Errorf("entry point node '%s' not found", a.workflow.entryPoint)
	}

	result, err := node(ctx, initialState)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to execute node '%s': %w", a.workflow.entryPoint, err)
	}

	return result, nil
}

// AInvoke runs the workflow asynchronously
func (a *App[T]) AInvoke(ctx context.Context, initialState T) (T, error) {
	// For now, just call Invoke since Go doesn't have native async/await like Python
	return a.Invoke(ctx, initialState)
}
