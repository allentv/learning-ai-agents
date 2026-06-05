"""Base tool implementation for LangGraph project."""

from abc import ABC, abstractmethod
from typing import Any, Dict

import structlog

logger = structlog.get_logger()


class BaseTool(ABC):
    """Base class for all tools in the LangGraph project."""

    def __init__(self, name: str, description: str):
        """Initialize the base tool.

        Args:
            name: The name of the tool
            description: A description of what the tool does
        """
        self.name = name
        self.description = description

    @abstractmethod
    async def execute(self, **kwargs: Any) -> Dict[str, Any]:
        """Execute the tool with the given parameters.

        Args:
            **kwargs: Tool parameters

        Returns:
            Dictionary containing the tool result
        """
        pass

    def get_schema(self) -> Dict[str, Any]:
        """Get the tool schema for documentation."""
        return {
            "name": self.name,
            "description": self.description,
        }
