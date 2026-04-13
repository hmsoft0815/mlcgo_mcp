# OpenStreetMap MCP Server

Statten Sie Ihr LLM mit realen Geodaten aus. Der OpenStreetMap MCP-Server stellt eine nahtlose Brücke zwischen modernen KI-Modellen und den umfassendsten offenen Kartendaten der Welt her.

## Warum OpenStreetMap MCP?

Geografie ist der Kontext für alles. Egal, ob Sie einen Assistenten bauen, der Schulen in der Nähe finden muss, ein Logistik-Tool für das Routing oder eine Standortanalyse-Engine – dieser Server liefert die präzisen, strukturierten Daten direkt aus OpenStreetMap.

## Kernfunktionen

- **Finden Sie alles, überall**: Nutzen Sie die **Overpass API**, um nach tausenden POI-Kategorien zu suchen – von E-Ladestationen bis hin zu lokalen Apotheken – mit hoher räumlicher Präzision.
- **Adressen in Insights umwandeln**: Verwenden Sie **Nominatim**, um menschenlesbare Adressen in maschinenlesbare Koordinaten zu transformieren (und umgekehrt).
- **Intelligente Navigation**: Integrieren Sie **OSRM-basiertes Routing**, um Reisezeiten und Distanzen für PKW, Fahrräder oder Fußgänger zu berechnen.
- **Fair Usage by Design**: Das integrierte Rate-Limiting stellt sicher, dass Ihre Integration die Infrastruktur der OpenStreetMap-Community respektiert.

## Einstieg

Der Server unterstützt sowohl **Stdio** für die lokale Desktop-Integration als auch **SSE** für Remote-Dienste und ermöglicht so den Einsatz in einer Vielzahl von Umgebungen.

## Schnellstart

### Claude Desktop
Fügen Sie Folgendes zu Ihrer `claude_desktop_config.json` hinzu:

```json
{
  "mcpServers": {
    "openstreetmap": {
      "command": "openstreetmap_mcp"
    }
  }
}
```

### Gemini-CLI
Fügen Sie den Server zu Ihrer `~/.gemini/settings.json` hinzu:

```json
{
  "mcpServers": {
    "openstreetmap": {
      "command": "openstreetmap_mcp"
    }
  }
}
```

### MCP-Tester
Ein neues Profil hinzufügen:

```bash
mcp-tester profile add osm -c "openstreetmap_mcp"
```
