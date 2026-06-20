"""API routes for the LangGraph project.

Provides a simple GET endpoint designed as a placeholder for future RAG
( Retrieval-Augmented Generation) functionality.
"""

from typing import Any

import structlog
from fastapi import APIRouter, HTTPException, Query

logger = structlog.get_logger()

router = APIRouter(prefix="/api", tags=["api"])


@router.get("/search")
async def search(q: str = Query(..., description="Search query")) -> dict[str, Any]:
    """Search the knowledge base.

    This is a simple placeholder endpoint that returns mock results.
    It is designed to be replaced with a RAG pipeline in the future.

    Args:
        q: The search query string.

    Returns:
        A dictionary containing the query, results, and metadata.
    """
    logger.info("Search request received", query=q)

    # TODO: Replace with actual RAG pipeline
    # For now, return mock results based on the query
    results = [
        {
            "id": 1,
            "title": f"Result for '{q}'",
            "snippet": f"This is a placeholder result for the query: {q}. "
            "In the future, this will return relevant documents from a vector store.",
            "score": 0.95,
        },
        {
            "id": 2,
            "title": f"Additional context for '{q}'",
            "snippet": f"More information related to '{q}' will be provided "
            "once the RAG pipeline is integrated.",
            "score": 0.82,
        },
    ]

    return {
        "query": q,
        "results": results,
        "metadata": {
            "total_results": len(results),
            "source": "mock",
        },
    }
