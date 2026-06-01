# Learning AI Agents

This project is a setup for running an AI agent powered by Ollama or llama.cpp.

## Architecture

- **Agent Service**: A Python-based agent built with FastAPI and Strands Agents SDK.
- **Model Service**: An Ollama instance that hosts and serves the LLM.
- **llama.cpp Service**: A CPU-only inference backend using GGUF models.

### Tech Stack

- **Python** (>=3.11)
- **FastAPI** — Web framework for the agent service
- **Strands Agents** — Agent framework with Ollama and OpenAI-compatible integration
- **Docker Compose** — Service orchestration
- **Ollama** — Model hosting and serving
- **llama.cpp** — CPU-only inference backend

## Getting Started

### Prerequisites

- Docker and Docker Compose
- An `.env` file with the following variables:
  - `OLAMA_API_KEY`: Your API key (if applicable).
  - `OLLAMA_MODEL`: The model to use (defaults to `ibm/granite4.1:3b`).
  - `MODEL_PROVIDER`: Choose between `ollama` (default) or `llamacpp` for CPU-only inference.

### Installation & Running

**Quick start** (recommended):

```zsh
make setup
```

This starts the services, pulls the LLM, and verifies everything is healthy.

To use a different model:

```zsh
make setup MODEL=llama3
```

**CPU-only inference** (no GPU needed):

```zsh
make setup-llamacpp
```

This uses llama.cpp for CPU-only inference with GGUF models.

**Manual steps** (if you prefer not to use `make`):

1. **Start the services**:

   ```zsh
   docker compose up -d
   ```

2. **Download the LLM model**:

   Ollama doesn't download models automatically on startup. You need to pull a model into the container before the agent can use it.

   **Pull the configured model** (set via `OLLAMA_MODEL` in your `.env` file):

   ```zsh
   docker compose exec model ollama pull ${OLLAMA_MODEL:-ibm/granite4.1:3b}
   ```

   **Or pull a specific model** directly:

   ```zsh
   docker compose exec model ollama pull llama3
   ```

   See [RUNBOOK.md](RUNBOOK.md) for listing, removing, and managing models.

   **For llama.cpp backend**: The GGUF model is automatically downloaded on first run. See [RUNBOOK.md](RUNBOOK.md) for manual download instructions.

3. **Access the Agent**:
   The agent is available at `http://localhost:8080`.

4. **Useful Commands**:

   See [RUNBOOK.md](RUNBOOK.md) for logs, troubleshooting, and operational procedures.

## Project Structure

```plaintext
learning-ai-agents/
├── agent.py              # Main agent application with FastAPI
├── docker-compose.yaml   # Docker Compose configuration
├── Dockerfile           # Agent service container definition
├── Makefile             # Build and deployment commands
├── pyproject.toml       # Python dependencies
├── uv.lock              # Dependency lock file
├── README.md            # This file
├── RUNBOOK.md           # Operational reference
├── main.py              # Simple entry point
└── llamacpp_models/     # GGUF models directory (local development)
```

## Configuration

The project uses a `docker-compose.yaml` file to manage services.

- The `model` service persists data in a Docker volume named `model_data` mapped to `/root/.ollama`.
- The `llamacpp` service persists GGUF models in a Docker volume named `llamacpp_models` mapped to `/models`.
- The `agent` service is built from the local `Dockerfile`.

### Environment Variables

- `MODEL_PROVIDER`: Choose between `ollama` (default) or `llamacpp`
- `OLLAMA_MODEL`: Ollama model to use (default: `ibm/granite4.1:3b`)
- `LLAMACPP_MODEL`: GGUF model filename (default: `Qwen3-0.6B-Q8_0.gguf`)
- `LLAMACPP_HF_REPO`: HuggingFace repository for GGUF model download
- `LLAMACPP_HF_FILE`: HuggingFace filename for GGUF model download
- `LLAMACPP_THREADS`: Number of CPU threads for llama.cpp (default: 4)

## Local Development

You can also run the agent locally without Docker (ensure Ollama is running separately):

```zsh
source .venv/bin/activate
uv run agent.py
```

Or using the Makefile approach:

```zsh
make setup  # This will set up the environment and start services
```

**Note**: The project uses `uv` for dependency management. The virtual environment is created automatically when you run `make setup` or `uv sync`.
