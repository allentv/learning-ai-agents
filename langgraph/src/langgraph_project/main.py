"""Main entry point for the LangGraph project."""

import asyncio
from typing import Any

import structlog
from langgraph.graph import END, StateGraph
from pydantic import BaseModel

from langgraph_project.agents.simple_agent import SimpleAgent
from langgraph_project.tools import search_api
from langgraph_project.utils.logging import setup_logging

# Configure structured logging
setup_logging()
logger = structlog.get_logger()


class AgentState(BaseModel):
    """State model for the agent."""

    query: str = ""
    response: str = ""
    metadata: dict[str, Any] = {}


# Tools available to the agent.
# Add new tool functions here to make them available to the agent.
AGENT_TOOLS = [search_api]


async def process_query(state: AgentState) -> AgentState:
    """Process a query using the agent.

    The agent is initialized with tools (e.g. search_api) so it can
    call external APIs when needed. pydantic-ai handles the tool call
    loop internally — the agent will call tools as many times as needed
    before producing a final text response.
    """
    logger.info("Processing query", query=state.query)

    agent = SimpleAgent(tools=AGENT_TOOLS)
    response = await agent.process(state.query)

    state.response = response
    logger.info("Query processed", response=response)

    return state


def create_workflow() -> StateGraph[AgentState]:
    """Create the LangGraph workflow.

    Current structure:
        process_query → END

    Tool calls are handled internally by pydantic-ai within the
    process_query node. If you need graph-level tool routing (e.g.
    for tools managed outside pydantic-ai), add a conditional edge
    here.
    """
    workflow: StateGraph[AgentState] = StateGraph(AgentState)

    workflow.add_node("process_query", process_query)
    workflow.set_entry_point("process_query")
    workflow.add_edge("process_query", END)

    return workflow


# Compile the graph at module level so LangGraph Studio / API server can discover it.
# The langgraph.json config references this as `main.py:graph`.
graph = create_workflow().compile()


async def main() -> None:
    """Main async function."""
    logger.info("Starting LangGraph project")

    # Example query
    initial_state = AgentState(query="Hello, how are you?")

    # Run the workflow
    result = await graph.ainvoke(initial_state)

    logger.info("Workflow completed", result=result)


if __name__ == "__main__":
    asyncio.run(main())
