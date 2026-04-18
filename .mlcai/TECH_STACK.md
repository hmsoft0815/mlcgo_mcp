# 🛠 Tech Stack & Constraints: mlcgo_mcp — Go MCP Server Hub

## Kern-Versionen
- **Sprache:** Go 1.24.2
- **MCP SDK:** `github.com/mark3labs/mcp-go` (v0.44.1), `github.com/modelcontextprotocol/go-sdk` (v1.4.0)
- **Datenbank:** SQLite (`modernc.org/sqlite` v1.45.0) in `memory-server`

## Bibliotheken (Erlaubt/Fixiert)
- **Diagramme:** `oss.terrastruct.com/d2` (v0.7.0) in `d2mcp`
- **Artifact Store:** `github.com/hmsoft0815/mlcartifact` (lokal via `replace`)
- **JSON Handling:** `github.com/invopop/jsonschema`, `github.com/mailru/easyjson`
- **Testing:** `github.com/stretchr/testify`

## Einschränkungen (Constraints)
- **Go Workspace:** Alle Module müssen im `go.work` registriert sein, um ein einheitliches Build-Environment zu gewährleisten.
- **Cross-Platform:** Binaries werden via Goreleaser (`.goreleaser.yml`) für Linux, macOS (Intel/M1) und Windows gebaut.
- **Lokale Abhängigkeiten:** `mlcartifact` muss relativ zum Hub-Verzeichnis liegen (`../../../mlcartifact`), da die `go.mod` Dateien statische `replace`-Direktiven nutzen.
- **Binary-Struktur:** Jeder Server hat sein eigenes Sub-Modul in `mcp/` mit eigenem `go.mod`.

## Styling-Regeln
- **Go Style:** Standard Go Formatierung (`gofmt`).
- **Dateigröße:** MLC Standard (max 250 Zeilen pro Datei).
- **Dokumentation:** GoDoc Kommentare für alle exportierten Symbole sind Pflicht.
- **Error Handling:** Explizite Error-Behandlung und Wrapping nach Go-Best-Practices.

---

## 📋 Meta

- **Zuletzt aktualisiert:** 2026-04-18
- **Aktualisiert von:** gemini-2.0-flash-exp
- **Status:** Aktuell
