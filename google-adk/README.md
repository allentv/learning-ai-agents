# Google ADK Go Project

This project uses the **official Google ADK for Go** from <https://github.com/google/adk-go>.

## What is Google ADK for Go?

The **Agent Development Kit (ADK) for Go** is an open-source, code-first Go toolkit for building, evaluating, and deploying sophisticated AI agents with flexibility and control.

**Official Resources:**

- **Repository**: <https://github.com/google/adk-go>
- **Documentation**: <https://google.github.io/adk-docs/>
- **Package**: `google.golang.org/adk`
- **Examples**: <https://github.com/google/adk-go/tree/main/examples>

## Features

- **Idiomatic Go**: Designed to feel natural and leverage the power of Go
- **Rich Tool Ecosystem**: Utilize pre-built tools, custom functions, or integrate existing tools
- **Code-First Development**: Define agent logic, tools, and orchestration directly in Go
- **Modular Multi-Agent Systems**: Design scalable applications by composing multiple specialized agents
- **Deploy Anywhere**: Easily containerize and deploy agents with cloud-native support

## Project Structure

```text
google-adk/
├── config/                   # Configuration management (using envconfig)
│   └── config.go
├── logging/                  # Logging system (using Uber Zap)
│   └── logging.go
├── main.go                   # Main application entry point
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency checksums
├── README.md                 # This file
└── RUNBOOK.md                # Operational reference
```

## Installation

1. **Install Go** (if not already installed):

   ```bash
   # On macOS
   brew install go

   # On Linux
   sudo apt-get install golang-go
   ```

2. **Install dependencies**:

   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:

   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

4. **Get Google API Key**:

   - Visit <https://ai.google.dev/gemini-api/docs/api-key>
   - Create an API key
   - Add it to your `.env` file: `GEMINI_API_KEY=your-api-key-here`

## Configuration

This project uses [envconfig](https://github.com/kelseyhightower/envconfig) to load configuration from environment variables. The `.env` file is automatically loaded at startup.

Available configuration options:

- `MODEL_PROVIDER`: Model provider to use (default: "gemini")
- `GEMINI_API_KEY`: Google Gemini API key (required)
- `GEMINI_MODEL`: Gemini model to use (default: "gemini-3.1-flash-lite")
- `OPENAI_API_KEY`: OpenAI API key (optional)
- `OPENAI_MODEL`: OpenAI model to use (default: "gpt-4o-mini")
- `LLAMACPP_URL`: llama.cpp server URL (default: "<http://localhost:12434/v1>")
- `LLAMACPP_MODEL`: llama.cpp model to use (default: "granite-4.0-h-micro-UD-Q4_K_XL.gguf")
- `LOG_LEVEL`: Logging level (default: "INFO")
- `APP_NAME`: Application name (default: "Google ADK Go Project")
- `APP_VERSION`: Application version (default: "0.1.0")

Example `.env` file:

```env
# Model Provider Configuration
MODEL_PROVIDER=gemini

# Gemini API Configuration (required)
GEMINI_API_KEY=your-gemini-api-key-here
GEMINI_MODEL=gemini-3.1-flash-lite

# OpenAI API Configuration (optional)
OPENAI_API_KEY=your-openai-api-key-here
OPENAI_MODEL=gpt-4o-mini

# llama.cpp Configuration (optional)
LLAMACPP_URL=http://localhost:12434/v1
LLAMACPP_MODEL=granite-4.0-h-micro-UD-Q4_K_XL.gguf

# Logging Configuration
LOG_LEVEL=INFO

# Application Configuration
APP_NAME=Google ADK Go Project
APP_VERSION=0.1.0
```

## Usage

### Running the Application

```bash
go run main.go
```

### Using Different Model Providers

**OpenAI:**

```bash
export MODEL_PROVIDER=openai
export OPENAI_API_KEY=your-key
go run main.go
```

**Gemini:**

```bash
export MODEL_PROVIDER=gemini
export GEMINI_API_KEY=your-key
go run main.go
```

**llama.cpp:**

```bash
export MODEL_PROVIDER=llamacpp
export LLAMACPP_URL=http://localhost:12434/v1
go run main.go
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o google-adk main.go
```

### Code Structure

- **Agent**: The `Agent` struct manages the interaction between tools and models
- **Model**: Interface-based model providers (OpenAI, Gemini, llama.cpp)
- **Tools**: Extensible tool system with registry management
- **Workflow**: State-based workflow orchestration
- **Logging**: Structured logging with configurable levels

## License

MIT
