# Use uv's Python base image
FROM ghcr.io/astral-sh/uv:python3.11-bookworm-slim

WORKDIR /app

# Copy uv files
COPY pyproject.toml uv.lock ./

# Install dependencies
RUN uv sync --frozen --no-cache

# Copy agent file
COPY agent.py ./

# Default Ollama host for Docker (agent connects to model service by name)
ENV OLLAMA_HOST=http://model:11434
ENV OLLAMA_MODEL=ibm/granite4.1:3b
ENV LLAMACPP_HOST=http://llamacpp:8080
ENV LLAMACPP_MODEL=Qwen3-0.6B-Q8_0.gguf
ENV MODEL_PROVIDER=ollama

# Expose port
EXPOSE 8080

# Run application
CMD ["uv", "run", "python", "agent.py"]
