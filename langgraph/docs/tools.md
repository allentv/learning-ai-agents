# Tools

The agent uses [pydantic-ai](https://ai.pydantic.dev/) for tool integration. Tools are plain async Python functions with type-annotated parameters — pydantic-ai automatically generates the JSON schema for the LLM.

## Available Tools

| Tool | Description | Endpoint |
|------|-------------|----------|
| `search_api` | Search the knowledge base via the external API | `GET /api/search?q=<query>` |

## Adding a New Tool

1. **Create the tool function** in `src/langgraph_project/tools/api_tools.py` (or a new file):

   ```python
   async def my_new_tool(param: str) -> dict[str, Any]:
       """Description of what the tool does.

       This docstring is used as the tool description for the LLM.

       Args:
           param: Description of the parameter.
       """
       # Implementation here
       return {"result": "..."}
   ```

2. **Export it** in `src/langgraph_project/tools/__init__.py`:

   ```python
   from langgraph_project.tools.api_tools import search_api, my_new_tool

   __all__ = ["search_api", "my_new_tool"]
   ```

3. **Register it** with the agent in `src/langgraph_project/main.py`:

   ```python
   from langgraph_project.tools import search_api, my_new_tool

   AGENT_TOOLS = [search_api, my_new_tool]
   ```

## How Tool Calls Work

1. The user sends a query to the agent
2. The LLM decides whether to call a tool based on the tool schemas
3. If a tool call is needed, pydantic-ai executes the tool function
4. The tool result is fed back to the LLM
5. The LLM may call more tools or produce a final text response
6. This loop continues until the LLM produces a final answer

This entire loop is handled automatically by pydantic-ai within the `process_query` node of the LangGraph graph.
