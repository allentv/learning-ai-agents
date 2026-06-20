"""Simple agent implementation using Pydantic AI."""

from collections.abc import Callable
from typing import Any

import structlog
from pydantic_ai import Agent
from pydantic_ai.messages import ModelMessage

from langgraph_project.utils.config import settings

logger = structlog.get_logger()


class SimpleAgent:
    """A simple agent that processes queries using Pydantic AI.

    Supports optional tool registration. Tools are plain async Python
    functions with type-annotated parameters — pydantic-ai generates
    the JSON schema for the LLM automatically.
    """

    def __init__(self, tools: list[Callable[..., Any]] | None = None) -> None:
        """Initialize the simple agent.

        Args:
            tools: Optional list of tool functions to register with the agent.
                   Each tool should be an async function with type-annotated parameters.
        """
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
            system_prompt="You are a helpful assistant. Provide concise and accurate responses. "
            "When you need to look up information, use the available tools to search "
            "the knowledge base before answering.",
            tools=tools or [],
        )

        if tools:
            logger.info("Agent initialized with tools", tool_count=len(tools))

    async def process(self, query: str) -> str:
        """Process a query and return the response.

        The agent may make one or more tool calls before producing
        a final text response. pydantic-ai handles the tool call
        loop internally.
        """
        logger.info("Agent processing query", query=query)

        try:
            result = await self.agent.run(query)
            response: str = result.output
            return response
        except Exception as e:
            logger.error("Error processing query", error=str(e), error_type=type(e).__name__)
            # Provide a helpful error message based on the exception type
            error_type_name = type(e).__name__
            if "ConnectError" in error_type_name or "connect" in str(e).lower():
                hint = (
                    " Cannot connect to the LLM server. "
                    f"Ensure your model provider is running (provider={settings.model_provider}). "
                    "If using llamacpp, run: mise run up  (or start the llama.cpp server). "
                    "If using OpenAI, set OPENAI_API_KEY in your .env file."
                )
            elif "timeout" in str(e).lower():
                hint = " The request timed out. The server may be overloaded or unreachable."
            else:
                hint = ""
            error_msg: str = f"Error processing query: {str(e)}{hint}"
            return error_msg

    def get_messages(self) -> list[ModelMessage]:
        """Get the conversation history."""
        # Note: Agent.messages may not be available in all pydantic-ai versions
        # For now, return empty list
        messages: list[ModelMessage] = []
        return messages
