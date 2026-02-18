// Package database provides a tool for querying PostgreSQL databases
package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// QueryTool allows executing SQL queries against the database
type QueryTool struct {
	db *sql.DB
}

// NewQueryTool creates a new database query tool
func NewQueryTool(connStr string) (*QueryTool, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &QueryTool{db: db}, nil
}

// Name implements tools.Tool
func (t *QueryTool) Name() string {
	return "query_database"
}

// Description implements tools.Tool
func (t *QueryTool) Description() string {
	return "Execute SQL queries against the PostgreSQL database. Only SELECT queries are allowed for safety."
}

// Schema implements tools.Tool
func (t *QueryTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "The SQL query to execute (SELECT only)",
			},
		},
		"required": []string{"query"},
	}
}

// Execute implements tools.Tool
func (t *QueryTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Get query from params
	queryStr, ok := params["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query parameter must be a string")
	}

	// Safety check: only allow SELECT queries
	// In a production system, use a proper SQL parser
	if len(queryStr) < 6 || queryStr[:6] != "SELECT" && queryStr[:6] != "select" {
		return nil, fmt.Errorf("only SELECT queries are allowed for safety")
	}

	// Execute query
	rows, err := t.db.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Collect results
	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{} to represent each column
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into our values
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			// Handle []byte as string for readability
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return map[string]interface{}{
		"columns": columns,
		"rows":    results,
		"count":   len(results),
	}, nil
}

// Close closes the database connection
func (t *QueryTool) Close() error {
	return t.db.Close()
}
