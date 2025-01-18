package requests

import "github.com/jmoiron/sqlx"

func Get(db *sqlx.DB, requestID int) (*RequestDTO, error) {
	var request RequestDTO

	query := `
		SELECT *
		FROM requests r
		WHERE r.id = $1;`

	err := db.Get(&request, query, requestID)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
