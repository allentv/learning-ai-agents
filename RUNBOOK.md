# Runbook — Learning AI Agents

Operational reference for the AI agent Docker Compose setup.
For setup instructions, see [README.md](README.md).

---

## Health Checks

**Check if the agent is healthy:**

```zsh
curl http://localhost:8080/ping
```

Expected response:

```json
{"status": "healthy"}
```

**Check if the Ollama model service is responding:**

```zsh
curl http://localhost:11434/api/tags
```

**Check if the llama.cpp service is responding:**

```zsh
curl http://localhost:8081/v1/models
```

---

## Logs

**Tail logs for all services:**

```zsh
docker compose logs -f
```

**Tail logs for a specific service:**

```zsh
docker compose logs -f agent
docker compose logs -f model
```

**View last N lines of logs:**

```zsh
docker compose logs --tail 50 agent
```

---

## Common Issues

### Port Already in Use

**Symptom:** `error: address already in use` for port 8080 or 11434.

**Resolution:**

1. Find the process using the port:

   ```zsh
   ss -tlnp sport = :8080
   lsof -i :8080
   ```

2. Stop the conflicting process, or change the port mapping in `docker-compose.yaml`:

   ```yaml
   ports:
     - "9090:8080"  # map to a different host port
   ```

3. Restart the services:

   ```zsh
   docker compose down && docker compose up -d
   ```

### Agent Can't Reach the Model Service

**Symptom:** Agent returns 500 errors or connection refused to Ollama.

**Resolution:**

1. Ensure the model service is running:

   ```zsh
   docker compose ps
   ```

2. The `agent` service connects to Ollama via `http://model:11434` (Docker service name). If running outside Docker or the agent runs locally, ensure Ollama is accessible at `http://localhost:11434`.

3. Verify the model is downloaded:

   ```zsh
   docker compose exec model ollama list
   ```

### Agent Can't Reach the llama.cpp Service

**Symptom:** Agent returns 500 errors or connection refused to llama.cpp.

**Resolution:**

1. Ensure the llama.cpp service is running:

   ```zsh
   docker compose ps
   ```

2. The `agent` service connects to llama.cpp via `http://llamacpp:8080` (Docker service name). If running outside Docker or the agent runs locally, ensure llama.cpp is accessible at `http://localhost:8081`.

3. Verify the GGUF model is loaded:

   ```zsh
   curl http://localhost:8081/v1/models
   ```

### Model Not Found

**Symptom:** Errors referencing a missing or unknown model.

**Resolution:**

1. Pull the model:

   ```zsh
   docker compose exec model ollama pull ibm/granite4.1:3b
   ```

2. Verify available models:

   ```zsh
   docker compose exec model ollama list
   ```

### llama.cpp Model Not Found

**Symptom:** `gguf_init_from_file: failed to open GGUF file` error when starting llama.cpp.

**Resolution:**

The GGUF model is downloaded automatically on first run. If you see this error:

1. Check if the model file exists:

   ```zsh
   docker volume inspect learning-ai-agents_llamacpp_models
   ```

2. Manually download the model:

   ```zsh
   make pull-llamacpp
   ```

3. Use a different GGUF model:

   ```zsh
   LLAMACPP_MODEL=Meta-Llama-3.2-1B-Instruct-Q4_K_M.gguf \
   LLAMACPP_HF_REPO=mradermacher/Meta-Llama-3.2-1B-Instruct-GGUF \
   LLAMACPP_HF_FILE=Meta-Llama-3.2-1B-Instruct.Q4_K_M.gguf \
   make setup-llamacpp
   ```

### Out of Memory / Slow Responses

**Symptom:** Container restarts unexpectedly, or responses take too long.

**Resolution:**

- Use a smaller model (e.g., `llama3:8b` instead of larger variants).
- Check Docker resource limits in Docker Desktop or via `docker stats`.
- Reduce `max_tokens` or `temperature` in `agent.py` for lighter inference.

---

## Data Persistence

**Ollama models** and state are stored in the `model_data` Docker volume, mapped to `/root/.ollama` inside the container.

**llama.cpp models** are stored in the `llamacpp_models` Docker volume, mapped to `/models` inside the container.

**Inspect the volumes:**

```zsh
docker volume inspect learning-ai-agents_model_data
docker volume inspect learning-ai-agents_llamacpp_models
```

**Check disk usage:**

```zsh
docker system df -v
```

**Nuclear option — clear all models and start fresh:**

```zsh
docker compose down -v
docker compose up -d
# Then re-pull your model
docker compose exec model ollama pull granite4.1:3b
# Or re-download the GGUF model
make pull-llamacpp
```

---

## Starting and Stopping

**Setup everything in one command:**

```zsh
make setup
make setup MODEL=llama3        # with a specific model
```

**Setup with llama.cpp (CPU-only, no GPU needed):**

```zsh
make setup-llamacpp
```

**All Makefile targets:**

| Command              | Description                              |
| -------------------- | ---------------------------------------- |
| `make setup`         | Start services, pull model, verify health|
| `make setup-llamacpp`| Start llama.cpp backend (CPU inference)  |
| `make pull`          | Pull/update the Ollama model             |
| `make pull-llamacpp` | Download/update the GGUF model           |
| `make ping`          | Check agent health                       |
| `make models`        | List models in Ollama                    |
| `make logs`          | Tail logs from all services              |
| `make stop`          | Stop all services                        |
| `make clean`         | Stop services and remove volumes         |

**Start services (detached):**

```zsh
docker compose up -d
```

**Rebuild and restart after code changes:**

```zsh
docker compose up --build
```

**Stop and remove containers:**

```zsh
docker compose down
```

**Stop and remove everything (including volumes):**

```zsh
docker compose down -v
```

---

## Local Development

To run the agent locally without Docker:

```zsh
source .venv/bin/activate
uv run agent.py
```

Ensure Ollama is running locally on port 11434 (install from [ollama.com](https://ollama.com)).  
The agent reads `OLAMA_API_KEY` from the environment if set.

---

## API Reference

| Endpoint      | Method | Description                      |
| ------------- | ------ | -------------------------------- |
| `/ping`       | GET    | Health check                     |
| `/invocations`| POST   | Invoke the agent with a prompt   |

**Example invocation:**

```zsh
curl -X POST http://localhost:8080/invocations \
  -H "Content-Type: application/json" \
  -d '{"input": {"prompt": "What is the capital of France?"}}'
```
