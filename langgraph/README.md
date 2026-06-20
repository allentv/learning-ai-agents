# LangGraph Project with Pydantic AI

A Python project built with LangGraph and Pydantic AI for creating intelligent agent workflows.

## Features

- **LangGraph**: Build stateful, multi-agent applications with LangGraph
- **Pydantic AI**: Type-safe AI agent development with Pydantic models
- **Tool Calls**: Agent can fetch information from external APIs using registered tools
- **API Service**: Lightweight FastAPI service providing endpoints for knowledge base search (RAG-ready)
- **Structured Logging**: Comprehensive logging with structlog
- **Configuration Management**: Environment-based configuration with Pydantic Settings
- **Async Support**: Full async/await support for efficient execution
- **Testing**: Comprehensive test suite with pytest

## Project Structure

```text
langgraph/
├── api/                                # FastAPI API service
│   ├── __init__.py
│   ├── app.py                          # FastAPI application
│   ├── routes.py                       # API endpoints (GET /api/search)
│   └── requirements.txt                # API-specific dependencies
├── src/
│   └── langgraph_project/
│       ├── __init__.py
│       ├── main.py                     # Main application entry point
│       ├── agents/
│       │   ├── __init__.py
│       │   └── simple_agent.py         # Simple agent with tool support
│       ├── tools/
│       │   ├── __init__.py             # Tool exports
│       │   └── api_tools.py            # API fetch tools (pydantic-ai compatible)
│       └── utils/
│           ├── __init__.py
│           ├── config.py               # Configuration settings
│           └── logging.py              # Logging configuration
├── tests/                              # Test files
├── pyproject.toml                      # Project dependencies
├── Dockerfile                          # App Dockerfile
├── Dockerfile.api                      # API service Dockerfile
├── docker-compose.yaml                 # Docker Compose (llamacpp + API + app)
├── langgraph.json                      # LangGraph Studio config
├── README.md                           # This file
└── .env.example                        # Environment variables example
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

# External API Configuration (used by agent tools)
API_BASE_URL=http://api:10000
API_TIMEOUT=30

# Logging Configuration
LOG_LEVEL=INFO
LOG_FORMAT=json  # or "text"

# Application Configuration
APP_NAME=LangGraph Project
APP_VERSION=0.1.0
```

## Usage

### Running with Docker Compose (Recommended)

1. **Download the model**:

   ```bash
   mise run download-model
   ```

2. **Start all services** (llamacpp + API + app):

   ```bash
   docker-compose up -d
   ```

   This starts three services:
   - **llamacpp**: Local LLM server on port `12434`
   - **api**: FastAPI knowledge base API on port `10000`
   - **app**: LangGraph agent application on port `8080`

3. **Run the application**:

   ```bash
   docker-compose exec app uv run python -m langgraph_project.main
   ```

### Running locally

```bash
python -m langgraph_project.main
```

### Using the agent directly with tools

```python
import asyncio
from langgraph_project.agents.simple_agent import SimpleAgent
from langgraph_project.tools import search_api

async def main():
    # Initialize agent with tools
    agent = SimpleAgent(tools=[search_api])

    # The agent will automatically call search_api when it needs information
    response = await agent.process("What is the latest news about AI?")
    print(response)

asyncio.run(main())
```

### Running the API service locally

```bash
uvicorn api.app:app --host 0.0.0.0 --port 10000 --reload
```

Test the API:
```bash
curl "http://localhost:10000/api/search?q=hello"
```

## Tools

The agent uses [pydantic-ai](https://ai.pydantic.dev/) for tool integration. Tools are plain async Python functions with type-annotated parameters — pydantic-ai automatically generates the JSON schema for the LLM.

### Available Tools

| Tool | Description | Endpoint |
|------|-------------|----------|
| `search_api` | Search the knowledge base via the external API | `GET /api/search?q=<query>` |

### Adding a New Tool

1. **Create the tool function** in `src/langgraph_project/tools/api_tools.py` (or a new file):

   ```python
   async def my_new_tool(param: str) -> dict[str, Any]:
       """Description of what the tool does.

       This docstring is used as the tool description for the LLM.

       Args:
           param: Description of the parameter.
       """
       # Implementation here
       return {"result": "..."}
   ```

2. **Export it** in `src/langgraph_project/tools/__init__.py`:

   ```python
   from langgraph_project.tools.api_tools import search_api, my_new_tool

   __all__ = ["search_api", "my_new_tool"]
   ```

3. **Register it** with the agent in `src/langgraph_project/main.py`:

   ```python
   from langgraph_project.tools import search_api, my_new_tool

   AGENT_TOOLS = [search_api, my_new_tool]
   ```

### How Tool Calls Work

1. The user sends a query to the agent
2. The LLM decides whether to call a tool based on the tool schemas
3. If a tool call is needed, pydantic-ai executes the tool function
4. The tool result is fed back to the LLM
5. The LLM may call more tools or produce a final text response
6. This loop continues until the LLM produces a final answer

This entire loop is handled automatically by pydantic-ai within the `process_query` node of the LangGraph graph.

## API Service

The project includes a lightweight FastAPI service designed as a placeholder for future RAG (Retrieval-Augmented Generation) functionality.

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/search?q=<query>` | Search the knowledge base |
| `GET` | `/health` | Health check |

### Example Response

```json
{
  "query": "hello",
  "results": [
    {
      "id": 1,
      "title": "Result for 'hello'",
      "snippet": "This is a placeholder result...",
      "score": 0.95
    }
  ],
  "metadata": {
    "total_results": 1,
    "source": "mock"
  }
}
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
uv run pytest tests/test_api_tools.py
```

## Project Dependencies

- **langgraph**: Framework for building stateful, multi-agent applications
- **pydantic-ai**: Type-safe AI agent development with native tool support
- **pydantic**: Data validation and settings management
- **openai**: OpenAI API client (used with llama.cpp)
- **httpx**: Async HTTP client for API tool calls
- **python-dotenv**: Environment variable management
- **structlog**: Structured logging

### Model Providers

The project supports multiple model providers:

- **llamacpp**: Local llama.cpp server (recommended for local development)
- **openai**: OpenAI API

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
