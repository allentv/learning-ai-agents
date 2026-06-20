Always use `structlog` for logging. Never use `print()` or stdlib `logging`.
Always use `httpx.AsyncClient` for HTTP calls. Never use `requests`.
All functions must have complete type annotations (args and return type).
Use Pydantic models for structured data. Avoid raw dicts.
Line length is 120 characters (Black + isort).
