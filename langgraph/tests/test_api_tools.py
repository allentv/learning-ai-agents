"""Tests for the API tools module."""

import httpx
import pytest

from langgraph_project.tools.api_tools import search_api


class TestSearchApi:
    """Test cases for the search_api tool."""

    @pytest.mark.asyncio
    async def test_search_api_success(self, httpx_mock: httpx.MockTransport) -> None:
        """Test successful API search call."""
        mock_response = {
            "query": "test query",
            "results": [
                {"id": 1, "title": "Result 1", "snippet": "Snippet 1", "score": 0.95},
            ],
            "metadata": {"total_results": 1, "source": "mock"},
        }

        def handler(request: httpx.Request) -> httpx.Response:
            assert "/api/search" in str(request.url)
            assert "q=test+query" in str(request.url) or "q=test%20query" in str(request.url)
            return httpx.Response(200, json=mock_response)

        transport = httpx.MockTransport(handler)
        # Patch the AsyncClient used inside search_api
        import langgraph_project.tools.api_tools as api_tools
        from unittest.mock import AsyncMock, patch

        mock_client = httpx.AsyncClient(transport=transport, base_url="http://testserver:8000")

        with patch.object(api_tools.settings, "api_base_url", "http://testserver:8000"):
            with patch("langgraph_project.tools.api_tools.httpx.AsyncClient") as MockClient:
                MockClient.return_value.__aenter__ = AsyncMock(return_value=mock_client)
                MockClient.return_value.__aexit__ = AsyncMock(return_value=False)

                result = await search_api("test query")

        assert result["query"] == "test query"
        assert len(result["results"]) == 1
        assert result["results"][0]["title"] == "Result 1"

    @pytest.mark.asyncio
    async def test_search_api_timeout(self) -> None:
        """Test API search timeout handling."""
        from unittest.mock import AsyncMock, patch

        def raise_timeout(*args, **kwargs):  # type: ignore[no-untyped-def]
            raise httpx.TimeoutException("Connection timed out")

        mock_client = AsyncMock()
        mock_client.__aenter__ = AsyncMock(return_value=mock_client)
        mock_client.__aexit__ = AsyncMock(return_value=False)
        mock_client.get = AsyncMock(side_effect=httpx.TimeoutException("Connection timed out"))

        import langgraph_project.tools.api_tools as api_tools

        with patch.object(api_tools.settings, "api_base_url", "http://testserver:8000"):
            with patch("langgraph_project.tools.api_tools.httpx.AsyncClient") as MockClient:
                MockClient.return_value.__aenter__ = AsyncMock(return_value=mock_client)
                MockClient.return_value.__aexit__ = AsyncMock(return_value=False)

                result = await search_api("test query")

        assert "error" in result
        assert "timed out" in result["error"]
        assert result["results"] == []

    @pytest.mark.asyncio
    async def test_search_api_http_error(self) -> None:
        """Test API search HTTP error handling."""
        from unittest.mock import AsyncMock, patch

        mock_response = httpx.Response(
            status_code=500,
            request=httpx.Request("GET", "http://testserver:8000/api/search"),
            text="Internal Server Error",
        )

        mock_client = AsyncMock()
        mock_client.__aenter__ = AsyncMock(return_value=mock_client)
        mock_client.__aexit__ = AsyncMock(return_value=False)
        mock_client.get = AsyncMock(side_effect=httpx.HTTPStatusError(
            "Server Error",
            request=httpx.Request("GET", "http://testserver:8000/api/search"),
            response=mock_response,
        ))

        import langgraph_project.tools.api_tools as api_tools

        with patch.object(api_tools.settings, "api_base_url", "http://testserver:8000"):
            with patch("langgraph_project.tools.api_tools.httpx.AsyncClient") as MockClient:
                MockClient.return_value.__aenter__ = AsyncMock(return_value=mock_client)
                MockClient.return_value.__aexit__ = AsyncMock(return_value=False)

                result = await search_api("test query")

        assert "error" in result
        assert "500" in result["error"]
        assert result["results"] == []

    @pytest.mark.asyncio
    async def test_search_api_connection_error(self) -> None:
        """Test API search connection error handling."""
        from unittest.mock import AsyncMock, patch

        mock_client = AsyncMock()
        mock_client.__aenter__ = AsyncMock(return_value=mock_client)
        mock_client.__aexit__ = AsyncMock(return_value=False)
        mock_client.get = AsyncMock(side_effect=httpx.ConnectError("Connection refused"))

        import langgraph_project.tools.api_tools as api_tools

        with patch.object(api_tools.settings, "api_base_url", "http://testserver:8000"):
            with patch("langgraph_project.tools.api_tools.httpx.AsyncClient") as MockClient:
                MockClient.return_value.__aenter__ = AsyncMock(return_value=mock_client)
                MockClient.return_value.__aexit__ = AsyncMock(return_value=False)

                result = await search_api("test query")

        assert "error" in result
        assert "Failed to connect" in result["error"]
        assert result["results"] == []
