# Docker Compose Setup for LangGraph with vLLM

This document describes how to run the LangGraph project with locally hosted LLMs using vLLM in Docker Compose.

## Prerequisites

- Docker and Docker Compose installed
- NVIDIA GPU (optional but recommended for better performance)
- At least 16GB RAM (for CPU-only mode)

## Quick Start

1. **Clone the repository** (if not already done):

   ```bash
   git clone https://github.com/allentv/learning-ai-agents.git
   cd learning-ai-agents/langgraph
   ```

2. **Create environment file**:

   ```bash
   cp .env.example .env
   ```

   Edit `.env` to set `MODEL_PROVIDER=model-runner` (`openai` for OpenAI API).

3. **Start services**:

   ```bash
   docker-compose up -d
   ```

4. **Check logs**:

   ```bash
   docker-compose logs -f
   ```

5. **Test the application**:

   ```bash
   # Run the test suite
   docker-compose exec app uv run pytest tests/

   # Or run a simple query
   docker-compose exec app uv run python -c "
   from langgraph_project.agents.simple_agent import SimpleAgent
   import asyncio
   async def test():
       agent = SimpleAgent()
       result = await agent.process('What is the capital of France?')
       print(result)
   asyncio.run(test())
   "
   ```

## GPU Configuration

For GPU acceleration, uncomment the GPU section in `docker-compose.yaml`:

```yaml
deploy:
  resources:
    reservations:
      devices:
        - driver: nvidia
          count: 1
          capabilities: [gpu]
```

You'll also need:

- NVIDIA Container Toolkit installed
- NVIDIA drivers on your host system

## Model Configuration

The default model is `hf.co/unsloth/granite-4.0-h-micro-GGUF:UD-Q4_K_XL` (model-runner) or `ibm-granite/granite-4.1-3b` (vLLM). To use a different model:

1. Update `MODEL_RUNNER_MODEL` or `VLLM_MODEL` in `.env`
2. Update the `provider.options.model` in the `model-runner` service or the `command` in the `vllm` service in `docker-compose.yaml`

## Troubleshooting

### Model Runner service fails to start

- Check that the model is available and accessible
- Ensure you have enough disk space for the model
- Check network connectivity to download the model if needed

### vLLM service fails to start

- Check GPU availability: `nvidia-smi`
- Ensure you have enough VRAM (at least 6GB for granite-4.1-3b)
- Try CPU-only mode by removing GPU configuration

### Model download issues

- The model will be downloaded automatically on first run
- Check disk space (model is ~3GB)
- Set `HF_HUB_CACHE` volume if you want to persist models

### Port conflicts

- vLLM runs on port 8000 by default
- App runs on port 8080 by default
- Change ports in `docker-compose.yaml` if needed

## Stopping Services

```bash
docker-compose down
```

To remove volumes (including downloaded models):

```bash
docker-compose down -v
```

## Development

For local development without Docker:

1. Install dependencies:

   ```bash
   uv sync
   ```

2. Run vLLM locally:

   ```bash
   python -m vllm.entrypoints.openai.api_server \
     --model ibm-granite/granite-4.1-3b \
     --host [IP_ADDRESS] \
     --port 8000
   ```

3. Run the application:

   ```bash
   MODEL_PROVIDER=vllm \
   VLLM_URL=http://localhost:8000/v1 \
   uv run python -m langgraph_project.main
   ```

## Performance Notes

- **CPU Mode**: Expect slower inference (seconds per token)
- **GPU Mode**: Much faster inference (milliseconds per token)
- **Model Size**: granite-4.1-3b is ~3GB, requires ~6GB RAM/VRAM
- **Context Length**: Default is 4096 tokens, can be increased with `--max-model-len`

## API Endpoints

- vLLM API: `http://localhost:8000/v1`
- Health check: `http://localhost:8000/health`
- OpenAI-compatible endpoint: `http://localhost:8000/v1/chat/completions`

## Environment Variables

| Variable          | Description                                | Default                              |
| ----------------- | ------------------------------------------ | ------------------------------------ |
| `MODEL_PROVIDER`  | Model provider (`openai` or `vllm`)        | `openai`                             |
| `VLLM_URL`        | vLLM API base URL                          | `http://localhost:8000/v1`           |
| `VLLM_MODEL`      | vLLM model name                            | `ibm-granite/granite-4.1-3b`         |
| `OPENAI_API_KEY`  | OpenAI API key (if using OpenAI)           | -                                    |
| `OPENAI_MODEL`    | OpenAI model name                          | `gpt-4o-mini`                        |

## License

This project is licensed under the MIT License.
