"""Tests for the simple agent."""

import os

import pytest

from langgraph_project.agents.simple_agent import SimpleAgent


class TestSimpleAgent:
    """Test cases for SimpleAgent."""

    @pytest.mark.skipif(not os.getenv("OPENAI_API_KEY"), reason="OPENAI_API_KEY environment variable not set")
    def test_agent_initialization(self):
        """Test that the agent initializes correctly."""
        agent = SimpleAgent()
        assert agent is not None
        assert agent.agent is not None

    @pytest.mark.asyncio
    @pytest.mark.skipif(not os.getenv("OPENAI_API_KEY"), reason="OPENAI_API_KEY environment variable not set")
    async def test_process_query(self):
        """Test processing a simple query."""
        agent = SimpleAgent()
        # Note: This test would require a valid OpenAI API key to run
        # For now, we'll just test the structure
        assert hasattr(agent, "process")
        assert callable(agent.process)
