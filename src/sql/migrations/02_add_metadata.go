package migrations

import (
	"log"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql"
	"github.com/jmoiron/sqlx"
)

func M02AddMetadata(db *sqlx.DB) {
	err := sql.AddColumnIfNotExists(db, "emails", "linkedin_link", "TEXT NOT NULL DEFAULT ''")
	if err != nil {
		log.Panic(err)
	}

	err = sql.AddColumnIfNotExists(db, "emails", "employee_count", "INTEGER NOT NULL DEFAULT 0")
	if err != nil {
		log.Panic(err)
	}
}
