<!-- mlc-dochub:begin — auto-managed, do not edit between these markers -->
## MLC Doc Hub — mlcgo-mcp / openstreetmap-mcp

**You are working in sub-project** `openstreetmap-mcp` of project `mlcgo-mcp`.
Every `mcp__mlc-dochub__*` call for this directory must pass
`project_id="mlcgo-mcp"` **together with** `sub_id="openstreetmap-mcp"`.

- **Source directory** (this cwd): `/mnt/data2tb/mlcgo_mcp/mcp/openstreetmap_mcp`
- **Parent project source:** `/mnt/data2tb/mlcgo_mcp`

The parent project has its own `.mlcai/` at `../../.mlcai/` — do not
write there unless you mean the parent. Sub-specific docs live at
`../../.mlcai/openstreetmap-mcp/`.


### Existing docs

Read any of these directly (native Read is fine) to gather context before you work:

- `../../.mlcai/openstreetmap-mcp/PRODUCT.md` — Product marketing page pointer (submodule meta)
- `../../.mlcai/INTEGRATION.md` — How this project fits into the larger system
- `../../.mlcai/TECH_STACK.md` — Stack & dependencies

### Recommended, not yet created

- `../../.mlcai/openstreetmap-mcp/INTEGRATION.md` — How this project fits into the larger system
- `../../.mlcai/openstreetmap-mcp/TECH_STACK.md` — Stack & dependencies
- `../../.mlcai/openstreetmap-mcp/API_CONTRACT.md` — API contract / endpoints
- `../../.mlcai/openstreetmap-mcp/DECISION_LOG.md` — Architecture decisions

### The rules themselves

The `mlc-dochub` server sends its full operating rules on connect, and most
clients place them in the system prompt — Claude Code and crush do (verified).
**Antigravity/`agy` does not** (verified 29.08.2026: it connects the server and
exposes its tools, but drops the server-level instructions). If you cannot see
them, this section is your copy:

- **Never write a `.mlcai/` file with a native editor.** Every create / update /
  delete goes through the MCP tools. They stamp the `## 📋 Meta` footer, append
  to the activity log and guard against concurrent edits. Reading with a native
  Read is fine and usually cheaper.
- **Pass `author="<your-model-name>"` on every write.**
- **Use the project id `mlcgo-mcp`.** Never invent `mlcgo-mcp-mac`, and
  never "fix" the registered path to your own filesystem — the registry is shared
  across every host, and machine-specific paths corrupt it for everyone else. A
  server-side path that does not exist on your machine is expected.
- **Three docs, three jobs.** `WORKLOG.md` is working memory — what is happening
  now and the next 1–3 steps, written at the END of a session, replaced wholesale
  by `update_worklog`. `PLAN.md` holds multi-phase plans, ticked off section by
  section with `update_section`. `BACKLOG.md` holds tickets, via the
  `*_backlog_item` tools. A vague todo → BACKLOG; a feature you are about to
  build → PLAN; "where I am right now" → WORKLOG.
- **Read before you write.** `get_doc` returns `last_modified`; pass it as
  `base_modified` on the next write so a concurrent edit is caught.
- **Writes are committed and pushed for you.** No git needed on `.mlcai/`. If a
  result says `[committed locally — PUSH FAILED …]`, tell the user.
- **If the tools are missing entirely, do not write** — read, and point the user
  at `task install-all` in the mlcintegration checkout (https://github.com/mlc911/mlcintegration).

`mcp__mlc-dochub__get_project_context` returns this project's metadata and style
guide; call it once before creating new docs.
<!-- mlc-dochub:end -->

## 🧠 Codebase Memory & Intelligence Engine (`cbm`)

Dieses Projekt unterstützt den `codebase-memory` MCP-Server (`cbm-server`) für blitzschnelle strukturelle Codebase-Navigation (<50µs In-Memory Graph, Tree-Sitter ASTs, Call-Hierarchien und Qdrant-Vektorsuche).

### 🛠️ Empfohlener Agent-Workflow:
1. **Orientierung bei Session-Start:** Rufe `get_repo_map` auf, um die Paketstruktur, Interfaces und Monorepo-Subprojekte kompakt (<400 Tokens) zu erfassen.
2. **Falls Repository noch nicht indexiert ist:** Rufe einmalig `index_repository(path=".")` auf.
3. **Funktions- & Aufrufketten analysieren:** Nutze `get_callers`, `get_callees` oder `trace_call_path` anstelle manueller Datei-Suchen.
4. **Vor Refactorings / Änderungen:** Prüfe Abhängigkeiten mit `simulate_refactoring` oder `get_impact_radius`.
5. **Git Diff Impact:** Nutze `detect_changes`, um den Blast Radius ungespeicherter Änderungen zu analysieren.
