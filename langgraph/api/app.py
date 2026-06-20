"""FastAPI application for the LangGraph project API.

This module creates and configures the FastAPI application that serves
the search endpoint (and potentially other API endpoints in the future).
"""

from fastapi import FastAPI

from api.routes import router

app = FastAPI(
    title="LangGraph API",
    description="API service for the LangGraph agent project. "
    "Provides endpoints for knowledge base search (RAG-ready).",
    version="0.1.0",
)

app.include_router(router)


@app.get("/health")
async def health_check() -> dict[str, str]:
    """Health check endpoint."""
    return {"status": "healthy"}
