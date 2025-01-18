package requests

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func (r *RequestDTO) WriteToDB(db *sqlx.DB) (int, error) {
	query := `
		INSERT INTO requests (timestamp, original_filename) 
		VALUES (:timestamp, :original_filename) 
		RETURNING id;
	`

	var id int
	rows, err := db.NamedQuery(query, r)
	if err != nil {
		// Assuming utils.Logger is set up correctly
		utils.Logger.Error("cant run query to create request", zap.Error(err))
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			utils.Logger.Error("cant get last id", zap.Error(err))
			return 0, err
		}
	}

	return id, nil
}
