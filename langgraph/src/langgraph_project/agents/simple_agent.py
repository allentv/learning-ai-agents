"""Simple agent implementation using Pydantic AI."""

import structlog
from pydantic_ai import Agent
from pydantic_ai.messages import ModelMessage

logger = structlog.get_logger()


class SimpleAgent:
    """A simple agent that processes queries using Pydantic AI."""

    def __init__(self):
        """Initialize the simple agent."""
        self.agent = Agent(
            model="openai:gpt-4o-mini",
            system_prompt="You are a helpful assistant. Provide concise and accurate responses.",
        )

    async def process(self, query: str) -> str:
        """Process a query and return the response."""
        logger.info("Agent processing query", query=query)

        try:
            result = await self.agent.run(query)
            response: str = result.data
            return response
        except Exception as e:
            logger.error("Error processing query", error=str(e))
            error_msg: str = f"Error processing query: {str(e)}"
            return error_msg

    def get_messages(self) -> list[ModelMessage]:
        """Get the conversation history."""
        messages: list[ModelMessage] = self.agent.messages
        return messages
