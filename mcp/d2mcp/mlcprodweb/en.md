# D2-MCP Server

Bring professional architectural diagramming to your LLM. The D2-MCP server enables AI models to design, generate, and export beautiful architecture diagrams using the D2 declarative language.

## Why D2-MCP?

Describing complex architectures in text is hard. D2-MCP allows your AI assistant to "think in diagrams," creating visual representations of software stacks, network topologies, and workflows effortlessly. With built-in support for multiple layout engines, your diagrams will always look professional and organized.

## Key Capabilities

- **Automated Diagram Generation**: Translate high-level architectural requirements into clean D2 code and rendered images (SVG/PNG).
- **Flexible Layouts**: Switch between engines like **Dagre**, **ELK**, **D2**, and **Tala** to find the perfect visual representation for your design.
- **Oracle Design Sessions**: Includes a specialized prompt to help the LLM act as an expert software architect during design sessions.
- **Production Ready**: Supports modern transport protocols (Stdio, SSE, Streamable HTTP), making it easy to integrate into both desktop and web-based AI tools.

## Getting Started

Deploy locally using **Stdio** for private use, or set up an **SSE/Streamable HTTP** server for shared organizational workflows.

## Quick Setup

### Claude Desktop
Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "d2mcp": {
      "command": "d2mcp"
    }
  }
}
```

### Gemini-CLI
Add to your `~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "d2mcp": {
      "command": "d2mcp"
    }
  }
}
```

### MCP-Tester
Add a new profile:

```bash
mcp-tester profile add d2 -c "d2mcp"
```
