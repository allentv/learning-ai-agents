Source code lives in `src/langgraph_project/`. Never place application code outside `src/`.
Tests live in `tests/`. Test files must be named `test_*.py`.
API code lives in `api/`. It has its own `requirements.txt`.
New tools go in `src/langgraph_project/tools/`. New agents go in `src/langgraph_project/agents/`.
Always export new public symbols in the relevant `__init__.py`.
