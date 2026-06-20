# Installation

## Prerequisites

- [mise](https://mise.jdx.dev/) — task runner and tool version manager
- [UV](https://docs.astral.sh/uv/) — fast Python package installer
- Python 3.14 (installed automatically by mise)

## Quick Start

1. **Install mise** (if not already installed):

   ```bash
   curl https://mise.run | sh
   ```

2. **Install UV** (if not already installed):

   ```bash
   curl -LsSf https://astral.sh/uv/install.sh | sh
   ```

3. **Set up the development environment**:

   ```bash
   mise run setup
   ```

   This installs dependencies and configures the project. Alternatively, you can run `uv sync --dev` directly.

4. **Set up environment variables**:

   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

## Available Tasks

This project uses [mise](https://mise.jdx.dev/) for task management. Tasks are defined as standalone scripts in `.mise/tasks/`. Run `mise tasks` to see all available tasks.

```bash
# List all tasks with descriptions
mise tasks

# Run a task
mise run <task-name>
