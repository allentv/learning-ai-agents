"""Tools module for the LangGraph project.

Tools are defined as plain async Python functions with type-annotated
parameters. Pydantic-AI automatically generates JSON schemas from the
function signatures for the LLM to use.

To add a new tool:
1. Create a function in this package (e.g., in api_tools.py)
2. Import and re-export it here
3. Pass it to SimpleAgent(tools=[...]) in the agent setup

Example:
    from langgraph_project.tools import search_api

    agent = SimpleAgent(tools=[search_api])
"""

from langgraph_project.tools.api_tools import search_api

__all__ = ["search_api"]
