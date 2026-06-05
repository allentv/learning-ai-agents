"""Main entry point for the LangGraph project."""

import asyncio
from typing import Any

import structlog
from langgraph.graph import END, StateGraph
from pydantic import BaseModel

from .agents.simple_agent import SimpleAgent
from .utils.logging import setup_logging

# Configure structured logging
setup_logging()
logger = structlog.get_logger()


class AgentState(BaseModel):
    """State model for the agent."""

    query: str = ""
    response: str = ""
    metadata: dict[str, Any] = {}


async def process_query(state: AgentState) -> AgentState:
    """Process a query using the agent."""
    logger.info("Processing query", query=state.query)

    agent = SimpleAgent()
    response = await agent.process(state.query)

    state.response = response
    logger.info("Query processed", response=response)

    return state


def create_workflow() -> StateGraph:
    """Create the LangGraph workflow."""
    workflow = StateGraph(AgentState)

    workflow.add_node("process_query", process_query)
    workflow.set_entry_point("process_query")
    workflow.add_edge("process_query", END)

    return workflow


async def main():
    """Main async function."""
    logger.info("Starting LangGraph project")

    # Create workflow
    workflow = create_workflow()
    app = workflow.compile()

    # Example query
    initial_state = AgentState(query="Hello, how are you?")

    # Run the workflow
    result = await app.ainvoke(initial_state)

    logger.info("Workflow completed", result=result)


if __name__ == "__main__":
    asyncio.run(main())
