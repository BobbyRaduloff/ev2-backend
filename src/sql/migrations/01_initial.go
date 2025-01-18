package migrations

import "github.com/jmoiron/sqlx"

var m01schema = `
CREATE TABLE IF NOT EXISTS requests (
	id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP,
	original_filename TEXT NOT NULL DEFAULT 'file.csv'
);
CREATE TABLE IF NOT EXISTS emails (
	id SERIAL PRIMARY KEY,
	email TEXT NOT NULL,
	is_valid INTEGER NOT NULL DEFAULT 0,
	is_nonpersonal INTEGER NOT NULL DEFAULT 0,
	is_disposable INTEGER NOT NULL DEFAULT 0,
	has_mx INTEGER NOT NULL DEFAULT 0,
	has_spf INTEGER NOT NULL DEFAULT 0,
	has_dmarc INTEGER NOT NULL DEFAULT 0,
	has_dkim INTEGER NOT NULL DEFAULT 0,
	mx TEXT NOT NULL DEFAULT '',
	handshake INTEGER NOT NULL DEFAULT 0,
	handshake_name TEXT NOT NULL DEFAULT 'FAILED',
	timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
	first_name TEXT NOT NULL DEFAULT '',
	last_name  TEXT NOT NULL DEFAULT '',
	title TEXT NOT NULL DEFAULT '',
	state TEXT NOT NULL DEFAULT '',
	city TEXT NOT NULL DEFAULT '',
	country TEXT NOT NULL DEFAULT '',
	company_name TEXT NOT NULL DEFAULT '',
	industry TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'QUEUED',
	request_id INT REFERENCES requests(id),
	UNIQUE (email, request_id)
);
CREATE TABLE IF NOT EXISTS filters (
	id SERIAL PRIMARY KEY,
  has_mx INTEGER,
	has_spf INTEGER,
	has_dkim INTEGER,
	has_dmarc INTEGER,
	isnt_disposable INTEGER,
	isnt_nonpersonal INTEGER,
	allow_failed INTEGER,
	allow_connected_but_requires_auth INTEGER,
	allow_connected_but_requires_tls INTEGER,
	allow_connected_but_catchall INTEGER,
	allow_connected_but_not_exists INTEGER,
	allow_connected_but_antispam INTEGER,
	allow_connected_and_exists INTEGER
);
`

func M01Initial(db *sqlx.DB) {
	db.MustExec(m01schema)

	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM filters")
	if err != nil {
		panic(err)
	}

	if count == 0 {
		insertQuery := `
		INSERT INTO filters (
			has_mx, has_spf, has_dkim, has_dmarc, isnt_disposable, isnt_nonpersonal, 
			allow_failed, allow_connected_but_requires_auth, allow_connected_but_requires_tls, 
			allow_connected_but_catchall, allow_connected_but_not_exists, 
			allow_connected_but_antispam, allow_connected_and_exists
		) VALUES (
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
		);`
		_, err := db.Exec(insertQuery)
		if err != nil {
			panic(err)
		}
	}
}
