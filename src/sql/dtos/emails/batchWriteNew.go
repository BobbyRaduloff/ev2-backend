package emails

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func BatchWriteNew(db *sqlx.DB, e []EmailDTO) error {

	query := `
		INSERT INTO emails (
			email,
			timestamp,
			request_id,
			first_name,
			last_name,
			title,
			state,
			city,
			country,
			company_name,
			industry,
			linkedin_link,
			employee_count,
			status
		) VALUES (
			:email,
			:timestamp,
			:request_id,
			:first_name,
			:last_name,
			:title,
			:state,
			:city,
			:country,
			:company_name,
			:industry,
			:linkedin_link,
			:employee_count,
			:status
		);`

	_, err := db.NamedExec(query, e)
	if err != nil {
		utils.Logger.Error("cant exec query for emails", zap.Error(err))
		return err
	}

	return nil
}
