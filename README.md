# mlcgo_mcp — Go MCP Server Hub

A collection of **Model Context Protocol (MCP)** servers written in Go. This repository serves as a central hub for specialized tools designed for AI agents.

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

## 🚀 Getting Started

### Prerequisites
- **Go 1.24+**
- **Task** (optional): `sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d`

### Build All Servers
```bash
# Using Task (recommended)
task build-all

# Without Task
go build -o ./bin/d2mcp ./mcp/d2mcp
go build -o ./bin/memory-server ./mcp/memory-server
# ... etc
```
Binaries will be placed in the `./bin/` directory.

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
