"""API tools for the LangGraph project.

Provides tool functions compatible with pydantic-ai that allow the agent
to fetch information from external APIs. These are plain async functions
with type-annotated parameters — pydantic-ai automatically generates the
JSON schema for the LLM from the function signature.
"""

from typing import Any

import httpx
import structlog

from langgraph_project.utils.config import settings

logger = structlog.get_logger()


async def search_api(query: str) -> dict[str, Any]:
    """Search the knowledge base via the external API.

    Retrieves relevant documents or information for the given query.
    Use this tool whenever the user asks a question that might be answered
    by the knowledge base or when you need to look up factual information.

    Args:
        query: The search query to find relevant information.

    Returns:
        A dictionary containing the query, a list of results with id, title,
        snippet, and score fields, and metadata about the search.
    """
    logger.info("Tool: search_api called", query=query)

    try:
        async with httpx.AsyncClient(timeout=settings.api_timeout) as client:
            response = await client.get(
                f"{settings.api_base_url}/api/search",
                params={"q": query},
            )
            response.raise_for_status()
            data: dict[str, Any] = response.json()
            logger.info("Tool: search_api succeeded", result_count=len(data.get("results", [])))
            return data
    except httpx.TimeoutException:
        logger.error("Tool: search_api timed out", query=query)
        return {
            "error": "The search request timed out. Please try again.",
            "query": query,
            "results": [],
        }
    except httpx.HTTPStatusError as e:
        logger.error("Tool: search_api HTTP error", status=e.response.status_code, query=query)
        return {
            "error": f"API returned status {e.response.status_code}: {e.response.text}",
            "query": query,
            "results": [],
        }
    except httpx.RequestError as e:
        logger.error("Tool: search_api request error", error=str(e), query=query)
        return {
            "error": f"Failed to connect to the API: {str(e)}",
            "query": query,
            "results": [],
        }
