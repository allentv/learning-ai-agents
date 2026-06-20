Always add tests for new functions and classes.
Use `async def test_*()` for async tests — pytest-asyncio auto-detects them.
Mock external HTTP calls with `pytest-httpx`. Never make real network requests in tests.
Fix code to pass tests — do not modify test assertions unless explicitly asked.
Run `mise run v` after all changes to confirm everything passes.
