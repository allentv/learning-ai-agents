# Skill: Add a New API Endpoint

Step-by-step guide for adding a new endpoint to the FastAPI API service.

## Steps

1. **Define the route** in `api/routes.py`:

```python
from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter()

class MyRequest(BaseModel):
    query: str

class MyResponse(BaseModel):
    result: str

@router.get("/api/my-endpoint", response_model=MyResponse)
async def my_endpoint(request: MyRequest) -> MyResponse:
    """Brief description of this endpoint."""
    # Implementation here
    return MyResponse(result="...")
```

2. **Include the router** in `api/app.py` if it's a new router, or add the endpoint to the existing router.

3. **Add tests** in `tests/test_api_<feature>.py`:
   - Use `httpx.AsyncClient` with the FastAPI `TestClient` or `pytest-httpx`
   - Test both success and error cases
   - Test input validation

4. **Update documentation** in `docs/api.md`.

5. **Run checks**: `mise run v`

## Conventions

- Use `GET` for read operations, `POST` for mutations
- Always define `response_model` for OpenAPI docs
- Use Pydantic models for request/response validation
- Keep routes thin — delegate business logic to service functions
- Use async handlers with `httpx.AsyncClient` for external calls
- Return appropriate HTTP status codes (200, 201, 400, 404, 500)
