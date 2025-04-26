package database

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// QueryManager loads and provides access to SQL queries
type QueryManager struct {
	queries map[string]string
}

// NewQueryManager loads all SQL queries from embedded files
func NewQueryManager() (*QueryManager, error) {
	qm := &QueryManager{
		queries: make(map[string]string),
	}

	// Read all SQL files
	files, err := sqlFiles.ReadDir("sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL directory: %w", err)
	}

	// Process each file
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			data, err := sqlFiles.ReadFile("sql/" + file.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", file.Name(), err)
			}

			// Parse named queries from this file
			fileQueries := parseNamedQueries(string(data))

			// Store each query with namespace from filename
			namespace := strings.TrimSuffix(file.Name(), ".sql")
			for name, query := range fileQueries {
				// Store as namespace.queryname (e.g., "chat.GetAllChatRoom")
				key := fmt.Sprintf("%s.%s", namespace, name)
				qm.queries[key] = query
			}
		}
	}

	return qm, nil
}

// Get returns a specific query by namespace and name
func (qm *QueryManager) Get(namespace, name string) (string, error) {
	key := fmt.Sprintf("%s.%s", namespace, name)
	query, exists := qm.queries[key]
	if !exists {
		return "", fmt.Errorf("query not found: %s", key)
	}
	return query, nil
}

// parseNamedQueries extracts individual queries from SQL content
func parseNamedQueries(content string) map[string]string {
	queries := make(map[string]string)
	lines := strings.Split(content, "\n")

	var currentName string
	var currentQuery strings.Builder

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "-- name:") {
			// Save previous query if we were building one
			if currentName != "" && currentQuery.Len() > 0 {
				queries[currentName] = strings.TrimSpace(currentQuery.String())
				currentQuery.Reset()
			}

			// Extract new query name
			nameParts := strings.SplitN(trimmedLine, "-- name:", 2)
			if len(nameParts) == 2 {
				currentName = strings.TrimSpace(nameParts[1])
			}
		} else if currentName != "" {
			// Add line to current query
			currentQuery.WriteString(line)
			currentQuery.WriteString("\n")
		}
	}

	// Save the last query
	if currentName != "" && currentQuery.Len() > 0 {
		queries[currentName] = strings.TrimSpace(currentQuery.String())
	}

	return queries
}
