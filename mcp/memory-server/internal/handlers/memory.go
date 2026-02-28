package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type MemoryHandler struct {
	db *sql.DB
}

func NewMemoryHandler(dbPath string) (*MemoryHandler, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS entities (
		name TEXT PRIMARY KEY,
		type TEXT
	);
	CREATE TABLE IF NOT EXISTS observations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entity_name TEXT,
		content TEXT,
		FOREIGN KEY(entity_name) REFERENCES entities(name) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS relations (
		from_name TEXT,
		to_name TEXT,
		type TEXT,
		PRIMARY KEY(from_name, to_name, type),
		FOREIGN KEY(from_name) REFERENCES entities(name) ON DELETE CASCADE,
		FOREIGN KEY(to_name) REFERENCES entities(name) ON DELETE CASCADE
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &MemoryHandler{db: db}, nil
}

// getField returns a string value from a map, checking both CamelCase and snake_case
func getField(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if val, ok := m[k].(string); ok {
			return val
		}
	}
	return ""
}

func (h *MemoryHandler) CreateEntities(args map[string]interface{}) (interface{}, error) {
	entities, ok := args["entities"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments: entities array missing")
	}

	for _, e := range entities {
		ent, ok := e.(map[string]interface{})
		if !ok { continue }
		
		name := getField(ent, "name", "entity_name")
		entType := getField(ent, "entityType", "entity_type", "type")
		
		if name == "" { continue }
		if entType == "" { entType = "unknown" }
		
		h.db.Exec("INSERT OR IGNORE INTO entities (name, type) VALUES (?, ?)", name, entType)

		// Handle observations inside entity
		if obs, ok := ent["observations"].([]interface{}); ok {
			for _, o := range obs {
				if str, ok := o.(string); ok {
					h.db.Exec("INSERT INTO observations (entity_name, content) VALUES (?, ?)", name, str)
				}
			}
		}
	}

	return "Entities processed", nil
}

func (h *MemoryHandler) CreateRelations(args map[string]interface{}) (interface{}, error) {
	relations, ok := args["relations"].([]interface{})
	if !ok { return "No relations provided", nil }

	for _, r := range relations {
		rel, ok := r.(map[string]interface{})
		if !ok { continue }
		
		from := getField(rel, "from", "source", "from_name")
		to := getField(rel, "to", "target", "to_name")
		relType := getField(rel, "relationType", "relation_type", "type")

		if from != "" && to != "" {
			if relType == "" { relType = "related_to" }
			h.db.Exec("INSERT OR IGNORE INTO entities (name, type) VALUES (?, 'unknown')", from)
			h.db.Exec("INSERT OR IGNORE INTO entities (name, type) VALUES (?, 'unknown')", to)
			h.db.Exec("INSERT OR IGNORE INTO relations (from_name, to_name, type) VALUES (?, ?, ?)", from, to, relType)
		}
	}

	return "Relations processed", nil
}

func (h *MemoryHandler) AddObservations(args map[string]interface{}) (interface{}, error) {
	observations, ok := args["observations"].([]interface{})
	if !ok { return "No observations provided", nil }

	for _, o := range observations {
		obs, ok := o.(map[string]interface{})
		if !ok { continue }
		
		name := getField(obs, "entityName", "entity_name", "name")
		contents, _ := obs["contents"].([]interface{})

		if name != "" {
			h.db.Exec("INSERT OR IGNORE INTO entities (name, type) VALUES (?, 'unknown')", name)
			for _, c := range contents {
				if str, ok := c.(string); ok {
					h.db.Exec("INSERT INTO observations (entity_name, content) VALUES (?, ?)", name, str)
				}
			}
		}
	}

	return "Observations processed", nil
}

func (h *MemoryHandler) SearchNodes(args map[string]interface{}) (interface{}, error) {
	query, _ := args["query"].(string)
	if query == "" { return "Empty query", nil }

	rows, err := h.db.Query(`
		SELECT DISTINCT e.name, e.type 
		FROM entities e
		LEFT JOIN observations o ON e.name = o.entity_name
		WHERE e.name LIKE ? OR o.content LIKE ? OR e.type LIKE ?`,
		"%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil { return nil, err }
	defer rows.Close()

	var results []string
	for rows.Next() {
		var name, entType string
		rows.Scan(&name, &entType)
		results = append(results, fmt.Sprintf("- %s (%s)", name, entType))
	}

	if len(results) == 0 { return "No matches found", nil }
	return "Search results:\n" + strings.Join(results, "\n"), nil
}

func (h *MemoryHandler) ReadGraph(args map[string]interface{}) (interface{}, error) {
	entRows, _ := h.db.Query("SELECT name, type FROM entities")
	defer entRows.Close()

	var sb strings.Builder
	sb.WriteString("KNOWLEDGE GRAPH CONTENT:\n")
	
	count := 0
	for entRows.Next() {
		count++
		var name, entType string
		entRows.Scan(&name, &entType)
		sb.WriteString(fmt.Sprintf("\nEntity: %s [%s]\n", name, entType))
		
		obsRows, _ := h.db.Query("SELECT content FROM observations WHERE entity_name = ?", name)
		for obsRows.Next() {
			var content string
			obsRows.Scan(&content)
			sb.WriteString(fmt.Sprintf("  - Observation: %s\n", content))
		}
		obsRows.Close()
	}

	if count == 0 { return "Memory is currently empty.", nil }
	return sb.String(), nil
}

func (h *MemoryHandler) DeleteEntities(args map[string]interface{}) (interface{}, error) {
	names, _ := args["entityNames"].([]interface{})
	for _, n := range names {
		if str, ok := n.(string); ok {
			h.db.Exec("DELETE FROM entities WHERE name = ?", str)
		}
	}
	return "Deleted", nil
}
