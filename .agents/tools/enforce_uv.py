#!/usr/bin/env python3
"""
UV Enforcement Utility Script

This script enforces the use of uv for all Python package installations
in the project and provides utilities for managing uv-based environments.
"""

import argparse
import subprocess
import sys
from pathlib import Path
from typing import List, Optional

# Add project root to path
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))


def check_uv_installed() -> bool:
    """Check if uv is installed and available."""
    try:
        result = subprocess.run(
            ["uv", "--version"],
            capture_output=True,
            text=True,
            check=True
        )
        print(f"✓ uv is installed: {result.stdout.strip()}")
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("✗ uv is not installed or not found in PATH")
        print("  Install uv with: curl -LsSf https://astral.sh/uv/install.sh | sh")
        return False


def create_uv_config() -> None:
    """Create uv configuration files."""
    uv_toml = project_root / "uv.toml"
    if not uv_toml.exists():
        uv_toml.write_text("""# UV Configuration
[tool.uv]
package-index = "pypi"
resolution = "highest"
python-preference = "only-managed"
""")
        print(f"✓ Created uv.toml at {uv_toml}")


def create_venv_with_uv(python_version: str = "3.10") -> None:
    """Create a virtual environment using uv."""
    venv_path = project_root / ".venv"
    
    if venv_path.exists():
        print(f"⚠️  Virtual environment already exists at {venv_path}")
        return

    try:
        subprocess.run(
            ["uv", "venv", str(venv_path), "--python", f"python{python_version}"],
            check=True,
            capture_output=True
        )
        print(f"✓ Created virtual environment at {venv_path}")
    except subprocess.CalledProcessError as e:
        print(f"✗ Failed to create virtual environment: {e}")
        if e.stderr:
            print(f"  Error: {e.stderr.decode()}")


def install_packages_with_uv(packages: List[str], dev: bool = False) -> None:
    """Install packages using uv."""
    if not check_uv_installed():
        return

    venv_path = project_root / ".venv"
    if not venv_path.exists():
        print("⚠️  No virtual environment found. Creating one...")
        create_venv_with_uv()

    cmd = ["uv", "pip", "install"]
    if dev:
        cmd.append("--dev")
    cmd.extend(packages)

    try:
        subprocess.run(cmd, check=True, capture_output=True)
        print(f"✓ Installed packages: {', '.join(packages)}")
    except subprocess.CalledProcessError as e:
        print(f"✗ Failed to install packages: {e}")
        if e.stderr:
            print(f"  Error: {e.stderr.decode()}")


def sync_dependencies() -> None:
    """Sync dependencies using uv."""
    if not check_uv_installed():
        return

    pyproject_path = project_root / "pyproject.toml"
    if not pyproject_path.exists():
        print("✗ pyproject.toml not found")
        return

    try:
        subprocess.run(
            ["uv", "sync"],
            check=True,
            capture_output=True,
            cwd=project_root
        )
        print("✓ Dependencies synced successfully")
    except subprocess.CalledProcessError as e:
        print(f"✗ Failed to sync dependencies: {e}")
        if e.stderr:
            print(f"  Error: {e.stderr.decode()}")


def check_uv_lock() -> None:
    """Check if uv.lock exists and is up to date."""
    lock_path = project_root / "uv.lock"
    if lock_path.exists():
        print("✓ uv.lock exists")
        # Check if it's up to date
        try:
            subprocess.run(
                ["uv", "lock", "--check"],
                check=True,
                capture_output=True,
                cwd=project_root
            )
            print("✓ uv.lock is up to date")
        except subprocess.CalledProcessError:
            print("⚠️  uv.lock may be out of date. Run 'uv lock' to update.")
    else:
        print("⚠️  uv.lock does not exist. Run 'uv lock' to create it.")


def main():
    parser = argparse.ArgumentParser(description="UV Enforcement Utility")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # Check command
    check_parser = subparsers.add_parser("check", help="Check if uv is installed")

    # Create venv command
    venv_parser = subparsers.add_parser("venv", help="Create virtual environment with uv")
    venv_parser.add_argument("--python", default="3.10", help="Python version to use")

    # Install command
    install_parser = subparsers.add_parser("install", help="Install packages with uv")
    install_parser.add_argument("packages", nargs="+", help="Packages to install")
    install_parser.add_argument("--dev", action="store_true", help="Install as dev dependencies")

    # Sync command
    subparsers.add_parser("sync", help="Sync dependencies with uv")

    # Lock command
    subparsers.add_parser("lock", help="Check uv.lock status")

    # Config command
    subparsers.add_parser("config", help="Create uv configuration files")

    args = parser.parse_args()

    if args.command == "check":
        check_uv_installed()

    elif args.command == "venv":
        create_venv_with_uv(args.python)

    elif args.command == "install":
        install_packages_with_uv(args.packages, args.dev)

    elif args.command == "sync":
        sync_dependencies()

    elif args.command == "lock":
        check_uv_lock()

    elif args.command == "config":
        create_uv_config()

    else:
        parser.print_help()


if __name__ == "__main__":
    main()
