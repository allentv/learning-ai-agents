import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, Any
from datetime import datetime, timezone
from strands import Agent

app = FastAPI(title="Strands Agent Server", version="1.0.0")

# Configure model backend based on MODEL_PROVIDER env var
# Supported: "ollama" or "llamacpp"
MODEL_PROVIDER = os.environ.get("MODEL_PROVIDER", "ollama").lower()

if MODEL_PROVIDER == "llamacpp":
    from strands.models.openai import OpenAIModel
    from openai import AsyncOpenAI

    llamacpp_host = os.environ.get("LLAMACPP_HOST", "http://localhost:8081")
    llamacpp_model = os.environ.get("LLAMACPP_MODEL", "Qwen3-0.6B-Q8_0.gguf")
    
    # Create a custom AsyncOpenAI client that works with llama.cpp
    custom_client = AsyncOpenAI(
        base_url=f"{llamacpp_host}/v1",
        api_key=os.environ.get("OPENAI_API_KEY", "sk-dummy"),
    )
    
    model = OpenAIModel(
        client=custom_client,
        model_id=llamacpp_model,
        params={
            "temperature": 0.85,
            "max_tokens": 4096,
        },
    )
    print(f"Using llama.cpp model: {llamacpp_model} at {llamacpp_host}")
else:
    from strands.models.ollama import OllamaModel

    ollama_host = os.environ.get("OLLAMA_HOST", "http://localhost:11434")
    ollama_model = os.environ.get("OLLAMA_MODEL", "ibm/granite4.1:3b")
    model = OllamaModel(model_id=ollama_model, host=ollama_host, temperature=0.85, max_tokens=4096)
    print(f"Using Ollama model: {ollama_model} at {ollama_host}")

strands_agent = Agent(model=model)

class InvocationRequest(BaseModel):
    input: Dict[str, Any]

class InvocationResponse(BaseModel):
    output: Dict[str, Any]

@app.post("/invocations", response_model=InvocationResponse)
async def invoke_agent(request: InvocationRequest):
    try:
        user_message = request.input.get("prompt", "")
        if not user_message:
            raise HTTPException(
                status_code=400,
                detail="No prompt found in input. Please provide a 'prompt' key in the input."
            )

        result = strands_agent(user_message)
        response = {
            "message": result.message,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "model": "strands-agent",
        }

        return InvocationResponse(output=response)

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Agent processing failed: {str(e)}")

@app.get("/ping")
async def ping():
    return {"status": "healthy"}

def main():
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080, log_level="debug")

if __name__ == "__main__":
    main()
