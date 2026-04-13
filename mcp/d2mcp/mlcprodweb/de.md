# D2-MCP Server

Bringen Sie professionelle Architektur-Diagramme in Ihr LLM. Der D2-MCP-Server ermöglicht es KI-Modellen, ansprechende Architekturdiagramme mit der deklarativen Sprache D2 zu entwerfen, zu generieren und zu exportieren.

## Warum D2-MCP?

Die Beschreibung komplexer Architekturen in Textform ist schwierig. D2-MCP erlaubt es Ihrem KI-Assistenten, "in Diagrammen zu denken" und visuelle Darstellungen von Software-Stacks, Netzwerktopologien und Workflows mühelos zu erstellen. Dank der integrierten Unterstützung für mehrere Layout-Engines sehen Ihre Diagramme immer professionell und organisiert aus.

## Kernfunktionen

- **Automatisierte Diagramm-Erstellung**: Übersetzen Sie High-Level-Architekturanforderungen in sauberen D2-Code und gerenderte Bilder (SVG/PNG).
- **Flexible Layouts**: Wechseln Sie zwischen Engines wie **Dagre**, **ELK**, **D2** und **Tala**, um die perfekte visuelle Darstellung für Ihr Design zu finden.
- **Oracle Design Sessions**: Enthält einen spezialisierten Prompt, der dem LLM hilft, in Design-Sitzungen als Experten-Softwarearchitekt zu agieren.
- **Produktionsbereit**: Unterstützt moderne Transportprotokolle (Stdio, SSE, Streamable HTTP), was die Integration in Desktop- und webbasierte KI-Tools vereinfacht.

## Einstieg

Verwenden Sie den Server lokal via **Stdio** für den privaten Gebrauch, oder setzen Sie einen **SSE-Server** für teamweite Zusammenarbeit auf.

## Schnellstart

### Claude Desktop
Fügen Sie Folgendes zu Ihrer `claude_desktop_config.json` hinzu:

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
Fügen Sie den Server zu Ihrer `~/.gemini/settings.json` hinzu:

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
Ein neues Profil hinzufügen:

```bash
mcp-tester profile add d2 -c "d2mcp"
```
