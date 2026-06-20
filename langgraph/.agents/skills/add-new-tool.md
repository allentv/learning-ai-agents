# Skill: Add a New Tool

Step-by-step guide for adding a new tool to the LangGraph project.

## Steps

1. **Create the tool function** in `src/langgraph_project/tools/` (create a new file or extend an existing one):

```python
from pydantic_ai import tool
from pydantic import BaseModel

class MyToolInput(BaseModel):
    query: str

@tool
def my_new_tool(args: MyToolInput) -> str:
    """Brief description of what this tool does."""
    # Implementation here
    return "result"
```

2. **Export the tool** in the module's `__init__.py`.

3. **Register the tool** with the relevant agent in `src/langgraph_project/agents/` by adding it to the agent's tool list.

4. **Add tests** in `tests/test_<tool_module>.py`:
   - Test the tool function directly
   - Mock any external HTTP calls with `pytest-httpx`
   - Use async tests where the tool is async

5. **Run checks**: `mise run v`

## Conventions

- Use `pydantic-ai`'s `@tool` decorator for type safety
- All inputs must be Pydantic models with typed fields
- All outputs should be typed (str, dict, Pydantic model)
- Use `httpx.AsyncClient` for any HTTP calls, never `requests`
- Use `structlog` for logging within tools
- Document the tool's purpose in its docstring
