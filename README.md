# Learning AI Agents

A comprehensive collection of AI agent implementations and experiments using various frameworks and technologies. This project serves as a learning resource for building intelligent agent systems.

## Project Structure

This repository contains multiple sub-projects, each exploring different AI agent frameworks and approaches:

### 1. AWS-Strands (`/AWS-Strands/`)

An AI agent service built with FastAPI and the Strands Agents SDK, featuring:

- **FastAPI** web framework for the agent service
- **Strands Agents** framework with Ollama and OpenAI-compatible integration
- **Docker Compose** for service orchestration
- **Ollama** for model hosting and serving
- **llama.cpp** for CPU-only inference backend

**Quick Start:**

```bash
cd AWS-Strands
make setup
```

### 2. LangGraph (`/langgraph/`)

A Python project built with LangGraph and Pydantic AI for creating intelligent agent workflows:

- **LangGraph** for building stateful, multi-agent applications
- **Pydantic AI** for type-safe AI agent development
- **Structured Logging** with structlog
- **Configuration Management** with Pydantic Settings
- **Async Support** for efficient execution
- **Comprehensive Testing** with pytest

**Quick Start:**

```bash
cd langgraph
uv sync
```

## Prerequisites

- Python 3.10+
- Docker and Docker Compose (for AWS-Strands)
- UV package manager (for LangGraph)
- Git for version control

## Development Setup

### Root Project Configuration

The root directory contains:

- `pyproject.toml` - Project-wide tool configurations (Black, isort, Ruff, MyPy, pytest)
- `.agents/` - AI agent rules and configurations
- `.venv/` - Python virtual environment
- `.vscode/` - VS Code workspace settings

### Tool Configurations

The project uses the following development tools with consistent configurations:

- **Black** - Code formatting (line length: 120)
- **isort** - Import sorting (compatible with Black)
- **Ruff** - Fast linting (selects E, W, F, I, B, C4, UP rules)
- **MyPy** - Static type checking (strict mode)
- **pytest** - Testing framework

## Getting Started

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd learning-ai-agents
   ```

2. **Set up the virtual environment:**

   ```bash
   python -m venv .venv
   source .venv/bin/activate  # On Windows: .venv\Scripts\activate
   ```

3. **Explore the sub-projects:**

   - Navigate to `AWS-Strands/` for FastAPI-based agent service
   - Navigate to `langgraph/` for LangGraph-based agent workflows

## Project Goals

This repository serves as a learning resource for:

- Understanding different AI agent frameworks
- Experimenting with various LLM integration patterns
- Building production-ready agent services
- Learning best practices in AI agent development

## License

This project is for educational purposes. Please refer to individual sub-projects for their specific licenses.

## Resources

- [LangGraph Documentation](https://langchain-ai.github.io/langgraph/)
- [Pydantic AI Documentation](https://ai.pydantic.dev/)
- [Strands Agents SDK](https://strandsagents.com/)
- [Ollama Documentation](https://ollama.com/)
- [llama.cpp GitHub](https://github.com/ggerganov/llama.cpp)
