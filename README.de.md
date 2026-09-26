# mlcgo_mcp — Go MCP Server Hub

> **[mlcgo.eu](https://mlcgo.eu)** — Werkzeuge, Bibliotheken und Handbücher


Eine Sammlung von **Model Context Protocol (MCP)** Servern, geschrieben in Go. Dieses Repository dient als zentraler Einstiegspunkt für spezialisierte Tools, die für KI-Agenten entwickelt wurden.

Copyright (c) 2026 Michael Lechner. Alle Rechte vorbehalten.
Lizenziert unter der MIT-Lizenz.

> 🇬🇧 [English Version](README.md)

---

## 📦 Enthaltene Server

| Server | Beschreibung | Integration |
|--------|--------------|-------------|
| `d2mcp` | Erstellt D2-Diagramme | Nutzt [mlcartifact](https://github.com/hmsoft0815/mlcartifact) |
| `memory-server` | Knowledge-Graph & Kontext | Persistentes KI-Gedächtnis |
| `openstreetmap_mcp` | Reale Kartendaten | OpenStreetMap-Anbindung |
| `task-manager` | KI-gestützte Aufgabenverwaltung | Projektsteuerung |

---

### Komplett-Installation (Empfohlen)

Um **alle drei Komponenten** (Hub, Artifact Store und Wollmilchsau) in einem Rutsch zu installieren:

```bash
curl -sfL https://raw.githubusercontent.com/hmsoft0815/mlcgo_mcp/main/scripts/install-all.sh | sh
```

### Einzel-Installation (Linux/macOS)

Nur die Hub-Server installieren:

```bash
curl -sfL https://raw.githubusercontent.com/hmsoft0815/mlcgo_mcp/main/scripts/install.sh | sh
```

### Fertige Binaries & Linux-Pakete

Lade die aktuelle Version als **ZIP/TAR** herunter oder installiere sie bequem via **.deb** oder **.rpm** (für Ubuntu/Debian/Fedora/openSUSE) direkt von der **[GitHub Releases](https://github.com/hmsoft0815/mlcgo_mcp/releases)** Seite.

---

## Claude Desktop Integration

Um diese Server in Claude Desktop zu nutzen, füge sie zu deiner `claude_desktop_config.json` hinzu:

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
*Hinweis: Falls die Binaries nicht in deinem PATH liegen, gib bitte den absoluten Pfad zum Befehl an.*

---

---

## ��️ Entwicklung

Dieses Repository nutzt einen **Go Workspace** (`go.work`), um mehrere unabhängige Module zu verwalten. Dies ermöglicht komfortables Testen und Entwickeln über alle Server hinweg.

```bash
# Alle Module aufräumen
task tidy
```

---

## 📜 Lizenz
MIT-Lizenz — Copyright (c) 2026 [Michael Lechner](https://github.com/hmsoft0815)

<!-- mlcai-private -->
## Projektdokumentation (`.mlcai/`)

`.mlcai/` ist ein **privates Git-Submodul**: interne Planung, Backlog und Arbeitsnotizen, gepflegt mit dem MLC Doc Hub. Es ist nicht öffentlich zugänglich — **ohne** `--recurse-submodules` klonen; für den Build wird es nicht gebraucht. Links nach `.mlcai/` funktionieren nur mit Zugriff (`git submodule update --init .mlcai`).
