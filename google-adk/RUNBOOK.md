# Runbook — Google ADK Go Project

Operational reference for the Google ADK Go implementation using the official Google ADK for Go.
For setup instructions, see [README.md](README.md).

**Official ADK Resources:**

- Repository: <https://github.com/google/adk-go>
- Documentation: <https://google.github.io/adk-docs/>
- Package: `google.golang.org/adk`

---

## Health Checks

**Check if the application is running:**

```bash
go run main.go
```

Expected output (with GOOGLE_API_KEY set):

```text
Starting Google ADK Go Project v0.1.0
Using model provider: gemini
Agent response: [Response from Gemini model]
Google ADK for Go application completed successfully
```

**Note:** You need to set the `GEMINI_API_KEY` environment variable with your Google AI API key. The `.env` file is automatically loaded using the envconfig package.

**Check if llama.cpp service is responding (if using llama.cpp provider):**

```bash
curl http://localhost:12434/v1/models
```

**Check if OpenAI API is accessible (if using OpenAI provider):**

```bash
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

---

## Logs

**Run with debug logging:**

```bash
export LOG_LEVEL=DEBUG
go run main.go
```

**View application output:**

The application logs to stdout by default. Use environment variables to control log level:

```bash
export LOG_LEVEL=INFO  # DEBUG, INFO, WARN, or ERROR
go run main.go
```

---

## Common Issues

### Model Provider Not Found

**Symptom:** `unsupported model provider: <provider>`

**Resolution:**

1. Check the `MODEL_PROVIDER` environment variable:

   ```bash
   echo $MODEL_PROVIDER
   ```

2. Ensure it's set to one of: `openai`, `gemini`, or `llamacpp`

3. Set the correct provider:

   ```bash
   export MODEL_PROVIDER=llamacpp
   ```

### API Key Not Set

**Symptom:** Authentication errors when using OpenAI or Gemini

**Resolution:**

1. Check if the API key is set:

   ```bash
   echo $OPENAI_API_KEY  # for OpenAI
   echo $GEMINI_API_KEY  # for Gemini
   ```

2. Set the API key:

   ```bash
   export OPENAI_API_KEY=your-key-here
   ```

3. Or add to `.env` file:

   ```env
   OPENAI_API_KEY=your-key-here
   ```

### llama.cpp Service Not Running

**Symptom:** Connection refused when using llama.cpp provider

**Resolution:**

1. Ensure llama.cpp is running:

   ```bash
   curl http://localhost:12434/v1/models
   ```

2. If not running, start it:

   ```bash
   # Example using llama.cpp server
   ./llama-server -m models/granite-4.0-h-micro-UD-Q4_K_XL.gguf --port 12434
   ```

3. Check the `LLAMACPP_URL` configuration:

   ```bash
   echo $LLAMACPP_URL
   ```

### Build Errors

**Symptom:** `go build` or `go run` fails with dependency errors

**Resolution:**

1. Clean the module cache:

   ```bash
   go clean -modcache
   ```

2. Tidy dependencies:

   ```bash
   go mod tidy
   ```

3. Try building again:

   ```bash
   go build -o google-adk main.go
   ```

---

## Development Workflow

### Adding a New Tool

1. Create a new tool file in `tools/`:

   ```go
   package tools

   import (
       "context"
       "fmt"
   )

   type MyTool struct {
       *BaseTool
   }

   func NewMyTool() *MyTool {
       return &MyTool{
           BaseTool: NewBaseTool("mytool", "Description of my tool"),
       }
   }

   func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
       // Implementation here
       return map[string]interface{}{
           "result": "success",
       }, nil
   }
   ```

2. Register the tool in your agent:

   ```go
   toolRegistry := tools.NewToolRegistry()
   toolRegistry.Register(tools.NewMyTool())
   ```

### Adding a New Model Provider

1. Extend the `Model` interface in `model/model.go`
2. Implement the new provider
3. Update the `New()` function to handle the new provider

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -ldflags="-s -w" -o google-adk main.go
```

---

## Deployment

### Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o google-adk main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/google-adk .
COPY --from=builder /app/.env .
CMD ["./google-adk"]
```

Build and run:

```bash
docker build -t google-adk .
docker run -it --rm google-adk
```

### Environment Variables

Ensure all required environment variables are set in production:

```bash
export MODEL_PROVIDER=llamacpp
export LLAMACPP_URL=http://llamacpp:12434/v1
export LLAMACPP_MODEL=granite-4.0-h-micro-UD-Q4_K_XL.gguf
export LOG_LEVEL=INFO
```
