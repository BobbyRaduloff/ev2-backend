package migrations

import "github.com/jmoiron/sqlx"

var m04schema = `
CREATE TABLE IF NOT EXISTS errors (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP,
	email TEXT NOT NULL,
	error TEXT NOT NULL
);
`

func M04CreateErrorTable(db *sqlx.DB) {
	db.MustExec(m04schema)
}
