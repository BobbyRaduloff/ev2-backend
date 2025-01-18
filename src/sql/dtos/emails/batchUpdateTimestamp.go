package emails

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func BatchUpdateEmailsTimestamp(db *sqlx.DB, emails []EmailDTO) error {
	ids := make([]interface{}, len(emails))
	for i, email := range emails {
		ids[i] = email.Id
	}

	query, args, err := sqlx.In(`UPDATE emails SET timestamp = NOW() WHERE id IN (?);`, ids)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	_, err = db.Exec(query, args...)
	if err != nil {
		utils.Logger.Error("failed to set timestamps", zap.Error(err))
		return err
	}

	return nil
}
