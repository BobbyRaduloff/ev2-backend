package requests

import "github.com/jmoiron/sqlx"

func UpdateCountsById(db *sqlx.DB, requestId int) error {
	updateQuery := `
		UPDATE requests
		SET total_count = (
				SELECT COUNT(*)
				FROM emails e
				WHERE e.request_id = $1
			),
			processed_count = (
				SELECT COUNT(*)
				FROM emails e
				WHERE e.request_id = $1 AND e.status = 'DONE'
			)
		WHERE id = $1;
	`
	_, err := db.Exec(updateQuery, requestId)
	return err
}
