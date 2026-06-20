# Skill: Add a New Agent

Step-by-step guide for adding a new agent to the LangGraph project.

## Steps

1. **Create the agent module** in `src/langgraph_project/agents/`:

```python
from langgraph.graph import StateGraph, MessagesState
from langgraph_project.utils.config import get_config

# Define typed state
class AgentState(MessagesState):
    """State for this agent."""
    pass

# Build the graph
graph = StateGraph(AgentState)

# Add nodes
graph.add_node("process", process_node)
graph.add_node("respond", respond_node)

# Add edges
graph.add_edge("process", "respond")

# Set entry and finish points
graph.set_entry_point("process")
graph.set_finish_point("respond")

# Compile
agent = graph.compile()
```

2. **Export the agent** in `src/langgraph_project/agents/__init__.py`.

3. **Add tests** in `tests/test_<agent_name>.py`:
   - Test individual nodes
   - Test the full graph execution
   - Use `pytest-asyncio` for async graph runs

4. **Register the agent** in `src/langgraph_project/main.py` if it's a top-level agent.

5. **Run checks**: `mise run v`

## Conventions

- Use LangGraph's `StateGraph` with typed state objects
- Each node should be a pure function that takes state and returns state updates
- Use `structlog` for logging within nodes
- Configure the LLM via `utils/config.py`
- Keep nodes focused — one responsibility per node
