package requests

import "time"

type RequestDTO struct {
	Id               int       `db:"id"`
	Timestamp        time.Time `db:"timestamp"`
	OriginalFilename string    `db:"original_filename"`
	TotalCount       int       `db:"total_count"`
	ProcessedCount   int       `db:"processed_count"`
}
