package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// AddColumnIfNotExists adds a column to a table if it does not already exist.
func AddColumnIfNotExists(db *sqlx.DB, tableName, columnName, columnType string) error {
	// Query to check if the column exists
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2;
	`
	var existingColumns []string

	err := db.Select(&existingColumns, query, tableName, columnName)
	if err != nil {
		return fmt.Errorf("failed to get table info: %w", err)
	}

	// Check if the column already exists
	for _, col := range existingColumns {
		if col == columnName {
			return nil
		}
	}

	// Add the column if it doesn't exist
	addColumnQuery := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnType)
	_, err = db.Exec(addColumnQuery)
	if err != nil {
		return fmt.Errorf("failed to add column: %w", err)
	}

	return nil
}
