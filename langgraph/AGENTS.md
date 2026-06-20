# AGENTS.md

Project-specific context and conventions for AI coding assistants working on this codebase.

## Project Overview

This is a Python project using **LangGraph** and **Pydantic AI** for building stateful, multi-agent workflows with tool support. It includes a FastAPI API service for RAG-ready knowledge base search.

## Tech Stack

- **Python** 3.14+
- **LangGraph** — stateful, multi-agent application framework
- **Pydantic AI** — type-safe AI agent development with native tool support
- **FastAPI** — lightweight API service
- **Pydantic** — data validation and settings management
- **OpenAI SDK** — LLM client (used with local llama.cpp)
- **httpx** — async HTTP client for API tool calls
- **structlog** — structured logging
- **Docker Compose** — local development with llama.cpp, API, and app services

## Project Structure

```text
langgraph/
├── src/langgraph_project/          # Main application package
│   ├── main.py                     # Entry point
│   ├── agents/                     # Agent definitions
│   ├── tools/                      # Tool implementations (pydantic-ai compatible)
│   └── utils/                      # Config, logging utilities
├── api/                            # FastAPI API service
│   ├── app.py                      # FastAPI application
│   └── routes.py                   # API endpoints
├── tests/                          # Test suite
├── docs/                           # Documentation
├── langgraph.json                  # LangGraph Studio config
└── pyproject.toml                  # Project metadata and tool config
```

## Code Style & Conventions

- **Formatter**: Black (line length 120)
- **Import sorting**: isort (black profile, line length 120)
- **Linting**: ruff (rules: E, F, I, W)
- **Type checking**: mypy strict mode — all functions must have type annotations
- **Async**: Prefer async/await throughout. Use `httpx.AsyncClient`, not `requests`.
- **Logging**: Use `structlog` for all logging. Do not use `print()` or stdlib `logging`.
- **Configuration**: Use Pydantic Settings from `utils/config.py`. Read env vars via `pydantic-settings`.

## Common Commands

All tasks run via `mise run <task>`. Use these shorthand aliases when invoking commands:

| Alias | Full Command | Description |
|-------|-------------|-------------|
| `v` | `mise run verify` | **Run all checks** (format-check + lint + type-check + test) |
| `f` | `mise run format` | Format code (black + isort) |
| `fc` | `mise run format-check` | Check formatting (no changes) |
| `l` | `mise run lint` | Lint with ruff |
| `lf` | `mise run lint-fix` | Auto-fix lint issues |
| `t` | `mise run test` | Run pytest suite |
| `tc` | `mise run type-check` | mypy strict |
| `r` | `mise run run` | Run the application |

**After any code change, run `mise run v` (verify) before committing.** This runs format-check, lint, type-check-strict, and test in sequence.

## Testing

- Framework: **pytest** with **pytest-asyncio** (auto mode)
- Test files: `tests/test_*.py`
- Async tests are automatically detected — use `async def test_*()`
- Mock HTTP calls with **pytest-httpx**
- Run `mise run verify` before committing

## Adding New Tools

Tools are pydantic-ai compatible functions defined in `src/langgraph_project/tools/`. To add a new tool:

1. Create or extend a module in `src/langgraph_project/tools/`
2. Decorate functions with `@tool` from pydantic-ai
3. Use type annotations and Pydantic models for inputs/outputs
4. Add corresponding tests in `tests/`
5. Register the tool in the relevant agent in `src/langgraph_project/agents/`

## Adding New Agents

1. Create a new file in `src/langgraph_project/agents/`
2. Define the agent using LangGraph's graph primitives
3. Register tools and configure the LLM via `utils/config.py`
4. Add tests in `tests/`

## Skills (Reusable Instructions)

Reusable task-specific guides are in `.agents/skills/`. Reference these when performing the corresponding task:

| Skill | File | When to use |
|-------|------|-------------|
| Add a new tool | `.agents/skills/add-new-tool.md` | Creating or extending agent tools |
| Add a new agent | `.agents/skills/add-new-agent.md` | Building a new LangGraph agent |
| Add an API endpoint | `.agents/skills/add-api-endpoint.md` | Adding FastAPI routes |
| Fix type errors | `.agents/skills/fix-type-errors.md` | Resolving mypy strict errors |

## Commands (Reusable Prompts)

Actionable prompt templates in `.agents/commands/`. Invoke these for common workflows:

| Command | File | Description |
|---------|------|-------------|
| `review` | `.agents/commands/review.md` | Review code for correctness, types, style, security |
| `explain` | `.agents/commands/explain.md` | Explain code purpose, architecture, and data flow |
| `scaffold` | `.agents/commands/scaffold.md` | Generate a new tool, agent, or API endpoint |
| `fix` | `.agents/commands/fix.md` | Fix all failing format/lint/type/test checks |

## Rules (Always-Active Constraints)

Short, enforceable directives in `.agents/rules/`. These apply to every code change:

| Rule | File | Key constraints |
|------|------|----------------|
| Code style | `.agents/rules/code-style.md` | structlog not print, httpx not requests, type annotations required |
| Project structure | `.agents/rules/project-structure.md` | Files go in correct dirs, exports in `__init__.py` |
| Testing | `.agents/rules/testing.md` | Always add tests, mock HTTP, fix code not assertions |

## Key Patterns

- **State management**: Use LangGraph's `StateGraph` with typed state objects
- **Tool integration**: Use pydantic-ai's `@tool` decorator for type-safe tool definitions
- **Configuration**: Environment variables are the source of truth; use `utils/config.py` for typed access
- **Error handling**: Use structured logging with context; raise specific exceptions
- **No hardcoded values**: Secrets and environment-specific values go in `.env` (see `.env.example`)
