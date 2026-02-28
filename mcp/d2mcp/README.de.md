# d2mcp — D2 Diagramm MCP-Server

**Hinweis: Dies ist eine modifizierte Version des ursprünglichen [d2mcp](https://github.com/i2y/d2mcp) Projekts von i2y.**  
Diese Version wurde von Michael Lechner angepasst und erweitert, um in das `mlcgo_mcp` Ökosystem zu passen, wobei der Fokus auf SVG-Ausgabe und die optionale Integration mit dem `mlcartifact` Storage-Service liegt.

Copyright (c) 2026 Michael Lechner. Alle Rechte vorbehalten.
Ursprüngliches Projekt Copyright (c) 2024 i2y.
Lizenziert unter der MIT-Lizenz.

---

## Warum D2 für KI-Agenten?

D2 ist eine moderne Skriptsprache für Diagramme, die Text in Grafiken verwandelt. Im Gegensatz zu pixelbasierten Formaten ermöglicht D2:
- **Inkrementeller Aufbau**: Nutzen Sie die "Oracle API", um Formen und Verbindungen Schritt für Schritt hinzuzufügen.
- **Kontext-Erhalt**: Da Diagramme textbasiert sind, kann der Agent den Status des Diagramms jederzeit "lesen".
- **Portabilität**: Erzeugt hochwertige SVGs, die überall funktionieren.

---

## Features

- **Exklusives SVG-Rendering**: Optimiert für Performance und Skalierbarkeit. PNG- und PDF-Exporte wurden entfernt, um den Server schlank zu halten.
- **Oracle API**: Inkrementelle Bearbeitung (Erstellen, Setzen, Löschen, Verschieben, Umbenennen) ohne das gesamte Diagramm neu rendern zu müssen.
- **[Optional] mlcartifact Integration**: Wenn der [mlcartifact Dienst](https://github.com/hmsoft0815/mlcartifact) läuft, speichert `d2mcp` Exporte automatisch als persistente Artefakte und gibt ein Referenz-Tag zurück.
- **20+ Themes**: Unterstützung für alle nativen D2-Themes.

---

## Tools

### Kern-Tools
- `d2_create`: Initialisiert eine neue Diagrammsitzung (leer oder mit Inhalt).
- `d2_export`: Rendert die aktuelle Sitzung als SVG. Falls `mlcartifact` aktiv ist, wird das Ergebnis als Datei gespeichert.
- `render_artifact`: Liest ein D2-Quell-Artefakt, rendert es zu SVG und speichert es als neues Artefakt.

### Oracle API (Inkrementell)
- `d2_oracle_create`: Form oder Verbindung hinzufügen.
- `d2_oracle_set`: Attribute ändern (Farben, Labels, Formen).
- `d2_oracle_delete`: Elemente entfernen.
- `d2_oracle_move`: Hierarchie reorganisieren.
- `d2_oracle_rename`: Schlüssel ändern.
- `d2_oracle_serialize`: Den vollständigen D2-Quelltext abrufen.

---

## Installation & Nutzung

Teil des [mlcgo_mcp](https://github.com/hmsoft0815/mlcgo_mcp) Hubs.

```bash
# Bauen
go build -o d2mcp ./mcp/d2mcp

# Starten (STDIO für Claude Desktop)
./d2mcp -transport=stdio
```

---

## Lizenz
MIT-Lizenz — Copyright (c) 2026 [Michael Lechner](https://github.com/hmsoft0815)
