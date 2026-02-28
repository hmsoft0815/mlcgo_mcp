# d2mcp — D2 Diagram MCP Server

A specialized **Model Context Protocol (MCP)** server for rendering and manipulating [D2 diagrams](https://d2lang.com/).

Copyright (c) 2026 Michael Lechner. All rights reserved.
Licensed under the MIT License.

---

## Why D2 for AI Agents?

D2 is a modern diagram scripting language that turns text into diagrams. Unlike pixel-based formats, D2 allows agents to:
- **Build diagrams incrementally**: Use the "Oracle API" to add shapes and connections one by one.
- **Maintain Context**: The text-based nature means the agent can "read" the diagram's state.
- **High Portability**: Renders high-quality SVG that works everywhere.

---

## Features

- **Exclusive SVG Rendering**: Optimized for weight and scalability. PNG and PDF exports have been removed to keep the server lean.
- **Oracle API**: Incremental editing (create, set, delete, move, rename) without re-rendering the whole source.
- **[Optional] mlcartifact Integration**: If the [mlcartifact service](https://github.com/hmsoft0815/mlcartifact) is running, `d2mcp` automatically saves exports as persistent artifacts and returns a reference tag.
- **20+ Themes**: Support for all native D2 themes.

---

## Tools

### Core Tools
- `d2_create`: Initialize a new diagram session (can be empty or with initial content).
- `d2_export`: Render the current session to an SVG. If `mlcartifact` is active, it saves the result as a file.
- `render_artifact`: Reads a D2 source artifact, renders it to SVG, and saves it as a new artifact.

### Oracle API (Incremental)
- `d2_oracle_create`: Add a shape or connection.
- `d2_oracle_set`: Modify attributes (colors, labels, shapes).
- `d2_oracle_delete`: Remove elements.
- `d2_oracle_move`: Reorganize hierarchy.
- `d2_oracle_rename`: Change keys.
- `d2_oracle_serialize`: Get the full D2 source text.

---

## Installation & Usage

Part of the [mlcgo_mcp](https://github.com/hmsoft0815/mlcgo_mcp) hub.

```bash
# Build
go build -o d2mcp ./mcp/d2mcp

# Run (STDIO for Claude Desktop)
./d2mcp -transport=stdio
```

---

## License
MIT License — Copyright (c) 2026 [Michael Lechner](https://github.com/hmsoft0815)
