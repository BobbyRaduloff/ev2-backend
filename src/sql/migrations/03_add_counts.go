package migrations

import (
	"log"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql"
	"github.com/jmoiron/sqlx"
)

func M03AddCounts(db *sqlx.DB) {
	err := sql.AddColumnIfNotExists(db, "requests", "total_count", "INTEGER NOT NULL DEFAULT 0")
	if err != nil {
		log.Panic(err)
	}

	err = sql.AddColumnIfNotExists(db, "requests", "processed_count", "INTEGER NOT NULL DEFAULT 0")
	if err != nil {
		log.Panic(err)
	}

	updateCountsQuery := `
		UPDATE requests
		SET total_count = sub.count,
			processed_count = sub.processed
		FROM (
			SELECT r.id,
				(SELECT COUNT(*) FROM emails e WHERE e.request_id = r.id) AS count,
				(SELECT COUNT(*) FROM emails e WHERE e.request_id = r.id AND e.status = 'DONE') AS processed
			FROM requests r
		) AS sub
		WHERE requests.id = sub.id;`

	_, err = db.Exec(updateCountsQuery)
	if err != nil {
		log.Panic(err)
	}
}
