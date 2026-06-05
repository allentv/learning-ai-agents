#!/usr/bin/env python3
"""
Agent Builder Utility Script

This script provides utilities for building and managing Python agents
according to the project's agentic coding standards.
"""

import argparse
import sys
from pathlib import Path
from typing import List, Optional

# Add project root to path
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))


def create_agent_template(agent_name: str, description: str, tools: List[str]) -> str:
    """Create a new agent template based on project standards."""
    template = f'''"""
{agent_name} Agent

{description}
"""

from typing import TypedDict, Annotated
from langgraph.graph import StateGraph, END
from langchain_core.messages import HumanMessage, AIMessage

from .tools import {', '.join(tools)}


class AgentState(TypedDict):
    """State for the {agent_name} agent."""
    messages: Annotated[list, lambda x, y: x + y]
    context: dict


def {agent_name.lower()}_agent(state: AgentState) -> AgentState:
    """Main agent logic for {agent_name}.
    
    Args:
        state: Current agent state containing messages and context
        
    Returns:
        Updated agent state
    """
    # TODO: Implement agent logic
    return state


def build_{agent_name.lower()}_workflow():
    """Build the workflow for {agent_name} agent."""
    workflow = StateGraph(AgentState)
    workflow.add_node("agent", {agent_name.lower()}_agent)
    workflow.add_edge("agent", END)
    workflow.set_entry_point("agent")
    return workflow.compile()


if __name__ == "__main__":
    app = build_{agent_name.lower()}_workflow()
    # Run the agent
'''
    return template


def create_tool_template(tool_name: str, description: str, input_schema: dict) -> str:
    """Create a new tool template based on project standards."""
    schema_lines = []
    for param, param_desc in input_schema.items():
        schema_lines.append(f'    {param}: str = Field(description="{param_desc}")')

    schema_str = '\n'.join(schema_lines) if schema_lines else '    pass'

    template = f'''"""
{tool_name} Tool

{description}
"""

from langchain_core.tools import tool
from pydantic import BaseModel, Field


class {tool_name}Input(BaseModel):
    """Input schema for {tool_name} tool."""
{schema_str}


@tool(args_schema={tool_name}Input)
def {tool_name.lower()}({', '.join(input_schema.keys())}: str) -> str:
    """{description}.
    
    Args:
{chr(10).join([f'        {k}: {v}' for k, v in input_schema.items()])}
        
    Returns:
        Result of the tool execution
    """
    # TODO: Implement tool logic
    return f"Processed: {{', '.join(input_schema.keys())}}"
'''
    return template


def main():
    parser = argparse.ArgumentParser(description="Agent Builder Utility")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # Create agent command
    agent_parser = subparsers.add_parser("create-agent", help="Create a new agent template")
    agent_parser.add_argument("name", help="Name of the agent")
    agent_parser.add_argument("--description", help="Description of the agent", default="")
    agent_parser.add_argument("--tools", nargs="+", help="Tools to include", default=[])

    # Create tool command
    tool_parser = subparsers.add_parser("create-tool", help="Create a new tool template")
    tool_parser.add_argument("name", help="Name of the tool")
    tool_parser.add_argument("--description", help="Description of the tool", default="")
    tool_parser.add_argument("--params", nargs="+", help="Tool parameters (format: name:description)", default=[])

    args = parser.parse_args()

    if args.command == "create-agent":
        template = create_agent_template(args.name, args.description, args.tools)
        output_path = project_root / "agents" / f"{args.name.lower()}.py"
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(template)
        print(f"Created agent template at: {output_path}")

    elif args.command == "create-tool":
        params = {}
        for param in args.params:
            if ":" in param:
                name, desc = param.split(":", 1)
                params[name] = desc
            else:
                params[param] = f"Parameter {param}"

        template = create_tool_template(args.name, args.description, params)
        output_path = project_root / "agents" / "tools" / f"{args.name.lower()}.py"
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(template)
        print(f"Created tool template at: {output_path}")

    else:
        parser.print_help()


if __name__ == "__main__":
    main()
