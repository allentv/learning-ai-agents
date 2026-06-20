# Skill: Fix Type Errors

Guide for resolving mypy strict mode type errors in this project.

## Running Type Checks

```bash
mise run tc    # Basic strict check
mise run v     # Run all checks (includes type-check)
```

## Common mypy Errors and Fixes

### `Missing return type [no-untyped-def]`

```python
# Bad
def process_data(data):
    return data.upper()

# Good
def process_data(data: str) -> str:
    return data.upper()
```

### `Missing type annotation [no-untyped-def]`

```python
# Bad
def my_function(x, y=5):
    pass

# Good
def my_function(x: str, y: int = 5) -> None:
    pass
```

### `Incompatible return type [return-value]`

Ensure return types match the declared annotation. Use `Optional[T]` for nullable returns.

### `Missing generic parameter [type-arg]`

```python
# Bad
def get_items() -> list:
    return []

# Good
def get_items() -> list[str]:
    return []
```

### `Untyped decorator [misc]`

Third-party decorators without types — suppress with `# type: ignore[misc]` and add a comment explaining why.

## Conventions

- All functions must have complete type annotations (args, return)
- Use `Optional[T]` (`T | None` in 3.10+) for nullable values
- Prefer `Pydantic BaseModel` over raw dicts for structured data
- Use `# type: ignore[specific-code]` sparingly with an explanation
- Never use `Any` unless absolutely necessary — prefer `object` or a specific type
