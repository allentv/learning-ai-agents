# Agent Patterns for Python Development

## Core Agent Architecture

### Agent Definition Pattern

```python
from langgraph.graph import StateGraph, END
from typing import TypedDict, Annotated
from langchain_core.messages import HumanMessage, AIMessage

class AgentState(TypedDict):
    messages: Annotated[list, lambda x, y: x + y]
    context: dict

def agent_node(state: AgentState) -> AgentState:
    """Process agent state and return updated state."""
    # Agent logic here
    return state

# Build the graph
workflow = StateGraph(AgentState)
workflow.add_node("agent", agent_node)
workflow.add_edge("agent", END)
workflow.set_entry_point("agent")
app = workflow.compile()
```

### Tool Integration Pattern

```python
from langchain_core.tools import tool
from langchain_core.messages import ToolMessage

@tool
def search_tool(query: str) -> str:
    """Search for information."""
    return f"Results for: {query}"

def tool_node(state: AgentState) -> AgentState:
    """Execute tools based on agent decisions."""
    messages = state["messages"]
    last_message = messages[-1]
    
    if hasattr(last_message, "tool_calls"):
        # Execute tool calls
        pass
    
    return state
```

## State Management

### Immutable State Pattern

```python
from dataclasses import dataclass
from typing import Optional

@dataclass(frozen=True)
class AgentConfig:
    model: str
    temperature: float
    max_tokens: int

@dataclass(frozen=True)
class AgentState:
    config: AgentConfig
    history: tuple
    metadata: dict
```

### State Validation Pattern

```python
from pydantic import BaseModel, validator
from typing import List

class ConversationState(BaseModel):
    messages: List[str]
    context: dict
    
    @validator('messages')
    def validate_messages(cls, v):
        if len(v) > 1000:
            raise ValueError('Too many messages')
        return v
```

## Error Handling

### Retry Pattern

```python
import asyncio
from typing import Callable, Any
from functools import wraps

def retry_with_backoff(
    max_retries: int = 3,
    base_delay: float = 1.0
) -> Callable:
    """Decorator for retrying functions with exponential backoff."""
    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs) -> Any:
            for attempt in range(max_retries):
                try:
                    return await func(*args, **kwargs)
                except Exception as e:
                    if attempt == max_retries - 1:
                        raise
                    delay = base_delay * (2 ** attempt)
                    await asyncio.sleep(delay)
        return wrapper
    return decorator
```

### Graceful Degradation Pattern

```python
from typing import Optional, Union

def safe_agent_execution(
    agent_func: Callable,
    fallback_func: Callable,
    state: AgentState
) -> AgentState:
    """Execute agent with fallback on failure."""
    try:
        return agent_func(state)
    except Exception as e:
        logger.warning(f"Agent failed: {e}, using fallback")
        return fallback_func(state)
```

## Testing Patterns

### Agent Testing Pattern

```python
import pytest
from unittest.mock import Mock, patch

def test_agent_state_transition():
    """Test agent state transitions."""
    initial_state = AgentState(messages=[], context={})
    
    with patch('agent_module.search_tool') as mock_search:
        mock_search.return_value = "test results"
        final_state = agent_node(initial_state)
        
        assert len(final_state["messages"]) > 0

@pytest.mark.asyncio
async def test_async_agent():
    """Test async agent functionality."""
    result = await async_agent_function()
    assert result is not None
```

## Performance Patterns

### Caching Pattern

```python
from functools import lru_cache
from typing import Any

@lru_cache(maxsize=100)
def cached_tool_call(tool_name: str, query: str) -> Any:
    """Cached tool execution."""
    # Expensive operation here
    return result
```

### Batch Processing Pattern

```python
from typing import List
import asyncio

async def batch_process_queries(queries: List[str]) -> List[str]:
    """Process multiple queries concurrently."""
    tasks = [process_query(q) for q in queries]
    return await asyncio.gather(*tasks)
```

## Best Practices

1. **Keep agents focused**: Each agent should have a single responsibility
2. **Validate inputs**: Always validate agent inputs before processing
3. **Log extensively**: Use structured logging for debugging
4. **Test thoroughly**: Write unit tests for all agent components
5. **Monitor performance**: Track agent execution times and success rates
6. **Use uv for dependencies**: Always use `uv` for package management and environment setup
