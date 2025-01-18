package sql

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToPostgres(host string, port int, user string, password string, dbname string, ssl bool) *sqlx.DB {
	// Connection string for PostgreSQL
	sslmode := "disable"
	if ssl {
		sslmode = "enable"
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	// Connect to the PostgreSQL database
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("can't connect to postgres db: %v\n", err)
	}

	return db
}
