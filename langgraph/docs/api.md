# API Service

The project includes a lightweight FastAPI service designed as a placeholder for future RAG (Retrieval-Augmented Generation) functionality.

## Running the API Service Locally

```bash
mise run api
```

This starts the API server with hot-reload on port `10000`.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/search?q=<query>` | Search the knowledge base |
| `GET` | `/health` | Health check |

## Example Response

```json
{
  "query": "hello",
  "results": [
    {
      "id": 1,
      "title": "Result for 'hello'",
      "snippet": "This is a placeholder result...",
      "score": 0.95
    }
  ],
  "metadata": {
    "total_results": 1,
    "source": "mock"
  }
}
```

## Testing the API

```bash
mise run api-test
```

This runs health checks and a sample search query against the locally running API.
