"""Tests for the simple agent."""

import os
from typing import Any
from unittest.mock import AsyncMock, patch

import pytest

from langgraph_project.agents.simple_agent import SimpleAgent
from langgraph_project.utils.config import settings


class TestSimpleAgent:
    """Test cases for SimpleAgent."""

    @pytest.mark.skipif(
        settings.model_provider == "openai" and not os.getenv("OPENAI_API_KEY"),
        reason="OPENAI_API_KEY environment variable not set",
    )
    def test_agent_initialization(self) -> None:
        """Test that the agent initializes correctly."""
        agent = SimpleAgent()
        assert agent is not None
        assert agent.agent is not None

    @pytest.mark.skipif(
        settings.model_provider == "openai" and not os.getenv("OPENAI_API_KEY"),
        reason="OPENAI_API_KEY environment variable not set",
    )
    def test_agent_initialization_with_tools(self) -> None:
        """Test that the agent initializes correctly with tools."""

        async def dummy_tool(query: str) -> dict[str, Any]:
            """A dummy tool for testing."""
            return {"result": query}

        agent = SimpleAgent(tools=[dummy_tool])
        assert agent is not None
        assert agent.agent is not None

    @pytest.mark.asyncio
    @pytest.mark.skipif(
        settings.model_provider == "openai" and not os.getenv("OPENAI_API_KEY"),
        reason="OPENAI_API_KEY environment variable not set",
    )
    async def test_process_query(self) -> None:
        """Test processing a simple query."""
        agent = SimpleAgent()
        # Note: This test would require a valid OpenAI API key or llama.cpp server to run
        # For now, we'll just test the structure
        assert hasattr(agent, "process")
        assert callable(agent.process)

    def test_llamacpp_mode(self) -> None:
        """Test that llama.cpp mode is properly configured."""
        # Set model provider to llamacpp for this test
        original_provider = settings.model_provider
        settings.model_provider = "llamacpp"

        try:
            agent = SimpleAgent()
            assert agent is not None
            assert agent.agent is not None
        finally:
            # Restore original provider
            settings.model_provider = original_provider

    def test_agent_no_tools_by_default(self) -> None:
        """Test that the agent initializes with no tools by default."""
        agent = SimpleAgent()
        assert agent is not None
