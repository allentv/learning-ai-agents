# Development

All development tasks are available as mise commands. Run `mise tasks` to see the full list.

## Code Formatting

```bash
mise run format          # Format code with black and isort
mise run format-check    # Check formatting without modifying files
```

## Linting

```bash
mise run lint            # Lint with ruff
mise run lint-fix        # Auto-fix linting issues
```

## Type Checking

```bash
mise run type-check                  # mypy strict mode
mise run type-check-strict           # strict mode with detailed errors
mise run type-check-verbose          # verbose output
```

## Testing

```bash
mise run test             # Run all tests
mise run test-coverage    # Run tests with coverage
```

## Run All Checks

```bash
mise run verify    # Runs format-check, lint, type-check-strict, and test
```

## Project Dependencies

- **langgraph**: Framework for building stateful, multi-agent applications
- **pydantic-ai**: Type-safe AI agent development with native tool support
- **pydantic**: Data validation and settings management
- **openai**: OpenAI API client (used with llama.cpp)
- **httpx**: Async HTTP client for API tool calls
- **python-dotenv**: Environment variable management
- **structlog**: Structured logging

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
