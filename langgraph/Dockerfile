FROM python:3.14-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc g++ \
    curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy pyproject.toml and uv.lock for dependency installation
COPY pyproject.toml uv.lock README.md ./

# Install uv and dependencies
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/
RUN echo "Installing dependencies with uv..." && \
    uv sync --frozen && \
    echo "Dependencies installed successfully."

# Copy application code
COPY src/ ./src/
COPY tests/ ./tests/

# Set environment variables
ENV PYTHONPATH=/app/src
ENV PYTHONUNBUFFERED=1

# Run the application
CMD ["uv", "run", "python", "-m", "langgraph_project.main"]
