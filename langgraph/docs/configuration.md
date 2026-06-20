# Configuration

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

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MODEL_PROVIDER` | Model provider (`openai` or `llamacpp`) | `llamacpp` |
| `OPENAI_API_KEY` | OpenAI API key (if using OpenAI) | - |
| `OPENAI_MODEL` | OpenAI model name | `gpt-4o-mini` |
| `LLAMACPP_URL` | llama.cpp API base URL | `http://llamacpp:12434/v1` |
| `LLAMACPP_MODEL` | llama.cpp model name | `granite-4.0-h-micro-UD-Q4_K_XL.gguf` |
| `API_BASE_URL` | External API base URL for agent tools | `http://api:10000` |
| `API_TIMEOUT` | API request timeout in seconds | `30` |
| `LOG_LEVEL` | Logging level | `INFO` |
| `LOG_FORMAT` | Log format (`json` or `text`) | `json` |
| `APP_NAME` | Application name | `LangGraph Project` |
| `APP_VERSION` | Application version | `0.1.0` |

## Model Providers

The project supports multiple model providers:

- **llamacpp**: Local llama.cpp server (recommended for local development)
- **openai**: OpenAI API
