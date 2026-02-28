# mlcgo_mcp — Go MCP Server Hub

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

## 🚀 Erste Schritte

### Voraussetzungen
- **Go 1.24+**
- **Task** (optional): `sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d`

### Alle Server bauen
```bash
# Mit Task (empfohlen)
task build-all

# Ohne Task
go build -o ./bin/d2mcp ./mcp/d2mcp
go build -o ./bin/memory-server ./mcp/memory-server
# ... etc
```
Die Binaries werden im Verzeichnis `./bin/` abgelegt.

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
