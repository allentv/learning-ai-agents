# Usage

## Running with Docker Compose (Recommended)

1. **Download the model**:

   ```bash
   mise run download-model
   ```

2. **Start all services** (llamacpp + API + app):

   ```bash
   mise run up
   ```

   This starts three services:
   - **llamacpp**: Local LLM server on port `12434`
   - **api**: FastAPI knowledge base API on port `10000`
   - **app**: LangGraph agent application on port `8080`

3. **Run the application**:

   ```bash
   mise run start
   ```

   This starts all services, waits for health checks, and runs the agent.

## Running Locally

```bash
mise run run
```

## Using the Agent Directly with Tools

```python
import asyncio
from langgraph_project.agents.simple_agent import SimpleAgent
from langgraph_project.tools import search_api

async def main():
    # Initialize agent with tools
    agent = SimpleAgent(tools=[search_api])

    # The agent will automatically call search_api when it needs information
    response = await agent.process("What is the latest news about AI?")
    print(response)

asyncio.run(main())
```

## Running the API Service Locally

```bash
mise run api
```

Test the API:
```bash
mise run api-test
