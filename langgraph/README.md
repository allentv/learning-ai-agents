# LangGraph Project with Pydantic AI

A Python project built with LangGraph and Pydantic AI for creating intelligent agent workflows.

## Features

- **LangGraph**: Build stateful, multi-agent applications with LangGraph
- **Pydantic AI**: Type-safe AI agent development with Pydantic models
- **Structured Logging**: Comprehensive logging with structlog
- **Configuration Management**: Environment-based configuration with Pydantic Settings
- **Async Support**: Full async/await support for efficient execution
- **Testing**: Comprehensive test suite with pytest

## Project Structure

```text
langgraph/
├── src/
│   └── langgraph_project/
│       ├── __init__.py
│       ├── main.py                 # Main application entry point
│       ├── agents/
│       │   ├── __init__.py
│       │   └── simple_agent.py     # Simple agent implementation
│       ├── tools/
│       │   ├── __init__.py
│       │   └── base_tool.py        # Base tool class
│       └── utils/
│           ├── __init__.py
│           ├── config.py           # Configuration settings
│           └── logging.py          # Logging configuration
├── tests/                          # Test files
├── pyproject.toml                  # Poetry dependencies
├── README.md                       # This file
└── .env.example                    # Environment variables example
```

## Installation

1. **Install UV** (if not already installed):

   ```bash
   curl -LsSf https://astral.sh/uv/install.sh | sh
   ```

2. **Install dependencies**:

   ```bash
   uv sync
   ```

3. **Activate the virtual environment**:

   ```bash
   source .venv/bin/activate
   ```

   **Note**: If you see a warning about `VIRTUAL_ENV` not matching the project environment path, you can either:
   - Deactivate your current virtual environment first: `deactivate`
   - Or use `uv run` directly without activating the environment

4. **Set up environment variables**:

   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Model Provider Configuration
MODEL_PROVIDER=llamacpp  # "openai" or "llamacpp"

# OpenAI API Configuration
OPENAI_API_KEY=your-openai-api-key-here
OPENAI_MODEL=gpt-4o-mini

# llama.cpp Configuration
LLAMACPP_URL=http://localhost:12434/v1
LLAMACPP_MODEL=granite-4.0-h-micro-UD-Q4_K_XL.gguf

# Logging Configuration
LOG_LEVEL=INFO
LOG_FORMAT=json  # or "text"

# Application Configuration
APP_NAME=LangGraph Project
APP_VERSION=0.1.0
```

## Usage

### Running with llama.cpp (Recommended)

1. **Download the model**:

   ```bash
   mise run download-model
   ```

2. **Start services with Docker Compose**:

   ```bash
   docker-compose up -d
   ```

3. **Run the application**:

   ```bash
   docker-compose exec app uv run python -m langgraph_project.main
   ```

### Running the main application locally

```bash
python -m langgraph_project.main
```

### Using the agent directly

```python
import asyncio
from langgraph_project.agents.simple_agent import SimpleAgent

async def main():
    agent = SimpleAgent()
    response = await agent.process("Hello, how are you?")
    print(response)

asyncio.run(main())
```

## Development

### Code Formatting

```bash
# Format code with black
uv run black src/ tests/

# Sort imports with isort
uv run isort src/ tests/
```

### Type Checking

```bash
# Run mypy for type checking
uv run mypy src/
```

### Linting

```bash
# Run ruff for linting
uv run ruff check src/ tests/
```

### Testing

```bash
# Run all tests
uv run pytest

# Run tests with coverage
uv run pytest --cov=src/langgraph_project

# Run specific test file
uv run pytest tests/test_specific.py
```

## Project Dependencies

- **langgraph**: Framework for building stateful, multi-agent applications
- **pydantic-ai**: Type-safe AI agent development
- **pydantic**: Data validation and settings management
- **openai**: OpenAI API client (used with llama.cpp)
- **python-dotenv**: Environment variable management
- **structlog**: Structured logging

### Model Providers

The project supports multiple model providers:

- **llamacpp**: Local llama.cpp server (recommended for local development)
- **openai**: OpenAI API
- **vllm**: vLLM server (for GPU-accelerated inference)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- LangGraph team for the excellent framework
- Pydantic team for the type-safe AI development tools
- OpenAI for the API and models
