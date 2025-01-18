package errors

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func BatchWriteNew(db *sqlx.DB, e []ErrorDTO) error {
	query := `
	INSERT INTO errors (
		email,
		timestamp,
		error
	) VALUES (
		:email,
		:timestamp, 
		:error
	)
	`

	_, err := db.NamedExec(query, e)
	if err != nil {
		utils.Logger.Error("cant exec query for errors", zap.Error(err))
		return err
	}

	return nil
}
