"""Simple agent implementation using Pydantic AI."""

import structlog
from pydantic_ai import Agent
from pydantic_ai.messages import ModelMessage

from langgraph_project.utils.config import settings

logger = structlog.get_logger()


class SimpleAgent:
    """A simple agent that processes queries using Pydantic AI."""

    def __init__(self) -> None:
        """Initialize the simple agent."""
        import os

        # Determine which model to use based on configuration
        if settings.model_provider == "llamacpp":
            model = f"openai-chat:{settings.llamacpp_model}"
            # Set the base URL for llama.cpp
            os.environ["OPENAI_BASE_URL"] = settings.llamacpp_url
            # llama.cpp doesn't require an API key, but the OpenAI client does
            # Set a dummy key if not already set
            if not os.environ.get("OPENAI_API_KEY"):
                os.environ["OPENAI_API_KEY"] = "llamacpp-dummy-key"
            logger.info(f"Using llama.cpp model: {settings.llamacpp_model} at {settings.llamacpp_url}")
        else:
            model = f"openai-chat:{settings.openai_model}"
            if settings.openai_api_key:
                os.environ["OPENAI_API_KEY"] = settings.openai_api_key
            logger.info(f"Using OpenAI model: {settings.openai_model}")

        self.agent = Agent(
            model=model,
            system_prompt="You are a helpful assistant. Provide concise and accurate responses.",
        )

    async def process(self, query: str) -> str:
        """Process a query and return the response."""
        logger.info("Agent processing query", query=query)

        try:
            result = await self.agent.run(query)
            response: str = result.output
            return response
        except Exception as e:
            logger.error("Error processing query", error=str(e))
            error_msg: str = f"Error processing query: {str(e)}"
            return error_msg

    def get_messages(self) -> list[ModelMessage]:
        """Get the conversation history."""
        # Note: Agent.messages may not be available in all pydantic-ai versions
        # For now, return empty list
        messages: list[ModelMessage] = []
        return messages
