#!/usr/bin/env python3
"""
Agent Validation Utility Script

This script validates Python agents against the project's agentic coding standards.
"""

import ast
import sys
from pathlib import Path
from typing import List, Dict, Any

# Add project root to path
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))


class AgentValidator:
    """Validator for Python agents against coding standards."""

    def __init__(self, file_path: Path):
        self.file_path = file_path
        self.errors: List[str] = []
        self.warnings: List[str] = []

    def validate(self) -> bool:
        """Run all validation checks."""
        if not self.file_path.exists():
            self.errors.append(f"File not found: {self.file_path}")
            return False

        with open(self.file_path, 'r') as f:
            content = f.read()

        try:
            tree = ast.parse(content)
        except SyntaxError as e:
            self.errors.append(f"Syntax error: {e}")
            return False

        self._check_imports(tree)
        self._check_type_hints(tree)
        self._check_docstrings(tree)
        self._check_agent_patterns(tree)

        return len(self.errors) == 0

    def _check_imports(self, tree: ast.AST) -> None:
        """Check import organization and usage."""
        imports = []
        for node in ast.walk(tree):
            if isinstance(node, (ast.Import, ast.ImportFrom)):
                imports.append(node)

        # Check for wildcard imports
        for imp in imports:
            if isinstance(imp, ast.ImportFrom):
                if imp.names and imp.names[0].name == '*':
                    self.errors.append("Wildcard imports are not allowed")

    def _check_type_hints(self, tree: ast.AST) -> None:
        """Check for type hints in function definitions."""
        for node in ast.walk(tree):
            if isinstance(node, ast.FunctionDef):
                # Check if function has type hints
                has_return_annotation = node.returns is not None
                has_param_annotations = any(
                    arg.annotation is not None for arg in node.args.args
                )

                if not has_return_annotation and node.name != '__init__':
                    self.warnings.append(f"Function '{node.name}' missing return type annotation")

                if not has_param_annotations and node.args.args:
                    self.warnings.append(f"Function '{node.name}' missing parameter type annotations")

    def _check_docstrings(self, tree: ast.AST) -> None:
        """Check for docstrings in public functions and classes."""
        for node in ast.walk(tree):
            if isinstance(node, (ast.FunctionDef, ast.ClassDef)):
                if not node.name.startswith('_'):  # Only check public members
                    docstring = ast.get_docstring(node)
                    if not docstring:
                        self.warnings.append(f"Missing docstring for '{node.name}'")

    def _check_agent_patterns(self, tree: ast.AST) -> None:
        """Check for agent-specific patterns."""
        has_state_class = False
        has_agent_function = False

        for node in ast.walk(tree):
            if isinstance(node, ast.ClassDef):
                if 'State' in node.name:
                    has_state_class = True
            elif isinstance(node, ast.FunctionDef):
                if 'agent' in node.name.lower():
                    has_agent_function = True

        if not has_state_class:
            self.warnings.append("No agent state class found")
        if not has_agent_function:
            self.warnings.append("No agent function found")

    def get_report(self) -> str:
        """Generate validation report."""
        report = f"Validation Report for {self.file_path}\n"
        report += "=" * 50 + "\n"

        if self.errors:
            report += "ERRORS:\n"
            for error in self.errors:
                report += f"  ❌ {error}\n"
        else:
            report += "No errors found.\n"

        if self.warnings:
            report += "\nWARNINGS:\n"
            for warning in self.warnings:
                report += f"  ⚠️  {warning}\n"
        else:
            report += "No warnings found.\n"

        return report


def main():
    import argparse

    parser = argparse.ArgumentParser(description="Validate Python agents against coding standards")
    parser.add_argument("file_path", help="Path to the Python agent file to validate")
    parser.add_argument("--verbose", "-v", action="store_true", help="Show detailed output")

    args = parser.parse_args()

    validator = AgentValidator(Path(args.file_path))
    is_valid = validator.validate()

    print(validator.get_report())

    if not is_valid:
        sys.exit(1)


if __name__ == "__main__":
    main()
