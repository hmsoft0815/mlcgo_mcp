# 🤖 AI & Integration Context: mlcgo_mcp — Go MCP Server Hub

## 1. Identität & Zweck
- **Kernaufgabe:** Zentraler Hub für spezialisierte Model Context Protocol (MCP) Server in Go.
- **Technischer Stack:** siehe `TECH_STACK.md`.
- **Hoster/Infrastruktur:** Lokal (Binary-Installation) oder via Docker; primär zur Integration in Claude Desktop oder andere MCP-Clients.

## 2. Die \"Nachbarschaft\" (System-Kontext)
- **Upstream (Wovon hänge ich ab?):**
  - [mlcartifact](https://github.com/hmsoft0815/mlcartifact) -> Lokales Repository (wird via `replace` eingebunden) -> Dient als Speicher für generierte D2-Diagramme (`d2mcp`).
- **Downstream (Wer nutzt mich?):**
  - **Claude Desktop** -> Nutzt die Server via `claude_desktop_config.json`.
  - Beliebige MCP-fähige AI-Clients -> Konsumieren Tools und Ressourcen der Sub-Server.
- **Shared Resources:**
  - Gemeinsames Go-Workspace (`go.work`) zur Entwicklung und zum Testen aller enthaltenen Server.

## 3. Schnittstellen-Vertrag
- **Primäre API:** Model Context Protocol (MCP).
- **Transport:** Standard-I/O (STDIO) oder Server-Sent Events (SSE).
- **Auth-Mechanismus:** Keine (lokale Ausführung unter User-Berechtigungen des Host-Prozesses).
- **Wichtige Datenmodelle:**
  - MCP Tools, Resources und Prompts (spezifisch pro Server, siehe Sub-Projekte).
- **API-Doku-Link:** Siehe die jeweiligen READMEs in `mcp/`.

## 4. Leitplanken & Regeln
- **Naming:** Go-Konventionen (camelCase für lokale, PascalCase für exportierte Symbole).
- **Testing:** Unit-Tests in `*_test.go` Dateien innerhalb der Sub-Module; Integrationstests via MCP-Client-Simulation.
- **Sicherheit:** Lokale Ausführung; Zugriffsbeschränkungen durch den Host-Prozess (z.B. Claude Desktop).

## 5. Aktueller Fokus (Status)
- **Bekannte Probleme:** Keine kritischen Probleme bekannt.
- **Nächste Schritte:** Erweiterung um weitere MCP-Server (z.B. `markitdown-go`) und Optimierung der SSE-Transport-Unterstützung.

---

## 6. Weitere Dokumentation (außerhalb .mlcai/)
- [README.md](../README.md) - Hauptdokumentation und Installationsanleitung.
- [README.de.md](../README.de.md) - Deutsche Version der Hauptdokumentation.
- [mcp/README.md](../mcp/README.md) - Übersicht der enthaltenen Server.
- [mcp/d2mcp/README.md](../mcp/d2mcp/README.md) - Doku für d2mcp.
- [mcp/memory-server/README.md](../mcp/memory-server/README.md) - Doku für memory-server.
- [mcp/openstreetmap_mcp/README.md](../mcp/openstreetmap_mcp/README.md) - Doku für openstreetmap_mcp.
- [mcp/task-manager/README.md](../mcp/task-manager/README.md) - Doku für task-manager.

---

## 📋 Meta

- **Zuletzt aktualisiert:** 2026-04-18
- **Aktualisiert von:** gemini-2.0-flash-exp
- **Status:** Aktuell
