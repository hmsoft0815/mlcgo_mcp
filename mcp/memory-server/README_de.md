# Memory MCP Server

Bietet einen persistenten Knowledge-Graph auf Basis von SQLite zum Speichern und Abrufen langfristiger Fakten und Beobachtungen.

## Tools

### 1. `mlc_memorize`
Speichert eine neue Information über eine Entität (Person, Projekt, Konzept).
- **Argumente**:
  - `entity` (erforderlich): Name der Entität.
  - `observation` (erforderlich): Die zu merkende Information.
  - `category` (optional): Gruppierung (z.B. `Setting`, `Preference`).

### 2. `mlc_search_nodes`
Sucht mithilfe von Stichworten nach gespeicherten Erinnerungen.
- **Argumente**:
  - `query` (erforderlich): Suchbegriff.

### 3. `mlc_read_graph`
Gibt einen vollständigen Dump des Knowledge-Graphs zurück.

## Speicherort

Die Daten werden in `~/.local/share/mcp-proxy/memory.db` gespeichert.

## Installation

Wird als Teil des Hauptprojekts gebaut:
```bash
task build
```
