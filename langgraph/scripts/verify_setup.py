#!/usr/bin/env python3
"""Script to verify the LangGraph project setup."""

import importlib.util
import sys
from pathlib import Path


def check_dependencies():
    """Check if required dependencies are available."""
    dependencies = [
        "langgraph",
        "pydantic_ai",
        "pydantic",
        "openai",
        "structlog",
    ]

    missing = []
    for dep in dependencies:
        spec = importlib.util.find_spec(dep)
        if spec is None:
            missing.append(dep)

    if missing:
        print(f"❌ Missing dependencies: {', '.join(missing)}")
        return False
    else:
        print("✅ All dependencies are available")
        return True


def check_project_structure():
    """Check if the project structure is correct."""
    base_path = Path(__file__).parent.parent
    required_files = [
        "src/langgraph_project/__init__.py",
        "src/langgraph_project/main.py",
        "src/langgraph_project/agents/simple_agent.py",
        "src/langgraph_project/utils/config.py",
        "src/langgraph_project/utils/logging.py",
        "pyproject.toml",
        "README.md",
    ]

    missing = []
    for file_path in required_files:
        full_path = base_path / file_path
        if not full_path.exists():
            missing.append(file_path)

    if missing:
        print(f"❌ Missing files: {', '.join(missing)}")
        return False
    else:
        print("✅ Project structure is correct")
        return True


def check_imports():
    """Check if imports work correctly."""
    try:
        from langgraph_project.agents.simple_agent import SimpleAgent  # noqa: F401
        from langgraph_project.utils.config import settings  # noqa: F401
        from langgraph_project.utils.logging import setup_logging  # noqa: F401

        print("✅ All imports work correctly")
        return True
    except ImportError as e:
        print(f"❌ Import error: {e}")
        return False


def main():
    """Main verification function."""
    print("🔍 Verifying LangGraph project setup...")
    print()

    checks = [
        ("Dependencies", check_dependencies),
        ("Project Structure", check_project_structure),
        ("Imports", check_imports),
    ]

    results = []
    for name, check_func in checks:
        print(f"Checking {name}...")
        result = check_func()
        results.append(result)
        print()

    if all(results):
        print("🎉 All checks passed! Project setup is complete.")
        return 0
    else:
        print("❌ Some checks failed. Please review the errors above.")
        return 1


if __name__ == "__main__":
    sys.exit(main())
