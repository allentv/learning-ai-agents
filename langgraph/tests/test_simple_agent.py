"""Tests for the simple agent."""

import os

import pytest

from langgraph_project.agents.simple_agent import SimpleAgent
from langgraph_project.utils.config import settings


class TestSimpleAgent:
    """Test cases for SimpleAgent."""

    @pytest.mark.skipif(
        settings.model_provider == "openai" and not os.getenv("OPENAI_API_KEY"),
        reason="OPENAI_API_KEY environment variable not set",
    )
    def test_agent_initialization(self):
        """Test that the agent initializes correctly."""
        agent = SimpleAgent()
        assert agent is not None
        assert agent.agent is not None

    @pytest.mark.asyncio
    @pytest.mark.skipif(
        settings.model_provider == "openai" and not os.getenv("OPENAI_API_KEY"),
        reason="OPENAI_API_KEY environment variable not set",
    )
    async def test_process_query(self):
        """Test processing a simple query."""
        agent = SimpleAgent()
        # Note: This test would require a valid OpenAI API key or vLLM server to run
        # For now, we'll just test the structure
        assert hasattr(agent, "process")
        assert callable(agent.process)

    def test_vllm_mode(self):
        """Test that vLLM mode is properly configured."""
        # Set model provider to vllm for this test
        original_provider = settings.model_provider
        settings.model_provider = "vllm"

        try:
            agent = SimpleAgent()
            assert agent is not None
            assert agent.agent is not None
        finally:
            # Restore original provider
            settings.model_provider = original_provider

    def test_model_runner_mode(self):
        """Test that model-runner mode is properly configured."""
        # Set model provider to model-runner for this test
        original_provider = settings.model_provider
        settings.model_provider = "model-runner"

        try:
            agent = SimpleAgent()
            assert agent is not None
            assert agent.agent is not None
        finally:
            # Restore original provider
            settings.model_provider = original_provider
