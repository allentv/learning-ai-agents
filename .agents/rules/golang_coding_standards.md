# Golang Coding Standards for Agentic Development

## General Principles

1. **Type Safety**: Always use explicit types for function parameters, return types, and struct fields
2. **Documentation**: Use GoDoc-style comments for all exported functions, types, and packages
3. **Error Handling**: Use explicit error returns with proper error checking (never ignore errors)
4. **Logging**: Use structured logging with appropriate log levels (e.g., `log/slog`)
5. **Testing**: Write unit tests for all new functionality using `go test`

## Code Style

### Imports

- Group imports in this order: standard library, third-party, local
- Use absolute imports for project modules
- Avoid wildcard imports (Go doesn't support them anyway)
- Use `goimports` for automatic import organization

### Naming Conventions

- Packages: `lowercase` (single word or short acronym)
- Structs/Interfaces: `PascalCase`
- Functions/Methods: `PascalCase` (exported) or `camelCase` (unexported)
- Variables: `camelCase`
- Constants: `UPPER_SNAKE_CASE` or `PascalCase` (depending on context)
- Private members: `lowercase` (unexported)

### Function Design

- Keep functions focused and single-purpose
- Maximum line length: 120 characters
- Use descriptive parameter names
- Return multiple values for error handling: `(result, error)`

## Agent-Specific Patterns

### Tool Definition

```go
package tools

import (
    "context"
    "fmt"
)

// ToolInput represents the input schema for a tool.
type ToolInput struct {
    Parameter string `json:"parameter" description:"Description of parameter"`
}

// ExampleTool processes a parameter and returns a result.
// It demonstrates proper typing and documentation.
func ExampleTool(ctx context.Context, input ToolInput) (string, error) {
    // Validate input
    if input.Parameter == "" {
        return "", fmt.Errorf("parameter cannot be empty")
    }
    
    // Process the parameter
    result := fmt.Sprintf("Processed: %s", input.Parameter)
    
    return result, nil
}
```

### Agent State Management

- Use structs with exported fields for agent state
- Consider using `sync.Mutex` for concurrent state access
- Keep state immutable where possible (return new instances instead of modifying)
- Use context for cancellation and deadlines

### Error Recovery

- Implement retry logic with exponential backoff using `time.Sleep`
- Provide meaningful error messages with context
- Log errors with structured logging
- Use `errors.Wrap` or `fmt.Errorf` with `%w` for error wrapping

## Quality Assurance

### Package Management

- **Always use `go mod` for dependency management**
- Never manually edit `go.mod` or `go.sum`
- Use `go mod tidy` to clean up dependencies
- Use `go get` to add new dependencies
- Use `go mod download` to download dependencies

### Pre-commit Hooks

- Run `go fmt` for formatting
- Run `go vet` for static analysis
- Run `golangci-lint` for comprehensive linting
- Run `go test` for testing
- Run `go mod tidy` to verify dependency management

### Code Review Checklist

- [ ] Proper error handling (no ignored errors)
- [ ] GoDoc comments present and accurate
- [ ] Context used for cancellation and deadlines
- [ ] Tests added or updated
- [ ] Logging appropriate
- [ ] No hardcoded values
- [ ] Thread-safe concurrent access (if applicable)

## Additional Go-Specific Guidelines

### Concurrency

- Use goroutines for concurrent operations
- Use channels for communication between goroutines
- Use `select` for handling multiple channels
- Always close channels when appropriate
- Use `sync.WaitGroup` for coordinating goroutines

### Memory Management

- Be mindful of memory allocations in hot paths
- Use `sync.Pool` for frequently allocated objects
- Avoid unnecessary pointer usage
- Let Go's garbage collector handle cleanup

### Performance

- Profile code with `go test -bench` and `pprof`
- Optimize based on actual performance measurements
- Consider algorithmic complexity before micro-optimizations
