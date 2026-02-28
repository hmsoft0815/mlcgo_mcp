# Memory MCP Server

Provides a persistent knowledge graph based on SQLite to store and retrieve long-term facts and observations.

## Tools

### 1. `mlc_memorize`
Stores a new fact about an entity (person, project, concept).
- **Arguments**:
  - `entity` (required): Name of the entity.
  - `observation` (required): The fact to remember.
  - `category` (optional): Grouping (e.g., `Setting`, `Preference`).

### 2. `mlc_search_nodes`
Searches for stored memories using keywords.
- **Arguments**:
  - `query` (required): Search term.

### 3. `mlc_read_graph`
Returns a complete dump of the knowledge graph.

## Storage Location

Data is stored in `~/.local/share/mcp-proxy/memory.db`.

## Installation

Built as part of the main project:
```bash
task build
```
