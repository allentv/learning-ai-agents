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
├── .mise/
│   └── tasks/                          # Mise task scripts (one file per task)
├── api/                                # FastAPI API service
│   ├── app.py                          # FastAPI application
│   ├── routes.py                       # API endpoints (GET /api/search)
│   └── requirements.txt                # API-specific dependencies
├── docs/                               # Documentation
│   ├── installation.md                 # Installation & prerequisites
│   ├── configuration.md                # Environment variables & model providers
│   ├── usage.md                        # Running with Docker or locally
│   ├── tools.md                        # Agent tool integration
│   ├── api.md                          # API service endpoints
│   ├── development.md                  # Dev workflow & contributing
│   └── docker.md                       # Docker Compose setup & GPU config
├── src/
│   └── langgraph_project/
│       ├── main.py                     # Main application entry point
│       ├── agents/
│       │   └── simple_agent.py         # Simple agent with tool support
│       ├── tools/
│       │   └── api_tools.py            # API fetch tools (pydantic-ai compatible)
│       └── utils/
│           ├── config.py               # Configuration settings
│           └── logging.py              # Logging configuration
├── tests/                              # Test files
├── pyproject.toml                      # Project dependencies
├── mise.toml                           # Mise config (tools + env only)
├── Dockerfile                          # App Dockerfile
├── Dockerfile.api                      # API service Dockerfile
├── docker-compose.yaml                 # Docker Compose (llamacpp + API + app)
├── langgraph.json                      # LangGraph Studio config
└── .env.example                        # Environment variables example
```

## Quick Start

```bash
# Install mise and UV
curl https://mise.run | sh
curl -LsSf https://astral.sh/uv/install.sh | sh

# Set up the project
mise run setup

# Configure environment
cp .env.example .env

# Run the application
mise run run
```

## Documentation

| Guide | Description |
|-------|-------------|
| [Installation](docs/installation.md) | Prerequisites, setup, and available tasks |
| [Configuration](docs/configuration.md) | Environment variables and model providers |
| [Usage](docs/usage.md) | Running with Docker or locally |
| [Tools](docs/tools.md) | Agent tool integration and adding new tools |
| [API Service](docs/api.md) | API endpoints and testing |
| [Development](docs/development.md) | Formatting, linting, testing, and contributing |
| [Docker](docs/docker.md) | Docker Compose setup, GPU config, and troubleshooting |

## License

This project is licensed under the MIT License - see the LICENSE file for details.
