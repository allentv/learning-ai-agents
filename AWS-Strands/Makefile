.PHONY: setup setup-llamacpp pull pull-llamacpp stop clean ping logs

MODEL ?= ibm/granite4.1:3b
LLAMACPP_MODEL ?= Qwen3-0.6B-Q8_0.gguf
AGENT_URL = http://localhost:8080
OLLAMA_URL = http://localhost:11434
LLAMACPP_URL = http://localhost:8081

## Start all services (Ollama + agent) and verify
setup:
	@docker compose up -d --build model agent
	@echo "⏳ Waiting for Ollama to be ready..."
	@until curl -sf $(OLLAMA_URL)/api/tags > /dev/null 2>&1; do sleep 2; done
	@echo "📥 Pulling model $(MODEL)..."
	docker compose exec model ollama pull $(MODEL)
	@echo "⏳ Waiting for agent to be ready..."
	@until curl -sf $(AGENT_URL)/ping > /dev/null 2>&1; do sleep 2; done
	@echo "✅ Setup complete! Agent is running at $(AGENT_URL)"

## Start services with llama.cpp backend (CPU-only, no GPU needed)
setup-llamacpp:
	@MODEL_PROVIDER=llamacpp docker compose up -d --build agent
	@docker compose up -d llamacpp
	@echo "⏳ Waiting for llama.cpp to be ready..."
	@until curl -sf $(LLAMACPP_URL)/v1/models > /dev/null 2>&1; do sleep 5; done
	@echo "⏳ Waiting for agent to be ready..."
	@until curl -sf $(AGENT_URL)/ping > /dev/null 2>&1; do sleep 2; done
	@echo "✅ Setup complete! Agent is running at $(AGENT_URL) with llama.cpp backend"

## Pull or update the LLM model (Ollama)
pull:
	docker compose exec model ollama pull $(MODEL)

## Pull or update the GGUF model (llama.cpp)
pull-llamacpp:
	@echo "📥 Downloading GGUF model $(LLAMACPP_MODEL)..."
	@mkdir -p llamacpp_models
	@curl -fSL -o "llamacpp_models/$(LLAMACPP_MODEL)" \
		"https://huggingface.co/Qwen/Qwen3-0.6B-GGUF/resolve/main/Qwen3-0.6B-Q8_0.gguf"
	@echo "✅ Model downloaded to llamacpp_models/$(LLAMACPP_MODEL)"

## Check agent health
ping:
	@curl -s $(AGENT_URL)/ping

## View logs
logs:
	docker compose logs -f

## Stop all services
stop:
	docker compose down

## Stop all services and remove volumes (clears downloaded models)
clean:
	docker compose down -v

## Show available models in Ollama
models:
	docker compose exec model ollama list
