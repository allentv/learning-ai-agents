# Python Coding Standards for Agentic Development

## General Principles

1. **Type Hints**: Always use type hints for function parameters, return types, and class attributes
2. **Docstrings**: Use Google-style docstrings for all public functions and classes
3. **Error Handling**: Use explicit exception handling with specific exception types
4. **Logging**: Use structured logging with appropriate log levels
5. **Testing**: Write unit tests for all new functionality

## Code Style

### Imports

- Group imports in this order: standard library, third-party, local
- Use absolute imports for project modules
- Avoid wildcard imports

### Naming Conventions

- Classes: `PascalCase`
- Functions/Methods: `snake_case`
- Variables: `snake_case`
- Constants: `UPPER_SNAKE_CASE`
- Private members: `_private_method`

### Function Design

- Keep functions focused and single-purpose
- Maximum line length: 120 characters
- Use descriptive parameter names
- Avoid mutable default arguments

## Agent-Specific Patterns

### Tool Definition

```python
from langchain_core.tools import tool
from pydantic import BaseModel, Field

class ToolInput(BaseModel):
    """Input schema for the tool."""
    parameter: str = Field(description="Description of parameter")

@tool(args_schema=ToolInput)
def example_tool(parameter: str) -> str:
    """Example tool with proper typing and documentation.
    
    Args:
        parameter: Description of parameter
        
    Returns:
        Description of return value
    """
    return f"Processed: {parameter}"
```

### Agent State Management

- Use TypedDict or Pydantic models for agent state
- Validate state transitions
- Keep state immutable where possible

### Error Recovery

- Implement retry logic with exponential backoff
- Provide meaningful error messages
- Log errors with context

## Quality Assurance

### Package Management

- **Always use `uv` for package installations**
- Never use `pip` directly for package installations
- Use `uv sync` to sync dependencies from `pyproject.toml`
- Use `uv lock` to manage the lock file
- Virtual environments should be created with `uv venv`

### Pre-commit Hooks

- Run `ruff check` for linting
- Run `black --check` for formatting
- Run `mypy` for type checking
- Run `pytest` for testing
- Run `uv sync --dry-run` to verify dependency sync

### Code Review Checklist

- [ ] Type hints present and correct
- [ ] Docstrings complete and accurate
- [ ] Error handling implemented
- [ ] Tests added or updated
- [ ] Logging appropriate
- [ ] No hardcoded values
