#!/usr/bin/env python3
"""Script to verify the LangGraph project setup."""

import importlib.util
import sys
from pathlib import Path

# Add src directory to Python path
src_path = Path(__file__).parent.parent / "src"
sys.path.insert(0, str(src_path))


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


def check_docker_compose_config():
    """Check if the configuration matches docker-compose.yaml."""
    from langgraph_project.utils.config import settings

    # Check if the model provider is set to llamacpp (as in docker-compose.yaml)
    if settings.model_provider != "llamacpp":
        print(f"❌ Model provider mismatch: Expected 'llamacpp', got '{settings.model_provider}'")
        return False

    # Check if llamacpp URL matches docker-compose.yaml
    expected_url = "http://llamacpp:12434/v1"
    if settings.llamacpp_url != expected_url:
        print(f"❌ llama.cpp URL mismatch: Expected '{expected_url}', got '{settings.llamacpp_url}'")
        return False

    # Check if llamacpp model matches docker-compose.yaml
    expected_model = "granite-4.0-h-micro-UD-Q4_K_XL.gguf"
    if settings.llamacpp_model != expected_model:
        print(f"❌ llama.cpp model mismatch: Expected '{expected_model}', got '{settings.llamacpp_model}'")
        return False

    # Check other settings from docker-compose.yaml
    if settings.log_level != "INFO":
        print(f"❌ Log level mismatch: Expected 'INFO', got '{settings.log_level}'")
        return False

    if settings.log_format != "json":
        print(f"❌ Log format mismatch: Expected 'json', got '{settings.log_format}'")
        return False

    if settings.app_name != "LangGraph Project":
        print(f"❌ App name mismatch: Expected 'LangGraph Project', got '{settings.app_name}'")
        return False

    if settings.app_version != "0.1.0":
        print(f"❌ App version mismatch: Expected '0.1.0', got '{settings.app_version}'")
        return False

    print("✅ Configuration matches docker-compose.yaml")
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
        ("Docker Compose Configuration", check_docker_compose_config),
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
