package errors

import "time"

type ErrorDTO struct {
	Id        int       `db:"id"`
	Timestamp time.Time `db:"timestamp"`
	Email     string    `db:"email"`
	Error     string    `db:"error"`
}
