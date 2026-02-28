# mlcgo_mcp — Go MCP Server Hub

A collection of **Model Context Protocol (MCP)** servers written in Go. This repository serves as a central hub for specialized tools designed for AI agents.

Copyright (c) 2026 Michael Lechner. All rights reserved.
Licensed under the MIT License.

> 🇩🇪 [Deutsche Version](README.de.md)

---

## 📦 Included Servers

| Server | Description | Integration |
|--------|-------------|-------------|
| `d2mcp` | Renders D2 diagrams | Saves to [mlcartifact](https://github.com/hmsoft0815/mlcartifact) |
| `memory-server` | Knowledge graph & context | Persistent AI memory |
| `openstreetmap_mcp` | Real-world map data | OpenStreetMap integration |
| `task-manager` | AI-driven task management | Project orchestration |

---

### One-Line Installation (Linux/macOS)

The fastest way to install all hub servers:

```bash
curl -sfL https://raw.githubusercontent.com/hmsoft0815/mlcgo_mcp/main/scripts/install.sh | sh
```

### Pre-built Binaries & Linux Packages

Download the latest version as a **ZIP/TAR**, or install via **.deb** or **.rpm** (for Ubuntu/Debian/Fedora/openSUSE) from the **[GitHub Releases](https://github.com/hmsoft0815/mlcgo_mcp/releases)** page.

---

## Claude Desktop Integration

To use these servers in Claude Desktop, add them to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "d2mcp": {
      "command": "d2mcp"
    },
    "memory-server": {
      "command": "memory-server"
    },
    "openstreetmap": {
      "command": "openstreetmap_mcp"
    },
    "task-manager": {
      "command": "task-manager"
    }
  }
}
```
*Note: If the binaries are not in your PATH, provide the absolute path to the command.*

---

---

## 🏗️ Development

This repository uses a **Go Workspace** (`go.work`) to manage multiple independent modules. This allows for seamless cross-module testing and shared development.

```bash
# Tidy all modules in the hub
task tidy
```

---

## 📜 License
MIT License — Copyright (c) 2026 [Michael Lechner](https://github.com/hmsoft0815)
